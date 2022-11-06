package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// ConfigDataJSON struct
type ConfigDataJSON struct {
	Service string          `json:"service"`
	Data    json.RawMessage `json:"data"`
}

// ConfigClient struct
type ConfigClient struct {
	uri      string
	service  string
	version  int
	cfg      []byte
	client   *http.Client
	callback UpdateCallback
	refresh  *time.Ticker
	done     chan bool
}

// UpdateCallback type
type UpdateCallback func([]byte)

const EMPTY_STRING = ""

var ErrEmptyServiceName = errors.New("empty service name")

// Connect function
func Connect(uri string, service string, version ...int) (*ConfigClient, error) {
	cl := &http.Client{
		Timeout: 5 * time.Second,
	}

	if service == EMPTY_STRING {
		return nil, ErrEmptyServiceName
	}

	ver := 0
	if len(version) > 0 {
		ver = version[0]
	}

	return &ConfigClient{
		uri:     uri,
		service: service,
		version: ver,
		client:  cl,
		done:    make(chan bool),
	}, nil
}

// SetServiceParams function
func (c *ConfigClient) SetServiceParams(service string, version ...int) error {
	if service == EMPTY_STRING {
		return ErrEmptyServiceName
	}

	c.service = service

	if len(version) > 0 && version[0] > 0 {
		c.version = version[0]
	}

	return nil
}

// CreateConfig function
func (c *ConfigClient) CreateConfig(ctx context.Context, data interface{}) error {
	return c.doPostOrPutRequest(ctx, http.MethodPost, data)
}

// readConfig function
func (c *ConfigClient) readConfig(ctx context.Context) ([]byte, error) {
	req, err := c.makeGetOrDeleteRequest(ctx, http.MethodGet)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request aborted with status: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// ReadConfigBytes function
func (c *ConfigClient) ReadConfigBytes(ctx context.Context) ([]byte, error) {
	var err error
	c.cfg, err = c.readConfig(ctx)

	return c.cfg, err
}

// ReadAndDecodeConfig function
func (c *ConfigClient) ReadAndDecodeConfig(ctx context.Context, data interface{}) error {
	cfgBytes, err := c.ReadConfigBytes(ctx)
	if err != nil {
		return err
	}

	return json.Unmarshal(cfgBytes, data)
}

// UpdateConfig function
func (c *ConfigClient) UpdateConfig(ctx context.Context, data interface{}) error {
	return c.doPostOrPutRequest(ctx, http.MethodPut, data)
}

// DeleteConfig function
func (c *ConfigClient) DeleteConfig(ctx context.Context) error {
	req, err := c.makeGetOrDeleteRequest(ctx, http.MethodDelete)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("request aborted with status: %s", resp.Status)
}

// AssignRefreshCallback function
func (c *ConfigClient) AssignRefreshCallback(period time.Duration, cb UpdateCallback) error {
	if c.callback != nil {
		return errors.New("callback is already assigned")
	}

	c.callback = cb
	c.refresh = time.NewTicker(period)

	go func() {
		for {
			select {
			case <-c.done:
				c.refresh.Stop()
				return
			case <-c.refresh.C:
				if cfgBytes, err := c.readConfig(context.Background()); err == nil {
					if !bytes.Equal(cfgBytes, c.cfg) {
						c.cfg = cfgBytes
						c.callback(c.cfg)
					}
				} else {
					log.Println(err)
				}
			}
		}
	}()

	return nil
}

func (c *ConfigClient) formatPostData(service string, data interface{}) ([]byte, error) {
	if service == "" {
		return nil, errors.New("empty service name")
	}

	if data == nil {
		return nil, errors.New("empty config data")
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	cfgData := &ConfigDataJSON{
		Service: service,
		Data:    dataBytes,
	}

	return json.Marshal(cfgData)
}

// doPostOrPutRequest function
func (c *ConfigClient) doPostOrPutRequest(ctx context.Context, method string, data interface{}) error {
	cfgData, err := c.formatPostData(c.service, data)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, c.uri, bytes.NewReader(cfgData))
	if err != nil {
		return err
	}
	defer req.Body.Close()

	req.Header.Add("Content-Type", "application/json")

	if _, err = c.client.Do(req); err != nil {
		return err
	}

	return nil
}

// makeGetOrDeleteRequest function
func (c *ConfigClient) makeGetOrDeleteRequest(ctx context.Context, method string) (*http.Request, error) {
	serviceUri := fmt.Sprintf("%s?service=%s&version=%d", c.uri, c.service, c.version)
	return http.NewRequestWithContext(ctx, method, serviceUri, nil)
}

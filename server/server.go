package server

import (
	"context"
	"errors"
	"fmt"
	"go-cloud-camp/internal/config"
	"go-cloud-camp/internal/handlers"
	"go-cloud-camp/internal/logging"
	"go-cloud-camp/internal/storage"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julienschmidt/httprouter"
)

// ConfigServer struct
type ConfigServer struct {
	cfg      *config.Config
	log      *logging.Logger
	storage  *storage.AppStorage
	router   *httprouter.Router
	listener net.Listener
	server   *http.Server
}

// Create function
func Create(path string) (*ConfigServer, error) {
	srv := &ConfigServer{}

	var err error

	// Create server config
	if srv.cfg, err = config.GetConfig(path); err != nil {
		return nil, err
	}

	// Create server logger
	if srv.log, err = logging.GetLogger(srv.cfg.Logging); err != nil {
		return nil, err
	}

	// Create storage
	if srv.storage, err = storage.Create(&srv.cfg.Storage, srv.log); err != nil {
		return nil, err
	}

	srv.log.Debug("create application router")
	srv.router = httprouter.New()

	srv.log.Debug("register router handlers")
	handlers.Create(srv.log, srv.storage).Register(srv.router)

	srv.log.Debug("create http server")
	srv.server = &http.Server{
		Handler:      srv.router,
		ReadTimeout:  srv.cfg.Listen.ReadTimeout,
		WriteTimeout: srv.cfg.Listen.WriteTimeout,
	}

	srv.log.Debug("create net listener")
	if srv.listener, err = net.Listen("tcp", srv.listenAddr()); err != nil {
		return nil, err
	}

	return srv, nil
}

// Run function
func (s *ConfigServer) Run() {
	s.log.Infof("start listening on %s", s.listenAddr())

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer close(stopCh)

	go s.startServer(stopCh)

	stop := <-stopCh

	fmt.Println()
	s.log.Info("interrupted with signal:", stop)

	s.stopServer()
}

// startserver function
func (s *ConfigServer) startServer(stopCh chan<- os.Signal) {
	defer s.listener.Close()
	defer s.storage.Close()

	err := s.server.Serve(s.listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Fatalln(err)
	}
	s.log.Debug("server stopped gracefully")
}

// stopserver function
func (s *ConfigServer) stopServer() {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Listen.ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Fatalln(err)
	}
	s.log.Debug("application shutdown")
}

// listenAddr function
func (s *ConfigServer) listenAddr() string {
	return fmt.Sprintf("%s:%s", s.cfg.Listen.BindIp, s.cfg.Listen.Port)
}

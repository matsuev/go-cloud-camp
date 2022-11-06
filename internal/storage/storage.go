package storage

import (
	"context"
	"go-cloud-camp/internal/common"
	"go-cloud-camp/internal/config"
	"go-cloud-camp/internal/logging"
	"go-cloud-camp/internal/storage/mongodb"
	"log"
)

// StorageBackend interface
type StorageBackend interface {
	CreateConfig(data *common.RequestData) error
	ReadConfig(string, int) ([]byte, error)
	UpdateConfig(data *common.RequestData) error
	DeleteConfig(string, int) error
	Close(context.Context) error
}

const (
	BACKEND_MONGODB = "mongodb"
)

// AppStorage struct
type AppStorage struct {
	logger  *logging.Logger
	backend StorageBackend
}

// Create function
func Create(cfg *config.StorageParams, log *logging.Logger) (*AppStorage, error) {
	var err error
	var backend StorageBackend

	switch cfg.Backend {
	case BACKEND_MONGODB:
		backend, err = mongodb.Create(cfg, log)
	default:
		backend = nil
	}

	if err != nil {
		return nil, err
	}

	return &AppStorage{
		logger:  log,
		backend: backend,
	}, nil
}

// Close function
func (s *AppStorage) Close() {
	if err := s.backend.Close(context.Background()); err != nil {
		log.Fatalln(err)
	}
}

// Create function
func (s *AppStorage) Create(data *common.RequestData) error {
	return s.backend.CreateConfig(data)
}

// Read function
func (s *AppStorage) Read(service string, version int) ([]byte, error) {
	return s.backend.ReadConfig(service, version)
}

// Update function
func (s *AppStorage) Update(data *common.RequestData) error {
	return s.backend.UpdateConfig(data)
}

// Delete function
func (s *AppStorage) Delete(service string, version int) error {
	return s.backend.DeleteConfig(service, version)
}

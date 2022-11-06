package common

import "errors"

// ErrAlreadyCreated
var ErrNotFound = errors.New("not found")
var ErrAlreadyCreated = errors.New("config already created")
var ErrEmptyServiceName = errors.New("empty service name")
var ErrNotValidJsonData = errors.New("not valid json data")
var ErrServiceNotFound = errors.New("service not found")
var ErrConfigIsUsed = errors.New("config is used")

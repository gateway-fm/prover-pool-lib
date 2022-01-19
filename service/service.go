package service

import (
	"crypto/sha256"
	"encoding/hex"
)

type IService interface {
	// HealthCheck check service health by
	// sending status request
	HealthCheck() error

	// Status return service current status
	Status() Status

	// ID return service unique ID
	ID() string

	// Address return service address
	Address() string
}

// TODO split address field to host and port

// BaseService represent basic service
// model implementation
type BaseService struct {
	id      string // service unique id - sha256(address)
	status  Status // service current status
	address string // service address to connect
}

// NewService create new BaseService with address and discovery
func NewService(address, transportProtocol string) IService {
	return &BaseService{
		id:      generateServiceID(address),
		status:  StatusUnHealthy,
		address: transportProtocol + address,
	}
}

// HealthCheck check service health by
// sending status request
func (n *BaseService) HealthCheck() error {
	// TODO implement basic http or tcp healthchecks
	return nil
}

// Status return BaseService current status
func (n *BaseService) Status() Status {
	return n.status
}

// ID return service unique ID
func (n *BaseService) ID() string {
	return n.id
}

// Address return service address
func (n *BaseService) Address() string {
	return n.address
}

// generateServiceID create BaseService unique id by
// hashing given address string
func generateServiceID(addr string) string {
	h := sha256.New()
	h.Write([]byte(addr))
	sum := h.Sum(nil)

	return hex.EncodeToString(sum)
}

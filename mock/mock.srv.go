package mock

import (
	"log"

	context "golang.org/x/net/context"
)

// Endpoint service that implements MockServer interface
type Endpoint struct{}

// Base is a handler for MockServer interface
func (e *Endpoint) Base(_ context.Context, bm *BaseMsg) (*EmptyMsg, error) {
	log.Println(bm)
	return &EmptyMsg{}, nil
}

// Echo is an echo handler for MockServer interface
func (e *Endpoint) Echo(_ context.Context, em *EchoMsg) (*EchoMsg, error) {
	return em, nil
}

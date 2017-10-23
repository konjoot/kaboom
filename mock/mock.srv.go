package mock

import context "golang.org/x/net/context"

// Endpoint service that implements MockServer interface
type Endpoint struct{}

// Base is a handler for MockServer interface
func (e *Endpoint) Base(context.Context, *BaseMsg) (*EmptyMsg, error) {
	return &EmptyMsg{}, nil
}

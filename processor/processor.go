package processor

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/golang/protobuf/proto"

	"google.golang.org/grpc"
)

// Process reads bytes from io.Reader and processes it as an unary gRPC request
// to the service, which available on the address (addr),
// to the specified method (method), then it writes bytes of the response
// into the io.Writer
func Process(in io.Reader, addr, method string, out io.Writer) error {

	opts := []grpc.DialOption{grpc.WithInsecure()}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return err
	}

	req := &Request{Payload: in}
	resp := &Response{Payload: out}

	return grpc.Invoke(context.Background(), method, req, resp, conn)
}

// Request struct is used to send raw bytes to the gRPC endpoints
// Implements proto.Marshaller interface
// Implements proto.Message interface
type Request struct {
	Payload io.Reader
}

// Marshal is a method of the proto.Marshal interface
func (in *Request) Marshal() ([]byte, error) {
	return ioutil.ReadAll(in.Payload)
}

// Reset is a method of the proto.Message interface
func (in *Request) Reset() { *in = Request{} }

// String is a method of the proto.Message interface
func (in *Request) String() string { return proto.CompactTextString(in) }

// ProtoMessage is a method of the proto.Message interface
func (*Request) ProtoMessage() {}

// Response struct is used to receive raw bytes from the gRPC endpoints
// Implements proto.Unmarshaler interface
// Implements proto.Message interface
type Response struct {
	Payload io.Writer
}

// Unmarshal is a method for proto.Unmarshaler interface
func (out *Response) Unmarshal(buf []byte) error {
	_, err := fmt.Fprint(out.Payload, buf)
	return err
}

// Reset is a method of the proto.Message interface
func (out *Response) Reset() { *out = Response{} }

// String is a method of the proto.Message interface
func (out *Response) String() string { return proto.CompactTextString(out) }

// ProtoMessage is a method of the proto.Message interface
func (*Response) ProtoMessage() {}

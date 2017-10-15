package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

var (
	addr string
)

func main() {

	opts := []grpc.DialOption{grpc.WithInsecure()}
	addr = "localhost:50051"

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	cli := &client{conn}

	out, err := cli.Call()
	if err != nil {
		log.Println(err)
		log.Printf("Kaboom:(\n => %v\n", out)
		os.Exit(1)
	}

	log.Printf("Kaboom:)\n => %v\n", out)
}

type client struct {
	conn *grpc.ClientConn
}

type RequestMessage struct{}

func (in *RequestMessage) Marshal() ([]byte, error) {
	// 0000 1000||0000 0001 -> 0x08, 0x01
	// 0001 0010 0000 0111 -> 0x12 0x07
	// 0x12 0x07||0x75 0x73 0x65 0x72 0x5F 0x69 0x64
	// 0001 1010 0100 1000 -> 0x1a, 0x08
	// 0x1a, 0x08||0x73 0x63, 0x6F, 0x70, 0x65, 0x5F, 0x69, 0x64
	return []byte{
		0x08, 0x01, // 1|0: 1
		0x12, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5F, 0x69, 0x64, // 2|2: user_id
		0x1a, 0x08, 0x73, 0x63, 0x6F, 0x70, 0x65, 0x5F, 0x69, 0x64, // 3|2: scope_id
	}, nil
}

func (in *RequestMessage) Reset()         { *in = RequestMessage{} }
func (in *RequestMessage) String() string { return proto.CompactTextString(in) }
func (*RequestMessage) ProtoMessage()     {}

type ResponseMessage struct {
	Response map[ProtoField]interface{}
}

func (out *ResponseMessage) Reset()         { *out = ResponseMessage{} }
func (out *ResponseMessage) String() string { return proto.CompactTextString(out) }
func (*ResponseMessage) ProtoMessage()      {}
func (out *ResponseMessage) MarshalText() (text []byte, err error) {
	for k, v := range out.Response {
		text = append(text, byte(k.Number))
		text = append(text, byte(k.Type))
		text = append(text, []byte(fmt.Sprint(v))...)
	}
	return text, nil
}

func (out *ResponseMessage) Unmarshal(buf []byte) error {
	if out.Response == nil {
		out.Response = make(map[ProtoField]interface{})
	}

	out.Response[ProtoField{Number: uint8(1), Type: uint8(2)}] = interface{}(1)
	return nil
}

type ProtoField struct {
	Number uint8
	Type   uint8
}

func (cli *client) Call() (*ResponseMessage, error) {
	ctx := context.Background()

	in := &RequestMessage{}

	out := &ResponseMessage{}
	opts := make([]grpc.CallOption, 0)
	err := grpc.Invoke(ctx, "/logos.sirius.rooms.Rooms/AddUser", in, out, cli.conn, opts...)
	if err != nil {
		return out, err
	}
	return out, nil
}

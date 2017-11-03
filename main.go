package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/konjoot/kaboom/config"
	"google.golang.org/grpc"
)

var (
	addr string
)

func main() {

	conf := config.New()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	addr = conf.Listen

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

func (cli *client) Call() (*ResponseMessage, error) {
	ctx := context.Background()

	// dec := json.NewDecoder(os.Stdin)
	var m = make(map[string]interface{})
	// if dec.More() { // read only first json
	// 	err := dec.Decode(&m)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	log.Printf("%v: %v\n", m, m)
	// }

	in := &RequestMessage{payload: m}

	out := &ResponseMessage{}
	opts := make([]grpc.CallOption, 0)
	err := grpc.Invoke(ctx, "/mock.Mock/Base", in, out, cli.conn, opts...)
	log.Println("err => ", err)
	if err != nil {
		return out, err
	}
	return out, nil
}

type RequestMessage struct {
	payload map[string]interface{}
}

func (in *RequestMessage) Marshal() ([]byte, error) {
	for key, val := range in.payload {
		log.Println(key)
		log.Println(val)
	}

	return []byte{
		0x28, 0x03,
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

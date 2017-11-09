package processor_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"github.com/konjoot/kaboom/config"
	"github.com/konjoot/kaboom/mock"
	"github.com/konjoot/kaboom/processor"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	conf := config.New()

	lis, err := net.Listen("tcp", conf.Listen)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	mock.RegisterMockServer(s, &mock.Endpoint{})

	started := make(chan struct{})
	go func() {
		defer s.GracefulStop()
		go close(started)

		if err := s.Serve(lis); err != nil {
			fmt.Println(err)
		}
	}()
	<-started

	os.Exit(m.Run())
}

func TestProcess(t *testing.T) {
	var (
		err  error
		bts  []byte
		conf = config.New()
	)
	for _, tc := range []struct {
		name   string
		in     io.Reader
		addr   string
		method string
		out    io.ReadWriter
		expBts []byte
		expErr error
	}{
		{
			name:   "SuccessAllFields",
			in:     bytes.NewReader([]byte{0x08, 0x96, 0x01, 0x12, 0x06, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67}),
			addr:   conf.Listen,
			method: "/mock.Mock/Echo",
			out:    &bytes.Buffer{},
			expBts: []byte{0x08, 0x96, 0x01, 0x12, 0x06, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67},
		},
		{
			name:   "SuccessFirstField",
			in:     bytes.NewReader([]byte{0x08, 0x96, 0x01}),
			addr:   conf.Listen,
			method: "/mock.Mock/Echo",
			out:    &bytes.Buffer{},
			expBts: []byte{0x08, 0x96, 0x01},
		},
		{
			name:   "SuccessLastField",
			in:     bytes.NewReader([]byte{0x12, 0x06, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67}),
			addr:   conf.Listen,
			method: "/mock.Mock/Echo",
			out:    &bytes.Buffer{},
			expBts: []byte{0x12, 0x06, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67},
		},
		{
			name:   "RandomInput",
			in:     bytes.NewReader([]byte{0x13, 0x16}),
			addr:   conf.Listen,
			method: "/mock.Mock/Echo",
			out:    &bytes.Buffer{},
			expBts: []byte{},
			expErr: errors.New("rpc error: code = Internal desc = grpc: error unmarshalling request: unexpected EOF"),
		},
		{
			name:   "BadEndpoint",
			in:     bytes.NewReader([]byte{0x13, 0x16}),
			addr:   "localhost:3000",
			method: "/mock.Mock/Echo",
			out:    &bytes.Buffer{},
			expBts: []byte{},
			expErr: errors.New("rpc error: code = Unavailable desc = grpc: the connection is unavailable"),
		},
		{
			name:   "WrongMethodName",
			in:     bytes.NewReader([]byte{0x13, 0x16}),
			addr:   conf.Listen,
			method: "/mock.Mock/WrongMethodName",
			out:    &bytes.Buffer{},
			expBts: []byte{},
			expErr: errors.New("rpc error: code = Unimplemented desc = unknown method WrongMethodName"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err = processor.Process(tc.in, tc.addr, tc.method, tc.out)
			t.Log("err => ", err)
			if fmt.Sprint(err) != fmt.Sprint(tc.expErr) {
				t.Error("Expected => ", tc.expErr)
			}

			bts, err = ioutil.ReadAll(tc.out)
			t.Logf("bts => % #x", bts)
			if string(bts) != string(tc.expBts) {
				t.Errorf("Expected => % #x", tc.expBts)
			}
		})
	}
}

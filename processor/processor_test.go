package processor_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"github.com/konjoot/kaboom/processor"

	"github.com/konjoot/kaboom/config"
	"github.com/konjoot/kaboom/mock"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	conf := config.New()

	lis, err := net.Listen("tcp", conf.Listen)
	if err != nil {
		fmt.Println(err)
		return
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	mock.RegisterMockServer(s, &mock.Endpoint{})

	started := make(chan struct{})
	go func() {
		defer s.GracefulStop()
		close(started)
		err := s.Serve(lis)
		if err != nil {
			fmt.Println(err)
		}
	}()
	<-started

	os.Exit(m.Run())
}

func TestProcess(t *testing.T) {

	var err error
	var bts []byte

	for _, tc := range []struct {
		name   string
		in     io.Reader
		addr   string
		method string
		out    io.ReadWriter
		expBts []byte
		expErr error
	}{} {
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

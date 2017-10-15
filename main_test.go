package main

import (
	"testing"
)

func TestRequestMessageMarshal(t *testing.T) {
	bts, err := (&RequestMessage{}).Marshal()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bts))
}

func TestResponseMessageUnmarshal(t *testing.T) {
	out := &ResponseMessage{}
	if err := out.Unmarshal([]byte{}); err != nil {
		t.Error(err)
	}

	t.Log(out.Response)
}

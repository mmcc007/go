package main

import (
	"context"
	"testing"

	pb "github.com/mmccabe/go/examples/helloworld/helloworld"
)

func TestHello(t *testing.T) {
	s := server{}

	// set up test cases
	helloTests := []struct {
		name     string
		expected string
	}{
		{
			name:     "world",
			expected: "Hello world",
		},
		{
			name:     "bob",
			expected: "Hello bob",
		},
	}

	for _, tt := range helloTests {
		req := &pb.HelloRequest{Name: tt.name}
		resp, err := s.SayHello(context.Background(), req)
		if err != nil {
			t.Errorf("HelloTest(%v) got unexpected error", err)
		}
		if resp.Message != tt.expected {
			t.Errorf("HelloText(%v)=%v, expected %v", tt.name, resp.Message, tt.expected)
		}
	}
}

package wsrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Args struct {
	A, B int
}

type Arith struct{}

func (a *Arith) Mul(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (a *Arith) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}

func TestWSRPC(t *testing.T) {
	server := NewServer(":8972", "/ws")

	server.Register(&Arith{})

	go server.Serve()
	defer server.Close()

	client, err := NewClient("http://localhost:8972", "ws://localhost:8972/ws")
	assert.NoError(t, err)
	err = client.Dial()
	assert.NoError(t, err)

	var reply int
	err = client.Call("Arith.Mul", &Args{2, 3}, &reply)
	assert.NoError(t, err)
	assert.Equal(t, 6, reply)

	err = client.Call("Arith.Add", &Args{2, 3}, &reply)
	assert.NoError(t, err)
	assert.Equal(t, 5, reply)

	client.Close()

}

func BenchmarkWSRPC(b *testing.B) {
	server := NewServer(":8972", "/ws")

	server.Register(&Arith{})

	go server.Serve()
	defer server.Close()

	client, err := NewClient("http://localhost:8972", "ws://localhost:8972/ws")
	assert.NoError(b, err)
	err = client.Dial()
	assert.NoError(b, err)
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var reply int
		err = client.Call("Arith.Mul", &Args{2, 3}, &reply)
		assert.NoError(b, err)
		assert.Equal(b, 6, reply)
	}

}

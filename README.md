# wsrpc
a tiny websocket rpc framework.

I like it because of its tiny and exquisite.

## Uasage

### define service

```go
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
```

### start a server

```go
server := NewServer(":8972", "/ws")
server.Register(&Arith{})
go server.Serve()
defer server.Close()
...
```

### start a client

```go
client, err := NewClient("http://localhost:8972", "ws://localhost:8972/ws")
if err != {
    panic(err)
}
err = client.Dial()
if err != {
    panic(err)
}

defer client.Close()

var reply int
err = client.Call("Arith.Mul", &Args{2, 3}, &reply)
...
```
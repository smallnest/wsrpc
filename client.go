package wsrpc

import (
	"net/rpc"
	"sync"

	"golang.org/x/net/websocket"
)

// WSRPCClient is a websocket based rpc client.
type WSRPCClient struct {
	origin, wsAddr string

	config *websocket.Config

	sync.RWMutex
	client *rpc.Client

	// plugins
	DialFunc       DialFunc
	BeforeCallFunc BeforeCallFunc
	AfterCallFunc  AfterCallFunc
	BeforeGoFunc   BeforeGoFunc
}

// NewClient returns a new WSRPCClient.
func NewClient(origin, wsAddr string) (*WSRPCClient, error) {
	config, err := websocket.NewConfig(wsAddr, origin)
	if err != nil {
		return nil, err
	}

	return &WSRPCClient{
		origin: origin,
		wsAddr: wsAddr,
		config: config,
	}, nil
}

func (c *WSRPCClient) Dial() error {
	// plugin
	if c.DialFunc != nil {
		if err := c.DialFunc(c.config); err != nil {
			return err
		}
	}

	ws, err := websocket.DialConfig(c.config)
	if err != nil {
		return err
	}

	c.client = rpc.NewClient(ws)

	return nil
}

func (c *WSRPCClient) reconnect() error {
	return c.Dial()
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (c *WSRPCClient) Call(serviceMethod string, args interface{}, reply interface{}) error {
	c.RWMutex.RLock()
	client := c.client
	c.RWMutex.RUnlock()

	// plugin
	if c.BeforeCallFunc != nil {
		if err := c.BeforeCallFunc(serviceMethod, args); err != nil {
			return err
		}
	}

	err := client.Call(serviceMethod, args, reply)

	// plugin
	if c.AfterCallFunc != nil {
		if err := c.AfterCallFunc(serviceMethod, args, reply, err); err != nil {
			return err
		}
	}

	if err != nil {
		c.RWMutex.Lock()
		if client == c.client {
			c.client.Close()
			c.reconnect()
		}
		c.RWMutex.Unlock()
	}

	return err
}

// Go invokes the function asynchronously.
func (c *WSRPCClient) Go(serviceMethod string, args interface{}, reply interface{}) *rpc.Call {
	c.RWMutex.RLock()
	client := c.client
	c.RWMutex.RUnlock()

	call := client.Go(serviceMethod, args, reply, nil)
	if call != nil && call.Error != nil {
		c.RWMutex.Lock()
		if client == c.client {
			c.client.Close()
			c.reconnect()
		}
		c.RWMutex.Unlock()
	}
	return call
}

// Close closes the client connection.
func (c *WSRPCClient) Close() error {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()

	return c.client.Close()
}

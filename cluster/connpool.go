package cluster

import (
	"errors"
	"pomelo-go/cluster/clusterpb"
	"sync"
	"sync/atomic"
)

type connPool struct {
	index uint32
	v     []clusterpb.MemberClientAgent
}

type rpcClient struct {
	sync.RWMutex
	isClosed bool
	pools    map[string]*connPool
}

func newConnArray(maxSize uint, addr string) (*connPool, error) {
	a := &connPool{
		index: 0,
		v:     make([]clusterpb.MemberClientAgent, maxSize),
	}
	if err := a.init(addr); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *connPool) init(addr string) error {
	for i := range a.v {

		memberAgent := clusterpb.NewMqttMemberClient(addr)
		err := memberAgent.Connect()
		if err != nil {
			return err
		}

		a.v[i] = memberAgent

	}
	return nil
}

func (a *connPool) Get() clusterpb.MemberClientAgent {
	next := atomic.AddUint32(&a.index, 1) % uint32(len(a.v))
	return a.v[next]
}

func (a *connPool) Close() {
	for i, c := range a.v {
		if c != nil {
			err := c.Close()
			if err != nil {
				// TODO: error handling
			}
			a.v[i] = nil
		}
	}
}

func newRPCClient() *rpcClient {
	return &rpcClient{
		pools: make(map[string]*connPool),
	}
}

func (c *rpcClient) getConnPool(addr string) (*connPool, error) {
	c.RLock()
	if c.isClosed {
		c.RUnlock()
		return nil, errors.New("rpc client is closed")
	}
	array, ok := c.pools[addr]
	c.RUnlock()
	if !ok {
		var err error
		array, err = c.createConnPool(addr)
		if err != nil {
			return nil, err
		}
	}
	return array, nil
}

func (c *rpcClient) createConnPool(addr string) (*connPool, error) {
	c.Lock()
	defer c.Unlock()
	array, ok := c.pools[addr]
	if !ok {
		var err error
		// TODO: make conn count configurable
		array, err = newConnArray(1, addr)
		if err != nil {
			return nil, err
		}
		c.pools[addr] = array
	}
	return array, nil
}

func (c *rpcClient) closePool() {
	c.Lock()
	if !c.isClosed {
		c.isClosed = true
		// close all connections
		for _, array := range c.pools {
			array.Close()
		}
	}
	c.Unlock()
}

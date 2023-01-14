package infra

import (
	"errors"
	"sync"
	"time"
)

type NodeShard interface {
	Validate() error
	Close()
}

type Dial func(host string) (NodeShard, error)

type Pool struct {
	services  []string
	ClientMap sync.Map
	Dial      Dial
	ch        chan interface{}
	Timeout   time.Duration
}

type Setting struct {
	Hosts []string //多ip:port
}

func NewDefaultPool(hosts []string, dial Dial) *Pool {
	pool := Pool{
		services: hosts,
		Dial:     dial,
		Timeout:  60 * time.Second,
		ch:       make(chan interface{}),
	}
	return &pool
}

func NewPool(setting *Setting, dial Dial) *Pool {
	pool := Pool{
		services: setting.Hosts,
		Dial:     dial,
		Timeout:  60 * time.Second,
		ch:       make(chan interface{}),
	}
	return &pool
}

func (p *Pool) Run() error {
	if p.Dial == nil {
		return errors.New("dial is nil")
	}

	p.checkAndUpdateMap()

	go p.checkService()

	if _, err := p.Get(); err != nil {
		return errors.New("detect all node failed, no one can be use")
	}

	return nil
}

func (p *Pool) SetTimeout(timeout time.Duration) {
	p.Timeout = timeout
}

func (p *Pool) Get() (NodeShard, error) {
	var client interface{}
	i := 0

Retry:
	p.ClientMap.Range(func(k, v interface{}) bool {
		client = v
		return false
	})
	// 获取到节点，返回
	if node, ok := client.(NodeShard); ok {
		return node, nil
	}

	// 获取失败，触发探测，重试一次
	if i < 1 {
		i++
		p.checkAndUpdateMap()
		goto Retry
	}

	return nil, errors.New("node all dead")
}

func (p *Pool) getServices() []string {
	return p.services
}

func (p *Pool) checkAndUpdateMap() {
	// 1.移除变更节点

	var err error

	hosts := p.getServices()

	// 2.检查hosts 是否变更， 删除不再新列表中的host

	for _, host := range p.services {
		if !stringInSlice(host, hosts) {
			client, ok := p.ClientMap.Load(host)
			if ok {
				if node, isOk := client.(NodeShard); isOk {
					node.Close()
				}
			}
			p.ClientMap.Delete(host)
		}
	}

	p.services = hosts

	// 3.判断所有已有节点是否可用，不可用则删除
	p.ClientMap.Range(func(host, client interface{}) bool {
		node, ok := client.(NodeShard)
		if ok {
			err = node.Validate()
		} else {
			err = errors.New("Type Assertion err")
		}

		if err != nil {
			client, ok := p.ClientMap.Load(host)
			if ok {
				if node, isOk := client.(NodeShard); isOk {
					node.Close()
				}
			}
			p.ClientMap.Delete(host) //remove bad node
		}
		return true
	})

	// 4.判断节点列表中不可用节点是否可用，可用加回map
	for _, host := range p.services {
		_, ok := p.ClientMap.Load(host)
		if !ok {
			conn, err := p.Dial(host)
			if err == nil {
				p.ClientMap.Store(host, conn)
			}
		}
	}

}

// checkService range services and validate, if err will remove from the map

func (p *Pool) checkService() {
	t := time.NewTimer(p.Timeout)
	for {
		select {
		case <-t.C:
			p.checkAndUpdateMap()
		case <-p.ch:
			p.checkAndUpdateMap()
		}
		t.Reset(p.Timeout)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

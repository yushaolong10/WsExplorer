package pool

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

//object pool
//use for connections storage

var (
	errPoolClosed = errors.New("pool closed")
)

type PoolObject interface {
	Close() error
}

type NewObject func() (PoolObject, error)

type Pool interface {
	Acquire() (PoolObject, error) // get object
	Release(PoolObject) error     // release object
	Shutdown() error              // close
}

type poolItem struct {
	object   PoolObject
	updateAt time.Time
}

type normalizePool struct {
	mutex       sync.Mutex
	pool        chan *poolItem
	closed      bool
	minOpen     int
	maxOpen     int
	numOpen     int
	factory     NewObject
	maxLifetime time.Duration
	timeout     time.Duration
}

func NewNormalizePool(minOpen, maxOpen int, maxLifetime, timeout time.Duration, factory NewObject) (Pool, error) {
	if maxOpen <= 0 || minOpen > maxOpen {
		return nil, fmt.Errorf("invalid config: min:%d,max:%d", minOpen, maxOpen)
	}
	p := &normalizePool{
		maxOpen:     maxOpen,
		minOpen:     minOpen,
		maxLifetime: maxLifetime,
		timeout:     timeout,
		factory:     factory,
		pool:        make(chan *poolItem, maxOpen),
	}
	for i := 0; i < minOpen; i++ {
		object, err := factory()
		if err != nil {
			return nil, err
		}
		p.pool <- &poolItem{object: object, updateAt: time.Now()}
		p.numOpen++
	}
	return p, nil
}

func (p *normalizePool) Acquire() (PoolObject, error) {
	if p.closed {
		return nil, errPoolClosed
	}
	select {
	case item := <-p.pool:
		p.numOpen--
		if !p.isExpired(item) {
			return item.object, nil
		}
	default:
	}
	return p.newObject()
}

func (p *normalizePool) newObject() (PoolObject, error) {
	p.mutex.Lock()
	if p.numOpen >= p.maxOpen {
		tick := time.NewTimer(p.timeout)
		select {
		case item := <-p.pool:
			p.numOpen--
			p.mutex.Unlock()
			return item.object, nil
		case <-tick.C:
			p.mutex.Unlock()
			tick.Stop()
			return nil, errors.New("time out because of exceed max num")
		}
	}
	//create new
	object, err := p.factory()
	if err != nil {
		p.mutex.Unlock()
		return nil, err
	}
	p.numOpen++
	p.mutex.Unlock()
	return object, nil
}

func (p *normalizePool) isExpired(item *poolItem) bool {
	if time.Now().Sub(item.updateAt) > p.maxLifetime {
		return true
	}
	return false
}

func (p *normalizePool) Release(object PoolObject) error {
	if p.closed {
		return errPoolClosed
	}
	p.mutex.Lock()
	tick := time.NewTimer(p.timeout)
	select {
	case p.pool <- &poolItem{object: object, updateAt: time.Now()}:
		p.numOpen++
		p.mutex.Unlock()
		return nil
	case <-tick.C:
		p.mutex.Unlock()
		tick.Stop()
		return errors.New("release time out")
	}
}

func (p *normalizePool) Shutdown() error {
	p.mutex.Lock()
	if p.closed {
		p.mutex.Unlock()
		return errPoolClosed
	}
	p.closed = true
	p.mutex.Unlock()
	//range
	close(p.pool)
	for item := range p.pool {
		item.object.Close()
		p.numOpen--
	}
	return nil
}

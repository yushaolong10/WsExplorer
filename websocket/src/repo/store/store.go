package store

import (
	"lib/logger"
	"lib/pool"
	"repo/store/client"
	"time"
)

var innerStorePool pool.Pool

func Init(addr []string, minOpen, maxOpen int, maxLifeTime, timeout int) error {
	newFunc := func() (object pool.ObjectItem, e error) {
		conn, err := client.NewStoreClient(addr[0])
		return conn, err
	}
	maxLifeTimeSec := time.Second * time.Duration(maxLifeTime)
	timeoutMs := time.Millisecond * time.Duration(timeout)
	p, err := pool.NewNormalizePool(minOpen, maxOpen, maxLifeTimeSec, timeoutMs, newFunc)
	if err != nil {
		return err
	}
	innerStorePool = p
	return nil
}

func Get(key string) (string, error) {
	object, err := innerStorePool.Acquire()
	if err != nil {
		logger.Error("[Get] acquire pool error. err:%s", err.Error())
		return "", err
	}
	conn := object.(*client.StoreClient)
	str, _, err := conn.Get(key)
	if err != nil {
		logger.Error("[Get] conn get error. err:%s", err.Error())
		return "", err
	}
	if err := innerStorePool.Release(conn); err != nil {
		logger.Error("[Get] release pool error. err:%s", err.Error())
	}
	return str, nil
}

func Delete(key string) error {
	object, err := innerStorePool.Acquire()
	if err != nil {
		logger.Error("[Delete] acquire pool error. err:%s", err.Error())
		return err
	}
	conn := object.(*client.StoreClient)
	_, err = conn.Delete(key)
	if err != nil {
		logger.Error("[Delete] conn delete error. err:%s", err.Error())
		return err
	}
	if err := innerStorePool.Release(conn); err != nil {
		logger.Error("[Delete] release pool error. err:%s", err.Error())
	}
	return nil
}

func Set(key, val string) error {
	object, err := innerStorePool.Acquire()
	if err != nil {
		logger.Error("[Set] acquire pool error. err:%s", err.Error())
		return err
	}
	conn := object.(*client.StoreClient)
	_, err = conn.Set(key, val)
	if err != nil {
		logger.Error("[Set] conn set error. err:%s", err.Error())
		return err
	}
	if err := innerStorePool.Release(conn); err != nil {
		logger.Error("[Set] release pool error. err:%s", err.Error())
	}
	return nil
}

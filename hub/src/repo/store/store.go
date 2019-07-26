package store

import (
	"lib/logger"
	"lib/pool"
	"repo/store/client"
	"time"
)

var innerStorePool pool.Pool

func Init(addr []string) error {
	newFunc := func() (object pool.PoolObject, e error) {
		conn, err := client.NewStoreClient(addr[0])
		return conn, err
	}
	p, err := pool.NewNormalizePool(20, 200, time.Second*10, time.Second*3, newFunc)
	if err != nil {
		return err
	}
	innerStorePool = p
	return nil
}

func Get(key string) (string, bool) {
	object, err := innerStorePool.Acquire()
	if err != nil {
		logger.Error("[Get] acquire pool error. err:%s", err.Error())
		return "", false
	}
	conn := object.(*client.StoreClient)
	str, has, _ := conn.Get(key)
	return str, has
}

func Delete(key string) bool {
	object, err := innerStorePool.Acquire()
	if err != nil {
		logger.Error("[Delete] acquire pool error. err:%s", err.Error())
		return false
	}
	conn := object.(*client.StoreClient)
	ret, _ := conn.Delete(key)
	return ret
}

func Set(key, val string) bool {
	object, err := innerStorePool.Acquire()
	if err != nil {
		logger.Error("[Set] acquire pool error. err:%s", err.Error())
		return false
	}
	conn := object.(*client.StoreClient)
	ret, _ := conn.Set(key, val)
	return ret
}

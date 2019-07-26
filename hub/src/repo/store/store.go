package store

import (
	"repo/store/client"
	"sync/atomic"
)


var innerStore *Store

func Init(addr []string) error {
	innerStore = &Store{
		host:addr,
		hostCount:uint32(len(addr)),
	}
	return nil
}

type Store struct {
	host []string
	hostCount uint32
	connCount uint32
}

func GetStore(key string) (string, bool) {
	add := atomic.AddUint32(&innerStore.connCount, 1)
	addr := innerStore.host[add % innerStore.connCount]
	storeClient, _ := getConnFromPool(addr)
	str,has,_ :=  storeClient.Get(key)
	return str, has
}

func getConnFromPool(addr string) (*client.StoreClient,error) {
	conn,_ := client.NewStoreClient(addr)
	return conn, nil
}

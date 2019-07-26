package client

import (
	"testing"
)

func TestStoreClient_SetGetDel(t *testing.T) {
	client, _ := NewStoreClient("10.209.214.21:8112")
	ret, _ := client.Set("abc", "21")
	t.Log("client set (abc:21) ", ret)

	str, ret, _ := client.Get("abc")
	t.Log("client get (abc) ", str, ret)

	ret, _ = client.Set("abcd", "29")
	t.Log("client set (abcd:29) ", ret)

	str, ret, _ = client.Get("ddd")
	t.Log("client get (ddd) ", str, ret)

	ret, _ = client.Delete("abc")
	t.Log("client delete (abc) ", ret)

	str, ret, _ = client.Get("abc")
	t.Log("client get (abc) ", str, ret)
}

/**
goos: linux
goarch: amd64
core: 4cpu
mem: 8G
BenchmarkStoreClient_PoolSet-4             10000            143254 ns/op
BenchmarkStoreClient_NewSet-4              10000            274057 ns/op
BenchmarkStoreClient_PoolGet-4             10000            159634 ns/op
BenchmarkStoreClient_NewGet-4              10000            283312 ns/op
*/

func BenchmarkStoreClient_PoolSet(b *testing.B) {
	b.StopTimer()
	client, _ := NewStoreClient("10.209.214.21:8112")
	b.StartTimer()
	for i := 1; i < b.N; i++ {
		client.Set("abc1234567890", "21")
	}
}

func BenchmarkStoreClient_NewSet(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 1; i < b.N; i++ {
		client, _ := NewStoreClient("10.209.214.21:8112")
		client.Set("abc1234567890", "21")
	}
}

func BenchmarkStoreClient_PoolGet(b *testing.B) {
	b.StopTimer()
	client, _ := NewStoreClient("10.209.214.21:8112")
	b.StartTimer()
	for i := 1; i < b.N; i++ {
		client.Get("abc1234567890")
	}
}

func BenchmarkStoreClient_NewGet(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 1; i < b.N; i++ {
		client, _ := NewStoreClient("10.209.214.21:8112")
		client.Get("abc1234567890")
	}
}

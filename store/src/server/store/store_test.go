package store

import (
	"testing"
)

func TestStoreNode_SetAndGet(t *testing.T) {
	//set
	var engine, _ = NewMapEngine(32, 32, 1)
	engine.Set("aa1", "111")
	engine.Set("aa2", "222")
	engine.Set("aa3", "333")
	//get
	v1, _ := engine.Get("aa1")
	v2, _ := engine.Get("aa1")
	v3, _ := engine.Get("aa3")
	v4, b := engine.Get("aa4")
	if !b {
		t.Logf("v4 not exist")
	}
	t.Logf("v1:%s,v2:%s,v3:%s,v4:%s", v1, v2, v3, v4)
}

func TestComputeAsciiSum(t *testing.T) {
	str := "ab"
	ret1 := computeAsciiSum(str)
	t.Logf("ret1:%v", ret1)
}

/*
goos: linux
goarch: amd64
BenchmarkMapEngine_SetNoDepth-4         20000000                63.5 ns/op
BenchmarkMapEngine_SetDepth1-4          20000000                84.3 ns/op
BenchmarkMapEngine_SetDepth2-4          20000000               102 ns/op
*/

func BenchmarkMapEngine_SetNoDepth(b *testing.B) {
	//set
	b.StopTimer()
	var engine, _ = NewMapEngine(32, 32, 0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		engine.Set("abcdefgh", "1234567")
	}
}

func BenchmarkMapEngine_SetDepth1(b *testing.B) {
	//set
	b.StopTimer()
	var engine, _ = NewMapEngine(32, 32, 1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		engine.Set("abcdefgh", "1234567")
	}
}

func BenchmarkMapEngine_SetDepth2(b *testing.B) {
	//set
	b.StopTimer()
	var engine, _ = NewMapEngine(32, 32, 2)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		engine.Set("abcdefgh", "1234567")
	}
}

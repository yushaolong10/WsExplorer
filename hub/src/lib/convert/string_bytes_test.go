package convert

import (
	"fmt"
	"testing"
)

func TestString2Bytes(t *testing.T) {
	str := "abc123"
	bytes := String2Bytes(str)
	fmt.Println("byte:", bytes)
}

func TestBytes2String(t *testing.T) {
	bytes := []byte{97, 98, 99, 49, 50, 51}
	str := Bytes2String(bytes)
	fmt.Println("str:", str)
}

/**
BenchmarkBytesString-4          1000000000               2.53 ns/op
BenchmarkBytesStringRaw-4       100000000               13.2 ns/op
*/

func BenchmarkBytesString(b *testing.B) {
	b.StopTimer()
	str := "abcdefghijklmnopqrstuvwxyz"

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a := String2Bytes(str)
		Bytes2String(a)
	}
}

func BenchmarkBytesStringRaw(b *testing.B) {
	b.StopTimer()
	str := "abcdefghijklmnopqrstuvwxyz"

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b := []byte(str)
		_ = string(b)
	}
}

package main

func main() {
	client, _ := NewStoreClient("localhost:8112")
	client.Set("abc", "21")
	client.Get("abc")
	client.Set("abce", "29")
	client.Get("ddd")
	client.Delete("abc")
	client.Get("abc")
}

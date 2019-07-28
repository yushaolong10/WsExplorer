package store

import "fmt"

func SetUniqIdGrpcHost(uniqId int, host string) error {
	key := fmt.Sprintf("%d%s", uniqId, "_host_grpc")
	return Set(key, host)
}

func GetUniqIdGrpcHost(uniqId int64) (string, error) {
	key := fmt.Sprintf("%d%s", uniqId, "_host_grpc")
	return Get(key)
}

func DeleteUniqIdGrpcHost(uniqId int) error {
	key := fmt.Sprintf("%d%s", uniqId, "_host_grpc")
	return Delete(key)
}

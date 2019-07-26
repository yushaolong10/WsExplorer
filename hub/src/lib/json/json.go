package json

import "github.com/json-iterator/go"

var jsonParser = jsoniter.ConfigCompatibleWithStandardLibrary

func Unmarshal(data []byte, v interface{}) error {
	return jsonParser.Unmarshal(data, v)
}

func Marshal(v interface{}) ([]byte, error) {
	return jsonParser.Marshal(v)
}

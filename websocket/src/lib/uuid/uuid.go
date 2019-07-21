package uuid

import "github.com/satori/go.uuid"

func MakeUUId() string {
	id, _ := uuid.NewV4()
	return id.String()
}

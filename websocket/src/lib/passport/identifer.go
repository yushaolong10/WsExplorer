package passport

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ActorAdmin   = 1
	ActorTeacher = 2
	ActorStudent = 3
	ActorParent  = 4
)

type Identifier struct {
	UniqId int
	Actor  int
}

func CheckAndGetInfoByToken(token string) (*Identifier, error) {
	slice := strings.Split(token, ":")
	if len(slice) != 3 {
		return nil, fmt.Errorf("token split length not equal 3")
	}
	uniqId, _ := strconv.Atoi(slice[1])
	actor, _ := strconv.Atoi(slice[2])
	i := Identifier{
		UniqId: uniqId,
		Actor:  actor,
	}
	return &i, nil
}

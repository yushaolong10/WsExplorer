package store

import (
	"fmt"
	"lib/convert"
	"lib/json"
	"strings"
)

type Engine interface {
	Get(key string) (string, bool)
	Set(key, val string) bool
	Delete(key string) bool
}

var storeEngine Engine

func Init(keyMaxLen, valMaxLen int, maxSplitDepth int) error {
	engine, err := NewMapEngine(keyMaxLen, valMaxLen, maxSplitDepth)
	if err != nil {
		return err
	}
	storeEngine = engine
	return nil
}

func Execute(cmdStr []byte) ([]byte, error) {
	if storeEngine == nil {
		return nil, fmt.Errorf("store engine not init")
	}
	cmd, err := parseCommand(convert.Bytes2String(cmdStr))
	//return
	var protocol *storeProtocol
	if err != nil {
		protocol = &storeProtocol{Err: 1, Msg: err.Error()}
	} else {
		protocol = cmd.Do(storeEngine)
	}
	return json.Marshal(protocol)
}

func parseCommand(cmd string) (Command, error) {
	if strings.HasPrefix(cmd, "SET") {
		return decodeSet(cmd)
	} else if strings.HasPrefix(cmd, "GET") {
		return decodeGet(cmd)
	} else if strings.HasPrefix(cmd, "DELETE") {
		return decodeDel(cmd)
	}
	return &emptyCmd{}, nil
}

func decodeSet(cmd string) (*setCmd, error) {
	cmdList := strings.Split(cmd, " ")
	if len(cmdList) < 3 {
		return nil, fmt.Errorf("[SET] command length less 3")
	}
	return &setCmd{Key: cmdList[1], Val: cmdList[2]}, nil
}

func decodeGet(cmd string) (*getCmd, error) {
	cmdList := strings.Split(cmd, " ")
	if len(cmdList) < 2 {
		return nil, fmt.Errorf("[GET] command length less 2")
	}
	return &getCmd{Key: cmdList[1]}, nil
}

func decodeDel(cmd string) (*delCmd, error) {
	cmdList := strings.Split(cmd, " ")
	if len(cmdList) < 2 {
		return nil, fmt.Errorf("[DELETE] command length less 2")
	}
	return &delCmd{Key: cmdList[1]}, nil
}

func computeAsciiSum(input string) int {
	length := len(input)
	count := 0
	for i := 0; i < length; i++ {
		count += int(input[i])
	}
	return count
}

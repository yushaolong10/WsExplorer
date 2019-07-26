package store

import (
	"fmt"
	"lib/logger"
)

const (
	initNodeCount = 32
)

//tree lock
func NewMapEngine(keyMaxLen, valMaxLen int, maxSplitDepth int) (Engine, error) {
	if keyMaxLen < 1 || keyMaxLen > 32 {
		return nil, fmt.Errorf("key_max_len error. must in range [1,32]")
	}
	if valMaxLen < 1 || valMaxLen > 32 {
		return nil, fmt.Errorf("val_max_len error. must in range [1,32]")
	}
	engine := mapEngine{
		keyMaxLen: keyMaxLen,
		valMaxLen: valMaxLen,
		nodeCount: initNodeCount,
		nodes:     make(map[int]*storeNode),
	}
	//init node
	for i := 0; i < initNodeCount; i++ {
		engine.nodes[i] = initStoreNode(i, i, 0, maxSplitDepth)
	}
	return &engine, nil
}

type mapEngine struct {
	keyMaxLen int
	valMaxLen int
	nodeCount int
	nodes     map[int]*storeNode
}

func (m *mapEngine) checkInputValid(k, v string) bool {
	return m.keyMaxLen > len(k) && m.valMaxLen > len(v)
}

func (m *mapEngine) checkKeyValid(k string) bool {
	return m.keyMaxLen > len(k)
}

func (m *mapEngine) GetKeyIndex(k string) int {
	count := computeAsciiSum(k)
	return count % m.nodeCount
}

func (m *mapEngine) GetNode(index int) (*storeNode, error) {
	if n, ok := m.nodes[index]; ok {
		return n, nil
	}
	return nil, fmt.Errorf("not found node by index:%d", index)
}

func (m *mapEngine) Set(k, v string) bool {
	if !m.checkInputValid(k, v) {
		logger.Error("[Set] key_len or val_len invalid. key:%s, val:%s", k, v)
		return false
	}
	index := m.GetKeyIndex(k)
	node, err := m.GetNode(index)
	if err != nil {
		logger.Error("[Set] found node error. err:%s", err.Error())
		return false
	}
	return node.Set(k, v)
}

func (m *mapEngine) Get(k string) (string, bool) {
	if !m.checkKeyValid(k) {
		logger.Error("[Get] key_len invalid. key:%s", k)
		return "", false
	}
	index := m.GetKeyIndex(k)
	node, err := m.GetNode(index)
	if err != nil {
		logger.Error("[Get] found node error. err:%s", err.Error())
		return "", false
	}
	return node.Get(k)
}

func (m *mapEngine) Delete(k string) bool {
	if !m.checkKeyValid(k) {
		logger.Error("[Delete] key_len invalid. key:%s", k)
		return false
	}
	index := m.GetKeyIndex(k)
	node, err := m.GetNode(index)
	if err != nil {
		logger.Error("[Delete] found node error. err:%s", err.Error())
		return false
	}
	return node.Delete(k)
}

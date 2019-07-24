package store

import "sync"

func initStoreNode(groupId int, midVal int, curDepth int, maxSplitDepth int) *storeNode {
	if midVal == 0 {
		midVal = 1
	}
	node := &storeNode{
		mux:           new(sync.Mutex),
		groupId:       groupId,
		midVal:        midVal,
		boundaryVal:   midVal << 1,
		maxSplitDepth: maxSplitDepth,
		curDepth:      curDepth,
		isEdge:        true,
		mapStore:      make(map[string]string),
	}
	node.Split()
	return node
}

//simple tree storage
type storeNode struct {
	mux     *sync.Mutex
	groupId int
	//tree node mod middle value
	midVal        int
	boundaryVal   int
	maxSplitDepth int
	//current tree node depth
	curDepth   int
	leftChild  *storeNode
	rightChild *storeNode
	mapStore   map[string]string
	isEdge     bool
}

func (m *storeNode) modDivideVal(sum int) int {
	return sum % m.boundaryVal
}

func (m *storeNode) Split() {
	if m.curDepth >= m.maxSplitDepth {
		return
	}
	m.leftChild = initStoreNode(m.groupId, m.boundaryVal, m.curDepth+1, m.maxSplitDepth)
	m.rightChild = initStoreNode(m.groupId, m.boundaryVal, m.curDepth+1, m.maxSplitDepth)
	m.isEdge = false
	return
}

func (m *storeNode) Set(k, v string) bool {
	if m.isEdge {
		m.mux.Lock()
		m.mapStore[k] = v
		m.mux.Unlock()
		return true
	}
	sum := computeAsciiSum(k)
	if m.modDivideVal(sum) > m.midVal {
		return m.rightChild.Set(k, v)
	}
	return m.leftChild.Set(k, v)
}

func (m *storeNode) Get(k string) (string, bool) {
	if m.isEdge {
		m.mux.Lock()
		v, ok := m.mapStore[k]
		m.mux.Unlock()
		return v, ok
	}
	sum := computeAsciiSum(k)
	if m.modDivideVal(sum) > m.midVal {
		return m.rightChild.Get(k)
	}
	return m.leftChild.Get(k)
}

func (m *storeNode) Delete(k string) bool {
	if m.isEdge {
		m.mux.Lock()
		_, ok := m.mapStore[k]
		delete(m.mapStore, k)
		m.mux.Unlock()
		return ok
	}
	sum := computeAsciiSum(k)
	if m.modDivideVal(sum) > m.midVal {
		return m.rightChild.Delete(k)
	}
	return m.leftChild.Delete(k)
}

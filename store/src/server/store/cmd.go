package store

type storeProtocol struct {
	Err  int                    `json:"err"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

type Command interface {
	Do(engine Engine) *storeProtocol
}

//default
type emptyCmd struct{}

func (empty *emptyCmd) Do(engine Engine) *storeProtocol {
	return &storeProtocol{Err: 1, Msg: "not support cmd"}
}

//set command
type setCmd struct {
	Key, Val string
}

func (s *setCmd) Do(engine Engine) *storeProtocol {
	ok := engine.Set(s.Key, s.Val)
	return &storeProtocol{Err: 0, Data: map[string]interface{}{"set": ok}}
}

//get command
type getCmd struct {
	Key string
}

func (g *getCmd) Do(engine Engine) *storeProtocol {
	val, ok := engine.Get(g.Key)
	return &storeProtocol{Err: 0, Data: map[string]interface{}{"value": val, "get": ok}}
}

//delete command
type delCmd struct {
	Key string
}

func (d *delCmd) Do(engine Engine) *storeProtocol {
	ok := engine.Delete(d.Key)
	return &storeProtocol{Err: 0, Data: map[string]interface{}{"delete": ok}}
}

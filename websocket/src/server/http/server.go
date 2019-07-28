package http

import (
	"config"
	"fmt"
	"github.com/gorilla/websocket"
	"lib/logger"
	"lib/passport"
	"lib/utils"
	"net/http"
	"repo/store"
	"runtime/debug"
	"server/connection"
	"strings"
	"time"
)

var grpcInternalAddr string

func localInit() {
	interIP, err := utils.GetInternalIP()
	if err != nil {
		panic(fmt.Sprintf("get internal ip err:%v", err))
	}
	grpcHost := strings.Split(config.Global.Grpc.Addr, ":")
	grpcInternalAddr = fmt.Sprintf("%s:%s", interIP, grpcHost[1])
}

func Start(errChan chan<- error) {
	localInit()
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/ws/conn", wsHandler)
	go func() {
		fmt.Println("http server start...")
		logger.Info("ListenAndServe addr: %s", config.Global.Http.Addr)
		errChan <- http.ListenAndServe(config.Global.Http.Addr, httpMux)
	}()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("[wsHandler] #PANIC# error. err:%v, stack:%s", err, string(debug.Stack()))
		}
	}()
	//check
	token := r.FormValue("token")
	if token == "" {
		errNeedToken(w)
		return
	}
	//token = "wsconnsecret:27:1"
	identifier, err := passport.CheckAndGetInfoByToken(token)
	if err != nil {
		errPassport(w)
		logger.Error("[wsHandler] token check error. token:%s, err:%s", token, err.Error())
		return
	}
	//upgrade
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errUpgrade(w)
		logger.Error("[wsHandler] ws conn upgrade error. token:%s, uniqId:%d, err:%s", token, identifier.UniqId, err.Error())
		return
	}
	//conn pool
	wsconn := &connection.WsConnInfo{
		UniqId:    identifier.UniqId,
		Actor:     identifier.Actor,
		WsConn:    conn,
		Timestamp: time.Now().Unix(),
	}
	if err := connection.AddWsConn2Pool(wsconn); err != nil {
		errAddPool(w)
		logger.Error("[wsHandler] ws conn add pool error. token:%s, uniqId:%d, err:%s", token, wsconn.UniqId, err.Error())
		return
	}
	if err := store.SetUniqIdGrpcHost(wsconn.UniqId, grpcInternalAddr); err != nil {
		errStoreSet(w)
		logger.Error("[wsHandler] ws conn set store error. token:%s, uniqId:%d, err:%s", token, wsconn.UniqId, err.Error())
		return
	}
	if err := connection.Monitor(wsconn); err != nil {
		errMonitor(w)
		logger.Error("[wsHandler] ws conn monitor error. token:%s, uniqId:%d, err:%s", token, wsconn.UniqId, err.Error())
		return
	}
	logger.Info("[wsHandler] user login success. token:%s, uniqId:%d, actor:%d", token, wsconn.UniqId, identifier.Actor)
}

package http

import (
	"lib/json"
	"lib/logger"
	"net/http"
)

func errNeedToken(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	resp := DefError(100001, "need token.")
	bytes, _ := json.Marshal(resp)
	_, err := w.Write(bytes)
	if err != nil {
		logger.Error("[errNeedToken] need token send err:%s", err.Error())
	}
}

func errPassport(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	resp := DefError(100002, "token invalid.")
	bytes, _ := json.Marshal(resp)
	_, err := w.Write(bytes)
	if err != nil {
		logger.Error("[errPassport] token invalid send err:%s", err.Error())
	}
}

func errUpgrade(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	resp := DefError(100003, "ws conn upgrade err.")
	bytes, _ := json.Marshal(resp)
	_, err := w.Write(bytes)
	if err != nil {
		logger.Error("[errUpgrade] ws conn upgrade send err:%s", err.Error())
	}
}

func errAddPool(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	resp := DefError(100004, "ws conn add pool err.")
	bytes, _ := json.Marshal(resp)
	_, err := w.Write(bytes)
	if err != nil {
		logger.Error("[errAddPool] ws conn add pool err:%s", err.Error())
	}
}

func errMonitor(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	resp := DefError(100005, "ws conn monitor err.")
	bytes, _ := json.Marshal(resp)
	_, err := w.Write(bytes)
	if err != nil {
		logger.Error("[errMonitor] ws conn monitor err:%s", err.Error())
	}
}

func errStoreSet(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	resp := DefError(100006, "ws store set err.")
	bytes, _ := json.Marshal(resp)
	_, err := w.Write(bytes)
	if err != nil {
		logger.Error("[errStoreSet] ws store set err:%s", err.Error())
	}
}

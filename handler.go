package rtsp

import (
	"encoding/json"
	"fmt"
	. "github.com/Monibuca/engine/v2"
	"net/http"
	"strings"
	"time"
)

type Resp struct {
	ErrorCode   int         `json:"ErrorCode"`
	Message     string      `json:"Message"`
	Data        interface{} `json:"Data"`
	RefreshTime int64       `json:"RefreshTime"`
}

func makeResp(errCode int, msg string, data interface{}) []byte {
	resp, _ := json.Marshal(Resp{
		ErrorCode:   errCode,
		Message:     msg,
		Data:        data,
		RefreshTime: time.Now().Unix(),
	})
	return resp
}

func makeJsonStrResp(errCode int, msg string, data string) []byte {
	resp := fmt.Sprintf(`{
    "ErrorCode": %d,
    "Message": "%s",
    "Data": "%s",
    "RefreshTime": %d
}`, errCode, msg, data, time.Now().Unix())
	return []byte(resp)
}

func ListAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	isAll := r.URL.Query().Get("isAll")

	var info []interface{}
	collection.Range(func(key, value interface{}) bool {
		rtsp := value.(*RTSP)
		pinfo := &rtsp.RTSPInfo
		if isAll == "1" {
			info = append(info, pinfo)
		} else {
			info = append(info, pinfo.URL)
		}
		return true
	})
	w.Write(makeResp(0, "ok", info))
}

func Pull(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	targetURL := r.URL.Query().Get("target")
	streamPath := r.URL.Query().Get("streamPath")
	if err := new(RTSP).PullStream(streamPath, targetURL); err == nil {
		w.Write(makeResp(0, "", streamPath))
	} else {
		var errCode = -1
		errMsg := err.Error()
		if strings.Contains(errMsg, "badname") {
			errCode = 1
		} else if strings.Contains(errMsg, "timeout") {
			errCode = 2
		}
		w.Write(makeResp(errCode, errMsg, nil))
	}
}

func Stop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if streamPath := r.URL.Query().Get("stream"); streamPath != "" {
		if s := FindStream(streamPath); s != nil {
			s.Cancel()
			w.Write(makeResp(0, "success", nil))
		} else {
			w.Write(makeResp(1, "no such stream", nil))
		}
	} else {
		w.Write(makeResp(-1, "param error", nil))
	}
}

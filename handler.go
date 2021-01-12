package rtsp

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		w.Write(makeResp(0, "", nil))
	} else {
		w.Write(makeResp(-1, err.Error(), nil))
	}
}

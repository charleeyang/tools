package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// 统一响应格式: { "code": 0, "data": {...}, "message": "ok" }
type apiResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

func writeJSON(w http.ResponseWriter, httpStatus int, body apiResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(body)
}

func ok(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, apiResponse{Code: 0, Data: data, Message: "ok"})
}

func okMsg(w http.ResponseWriter, data interface{}, msg string) {
	writeJSON(w, http.StatusOK, apiResponse{Code: 0, Data: data, Message: msg})
}

func fail(w http.ResponseWriter, httpStatus, code int, msg string) {
	writeJSON(w, httpStatus, apiResponse{Code: code, Message: msg})
}

// list 分页响应
type listData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"pageSize"`
}

func decodeBody(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func queryInt(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func pathInt(r *http.Request, key string) int64 {
	v := r.PathValue(key)
	n, _ := strconv.ParseInt(v, 10, 64)
	return n
}

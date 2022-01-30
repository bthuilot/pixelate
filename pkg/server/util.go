package server

import "net/http"

func TextResponse(w http.ResponseWriter, body string, statusCode int) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(body))
}
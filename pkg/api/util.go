package api

import (
	"encoding/json"
	"net/http"
)

func TextResponse(w http.ResponseWriter, body string, statusCode int) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(body))
}

func JsonResponse(w http.ResponseWriter, obj interface{}, statusCode int) {
	jsonText, err := json.Marshal(obj)
	if err != nil {
		TextResponse(w, "api error", http.StatusInternalServerError)
	} else {
		w.WriteHeader(statusCode)
		_, _ = w.Write(jsonText)
	}
}

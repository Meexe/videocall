package utils

import (
	"encoding/json"
	"net/http"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func WsMessage(t string, payload interface{}) map[string]interface{} {
	return map[string]interface{}{"type": t, "payload": payload}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")             // Delete after configuring reverse proxy
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // And this
	json.NewEncoder(w).Encode(data)
}

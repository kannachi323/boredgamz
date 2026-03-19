package utils

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //remove this in production
	},
}

func UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return Upgrader.Upgrade(w, r, nil)
}


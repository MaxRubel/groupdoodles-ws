package main

import (
	"fmt"
	"log"
	"net/http"

	wsHandler "github.com/MaxRubel/groupdoodles-ws/pkg"
)

func main() {
	http.HandleFunc("/ws", wsHandler.HandleWebSocket)
	fmt.Println("Go Websocket server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

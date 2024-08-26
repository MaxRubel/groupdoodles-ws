package structs

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	RoomId  string `json:"roomId"`
	Clients map[string]Client
	mu      sync.RWMutex
}

func (r *Room) AddClient(clientId string, ws *websocket.Conn) error {
	if r == nil {
		return errors.New("nil pointer, client is not defined or accesible")
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.Clients[clientId] = Client{
		ClientId: clientId,
		Conn:     ws,
	}
	
	fmt.Printf("added client %s to room %s \n", clientId, r.RoomId)
	return nil
}

func (r *Room) RemoveClient(clientId string) error {
	if r == nil {
		return errors.New("nil pointer, client is not defined or accesible")
	}
	
	r.mu.Lock() 
	defer r.mu.Unlock()
	
	delete(r.Clients, clientId)

	fmt.Printf("removed client %s from room %s \n", clientId, r.RoomId)
	return nil
}

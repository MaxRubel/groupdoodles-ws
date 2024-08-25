package structs

import (
	"errors"

	"github.com/gorilla/websocket"
)

type Room struct {
	RoomId  string `json:"roomId"`
	Clients map[string]Client
}

func (r *Room) AddClient(clientId string, ws *websocket.Conn) error {
	if r == nil {
		return errors.New("nil pointer, client is not defined or accesible")
	}
	r.Clients[clientId] = Client{
		ClientId: clientId,
		Conn:     ws,
	}
	return nil
}

func (r *Room) RemoveClient(clientId string) error {
	if r == nil {
		return errors.New("nil pointer, client is not defined or accesible")
	}
	delete(r.Clients, clientId)
	return nil
}

package negotiations

import (
	"encoding/json"
	"fmt"

	"github.com/MaxRubel/groupdoodles-ws/pkg/structs"
	"github.com/gorilla/websocket"
)

type OutgoingOffer struct{
	Type     string      `json:"type"`
	To       string      `json:"to"`
	From     string      `json:"from"`
	Room     string      `json:"room"`
	Data     interface{} `json:"data"`
}

func HandleOffer(msg structs.IncomingMessage){
	roomId := msg.Room
	senderId := msg.From
	recipient := msg.To
	offer := msg.Data

	room, err := structs.GetRoom(roomId)

	if err != nil {
		fmt.Println(err)
		return
	}
	
	client, ok := room.Clients[recipient]
	if !ok {
		fmt.Println("client not found in room")
		return
	}

	outMsg := OutgoingOffer{
		Type: "offer",
		To:   recipient,
		From: senderId,
		Room: roomId,
		Data: offer,
	}

	jsonMsg, err := json.Marshal(outMsg)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = client.Conn.WriteMessage(websocket.TextMessage, jsonMsg)
	fmt.Println("sending response")
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
}

func HandleAnswer(msg structs.IncomingMessage){
	roomId := msg.Room
	senderId := msg.From
	recipient := msg.To
	answer := msg.Data

	room, err := structs.GetRoom(roomId)

	if err != nil {
		fmt.Println(err)
		return
	}
	
	client, ok := room.Clients[recipient]
	fmt.Println(client)
	if !ok {
		fmt.Println("client not found in room")
		return
	}

	outMsg := OutgoingOffer{
		Type: "answer",
		To:   recipient,
		From: senderId,
		Room: roomId,
		Data: answer,
	}

	jsonMsg, err := json.Marshal(outMsg)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = client.Conn.WriteMessage(websocket.TextMessage, jsonMsg)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
}
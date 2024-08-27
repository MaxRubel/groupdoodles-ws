package wsHandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/MaxRubel/groupdoodles-ws/pkg/negotiations"
	"github.com/MaxRubel/groupdoodles-ws/pkg/structs"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Parse initial data from URL
	initialData, err := parseInitialData(r, conn)
	if err != nil {
		log.Println("Error parsing initial data:", err)
		return
	}

	log.Printf("New connection: Host: %v, RoomId: %s, ClientId: %s", initialData.Host, initialData.RoomId, initialData.ClientId)

	room, _ := structs.GetRoom(initialData.RoomId)
	
	structs.SendRoomAsJSON(conn, room)

	cleanup := func() {
		log.Printf("Connection closed for client %s in room %s", initialData.ClientId, initialData.RoomId)
		room, err := structs.GetRoom(initialData.RoomId)
		
		if err != nil {
			fmt.Println(err)
			return
		}
		
		if room != nil {
			err := room.RemoveClient(initialData.ClientId)
			if err != nil {
				log.Printf("Error removing client from room: %v", err)
			}
			if len(room.Clients) == 0 {
				structs.DeleteRoom(initialData.RoomId)
			}
		}
		conn.Close()
	}

	defer cleanup()


	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var message structs.IncomingMessage
		if err := json.Unmarshal(p, &message); err != nil {
			log.Println("Error parsing message:", err)
			continue
		}

		// Handle the message based on its type
		switch message.Type {
		case "offer":
			log.Println("received offer")
			negotiations.HandleOffer(message)
		case "answer":
			log.Println("received answer")
			negotiations.HandleAnswer(message)
		case "iceCandidate":
			fmt.Println("got some ice candidate data", message.Data)
			// Handle action
		default:
			log.Printf("Received unknown message type: %s", message.Type)
		}
	}

}

func parseInitialData(r *http.Request, conn *websocket.Conn) (structs.InitialData, error) {
	var initialData structs.InitialData

	encodedData := r.URL.Query().Get("data")
	if encodedData == "" {
		return initialData, fmt.Errorf("no initial data provided")
	}

	decodedData, err := url.QueryUnescape(encodedData)
	if err != nil {
		return initialData, fmt.Errorf("error decoding data: %v", err)
	}

	err = json.Unmarshal([]byte(decodedData), &initialData)
	if err != nil {
		return initialData, fmt.Errorf("error parsing JSON: %v", err)
	}

	if initialData.Host{
		structs.AddRoom(initialData.RoomId)
		fmt.Println("host has joined")
	} else {
		fmt.Println("gues has joined")
	}
	
	room := structs.AllRooms[initialData.RoomId]
	fmt.Println(room)
	
	err = room.AddClient(initialData.ClientId, conn)

	if err != nil {
		fmt.Println(err)
	}

	return initialData, nil
}
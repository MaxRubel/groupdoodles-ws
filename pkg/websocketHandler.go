package wsHandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

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

type InitialData struct {
	Host     bool   `json:"host"`
	RoomId   string `json:"roomId"`
	ClientId string `json:"clientId"`
}

type RegularMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
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

	// Send a welcome message upon successful connection
	welcomeMsg := []byte("You've successfully connected to the WebSocket server.")
	if err := conn.WriteMessage(websocket.TextMessage, welcomeMsg); err != nil {
		log.Println("Error sending welcome message:", err)
		return
	}

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

	// Handle WebSocket connection
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Parse regular messages
		var message RegularMessage
		if err := json.Unmarshal(p, &message); err != nil {
			log.Println("Error parsing message:", err)
			continue
		}

		// Handle the message based on its type
		switch message.Type {
		case "chat":
			log.Printf("Received chat message: %s", message.Content)
			// Handle chat message
		case "action":
			log.Printf("Received action: %s", message.Content)
			// Handle action
		default:
			log.Printf("Received unknown message type: %s", message.Type)
		}

		// Echo the message back (you can modify this behavior as needed)
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}

}

func parseInitialData(r *http.Request, conn *websocket.Conn) (InitialData, error) {
	var initialData InitialData

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
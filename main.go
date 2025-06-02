package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kenzokravin/tic-tac-toe/rooms"
)

var roomController = rooms.CreateRoomController() //Creating room controller.

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected")

	player := rooms.Player{ID: uuid.New(), Name: "anon_player", Turn: false, Hand: []rooms.Card{}, Conn: conn} //Creating new player variable with ID and default values.

	rooms.JoinRoom(roomController, player) //Adding player to available room  with room controller.

	player.StartWriter() //Start writer for player.

	for { //Reading messages from clients.
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Client disconnected")
			break
		}
		fmt.Println("Message:", string(msg))
		// Echo message back
		conn.WriteMessage(websocket.TextMessage, []byte(string(msg)))
	}
}

func main() {

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

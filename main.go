package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kenzokravin/tic-tac-toe/rooms"
)

var roomController = rooms.CreateRoomController() //Creating room controller.
var pConMap = make(map[*websocket.Conn]uuid.UUID) //Key is player id, value is connection.

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

	player := &rooms.Player{ //Creating new player variable with ID and default values. This is a pointer value.
		ID:        uuid.New(),            //Player ID
		Name:      "anon_player",         //Init Player display name.
		Turn:      false,                 //Setting turn to false.
		Hand:      []*rooms.Card{},       //Init player's hand.
		Conn:      conn,                  //Player's ws connection.
		SendQueue: make(chan string, 16), // Init send queue with buffer of 16 messages.
	}

	pConMap[conn] = player.ID //Inserting into pConMap for retrieval when messaged.

	player.StartWriter() //Start writer for player.

	rooms.JoinRoom(&roomController, player) //Adding player to available room  with room controller.

	fmt.Println("Rooms: ", &roomController.Rooms)

	for { //Reading messages from clients.
		_, msg, err := conn.ReadMessage()
		if err != nil {

			rooms.FindRoomByPlayer(player).RemovePlayerFromRoom(player) //Removing player from room.

			player.Close() //Close player connection.

			fmt.Println("Client disconnected")

			break
		}
		fmt.Println("Message:", string(msg))

		//Where we need to handle player messages.
		//A way to do this:
		//Create a find player in rooms

		var clientMsg rooms.PlayerMessage
		err = json.Unmarshal(msg, &clientMsg)
		if err != nil {
			// handle JSON parse error
		}

		rooms.ManagePlayerMessage(player, &clientMsg)

		// Echo message back
		//conn.WriteMessage(websocket.TextMessage, []byte(string(msg)))
	}
}

func main() {
	rooms.CreateCards()               //Creating cards.
	roomController.StartRoomCleaner() //Starting room cleaner.

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

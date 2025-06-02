package rooms

import (
	"fmt"

	"github.com/google/uuid"
)

type RoomController struct {
	Rooms []*Room
}

var plRoomMap = make(map[uuid.UUID]*Room) //Used for retrieving room by playerID. Global.

func CreateRoomController() RoomController {

	rooms := []*Room{} //creating room list.

	rc := RoomController{ //Creating room controller instance.
		Rooms: rooms,
	}

	return rc
}

func CreateRoom() Room { //Creating the room and gameboard.

	roomid := uuid.New()   //Creating the room id.
	state := "Not Started" // Setting state.
	pop := 0               //Setting population to 0

	gameboard := CreateBoard() //Creating gameboard.

	players := []*Player{}

	crRoom := Room{ID: roomid, State: state, Pop: pop, Full: false, Board: &gameboard, Players: players} //create room instance.

	fmt.Println("New Room Created.")

	return crRoom

}

func JoinRoom(rmControl *RoomController, player *Player) {

	availableRooms := false

	for i := 0; i < len(rmControl.Rooms); i++ {
		room := rmControl.Rooms[i]
		if !room.Full && room.State == "Not Started" {
			if JoinSpecificRoom(room, player) {
				return
			}
		}
	}

	if !availableRooms { //if no available rooms, create new room and join.

		crRoom := CreateRoom()

		rmControl.Rooms = append(rmControl.Rooms, &crRoom)

		JoinSpecificRoom(&crRoom, player)
	}

}

func JoinSpecificRoom(room *Room, player *Player) bool { //Add player to room.

	if room.Full { //If room is full then don't add.

		return false
	}

	room.Players = append(room.Players, player) //Add player to room.
	room.State = "Not Started"
	room.Pop += 1 //Increase room population.

	plRoomMap[player.ID] = room //inserting player id and room id into map.

	fmt.Println("Player joined room:", *room)

	if room.Pop == 2 { //If room has two players already, change status to full.
		room.Full = true
		room.State = "Starting Room"

		//We can start the game here as the room is now full.
		StartRoomGame(room)
	}

	return true

}

func FindRoomByPlayer(player *Player) *Room {

	pRoom := plRoomMap[player.ID] //Fetch room pointer from map.

	return pRoom //return room pointer.

}

func ManagePlayerMessage(player *Player, pMsg *PlayerMessage) { //Manages player actions/messages.

	plRoom := FindRoomByPlayer(player) //Finding player room.

	//Check message type and send to room if required.
	switch action := pMsg.Action; action {
	case "play_card":
		PlayCard(plRoom, player, pMsg)

	}

}

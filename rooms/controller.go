package rooms

import "fmt"

type RoomController struct {
	Rooms []*Room
}

func CreateRoomController() RoomController {

	rooms := []*Room{} //creating room list.

	rc := RoomController{ //Creating room controller instance.
		Rooms: rooms,
	}

	return rc
}

func CreateRoom() Room { //Creating the room and gameboard.

	roomid := "replace_with_UUID" //Creating the room id.
	state := "Not Started"        // Setting state.
	pop := 0                      //Setting population to 0

	gameboard := CreateBoard() //Creating gameboard.

	crRoom := Room{ID: roomid, State: state, Pop: pop, Full: false, Board: gameboard} //create room instance.

	fmt.Println("New Room Created.")

	return crRoom

}

func JoinRoom(rmControl *RoomController, player Player) {

	availableRooms := false

	for _, room := range rmControl.Rooms {
		if !room.Full && room.State == "Not Started" {
			if JoinSpecificRoom(room, player) { //If successful in joining room, exit func.
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

func JoinSpecificRoom(room *Room, player Player) bool { //Add player to room.

	if room.Pop == 2 { //If room has two players already, change status to full.
		room.Full = true
		room.State = "Starting Room"

		//We can start the game here as the room is now full.
		StartRoomGame(*room)
	}

	if room.Full { //If room is full then don't add.

		return false
	}

	room.Players = append(room.Players, player) //Add player to room.
	room.State = "Not Started"
	room.Pop += 1 //Increase room population.

	fmt.Println("Player joined room:", *room)

	return true

}

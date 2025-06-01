package rooms

type RoomController struct {
	Rooms []Room
}

func CreateRoomController() RoomController {
	rooms := []Room{}

	rc := RoomController{
		Rooms: rooms,
	}

	return rc
}

func CreateRoom() Room { //Creating the room and gameboard.
	roomid := "replace_with_UUID"
	state := "Not Started"
	pop := 0

	gameboard := CreateBoard()

	crRoom := Room{ID: roomid, State: state, Pop: pop, Full: false, Board: gameboard}

	return crRoom

}

func JoinRoom(rmControl RoomController, player Player) {

	availableRooms := false

	for i := 0; i < len(rmControl.Rooms); i++ {

		if rmControl.Rooms[i].State == "Not Started" && !rmControl.Rooms[i].Full { //Check if game has not started and room is not full.

			if JoinSpecificRoom(rmControl.Rooms[i], player) { //If successful in joining room, exit func.
				return
			}

		}

	}

	if !availableRooms { //if no available rooms, create new room and join.

		crRoom := CreateRoom()

		rmControl.Rooms = append(rmControl.Rooms, crRoom)

		JoinSpecificRoom(crRoom, player)
	}

}

func JoinSpecificRoom(room Room, player Player) bool { //Add player to room.

	if room.Pop == 2 { //If room has two players already, change status to full.
		room.Full = true
	}

	if room.Full { //If room is full then don't add.

		return false
	}

	room.Players = append(room.Players, player) //Add player to room.
	room.Pop += 1                               //Increase room population.

	return true

}

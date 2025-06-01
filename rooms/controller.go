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

func CreateRoom() Room {
	roomid := "replace_with_UUID"
	state := "Not Started"
	capacity := 0

	crRoom := Room{ID: roomid, State: state, Capacity: capacity, Full: false}

	return crRoom

}

func JoinRoom(rmControl RoomController, player Player) {

	for i := 0; i < len(rmControl.Rooms); i++ {

		if rmControl.Rooms[i].State == "Not Started" && !rmControl.Rooms[i].Full {

			//Add player to room.

		}

	}

}

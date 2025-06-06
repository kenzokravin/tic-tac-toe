package rooms

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RoomController struct {
	Rooms []*Room
	Mu    sync.Mutex
}

var plRoomMap = make(map[uuid.UUID]*Room) //Used for retrieving room by playerID. Global.
var plRoomMapMu sync.RWMutex              //Read-Write Mutex allows multiple readers, one write.

const roomCleanerFreq int = 10 //in minutes.

func CreateRoomController() RoomController {

	rooms := []*Room{} //creating room list.

	rc := RoomController{ //Creating room controller instance.
		Rooms: rooms,
	}

	return rc
}

func (rm *RoomController) CreateRoom() *Room { //Creating the room and gameboard.

	rm.Mu.Lock()         //Locking the thread
	defer rm.Mu.Unlock() //Defering unlock until after new room.

	roomid := uuid.New()   //Creating the room id.
	state := "Not Started" // Setting state.
	pop := 0               //Setting population to 0

	lastActive := time.Now()

	gameboard := CreateBoard() //Creating gameboard.

	players := []*Player{}

	crRoom := &Room{ID: roomid, State: state, Pop: pop, Full: false, Board: &gameboard, Players: players, LastActive: lastActive} //create room instance.

	fmt.Println("New Room Created.")

	rm.Rooms = append(rm.Rooms, crRoom)

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

		// crRoom := CreateRoom()

		// rmControl.Rooms = append(rmControl.Rooms, &crRoom)

		crRoom := rmControl.CreateRoom()

		JoinSpecificRoom(crRoom, player)
	}

}

func JoinSpecificRoom(room *Room, player *Player) bool { //Add player to room.

	room.Mu.Lock()

	plRoomMapMu.Lock()

	if room.Full { //If room is full then don't add.

		return false
	}

	room.Players = append(room.Players, player) //Add player to room.
	room.State = "Not Started"
	room.Pop += 1 //Increase room population.

	plRoomMap[player.ID] = room //inserting player id and room id into map.

	fmt.Println("Player joined room:", room)

	if room.Pop == 2 { //If room has two players already, change status to full.
		room.Full = true
		room.State = "Starting Room"

		//We can start the game here as the room is now full.

	}
	room.Mu.Unlock()
	plRoomMapMu.Unlock()

	if room.Full { //If Full, start game using goroutine.
		go func() {

			time.Sleep(1 * time.Second)
			StartRoomGame(room)

		}()
	}

	return true

}

func FindRoomByPlayer(player *Player) *Room {

	plRoomMapMu.RLock() // read-lock

	pRoom := plRoomMap[player.ID] //Fetch room pointer from map.

	plRoomMapMu.RUnlock() // unlock after read

	return pRoom //return room pointer.

}

func ManagePlayerMessage(player *Player, pMsg *PlayerMessage) { //Manages player actions/messages.

	plRoom := FindRoomByPlayer(player) //Finding player room.

	plRoom.ManagePlActionInRm(player, pMsg)

}

func (rc *RoomController) StartRoomCleaner() { //Function used to remove inactive rooms.
	go func() {

		ticker := time.NewTicker(time.Duration(roomCleanerFreq) * time.Minute) //Every roomCleanerFreq mins, run room cleaner.
		defer ticker.Stop()                                                    //Defer stop until after function.

		for {
			fmt.Println("Room Cleaner Cycle Executing...")
			<-ticker.C // Wait roomCleanerFreq minutes

			rc.Mu.Lock() //Locking mutex.

			for _, room := range rc.Rooms {

				if time.Since(room.LastActive) >= time.Duration(roomCleanerFreq)*time.Minute { //Checks if room has been inactive for more than roomCleanerFreq minutes.

					rc.RemoveRoom(room) // Clean up or remove the room
				}

			}
			rc.Mu.Unlock() //unlocking after room cleaning.
		}

	}()

}

func (rc *RoomController) RemoveRoom(room *Room) { //Function to remove data from slice.

	for i, rm := range rc.Rooms {

		if rm.ID == room.ID { //If Id's match:

			RemovePlayerRoomMapEntries(rm) //Removing from player room map.

			rc.Rooms = append(rc.Rooms[:i], rc.Rooms[i+1:]...) //Creates a new slice using everything before i (:i) and after i+1 (i+1)...
			fmt.Println("Room Removed.")
			break //Stop function after deletion.
		}

	}

}

func RemovePlayerRoomMapEntries(room *Room) { //Remove player room map entries.

	plRoomMapMu.Lock() //Lock mutex.

	for _, pl := range room.Players { //For each player in room.

		delete(plRoomMap, pl.ID) //Remove from player room map.

	}

	plRoomMapMu.Unlock() //Unlock mutex.

}

func (room *Room) RemovePlayerFromRoom(player *Player) {

	room.Mu.Lock() //Lock Mutex

	nPlayers := []*Player{}

	for _, pl := range room.Players { //For each player in room.

		if pl.ID != player.ID { //If ID does not match

			nPlayers = append(nPlayers, pl)
		}

	}

	room.Pop-- //Decrease room population.

	fmt.Println("Player removed:", player)

	room.Players = nPlayers //Update player list.

	//Might need to check room state and start game end process if player count is <=1

	room.Mu.Unlock() // Unlock Mutex

}

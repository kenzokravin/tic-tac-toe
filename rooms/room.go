package rooms

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID         uuid.UUID
	State      string
	Pop        int
	Full       bool
	Board      *Board
	Players    []*Player
	LastActive time.Time
	Mu         sync.Mutex
}

func StartRoomGame(room *Room) {

	room.Mu.Lock()

	room.State = "In Progress" //Setting Game state to playing.

	for i := 0; i < room.Pop; i++ { //Drawing Start Cards for players.
		DrawStartCards(room.Players[i])

		fmt.Println("The following player has cards: ", &room.Players[i])

	}

	room.Players[0].Turn = true //Allowing first player to have their turn.

	//Start timer?

	//Send message to players game has started and whose turn it is.
	for i := 0; i < room.Pop; i++ {

		//msg := `{"type":"game_start"}`

		msg := GameMessage{ //Create game message to send to clients.
			Type:     "game_start",         //Setting type to game_start
			AddCards: room.Players[i].Hand, //sending cards to add.
		}

		fmt.Println("Sending start message players:")
		fmt.Printf("Player %d hand: %+v\n", i, room.Players[i].Hand)

		SendMessageToPlayer(room.Players[i], ConvertMsgToJson(&msg)) //Add Message to send queue and convert to json compatible.

	}

	room.Mu.Unlock()

}

func (r *Room) ManagePlActionInRm(player *Player, pMsg *PlayerMessage) { //Manages player actions/messages and uses mutex for thread-safety.

	r.Mu.Lock() //Locking the room mutex

	fmt.Println("Managing Player action.")
	r.LastActive = time.Now() //Update activity to show room is active.

	//Check message type and send to room if required.
	switch action := pMsg.Action; action {
	case "play_card": //If user is playing a card.
		fmt.Println("Managing Player action - switch case")
		PlayCard(r, player, pMsg)

	}

	r.EndTurn() //End Turn after action.

	r.Mu.Unlock() //Unlock after func has completed.

}

func (room *Room) FlipTurns() { //Method that flips player turns in room.
	for i := 0; i < room.Pop; i++ {
		room.Players[i].Turn = !room.Players[i].Turn
	}

}

func (room *Room) EndTurn() { //Method to send game state to all players.

	room.FlipTurns() //Flipping player turns after card has been played. Only allows 1 card per turn (might increase for balancing).

	msg := GameMessage{ //Create game message to send to clients.
		Type:       "game_state", //Setting type to game_start
		BoardState: room.SendBoardState(),
	}

	for i := 0; i < room.Pop; i++ {

		fmt.Println("Sending start message players:")
		fmt.Printf("Player %d hand: %+v\n", i, room.Players[i].Hand)

		SendMessageToPlayer(room.Players[i], ConvertMsgToJson(&msg)) //Add Message to send queue and convert to json compatible.

	}
}

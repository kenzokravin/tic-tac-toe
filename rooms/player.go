package rooms

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Player struct.
type Player struct {
	ID        uuid.UUID       //Unique ID
	Name      string          //Display name
	Turn      bool            //Tracks if able to place
	Hand      []*Card         //Tracks Cards in hand (used for validating actions)
	Conn      *websocket.Conn //The client's connection.
	SendQueue chan string     //Queue for writing messages to client.
	Mu        sync.Mutex      //Player connection mutex.
}

type GameMessage struct { //Game message for communicating turns to players.
	Type         string        `json:"type"`                      //Game message type (i.e. setup, turn etc)
	AddCards     []*Card       `json:"cards_to_add,omitempty"`    //Cards to add to hand.
	RemoveCards  []*Card       `json:"cards_to_remove,omitempty"` //Cards to remove from hand.
	TargetSlotID int           `json:"target_slot,omitempty"`     //The id of the target slot, used to convey target slots from enemy moves (i.e. placing a mark.)
	BoardState   []*MarkEffect `json:"board_state,omitempty"`     //Cards to add to hand.
}

type PlayerMessage struct { //Message struct for when players send messages.
	Action       string `json:"action"`                //Used to figure out message type (i.e. Play card or send chat etc)
	CardName     string `json:"card_name,omitempty"`   //Name of card used, if no card then omit.
	TargetSlotID int    `json:"target_slot,omitempty"` //The id of the target slot.
}

var defPlayer *Player = nil //Pointing to a null player. This is used to init card effects.

func ConvertMsgToJson(msg *GameMessage) string {

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		// handle error
		fmt.Println("Error marshaling:", err)

	}

	return string(jsonMsg)

}

func SendMessageToPlayer(player *Player, msg string) { //function to send a message to the desired player.
	player.SendQueue <- msg //Adding msg to player's msg queue.

}

func (p *Player) StartWriter() { //Method to start writer queue.
	fmt.Println("Start msg writer for", p.ID)
	go func() { //Starts go routine that constantly runs for player until disconnect.
		for msg := range p.SendQueue {

			p.Mu.Lock() //Lock mutex.
			fmt.Println("Sent msg")
			err := p.Conn.WriteMessage(websocket.TextMessage, []byte(msg)) //Writes message to player.
			p.Mu.Unlock()                                                  //Unlock after sending.
			if err != nil {
				fmt.Println("Write error:", err)
				break // exit if there's an error (e.g. client disconnects)
			}

		}
	}()
}

func (p *Player) Close() {
	p.Mu.Lock()

	close(p.SendQueue)
	p.Conn.Close()

	p.Mu.Unlock()

	fmt.Println("Closed player:", p.ID)
}

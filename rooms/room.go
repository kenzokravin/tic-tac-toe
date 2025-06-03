package rooms

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Room struct {
	ID      uuid.UUID
	State   string
	Pop     int
	Full    bool
	Board   *Board
	Players []*Player
}

type Card struct {
	Type        string  //Card type (i.e. playable)
	Name        string  //Card name (should be unique for each card)
	Description string  //Card description
	Rarity      float64 //Card Rarity, all card rarities should add to 1.0
	GraphicPath string
	MarkerPath  string
	ImpactType  string
	ImpactShape string
	MarkEffect  *MarkEffect //The effect the card has on the slots.
}

// Slot struct.
type Slot struct {
	ID      int //The number ID of the slot.
	Row     int //The slot row.
	Col     int //The slot column.
	Owner   *Player
	Effects []*MarkEffect //The effects currently on the slot.
}

type MarkEffect struct { //Mark Effects are the effects of the marks (These typically involving adding or subtracting health)
	Owner       *Player //The owner of the mark.
	Health      int     // The amount of health a mark has.
	GraphicPath string  //The path of the graphic (mark) to show.
	Damage      int     //How much damage the effect does to the slot.
}

// Player struct.
type Player struct {
	ID        uuid.UUID       //Unique ID
	Name      string          //Display name
	Turn      bool            //Tracks if able to place
	Hand      []*Card         //Tracks Cards in hand (used for validating actions)
	Conn      *websocket.Conn //The client's connection.
	SendQueue chan string     //Queue for writing messages to client.
}

type Board struct {
	Slots []*Slot
}

type GameMessage struct { //Game message for communicating turns to players.
	Type         string  `json:"type"`                      //Game message type (i.e. setup, turn etc)
	AddCards     []*Card `json:"cards_to_add,omitempty"`    //Cards to add to hand.
	RemoveCards  []*Card `json:"cards_to_remove,omitempty"` //Cards to remove from hand.
	TargetSlotID int     `json:"target_slot,omitempty"`     //The id of the target slot, used to convey target slots from enemy moves (i.e. placing a mark.)
}

type PlayerMessage struct { //Message struct for when players send messages.
	Action       string `json:"action"`                //Used to figure out message type (i.e. Play card or send chat etc)
	CardName     string `json:"card_name,omitempty"`   //Name of card used, if no card then omit.
	TargetSlotID int    `json:"target_slot,omitempty"` //The id of the target slot.
}

var cards = []Card{}        //An array that stores all possible card types.
var defPlayer *Player = nil //Pointing to a null player. This is used to init card effects.

func CreateCards() { //Creating all possible cards.

	crdDefault := Card{Type: "Null", Name: "Default", //Default card, used for nothing but testing.
		Description: "This card does nothing.",
		Rarity:      1.0,
		GraphicPath: "src/card_test_mark.png",
		MarkerPath:  "src/naught.svg",
		ImpactType:  "singular",
		MarkEffect: &MarkEffect{
			Owner:       defPlayer,
			Health:      1,
			GraphicPath: "src/naught.svg",
			Damage:      0,
		},
	}

	cards = append(cards, crdDefault)

	crdMark := Card{Type: "Playable", //The mark card, used to place a mark.
		Name:        "Mark",
		Description: "Place a mark in a square.",
		Rarity:      1.0,
		GraphicPath: "src/card_test_mark.png",
		MarkerPath:  "src/naught.svg",
		ImpactType:  "singular",
		MarkEffect: &MarkEffect{
			Owner:       defPlayer,
			Health:      1,
			GraphicPath: "src/naught.svg",
			Damage:      0,
		},
	}

	cards = append(cards, crdMark) //Adding to card list.

}

func CreateBoard() Board { //Creating and returning the Board filled with slots.

	board := Board{Slots: []*Slot{}}

	for i := 0; i < 3; i++ {
		for z := 0; z < 3; z++ {

			id := i*3 + z

			nSlot := Slot{ID: id, Row: i, Col: z}

			board.Slots = append(board.Slots, &nSlot)

		}
	}

	fmt.Println("Created game board in room.")

	return board

}

func StartRoomGame(room *Room) {

	room.State = "In Progress" //Setting Game state to playing.

	for i := 0; i < room.Pop; i++ { //Drawing Start Cards for players.
		DrawStartCards(room.Players[i])

		fmt.Println("Player has cards:", &room.Players[i])

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

		SendMessageToPlayer(room.Players[i], ConvertMsgToJson(&msg)) //Add Message to send queue and convert to json compatible.

	}

}

func DrawStartCards(player *Player) { //Drawing start cards.

	for i := 0; i < 3; i++ { //Draw 3 cards.

		player.Hand = append(player.Hand, DrawCard()) //Add cards to player's hand.

	}

}

func DrawCard() *Card { //Draws a card from the initialized cards using chance (math/rand). Must make sure we are dealing with card copies or not.

	cmpRarity := 0.0 //compound rarity used to check values.
	prevRarity := 0.0

	chance := rand.Float64() //chance value that is the card drawn. returns a value between [0.0,1.0]

	devCard := Card{Type: "Playable", Name: "Mark",
		Description: "This is a dev mark card.", Rarity: 1.0, GraphicPath: "src/card_test_mark.png", MarkerPath: "src/naught.svg"} //ONLY FOR DEVELOPMENT PURPOSES.

	for i := 0; i < len(cards); i++ {

		if i == 0 {
			prevRarity = 0.0
		} else {
			prevRarity += cards[i-1].Rarity //adds previous card rarity to val, creating the lower bounds.
		}

		if chance > float64(prevRarity) && chance < float64(cmpRarity+cards[i].Rarity) { //Check if chance is between range.

			fmt.Println("Card Drawn.")

			return &cards[i]

		}

		cmpRarity += cards[i].Rarity //Adds the current card rarity to value, creating the upper bounds.

		//return cards[i]

	}

	return &devCard //ONLY FOR DEVELOPMENT PURPOSES.

}

func PlayCard(room *Room, player *Player, pMsg *PlayerMessage) { //plays a card.

	if !player.Turn { //Check if player's turn and break function if not.
		fmt.Println("ERROR: Not Player's Turn.") //Send error to console. Will have to change to send the client an error message.
		return
	}

	isCardAvailable := false

	playedCard := &cards[0] //Creating a pointer to the default card.

	for i := 0; i < len(player.Hand); i++ { //Checking if card is in player's hand.
		if pMsg.CardName == player.Hand[i].Name {
			isCardAvailable = true

			playedCard = player.Hand[i]

		}

	}

	if !isCardAvailable {
		fmt.Println("ERROR: Cannot play Card that is not in hand.")
		return
	}

	switch playedCard.ImpactType { //Determine which slots to effect using impact type.
	case "singular": //This means a singular slot is effected.
		room.Board.Slots[pMsg.TargetSlotID].AddEffectToSlot(playedCard.MarkEffect)
	case "multiple": //Means multiple slots get affected.

	}

	room.FlipTurns() //Flipping player turns after card has been played.

}

func (sl *Slot) AddEffectToSlot(mEffect *MarkEffect) { //Method that adds the effect to the slot.

	sl.Effects = append(sl.Effects, mEffect)

}

func (b *Board) GetAffectedSlots(shape string, tarSlotID int) []*Slot {

	retSlots := []*Slot{} //Creating return slot array.

	tSlot := b.ReturnSlotFromID(tarSlotID)

	switch shape {
	case "lines": //If the shape is similar to a bomberman.
		for i := 0; i < len(b.Slots); i++ { //Cycle through slots to determine if affected or not.
			switch b.Slots[i].ID {
			case tarSlotID: //If Same.
				retSlots = append(retSlots, b.Slots[i])
				continue
			case tarSlotID - 3: //If above
				retSlots = append(retSlots, b.Slots[i])
				continue
			case tarSlotID - 6: //If two above
				retSlots = append(retSlots, b.Slots[i])
				continue
			case tarSlotID + 3: //If below
				retSlots = append(retSlots, b.Slots[i])
				continue
			case tarSlotID + 6: //If two above
				retSlots = append(retSlots, b.Slots[i])
				continue
			case tarSlotID - 1: //If to the left
				retSlots = append(retSlots, b.Slots[i])
				continue
			case tarSlotID - 2: //If two to the left.
				retSlots = append(retSlots, b.Slots[i])
				continue
			case tarSlotID + 1: //If to the right
				retSlots = append(retSlots, b.Slots[i])
				continue
			case tarSlotID + 2: //If two to the right.
				retSlots = append(retSlots, b.Slots[i])
				continue
			}
		}
	case "radius": //If the shape is a radius (1 slot around target Slot)
		for i := 0; i < len(b.Slots); i++ {

			if b.Slots[i].ID == tarSlotID {
				retSlots = append(retSlots, b.Slots[i])
				continue
			}

			if b.Slots[i].Row == tSlot.Row { //If the row is equal
				if tSlot.Col == 1 { //If the target slot column is the middle column, add all in row.
					retSlots = append(retSlots, b.Slots[i])
					continue
				} else if tSlot.Col == 0 { //If target is in left column.
					if b.Slots[i].ID == tarSlotID+1 {
						retSlots = append(retSlots, b.Slots[i])
						continue
					}
				} else if tSlot.Col == 2 { //If target is in right column
					if b.Slots[i].ID == tarSlotID-1 {
						retSlots = append(retSlots, b.Slots[i])
						//Need to append tarSlot as well potentially.
						continue
					}
				}
			}

			if b.Slots[i].Col == tSlot.Col {
				if tSlot.Row == 1 { //If the target slot column is the middle column, add all in row.
					retSlots = append(retSlots, b.Slots[i])
					continue
				} else if tSlot.Row == 0 { //If target is in top row
					if b.Slots[i].ID == tarSlotID+3 { //If ID matches ID in r below.
						retSlots = append(retSlots, b.Slots[i])
						continue
					}
				} else if tSlot.Row == 2 { //If target is in bottom row
					if b.Slots[i].ID == tarSlotID-3 { //If ID matches ID in r above.
						retSlots = append(retSlots, b.Slots[i])
						continue
					}
				}
			}

			if b.Slots[i].Row == 1 && tSlot.Row != 1 { //If slot is in middle row and target is not.
				if tSlot.Col == 1 { //If target is in middle column
					if b.Slots[i].Col != 1 { //If the slot is not in the centre column, add to retSlots.
						retSlots = append(retSlots, b.Slots[i])
						continue
					}
				}
			}

			if b.Slots[i].Row != 1 && tSlot.Row == 1 {
				if tSlot.Col == 1 { //If target is in middle column
					if b.Slots[i].Col != 1 { //If the slot is not in the centre column, add to retSlots.
						retSlots = append(retSlots, b.Slots[i])
						continue
					}
				}
			}

		}

	}

	return retSlots //returning slot array.
}

func ProcessSlots(room *Room) {

	for i := 0; i < len(room.Board.Slots); i++ { //Cycle through slots 0-9
		for z := 0; z < len(room.Board.Slots[i].Effects); z++ { //Cycle through effects.

		}
	}
}

func (b *Board) ReturnSlotFromID(slotID int) *Slot { //Method that returns the slot from slotID

	var retSlot *Slot = nil

	for i := 0; i < len(b.Slots); i++ {
		if b.Slots[i].ID == slotID {
			return b.Slots[i]

		}
	}

	return retSlot

}

func SendMessageToPlayer(player *Player, msg string) { //function to send a message to the desired player.
	player.SendQueue <- msg //Adding msg to player's msg queue.

}

func ConvertMsgToJson(msg *GameMessage) string {

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		// handle error
		fmt.Println("Error marshaling:", err)

	}

	return string(jsonMsg)

}

func (p *Player) StartWriter() { //Method to start writer queue.
	fmt.Println("Start msg writer for", p.ID)
	go func() { //Starts go routine that constantly runs for player until disconnect.
		for msg := range p.SendQueue {
			err := p.Conn.WriteMessage(websocket.TextMessage, []byte(msg)) //Writes message to player.
			if err != nil {
				fmt.Println("Write error:", err)
				break // exit if there's an error (e.g. client disconnects)
			}
		}
	}()
}

func (room *Room) FlipTurns() { //Method that flips player turns in room.
	for i := 0; i < room.Pop; i++ {
		room.Players[i].Turn = !room.Players[i].Turn
	}

}

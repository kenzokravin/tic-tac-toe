package rooms

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"

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
	Mu      sync.Mutex
}

type Card struct {
	Type        string  //Card type (i.e. attack)
	Name        string  //Card name (should be unique for each card)
	Description string  //Card description
	Rarity      float64 //Card Rarity, all card rarities should add to 1.0
	GraphicPath string
	MarkerPath  string
	ImpactType  string      //Impact Type decides if many or singular slots are effected.
	ImpactShape string      //Impact Shape is the shape of the effect. (i.e. does it strike rows or a radius all around etc)
	MarkEffect  *MarkEffect //The effect the card has on the slots.
}

// Slot struct. A board is composed of 9 slots.
type Slot struct {
	ID      int //The number ID of the slot.
	Row     int //The slot row.
	Col     int //The slot column.
	Owner   *Player
	Effects []*MarkEffect //The effects currently on the slot.
}

type MarkEffect struct { //Mark Effects are the effects of the marks (These typically involving adding or subtracting health). Each card has a mark (effect).
	Owner         *Player //The owner of the mark.
	Health        int     // The amount of health a mark has.
	GraphicPath   string  //The path of the graphic (mark) to show.
	Damage        int     //How much damage the effect does to the slot.
	IsDestroyable bool    //If Mark is destroyable.
	IsStackable   bool    //If Mark can be added to effect stack.
	IsBlocking    bool    //If Mark prevents mark placements
	DamageType    string  //Used to check if damage is pure (cannot be blocked)
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

var cards = []*Card{}    //An array that stores all possible card types.
var cardsMu sync.RWMutex //Read-Write Mutex allows multiple readers, one write.

var defPlayer *Player = nil //Pointing to a null player. This is used to init card effects.

func CreateCards() { //Creating all possible cards.

	cardsMu.Lock()
	defer cardsMu.Unlock()

	crdDefault := Card{Type: "Null", Name: "Default", //Default card, used for nothing but testing.
		Description: "This card does nothing.",
		Rarity:      0.0,
		GraphicPath: "src/card_test_mark.png",
		MarkerPath:  "src/naught.svg",
		ImpactType:  "singular",
		ImpactShape: "null",
		MarkEffect: &MarkEffect{
			Owner:         defPlayer,
			Health:        1,
			GraphicPath:   "src/naught.svg",
			Damage:        0,
			IsDestroyable: false,
			IsStackable:   false,
		},
	}

	cards = append(cards, &crdDefault)

	crdMark := Card{Type: "attack", //The mark card, used to place a mark.
		Name:        "Mark",
		Description: "Place a mark in a square.",
		Rarity:      0.5,
		GraphicPath: "src/card_test_mark.png",
		MarkerPath:  "src/naught.svg",
		ImpactType:  "singular",
		ImpactShape: "null",
		MarkEffect: &MarkEffect{
			Owner:         defPlayer,
			Health:        1,
			GraphicPath:   "src/naught.svg",
			Damage:        0, //Cannot damage card.
			IsDestroyable: true,
			IsStackable:   true,
			IsBlocking:    true,
			DamageType:    "place", //"place" means it can be blocked and the slot must not have a blocking mark effect (i.e. an opponent mark, but it can have an invisible mark.)
		},
	}

	cards = append(cards, &crdMark) //Adding to card list.

	crdBomb := Card{Type: "attack", //The mark card, used to place a mark.
		Name:        "Bomb",
		Description: "Destroys all marks in a 1 slot radius.",
		Rarity:      0.5,
		GraphicPath: "src/card_test_mark.png",
		MarkerPath:  "src/naught.svg",
		ImpactType:  "multiple",
		ImpactShape: "radius",
		MarkEffect: &MarkEffect{
			Owner:         defPlayer,
			Health:        0,
			GraphicPath:   "src/naught.svg",
			Damage:        100,
			IsDestroyable: false,
			IsStackable:   false,
			IsBlocking:    false,
			DamageType:    "pure", //"pure" means it cannot be blocked or affected by protective buffs.
		},
	}

	cards = append(cards, &crdBomb) //Adding to card list.

	crdDyn := Card{Type: "attack", //The mark card, used to place a mark.
		Name:        "Dynamite",
		Description: "Destroys all marks in the same row and column.",
		Rarity:      0.5,
		GraphicPath: "src/card_test_mark.png",
		MarkerPath:  "src/naught.svg",
		ImpactType:  "multiple",
		ImpactShape: "lines",
		MarkEffect: &MarkEffect{
			Owner:         defPlayer,
			Health:        0,
			GraphicPath:   "src/naught.svg",
			Damage:        100,
			IsDestroyable: false,
			IsStackable:   false,
			IsBlocking:    false,
			DamageType:    "pure",
		},
	}

	cards = append(cards, &crdDyn) //Adding to card list.

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

	room.Mu.Lock()
	defer room.Mu.Unlock()

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

func (r *Room) ManagePlActionInRm(player *Player, pMsg *PlayerMessage) { //Manages player actions/messages and uses mutex for thread-safety.

	r.Mu.Lock()         //Locking the room mutex
	defer r.Mu.Unlock() //Unlock after func has completed.

	//Check message type and send to room if required.
	switch action := pMsg.Action; action {
	case "play_card": //If user is playing a card.
		PlayCard(r, player, pMsg)

	}

}

func DrawStartCards(player *Player) { //Drawing start cards.

	for i := 0; i < 3; i++ { //Draw 3 cards.

		player.Hand = append(player.Hand, DrawCard()) //Add cards to player's hand.

	}

}

func DrawCard() *Card { //Draws a card from the initialized cards using chance (math/rand). Must make sure we are dealing with card copies or not.

	cardsMu.RLock()
	defer cardsMu.RUnlock()

	// cmpRarity := 0.0 //compound rarity used to check values.
	// prevRarity := 0.0

	chance := rand.Float64() //chance value that is the card drawn. returns a value between [0.0,1.0]

	total := 0.0
	for _, c := range cards { //Normalizing the rarity to ensure it equals 1.0
		total += c.Rarity
	}
	for i := range cards {
		cards[i].Rarity /= total
	}

	var cumulative float64 //Using weighted rarity.
	for i := 0; i < len(cards); i++ {
		cumulative += cards[i].Rarity
		if chance <= cumulative {
			fmt.Println("Card Drawn.")
			return cards[i]
		}
	}

	devCard := Card{Type: "Playable", Name: "Mark",
		Description: "This is a dev mark card.", Rarity: 1.0, GraphicPath: "src/card_test_mark.png", MarkerPath: "src/naught.svg"} //ONLY FOR DEVELOPMENT PURPOSES, Only create if no cards are available (Should never happen)

	return &devCard //ONLY FOR DEVELOPMENT PURPOSES.

}

func PlayCard(room *Room, player *Player, pMsg *PlayerMessage) { //Plays a card.

	if !player.Turn { //Check if player's turn and break function if not.
		fmt.Println("ERROR: Not Player's Turn.") //Send error to console. Will have to change to send the client an error message.
		return
	}

	isCardAvailable := false

	playedCard := cards[0] //Creating a pointer to the default card.

	for i := 0; i < len(player.Hand); i++ { //Checking if card is in player's hand.
		if pMsg.CardName == player.Hand[i].Name {
			isCardAvailable = true      //Set card availability to true.
			playedCard = player.Hand[i] //Set the played card to the reference of the card.
			break                       //break out of loop.
		}
	}

	if !isCardAvailable { //If card not available, break function.
		fmt.Println("ERROR: Cannot play Card that is not in hand.")
		return
	}

	switch playedCard.Type { //Checking card type.
	case "attack": //If card is an attack type. (i.e. damages other marks, places marks etc)
		switch playedCard.ImpactType { //Determine which slots to effect using impact type.
		case "singular": //This means a singular slot is effected.
			//Might have to calculate if it can be played or not (i.e. is slot valid)
			room.Board.Slots[pMsg.TargetSlotID].AddEffectToSlot(playedCard.MarkEffect) //Add card effect to slot.
		case "multiple": //Means multiple slots get affected.
			slotsToAffect := room.Board.GetAffectedSlots(playedCard.ImpactShape, pMsg.TargetSlotID) //Retrieving slots to affect.

			room.Board.ApplyDamageToSlotsFromCard(slotsToAffect, playedCard.MarkEffect)

			for i := 0; i < len(slotsToAffect); i++ { //Cycle through slots and add effect.
				slotsToAffect[i].AddEffectToSlot(playedCard.MarkEffect) //Adding effect.
			}

		}
	case "buff": //If card is a buff type (i.e. effects that add health.)

	}

	room.FlipTurns() //Flipping player turns after card has been played. Only allows 1 card per turn (might increase for balancing).

}

func (sl *Slot) IsOpen()

func (sl *Slot) AddEffectToSlot(mEffect *MarkEffect) { //Method that adds the effect to the slot.

	if !mEffect.IsStackable { //If effect cannot be stacked, do not add to stack.
		return
	}

	sl.Effects = append(sl.Effects, mEffect)

}

func (b *Board) ApplyDamageToSlotsFromCard(slots []*Slot, mEffect *MarkEffect) { //Applies damage to slots.

	if mEffect.Damage <= 0 { //If effect cannot damage, return from function.
		return
	}

	for i := 0; i < len(slots); i++ { //Cycle through slots.
		newEffects := []*MarkEffect{} // Adjust type as needed

		for _, eff := range slots[i].Effects {

			if !eff.IsDestroyable { //If not destroyable, skip damage taking (prevents destruction of marks that are immune to damage)
				continue
			}

			eff.Health -= mEffect.Damage
			if eff.Health > 0 {
				newEffects = append(newEffects, eff)
			}
			// Otherwise: Effect is dead, so exclude it
		}

		slots[i].Effects = newEffects // Replace with filtered effects.
	}

}

func (b *Board) GetAffectedSlots(shape string, tarSlotID int) []*Slot { //Method that retrieves an array of slots to be affected by the card and it's impact shape.

	retSlots := []*Slot{} //Creating return slot array.

	tSlot := b.ReturnSlotFromID(tarSlotID) //Get a pointer to the targetSlot, to use row and col data.

	switch shape {
	case "lines": //If the shape is similar to a bomberman.
		for i := 0; i < len(b.Slots); i++ { //Cycle through slots to determine if affected or not.
			dRow := abs(b.Slots[i].Row - tSlot.Row) //Getting difference in rows.
			dCol := abs(b.Slots[i].Col - tSlot.Col) //Getting difference in columns.

			if dRow < 1 && dCol <= 2 {
				retSlots = append(retSlots, b.Slots[i])
				continue
			}

			if dRow <= 2 && dCol < 1 {
				retSlots = append(retSlots, b.Slots[i])
				continue
			}
		}
	case "radius": //If the shape is a radius (1 slot around target Slot)
		for i := 0; i < len(b.Slots); i++ {

			dRow := abs(b.Slots[i].Row - tSlot.Row) //Getting difference in rows.
			dCol := abs(b.Slots[i].Col - tSlot.Col) //Getting difference in columns.

			if dRow <= 1 && dCol <= 1 {
				retSlots = append(retSlots, b.Slots[i])
				continue
			}

		}

	}

	return retSlots //returning slot array.
}

func (b *Board) CheckBoardWin() { //Method that checks if player has a valid line of marks.

}

func (b *Board) RemoveDeadMarks() {

	for i := 0; i < len(b.Slots); i++ { //Cycle through slots 0-9
		newEffects := []*MarkEffect{} // Adjust type as needed

		for _, eff := range b.Slots[i].Effects {
			if eff.Health <= 0 { //If health is less or equal to 0, skip over and leave.
				continue
			} else {
				newEffects = append(newEffects, eff) //If mark has health, keep.
			}
		}

		b.Slots[i].Effects = newEffects //Reassigning new effects array.

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

func (room *Room) FlipTurns() { //Method that flips player turns in room.
	for i := 0; i < room.Pop; i++ {
		room.Players[i].Turn = !room.Players[i].Turn
	}

}

//----------------------------------------------------------------------------------------
//---------------------------------Utility Functions--------------------------------------
//----------------------------------------------------------------------------------------

func abs(x int) int { //returns the absolute value.
	if x < 0 {
		return -x
	}
	return x
}

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
			err := p.Conn.WriteMessage(websocket.TextMessage, []byte(msg)) //Writes message to player.
			if err != nil {
				fmt.Println("Write error:", err)
				break // exit if there's an error (e.g. client disconnects)
			}
		}
	}()
}

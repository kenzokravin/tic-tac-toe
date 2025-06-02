package rooms

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Room struct {
	ID      string
	State   string
	Pop     int
	Full    bool
	Board   Board
	Players []*Player
}

var cards = []Card{} //An array that stores all possible card types.

type Card struct {
	Type        string  //Card type (i.e. playable)
	Name        string  //Card name
	Description string  //Card description
	Rarity      float64 //Card Rarity, all card rarities should add to 1.0
	GraphicPath string
	MarkerPath  string
	ImpactType  string
}

// Slot struct.
type Slot struct {
	ID  int
	Row int
	Col int
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
	Type        string  `json:"type"` //Game message type (i.e. setup, turn etc)
	AddCards    []*Card `json:"cards_to_add,omitempty"`
	RemoveCards []*Card `json:"cards_to_remove,omitempty"` //Card Hand, sends to the player what their current hand is.
}

func CreateCards() { //Creating all possible cards.

	crdMark := Card{Type: "Playable", Name: "Mark",
		Description: "Place a mark in a square.", Rarity: 1.0, GraphicPath: "/src/card_test.png", MarkerPath: "src/naught.svg", ImpactType: "singular"}

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

func DrawStartCards(player *Player) {

	for i := 0; i < 3; i++ { //Draw 3 cards.

		player.Hand = append(player.Hand, DrawCard()) //Add cards to player's hand.

	}

}

func DrawCard() *Card { //draws a card from the initialized cards using chance (math/rand).

	cmpRarity := 0.0 //compound rarity used to check values.
	prevRarity := 0.0

	chance := rand.Float64() //chance value that is the card drawn. returns a value between [0.0,1.0]

	devCard := Card{Type: "Playable", Name: "Mark",
		Description: "Place a mark in a square.", Rarity: 1.0, GraphicPath: "/src/card_test.png", MarkerPath: "src/naught.svg"} //ONLY FOR DEVELOPMENT PURPOSES.

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

func PlayCard(card Card) { //plays a card.

}

func SendMessageToPlayer(player *Player, msg string) { //function to send a message to the desired player.
	player.SendQueue <- msg //Adding msg to player's msg queue.

	// fmt.Printf("About to send to player %s\n", player.ID)
	// fmt.Printf("SendQueue length: %d\n", len(player.SendQueue))
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

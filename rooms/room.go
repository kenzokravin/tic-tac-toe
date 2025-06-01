package rooms

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

type Room struct {
	ID       string
	State    string
	Capacity int
	Full     bool
}

var cards = []Card{} //An array that stores all possible card types.

type Card struct {
	Type        string  //Card type (i.e. playable)
	Name        string  //Card name
	Description string  //Card description
	Rarity      float64 //Card Rarity, all card rarities should add to 1.0
	GraphicPath string
	MarkerPath  string
}

type Slot struct {
	ID  int
	Row int
	Col int
}

type Player struct {
	ID   uuid.UUID
	Name string
}

func CreateCards() { //Creating all possible cards.

	crdMark := Card{Type: "Playable", Name: "Mark",
		Description: "Place a mark in a square.", Rarity: 1.0, GraphicPath: "/src/card_test.png", MarkerPath: "src/naught.svg"}

	cards = append(cards, crdMark) //Adding to card list.

}

func DrawCard() Card {

	cmpRarity := 0.0 //compound rarity used to check values.
	prevRarity := 0.0

	chance := rand.Float64() //chance value that is the card drawn. returns a value between [0.0,1.0]

	for i := 0; i < len(cards); i++ {

		if i == 0 {
			prevRarity = 0.0
		} else {
			prevRarity += cards[i-1].Rarity //adds previous card rarity to val, creating the lower bounds.
		}

		if chance > float64(prevRarity) && chance < float64(cmpRarity+cards[i].Rarity) { //Check if chance is between range.

			fmt.Println("Card Drawn.")

			return cards[i]

		}

		cmpRarity += cards[i].Rarity //Adds the current card rarity to value, creating the upper bounds.

	}

}

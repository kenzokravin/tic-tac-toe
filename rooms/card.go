package rooms

import "sync"

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

type MarkEffect struct { //Mark Effects are the effects of the marks (These typically involving adding or subtracting health). Each card has a mark (effect).
	Owner         *Player //The owner of the mark.
	Health        int     // The amount of health a mark has.
	GraphicPath   string  //The path of the graphic (mark) to show.
	Damage        int     //How much damage the effect does to the slot.
	IsDestroyable bool    //If Mark is destroyable.
	IsStackable   bool    //If Mark can be added to effect stack.
	IsBlocking    bool    //If Mark prevents mark placements
	DamageType    string  //Used to check if damage is pure (cannot be blocked)
	IsWinEffect   bool    //If the effect can be considered as winnable (i.e. it is a valid mark.)
	IsDisplayable bool
}

var cards = []*Card{}    //An array that stores all possible card types.
var cardsMu sync.RWMutex //Read-Write Mutex allows multiple readers, one write.

func CreateCards() []*Card { //Creating all possible cards.

	cardsMu.Lock()
	cardsToRet := []*Card{}

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

	cardsToRet = append(cardsToRet, &crdDefault)

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
			IsWinEffect:   true,
		},
	}

	cardsToRet = append(cardsToRet, &crdMark) //Adding to card list.

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
			IsWinEffect:   false,
		},
	}

	cardsToRet = append(cardsToRet, &crdBomb) //Adding to card list.

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
			IsWinEffect:   false,
		},
	}

	cardsToRet = append(cardsToRet, &crdDyn) //Adding to card list.

	cards = cardsToRet

	cardsMu.Unlock()

	return cardsToRet

}

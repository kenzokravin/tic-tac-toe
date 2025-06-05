package rooms

import (
	"fmt"
	"math/rand"
)

type Board struct {
	Slots []*Slot
}

// Slot struct. A board is composed of 9 slots.
type Slot struct {
	ID      int //The number ID of the slot.
	Row     int //The slot row.
	Col     int //The slot column.
	Owner   *Player
	Effects []*MarkEffect //The effects currently on the slot.
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

func DrawStartCards(player *Player) { //Drawing start cards.

	cardsMu.RLock()

	for i := 0; i < 3; i++ { //Draw 3 cards.

		player.Hand = append(player.Hand, DrawCard()) //Add cards to player's hand.

	}
	cardsMu.RUnlock()

}

func DrawCard() *Card { //Draws a card from the initialized cards using chance (math/rand). Must make sure we are dealing with card copies or not.

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

func (sl *Slot) AddEffectToSlot(mEffect *MarkEffect) { //Method that adds the effect to the slot.

	if !mEffect.IsStackable { //If effect cannot be stacked, do not add to stack.
		return
	}

	fmt.Println("Adding slot effect.")

	sl.Effects = append(sl.Effects, mEffect)

}

func PlayCard(room *Room, player *Player, pMsg *PlayerMessage) { //Plays a card.

	if !player.Turn { //Check if player's turn and break function if not.
		fmt.Println("ERROR: Not Player's Turn.") //Send error to console. Will have to change to send the client an error message.
		return
	}

	isCardAvailable := false

	//This throws an error as cards need to be created.
	playedCard := cards[0] //Creating a pointer to the default card.

	fmt.Println("Card name from data.", pMsg.CardName)

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

	fmt.Println("Managing Player action - Playing card.")

	switch playedCard.Type { //Checking card type.
	case "attack": //If card is an attack type. (i.e. damages other marks, places marks etc)
		switch playedCard.ImpactType { //Determine which slots to effect using impact type.
		case "singular": //This means a singular slot is effected.
			//Might have to calculate if it can be played or not (i.e. is slot valid)
			fmt.Println("Playing singular card.")
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

func (b *Board) CheckBoardWin() { //Method that checks if player has a valid line of marks.

}

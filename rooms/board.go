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
	ID      int           //The number ID of the slot.
	Row     int           //The slot row.
	Col     int           //The slot column.
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

func (sl *Slot) AddEffectToSlot(mEffect *MarkEffect, player *Player) { //Method that adds the effect to the slot.

	if !mEffect.IsStackable { //If effect cannot be stacked, do not add to stack.
		fmt.Println("Effect is not stackable.")
		return
	}

	fmt.Println("Adding slot effect.")

	mEffect.Owner = player.ID

	sl.Effects = append(sl.Effects, mEffect)

}

func PlayCard(room *Room, player *Player, pMsg *PlayerMessage) { //Plays a card.

	if !player.Turn { //Check if player's turn and break function if not.
		fmt.Println("ERROR: Not Player's Turn.") //Send error to console. Will have to change to send the client an error message.

		msg := GameMessage{ //Create game message to send to clients.
			Type: "error", //Setting type to error
		}

		SendMessageToPlayer(player, ConvertMsgToJson(&msg)) //Add Message to send queue and convert to json compatible.

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
			fmt.Println("Found Available Card: ", playedCard)
			break //break out of loop.
		}
	}

	if !isCardAvailable { //If card not available, break function.
		fmt.Println("ERROR: Cannot play Card that is not in hand.")
		return
	}

	fmt.Println("Target Slot is: ", pMsg.TargetSlotID)

	if pMsg.TargetSlotID >= 9 || pMsg.TargetSlotID < 0 {
		fmt.Println("ERROR: Invalid Target Slot. ID out of bounds.")
		return
	}

	switch playedCard.Type { //Checking card type.
	case "attack": //If card is an attack type. (i.e. damages other marks, places marks etc)
		switch playedCard.ImpactType { //Determine which slots to effect using impact type.
		case "singular": //This means a singular slot is effected.
			//Might have to calculate if it can be played or not (i.e. is slot valid)
			fmt.Println("Playing singular card.")
			room.Board.Slots[pMsg.TargetSlotID].AddEffectToSlot(playedCard.MarkEffect, player) //Add card effect to slot.
		case "multiple": //Means multiple slots get affected.
			fmt.Println("Playing multiple card.")
			slotsToAffect := room.Board.GetAffectedSlots(playedCard.ImpactShape, pMsg.TargetSlotID) //Retrieving slots to affect.

			room.Board.ApplyDamageToSlotsFromCard(slotsToAffect, playedCard.MarkEffect)

			fmt.Println("slots to affect: ", slotsToAffect)

			for i := 0; i < len(slotsToAffect); i++ { //Cycle through slots and add effect.
				slotsToAffect[i].AddEffectToSlot(playedCard.MarkEffect, player) //Adding effect.
			}

		}
	case "buff": //If card is a buff type (i.e. effects that add health.)

	}

	msg := GameMessage{ //Create game message to send to clients.
		Type:         "play_card_success", //Setting type to successful card play.
		TargetSlotID: &pMsg.TargetSlotID,
	}

	SendMessageToPlayer(player, ConvertMsgToJson(&msg))

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

func (rm *Room) CheckBoardWin(b *Board) {
	for p := 0; p < len(rm.Players); p++ {
		winningSlots := []*Slot{}

		for i := 0; i < len(b.Slots); i++ {
			for z := 0; z < len(b.Slots[i].Effects); z++ {
				if b.Slots[i].Effects[z].IsWinEffect && b.Slots[i].Effects[z].Owner == rm.Players[p].ID {
					winningSlots = append(winningSlots, b.Slots[i])
				}
			}
		}

		for i := 0; i < len(winningSlots); i++ {
			if SlotWinCheck(winningSlots, i, "null", 1) {
				fmt.Println("Player", rm.Players[p].Name, "wins!")
				return
			}
		}
	}
}

func SlotWinCheck(winningSlots []*Slot, i int, winType string, valSlots int) bool { //Might need to test, but should allow for a win.
	winningType := winType

	for z := i + 1; z < len(winningSlots); z++ {
		dRow := abs(winningSlots[i].Row - winningSlots[z].Row)
		dCol := abs(winningSlots[i].Col - winningSlots[z].Col)

		foundValid := false

		if winningType == "col" || winningType == "null" {
			if dRow == 1 && dCol == 0 {
				winningType = "col"
				valSlots++
				foundValid = true
			}
		}
		if winningType == "row" || winningType == "null" {
			if dRow == 0 && dCol == 1 {
				winningType = "row"
				valSlots++
				foundValid = true
			}
		}
		if winningType == "diag" || winningType == "null" {
			if dRow == 1 && dCol == 1 {
				winningType = "diag"
				valSlots++
				foundValid = true
			}
		}

		if !foundValid {
			continue
		}

		if valSlots == 3 {
			fmt.Println("Game won.")
			return true
		}

		if SlotWinCheck(winningSlots, z, winningType, valSlots) {
			return true
		}
	}

	return false
}

func (rm *Room) SendBoardState() []*MarkEffect {

	displayMarks := []*MarkEffect{}

	for _, sl := range rm.Board.Slots {

		for i := len(sl.Effects) - 1; i <= 0; i-- { //Reverse Searching to find top most graphic.

			if sl.Effects[i].IsDisplayable {
				displayMarks = append(displayMarks, sl.Effects[i])
			}

		}

	}

	return displayMarks

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

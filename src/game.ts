import * as PIXI from "pixi.js";
import { send } from "./client";
import eventBus from "./client";



(async () => {

  PIXI.settings.RESOLUTION = window.devicePixelRatio;

  const app = new PIXI.Application<HTMLCanvasElement>({ background: 0xffffff, resizeTo: window ,autoDensity:true});

  document.body.appendChild(app.view);
  const container = new PIXI.Container();

  app.stage.sortableChildren = true; //Allows z-index to be used.

  app.stage.addChild(container);

  ///------------------------ INIT VAR ----------------------

  //Creating the tic-tac-toe board.
  //The slot interface ensures we have all data for each slot.
  interface Slot {
    id: number;
    x:  number;
    y:  number;
    row: number;
    col: number;
    colour: number;
    slotGraphic: PIXI.Graphics;
    markerGraphic: PIXI.Sprite;
  }

  let slotSize = window.innerWidth*0.04;

  let board: {
    slots: Slot[],
    x: number,
    y: number
  } = {
    slots: [],
    x: window.innerWidth/2,
    y: window.innerHeight/2
  };

  let slotCounter = 0;

  //Creating Card Hand
  let cardSpriteScaler = 1;
  let cardHandSpace = window.innerWidth * 0.01;

  interface Card {
    name:string,
    description:string,
    selected:boolean,
    graphicPath: string,
    markerSprite: PIXI.Sprite,
    sprite: PIXI.Sprite,
    targetX:number,
    targetY:number,
    x:number,
    y:number
  }

//   type Card struct {
// 	Type        string  //Card type (i.e. attack)
// 	Name        string  //Card name (should be unique for each card)
// 	Description string  //Card description
// 	Rarity      float64 //Card Rarity, all card rarities should add to 1.0
// 	GraphicPath string
// 	MarkerPath  string
// 	ImpactType  string      //Impact Type decides if many or singular slots are effected.
// 	ImpactShape string      //Impact Shape is the shape of the effect. (i.e. does it strike rows or a radius all around etc)
// 	MarkEffect  *MarkEffect //The effect the card has on the slots.
// }



  const cardHand: Card[] = [];
  let selectedCard: Card | undefined;
  let handHeight = window.innerHeight *0.325;
  let cardSelectRaise = window.innerHeight * 0.032;

  let cardDiscPoint = { //The point where cards go to be discarded.
    x:window.innerWidth,
    y:window.innerHeight/2,
  }


  await PIXI.Assets.load('src/Inter-VariableFont_opsz,wght.ttf');

  const font = new FontFace('Inter', 'url(src/Inter-VariableFont_opsz,wght.ttf)');
  await font.load();
  document.fonts.add(font);

  const style = new PIXI.TextStyle({
    fontFamily: 'Inter',
    fontSize: 16,
    fill: '#000000',
  });

  let crdText = new PIXI.Text( "hi",style);

  let textPadding = 8;
  let descBox = new PIXI.Graphics();
  

  crdText.x = textPadding;
  crdText.y = textPadding;

  let descContainer = new PIXI.Container();

  descContainer.position.x = window.innerWidth/2 + 200;

  


  //------------------------------- INIT Functions --------------------------


  ScaleSize();
  CentreBoard();
  SetSlotListeners();

  // ----------------------------- Functions ------------------------------

 

  //Function used to centre the board in space.
  function CentreBoard() {
    board.x = window.innerWidth/2;
    board.y = window.innerHeight/2;

    slotCounter = 0;

    let slotMarkers = []; //Creating slot markers array to allow for markers to be moved as well.

    for (const slot of board.slots) {
      slotMarkers.push(slot.markerGraphic); //Adding markers into slotMarkers array.

      app.stage.removeChild(slot.slotGraphic);
      slot.slotGraphic.destroy(); // destroy graphics to resize.
    }
    board.slots = []; // Reset array

    
    for (let i = 0; i < 3; i++) {
      for (let z = 0; z < 3; z++) {

        const id = i * 3 + z;
        const colour = slotCounter % 2 === 0 ? 0xd3d3d3 : 0xe8e8e8;
        let x =  ((board.x - 1.5 * slotSize) + slotSize * z);
        let y =  ((board.y - 1.5 * slotSize) + slotSize * i);
        const row = i;
        const col = z;
        const markerGraphic = slotMarkers[id];

       

        const slotGraphic = new PIXI.Graphics();

        slotGraphic.beginFill(colour);
        slotGraphic.drawRect(
          x,
          y,
          slotSize,
            slotSize
        );
        slotGraphic.endFill();

        slotGraphic.zIndex = 0;

        slotGraphic.cursor = 'pointer';

        app.stage.addChild(slotGraphic);

         if(markerGraphic !== undefined) {
          markerGraphic.zIndex = 2;

          markerGraphic.position.x = x + (slotSize*0.5) - markerGraphic.width/2; 
          markerGraphic.position.y = y + (slotSize*0.5) - markerGraphic.width/2;
        }

        const slot:Slot = {
          id,
          x,
          y,
          row,
          col,
          colour,
          slotGraphic,
          markerGraphic
        };

        board.slots.push(slot);
        slotCounter++;
        }
    }
  }

  //Function used to dynamically resize the playing board and cards.
  function ScaleSize() {
    if(window.innerWidth < 255) {
        slotSize = window.innerWidth*0.14;
        cardSpriteScaler = 0.18;
         cardHandSpace =  window.innerWidth * 0.01;

      } else if (window.innerWidth >= 255 && window.innerWidth < 370)
      {
        slotSize = window.innerWidth*0.12;
        cardSpriteScaler = 0.2;
        cardHandSpace =  window.innerWidth * 0.01;
      }
      else if (window.innerWidth >= 370 && window.innerWidth < 512)
      {
        slotSize = window.innerWidth*0.1;
        cardSpriteScaler = 0.24;
        cardHandSpace =  window.innerWidth * 0.01;
      }
      else if (window.innerWidth >= 512 && window.innerWidth < 900)
      {
        slotSize = window.innerWidth*0.075;
        cardSpriteScaler = 0.28;
        cardHandSpace =  window.innerWidth * 0.01;
      }
      else if (window.innerWidth >= 900 && window.innerWidth < 1000)
      {
        slotSize = window.innerWidth*0.07;
        cardSpriteScaler = 0.3;
        cardHandSpace =  window.innerWidth * 0.01;
      }
      else if (window.innerWidth >= 1000 && window.innerWidth < 1200)
      {
        slotSize = window.innerWidth*0.06;
        cardSpriteScaler = 0.3;
        cardHandSpace =  window.innerWidth * 0.01;
      }
      else 
      {
        slotSize = window.innerWidth*0.05;
        cardSpriteScaler = 0.3;
        cardHandSpace =  window.innerWidth * 0.01;
      }

  }

  function SetSlotListeners() {
    for (const slot of board.slots) {
      slot.slotGraphic.eventMode = "dynamic";
      slot.slotGraphic.on('mouseup', () => {
        console.log('Mouse released on a slot: ' + slot.id);

        PlayCard(slot);
      });
    }
  }

  //Used to create the card from the received message from server.
  async function DrawCard(data:JSON) {

    //Load Cards Textures
    const texture = await PIXI.Assets.load(data.GraphicPath);
    const sprite = new PIXI.Sprite(texture);
    sprite.anchor.set(.5,.5);
    sprite.scale.set(cardSpriteScaler);
    sprite.cursor = 'pointer';
    app.stage.addChild(sprite);

    const markTex = await PIXI.Assets.load(data.MarkerPath);
    const markSprite = new PIXI.Sprite(markTex);
    markSprite.scale.set(0.3);

    //Adding card data.
    let name = data.Name;
    let description = data.Description;
    let selected = false;
    let graphicPath = data.GraphicPath;
    let markerSprite = markSprite;
    let targetX = 0;
    let targetY = 0;
    let x=0;
    let y = 0;




    const card:Card = {
         name,
         description,
         selected,
         graphicPath,
         markerSprite,
         sprite,
         targetX,
         targetY,
         x,
         y
    };

//  type Card struct {
// 	Type        string  //Card type (i.e. attack)
// 	Name        string  //Card name (should be unique for each card)
// 	Description string  //Card description
// 	Rarity      float64 //Card Rarity, all card rarities should add to 1.0
// 	GraphicPath string
// 	MarkerPath  string
// 	ImpactType  string      //Impact Type decides if many or singular slots are effected.
// 	ImpactShape string      //Impact Shape is the shape of the effect. (i.e. does it strike rows or a radius all around etc)
// 	MarkEffect  *MarkEffect //The effect the card has on the slots.
// }

    cardHand.push(card);

    DeselectAll();

    //Adding event listener for card select.
    sprite.eventMode = 'dynamic';
    sprite.on('mouseup', () => {
      console.log('Mouse down on a card');

      SelectCard(card);

    });

    sprite.on('mouseenter', () => {
      console.log('Mouse enter on a card');

      CardHoverEnter();

    });

    sprite.on('mouseexit', () => {
      console.log('Mouse enter on a card');

      CardHoverExit();

    });

    CentreHand();

  }

  //Used to Centre card hand.
  function CentreHand() {
   let cardCounter = 0;
    let cardHandLength = cardHand.length;

    const threshold = 3;
    const maxSpacing = window.innerWidth * 0.01;
    const minSpacing = window.innerWidth * 0.00001;
    const decayFactor = 0.6;

    let cardHandSpace;

    if (cardHandLength <= threshold) {
      cardHandSpace = maxSpacing;
    } else {
      const excess = cardHandLength - threshold;
      cardHandSpace = Math.max(minSpacing, maxSpacing * Math.pow(decayFactor, excess));
    }

    const cardWidth = cardHand[0]?.sprite.width || 100; // default if sprite width is unavailable

    // Total width from center to center
    const totalWidth = (cardWidth + cardHandSpace) * (cardHandLength - 1);

    // Starting X for the first card (anchor-centered)
    const startCardPosition = window.innerWidth / 2 - totalWidth / 2;

    for (const card of cardHand) {
      card.targetX = startCardPosition + cardCounter * (cardWidth + cardHandSpace);
      card.targetY = window.innerHeight / 2 + handHeight;
      cardCounter++;
    }

    
  }

  //Function to select card.
  function SelectCard(card:Card) {
    console.log("In Select: " + selectedCard);

    if(selectedCard !== undefined) {

      if(selectedCard == card) {
        DeselectCard(selectedCard);
        return;
      }


      DeselectCard(selectedCard);
    }


    card.selected = !card.selected;

    selectedCard = card;

    crdText.text = card.description; //displays text of card.
    crdText.style = style;
    crdText.scale.set(1);

    //app.stage.addChild(crdText);

    descBox = new PIXI.Graphics();
    descBox.lineStyle(2, 0x444444);
    
    descBox.beginFill(0xffffff);
    descBox.drawRect(0, 0, crdText.width + textPadding * 2, crdText.height + textPadding * 2);
    descBox.endFill();


    crdText.x = textPadding;
    crdText.y = textPadding;

    // descBox.lineStyle(5, 0xffffff); // border thickness and color
    // descBox.beginFill(0xe8e8e8);    // background color
    // descBox.drawRect(
    //     0, 0,
    //     crdText.width + textPadding * 2,
    //     crdText.height + textPadding * 2
    // );
    // descBox.endFill();

    descContainer = new PIXI.Container();
    descContainer.addChild(descBox);
    descContainer.addChild(crdText);

    descContainer.position.x = window.innerWidth/2 + 200;
    descContainer.position.y = window.innerHeight/2 - (descBox.height/2);

    // Add to stage
    app.stage.addChild(descContainer);

    if(card.selected == true) {
      //if already selected, move down.

      //card.sprite.position.y -= 50;
      card.targetY -= cardSelectRaise;
     
      card.sprite.scale.set(cardSpriteScaler*1.1);


    } else {

      //card.sprite.position.y += 50;
      card.targetY += cardSelectRaise;
      card.sprite.scale.set(cardSpriteScaler);
      
    }

    for (const card of cardHand) {
      console.log(card);
     }

  }

  //Function plays when mouse hovers over card.
  function CardHoverEnter() {

  }

  //Plays when mouse hover exits card.
  function CardHoverExit() {

  }

  //Deselect card logic.
  function DeselectCard(card:Card) {
    console.log("In deselect: " + selectedCard);

    //app.stage.removeChild(crdText);
    descBox.destroy();
    app.stage.removeChild(descContainer);
    descContainer.destroy();

    card.targetY += cardSelectRaise;
    card.sprite.scale.set(cardSpriteScaler);
    //card.sprite.position.y +=50;
    card.selected = false;
    selectedCard = undefined;

    for (const card of cardHand) {
      console.log(card);
     }

  }

  function DeselectAll() {
    for (const card of cardHand) {
      card.selected = false;
    }

    selectedCard = undefined;
  }

  //The play card logic.
  async function PlayCard(slot:Slot) {

    if(selectedCard === undefined) {
      return; //If no selected card, return.
    }

    send({ action: "play_card",type: "play_card", card_name:selectedCard.name,description: selectedCard.description,graphicPath:selectedCard.graphicPath,target_slot:slot.id}); //sending played card to server.


  }

  function PlayCardSuccess(data: JSON) {

    if(selectedCard === undefined) {
      return; //If no selected card, return.
    }

    if(selectedCard.markerSprite !== undefined) {

      let slotID = data.target_slot; //retreive target slot.

       if (slotID === undefined) {
        console.log("Err: SlotID undefined.")
        return;
      }

      let slot = board.slots[slotID]; //access display slot.

      if (slot === undefined) {
        return;
      }

      //This should be confirmed by server.

      slot.markerGraphic = selectedCard.markerSprite;

      //Adjusting zIndex
      slot.markerGraphic.zIndex = 1;

      //Centering slot marker based on width and slot size.
      slot.markerGraphic.position.x = slot.x + (slotSize*0.5) - slot.markerGraphic.width/2; 
      slot.markerGraphic.position.y = slot.y + (slotSize*0.5) - slot.markerGraphic.width/2;

      //Adding to scene.
      app.stage.addChild(slot.markerGraphic);


    }



    descBox.destroy(); 
    app.stage.removeChild(descContainer);
    descContainer.destroy();

    app.stage.removeChild(crdText);
   
    RemoveCard(selectedCard);//Removing card after it has been played.



  }

  async function AddMarkerGraphic(slot:Slot,mGraphicPath:string) { //Func used to add graphic to board using path.

    if (slot === undefined) {
      return;
    }

    //Creating sprite from graphic path.
    const markTex = await PIXI.Assets.load(mGraphicPath);
    const markSprite = new PIXI.Sprite(markTex);
    markSprite.scale.set(0.3);

    //This should be confirmed by server.

    slot.markerGraphic = markSprite;

    //Adjusting zIndex
    slot.markerGraphic.zIndex = 1;

    //Centering slot marker based on width and slot size.
    slot.markerGraphic.position.x = slot.x + (slotSize*0.5) - slot.markerGraphic.width/2; 
    slot.markerGraphic.position.y = slot.y + (slotSize*0.5) - slot.markerGraphic.width/2;

    //Adding to scene.
    app.stage.addChild(slot.markerGraphic);

  }

  //Remove card function.
  function RemoveCard(card:Card) {

    if (card === undefined) {
      return; //If the card is undefined, return.
    }

    let deleteCardIndex = cardHand.indexOf(card); //Finding index of card to delete.

    app.stage.removeChild(card.sprite);
    card.sprite.destroy(); //Destroy sprite.

    cardHand.splice(deleteCardIndex,1); //remove from cardHand array.

    selectedCard = undefined; //Changing selected to null.

    CentreHand(); //Cenre hand after deletion.

  }

  //Lerp function.
  function lerp(a:number,b:number,t:number) {
    return a + (b - a) * t;
  }

  function StartGame(data:JSON) {

   // console.log(data.cards_to_add);


    for (let i=0;i<data.cards_to_add.length;i++) { //Drawing starting cards.
      //console.log(data.cards_to_add[i].GraphicPath);
      DrawCard(data.cards_to_add[i]);
    }


  }

  function UpdateBoard(data:JSON) {

    let slotsToUpdate = data.board_state;

    if (slotsToUpdate === undefined) { //If no slotsToUpdate, return.
      return;
    }

    for (let i = 0; i < board.slots.length; i++) {
      const clSlot = board.slots[i];

      for (let z = 0; z < slotsToUpdate.length; z++) {
        const sSlot = slotsToUpdate[z];

        if (clSlot.id == sSlot.ID) {

          if (sSlot.Effects == null) { //skips slot if effects is null.
            console.log("Skipping slot: Update empty.");
            continue;
          }

          

           for (let mEffId = (sSlot.Effects.length - 1); mEffId >= 0; mEffId--) {
            if (sSlot.Effects[mEffId].IsDisplayable == true) {
              console.log("Updating Slot Graphic");
              AddMarkerGraphic(clSlot,sSlot.Effects[mEffId].GraphicPath);

              //clSlot.markerGraphic = sSlot.Effects[mEffId].GraphicPath;
              break;
            }
           }
           
        }
      }
    }
  }




  // ------------------ EVENT LISTENERS ------------------------


  window.addEventListener('resize', () => {
    ScaleSize();
    CentreBoard();
  
    DeselectAll();

    CentreHand();
    SetSlotListeners();
    DeselectAll();

  });


  // Movement controls
  window.addEventListener("keydown", (e) => {
    const speed = 10;
    switch (e.key) {
      case "ArrowUp":
        //send({ type: "draw_card", cardName:"mark",description: "Place a mark.",graphicPath:"src/card_test_design2.png",markerPath:"src/naught.svg"});
      
        break;
      case "ArrowDown":
       send({ type: "draw_card", cardName:"remove",description: "Remove a random opponent mark.",graphicPath:"src/card_ttt_test3.png",markerPath:"src/cross.svg"});
      
        break;
      case "ArrowLeft":
       
        break;
      case "ArrowRight":
        
        break;
    }

   
  });

  eventBus.addEventListener("wsMessage", (event: Event) => {
    const customEvent = event as CustomEvent;
    const data = customEvent.detail;

    const jsonData = JSON.parse(data); //Reading JSON message.

    

    switch (jsonData.type) { //Determining message type and how to react.
      case "game_start":
        console.log(jsonData);
        StartGame(jsonData);
        break;
      case "draw_card":
        console.log("Received Draw Card Message.");
        DrawCard(jsonData);
        break;
      case "turn_start":
        break;
      case "play_card_success":
        PlayCardSuccess(jsonData);
        break;
      case "game_state":
        UpdateBoard(jsonData);
        break;
      

    }
  });

  

  app.ticker.add((ticker) =>
    {
        for (const card of cardHand) {
        // Animating Card Position changes.
          if(Math.abs(card.targetX - card.sprite.x) < 0.5 && Math.abs(card.targetY - card.sprite.y) < 0.5 ) {
            //Jumps to position to stop floating point lerp issues.
            card.sprite.x = card.targetX;
            card.sprite.y = card.targetY;
          } else {
            //Lerping position.
            card.sprite.x = lerp(card.sprite.x, card.targetX, 0.1); 
            card.sprite.y = lerp(card.sprite.y, card.targetY, 0.1);
          }
         
          // // Rotation
          // card.sprite.rotation = lerp(card.sprite.rotation, card.targetRotation, 0.1);
        }
        //console.log("Ticker val: " + app.ticker.deltaTime);
    });

}
)()





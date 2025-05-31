import * as PIXI from "pixi.js";
import { send } from "./client";
import eventBus from "./client";



(async () => {


  const app = new PIXI.Application<HTMLCanvasElement>({ background: 0xffffff, resizeTo: window });

  document.body.appendChild(app.view);
  const container = new PIXI.Container();

  app.stage.addChild(container);

  ///------------------------ INIT VAR ----------------------

  //Creating the tic-tac-toe board.
  //The slot interface ensures we have all data for each slot.
  interface Slot {
    id: number;
    x:number;
    y:number;
    row: number;
    col: number;
    colour: number;
    slotGraphic: PIXI.Graphics;
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
    sprite: PIXI.Sprite,
    targetX:number,
    targetY:number,
    x:number,
    y:number
  }

  const cardHand: Card[] = [];
  let selectedCard: Card | undefined;

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

    for (const slot of board.slots) {
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

        const slotGraphic = new PIXI.Graphics();

        slotGraphic.beginFill(colour);
        slotGraphic.drawRect(
          x,
          y,
          slotSize,
            slotSize
        );
        slotGraphic.endFill();

        app.stage.addChild(slotGraphic);

        const slot:Slot = {
          id,
          x,
          y,
          row,
          col,
          colour,
          slotGraphic
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
        cardSpriteScaler = 0.2;
         cardHandSpace =  window.innerWidth * 0.01;

      } else if (window.innerWidth >= 255 && window.innerWidth < 370)
      {
        slotSize = window.innerWidth*0.12;
        cardSpriteScaler = 0.24;
        cardHandSpace =  window.innerWidth * 0.01;
      }
      else if (window.innerWidth >= 370 && window.innerWidth < 512)
      {
        slotSize = window.innerWidth*0.1;
        cardSpriteScaler = 0.3;
        cardHandSpace =  window.innerWidth * 0.01;
      }
      else if (window.innerWidth >= 512 && window.innerWidth < 900)
      {
        slotSize = window.innerWidth*0.075;
        cardSpriteScaler = 0.32;
        cardHandSpace =  window.innerWidth * 0.01;
      }
      else if (window.innerWidth >= 900 && window.innerWidth < 1000)
      {
        slotSize = window.innerWidth*0.07;
        cardSpriteScaler = 0.35;
        cardHandSpace =  window.innerWidth * 0.01;
      }
      else if (window.innerWidth >= 1000 && window.innerWidth < 1200)
      {
        slotSize = window.innerWidth*0.06;
        cardSpriteScaler = 0.35;
        cardHandSpace =  window.innerWidth * 0.01;
      }
      else 
      {
        slotSize = window.innerWidth*0.05;
        cardSpriteScaler = 0.35;
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

    //Load Cards.
    const texture = await PIXI.Assets.load(data.graphicPath);
    const sprite = new PIXI.Sprite(texture);
    sprite.scale.set(cardSpriteScaler);
    app.stage.addChild(sprite);


    //Adding card data.
    let name = "must add card name.";
    let description = "must add card desc.";
    let selected = false;
    let graphicPath = data.graphicPath;
    let targetX = 0;
    let targetY = 0;
    let x=0;
    let y = 0;

    const card:Card = {
         name,
         description,
         selected,
         graphicPath,
         sprite,
         targetX,
         targetY,
         x,
         y
    };

    cardHand.push(card);

    //Adding event listener for card select.
    sprite.eventMode = 'dynamic';
    sprite.on('mousedown', () => {
      console.log('Mouse released on a slot: ' + sprite);

      SelectCard(card);

    });

    CentreHand();

  }

  //Used to Centre card hand.
  function CentreHand() {
    let cardCounter = 0;
    let cardLength = cardHand.length;
    let startCardPosition = 0;

    const overlapFactor = 0.75;
    const baseSpacing = window.innerWidth * 0.01; // or tweak as needed
    cardHandSpace = baseSpacing / (1 + Math.log(cardLength));

    const maxSpacing = window.innerWidth * 0.01; // px spacing for small hands
    const minSpacing = window.innerWidth * 0.00001; // minimum spacing allowed when overlapping
    const threshold = 3;   // number of cards before overlap starts

    if (cardLength <= threshold) {
      cardHandSpace = maxSpacing;
    } else {
      // Smooth decay only after threshold
      const excess = cardLength - threshold;
      const decayFactor = 0.6; // lower = more overlap per extra card
      cardHandSpace = Math.max(
        minSpacing,
        maxSpacing * Math.pow(decayFactor, excess)
      );
    }

    for (const card of cardHand) {
      card.sprite.scale.set(cardSpriteScaler);

      if (cardCounter==0) {
        startCardPosition = window.innerWidth/2 - (((card.sprite.width + cardHandSpace) * cardLength)/2);
      }

      //card.sprite.position.set(startCardPosition + (cardCounter*(card.sprite.width+cardHandSpace)),window.innerHeight/2 + 250);

      card.targetX = startCardPosition + (cardCounter*(card.sprite.width+cardHandSpace));
      card.targetY = window.innerHeight/2 + 200;

      cardCounter++;
    }
  }

  //Function to select card.
  function SelectCard(card:Card) {

    if(selectedCard !== undefined) {

      if(selectedCard == card) {
        DeselectCard(selectedCard);
        return;
      }


      DeselectCard(selectedCard);
    }


    card.selected = !card.selected;

    selectedCard = card;

    if(card.selected == true) {
      //if already selected, move down.

      //card.sprite.position.y -= 50;
      card.targetY -= 50;


    } else {

      //card.sprite.position.y += 50;
      card.targetY += 50;
    }

  }

  function DeselectCard(card:Card) {
    
    card.selected = false;
    //card.sprite.position.y += 50;
    card.targetY += 50;
    selectedCard = undefined;

  }

  //The play card logic, where all cards are played.
  function PlayCard(slot:Slot) {

    if(selectedCard === undefined) {
      return;
    }



    console.log("Card Played.");

  }

  function lerp(a:number,b:number,t:number) {
    return a + (b - a) * t;
  }



  // ------------------ EVENT LISTENERS ------------------------


  window.addEventListener('resize', () => {
    ScaleSize();
    CentreBoard();
    CentreHand();
    SetSlotListeners();

  });


  // Movement controls
  window.addEventListener("keydown", (e) => {
    const speed = 10;
    switch (e.key) {
      case "ArrowUp":
        send({ type: "draw_card", cardName:"mark",graphicPath:"src/card_test.png"});
       // player.y -= speed;
        break;
      case "ArrowDown":
      //  player.y += speed;
        break;
      case "ArrowLeft":
       // player.x -= speed;
        break;
      case "ArrowRight":
        //player.x += speed;
        break;
    }

   
  });

  eventBus.addEventListener("wsMessage", (event: Event) => {
    const customEvent = event as CustomEvent;
    const data = customEvent.detail;

    const jsonData = JSON.parse(data); //Reading JSON message.

    switch (jsonData.type) { //Determining message type and how to react.
      case "draw_card":
        console.log("Received Draw Card Message.");
        DrawCard(jsonData);
        break;

    }
  });

  

  app.ticker.add((ticker) =>
    {
        for (const card of cardHand) {
        // Position
        if(Math.abs(card.targetX - card.sprite.x) < 0.5 && Math.abs(card.targetY - card.sprite.y) < 0.5 ) {
          card.sprite.x = card.targetX;
          card.sprite.y = card.targetY;
        } else {
          card.sprite.x = lerp(card.sprite.x, card.targetX, 0.1); 
          card.sprite.y = lerp(card.sprite.y, card.targetY, 0.1);
        }
         

          

          // // Rotation
          // card.sprite.rotation = lerp(card.sprite.rotation, card.targetRotation, 0.1);
        }
        //console.log("Ticker val: " + app.ticker.deltaTime);
    });

})()


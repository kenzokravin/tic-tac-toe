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
    graphicPath: string,
    sprite: PIXI.Sprite
  }

  const cardHand: Card[] = [];

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

    //Adding event listener for drag.
    sprite.eventMode = 'dynamic';
    sprite.on('mousedown', () => {
      console.log('Mouse released on a slot: ' + sprite);
    });

    //Adding card data to hand.

    let name = "must add card name.";
    let graphicPath = data.graphicPath;

    const card:Card = {
         name,
         graphicPath,
         sprite
    };

    cardHand.push(card);

    CentreHand();

  }

  function CentreHand() {
    let cardCounter = 0;
    let cardLength = cardHand.length;
    let startCardPosition = 0;

    for (const card of cardHand) {
      //app.stage.removeChild(card.sprite);
      //card.sprite.destroy(); // destroy graphics to resize.
      card.sprite.scale.set(cardSpriteScaler);

      if (cardCounter==0) {
        startCardPosition = window.innerWidth/2 - (((card.sprite.width + cardHandSpace) * cardLength)/2);
      }

      card.sprite.position.set(startCardPosition + (cardCounter*(card.sprite.width+cardHandSpace)),window.innerHeight/2 + 250);



      cardCounter++;
    }
    

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
        
        //console.log("Ticker val: " + app.ticker.deltaTime);
    });

})()


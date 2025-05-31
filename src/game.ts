import * as PIXI from "pixi.js";
import { send } from "./client";
import eventBus from "./client";



(async () => {
  const app = new PIXI.Application<HTMLCanvasElement>({ background: 0xffffff, resizeTo: window });

  document.body.appendChild(app.view);
  const container = new PIXI.Container();

  app.stage.addChild(container);

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


  ScaleSlotSize();

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

  CentreBoard();
  SetSlotListeners();

  //Function used to centre the board in space.
  function CentreBoard() {
    board.x = window.innerWidth/2;
    ScaleSlotSize();
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

  //Function used to dynamically resize the playing board.
  function ScaleSlotSize() {
    if(window.innerWidth < 255) {
        slotSize = window.innerWidth*0.14;

      } else if (window.innerWidth >= 255 && window.innerWidth < 370)
      {
        slotSize = window.innerWidth*0.12;
      }
      else if (window.innerWidth >= 370 && window.innerWidth < 512)
      {
        slotSize = window.innerWidth*0.1;
      }
      else if (window.innerWidth >= 512 && window.innerWidth < 900)
      {
        slotSize = window.innerWidth*0.075;
      }
      else if (window.innerWidth >= 900 && window.innerWidth < 1000)
      {
        slotSize = window.innerWidth*0.07;
      }
      else if (window.innerWidth >= 1000 && window.innerWidth < 1200)
      {
        slotSize = window.innerWidth*0.06;
      }
      else 
      {
        slotSize = window.innerWidth*0.05;
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
  async function DrawCard(path:string) {

    const texture = await PIXI.Assets.load(path);
    const sprite = new PIXI.Sprite(texture);
    app.stage.addChild(sprite);
  }



  // ------------------ EVENT LISTENERS ------------------------


  window.addEventListener('resize', () => {

    CentreBoard();
    SetSlotListeners();

  });


  // Movement controls
  window.addEventListener("keydown", (e) => {
    const speed = 10;
    switch (e.key) {
      case "ArrowUp":
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

    send({ type: "draw_card", cardName:"mark",graphicPath:"src/card_test.png"});
  });

  eventBus.addEventListener("wsMessage", (event: Event) => {
    const customEvent = event as CustomEvent;
    const data = customEvent.detail;

    console.log("Game received:", data);
    
  });


  app.ticker.add((ticker) =>
    {
        
        //console.log("Ticker val: " + app.ticker.deltaTime);
    });

})()


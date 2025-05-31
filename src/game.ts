import * as PIXI from "pixi.js";
import { send } from "./client";

const app = new PIXI.Application({
  resizeTo: window,
  backgroundColor: 0xffffff,
});
document.body.appendChild(app.view);

// Player circle
const player = new PIXI.Graphics();
player.beginFill(0x00ff00);
player.drawCircle(0, 0, 20);
player.endFill();
player.x = window.innerWidth / 2;
player.y = window.innerHeight / 2;
//app.stage.addChild(player);


//Creating the tic-tac-toe board.
let slotSize = window.innerWidth*0.04;

ScaleSlotSize();

let board: {
  slots: PIXI.Graphics[],
  x: number,
  y: number
} = {
  slots: [],
  x: window.innerWidth/2,
  y: window.innerHeight/2
};

let slotCounter = 0;

for (let i=0;i < 3;i++) {
  for (let z = 0; z <3;z++) {
    let slot = new PIXI.Graphics();

    if (slotCounter%2 ==0) {
      slot.beginFill(0xd3d3d3);
    } else {
      slot.beginFill(0x888888);
    }

    slot.drawRect((board.x - (1.5*slotSize)) + slotSize*z ,board.y- (1.5*slotSize) +(i*slotSize) ,slotSize,slotSize);
    

    slot.endFill();

    app.stage.addChild(slot);

    slotCounter++;
    board.slots.push(slot);
    
  }
}




function CentreBoard() {
  board.x = window.innerWidth/2;

  console.log("Width: " + window.innerWidth);

  ScaleSlotSize();

  board.y = window.innerHeight/2;

  slotCounter = 0;

  for (const slot of board.slots) {
    app.stage.removeChild(slot);
    slot.destroy(); // destroy graphics to resize.
  }
  board.slots = []; // Reset array

  
      for (let i = 0; i < 3; i++) {
    for (let z = 0; z < 3; z++) {
      const slot = new PIXI.Graphics();

      slot.beginFill(slotCounter % 2 === 0 ? 0xd3d3d3 : 0x888888);
      slot.drawRect(
        (board.x - 1.5 * slotSize) + slotSize * z,
        (board.y - 1.5 * slotSize) + slotSize * i,
        slotSize,
        slotSize
      );
      slot.endFill();

      app.stage.addChild(slot);
      board.slots.push(slot);
      slotCounter++;
    }
  }
}

function ScaleSlotSize() {
  if(window.innerWidth < 255) {
      slotSize = window.innerWidth*0.12;

    } else if (window.innerWidth >= 255 && window.innerWidth < 370)
    {
      slotSize = window.innerWidth*0.1;
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
      slotSize = window.innerWidth*0.04;
    }

}



// ------------------ EVENT LISTENERS ------------------------


window.addEventListener('resize', () => {

  CentreBoard();

});


// Movement controls
window.addEventListener("keydown", (e) => {
  const speed = 10;
  switch (e.key) {
    case "ArrowUp":
      player.y -= speed;
      break;
    case "ArrowDown":
      player.y += speed;
      break;
    case "ArrowLeft":
      player.x -= speed;
      break;
    case "ArrowRight":
      player.x += speed;
      break;
  }

  send({ type: "move", x: player.x, y: player.y });
});

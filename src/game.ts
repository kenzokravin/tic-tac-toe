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
let slotSize = 50;
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

    slot.drawRect((board.x - (1.5*slotSize)) + slotSize*z,board.y- (1.5*slotSize) +(i*slotSize),slotSize,slotSize);

    slot.endFill();

    app.stage.addChild(slot);

    slotCounter++;
    board.slots.push(slot);
    
  }
}






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

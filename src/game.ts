import * as PIXI from "pixi.js";
import { send } from "./client";

const app = new PIXI.Application({
  resizeTo: window,
  backgroundColor: 0x222222,
});
document.body.appendChild(app.view);

// Player circle
const player = new PIXI.Graphics();
player.beginFill(0x00ff00);
player.drawCircle(0, 0, 20);
player.endFill();
player.x = window.innerWidth / 2;
player.y = window.innerHeight / 2;
app.stage.addChild(player);

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

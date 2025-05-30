const socket = new WebSocket("ws://localhost:8080/ws");

socket.addEventListener("open", () => {
  console.log("✅ WebSocket connected");
});

socket.addEventListener("message", (event) => {
  console.log("📨 Server:", event.data);
});

socket.addEventListener("close", () => {
  console.warn("🔌 WebSocket disconnected");
});

socket.addEventListener("error", (err) => {
  console.error("❌ WebSocket error:", err);
});

export function send(data: object) {
  if (socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(data));
  }
}

export { socket };

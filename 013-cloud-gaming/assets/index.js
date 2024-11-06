let peerConn;
let signalCh;
let dataCh;
let videoE;

window.onload = async function () {
  const startBtn = document.getElementById("snake_game_btn");
  startBtn.addEventListener("click", start);

  videoE = document.getElementById("video");
};

async function start() {
  signalCh = new WebSocket("ws://localhost:8080/ws");

  signalCh.addEventListener("open", async () => {
    peerConn = new RTCPeerConnection();
  });

  signalCh.addEventListener("message", async (event) => {
  });
}

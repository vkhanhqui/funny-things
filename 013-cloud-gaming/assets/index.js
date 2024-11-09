let peerConn;
let signalCh;
let videoE;

window.onload = async function () {
  const startBtn = document.getElementById("snake_game_btn");
  startBtn.addEventListener("click", start);
  videoE = document.getElementById("video");
};

async function start() {
  signalCh = new WebSocket("ws://localhost:8080/ws");
  signalCh.addEventListener("open", onOpen);
  signalCh.addEventListener("message", onMessage);
}

// region: WebRTC support functions

const onOpen = async () => {
  peerConn = new RTCPeerConnection();
  peerConn.onicecandidate = handleIceCandidateEvent;

  const offer = await peerConn.createOffer({
    offerToReceiveVideo: true,
  });
  await peerConn.setLocalDescription(offer);
  signalCh.send(JSON.stringify({ type: "offer", value: offer.sdp }));
};

const onMessage = async (event) => {
  const message = JSON.parse(event.data);
  if (message.type === "ice") {
    handleIceMessage(message.value);
  }
  if (message.type === "answer") {
    handleAnswerMessage(message.value);
  }
};

const handleIceCandidateEvent = (event) => {
  if (event.candidate) {
    signalCh.send(
      JSON.stringify({ type: "ice", value: event.candidate.candidate })
    );
  }
};

const handleIceMessage = async (iceValue) => {
  const iceCandidate = new RTCIceCandidate({
    ...iceValue,
    sdpMLineIndex: 0,
    sdpMid: "0",
  });
  try {
    await peerConn.addIceCandidate(iceCandidate);
  } catch (error) {
    console.log("Error adding ICE candidate:", error);
  }
};

const handleAnswerMessage = async (answerValue) => {
  const remoteDescription = new RTCSessionDescription({
    sdp: answerValue,
    type: "answer",
  });
  await peerConn.setRemoteDescription(remoteDescription);
};

let localStream;
let peerConnections = {}; // id အလိုက် peer တည်ဆောက်မှုများကို သိမ်းဆည်းရန်
let wsSocket, userRole;
let isScreenSharing = false;

const iceServersConfig = { iceServers: [{ urls: ["stun:://google.com"] }] };

// ⚙️ assets/main.js ရဲ့ အပေါ်ပိုင်း initializeLiveClassroom() ဖန်ရှင်တွင် ဤသို့ ပြင်ပါ [source: 2]

function initializeLiveClassroom() {
    userRole = document.getElementById("roleInput").value;

    document.getElementById("entrance-gate").classList.add("d-none");
    document.getElementById("classroom-view").classList.remove("d-none");
    document.getElementById("live-controls-toolbar").classList.remove("d-none");

    document.getElementById("localUserNameTag").innerText = userRole === "teacher" ? "👨‍🏫 Teacher (Host)" : "🎓 Student (Peer)";

    // 🌟 ဂျီမကျစေရန် ကျောင်းသားဖြစ်ပါက Private Chat ရော Hand Raising ခလုတ်ပါ တစ်ခါတည်း ဖွင့်ပေးခြင်း [source: 2]
    if (userRole === "student") {
        document.getElementById("private-student-panel").classList.remove("d-none");
        document.getElementById("btn-hand").classList.remove("d-none");
    }

    setupWebSocket();
}

function setupWebSocket() {
    wsSocket = new WebSocket(`ws://${window.location.host}/ws/join?role=${userRole}`);

    wsSocket.onopen = async () => {
        try {
            localStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
            document.getElementById("localVideo").srcObject = localStream;
        } catch (e) {
            console.warn("Camera hardware bypass mode activate.");
            localStream = new MediaStream();
        }
        // အခန်းထဲရောက်ကြောင်း အားလုံးဆီ သတင်းပို့ခြင်း
        sendSignalingMessage({ type: "ready-to-stream" });
    };

    // 📄 assets/main.js ၏ wsSocket.onmessage တစ်ခုလုံးအား ဤကုဒ်ဖြင့် လဲလှယ်ပါ [source: 1.2.7]

    wsSocket.onmessage = async (event) => {
        const message = JSON.parse(event.data);

        if (message.type === "student-count-update") {
            document.getElementById("studentCountDisplay").innerText = `👥 Online Students: ${message.count}`;
            return;
        }

        const senderId = message.senderId;

        if (message.type === "ready-to-stream") {
            await initiateWebRTCCall(senderId, message.senderRole);
        } else if (message.type === "webrtc-offer") {
            await handleIncomingOffer(senderId, message.payload, message.senderRole);
        } else if (message.type === "webrtc-answer") {
            await handleIncomingAnswer(senderId, message.payload);
        } else if (message.type === "ice-candidate") {
            await handleIncomingCandidate(senderId, message.payload);
        } else if (message.type === "chat") {
            appendChatMessageToBox(message.senderName, message.content);
        } else if (message.type === "private-chat" && userRole === "student") {
            appendPrivateChatMessageToBox(message.senderName, message.content);
        } else if (message.type === "peer-disconnected") {
            removePeerVideoFrame(senderId);

            // 🌟 [အသစ်ထပ်တိုး] ဆရာမဘက်က Host ပိတ်လိုက်ပါက ကျောင်းသားအားလုံးကို အလိုအလျောက် မောင်းထုတ်ခြင်း [source: 1.1.4]
        } else if (message.type === "terminate-class") {
            alert("🚨 ဆရာမမှ အတန်းအား အပြီးသတ် ပိတ်သိမ်းလိုက်ပါပြီဗျာ။ အပြင်သို့ ပြန်လည် ပို့ဆောင်ပါမည်။");
            window.location.reload(); // Page ကို Reload ချပြီး Login Gate သို့ ပြန်ပို့ခြင်း
        }
    };

}

// ✋ လက်ထောင်လိုက်ကြောင်း ဆာဗာမှတစ်ဆင့် Broadcast ပို့ခြင်း [source: 2]
function signalRaiseHandAction() {
    // 🌟 ဆာဗာမှတစ်ဆင့် အခန်းထဲရှိ လူအားလုံးဆီသို့ လက်ထောင်ကြောင်း Broadcast ပို့ခိုင်းခြင်း [source: 1.2.7]
    sendSignalingMessage({
        type: "raise-hand"
    });

    // ကိုယ့် Chat Box ထဲတွင်လည်း အသိပေးချက် ပြသရန် [source: 1.1.4]
    appendChatMessageToBox("Me", "You raised your hand ✋");
}

function createRTCObject(peerId, peerRole) {
    const pc = new RTCPeerConnection(iceServersConfig);
    peerConnections[peerId] = pc;

    if (localStream) {
        localStream.getTracks().forEach(track => pc.addTrack(track, localStream));
    }

    pc.ontrack = (event) => {
        let videoEl = document.getElementById(`video-${peerId}`);
        if (!videoEl) {
            const grid = document.getElementById("video-grid-matrix");
            const card = document.createElement("div");
            card.id = `card-${peerId}`;
            card.className = "video-card-node";

            videoEl = document.createElement("video");
            videoEl.id = `video-${peerId}`;
            videoEl.autoplay = true;
            videoEl.playsInline = true;

            const label = document.createElement("div");
            label.className = "video-label";
            label.innerText = peerRole === "teacher" ? "👨‍🏫 Teacher" : "🎓 Student";

            card.appendChild(videoEl);
            card.appendChild(label);
            grid.appendChild(card);
        }
        videoEl.srcObject = event.streams[0];
    };

    pc.onicecandidate = (event) => {
        if (event.candidate) {
            sendSignalingMessage({ type: "ice-candidate", targetId: peerId, payload: event.candidate });
        }
    };

    return pc;
}

async function initiateWebRTCCall(peerId, peerRole) {
    const pc = createRTCObject(peerId, peerRole);
    const offer = await pc.createOffer();
    await pc.setLocalDescription(offer);
    sendSignalingMessage({ type: "webrtc-offer", targetId: peerId, payload: offer });
}

async function handleIncomingOffer(senderId, offer, senderRole) {
    const pc = createRTCObject(senderId, senderRole);
    await pc.setRemoteDescription(new RTCSessionDescription(offer));
    const answer = await pc.createAnswer();
    await pc.setLocalDescription(answer);
    sendSignalingMessage({ type: "webrtc-answer", targetId: senderId, payload: answer });
}

async function handleIncomingAnswer(senderId, answer) {
    const pc = peerConnections[senderId];
    if (pc) await pc.setRemoteDescription(new RTCSessionDescription(answer));
}

async function handleIncomingCandidate(senderId, candidate) {
    const pc = peerConnections[senderId];
    if (pc) await pc.addIceCandidate(new RTCIceCandidate(candidate));
}

function removePeerVideoFrame(peerId) {
    const card = document.getElementById(`card-${peerId}`);
    if (card) card.remove();
    if (peerConnections[peerId]) {
        peerConnections[peerId].close();
        delete peerConnections[peerId];
    }
}

function sendSignalingMessage(data) {
    if (wsSocket && wsSocket.readyState === WebSocket.OPEN) {
        wsSocket.send(JSON.stringify(data));
    }
}

// 🖥️ SCREEN SHARING SYSTEM
async function toggleScreenSharing() {
    try {
        const shareBtn = document.getElementById("btn-share");
        if (!isScreenSharing) {
            const screenStream = await navigator.mediaDevices.getDisplayMedia({ video: true });
            const screenTrack = screenStream.getVideoTracks()[0];

            for (let id in peerConnections) {
                const sender = peerConnections[id].getSenders().find(s => s.track && s.track.kind === "video");
                if (sender) await sender.replaceTrack(screenTrack);
            }

            document.getElementById("localVideo").srcObject = screenStream;
            screenTrack.onended = () => { revertToCameraStream(); };

            shareBtn.innerHTML = `<i class="las la-stop-circle"></i>`;
            shareBtn.className = "btn btn-action-node btn-danger text-white";
            isScreenSharing = true;
        } else {
            revertToCameraStream();
        }
    } catch (err) { console.error(err); }
}

async function revertToCameraStream() {
    const cameraTrack = localStream.getVideoTracks()[0];
    for (let id in peerConnections) {
        const sender = peerConnections[id].getSenders().find(s => s.track && s.track.kind === "video");
        if (sender) await sender.replaceTrack(cameraTrack);
    }
    document.getElementById("localVideo").srcObject = localStream;
    document.getElementById("btn-share").innerHTML = `<i class="las la-desktop"></i>`;
    document.getElementById("btn-share").className = "btn btn-action-node btn-warning text-white";
    isScreenSharing = false;
}

// CHAT FUNCTIONS
function transmitLiveChatMessage() {
    const input = document.getElementById("chatInput");
    const text = input.value.trim();
    if (!text) return;
    appendChatMessageToBox("Me", text);
    sendSignalingMessage({ type: "chat", content: text });
    input.value = "";
}

function appendChatMessageToBox(sender, content) {
    const container = document.getElementById("chat-messages");
    const div = document.createElement("div");
    div.className = sender === "Me" ? "chat-bubble-node bubble-me" : "chat-bubble-node bubble-remote";
    div.innerHTML = `<small class="d-block opacity-75 fw-bold mb-1">${sender}</small>${content}`;
    container.appendChild(div);
    container.scrollTop = container.scrollHeight;
}

function transmitPrivateChatMessage() {
    const input = document.getElementById("privateChatInput");
    const text = input.value.trim();
    if (!text) return;
    appendPrivateChatMessageToBox("Me", text);
    sendSignalingMessage({ type: "private-chat", content: text });
    input.value = "";
}

function sendPrivateEmoji(emojiIcon) {
    appendPrivateChatMessageToBox("Me", emojiIcon);
    sendSignalingMessage({ type: "private-chat", content: emojiIcon });
}

function appendPrivateChatMessageToBox(sender, content) {
    const container = document.getElementById("private-chat-messages");
    if (!container) return;
    const div = document.createElement("div");
    div.className = sender === "Me" ? "chat-bubble-node bubble-me" : "chat-bubble-node bubble-remote";
    div.innerHTML = `<small class="d-block text-success fw-bold mb-1">${sender}</small>${content}`;
    container.appendChild(div);
    container.scrollTop = container.scrollHeight;
}

function toggleAudioMicrophone() {
    const track = localStream.getAudioTracks()[0];
    if (track) {
        track.enabled = !track.enabled;
        document.getElementById("btn-mic").className = track.enabled ? "btn btn-action-node btn-light text-primary" : "btn btn-action-node btn-danger text-white";
    }
}
function toggleVideoCamera() {
    const track = localStream.getVideoTracks()[0];
    if (track) {
        track.enabled = !track.enabled;
        document.getElementById("btn-video").className = track.enabled ? "btn btn-action-node btn-light text-primary" : "btn btn-action-node btn-danger text-white";
    }
}

// 📄 assets/main.js ၏ အောက်ခြေဆုံးတွင် အောက်ပါ ကုဒ်အသစ်များကို ထပ်ဖြည့်ပေးပါ [source: 1.4.5]

// 🌟 [အသစ်ထပ်တိုး] ဆရာမဘက်က Host အပြီးပိတ်သိမ်းမည့် ခလုတ်စနစ် [source: 1.1.4]
function terminateClassroomByHost() {
    if (userRole === "teacher") {
        if (confirm("⚠️ သင်သည် ဤအတန်းကို အပြီးပိတ်ပြီး ကျောင်းသားအားလုံးအား မောင်းထုတ်လိုပါသလား?")) {
            sendSignalingMessage({ type: "terminate-class" });
            window.location.reload();
        }
    }
}

// 🌟 [အသစ်ထပ်တိုး] Message input ကွက်များတွင် Enter ခေါက်ရုံဖြင့် dynamic ပို့ဆောင်နိုင်မည့် စနစ် [source: 1.4.5, 1.4.6]
document.addEventListener("DOMContentLoaded", () => {
    const chatInput = document.getElementById("chatInput");
    const privateChatInput = document.getElementById("privateChatInput");

    if (chatInput) {
        chatInput.addEventListener("keydown", (e) => {
            if (e.key === "Enter" && !e.shiftKey) {
                e.preventDefault(); // browser form submit ဖြစ်ခြင်းကို တားဆီးရန် [source: 1.4.5]
                transmitLiveChatMessage();
            }
        });
    }

    if (privateChatInput) {
        privateChatInput.addEventListener("keydown", (e) => {
            if (e.key === "Enter" && !e.shiftKey) {
                e.preventDefault();
                transmitPrivateChatMessage();
            }
        });
    }
});


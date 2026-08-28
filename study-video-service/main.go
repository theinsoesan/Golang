package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Room string
	Role string
}

type RoomHub struct {
	Mutex sync.RWMutex
	Rooms map[string]map[*Client]bool
}

var hub = RoomHub{Rooms: make(map[string]map[*Client]bool)}
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func main() {
	r := gin.Default()
	r.Static("/assets", "./assets")
	r.LoadHTMLFiles("./templates/video-call.html")

	r.GET("/", func(c *gin.Context) { c.HTML(200, "video-call.html", nil) })
	r.GET("/ws/join", handleWebSocketSignaling)

	log.Println("🚀 Classroom Server running on http://127.0.0.1:5000")
	r.Run(":5000")
}

func handleWebSocketSignaling(c *gin.Context) {
	roomCode := "ROOM-XYZ"
	role := c.Query("role")
	if role == "" {
		role = "student"
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	peerID := fmt.Sprintf("peer-%d", rand.Intn(100000))
	client := &Client{ID: peerID, Conn: conn, Room: roomCode, Role: role}

	hub.Mutex.Lock()
	if hub.Rooms[roomCode] == nil {
		hub.Rooms[roomCode] = make(map[*Client]bool)
	}
	hub.Rooms[roomCode][client] = true
	hub.Mutex.Unlock()

	broadcastStudentCount(roomCode)

	defer func() {
		hub.Mutex.Lock()
		if hub.Rooms[roomCode] != nil {
			delete(hub.Rooms[roomCode], client)
		}
		hub.Mutex.Unlock()

		// လူထွက်သွားလျှင် အခြားသူများဆီမှ ဗီဒီယိုကွက်ကို ဖျက်ခိုင်းခြင်း
		disconnectMsg := map[string]interface{}{"type": "peer-disconnected", "senderId": peerID}
		broadcastToRoom(roomCode, client, disconnectMsg)

		broadcastStudentCount(roomCode)
		conn.Close()
	}()

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		msg["senderId"] = peerID
		msg["senderRole"] = role
		msg["senderName"] = fmt.Sprintf("%s (%s)", role, peerID)

		msgType, _ := msg["type"].(string)
		targetID, _ := msg["targetId"].(string)

		hub.Mutex.RLock()
		clientsInRoom := hub.Rooms[roomCode]

		for peer := range clientsInRoom {
			if msgType == "private-chat" {
				if peer.Role == "student" && peer != client {
					_ = peer.Conn.WriteJSON(msg)
				}
				continue
			}

			if msgType == "terminate-class" {
				if peer != client {
					_ = peer.Conn.WriteJSON(msg)
				}
				continue
			}

			// 🌟 [အသစ်ထပ်တိုး] ကျောင်းသား လက်ထောင်ပါက အခန်းထဲရှိ အခြားသူအားလုံး (ဆရာမ အပါအဝင်) ထံသို့ Broadcast တွန်းပို့ခြင်း [source: 1.3.5]
			if msgType == "raise-hand" {
				if peer != client {
					_ = peer.Conn.WriteJSON(msg)
				}
				continue
			}

			// WebRTC Target Specific Signaling Routing
			if targetID != "" {
				if peer.ID == targetID {
					_ = peer.Conn.WriteJSON(msg)
				}
			} else {
				if peer != client {
					_ = peer.Conn.WriteJSON(msg)
				}
			}
		}
		hub.Mutex.RUnlock()
	}
}

func broadcastToRoom(roomCode string, self *Client, msg map[string]interface{}) {
	hub.Mutex.RLock()
	defer hub.Mutex.RUnlock()
	for peer := range hub.Rooms[roomCode] {
		if peer != self {
			_ = peer.Conn.WriteJSON(msg)
		}
	}
}

func broadcastStudentCount(roomCode string) {
	hub.Mutex.RLock()
	defer hub.Mutex.RUnlock()
	clientsInRoom := hub.Rooms[roomCode]

	studentCount := 0
	for peer := range clientsInRoom {
		if peer.Role == "student" {
			studentCount++
		}
	}

	msg := map[string]interface{}{
		"type":  "student-count-update",
		"count": studentCount,
	}

	for peer := range clientsInRoom {
		_ = peer.Conn.WriteJSON(msg)
	}
}

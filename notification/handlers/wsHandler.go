package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	authMiddleware "github.com/rzaf/youtube-clone/auth/middlewares"
	user_pb "github.com/rzaf/youtube-clone/database/pbs/user-pb"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin: func(r *http.Request) bool {
	// 	return true
	// },
}

var userChannels = make(map[int64]chan []byte)
var userChannelsMutex = &sync.Mutex{}

func reader(ws *websocket.Conn) {
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(s string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn, userChannel chan []byte) {
	pingTicker := time.NewTicker(pingPeriod)
	log.Println("pingTicker")
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case message, ok := <-userChannel:
			if !ok {
				return
			}
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("error sending notification:", err)
				return
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("error", err)
				return
			}
		}
	}
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(authMiddleware.AuthUser("user")).(*user_pb.CurrentUserData)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	userChannelsMutex.Lock()
	if _, exists := userChannels[currentUser.Id]; !exists {
		userChannels[currentUser.Id] = make(chan []byte)
	}
	userChannel := userChannels[currentUser.Id]
	userChannelsMutex.Unlock()

	defer func() {
		userChannelsMutex.Lock()
		close(userChannel)
		delete(userChannels, currentUser.Id)
		userChannelsMutex.Unlock()
	}()

	go writer(c, userChannel)
	reader(c)
}

func SendWsMessage(id string, userId int64, title, message string) {
	userChannelsMutex.Lock()
	userChannel, exists := userChannels[userId]
	userChannelsMutex.Unlock()

	if exists {
		jsonMessage, err := json.Marshal(map[string]string{
			"id":      id,
			"title":   title,
			"message": message,
		})
		if err != nil {
			log.Println("error json.Marshal: ", err)
			return
		}
		userChannel <- jsonMessage

	} else {
		log.Printf("User %d not connected. Notification dropped.", userId)
	}
}

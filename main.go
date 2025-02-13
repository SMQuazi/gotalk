package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type UserConnectionMap map[string]*websocket.Conn
type Message struct {
	Users UserConnectionMap `json:"users"`
}

var usersConnectionsMap = make(UserConnectionMap)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Server started ")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("authorization")
	if user == "" {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	if _, ok := usersConnectionsMap[user]; ok {
		http.Error(w, "User already connected", http.StatusUnauthorized)
		return
	}

	conn, connErr := upgrader.Upgrade(w, r, nil)
	loginUser(user, conn)
	if connErr != nil {
		log.Fatal(connErr)
	}
	defer closeConnection(conn, user)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("User '%s' sent: %s\n", user, msg)
	}

}

func loginUser(user string, conn *websocket.Conn) {
	fmt.Printf("User '%s' connected.\n", user)
	usersConnectionsMap[user] = conn
	message := Message{Users: usersConnectionsMap}
	json, _ := json.Marshal(message)
	for user := range usersConnectionsMap {
		usersConnectionsMap[user].WriteMessage(websocket.TextMessage, json)
	}
}

func closeConnection(conn *websocket.Conn, user string) {
	delete(usersConnectionsMap, user)
	conn.Close()
	fmt.Println("Client disconnected")
}

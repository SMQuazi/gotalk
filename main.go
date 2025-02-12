package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var users = make(map[string]bool)

type Message struct {
	users   map[string]bool
	message string
}

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

func closeConnection(conn *websocket.Conn, user string) {
	delete(users, user)
	conn.Close()
	fmt.Println("Client disconnected")
}

func loginUser(w http.ResponseWriter, r *http.Request) (string, error) {
	user := r.URL.Query().Get("authorization")
	if user == "" {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return "", errors.New("user not found")
	}
	if val, ok := users[user]; ok && val {
		http.Error(w, "User already connected", http.StatusUnauthorized)
		return "", errors.New("user already connected")
	}

	fmt.Printf("User '%s' connected.\n", user)
	users[user] = true
	return user, nil
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	user, loginError := loginUser(w, r)
	if user == "" && loginError != nil {
		log.Println(loginError)
		return
	}

	conn, connErr := upgrader.Upgrade(w, r, nil)
	if connErr != nil {
		log.Fatal(connErr)
	}
	defer closeConnection(conn, user)

	message := Message{
		users:   users,
		message: "",
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		message.message = string(p)

		writeErr := conn.WriteMessage(messageType, []byte(fmt.Sprintf("%#v", message)))
		if writeErr != nil {
			break
		}
		fmt.Println(string(p))

	}
}

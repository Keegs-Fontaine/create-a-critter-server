package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type UserPayload struct {
	UserName string `json:"userName"`
}

type User struct {
	UserName string
}

func main() {
	fmt.Print("Go SERVER!!!")

	users := []User{}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()

		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)
			err = c.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	})

	http.HandleFunc("POST /user", func(w http.ResponseWriter, r *http.Request) {
		var u UserPayload

		err := json.NewDecoder(r.Body).Decode(&u)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:5173")
		fmt.Println("I AM READING FROM USER ", u.UserName)

		users = append(users, User{UserName: u.UserName})

		w.Write([]byte("fuck"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

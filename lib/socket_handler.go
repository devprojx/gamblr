package lib

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	var cacheGames []Game
	c := make(chan Game)

	go func() {
		for {
			l := len(cacheGames)
			games := FetchGames()
			if l > 0 {
				time.Sleep(30 * time.Second)
			}

			for _, game := range games {
				//Send all the games through the channel if cache is empty
				if l == 0 {
					cacheGames = append(cacheGames, game)
					c <- game
					continue
				}

				//Check if cached data is outdated
				for idx, cache := range cacheGames {
					if cache.Name == game.Name && (cache.Date != game.Date || cache.TimeOfDay != game.TimeOfDay) {
						cacheGames[idx] = game
						c <- game
					}
				}
			}

		}
	}()

	//Websocket Writer
	go func() {
		for n := range c {
			// pids := strings.Join(n.PIDs, " ")
			f, err := json.Marshal(n)
			if err != nil {
				log.Println("[error] error parsing json: ", err)
				return
			}

			ws.WriteMessage(websocket.TextMessage, f)
		}
	}()

	//Websocket Reader
	go func(c chan Game, conn *websocket.Conn) {
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				close(c)
				log.Println("[disconnect]: ", err)
				return
			}
			log.Println(string(p))

			if err := conn.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				return
			}

		}
	}(c, ws)
}

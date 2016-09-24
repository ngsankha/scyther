package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/rpc"
)

type ScytherMessage struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

var HttpAddr = flag.String("http", ":13921", "Host:Port")
var RcpAddr = flag.String("rcp", ":13922", "RcpHost:Port")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func authorize(w http.ResponseWriter, r *http.Request) {
	var rpcClient *rpc.Client
	isConnected := false
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		res := ScytherMessage{}
		_ = json.Unmarshal(message, &res)
		log.Printf("Received message type %s, with data %s\n", res.Type, res.Value)

		if res.Type == "handshake" && res.Value == "hello scyther native" {
			isConnected = true
			data := ScytherMessage{"handshake", "hello scyther web"}
			ret, _ := json.Marshal(&data)
			err = conn.WriteMessage(mt, ret)
			if err != nil {
				log.Println("write:", err)
			}
			continue
		}

		if isConnected {
			log.Printf("I am connected %s", res.Type)
			switch res.Type {
			case "request":
				go relay(conn, res.Value, mt)
			case "privilege":
				rpcClient, err = rpc.Dial("tcp", *RcpAddr)
				if err != nil {
					log.Fatal(err)
				}
				go get_privilege(conn, res.Value, mt, rpcClient)
			default:
				data := ScytherMessage{"response", "No method implemented"}
				ret, _ := json.Marshal(&data)
				err = conn.WriteMessage(mt, ret)
				if err != nil {
					log.Println("write:", err)
				}
			}
		} else {
			data := ScytherMessage{"response", "Unauthorized"}
			ret, _ := json.Marshal(&data)
			err = conn.WriteMessage(mt, ret)
			if err != nil {
				log.Println("write:", err)
			}
		}
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/auth", authorize)
	log.Printf("Server started on %s\n", *HttpAddr)
	log.Fatal(http.ListenAndServe(*HttpAddr, nil))
}

package main

import (
  "flag"
  "log"
  "net/http"
  "encoding/json"
  "github.com/gorilla/websocket"
)

type ScytherMessage struct {
  Type string  `json:"type"`
  Value string  `json:"value"`
}

var HttpAddr = flag.String("http", ":13921", "Host:Port")

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func authorize(w http.ResponseWriter, r *http.Request) {
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
    log.Printf("recv: %s", message)

    res := ScytherMessage{}
    _ = json.Unmarshal(message, &res)

    if( res.Type == "handshake" && res.Value == "hello scyther native" ) {
      data := ScytherMessage{"handshake", "hello scyther web"}
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

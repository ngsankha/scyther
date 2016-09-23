package main

import (
  "log"
  "encoding/json"
  "github.com/gorilla/websocket"
)

func relay(ws *websocket.Conn, message string, messageType int) {
  log.Printf("Relaying data\n"); 
  res := ScytherMessage{"response", "Good"}
  ret, _ := json.Marshal(&res)
  err := ws.WriteMessage(messageType, ret)
  if err != nil {
    log.Println("write:", err)
  }
}

func get_privilige(ws *websocket.Conn, message string, messageType int) {
  log.Printf("Getting privileges\n"); 
  res := ScytherMessage{"response", "Good"}
  ret, _ := json.Marshal(&res)
  err := ws.WriteMessage(messageType, ret)
  if err != nil {
    log.Println("write:", err)
  }
}

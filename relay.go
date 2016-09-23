package main

import (
  "log"
  "net/rpc"
  "encoding/json"
  "github.com/gorilla/websocket"
)

type Args struct {}

type BLAdvertisement struct {
    ID string
    Name string
    LocalName string
    TxPowerLevel int
    ManufacturerData []byte
}

func relay(ws *websocket.Conn, message string, messageType int) {
  log.Printf("Relaying data\n"); 
  res := ScytherMessage{"response", "Bad"}
  ret, _ := json.Marshal(&res)
  err := ws.WriteMessage(messageType, ret)
  if err != nil {
    log.Println("write:", err)
  }
}

func get_privilige(ws *websocket.Conn, message string, messageType int, rpcClient *rpc.Client) {
  log.Printf("Getting privileges\n")
  reply := new([]BLAdvertisement)
  args := Args{}
  err := rpcClient.Call("Bluetooth.Peripheral", args, &reply)
  log.Printf("%d\n", reply)
  res := ScytherMessage{"response", "Bad"}
  ret, _ := json.Marshal(&res)
  err = ws.WriteMessage(messageType, ret)
  if err != nil {
    log.Println("write:", err)
  }
}

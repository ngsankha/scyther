package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/rpc"
)

type Args struct{}

type BLAdvertisement struct {
	ID               string
	Name             string
	LocalName        string
	TxPowerLevel     int
	ManufacturerData []byte
}

func relay(ws *websocket.Conn, message string, messageType int) {
	log.Printf("Relaying data\n")
	res := ScytherMessage{"response", "Bad"}
	ret, _ := json.Marshal(&res)
	err := ws.WriteMessage(messageType, ret)
	if err != nil {
		log.Println("write:", err)
	}
}

func get_privilege(ws *websocket.Conn, message string, messageType int, rpcClient *rpc.Client) {
	log.Printf("Getting privileges\n")
	reply := new([]BLAdvertisement)
	args := Args{}
	err := rpcClient.Call("Bluetooth.Peripheral", args, &reply)
	// log.Println(*reply)
	d, _ := json.Marshal(reply)
	n := len(d)
	res := ScytherMessage{"response", string(d[:n])}
	d, _ = json.Marshal(res)
	err = ws.WriteMessage(messageType, d)
	if err != nil {
		log.Println("write:", err)
	}
}

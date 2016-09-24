package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

type Args struct{}
type BLAdvertisement struct {
	ID               string
	Name             string
	LocalName        string
	TxPowerLevel     int
	ManufacturerData []byte
}
type Bluetooth int

var blAds []BLAdvertisement

func (t *Bluetooth) Peripheral(args *Args, reply *[]BLAdvertisement) error {
	*reply = blAds
	return nil
}

func onStateChanged(d gatt.Device, s gatt.State) {
	fmt.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("scanning...")
		d.Scan([]gatt.UUID{}, false)
		return
	default:
		d.StopScanning()
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	blAd := BLAdvertisement{p.ID(), p.Name(), a.LocalName, a.TxPowerLevel, a.ManufacturerData}
	blAds = append(blAds, blAd)
	fmt.Printf("\nPeripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())
	fmt.Println("  Local Name        =", a.LocalName)
	fmt.Println("  TX Power Level    =", a.TxPowerLevel)
	fmt.Println("  Manufacturer Data =", a.ManufacturerData)
	fmt.Println("  Service Data      =", a.ServiceData)
}

func main() {
	blAds = make([]BLAdvertisement, 0, 20)
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}
	// Register handlers.
	d.Handle(gatt.PeripheralDiscovered(onPeriphDiscovered))
	d.Init(onStateChanged)

	bl := new(Bluetooth)
	rpc.Register(bl)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":13922")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go rpc.Accept(l)
	fmt.Println("loaded bluetooth plugin")

	select {}
}

package main

import (
  "log"
  "time"
  "github.com/fulr/rfm69"
  "os"
  "net/http"
  _"encoding/json"
  "github.com/gorilla/mux"
)

type Device struct {
  radio *rfm69.Device
}
type Data struct {
  *rfm69.Data
}

var networkID = byte(100)
var nodeID = byte(0x01)
var rx = make(chan *rfm69.Data, 5)
var sigint = make(chan os.Signal, 1)
var device Device

func receiveData() func(d *rfm69.Data) {
  f := func(d *rfm69.Data) {
    log.Println("Data received:", d)
    rx <- d
  }
  return f
}

func receiveACK(fromAddress byte, delay int64) bool {
  running := true
  start := time.Now()
  log.Println("ReceiveACK: entry",", time:",start)

  for running {
    select {
      case data := <-rx:
	log.Println("got data from", data.FromAddress, ", RSSI", data.Rssi)
        if data.ToAddress != nodeID {
          break
        }
        if data.FromAddress == fromAddress && data.SendAck {
           return true
        }
      default:
    }
    time.Sleep(1 * time.Millisecond)
    period := time.Now().Sub(start).Nanoseconds() / 1000000
    if period > delay {
        running = false;
    }
  }
  log.Println("receiveACK: exit")
  return false;
}

func (rfm *Device) sendWithRetry(d *rfm69.Data, numRetry int, retryDelay int64) bool {
  rfm.radio.OnReceive = receiveData()
  for i:=0; i<numRetry; i++ {
    log.Println("Retry: ",i)
    rfm.radio.Send(d)
    log.Println("Sent#",i)
    if receiveACK(d.ToAddress, retryDelay) {
      return true
    }
  }
  return false
}

func main() {  
  log.Println("Hello World")

  rfm,err := rfm69.NewDevice(0x01, 0x00, true)
  if err != nil {
    panic(err)
  }
  defer rfm.Close()
/*
  device := &Device{
              radio: rfm,
  }
*/

  log.Println("set frequency...")
  err = rfm.SetFrequency("915")
  if err != nil {
    panic(err)
  }

  log.Println("set power level...")
  err = rfm.SetPowerLevel(byte(31))
  if err != nil {
    panic(err)
  }

  log.Println("Set encryption key")
  err = rfm.Encrypt([]byte("1234567890111222"))
  if err != nil {
    panic(err)
  }
/*
  device.sendWithRetry(&rfm69.Data{
    ToAddress: byte(2),
    Data: []byte("Hello there. What's up?"),
    RequestAck: true,
  }, 3, 1000)
*/

  router := mux.NewRouter()
  router.HandleFunc("/things", GetThings).Methods("GET")
  router.HandleFunc("/thing/{id}", GetThing).Methods("GET")

  log.Println("Starting REST server at port 8000")
  log.Fatal(http.ListenAndServe(":8000", router))
}

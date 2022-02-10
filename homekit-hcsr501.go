package main

import (
  "github.com/brutella/hc"
  "github.com/brutella/hc/accessory"
  "github.com/brutella/hc/characteristic"
  "github.com/brutella/hc/service"
  "github.com/stianeikeland/go-rpio/v4"

  "flag"
  "fmt"
  "log"
  "os"
  "time"
)

var sensorPin int
var developmentMode bool

func init() {
  flag.IntVar(&sensorPin, "pin", 23, "sensor pin your HC-SR501 is connected to, an int")
  flag.BoolVar(&developmentMode, "dev", false, "turn on development mode to return a random reading, boolean")
  flag.Parse()

  if developmentMode == true {
    log.Println("Development mode on, ignoring sensor and returning random values...")  
  }
}

func main() {
  info := accessory.Info{
    Name:             "HC-SR501",
    SerialNumber:     "18420",
    Manufacturer:     "Kuongshun",
    Model:            "HC-SR501",
    FirmwareRevision: "1.0.0",
  }
  acc := accessory.New(
    info,
    10,  // Sensor
  )
  motion := service.NewMotionSensor()
  motionDetected := characteristic.NewMotionDetected()
  motion.Service.AddCharacteristic(motionDetected.Characteristic)
  acc.AddService(motion.Service)
  config := hc.Config{
    // Change the default Apple Accessory Pin if you wish
    Pin: "00102003",
    // Port: "12345",
    // StoragePath: "./db",
  }
  t, err := hc.NewIPTransport(config, acc)
  if err != nil {
    log.Fatal(err)
  }

  go func() {
    pin := rpio.Pin(23)
    // Open and map memory to access gpio, check for errors
    if err := rpio.Open(); err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Unmap gpio memory when done
    defer rpio.Close()

    // Set pin to input mode
    pin.Input()

    for {
      res := pin.Read()
      if res > 0 {
        fmt.Println("Motion detected!")
        motion.MotionDetected.SetValue(true)
      } else {
        motion.MotionDetected.SetValue(false)
      }
      time.Sleep(time.Second / 10)
    }
  }()

  hc.OnTermination(func() {
    <-t.Stop()
  })

  t.Start()
}

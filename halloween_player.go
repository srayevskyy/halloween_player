package main

import (
	"flag"
	"fmt"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
	"log"
	"math/rand"
	"os/exec"
	"path/filepath"
)

func CheckError(message string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

func main() {
	resource_path := "/home/pi/Halloween/resources/"
	resource_extension := "wav"
	sensor_pin := 23

	// file resource init
	files, err := filepath.Glob(resource_path + "/*" + resource_extension)
	CheckError("Error wile scanning resource directory", err)
	log.Printf("Found %d resource file(s) in %s", len(files), resource_path)

	// random number generator init
	rand.Seed(42)

	// GPIO init
	flag.Parse()
	err = embd.InitGPIO()
	CheckError("Error initlalizing GPIO", err)
	defer embd.CloseGPIO()

	btn, err := embd.NewDigitalPin(sensor_pin)
	CheckError(fmt.Sprintf("Error wile accessing GPIO pin %d", sensor_pin), err)
	defer btn.Close()

	err = btn.SetDirection(embd.In)
	CheckError(fmt.Sprintf("Error setting mode for GPIO pin %d", sensor_pin), err)

	btn.ActiveLow(false)

	quit := make(chan interface{})

	log.Printf("Waiting for sensor to trigger...")

	err = btn.Watch(embd.EdgeBoth,
		func(btn embd.DigitalPin) {
			quit <- btn
		})
	CheckError("Error wile watching for sensor to trigger", err)

	log.Printf("Sensor %v has been triggered.\n", <-quit)

	cmd := exec.Command("omxplayer", files[rand.Intn(len(files))])
	err = cmd.Run()
	CheckError("Cannot play audio file", err)

}

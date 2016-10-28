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
	"time"
)

func CheckError(message string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

func main() {
	resource_path := "/home/pi/Halloween/resources/"
	resource_extension := "wav"
	sensor_pin := 21
	sleep_interval := 100 * time.Millisecond

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

	err = btn.SetDirection(embd.In)
	CheckError(fmt.Sprintf("Error setting mode for GPIO pin %d", sensor_pin), err)
	defer btn.Close()

	btn.ActiveLow(false)

	for {

		state, _ := btn.Read()
		log.Printf("Sensor state: %d", state)

		if state == 1 {
			log.Printf("Waiting for sensor to calm down...")
			for state == 1 {
				state, _ = btn.Read()
				if state == 1 {
					time.Sleep(sleep_interval)
				}
			}
			log.Printf("Sensor is calm.")
		}

		// additional check whether the sensor is calm
		state, _ = btn.Read()
		if state == 0 {
			// log.Printf("Sensor state before second loop: %d", state)

			log.Printf("Waiting for sensor to trigger...")

			for state == 0 {
				state, _ = btn.Read()
				if state == 0 {
					time.Sleep(sleep_interval)
				}
			}

			file_to_play := files[rand.Intn(len(files))]
			log.Printf("Motion sensor has been triggered, playing scary sound %s", file_to_play)
			cmd := exec.Command("omxplayer", file_to_play)
			err = cmd.Run()
			CheckError("Cannot play audio file", err)
		}
	}
}

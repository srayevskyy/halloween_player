package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"math/big"
	"os"
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

	log.SetOutput(&lumberjack.Logger{
		Filename:   os.Getenv("HOME") + "/halloween_player/halloween_player.log",
		MaxSize:    10, // megabytes
		MaxBackups: 20, // number of logfiles to keep
		MaxAge:     1,  // days to keep logfiles, wins over MaxBackups
		LocalTime:  true,
	})

	resource_path := os.Getenv("HOME") + "/Halloween/resources/"
	resource_extension := "wav"
	sensor_pin := 21
	sleep_interval := 100 * time.Millisecond
	sound_probability_percent := 5

	// file resource init
	files, err := filepath.Glob(resource_path + "/*" + resource_extension)
	CheckError("Error wile scanning resource directory", err)
	log.Printf("Found %d resource file(s) in %s", len(files), resource_path)

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

			nBig, err := rand.Int(rand.Reader, big.NewInt(100))
			CheckError("Error generating random number for scary sound probability", err)
			log.Printf("Generated number for sound probability: %d", nBig.Int64())
			if nBig.Int64() < int64(sound_probability_percent) {
				log.Printf("OK, let's play some scary sound")
				nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(files))))
				CheckError("Error generating random number of scary sound", err)
				file_to_play := files[nBig.Int64()]
				log.Printf("Motion sensor has been triggered, playing scary sound %s", file_to_play)
				cmd := exec.Command("omxplayer", file_to_play)
				err = cmd.Run()
				CheckError("Cannot play audio file", err)
			} else {
				log.Printf("NO, let's not scare anyone this time")
			}
		}
	}
}

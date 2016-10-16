package main

import (
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

	// init random number generator
	rand.Seed(42)

	files, err := filepath.Glob(resource_path + "/*" + resource_extension)
	CheckError("Error wile scanning resource directory", err)
	log.Printf("Found %d resource file(s) in %s", len(files), resource_path)

	cmd := exec.Command("omxplayer", files[rand.Intn(len(files))])
	err = cmd.Run()
	CheckError("Cannot play audio file", err)

}

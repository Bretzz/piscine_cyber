package main

import (
	"log"
	"os"

	"scorpion/pkg/exif"
)

func main() {
	// fmt.Printf("Hello World\n")

	const path = "/home/topiana-/live_projects/Cyber_piscine/01 - arachnida/scorpion/pictures/20250616_153403.jpg"
	// const path = "/home/topiana-/live_projects/Cyber_piscine/01 - arachnida/scorpion/pictures/aftermarket-update-logo-white-bg.jpg"

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Error while reading file: ", err.Error())
	}

	err = exif.Parse(data)
	if err != nil {
		log.Fatal(err)
	}

	exif.Display()
}

package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	FaceCascadePath = "cascades/haarcascade_frontalface_alt.xml"
	EyeCascadePath  = "cascades/haarcascade_eye.xml"
	GlassesPath     = "assets/glasses.png"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Cannot parse env")
	}
	// Bot SECTION
	thugbot, err := NewThugbot(os.Getenv("SLACK_TOKEN"), "cascades/")
	if err != nil {
		log.Fatalln(err)
	}
	thugbot.Run()
}

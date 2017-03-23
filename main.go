package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// FaceCascadePath is path to haar classifier xml to detect faces
// EyeCascadePath is path ot haar classifier xml to detect eyes
// GlassesPath is path to the glasses image that will be drawn over original image
var (
	FaceCascadePath = "cascades/haarcascade_frontalface_alt.xml"
	EyeCascadePath  = "cascades/haarcascade_eye.xml"
	GlassesPath     = "assets/glasses.png"
	TextPath        = "assets/impact-regular.ttf"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Cannot parse env")
	}
	// Bot SECTION
	thugbot, err := NewThugbot(os.Getenv("SLACK_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}
	thugbot.Run()
}

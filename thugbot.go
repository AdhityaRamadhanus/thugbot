package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	thugly "github.com/AdhityaRamadhanus/thug.ly"
	"github.com/disintegration/imaging"
	"github.com/lazywei/go-opencv/opencv"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

const (
	helpText = "How can I help you? mention me and upload your photos, i will do some magic"
)

var (
	greetingPattern  = "hi bot|hello bot"
	greetingPrefixes = []string{"Hi", "Hello", "Howdy", "Wazzzup", "Hey"}
)

// Thugbot is the main struct consist of all the sub component including slack api, real-time messaing api and face detector
type Thugbot struct {
	ID       string
	RTM      *slack.RTM
	SlackApi *slack.Client
	Detector *thugly.Detector
}

// NewThugbot create new Thug bot
func NewThugbot(slackToken string) (*Thugbot, error) {
	thugbot := &Thugbot{
		SlackApi: slack.New(slackToken),
		Detector: thugly.NewDetector(FaceCascadePath, EyeCascadePath),
	}

	return thugbot, nil
}

func (t *Thugbot) messageHandler(ev *slack.MessageEvent) {
	// Simple Router for simple task
	if ev.Type == "message" &&
		ev.SubType == "file_share" &&
		strings.Contains(ev.Text, "<@"+t.ID+">") {
		go t.thuglify(ev)
		return
	}

	if ev.Type == "message" &&
		(strings.HasPrefix(strings.ToLower(ev.Text), "hi bot") ||
			strings.HasPrefix(strings.ToLower(ev.Text), "hello bot")) {
		go t.helloWorld(ev)
		return
	}

	if ev.Type == "message" && strings.HasPrefix(strings.ToLower(ev.Text), "bot help") {
		go t.help(ev)
		return
	}
}

func (t *Thugbot) thuglify(ev *slack.MessageEvent) (err error) {
	t.RTM.SendMessage(t.RTM.NewOutgoingMessage("Processing image, please wait ...", ev.Channel))
	defer func() {
		if err != nil {
			log.Println(err)
			t.RTM.SendMessage(t.RTM.NewOutgoingMessage("Failed to thuglify "+ev.File.Name, ev.Channel))
		}
	}()
	baseImage := DownloadImage(ev.File, os.Getenv("SLACK_TOKEN"))
	persons := t.Detector.DetectPersons(opencv.FromImage(baseImage))
	glassesImage, err := thugly.LoadImage(GlassesPath)
	if err != nil {
		return errors.Wrap(err, "Failed to load the glsses image")
	}
	bounds := baseImage.Bounds()
	canvas := image.NewRGBA(bounds)
	draw.Draw(canvas, bounds, baseImage, bounds.Min, draw.Over)
	for _, person := range persons {
		eyesReg := person.GetEyesRect()
		fitGlassImg := imaging.Resize(glassesImage, eyesReg.Dx(), eyesReg.Dy(), imaging.Lanczos)
		draw.Draw(
			canvas,
			eyesReg,
			fitGlassImg,
			bounds.Min,
			draw.Over)
	}
	outputFileName := strings.Split(ev.File.Title, ".")[0] + "-thuglify.jpg"
	outputImage, err := os.Create(outputFileName)
	if err != nil {
		return errors.Wrap(err, "Cannot open output image")
	}
	if err := jpeg.Encode(outputImage,
		canvas,
		&jpeg.Options{
			Quality: jpeg.DefaultQuality,
		}); err != nil {
		return errors.Wrap(err, "Cannot encode jpeg output")
	}
	defer func() {
		os.Remove(outputFileName)
	}()
	// Upload File to Slack
	params := slack.FileUploadParameters{
		Title: outputFileName,
		File:  outputFileName,
		Channels: []string{
			ev.Channel,
		},
	}
	_, err = t.SlackApi.UploadFile(params)
	if err != nil {
		return errors.Wrap(err, "Cannot upload output image")
	}
	return nil
}

func (t *Thugbot) helloWorld(ev *slack.MessageEvent) (err error) {
	rand.Seed(time.Now().UnixNano())
	msg := greetingPrefixes[rand.Intn(len(greetingPrefixes))] + " <@" + ev.User + ">!"
	t.RTM.SendMessage(t.RTM.NewTypingMessage(ev.Channel))
	t.RTM.SendMessage(t.RTM.NewOutgoingMessage(msg, ev.Channel))
	return nil
}

func (t *Thugbot) help(ev *slack.MessageEvent) (err error) {
	t.RTM.SendMessage(t.RTM.NewOutgoingMessage(helpText, ev.Channel))
	return nil
}

// Run is main loop function to retrive slack message and process it
func (t *Thugbot) Run() {
	t.RTM = t.SlackApi.NewRTM()
	go t.RTM.ManageConnection()

	for msg := range t.RTM.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			t.ID = ev.Info.User.ID
			log.Println(ev.Info.User.ID, "Connected")
		case *slack.MessageEvent:
			t.messageHandler(ev)
		case *slack.RTMError:
			log.Println(ev.Error())
		case *slack.InvalidAuthEvent:
			log.Println("Failed to Authenticate")
			return
		default:
		}
	}
}

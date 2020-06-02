package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"github.com/micmonay/keybd_event"
)

var keywordConfig = map[string]int{
	"launch elite dangerous": keybd_event.VK_Z,
	"landing gear":           keybd_event.VK_X,
	"docking permissions":    keybd_event.VK_C,
}

var keyboard keybd_event.KeyBonding

func recognizedHandler(event speech.SpeechRecognitionEventArgs) {
	defer event.Close()
	fmt.Println("Recognized:", event.Result.Text)
	for k, v := range keywordConfig {
		if strings.Contains(strings.ToLower(event.Result.Text), k) {
			fmt.Println("Executing:", k)
			pressKey(v)
		}
	}
}

func cancelledHandler(event speech.SpeechRecognitionCanceledEventArgs) {
	defer event.Close()
	fmt.Println("Received a cancellation: ", event.ErrorDetails)
}

func pressKey(key int) {
	// Select keys to be pressed
	keyboard.SetKeys(key)

	// Press the selected keys
	err := keyboard.Launching()
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Starting up...")

	fmt.Println("Configuring keyboard presses...")
	var err error
	keyboard, err = keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// For Linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	fmt.Println("Configuring Azure speech interpretation...")
	subscription := os.Getenv("AZURE_SUBSCRIPTION_ID")
	region := os.Getenv("AZURE_REGION")

	audioConfig, err := audio.NewAudioConfigFromDefaultMicrophoneInput()
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer audioConfig.Close()
	config, err := speech.NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer config.Close()
	speechRecognizer, err := speech.NewSpeechRecognizerFromConfig(config, audioConfig)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer speechRecognizer.Close()
	fmt.Println("Listening for commands...")
	speechRecognizer.Recognized(recognizedHandler)
	speechRecognizer.Canceled(cancelledHandler)
	speechRecognizer.StartContinuousRecognitionAsync()
	defer speechRecognizer.StopContinuousRecognitionAsync()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

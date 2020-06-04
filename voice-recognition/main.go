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
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/sphero/bb8"
	"gobot.io/x/gobot/platforms/sphero/ollie"
)

var keywordConfig = map[string]int{
	"weapons":     keybd_event.VK_T,
	"hard points": keybd_event.VK_T,
	"hardpoints":  keybd_event.VK_T,

	"landing gear": keybd_event.VK_F,

	"docking permissions": keybd_event.VK_Q,

	"super cruise": keybd_event.VK_E,
	"supercruise":  keybd_event.VK_E,
	"warp drive":   keybd_event.VK_E,
	"light speed":  keybd_event.VK_E,
	"lightspeed":   keybd_event.VK_E,
	"drop":         keybd_event.VK_E,

	"take us up": keybd_event.VK_R,
	"get going":  keybd_event.VK_R,
	"lift off":   keybd_event.VK_R,
	"launch":     keybd_event.VK_R,

	"communications panel": keybd_event.VK_W,
	"comms panel":          keybd_event.VK_W,
	"communications":       keybd_event.VK_W,
	"communication":        keybd_event.VK_W,

	"ship panel": keybd_event.VK_D,

	"control panel": keybd_event.VK_S,

	"navigation panel": keybd_event.VK_A,
	"nav panel":        keybd_event.VK_A,
	"navigation":       keybd_event.VK_A,
}

var keyboard keybd_event.KeyBonding
var bb8Id = "BB-9B05"
var bb8Conn *bb8.BB8Driver
var batteryStates = []string{"charging", "battery ok", "battery low", "battery critical"}

func recognizedHandler(event speech.SpeechRecognitionEventArgs) {
	defer event.Close()
	fmt.Println("Recognized:", event.Result.Text)
	for k, v := range keywordConfig {
		if strings.Contains(strings.ToLower(event.Result.Text), k) {
			fmt.Println("Executing:", k)
			pressKey(v)
			bb8Flash("success")
			return
		}
	}

	bb8Flash("error")
	return
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

func bb8Work() {
	gobot.Every(time.Second*30, func() {
		bb8Conn.GetPowerState(func(powerState ollie.PowerStatePacket) {
			fmt.Println("Battery: " + batteryStates[powerState.PowerState-1])

			if powerState.PowerState-1 < 2 {
				angle := uint16(gobot.Rand(200))
				bb8Conn.Roll(1, angle)
				time.Sleep(time.Second)
				bb8Conn.Roll(0, angle)
			}
		})
	})
}

func bb8Flash(flashType string) {
	r, g, b := uint8(0), uint8(0), uint8(0)

	if flashType == "error" {
		r, g, b = 197, 101, 4
	} else if flashType == "success" {
		r, g, b = 2, 156, 254
	}

	bb8Conn.SetRGB(r, g, b)
	time.Sleep(time.Second * 2)
	bb8Conn.SetRGB(0, 0, 0)
}

func main() {
	fmt.Println("Starting up...")

	fmt.Println("Configuring BB-8 integration...")
	bleAdaptor := ble.NewClientAdaptor(bb8Id)
	bb8Conn = bb8.NewDriver(bleAdaptor)
	bb8Conn.SetRGB(0, 0, 0)
	robot := gobot.NewRobot("bb",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{bb8Conn},
		bb8Work,
	)
	robot.Start()

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

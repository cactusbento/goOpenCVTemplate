package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bendahl/uinput"
	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	hook "github.com/robotn/gohook"

	"gocv.io/x/gocv"
)

const (
	GAME = "Enter the Gungeon"
	IMG_DIR = "img/"
)

var (
	X, Y, W, H int 
)

func main() {
	out, err := exec.Command("xdotool", "search", "--name", "--onlyvisible", GAME).Output()
	if err != nil { log.Fatalln("Could not find", GAME) }
	XID, _ := strconv.Atoi( strings.TrimSuffix(string(out), "\n") )

	fmt.Println("Found", GAME, "\nXID:", XID)
	exec.Command("xdotool", "windowactivate", fmt.Sprint(XID)).Run()

	// Virtual Input Devices Setup
	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("goKeyboard")) 
	if err != nil { log.Fatalln("Could not create keyboard device") }
	defer keyboard.Close()

	mouse, err := uinput.CreateMouse("/dev/uinput", []byte("goMouse")) 
	if err != nil { log.Fatalln("Could not create mouse device") }
	defer mouse.Close()
	
	// OpenCV Window
	cvWin := gocv.NewWindow("Bot Vision")
	defer cvWin.Close()

	h := hook.Start()
	defer hook.End()
	
	counter := 0
	negFile, _ := os.OpenFile("annotations/neg.txt", os.O_RDWR|os.O_CREATE, 0755)
	defer negFile.Close()

	go func() {
		for event := range h {
			if event.Kind == hook.KeyHold {
				fileName := fmt.Sprintf("%04d.jpg", counter )

				// POSITIVE
				if event.Keycode == hook.Keycode["v"] {
					fmt.Println("[Screenshot] POSITIVE:", IMG_DIR+"pos/"+fileName)
					go snap("pos/"+fileName)
					counter += 1
				}
				
				// NEGATIVE
				if event.Keycode == hook.Keycode["b"] {
					fmt.Println("[Screenshot] NEGATIVE:", IMG_DIR+"neg/"+fileName)
					go snap("neg/"+fileName)
					negFile.WriteString(IMG_DIR+"neg/"+fileName+"\n")
					counter += 1
				}
			}
		}
	}()

	cascade := gocv.NewCascadeClassifier()
	defer cascade.Close()

	model_present := false
	if _, err := os.Stat("cascade/cascade.xml"); !errors.Is(err, os.ErrNotExist) {
		cascade.Load("cascade/cascade.xml")
		model_present = true
	}
	
	
	for {
		X, Y, W, H = XIDBounds(XID)

		img, err := screenshot.CaptureRect(image.Rect(X, Y, X+W, Y+H))
		if err != nil { log.Fatalf("Failed to capture screen section: %v\n", err) }

		cvMat, err := gocv.ImageToMatRGB(img)
		if err != nil { log.Fatalf("Cannot convert image.Image ro gocv.Mat") }

		if model_present {
			finds := cascade.DetectMultiScale(cvMat)
			for _, v := range finds {
				wg := new(sync.WaitGroup)
				wg.Add(1)
				go func() {
					defer wg.Done()
					for i := v.Min.X; i < v.Max.X; i++ {
						img.Set(i, v.Min.Y, color.RGBA{R: 255, A: 255})
						img.Set(i, v.Max.Y, color.RGBA{R: 255, A: 255})
					}
				}()

				for i := v.Min.Y; i < v.Max.Y; i++ {
					img.Set(v.Min.X, i, color.RGBA{R: 255, A: 255})
					img.Set(v.Max.X, i, color.RGBA{R: 255, A: 255})
				}
				wg.Wait()
			}
		}

		cvMat, err = gocv.ImageToMatRGB(img)
		if err != nil { log.Fatalf("Cannot convert image.Image ro gocv.Mat") }

		cvWin.IMShow(cvMat)
		cvWin.WaitKey(1)
	}

}

func snap(fileName string) {
	scrsht, err := screenshot.CaptureRect(image.Rect(X, Y, X+W, Y+H))
	if err != nil {
	   log.Println("FAILED TO CAPTURE SCREENSHOT")
	}

	nImg, _ := os.Create(IMG_DIR + fileName)
	defer nImg.Close()
	jpeg.Encode(nImg, scrsht, nil)

}

func moveMouse(mouse uinput.Mouse, X, Y int) {
	CX, CY := robotgo.GetMousePos()

	for CX != X || CY != Y {
		CX, CY = robotgo.GetMousePos()

		if CX < X {
			mouse.MoveRight(1)
		}

		if CX > X {
			mouse.MoveLeft(1) 
		}

		if CY < Y {
			mouse.MoveDown(1) 
		}
		
		if CY > Y {
			mouse.MoveUp(1)
		}
		time.Sleep(time.Millisecond * 2)
	}

}


func XIDBounds(xid int) (X, Y, W, H int) {
	out, err := exec.Command("xwininfo", "-id", fmt.Sprint(xid)).Output()
	if err != nil { 
		log.Println(err)
		X, Y, W, H = 0, 0, 0, 0
		return
	}

	wininfo := strings.Split(string(out), "\n")
	info := make([]int, 0)

	for _, v := range wininfo {
		if strings.Contains(v, "Absolute") || strings.Contains(v, "Width") || strings.Contains(v, "Height") {
			kv := strings.SplitAfter(v, ": ")
			kvv := strings.TrimSpace(kv[len(kv)-1])

			extracted, err := strconv.Atoi(kvv)
			if err != nil {log.Fatalln(err)}

			info = append(info, extracted)
		}
	}
	
	//log.Println(info)
	return info[0], info[1] , info[2], info[3]
}

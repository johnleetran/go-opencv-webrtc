package main

import (
	"fmt"
	"go-yuv2webRTC/screenshot"
	"go-yuv2webRTC/webrtc"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"gocv.io/x/gocv"

	_ "image/jpeg"
	_ "image/png"
)

var screenWidth int
var screenHeight int
var resizeWidth int
var resizeHeight int
var webRTC *webrtc.WebRTC
var imageFilePath string

var cvImg gocv.Mat
var sigmaX float64 = 0.0
var sigmaY float64 = 0.0

func init() {

	//initialize image stuff
	imageFilePath = "/Users/john/go/src/go-yuv2webRTC/IMG_1566.png"
	cvImg = gocv.IMRead(imageFilePath, gocv.IMReadUnchanged)

	screenWidth, screenHeight = screenshot.GetScreenSize(cvImg.Cols(), cvImg.Rows())
	resizeWidth, resizeHeight = screenWidth/2, screenHeight/2
	webRTC = webrtc.NewWebRTC()
	// start screenshot loop, wait for connection
	go screenshotLoop()
}

func main() {
	fmt.Println("http://localhost:8000")

	router := mux.NewRouter()
	router.HandleFunc("/", getWeb).Methods("GET")
	router.HandleFunc("/session", postSession).Methods("POST")

	http.ListenAndServe(":8000", router)
}

func getWeb(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadFile("./index.html")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(bs)
}

func postSession(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	r.Body.Close()

	localSession, err := webRTC.StartClient(string(bs), resizeWidth, resizeHeight)
	if err != nil {
		log.Fatalln(err)
	}

	w.Write([]byte(localSession))
}

func screenshotLoop() {
	// prevX := -1.0
	// prevY := -1.0
	for {
		if webRTC.IsConnected() {
			blurImage := gocv.NewMat()
			sigmaX = webRTC.ClientDataChannelMessage.SigmaX
			sigmaY = webRTC.ClientDataChannelMessage.SigmaY

			fmt.Println("x: %f y: %f", sigmaX, sigmaY)
			if sigmaX == 0 || sigmaY == 0 {
				blurImage = cvImg.Clone()
			} else {
				gocv.GaussianBlur(cvImg, &blurImage, image.Point{0, 0}, sigmaX, sigmaY, gocv.BorderDefault)
			}

			rbgImage, _ := blurImage.ToImage()
			img := resize.Resize(uint(resizeWidth), uint(resizeHeight), rbgImage, resize.Bilinear).(*image.RGBA)

			yuv := screenshot.RgbaToYuv(img)
			webRTC.ImageChannel <- yuv
			// prevX = sigmaX
			// prevY = sigmaY
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// func screenshotLoopReadingFile() {
// 	for {
// 		if webRTC.IsConnected() {
// 			//rgbaImg := screenshot.GetScreenshot(0, 0, screenWidth, screenHeight, resizeWidth, resizeHeight)

// 			imgfile, err := os.Open(imageFilePath)

// 			if err != nil {
// 				fmt.Println(err)
// 				os.Exit(1)
// 			}

// 			rgbaImg, _, err := image.Decode(imgfile)

// 			if err != nil {
// 				fmt.Println(err)
// 				os.Exit(1)
// 			}
// 			img := resize.Resize(uint(resizeWidth), uint(resizeHeight), rgbaImg, resize.Bilinear).(*image.RGBA)

// 			yuv := screenshot.RgbaToYuv(img)
// 			webRTC.ImageChannel <- yuv
// 			imgfile.Close()
// 		}
// 		//time.Sleep(10 * time.Millisecond)
// 	}
// }

// func screenshotLoopOriginal() {
// 	for {
// 		if webRTC.IsConnected() {
// 			rgbaImg := screenshot.GetScreenshot(0, 0, screenWidth, screenHeight, resizeWidth, resizeHeight)
// 			yuv := screenshot.RgbaToYuv(rgbaImg)
// 			webRTC.ImageChannel <- yuv
// 		}
// 		time.Sleep(10 * time.Millisecond)
// 	}
// }

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/dhowden/raspicam"
	"github.com/gorilla/mux"
)

const OUTPUT_DIRECTORY = "output"

func setup() {
	path := OUTPUT_DIRECTORY
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Creating output image directory")
		os.Mkdir(path, os.ModeDir|os.ModePerm)
	}
}

type CameraImage struct {
	Name   string
	Lock   sync.Mutex
	Buffer bytes.Buffer
}

type CameraImages struct {
	CurrentImage *CameraImage
	LoadingImage *CameraImage
	Lock         sync.Mutex
	Count        int
}

var CameraImagesBroker CameraImages

func CaptureImages() {

	for {
		timestamp := time.Now().Unix()
		photo_name := fmt.Sprintf("%d", timestamp)

		if CameraImagesBroker.LoadingImage == nil {
			CameraImagesBroker.LoadingImage = new(CameraImage)
		}
		CameraImagesBroker.LoadingImage.Lock.Lock()
		CameraImagesBroker.LoadingImage.Name = photo_name
		w := bufio.NewWriter(&CameraImagesBroker.LoadingImage.Buffer)
		s := raspicam.NewStill()
		errCh := make(chan error)
		go func() {
			for x := range errCh {
				fmt.Fprintf(os.Stderr, "%v\n", x)
			}
		}()

		raspicam.Capture(s, w, errCh)

		CameraImagesBroker.LoadingImage.Lock.Unlock()
		bWasNil := true
		var oldImage *CameraImage
		if CameraImagesBroker.CurrentImage != nil {
			bWasNil = false
			oldImage = CameraImagesBroker.CurrentImage
			oldImage.Lock.Lock()

		}

		CameraImagesBroker.Lock.Lock()

		// Swap out the images while broker is locked
		CameraImagesBroker.CurrentImage = CameraImagesBroker.LoadingImage
		CameraImagesBroker.LoadingImage = nil
		CameraImagesBroker.Count += 1
		CameraImagesBroker.Lock.Unlock()

		if !bWasNil {
			oldImage.Lock.Unlock()
		}
	}
}

func getLatestImage(w http.ResponseWriter, r *http.Request) {
	if CameraImagesBroker.CurrentImage == nil {
		// TODO: Something better than this
		return
	}

	CameraImagesBroker.CurrentImage.Lock.Lock()
	// fmt.Println("Getting Image " + CameraImagesBroker.CurrentImage.Name)
	w.Write(CameraImagesBroker.CurrentImage.Buffer.Bytes())
	CameraImagesBroker.CurrentImage.Lock.Unlock()
}

type StatusResponse struct {
	NumberImages int
	Status       string
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	var resp StatusResponse

	CameraImagesBroker.Lock.Lock()
	resp.NumberImages = CameraImagesBroker.Count
	CameraImagesBroker.Lock.Unlock()

	resp.Status = "ok"
	json.NewEncoder(w).Encode(&resp)
}

func getHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func main() {
	fmt.Println("Starting camera service")

	setup()

	go CaptureImages()

	fmt.Println("setting up router")

	r := mux.NewRouter()
	r.HandleFunc("/", getHomePage)
	r.HandleFunc("/image", getLatestImage)
	r.HandleFunc("/status", getStatus)
	http.Handle("/", r)
	fmt.Println("About to start server")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dhowden/raspicam"
)

const OUTPUT_DIRECTORY = "output"

func setup() {
	path := OUTPUT_DIRECTORY
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Creating output image directory")
		os.Mkdir(path, os.ModeDir|os.ModePerm)
	}
}

func main() {
	fmt.Println("Starting camera service")

	setup()

	for {
		timestamp := time.Now().Unix()

		filename := fmt.Sprintf("%s/%d.tmp", OUTPUT_DIRECTORY, timestamp)
		final_fn := fmt.Sprintf("%s/%d.jpg", OUTPUT_DIRECTORY, timestamp)
		lock_filename := fmt.Sprintf("%s/%d.lock", OUTPUT_DIRECTORY, timestamp)
		lock_file, err := os.Create(lock_filename)
		if err != nil {
			fmt.Println("Failed to create lock file")
			panic("Failed to create lock")
			// return
		}

		lock_file.Close()

		f, err := os.Create(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "create file: %v", err)
			panic("Failed to create file")
			// return
		}
		// defer f.Close()

		s := raspicam.NewStill()
		errCh := make(chan error)
		go func() {
			for x := range errCh {
				fmt.Fprintf(os.Stderr, "%v\n", x)
			}
		}()
		err = os.Remove(lock_filename)
		if err != nil {
			fmt.Println("Failed to remove lock")
			panic("Failed to remove lock")

		}
		// log.Println("Capturing image...")
		raspicam.Capture(s, f, errCh)
		f.Close()

		os.Rename(filename, final_fn)
		os.Chmod(final_fn, 0666)

	}

}

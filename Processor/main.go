package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"
)

var FILE_LIMIT = float64(200) // Num MB limit

func removeExcessFiles() bool {
	foldername := "./output"
	// fmt.Println(foldername)
	files, err := ioutil.ReadDir(foldername)
	if err != nil {
		log.Fatal(err)
		return false
	}

	var filesize int64

	for _, fname := range files {
		path := fmt.Sprintf("%s/%s", foldername, fname.Name())
		// fmt.Println("Filename: " + path)
		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
			return false
		}
		info, err := f.Stat()
		if err != nil {
			log.Fatal(err)
			return false
		}
		filesize += info.Size()
		// fmt.Println(f.Name())
	}

	fs := float64(filesize) / math.Pow10(6)
	// fmt.Printf("File Size : %fM\n", fs)

	if fs > FILE_LIMIT {
		path := fmt.Sprintf("%s/%s", foldername, files[0].Name())
		fmt.Println("Removing " + path)
		os.Remove(path)
		return true
	}

	return false
}

func main() {
	fmt.Println("Starting processor")
	for {
		ret := removeExcessFiles()
		if ret {
			time.Sleep(1 * time.Second)
		} else {
			time.Sleep(10 * time.Second)
		}
	}
}

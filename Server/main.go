package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func FilterDirsGlob(dir, suffix string) ([]string, error) {
	return filepath.Glob(filepath.Join(dir, suffix))
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func getLatestImage(w http.ResponseWriter, r *http.Request) {
	foldername := "./output"

	files, err := FilterDirsGlob(foldername, "*.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(files) < 1 {
		fmt.Println("No files to send")
		return
	}

	retFile := files[len(files)-1]
	// fmt.Println(retFile)
	http.ServeFile(w, r, retFile)
}

func getNamedImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fn := vars["image_name"]
	// fmt.Println("Getting " + fn)
	http.ServeFile(w, r, "output/"+fn)
}

func main() {
	fmt.Println("Starting Server")
	r := mux.NewRouter()
	r.HandleFunc("/", indexPage)
	r.HandleFunc("/latest", getLatestImage)
	r.HandleFunc("/image/{image_name}", getNamedImage)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

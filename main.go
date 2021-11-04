package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/stevenferrer/croptop"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./dist")))
	mux.HandleFunc("/crop", crop)

	addr := ":3002"
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Printf("Listening on %s\n", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	log.Println("Exited")
}

// Opts is crop options
type Opts struct {
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Rotate int     `json:"rotate"`
	ScaleX int     `json:"scaleX"`
	ScaleY int     `json:"scaleY"`
}

func crop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Printf("error parsing form: %v\n", err)
		http.Error(w, "bad request", 400)
		return
	}

	mpOpts, _, err := r.FormFile("opts")
	if err != nil {
		log.Printf("error getting opts: %v\n", err)
		http.Error(w, "internal server error", 500)
		return
	}

	var opts Opts
	err = json.NewDecoder(mpOpts).Decode(&opts)
	if err != nil {
		log.Printf("error decoding options: %v\n", err)
		http.Error(w, "internal server error", 500)
		return
	}

	mpImage, _, err := r.FormFile("image")
	if err != nil {
		log.Printf("error getting image: %v\n", err)
		http.Error(w, "internal server error", 500)
		return
	}

	img, err := croptop.Decode(mpImage)
	if err != nil {
		log.Printf("error decoding image: %v\n", err)
		http.Error(w, "internal server error", 500)
		return
	}

	err = img.Height(opts.Height).Width(opts.Width).
		OffsetX(opts.X).OffsetY(opts.Y).Crop().Encode(w)
	if err != nil {
		log.Printf("error encoding image: %v\n", err)
		http.Error(w, "internal server error", 500)
		return
	}
}

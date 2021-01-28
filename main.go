package main

import (
	"encoding/json"
	"image"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/disintegration/imaging"
	"github.com/go-chi/chi"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./dist")))
	mux.HandleFunc("/crop", crop)

	r := chi.NewRouter()
	fileSrv := http.FileServer(http.Dir("./dist"))
	r.Method(http.MethodGet, "/", fileSrv)
	r.Post("/crop", crop)

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
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Printf("error parsing form: %v\n", err)
		http.Error(w, "bad request", 400)
		return
	}

	mpOpts, _, err := r.FormFile("opts")
	if err != nil {
		log.Printf("error getting meta: %v\n", err)
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
		log.Printf("error getting meta: %v\n", err)
		http.Error(w, "internal server error", 500)
		return
	}

	img, err := imaging.Decode(mpImage)
	if err != nil {
		log.Printf("error decoding image: %v\n", err)
		http.Error(w, "internal server error", 500)
		return
	}

	width := int(math.Round(opts.Width))
	height := int(math.Round(opts.Height))
	offsetX := int(math.Round(opts.X))
	offsetY := int(math.Round(opts.Y))
	rect := image.Rect(0, 0, width, height).
		Add(image.Pt(offsetX, offsetY))
	img = imaging.Crop(img, rect)

	err = imaging.Encode(w, img, imaging.JPEG)
	if err != nil {
		log.Printf("error encoding image: %v\n", err)
		http.Error(w, "internal server error", 500)
		return
	}
}

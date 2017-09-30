package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	listenAddr = flag.String("listen", ":8080", "http listen address")
	dir        = flag.String("dir", ".", "working directory")
)

func main() {
	flag.Parse()
	if err := os.Chdir(*dir); err != nil {
		log.Fatalf("unable to change into working directory: %v", *dir)
	}
	if err := getChannels(); err != nil {
		log.Fatalf("unable to get channels; %v", err)
	}
	log.Printf("Launching http server at %v", *listenAddr)
	log.Fatal(serve(*listenAddr))
}
func serve(listenAddr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/api", handleAPI)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	var staticFiles = []string{
		"/favicon.ico", "/sw.js",
	}
	for _, f := range staticFiles {
		mux.HandleFunc(f, handleStaticFile)
	}

	svr := &http.Server{
		Addr:           listenAddr,
		Handler:        mux,
		ReadTimeout:    5 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	return svr.ListenAndServe()

}
func handleStaticFile(w http.ResponseWriter, r *http.Request) {
	fn := filepath.Join("static", r.URL.Path)
	http.ServeFile(w, r, fn)
}
func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.Header().Add("X-XSS-Protection", "1; mode=block")
	w.Header().Add("X-Frame-Options", "SAMEORIGIN")
	w.Header().Add("X-UA-Compatible", "IE=edge")
	http.ServeFile(w, r, "index.html")
}

type apiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (a *apiResponse) render(w http.ResponseWriter) {
	if a == nil {
		a = &apiResponse{}
	}

	if a.Data == nil {
		a.Code = http.StatusServiceUnavailable
	}
	if a.Data != nil && a.Code == 0 {
		a.Code = http.StatusOK
	}
	if a.Code == 0 {
		a.Code = http.StatusInternalServerError
	}
	if a.Message == "" {
		a.Message = http.StatusText(a.Code)
	}
	dat, err := json.Marshal(a)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, []byte(`{"code":500}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, dat)
}
func handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET")

	if r.Method == "OPTIONS" {
		return // preflight check
	}
	if r.Method != "GET" {
		return // bad
	}

	resource := r.URL.Query().Get("q")
	//entity := r.URL.Query().Get("e")

	if resource == "channels" {
		api := &apiResponse{
			Data: channels,
		}
		api.render(w)
		return
	}
}

type channel struct {
	Name string
	ID   string
}

var channels []channel

func getChannels() error {
	dat, err := ioutil.ReadFile("channels.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(dat, &channels)
	if err != nil {
		return err
	}
	if len(channels) < 1 {
		return errors.New("no channels found in channels.json")
	}
	return nil
}

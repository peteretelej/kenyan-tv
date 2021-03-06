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
	if err := LoadChannels(); err != nil {
		log.Fatalf("unable to get channels; %v", err)
	}
	log.Printf("Launching http server at %v", *listenAddr)
	log.Fatal(serve(*listenAddr))
}
func serve(listenAddr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/polymer", handlePolymer)
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
func renderHTML(w http.ResponseWriter, r *http.Request, filename string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.Header().Add("X-XSS-Protection", "1; mode=block")
	w.Header().Add("X-Frame-Options", "SAMEORIGIN")
	w.Header().Add("X-UA-Compatible", "IE=edge")
	http.ServeFile(w, r, filename)
}
func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	renderHTML(w, r, filepath.Join("vue", "index.html"))
}
func handlePolymer(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, r, filepath.Join("polymer", "index.html"))
}

// APIResponse is the /api response structure, rendered as json
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (a *APIResponse) render(w http.ResponseWriter) {
	if a == nil {
		a = &APIResponse{}
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
	fmt.Fprint(w, string(dat))
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
		api := &APIResponse{
			Data: channels,
		}
		api.render(w)
		return
	}
	if resource == "intlchannels" {
		api := &APIResponse{
			Data: intlchannels,
		}
		api.render(w)
		return
	}
}

type channel struct {
	Name       string `json:"name"`
	ID         string `json:"id"` // youtube channel ID
	Livestream string `json:"livestream"`
}

func (ch *channel) prep() {
	switch ch.Name {
	case "K24":
		ch.Livestream = "https://livestream.com/accounts/17606245/events/4832042/player?width=200&amp;height=200&amp;autoPlay=false&amp;mute=false"
	default:
		ch.Livestream = fmt.Sprintf("https://youtube.com/embed/live_stream?channel=%s&autoplay=1", ch.ID)
	}
}

var channels, intlchannels []channel

// LoadChannels loads channels from their respective json files
func LoadChannels() error {
	chs, err := getChannels("channels.json")
	if err != nil {
		return err
	}
	channels = chs
	log.Printf("%v", channels)
	if intlchannels, err = getChannels("intlchannels.json"); err != nil {
		return err
	}
	return nil

}
func getChannels(filename string) ([]channel, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var chs []channel
	err = json.Unmarshal(dat, &chs)
	if err != nil {
		return nil, err
	}
	if len(chs) < 1 {
		return nil, errors.New("no channels found in channels.json")
	}
	var channels []channel
	for _, ch := range chs {
		ch.prep()
		channels = append(channels, ch)
	}
	return channels, nil
}

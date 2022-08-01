package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"net/http"

	"go.uber.org/zap"
)

type NewURLRequester struct {
	Url string
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const HOST_URL = "https://yock.to"
type URLStore map[string]string

func (store *URLStore) GenerateURLKey() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	key := string(b)

	if _, ok := (*store)[key]; ok {
		//collision
		return store.GenerateURLKey()
	}
	return key
}

func (store *URLStore) apiHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	urlReq := &NewURLRequester{}
	if err := json.NewDecoder(req.Body).Decode(urlReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if urlReq.Url == "" {
		http.Error(w, "No URL provided", http.StatusInternalServerError)
		return
	}

	if len(urlReq.Url) > 64000 {
		http.Error(w, "Error: URL TOO LONG", http.StatusBadRequest)
		return
	}

	key := store.GenerateURLKey()

	(*store)[key] = urlReq.Url

	w.Header().Add("Content-Type", "application/json")

    final := fmt.Sprintf("%s/%s", HOST_URL, key)

	response := map[string]string{"url": final}
	json.NewEncoder(w).Encode(response)

}

func (store *URLStore) indexHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if url, ok := (*store)[req.URL.Path[1:]]; ok {
		http.Redirect(w, req, url, http.StatusSeeOther)
		return
	}

}

func main() {
	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatal(err)
	}

	sugar := logger.Sugar()

	sugar.Debug("Reading built-in FS")
	fsys, err := fs.Sub(res, "public")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(fsys))
	sugar.Debug("Reading built-in FS")
	mux := http.NewServeMux()

	store := URLStore{}
	mux.HandleFunc("/api/shorten", store.apiHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := store[r.URL.Path[1:]]; ok {
			store.indexHandler(w, r)
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})

	middleware := LoggingMiddleware(sugar)

	if err := http.ListenAndServe(":8090", middleware(mux)); err != nil {
		log.Panic(err)
	}
}

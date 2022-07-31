package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type NewURLRequester struct {
	Url string
}

func GenerateURLKey(url string) string {
	h := sha1.New()
	h.Write([]byte(url))
	bs := h.Sum(nil)

	final := hex.EncodeToString(bs)
	return final[:6]

}

type URLStore map[string]string

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

	key := GenerateURLKey(urlReq.Url)

	(*store)[key] = urlReq.Url

	w.Header().Add("Content-Type", "application/json")

	final := fmt.Sprintf("http://%s/%s", req.Host, key)

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
	fileServer := http.FileServer(http.Dir("./public"))

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

	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Panic(err)
	}
}

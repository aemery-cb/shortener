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

	w.Header().Add("Content-Type", "text/plain")

	final := fmt.Sprintf("https://%s/%s\n", req.Host, key)
	fmt.Fprintf(w, "%s\n", final)

}

func (store *URLStore) indexHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method != "GET" {
		return
	}

	if req.URL.Path == "" {
		fmt.Println("index probably")
		return
	}

	if url, ok := (*store)[req.URL.Path[1:]]; ok {
		http.Redirect(w, req, url, http.StatusSeeOther)
		return
	}

	http.Error(w, "Not Found", http.StatusNotFound)
}

func main() {

	store := URLStore{}
	http.HandleFunc("/api/shorten", store.apiHandler)
	http.HandleFunc("/", store.indexHandler)

	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Panic(err)
	}
}

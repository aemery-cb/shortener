package ui

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed public
var Public embed.FS
var EmbedUI = false

// Get the subtree of the embedded files with `build` directory as a root.
func BuildHTTPFS() http.FileSystem {

	if EmbedUI {

		build, err := fs.Sub(Public, "public")
		if err != nil {
			log.Fatal(err)
		}
		return http.FS(build)
	}
	fmt.Println("using os")
	return http.FS(os.DirFS("./ui/public"))
}

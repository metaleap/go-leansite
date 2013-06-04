package appengine

import (
	"net/http"
	"os"
	"path/filepath"

	leansite "github.com/metaleap/go-leansite"
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	leansite.Init(filepath.Dir(cwd))
	http.Handle("/", leansite.Router)
}

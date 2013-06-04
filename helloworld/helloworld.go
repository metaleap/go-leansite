package helloworld

import (
	"net/http"
	"os"

	leansite "github.com/metaleap/go-leansite"
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	leansite.Init(cwd)
	http.Handle("/", leansite.Router)
}

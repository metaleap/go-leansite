package helloworld

import (
	"net/http"

	leansite "github.com/metaleap/go-leansite"
)

func init() {
	leansite.Init("")
	http.Handle("/", leansite.Router)
}

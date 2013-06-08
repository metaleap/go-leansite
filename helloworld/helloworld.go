//	This is the app handler for running the helloworld Site on Google App Engine
package helloworld

import (
	"net/http"

	leansite "github.com/metaleap/go-leansite"
)

func init() {
	leansite.Init("")
	http.Handle("/", leansite.Router)
}

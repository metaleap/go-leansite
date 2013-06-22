// App Engine demo app at http://go-leansite-helloworld.appspot.com
package helloworld

import (
	"net/http"

	leansite "github.com/metaleap/go-leansite"
)

func init() {
	leansite.Init("")
	http.Handle("/", leansite.Router)
}

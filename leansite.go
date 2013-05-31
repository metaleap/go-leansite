package leansite

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"

	uio "github.com/metaleap/go-util/io"
)

var (
	DirPath  string
	DirWatch *uio.Watcher
	SiteData struct {
		TopNav NavItems
		Blogs  map[string]BlogNavItems
	}

	fileServer    http.Handler
	tmpl          *template.Template
	pageTemplates = map[string]*template.Template{}
)

func dir(names ...string) string {
	return filepath.Join(append([]string{DirPath}, names...)...)
}

func ListenAndServe(dirPath string) (err error) {
	SiteData.Blogs = map[string]BlogNavItems{}
	DirPath = dirPath
	if DirWatch, err = uio.NewWatcher(); err != nil {
		return
	} else {
		defer DirWatch.Close()
	}

	//	Load and watch templates
	DirWatch.WatchDir(dir("templates"), true, reloadTemplates)

	//	Load and watch top.nav
	DirWatch.WatchFiles(dir("contents"), "top.nav", true, func(filePath string) {
		data, err := ioutil.ReadFile(filePath)
		if err == nil {
			err = json.Unmarshal(data, &SiteData.TopNav)
		}
		if err != nil {
			SiteData.TopNav = append(NavItems{}, newNavItem("#", err.Error()))
		}
	})

	//	Activate *actual* file monitoring
	go DirWatch.Go()

	//	Listen and serve
	fileServer = http.FileServer(http.Dir(dir("static")))
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(serveTemplatedContent)
	s := &http.Server{
		Addr:    ":8080",
		Handler: r,
		//	http://stackoverflow.com/questions/10971800/golang-http-server-leaving-open-goroutines
		ReadTimeout: 2 * time.Minute,
	}
	return s.ListenAndServe()
}

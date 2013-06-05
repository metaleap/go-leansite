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
	Router   *mux.Router
	SiteData struct {
		TopNav NavItems
		Blogs  map[string]BlogNav

		mainTemplate  *template.Template
		pageTemplates map[string]*template.Template
	}

	fileServer http.Handler
)

func dir(names ...string) string {
	if len(DirPath) == 0 {
		return filepath.Join(names...)
	}
	return filepath.Join(append([]string{DirPath}, names...)...)
}

func Init(dirPath string) (err error) {
	SiteData.Blogs = map[string]BlogNav{}
	SiteData.pageTemplates = map[string]*template.Template{}
	DirPath = dirPath
	DirWatch, err = uio.NewWatcher()

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
	Router = mux.NewRouter()
	Router.PathPrefix("/").HandlerFunc(serveTemplatedContent)
	return
}

func ListenAndServe(addr string) (err error) {
	defer DirWatch.Close()
	s := &http.Server{
		Addr:    addr,
		Handler: Router,
		//	http://stackoverflow.com/questions/10971800/golang-http-server-leaving-open-goroutines
		ReadTimeout: 2 * time.Minute,
	}
	return s.ListenAndServe()
}

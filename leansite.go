package leansite

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"

	"github.com/metaleap/go-util/fs"
)

var (
	//	To be set via Init(), see func Init() docs
	DirPath string

	//	Various file-system watchers, initialized by Init()
	DirWatch *ufs.Watcher

	//	Our request router, initialized by Init()
	Router *mux.Router

	//	Some site-specific data loaded from DirPath
	SiteData struct {
		//	Top-level navigation links loaded from contents/top.nav JSON file
		TopNav NavItems

		//	A map of blogs. Populated by PageContext.GetBlogArchive(), which is called from a template
		Blogs map[string]BlogNav

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

//	Stand-alone: call this in your main() before calling ListenAndServe().
//
//	App Engine: call this in your init() and DON'T call ListenAndServe().
//
//	dirPath: the site's base directory which contains folders "static", "contents" and "templates".
func Init(dirPath string) (err error) {
	SiteData.Blogs = map[string]BlogNav{}
	SiteData.pageTemplates = map[string]*template.Template{}
	DirPath = dirPath
	DirWatch, err = ufs.NewWatcher()

	//	Load and watch templates
	DirWatch.WatchIn(dir("templates"), "*.html", true, reloadTemplates)

	//	Load and watch top.nav
	DirWatch.WatchIn(dir("contents"), "top.nav", true, func(filePath string) {
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

//	addr: see http.Server.Addr
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

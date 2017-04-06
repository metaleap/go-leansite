# leansite
--
    import "github.com/metaleap/go-leansite"

A super-minimalistic "dynamic-web-page" server written in Go, just to explore
the net/http, html/template, and Gorilla web packages.

The folder **helloworld** represents a sample LeanSite and can be used as a seed
site:

- run it standalone via `go run
$GOPATH/src/github.com/metaleap/go-leansite/helloworld/go-leansite-helloworld/main.go`

- run it on your Google App Engine SDK via `dev_appserver.py
$GOPATH/src/github.com/metaleap/go-leansite/helloworld`

- see it in action running on Google App Engine at
http://go-leansite-helloworld.appspot.com/

## Usage

```go
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
	}
)
```

#### func  Init

```go
func Init(dirPath string) (err error)
```
Stand-alone: call this in your main() before calling ListenAndServe().

App Engine: call this in your init() and DON'T call ListenAndServe().

dirPath: the site's base directory which contains folders "static", "contents"
and "templates".

#### func  ListenAndServe

```go
func ListenAndServe(addr string) (err error)
```
addr: see http.Server.Addr

#### type BlogNav

```go
type BlogNav struct {
	//	A chronological listing of blog entries
	Nav BlogNavItems
}
```

A chronological listing of blog entries

#### func (*BlogNav) ShowYear

```go
func (me *BlogNav) ShowYear(year string) (dif bool)
```
Returns true if year is different than the value passed when this method was
last called

#### type BlogNavItem

```go
type BlogNavItem struct {
	//	Embedded navigation info (Href, Caption)
	NavItem

	//	Date posted
	Year, Month, Day string
}
```

Represents a blog entry

#### type BlogNavItems

```go
type BlogNavItems []BlogNavItem
```

A collection of blog entries

#### func (BlogNavItems) Len

```go
func (me BlogNavItems) Len() int
```
Implements sort.Interface.Len()

#### func (BlogNavItems) Less

```go
func (me BlogNavItems) Less(i, j int) bool
```
Implements sort.Interface.Less()

#### func (BlogNavItems) Swap

```go
func (me BlogNavItems) Swap(i, j int)
```
Implements sort.Interface.Swap()

#### type NavItem

```go
type NavItem struct {
	//	Link attributes
	CssClass, Href, Caption, Desc string
}
```

A navigation link such as used in the sample site's top-navbar

#### func (*NavItem) IsActive

```go
func (me *NavItem) IsActive(pc *PageContext) bool
```
Returns true if this NavItem points to the resource represented by pc, or an
ancestor resource.

#### type NavItems

```go
type NavItems []NavItem
```

A collection of navigation links

#### type PageContext

```go
type PageContext struct {
	//	The underlying http.Request
	R *http.Request

	//	Host-relative Request URL path
	Path string

	//	Extracted from either SiteData.TopNav or the first <h2> occurence in the final HTML output
	PageTitle string

	//	Final HTML output
	HtmlContent template.HTML

	//	Always equivalent to SiteData.TopNav
	TopNav NavItems

	//	Contains the current year, for auto-updating copyright notices
	Year int
}
```

Represents a single page request, used as pipeline in main template

#### func  NewPageContext

```go
func NewPageContext(r *http.Request, path string) (me *PageContext)
```
Creates a new PageContext for the specified http.Request and host-relative
request URL path

#### func (*PageContext) GetBlogArchive

```go
func (me *PageContext) GetBlogArchive(path string) *BlogNav
```
Returns a BlogNav for the specified path. For example, GetBlogArchive("blog")
maps to "contents/blog/"

--
**godocdown** http://github.com/robertkrimen/godocdown

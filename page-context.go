package leansite

import (
	"html/template"
	"net/http"
	"strings"
	"time"
)

//	A collection of navigation links
type NavItems []NavItem

//	Represents a single page request, used as pipeline in main template
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

//	Creates a new PageContext for the specified http.Request and host-relative request URL path
func NewPageContext(r *http.Request, path string) (me *PageContext) {
	me = &PageContext{R: r, Path: path, Year: time.Now().Year()}
	me.TopNav = SiteData.TopNav
	for _, nav := range me.TopNav {
		if nav.IsActive(me) {
			me.PageTitle = nav.Caption
		}
	}
	return
}

//	A navigation link such as used in the sample site's top-navbar
type NavItem struct {
	//	Link attributes
	CssClass, Href, Caption, Desc string
}

func newNavItem(href, caption string) (me NavItem) {
	me = NavItem{Href: href, Caption: caption}
	return
}

//	Returns true if this NavItem points to the resource represented by pc, or an ancestor resource.
func (me *NavItem) IsActive(pc *PageContext) bool {
	if len(me.Href) == 0 {
		return len(pc.Path) == 0
	}
	return strings.HasPrefix(pc.Path, me.Href)
}

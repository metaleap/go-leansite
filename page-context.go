package leansite

import (
	"html/template"
	"net/http"
	"strings"
	"time"
)

type NavItems []NavItem

type PageContext struct {
	R               *http.Request
	Path, PageTitle string
	HtmlContent     template.HTML
	TopNav          NavItems
	Year            int
}

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

type NavItem struct {
	CssClass, Href, Caption, Desc string
}

func newNavItem(href, caption string) (me NavItem) {
	me = NavItem{Href: href, Caption: caption}
	return
}

func (me *NavItem) IsActive(pc *PageContext) bool {
	if len(me.Href) == 0 {
		return len(pc.Path) == 0
	}
	return strings.HasPrefix(pc.Path, me.Href)
}

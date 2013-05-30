package leansite

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	siteTopNav []*PageTopNavItem
)

type PageContext struct {
	R           *http.Request
	Path        string
	HtmlContent template.HTML
	TopNav      []*PageTopNavItem
}

func NewPageContext(r *http.Request, path string) (me *PageContext) {
	me = &PageContext{R: r, Path: path}
	if len(siteTopNav) == 0 {
		data, err := ioutil.ReadFile(filepath.Join(dir("contents"), "top.nav"))
		if err == nil {
			err = json.Unmarshal(data, &siteTopNav)
		}
		if err != nil {
			siteTopNav = append(siteTopNav, newPageTopNavItem("#", "Nav Error", err.Error(), ""))
		}
	}
	me.TopNav = siteTopNav
	return
}

func (me *PageContext) GetBlog(path string) (entries []*PageTopNavItem) {
	return
}

type PageTopNavItem struct {
	CssClass, Href, Caption, Desc string
}

func newPageTopNavItem(href, caption, desc, cssClass string) (me *PageTopNavItem) {
	me = &PageTopNavItem{Href: href, Caption: caption, CssClass: cssClass, Desc: desc}
	return
}

func (me *PageTopNavItem) IsActive(pc *PageContext) bool {
	if len(me.Href) == 0 {
		return len(pc.Path) == 0
	}
	return strings.HasPrefix(pc.Path, me.Href)
}

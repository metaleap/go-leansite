package leansite

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	uio "github.com/metaleap/go-util/io"
)

type BlogNavItems []BlogNavItem

func (me BlogNavItems) Len() int { return len(me) }

func (me BlogNavItems) Less(i, j int) bool { return me[j].Href < me[i].Href }

func (me BlogNavItems) Swap(i, j int) { me[i], me[j] = me[j], me[i] }

type BlogNavItem struct {
	NavItem
	Year, Month, Day string
}

func (me *PageContext) GetBlogArchive(path string) BlogNavItems {
	if _, ok := SiteData.Blogs[path]; !ok {
		dirPath := dir("contents", path)
		handler := func() {
			items := BlogNavItems{}
			uio.NewDirWalker(true, nil, func(_ *uio.DirWalker, fullPath string, _ os.FileInfo) bool {
				if filepath.Dir(fullPath) != dirPath {
					vpath := fullPath[:len(fullPath)-len(filepath.Ext(fullPath))]
					vpath = vpath[len(dirPath):]
					navItem := BlogNavItem{}
					navItem.Href, navItem.Caption = filepath.ToSlash(vpath), filepath.Base(vpath)
					if src := uio.ReadTextFile(fullPath, false, ""); len(src) > 0 {
						if pos1, pos2 := strings.Index(src, "<h2>"), strings.Index(src, "</h2>"); pos1 >= 0 && pos2 > pos1 {
							src = src[:pos2]
							navItem.Caption = src[pos1+4:]
						}
					}
					if pathItems := strings.Split(navItem.Href, "/"); len(pathItems) > 0 {
						if navItem.Year = pathItems[1]; len(pathItems) > 1 {
							if navItem.Month = pathItems[2]; len(pathItems) > 2 {
								navItem.Day = pathItems[3]
							}
						}
					}
					items = append(items, navItem)
				}
				return true
			}).Walk(dirPath)
			sort.Sort(items)
			SiteData.Blogs[path] = items
		}
		uio.NewDirWalker(true, func(_ *uio.DirWalker, fullPath string, _ os.FileInfo) bool {
			DirWatch.WatchDir(fullPath, false, handler)
			return true
		}, nil).Walk(dirPath)
		handler()
	}
	return SiteData.Blogs[path]
}

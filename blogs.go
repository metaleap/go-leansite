package leansite

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/metaleap/go-util/fs"
)

//	A collection of blog entries
type BlogNavItems []BlogNavItem

//	Implements sort.Interface.Len()
func (me BlogNavItems) Len() int { return len(me) }

//	Implements sort.Interface.Less()
func (me BlogNavItems) Less(i, j int) bool { return me[j].Href < me[i].Href }

//	Implements sort.Interface.Swap()
func (me BlogNavItems) Swap(i, j int) { me[i], me[j] = me[j], me[i] }

//	Represents a blog entry
type BlogNavItem struct {
	//	Embedded navigation info (Href, Caption)
	NavItem

	//	Date posted
	Year, Month, Day string
}

//	Returns true if year is different than the value passed when this method was last called
func (me *BlogNav) ShowYear(year string) (dif bool) {
	if dif = (year != me.lastYear); dif {
		me.lastYear = year
	}
	return
}

//	A chronological listing of blog entries
type BlogNav struct {
	//	A chronological listing of blog entries
	Nav BlogNavItems

	lastYear string
}

//	Returns a BlogNav for the specified path.
//	For example, GetBlogArchive("blog") maps to "contents/blog/"
func (me *PageContext) GetBlogArchive(path string) *BlogNav {
	if _, ok := SiteData.Blogs[path]; !ok {
		dirPath := dir("contents", path)
		handler := func(_ string) {
			items := BlogNavItems{}
			ufs.WalkAllFiles(dirPath, func(fullPath string) bool {
				if filepath.Dir(fullPath) != dirPath {
					vpath := fullPath[:len(fullPath)-len(filepath.Ext(fullPath))]
					vpath = vpath[len(dirPath):]
					navItem := BlogNavItem{}
					navItem.Href, navItem.Caption = filepath.ToSlash(vpath), filepath.Base(vpath)
					if src := ufs.ReadTextFile(fullPath, false, ""); len(src) > 0 {
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
			})
			sort.Sort(items)
			SiteData.Blogs[path] = BlogNav{Nav: items}
		}
		ufs.WalkAllDirs(dirPath, func(fullPath string) bool {
			DirWatch.WatchIn(fullPath, "*", false, handler)
			return true
		})
		handler("")
	}
	copy := SiteData.Blogs[path]
	return &copy
}

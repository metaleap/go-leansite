package leansite

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	markdown "github.com/russross/blackfriday"

	"github.com/metaleap/go-util-fs"
	"github.com/metaleap/go-util-str"
)

func reloadTemplates(_ string) {
	fileNames := []string{filepath.Join(dir("templates"), "main.html")}
	ufs.WalkFilesIn(dir("templates"), func(fullPath string) bool {
		if !strings.HasSuffix(fullPath, string(filepath.Separator)+"main.html") {
			fileNames = append(fileNames, fullPath)
		}
		return true
	})
	var err error
	SiteData.mainTemplate, err = template.ParseFiles(fileNames...)
	if err != nil {
		SiteData.mainTemplate, err = template.New("error").Parse(fmt.Sprintf("ERROR loading templates: %+v", err))
	}
}

func serveTemplatedContent(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.Trim(r.URL.Path, "/")
	//	First handle static files (robots.txt / sitemap.xml / favicon.ico etc.) etc via 'static' folder
	if ufs.FileExists(filepath.Join(dir("static"), urlPath)) {
		fileServer.ServeHTTP(w, r)
		return
	}

	var (
		err        error
		filePath   string
		fileData   []byte
		isMarkdown bool
	)
	if filePath = filepath.Join(dir("contents"), urlPath) + ".html"; !ufs.FileExists(filePath) {
		if filePath = filepath.Join(dir("contents"), urlPath, "index.html"); !ufs.FileExists(filePath) {
			isMarkdown = true
			if filePath = filepath.Join(dir("contents"), urlPath) + ".md"; !ufs.FileExists(filePath) {
				filePath = filepath.Join(dir("contents"), urlPath, "index.md")
			}
		}
	}
	pc := NewPageContext(r, urlPath)
	if len(filePath) > 0 && ufs.FileExists(filePath) {
		if fileData, err = ioutil.ReadFile(filePath); err == nil {
			var tmpl *template.Template
			if pos := bytes.Index(fileData, []byte("{{")); pos >= 0 {
				if bytes.Index(fileData, []byte("}}")) > pos {
					if _, ok := SiteData.pageTemplates[filePath]; !ok {
						SiteData.pageTemplates[filePath] = nil
						DirWatch.WatchIn(filepath.Dir(filePath), ustr.Pattern(filepath.Base(filePath)), false, func(fullPath string) {
							SiteData.pageTemplates[fullPath] = nil
						})
					}
					if tmpl = SiteData.pageTemplates[filePath]; tmpl == nil {
						var err error
						tmpl, err = template.ParseFiles(filePath)
						if err == nil {
							SiteData.pageTemplates[filePath] = tmpl
						} else {
							tmpl, err = template.New("pt_error").Parse(fmt.Sprintf("ERROR loading template %s:\t%+v", filePath, err))
						}
					}
				}
			}
			if tmpl != nil {
				var buf bytes.Buffer
				if err := tmpl.Execute(&buf, pc); err != nil {
					fileData = []byte(err.Error())
				} else {
					fileData = buf.Bytes()
				}
			}
			if isMarkdown {
				pc.HtmlContent = template.HTML(markdown.MarkdownCommon(fileData))
			} else {
				pc.HtmlContent = template.HTML(fileData)
			}
		}
	} else {
		pc.HtmlContent = "404 Not Found"
	}
	pos1, pos2 := strings.Index(string(pc.HtmlContent), "<h2>"), strings.LastIndex(string(pc.HtmlContent), "</h2>")
	if pos1 >= 0 && pos2 > pos1 {
		docTitle := string(pc.HtmlContent[:pos2])
		if pos2 = strings.Index(docTitle, "</h2>"); pos2 < 0 {
			pc.PageTitle = docTitle[pos1+4:]
		}
	}
	if err == nil {
		err = SiteData.mainTemplate.Execute(w, pc)
	} else {
		w.Write([]byte(err.Error()))
	}
}

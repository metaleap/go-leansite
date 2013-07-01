// Standalone web server
package main

import (
	"flag"
	"log"

	"github.com/go-utils/ugo"

	leansite "github.com/metaleap/go-leansite"
)

func main() {
	flag.Parse()
	dirPath := *flag.String("dir", ugo.GopathSrcGithub("metaleap", "go-leansite", "helloworld"), "Root directory path containing the static, contents, templates etc. folders.")
	leansite.Init(dirPath)
	log.Println("Listening any moment now...")
	log.Fatal(leansite.ListenAndServe(":8008"))
}

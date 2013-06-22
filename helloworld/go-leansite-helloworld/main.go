// Standalone web server
package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/go-utils/ugo"

	leansite "github.com/metaleap/go-leansite"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	dirPath := *flag.String("dir", ugo.GopathSrcGithub("metaleap", "go-leansite", "helloworld"), "Root directory path containing the static, contents, templates etc. folders.")
	leansite.Init(dirPath)
	log.Fatal(leansite.ListenAndServe(":8008"))
}

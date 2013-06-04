package main

import (
	"flag"
	"log"
	"runtime"

	leansite "github.com/metaleap/go-leansite"
	ugo "github.com/metaleap/go-util"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	dirPath := *flag.String("dir", ugo.GopathSrcGithub("metaleap", "go-leansite", "helloworld"), "Root directory path containing the static, contents, templates etc. folders.")
	leansite.Init(dirPath)
	log.Fatal(leansite.ListenAndServe(":8008"))
}

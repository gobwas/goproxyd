package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/goproxy/goproxy"
)

func main() {
	var (
		root     string
		addr     string
		readonly bool
	)
	flag.StringVar(&root,
		"root", "/cache/download",
		"cache download root",
	)
	flag.BoolVar(&readonly,
		"ro", true,
		"readonly mode",
	)
	flag.StringVar(&addr,
		"addr", "0.0.0.0:8080",
		"addr to bind to",
	)
	flag.Parse()

	proxy := goproxy.New()
	proxy.Cacher = &cacher{
		root:     root,
		readonly: readonly,
	}

	log.Println("listening", addr)
	log.Fatal(http.ListenAndServe(addr, logHandler(proxy)))
}

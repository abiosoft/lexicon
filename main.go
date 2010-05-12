// Copyright 2010 Abiola Ibrahim. All rights reserved.
// CSC 332 Lexicon Assignment
// Analysis of Java Programming Language using Go (http://golang.org)
package main

import (
	"http"
	"io"
	"io/ioutil"
	"strings"
	"lexer"
	"os"
)

func Server(c *http.Conn, req *http.Request) {
	path := req.FormValue("file")
	file, err := os.Open(path, os.O_RDONLY, 0666)
	lex := new(lexer.Lexicon)
	if err != nil{
		lex.Init("hello.java")
	}else{
		file.Close()
		lex.Init(path)
	}
	lex.FilterStrings()
	lex.AddSpaces()
	lex.Analyze()
	io.WriteString(c, lex.Html())
}

func ProcessServer(c *http.Conn, req *http.Request) {
	path := req.FormValue("file")
	file, err := os.Open(path, os.O_RDONLY, 0666)
	if err != nil{
		path = "hello.java"
	}else{
		file.Close()
	}
	data, _ := ioutil.ReadFile(path)
	io.WriteString(c, string(data))
}

func FileServer(c *http.Conn, req *http.Request) {
	if req.URL.Path == "/" {
		http.Redirect(c, "index.html", 307)
		return
	} else if strings.HasSuffix(req.URL.Path, "html") {
		data, _ := ioutil.ReadFile(req.URL.Path[1:])
		io.WriteString(c, string(data))
		return
	}
	http.ServeFile(c, req, req.URL.Path[1:])
}

func main() {
	http.HandleFunc("/process", Server)
	http.HandleFunc("/file", ProcessServer)
	http.HandleFunc("/", FileServer)
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		println("Fatal error occured, program needs to close")
	}
}

package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"path"
	"strings"

	"github.com/russross/blackfriday"
)

type FastCGIServer struct{}

var outputTemplate = template.Must(template.New("base").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>{{ .Title }}</title>
</head>
<body>
	{{ .Body }}
</body>
</html>`))


func (s FastCGIServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if !strings.HasSuffix(req.URL.Path, ".md") {
		http.Error(resp, "Not Found", 404)
		return
	}

	fcgiEnv := fcgi.ProcessEnv(req)
	fullPath := path.Join(fcgiEnv["DOCUMENT_ROOT"], fcgiEnv["DOCUMENT_URI"])
	_, fileName := path.Split(req.URL.Path)


	_, err := os.Stat(fullPath)
	if err != nil {
		http.Error(resp, "Not Found", 404)
		return
	}

	input, err := ioutil.ReadFile(fullPath)
	if err != nil {
		http.Error(resp, "Internal Server Error", 500)
		return
	}

	output := blackfriday.Run(input)

	resp.Header().Set("Content-Type", "text/html")
	outputTemplate.Execute(resp, struct {
		Title string
		Body template.HTML
	}{
		Title: fmt.Sprintf("%s - [%s]", fileName, req.URL.Path),
		Body: template.HTML(string(output)),
	})
	fmt.Printf("GET %s 200\n", req.URL.Path)
}


func main() {
	fmt.Println("Markdown FastCGI Service Listen at 127.0.0.1:9001")
	listener, _ := net.Listen("tcp", "127.0.0.1:9001")
	srv := new(FastCGIServer)
	fcgi.Serve(listener, srv)
}


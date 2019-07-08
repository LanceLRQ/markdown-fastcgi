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
	"regexp"
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

	fcgiEnv := fcgi.ProcessEnv(req)
	fullPath := path.Join(fcgiEnv["DOCUMENT_ROOT"], fcgiEnv["DOCUMENT_URI"])
	_, fileName := path.Split(fcgiEnv["DOCUMENT_URI"])


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
	err = outputTemplate.Execute(resp, struct {
		Title string
		Body template.HTML
	}{
		Title: fmt.Sprintf("%s - [%s]", fileName, fcgiEnv["DOCUMENT_URI"]),
		Body: template.HTML(string(output)),
	})
	if err != nil {
		http.Error(resp, "Internal Server Error", 500)
		return
	}

	fmt.Printf("GET %s 200\n", req.URL.Path)
}


func main() {
	networkMode := "tcp"
	hostAddr := "127.0.0.1:9001"
	for i := 0; i < len(os.Args); i++ {
		p := os.Args[i]
		switch p {
		case "-l", "--listen":
			if i+1 >= len(os.Args) {
				break
			}
			i++
			addr := os.Args[i]
			indexUnix := strings.Index(addr, "unix://")
			if indexUnix > -1 {
				hostAddr = addr[indexUnix + 6:]
				networkMode = "unix"
			} else {
				a, err := regexp.MatchString("((?:(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d)))\\.){3}(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d))))(:\\d)", addr)
				if !a || err != nil {
					panic("Wrong IP Address!")
				}
				networkMode = "tcp"
				hostAddr = addr
			}
		}
	}

	fmt.Printf("Markdown FastCGI Service Listen at %s://%s\n", networkMode, hostAddr)

	listener, _ := net.Listen(networkMode, hostAddr)
	srv := new(FastCGIServer)
	err := fcgi.Serve(listener, srv)
	if err != nil {
		panic(err)
	}
}


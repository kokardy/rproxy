package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
)

type Director func(*http.Request)

func NewDirector(scheme, host string) (dir Director) {

	dir = func(request *http.Request) {
		request.URL.Scheme = scheme
		request.URL.Host = host
	}
	return
}

type Modifier func(*http.Response) error

type Converter func(io.Reader) io.ReadCloser

func NewRegConverter(ori, dest string) (c Converter) {
	reg1 := regexp.MustCompile(ori)
	c = func(r io.Reader) (rc io.ReadCloser) {
		buf := bytes.NewBuffer(nil)
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			line := sc.Text()
			new_line := reg1.ReplaceAllString(line, dest)
			buf.Write([]byte(new_line + "\n"))
		}
		rc = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
		return
	}
	return
}

func NewModifier(c Converter) (mod Modifier) {
	mod = func(res *http.Response) error {
		for _, ct := range res.Header["Content-type"] {
			if ct == "text/html" {
				original := res.Body
				modified := c(original)
				res.Body = modified
				return nil
			}
		}
		return nil
	}
	return
}

func main() {
	//rproxy http host listenport original replace

	var remoteScheme string
	var remoteHost string
	var addr string
	var ori, dest string

	flag.StringVar(&remoteScheme, "scheme", "http", "remote scheme: -scheme http or -s https. http default")
	flag.StringVar(&remoteHost, "rhost", "127.0.0.1:80", "remote address: -rhost www.google.com:80  default 127.0.0.1:80")
	flag.StringVar(&addr, "addr", ":8080", "address: -addr :8080 default 8080")
	flag.StringVar(&ori, "ori", "", `original Regexp: -ori href=\\"https?://(.*)/\\"`)
	flag.StringVar(&dest, "dest", "", `modifiled: -dest href=\\"https?://$1/\\"`)

	flag.Parse()

	fmt.Println("a", addr)

	dir := NewDirector(remoteScheme, remoteHost)
	conv := NewRegConverter(ori, dest)
	mod := NewModifier(conv)
	rp := &httputil.ReverseProxy{
		Director:       dir,
		ModifyResponse: mod,
	}
	server := http.Server{
		Addr:    addr,
		Handler: rp,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}

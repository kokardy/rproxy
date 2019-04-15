package main

import (
	"bufio"
	"bytes"
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
			buf.Write([]byte(new_line))
		}
		rc = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
		return
	}
	return
}

func NewModifier(c Converter) (mod Modifier) {
	mod = func(res *http.Response) error {

		original := res.Body
		modified := c(original)
		res.Body = modified
		return nil
	}
	return
}

func main() {
	//rproxy http host listenport original replace

	var scheme string
	var host string
	var listenport string
	var ori, dest string

	dir := NewDirector(scheme, host)
	conv := NewRegConverter(ori, dest)
	mod := NewModifier(conv)
	rp := &httputil.ReverseProxy{
		Director:       dir,
		ModifyResponse: mod,
	}
	server := http.Server{
		Addr:    listenport,
		Handler: rp,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}

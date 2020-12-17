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
	"strconv"
	"strings"
)

//Director はリクエストの向きを変えてリバースプロキシのターゲットを設定する
type Director func(*http.Request)

//NewDirector は新しいDirectorをスキームとホスト名を指定して生成する
func NewDirector(scheme, host string) (dir Director) {
	dir = func(request *http.Request) {
		request.URL.Scheme = scheme
		request.URL.Host = host
	}
	return
}

//Converter はio.Readerから読み込んで中身を書き換えてio.ReadCloserにして返す
type Converter func(io.Reader) (io.ReadCloser, int)

//NewRegConverter 正規表現で中身を書き換えるConveterを生成する
func NewRegConverter(ori, dest string) (c Converter) {
	reg1 := regexp.MustCompile(ori)
	c = func(r io.Reader) (rc io.ReadCloser, contentLength int) {
		buf := bytes.NewBuffer(nil)
		sc := bufio.NewReader(r)
		preline := make([]byte, 0, 0)
		for {
			line, isPrefix, err := sc.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err.Error())
			}
			if isPrefix {
				preline = append(preline, line...)
				continue
			} else {
				line = append(preline, line...)
				preline = []byte{}

				new_line := reg1.ReplaceAllString(string(line), dest)
				buf.Write([]byte(new_line + "\n"))
			}
		}
		rc = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
		contentLength = len(buf.Bytes())
		return
	}
	return
}

//Modifier はリバースプロキシで中継するレスポンスを書き換える
type Modifier func(*http.Response) error

//NewModifier はModifierをConverterから生成する
func NewModifier(c Converter) (mod Modifier) {
	mod = func(res *http.Response) error {
		for _, ct := range res.Header["Content-Type"] {
			if strings.Contains(ct, "text/html") {
				original := res.Body
				modified, contentLength := c(original)
				res.Body = modified
				res.Header.Set("Content-Length", strconv.Itoa(contentLength))
				return nil
			}
		}
		return nil
	}
	return
}

func main() {

	//*****************コマンドラインオプション START
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
	//******************コマンドラインオプション END

	dir := NewDirector(remoteScheme, remoteHost)
	conv := NewRegConverter(ori, dest)
	mod := NewModifier(conv)
	//リバースプロキシ生成
	rp := &httputil.ReverseProxy{
		Director:       dir,
		ModifyResponse: mod,
	}

	//Server生成してリバースプロキシをHandlerとしてセットする
	server := http.Server{
		Addr:    addr,
		Handler: rp,
	}

	//Listenしているポートを表示する
	fmt.Println("Listen: ", addr)
	fmt.Println("Press Ctrl+C to stop this server.")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}

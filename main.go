package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

var CACHE_DIR string
var HOST = "127.0.0.1"
var PORT = 9988

var seq = NewSequence()
var notFound = []byte("not found")

func init() {
	gopath, _ := os.LookupEnv("GOPATH")
	if gopath == "" {
		panic("the env GOPATH is empty, please set it and rerun again")
	}
	log.Printf("Initialization completed, GOPATH is %s\n", gopath)
	CACHE_DIR = strings.ReplaceAll(gopath, "\\", "/") + "/pkg/mod/cache/download"
	log.Printf("Please set GOPROXY to http://%s:%d\n", HOST, PORT)
}

func main() {
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", HOST, PORT),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uri, err := url.PathUnescape(r.RequestURI)
			if err != nil {
				panic(err)
			}
			var buffer = bufio.NewWriter(os.Stderr)
			defer buffer.Flush()
			logger := log.New(buffer, seq.Next()+"-", 1|2)
			logger.Printf("<- %s", uri)
			var resp = make([]byte, 0)
			switch {
			case strings.Contains(uri, "/@v"):
				if !strings.HasSuffix(uri, "list") && !strings.HasSuffix(uri, "info") &&
					!strings.HasSuffix(uri, "mod") && !strings.HasSuffix(uri, "zip") {
					w.WriteHeader(http.StatusNotFound)
					break
				}
				var dest = CACHE_DIR + uri
				if info, err := os.Stat(dest); err != nil || info == nil {
					w.WriteHeader(http.StatusNotFound)
					resp = notFound
					break
				}
				resp, err = os.ReadFile(dest)
				if err != nil {
					panic(err)
				}
			case strings.HasSuffix(uri, "/@latest"):
				var dest = CACHE_DIR + strings.TrimSuffix(uri, "/@latest")
				info, err := os.Stat(dest + "/@v/list")
				if err != nil || info == nil {
					w.WriteHeader(http.StatusNotFound)
					resp = notFound
					break
				}
				b, err := os.ReadFile(dest + "/@v/list")
				if err != nil {
					panic(err)
				}
				vs := strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
				var ver string
				for i := len(vs) - 1; i >= 0; i-- {
					info, err := os.Stat(dest + "/@v/" + vs[i] + ".info")
					if err == nil && info != nil {
						ver = vs[i]
						break
					}
				}
				if ver != "" {
					w.WriteHeader(http.StatusNotFound)
					resp = notFound
					break
				}
				resp, err = os.ReadFile(dest + "/@v/" + ver + ".info")
				if err != nil {
					panic(err)
				}
			default:
				w.WriteHeader(http.StatusNotFound)
			}
			w.Write(resp)
			if len(resp) > 512 {
				logger.Printf("-> %s", "blob")
			} else {
				logger.Printf("-> %s", strings.ReplaceAll(string(resp), "\n", " "))
			}
		}),
	)
	if err != nil {
		panic(err)
	}
}

type sequence struct {
	sync.Mutex
	id int64
}

func NewSequence() *sequence {
	return &sequence{}
}

func (seq *sequence) Next() string {
	seq.Lock()
	defer seq.Unlock()
	if seq.id < 99999999 {
		seq.id++
	} else {
		seq.id = 1
	}
	return fmt.Sprintf("%08d", seq.id)
}

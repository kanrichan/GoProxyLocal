package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var BASE_DIR string
var HOST = "127.0.0.1"
var PORT = 8080

func init() {
	gopath, _ := os.LookupEnv("GOPATH")
	if gopath == "" {
		panic("the env GOPATH is empty, please set it and rerun again")
	}
	BASE_DIR = strings.ReplaceAll(gopath, "\\", "/") + "/pkg/mod"
	log.Printf("Initialization completed, BASE_DIR is %s\n", BASE_DIR)
	log.Printf("Please set GOPROXY to http://%s:%d\n", HOST, PORT)
}

func main() {
	http.ListenAndServe(fmt.Sprintf("%s:%d", HOST, PORT),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uri, err := url.PathUnescape(r.RequestURI)
			if err != nil {
				panic(err)
			}
			log.Printf("<- %s", uri)
			switch {
			case strings.HasSuffix(uri, "@v/list"):
				repo := strings.Split(uri, `/@v`)[0]
				list, err := version(repo[1:])
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(fmt.Sprintf(`not found: module %s: no matching versions for query "latest"`, repo[1:])))
					return
				}
				w.Write([]byte(strings.Join(list, "\n")))
				return
			case strings.HasSuffix(uri, "/@latest"):
				repo := strings.Split(uri, `/@latest`)[0]
				list, err := version(repo[1:])
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(fmt.Sprintf(`not found: module %s: no matching versions for query "latest"`, repo[1:])))
					return
				}
				w.Write([]byte(fmt.Sprintf(`{"Version":"%s","Time":"2021-10-08T14:36:13Z"}`, list[len(list)-1])))
				return
			}
		}),
	)
}

func version(module string) ([]string, error) {
	repoelem := strings.Split(module, `/`)
	parent := strings.Join(repoelem[:len(repoelem)-1], `/`)
	name := repoelem[len(repoelem)-1]
	fd, err := os.ReadDir(BASE_DIR + "/" + parent)
	if err == nil {
		var list = make([]string, 0)
		for i := range fd {
			if !fd[i].IsDir() {
				continue
			}
			if !strings.HasPrefix(fd[i].Name(), name) {
				continue
			}
			if !strings.Contains(fd[i].Name(), "@") {
				continue
			}
			list = append(list, strings.Split(fd[i].Name(), "@")[1])
		}
		if len(list) > 0 {
			return list, nil
		}
	}
	return nil, errors.New("not found")
}

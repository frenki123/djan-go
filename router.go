package godjan

import (
	"fmt"
	"log"
	"net/http"
)

type Router struct {
	mux    *http.ServeMux
	prefix string
}

func NewRouter(prefix string) Router {
	if prefix != "" {
		if prefix[0] != '/' {
			log.Fatalf("Prefix '%s' doesn't start with char '/'", prefix)
		}
	}
	mux := http.NewServeMux()
	router := Router{
		mux:    mux,
		prefix: prefix,
	}
	return router
}

func RouterWithUrls(prefix string, urlpaterns []Path) Router {
	router := NewRouter(prefix)
	for _, url := range urlpaterns {
		router.GET(url.path, url.handler)
	}
	return router
}

func (router *Router) GET(url string, handleFunc http.HandlerFunc) {
	router.addHandler(http.MethodGet, url, handleFunc)
}
func (router *Router) POST(url string, handleFunc http.HandlerFunc) {
	router.addHandler(http.MethodPost, url, handleFunc)
}

func (router *Router) addHandler(method string, url string, handleFunc http.HandlerFunc) {
	if url[0] != '/' {
		log.Fatalf("URL '%s' doesn't start with char '/'", url)
	}
	url = fmt.Sprintf("%s%s", router.prefix, url)
	if url == "/" {
		url = "/{$}"
	}
	url = fmt.Sprintf("%s %s", method, url)
	router.mux.HandleFunc(url, handleFunc)
}

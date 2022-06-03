package docs

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

type DocsEndpoint struct {
	AllRoutes []RouteInfo
	Name      string
	Host      string
}
type RouteInfo struct {
	Path    string
	URL     string
	Methods []string
	Name    string
}

var docsData DocsEndpoint

func docsEndpoint(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("docs/template.html")
	if err != nil {
		http.Error(w, "can't get docs template", 500)
		return
	}

	err = tmpl.Execute(w, docsData)
	if err != nil {
		http.Error(w, fmt.Sprintf("can't get execute docs template with data %+v", docsData), 500)
		return
	}
}

func Middleware(router *mux.Router, appHost string, docsPathOptional ...string) func(http.Handler) http.Handler {
	appURL := getURLFromHost(appHost)

	docsData = DocsEndpoint{
		Host: appURL,
		Name: appHost,
	}
	router.Walk(func(route *mux.Route, router *mux.Router, ancsestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		name := route.GetName()
		queries, _ := route.GetQueriesTemplates()
		queriesString := formatQueriesTemplates(queries)

		docsData.AllRoutes = append(docsData.AllRoutes, RouteInfo{
			Path:    path + queriesString,
			URL:     appURL + path + queriesString,
			Methods: methods,
			Name:    name,
		})
		return nil
	})

	docsPath := getOptionalParam(docsPathOptional, "/docs")
	router.HandleFunc(docsPath, docsEndpoint).Methods("GET")
	log.Printf("[Docs] docs path %s", appURL+docsPath)

	// log.Printf("[Docs] routes inited %+v\n", docsData)
	return func(h http.Handler) http.Handler {
		return h
	}
}

func getURLFromHost(host string) string {
	url := host
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	return url
}

func formatQueriesTemplates(queries []string) string {
	if len(queries) != 0 {
		queries[0] = "?" + queries[0]
	}
	return strings.Join(queries, "&")
}

func getOptionalParam(params []string, def string) string {
	if len(params) == 0 {
		params = append(params, def)
	}
	return params[0]
}

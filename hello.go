package hello

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"encoding/json"

	"github.com/golang/gddo/httputil/header"
	gae "google.golang.org/appengine"
)

const htmlIndexTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>Yuksnort!</title>
    <meta name="description" content="">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta/css/bootstrap.min.css" integrity="sha384-/Y6pD6FV/Vv2HJnA6t+vslU6fwYXjCFtcEpHbNJ0lyAFsXTsjBbfaDjzALeQsN6M" crossorigin="anonymous">
    <link rel="stylesheet" href="/css/app.css">
</head>
{{ define "leftCol" }}<div class="col-2">{{ . }}</div>{{ end }}
{{ define "rightCol" }}<div class="col">{{ . }}</div>{{ end }}
<body>
    <div class="container-fluid">
    	<div class="row header">
    		Google App Engine
    	</div>
    	<div class="row">
    		{{ template "leftCol" "AppID" }}{{ template "rightCol" .AppId }}
    	</div>
    	<div class="row">
    		{{ template "leftCol" "Datacenter" }}{{ template "rightCol" .Datacenter }}
    	</div>
    	<div class="row">
    		{{ template "leftCol" "InstanceID" }}{{ template "rightCol" .InstanceId }}
    	</div>
    	<div class="row">
    		{{ template "leftCol" "VersionID" }}{{ template "rightCol" .VersionId }}
    	</div>

		<div class="row header">Request</div>
    	<div class="row">
    		{{ template "leftCol" "Host" }}{{ template "rightCol" .Request.Host }}
    	</div>
    	<div class="row">
    		{{ template "leftCol" "Method" }}{{ template "rightCol" .Request.Method }}
    	</div>
    	<div class="row">
    		{{ template "leftCol" "Path" }}{{ template "rightCol" .Request.URL.Path }}
    	</div>
    	<div class="row">
    		{{ template "leftCol" "Proto" }}{{ template "rightCol" .Request.Proto }}
    	</div>
    	<div class="row">
    		{{ template "leftCol" "Remote Address" }}{{ template "rightCol" .Request.RemoteAddr }}
    	</div>
    	<div class="row">
    		{{ template "leftCol" "Request URI" }}{{ template "rightCol" .Request.RequestURI }}
    	</div>
    	<div class="row">
    		{{ template "leftCol" "Content Length" }}{{ template "rightCol" .Request.ContentLength }}
    	</div>

		<div class="row header">Request Headers</div>

    	{{ range $key, $value := .Request.Header }}
    	<div class="row">
			{{ template "leftCol" $key }}{{ template "rightCol" $value }}
		</div>
		{{ end }}

		<div class="row  header">Environment</div>
		{{ range $el := .Env }}
			<div class="row">{{ $el }}</div>
		{{ end }}

    </div>
</body>
</html>`

const plainIndexTemplate = `TEST JIG

Google App Engine
-----------------
AppID            {{ .AppId }}
Datacenter       {{ .Datacenter }}
InstanceID       {{ .InstanceId }}
VersionID        {{ .VersionId }}

Request
-------
Host:            {{ .Request.Host }}
Method:          {{ .Request.Method }}
URL:             {{ .Request.URL.Path }}
Proto:           {{ .Request.Proto }}
Remote Address:  {{ .Request.RemoteAddr }}
Request URI:     {{ .Request.RequestURI }}
Content Length:  {{ .Request.ContentLength }}
Headers:{{ range $key, $value := .Request.Header }}
       {{ $key }} :  {{ $value }}{{ end }}

Environment
-----------
{{ range $el := .Env }}{{ $el }}
{{ end }}
`

/*
 * My daily reading list
 */

type ListItem struct {
	Title string `json:"title"`
	Site  string `json:"site"`
}

var myList = []ListItem{
	ListItem{
		Title: "KC Star",
		Site:  "http://www.kansascity.com",
	},
	ListItem{
		Title: "Real Clear Politics",
		Site:  "http://www.realclearpolitics.com",
	},
	ListItem{
		Title: "Real Clear Defense",
		Site:  "http://www.realcleardefense.com",
	},
	ListItem{
		Title: "Real Clear Science",
		Site:  "http://www.realclearscience.com",
	},
	ListItem{
		Title: "Infoworld",
		Site:  "http://infoworld.com",
	},
	ListItem{
		Title: "Networkworld",
		Site:  "http://networkworld.com",
	},
	ListItem{
		Title: "Wall Street Journal",
		Site:  "https://www.wsj.com",
	},
	ListItem{
		Title: "Hacker News",
		Site:  "https://news.ycombinator.com",
	},
}

func init() {
	http.HandleFunc("/", handler)
}

func testjig(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	var err error

	offers := []string{"text/plain", "text/html"}
	contentType := NegotiateContentType(r, offers, "text/plain")

	if contentType == "text/html" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t, err = template.New("index").Parse(htmlIndexTemplate)
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		t, err = template.New("index").Parse(plainIndexTemplate)
	}

	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}

	ctx := gae.NewContext(r)
	data := struct {
		Title      string
		Headers    map[string][]string
		Env        []string
		Request    *http.Request
		AppId      string
		Datacenter string
		InstanceId string
		VersionId  string
	}{
		Title:      "hello, world",
		Headers:    r.Header,
		Env:        os.Environ(),
		Request:    r,
		AppId:      gae.AppID(ctx),
		Datacenter: gae.Datacenter(ctx),
		InstanceId: gae.InstanceID(),
		VersionId:  gae.VersionID(ctx),
	}

	err = t.Execute(w, data)
}

func list(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(myList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(b)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := strings.ToLower(r.URL.Path)
	if path == "/testjig" {
		testjig(w, r)
		return
	}

	if path == "/api/v1/list" {
		list(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

// from github.com/golang/gddo/httputil/Negotiate.go
//
// NegotiateContentType returns the best offered content type for the request's
// Accept header. If two offers match with equal weight, then the more specific
// offer is preferred.  For example, text/* trumps */*. If two offers match
// with equal weight and specificity, then the offer earlier in the list is
// preferred. If no offers match, then defaultOffer is returned.
func NegotiateContentType(r *http.Request, offers []string, defaultOffer string) string {
	bestOffer := defaultOffer
	bestQ := -1.0
	bestWild := 3
	specs := header.ParseAccept(r.Header, "Accept")
	for _, offer := range offers {
		for _, spec := range specs {
			switch {
			case spec.Q == 0.0:
				// ignore
			case spec.Q < bestQ:
				// better match found
			case spec.Value == "*/*":
				if spec.Q > bestQ || bestWild > 2 {
					bestQ = spec.Q
					bestWild = 2
					bestOffer = offer
				}
			case strings.HasSuffix(spec.Value, "/*"):
				if strings.HasPrefix(offer, spec.Value[:len(spec.Value)-1]) &&
					(spec.Q > bestQ || bestWild > 1) {
					bestQ = spec.Q
					bestWild = 1
					bestOffer = offer
				}
			default:
				if spec.Value == offer &&
					(spec.Q > bestQ || bestWild > 0) {
					bestQ = spec.Q
					bestWild = 0
					bestOffer = offer
				}
			}
		}
	}
	return bestOffer
}

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/peterhellberg/kop/list"
	"github.com/peterhellberg/kop/rpc"
	"github.com/peterhellberg/kop/store/file"
	"github.com/peterhellberg/kop/store/memory"
)

const defaultPort = "12432"

func main() {
	var name string

	flag.StringVar(&name, "name", "list", "The name of the cache file to store the list in")
	flag.Parse()

	server := rpc.NewServer()

	store, err := file.Store(name, memory.Store())
	if err != nil {
		log.Fatalf("file.Store error: %v\n", err)
	}

	svc := list.New(store)

	rpc.RegisterList(server, svc)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			server.NotFound.ServeHTTP(w, r)
			return
		}

		if res, err := svc.Items(r.Context(), rpc.ItemsRequest{}); err == nil {
			t.Execute(w, res)
		}
	})

	http.Handle(server.Basepath, server)

	log.Fatal(http.ListenAndServe(":"+port(), nil))
}

func port() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}

	return defaultPort
}

var t = template.Must(template.New("").Funcs(template.FuncMap{"join": strings.Join}).Parse(`<!doctype html>
<html lang="en" data-theme="light">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<meta name="apple-mobile-web-app-capable" content="yes">
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@next/css/pico.classless.min.css">
		<title>📝 {{ join .Items ", " }}</title>
	</head>
	<body>
		<main>
			<article>
				<header>
					<h1><a href="/" style="text-decoration: none; color: #2d3138">Köp 📝</a></h1>
				</header>
				<fieldset>
				{{ range .Items }}
					<h2>
						<label>
							<input type="checkbox" name="{{.}}" />
							{{.}}
						</label>
					</h2>
					<hr>
				{{ end }}
				</fieldset>
			</article>
		</main>
	</body>
</html>`))

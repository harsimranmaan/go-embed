package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

// The go embed directive statement must be outside of function body

// Embed the file content as string.
//go:embed title.txt
var title string

// Embed the entire directory.
//go:embed templates
var indexHTML embed.FS

//go:embed static
var staticFiles embed.FS

func main() {

	// Note the call to ParseFS instead of Parse
	t, err := template.ParseFS(indexHTML, "templates/index.html.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	// http.FS can be used to create a http Filesystem
	var staticFS = http.FS(staticFiles)
	fs := http.FileServer(staticFS)

	// Serve static files
	http.Handle("/static/", fs)
	// Handle all other requests
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var path = req.URL.Path
		log.Println("Serving request for path", path)
		w.Header().Add("Content-Type", "text/html")

		// respond with the output of template execution
		t.Execute(w, struct {
			Title    string
			Response string
		}{Title: title, Response: path})

	})

	log.Println("Listening on :3000...")
	// start the server
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

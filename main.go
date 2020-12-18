package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

func main() {

	// Embed the file content as string
	//go:embed title.txt
	var title string

	// Embed the entire directory. The path is relative to the package directory
	//go:embed templates
	var indexHTML embed.FS

	// Note the call to ParseFS instead of Parse
	t, err := template.ParseFS(indexHTML, "templates/index.html.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	//go:embed static
	var staticFiles embed.FS

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

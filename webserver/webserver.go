package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Post struct {
	ID       string
	PostedAt string
	Title    string
	Body     string
}

var postTemplate = template.Must(template.New("").Parse(`
<html>
  <body>
    <h1>{{.Title}}</h1>
	<h3>{{.PostedAt}}</h3>
	<p>{{.Body}}</p>
  </body>
</html>
	`))

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	http.HandleFunc("POST /posts", func(w http.ResponseWriter, r *http.Request) {
		contents, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id, err := createPost(r.FormValue("title"), string(contents))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "created post with id %q\n", id)
	})

	http.HandleFunc("GET /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, err := getPost(r.PathValue("id"))
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(w, "no such post", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var buf bytes.Buffer
		if err := postTemplate.Execute(&buf, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buf.WriteTo(w)
	})

	log.Printf("serving")
	log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil))
}

func createPost(title, body string) (id string, err error) {
	mu.Lock()
	defer mu.Unlock()
	id = strconv.Itoa(nextID)
	nextID++
	posts[id] = &Post{
		ID:       id,
		PostedAt: time.Now().Format(time.DateTime),
		Title:    title,
		Body:     body,
	}
	return id, nil
}

func getPost(id string) (*Post, error) {
	mu.Lock()
	defer mu.Unlock()
	if p, ok := posts[id]; ok {
		return p, nil
	}
	return nil, fs.ErrNotExist
}

var (
	mu     sync.Mutex
	nextID int
	posts  = map[string]*Post{}
)

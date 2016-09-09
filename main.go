/*
A simple web server designed to provide images. Written in Go!
*/

package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"net/http"
	"html/template"
	"regexp"
)

const usage string = "no arguments needed"
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))


type Page struct {
	Title string
	Body []byte
}

// function save() returns an error
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// function loadPage(string), returns (page pointer, error)
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}


// init function is automatically called when entering the file
// check for arguments
func init() {
	fmt.Printf("Init func\n")
	if len(os.Args) != 1 {
		fmt.Printf("%s\n%s\n",usage,os.Args[0])
		os.Exit(1)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there! <b>It's working!</b>, %s", r.URL.Path[1:])
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	// Slice containing the path of capacity length of the path
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)

	// On successful load
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Image Page, %s", r.URL.Path[1:])
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "upload page, %s", r.URL.Path[1:])
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func registerHandlers() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/img/", imageHandler)
	http.HandleFunc("/upload/", uploadHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
}

func main() {
	registerHandlers()
	http.ListenAndServe(":8080", nil)
	return
}


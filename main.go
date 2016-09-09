/*
A simple web server designed to provide images. Written in Go!
*/

package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"net/http"
)

const usage string = "no arguments needed"

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

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
	return
}


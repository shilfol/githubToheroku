package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"text/template"
)

type Page struct {
	Title string
	Body  []byte
}

const lenPath = len("/view/")

var templates = make(map[string]*template.Template)

var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

func init() {
	for _, tmpl := range []string{"edit", "view"} {
		t := template.Must(template.ParseFile(tmpl + ".html"))
		templates[tmpl] = t
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (title string, err os.Error) {
	title = r.URL.Path[lenPath:]
	if !titleValidator.MatchString(title) {
		http.NotFound(w, r)
		err = os.NewError("Invalid Page Title")
	}
	return
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)

}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err = templates[tmpl].Execute(w, p)
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title = r.URL.Path[lenPath:]
		if !titleValidator.MatchString(title) {
			http.NotFound(w, r)
			return
		}
		fn(w, r, title)
	}

}

func (p *Page) save() os.Error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 600)
}

func loadPage(title string) *Page {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err

	}
	return &Page{Title: title, Body: body}
}

func main() {

	http.HandleFunc("/view/", makeHandler(viewhandler))
	http.HandleFunc("/edit/", makeHandler(edithandler))
	http.HandleFunc("/save/", makeHandler(savehandler))
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", r.URL.Path[1:])

}

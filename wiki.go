package main

import (
		"html/template"
		"io/ioutil"
		"net/http"
		"regexp"
		"strings"
		"fmt"
		"log"

		"github.com/spf13/viper"
)

var envPrefix = "wiki"
var templatePath = "tmpl/*"
var dataPath = "data/"
var templates map[string]*template.Template
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
		Title string
		Body []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(dataPath+filename, p.Body, 0600)
}

func setupConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("app")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error reading config file: %s \n", err))
	}
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func loadTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates["view"] = template.Must(template.ParseFiles("tmpl/layout.html", "tmpl/sidebar.html", "tmpl/view.html"))
	templates["edit"] = template.Must(template.ParseFiles("tmpl/layout.html", "tmpl/sidebar.html", "tmpl/edit.html"))
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(dataPath+filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, ok := templates[tmpl]
	if !ok {
		http.Error(w, "Templte does not exist: "+tmpl, http.StatusNotFound)
	}
	err := t.ExecuteTemplate(w, "layout", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/home", http.StatusFound)
}

func init() {
	setupConfig()
	loadTemplates()
}

func main() {
	// Handle our view / edit / save routes
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	// Serve static assets from teh static folder
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Redirect requests for the root "/" to /view/home
	http.HandleFunc("/", rootHandler)

	log.Fatal(http.ListenAndServe(":"+viper.GetString("server.port"), nil))
}


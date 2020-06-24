package main

import (
	"bytes"
	"github.com/devprojx/gamblr/lib"
	"github.com/zserge/lorca"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

//Page defines the struct of a webpage
type Page struct {
	Title string
	Js    []string
	CSS   []string
}

var tmpl *template.Template
var workDir string
var templatesPath string
var page Page

func init() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Printf("[error] unable to get working directory: %s", err)
	}

	templatesPath = "./views/"

	//Set Javascript and CSS paths
	page = Page{
		Js: []string{
			"/public/js/http.js",
			"/public/js/app.js",
		},
		CSS: []string{
			"/public/css/main.css",
			"https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.11.2/css/all.min.css",
			"https://cdn.jsdelivr.net/npm/bulma@0.9.0/css/bulma.min.css",
		},
	}

	log.Print("[info] templates loading...")
	tmpl, err = template.ParseGlob(filepath.Join(workDir, templatesPath+"*"))
	if err != nil {
		log.Printf("[error] unable to load tempalates: %s", err)
	}

	_ = tmpl
}

func initBindings(ui lorca.UI) {
	ui.Bind("saveSettings", lib.SaveSettings)
	ui.Bind("loadSettings", lib.LoadSettings)
}

func main() {

	ui, err := lorca.New("", "", 640, 480)
	if execErr, ok := err.(*exec.Error); ok {
		lorca.PromptDownload()
		log.Fatalf("Chrome could not be started. Do you have Chrome installed? %s\n", execErr)
	} else if err != nil {
		log.Fatalf("Ops. Something went wrong while starting the application: %s", err)
	}
	defer ui.Close()

	initBindings(ui)

	page.Title = "Gamblr"

	var tpl bytes.Buffer
	err = tmpl.ExecuteTemplate(&tpl, "games", &page)
	if err != nil {
		log.Println("[error]: loading template ", err)
	}

	// Load HTML after Go functions are bound to JS
	ui.Load("data:text/html," + url.PathEscape(tpl.String()))

	go func() {
		http.HandleFunc("/", lib.SocketHandler)

		if err := http.ListenAndServe(":1234", nil); err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}()

	<-ui.Done()
}

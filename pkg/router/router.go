package router

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Page struct {
	relPath string
	route   string
}

func Serve() {
	pages := preFlight()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, page := range pages {
			if r.URL.Path == page.route {
				http.ServeFile(w, r, filepath.Join("./templates/", page.relPath))
				return
			}
		}
		if strings.HasPrefix(r.URL.Path, "/static/") {
			//serve the static files
			http.ServeFile(w, r, filepath.Join(".", r.URL.Path))
			return
		}
		http.ServeFile(w, r, "./templates/notfound.html")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func preFlight() []Page {
	//✈✈✈ Pre-flight checks ✈✈✈
	var pages []Page = getTemplatesToServe()
	if len(pages) == 0 {
		log.Fatal("No pages found in templates folder")
	}
	validateTemplateCorrectFormat(pages)
	return pages
}

func getTemplatesToServe() []Page {
	//get root directory of project
	f, _ := os.Getwd()
	files, err := ioutil.ReadDir(f + "/templates")
	if err != nil {
		log.Fatal(err)
	}

	var pages []Page
	for _, file := range files {
		if file.IsDir() {
			subdirFiles, err := ioutil.ReadDir(filepath.Join("./templates", file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			for _, subdirFile := range subdirFiles {
				pages = append(
					pages,
					Page{
						filepath.Join(file.Name(), subdirFile.Name()), "/" + strings.Replace(filepath.ToSlash(filepath.Join(file.Name(), subdirFile.Name())), ".html", "", -1)})
			}
		} else {
			pages = append(pages, Page{file.Name(), "/" + strings.Replace(file.Name(), ".html", "", -1)})
		}
	}
	return pages
}

func validateTemplateCorrectFormat(pages []Page) {
	for _, page := range pages {
		if filepath.Ext(page.relPath) != ".html" {
			log.Fatal("Page does not end in .html: " + page.relPath)
		}
	}
}

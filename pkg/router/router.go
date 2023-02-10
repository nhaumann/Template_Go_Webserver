package router

import (
	"log"
	"net/http"
	"os"
	fp "path/filepath"
	str "strings"
)

type Page struct {
	relPath string
	route   string
}

const TEMPLATES_PATH = "./templates"

const NOT_FOUND_PATH = TEMPLATES_PATH + "/notfound.html"

func ServeTemplatesAndStyles() {
	pages := preFlight()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, page := range pages {
			if r.URL.Path == page.route {
				http.ServeFile(w, r, fp.Join(TEMPLATES_PATH, page.relPath))
				return
			}
		}
		if str.HasPrefix(r.URL.Path, "/static/") {
			http.ServeFile(w, r, fp.Join(".", r.URL.Path))
			return
		}
		http.ServeFile(w, r, NOT_FOUND_PATH)
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
	f, _ := os.Getwd()
	files, err := os.ReadDir(f + TEMPLATES_PATH)

	if err != nil {
		log.Fatal(err)
	}

	var pages []Page
	for _, file := range files {
		if file.IsDir() {
			subdirFiles, err := os.ReadDir(fp.Join(TEMPLATES_PATH, file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			for _, subdirFile := range subdirFiles {
				pages = append(
					pages,
					Page{
						fp.Join(file.Name(), subdirFile.Name()),
						"/" + str.Replace(fp.ToSlash(fp.Join(file.Name(), subdirFile.Name())), ".html", "", -1),
					},
				)
			}
		} else {
			pages = append(pages, Page{file.Name(), "/" + str.Replace(file.Name(), ".html", "", -1)})
		}
	}
	return pages
}

func validateTemplateCorrectFormat(pages []Page) {
	for _, page := range pages {
		if fp.Ext(page.relPath) != ".html" {
			log.Fatal("Page does not end in .html: " + page.relPath)
		}
	}
}

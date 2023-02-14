package webserver

import (
	"log"
	"net/http"
	"os"
	fp "path/filepath"
	str "strings"
)

type WebServerConfig struct {
	WebPort             string
	StaticContentPrefix string
	WebPath             string
	TemplatesPath       string
	NotFoundPath        string
	WebRoot             string
}

type Resource struct {
	relPath string
	route   string
}

// refactor to use config
func Serve(config WebServerConfig) {
	resources := runPreflightAndGetRoutes(config)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "" {
			http.ServeFile(w, r, fp.Join(config.WebPath, config.TemplatesPath, "index.html"))
			return
		}
		for _, page := range resources {
			if r.URL.Path == page.route {
				http.ServeFile(w, r, fp.Join(config.WebPath, config.TemplatesPath, page.relPath))
				return
			}
		}
		if str.HasPrefix(r.URL.Path, config.StaticContentPrefix) {
			http.ServeFile(w, r, fp.Join(config.WebPath, r.URL.Path))
			return
		}
		http.ServeFile(w, r, fp.Join(config.WebPath, config.NotFoundPath))
	})
	log.Fatal(http.ListenAndServe(config.WebPort, nil))
}

func runPreflightAndGetRoutes(config WebServerConfig) []Resource {
	//✈✈✈ Pre-flight checks ✈✈✈
	var resources []Resource = getTemplatesToServe(config)
	if len(resources) == 0 {
		log.Fatal("No pages found in templates folder")
	}
	validateTemplateCorrectFormat(resources)
	return resources
}

func getTemplatesToServe(config WebServerConfig) []Resource {
	f, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir(f + config.WebPath + config.TemplatesPath)

	if err != nil {
		log.Fatal(err)
	}

	var pages []Resource
	for _, file := range files {
		if file.IsDir() {
			subdirFiles, err := os.ReadDir(fp.Join(config.WebPath, config.TemplatesPath, file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			//code smell? maybe - readability improvement conflicts with Go styling
			for _, subdirFile := range subdirFiles {
				pages = append(
					pages,
					Resource{
						fp.Join(file.Name(), subdirFile.Name()),
						"/" + str.Replace(fp.ToSlash(fp.Join(file.Name(), subdirFile.Name())), ".html", "", -1),
					},
				)
			}
		} else {
			pages = append(pages, Resource{file.Name(), "/" + str.Replace(file.Name(), ".html", "", -1)})
		}
	}
	return pages
}

func validateTemplateCorrectFormat(pages []Resource) {
	for _, page := range pages {
		if fp.Ext(page.relPath) != ".html" {
			log.Fatal("file does not contain correct HTML extension: " + page.relPath)
		}
	}
}

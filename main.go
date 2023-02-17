package main

import (
	"bufio"
	"flag"
	"fmt"
	websrvr "glossary/pkg/webserver"
	"log"
	"os"

	types "glossary/data/types"

	psql "glossary/data/pssql"

	"github.com/joho/godotenv"
)

func main() {

	//This was me just playing around with concurrency and a random bits of code I wanted to try out.
	//cce.Setup()

	//get config from env

	// Define flags for the "add" command
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	term := addCmd.String("term", "", "term to add")
	definition := addCmd.String("definition", "", "definition of the term")

	// Parse command-line arguments
	if len(os.Args) < 2 {
		log.Fatal("command required: serve, list, or add")
	}

	//open db
	db, err := psql.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch os.Args[1] {
	case "serve":

		//get config from env or use docker args

		config := getConfigFromEnv()

		fmt.Println("Serving Glossary on port: " + config.WEB_PORT)
		websrvr.Serve(websrvr.WebServerConfig{
			WebPort:             config.WEB_PORT,
			StaticContentPrefix: config.STATIC_CONTENT_PREFIX,
			WebPath:             config.WEB_PATH,
			TemplatesPath:       config.TEMPLATES_PATH,
			NotFoundPath:        config.NOT_FOUND_FILE_NAME_PATH,
			WebRoot:             config.WEB_ROOT,
		})

		// go hcsrver.Serve(hcsrver.HealthCheckServerConfig{
		// 	WebPort: config.HEALTCHECK_SERVER_PORT,
		// })

		// go hcclnt.Connect(hcsrver.HealthCheckServerConfig{
		// 	WebPort: config.HEALTCHECK_SERVER_PORT,
		// })

		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')

	case "list":
		items, err := psql.GetGlossaryItems(db)
		if err != nil {
			log.Fatal(err)
		}
		for _, item := range items {
			fmt.Println("Term:", item.Term, " Definition:", item.Definition)
		}

	case "add":
		addCmd.Parse(os.Args[2:])
		if *term == "" || *definition == "" {
			log.Fatal("both --term and --definition flags are required")
		} else {

			id, err := psql.InsertGlossaryItem(db, types.GlossaryItem{
				Term:       *term,
				Definition: *definition,
			})

			if err != nil || id == 0 {
				log.Fatal(err)
			}

			fmt.Println("Added GlossaryItem with Term: " + *term + " and Definition: " + *definition)

		}

	case "reset":
		err := psql.ResetDB(db)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("DB Reset")

	default:
		log.Fatal("command required: serve, list, or add")
	}

	//wait for user input to exit

}

func getConfigFromEnv() ApplicationConfig {
	//load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	//get config from env
	WEB_PORT := os.Getenv("WEB_PORT")
	HEALTCHECK_SERVER_PORT := os.Getenv("HEALTCHECK_SERVER_PORT")
	STATIC_CONTENT_PREFIX := os.Getenv("STATIC_CONTENT_PREFIX")
	WEB_PATH := os.Getenv("WEB_PATH")
	TEMPLATES_PATH := os.Getenv("TEMPLATES_PATH")
	NOT_FOUND_FILE_NAME_PATH := os.Getenv("NOT_FOUND_FILE_NAME_PATH")
	WEB_ROOT := os.Getenv("WEB_ROOT")

	//set config
	return ApplicationConfig{
		WEB_PORT,
		HEALTCHECK_SERVER_PORT,
		STATIC_CONTENT_PREFIX,
		WEB_PATH,
		TEMPLATES_PATH,
		NOT_FOUND_FILE_NAME_PATH,
		WEB_ROOT,
	}
}

type ApplicationConfig struct {
	WEB_PORT                 string
	HEALTCHECK_SERVER_PORT   string
	STATIC_CONTENT_PREFIX    string
	WEB_PATH                 string
	TEMPLATES_PATH           string
	NOT_FOUND_FILE_NAME_PATH string
	WEB_ROOT                 string
}

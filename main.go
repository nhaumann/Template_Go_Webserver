package main

import (
	"fmt"
	"log"
	hcclnt "webserver/pkg/healthcheck/healthcheck_client"
	hcsrver "webserver/pkg/healthcheck/healthcheck_server"
	websrvr "webserver/pkg/webserver"

	"bufio"
	"os"

	env "github.com/joho/godotenv"
)

func main() {

	//This was me just playing around with concurrency and a random bits of code I wanted to try out.
	//cce.Setup()

	//get config from env
	config := getConfigFromEnv()

	fmt.Println(config)

	go websrvr.Serve(websrvr.WebServerConfig{
		WebPort:             config.WEB_PORT,
		StaticContentPrefix: config.STATIC_CONTENT_PREFIX,
		WebPath:             config.WEB_PATH,
		TemplatesPath:       config.TEMPLATES_PATH,
		NotFoundPath:        config.NOT_FOUND_FILE_NAME_PATH,
		WebRoot:             config.WEB_ROOT,
	})

	go hcsrver.Serve(hcsrver.HealthCheckServerConfig{
		WebPort: config.HEALTCHECK_SERVER_PORT,
	})

	go hcclnt.Connect(hcsrver.HealthCheckServerConfig{
		WebPort: config.HEALTCHECK_SERVER_PORT,
	})

	//wait for user input to exit
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

}

func getConfigFromEnv() ApplicationConfig {
	//load env
	err := env.Load(".env")
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

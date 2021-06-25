package main

import (
	"log"
	"os"

	"github.com/emacampolo/link-tracker/cmd/shell/server/handler"
	"github.com/emacampolo/link-tracker/internal/link"
	"github.com/emacampolo/link-tracker/internal/platform/shell"
)

/*
	Example:

	1. CREATE link:example.com password:secret
	2. REDIRECT id:1 password:secret
	3. METRICS id:1
	4. INACTIVATE id: 1
*/

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	linkRepository := link.NewInMemoryRepository()
	linkService := link.NewService(linkRepository)
	linkHandler := handler.NewLink(linkService)

	application := shell.New(os.Stdin, os.Stdout)

	application.AddHandlerFunc("CREATE", linkHandler.Create())
	application.AddHandlerFunc("REDIRECT", linkHandler.Redirect())
	application.AddHandlerFunc("METRICS", linkHandler.Metrics())
	application.AddHandlerFunc("INACTIVATE", linkHandler.Inactivate())

	return application.Run()
}

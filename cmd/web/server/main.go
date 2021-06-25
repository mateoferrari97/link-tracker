package main

import (
	"log"

	"github.com/mateoferrari97/link-tracker/cmd/web/server/handler"
	"github.com/mateoferrari97/link-tracker/internal/link"
	"github.com/mateoferrari97/link-tracker/internal/platform/web"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	linkRepository := link.NewInMemoryRepository()
	linkService := link.NewService(linkRepository)
	linkHandler := handler.NewLink(linkService)

	application := web.New()

	application.Method("POST", "/link", linkHandler.Create())
	application.Method("GET", "/link/{id}", linkHandler.Redirect())
	application.Method("GET", "/link/{id}/metrics", linkHandler.Metrics())
	application.Method("POST", "/link/{id}/inactivate", linkHandler.Inactivate())

	return application.Run()
}

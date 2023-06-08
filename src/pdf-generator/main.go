package main

import (
	"github.com/Limpid-LLC/pdf-generator-service/internal"
	"github.com/Limpid-LLC/saiService"
)

func main() {
	svc := saiService.NewService("htmlToPDF converter")
	is := internal.InternalService{Context: svc.Context}

	svc.RegisterConfig("config.yml")

	svc.RegisterInitTask(is.Init)

	svc.RegisterHandlers(
		is.NewHandler(),
	)

	svc.Start()
}

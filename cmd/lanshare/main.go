package main

import (
	"fmt"
	"os"

	"lanshare/config"
	"lanshare/internal/app"
	"lanshare/internal/db"
	"lanshare/internal/mdns"
)

////go:embed web/dist/*
//var frontend embed.FS

func main() {

	database := db.NewService()
	err := database.Start()
	if err != nil {
		panic(err)
	}
	fmt.Println("Datubāze veiksmīgi startēta")

	currentSettings, err := app.SettingsStartup(database)
	if err != nil {
		fmt.Printf("Startēšana pārtraukta: %v\n", err)
		return
	}

	application := app.NewApp(database, currentSettings)
	_ = application

	if currentSettings != nil && currentSettings.StoragePath != "" {
		if err := os.MkdirAll(currentSettings.StoragePath, 0755); err != nil {
			fmt.Printf("Brīdinājums: neizdevās izveidot glabāšanas mapi '%s': %v\n", currentSettings.StoragePath, err)
		}
	}

	go func() {
		application.Router.Serve()
	}()

	mdnsServer := mdns.NewMDNsServer(config.DomainName)
	go func() {
		mdnsServer.Serve(config.DomainName)
	}()
	select {}
}

package main

import "github.com/gonfff/mockster/app"

func main() {
	// 	// run main application
	application := app.NewApp()
	application.Setup()
	application.Start()

}

package main

import "github.com/gonfff/mockster/app/cmd"

func main() {
	// 	// run main application
	app := cmd.NewApp()
	app.Setup()
	app.Start()

}

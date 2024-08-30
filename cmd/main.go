package main

import "github.com/ZnNr/notes-keeper.git/intenal/app"

const configPath = "./config/config.yaml"

func main() {
	app.Run(configPath)
}

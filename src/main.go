package main

import (
	"github.com/daniel5268/meliChallenge/src/app"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	app.NewApp().StartApp()
}

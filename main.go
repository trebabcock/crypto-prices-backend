package main

import (
	app "crypto-prices/app"
)

func main() {
	app := &app.App{}
	app.Init()
	app.Run(":2814")
}

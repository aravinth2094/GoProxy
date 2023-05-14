package main

import (
	"os"

	"github.com/aravinth2094/GoProxy/app"
)

func main() {
	app.CreateApp().Run(os.Args)
}

package main

import (
	"github.com/Slava02/URL_shortner/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

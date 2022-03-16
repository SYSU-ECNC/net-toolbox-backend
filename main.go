package main

import (
	"toolBox/router"
)

func main() {
	r := router.SetupRouters()

	r.Run(":8888")
}

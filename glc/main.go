package main

import (
	"glc/gweb/http"
	"glc/web/router"
)

func main() {
	http.StartHttpServer(router.Register)
}

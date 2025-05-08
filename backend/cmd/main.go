package main

import (
	"github.com/othello-api/internal/routers"
)

func main() {
	r := router.SetupRouter()
	r.Run(":8080")
}

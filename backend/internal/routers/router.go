package router

import (
	"github.com/gin-gonic/gin"
	"github.com/othello-api/internal/handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/api/games/:gameID", handler.GetGameState)
	r.POST("/api/games", handler.CreateGame)
	r.POST("/api/games/:gameID/moves", handler.MakeMove)
	r.GET("/api/games/IDs", handler.GetGameIDs)
	return r
}

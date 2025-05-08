package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	othello "github.com/othello-api/internal/utils"
	"net/http"
	"sync"
)

// ゲームの状態をメモリに保持
// TODO: DB で動作するようにする
var (
	games = make(map[string]*othello.Othello)
	mutex = &sync.Mutex{}
)

type MoveRequest struct {
	Row int `json:"row" binding:"required,gte=0,lt=8"`
	Col int `json:"col" binding:"required,gte=0,lt=8"`
}

// 新しいゲームを作成
func CreateGame(c *gin.Context) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	uuidStr := uuid.String()
	o := othello.NewGame(uuidStr)
	mutex.Lock()
	games[o.GameID] = o
	mutex.Unlock()
	responseData := gin.H{
		"message":      "New game created successfly!",
		"gameId":       o.GameID,
		"initialBoard": o.Board,
	}
	c.JSON(http.StatusOK, responseData)
}

// ゲームの状態を取得
func GetGameState(c *gin.Context) {
	gameID := c.Param("gameID")
	mutex.Lock()
	gameStatus, exists := games[gameID]
	mutex.Unlock()
	if !exists {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error":       "Game not found",
				"requestedId": gameID,
			},
		)
		panic("Cannot find Game")
	}
	responseData := gin.H{
		"gameID":        gameStatus.GameID,
		"initialBoard":  gameStatus.Board,
		"currentPlayer": gameStatus.Player,
		"validMove":     gameStatus.ValidMoves,
		"winner":        gameStatus.Winner,
	}
	c.JSON(http.StatusOK, responseData)
}

// 手番を移動
func MakeMove(c *gin.Context) {
	gameID := c.Param("gameID")
	var req MoveRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
		return
	}
	mutex.Lock()
	gameStatus, exists := games[gameID]
	mutex.Unlock()
	if !exists {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error":       "Game not found",
				"requestedId": gameID,
			},
		)
		panic("Cannot find Game")
	}
	if err := gameStatus.MakeMove(req.Row, req.Col); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mutex.Lock()
	games[gameID] = gameStatus
	mutex.Unlock()
	responseData := gin.H{
		"gameID":        gameStatus.GameID,
		"initialBoard":  gameStatus.Board,
		"currentPlayer": gameStatus.Player,
		"validMove":     gameStatus.ValidMoves,
		"winner":        gameStatus.Winner,
	}
	c.JSON(http.StatusOK, responseData)
}

// 管理している Game ID を全て表示
func GetGameIDs(c *gin.Context) {
	mutex.Lock()
	gameIDs := make([]string, 0, len(games))
	for key := range games {
		gameIDs = append(gameIDs, key)
	}
	mutex.Unlock()
	responseData := gin.H{
		"gameIDs": gameIDs,
	}
	c.JSON(http.StatusOK, responseData)
}

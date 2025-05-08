package othello

import (
	"fmt"
)

const (
	BoardSize int = 8 // 盤面サイズ
	Empty     int = 0 // 空
	Black     int = 1 // 黒石
	White     int = 2 // 白石
)

type Othello struct {
	Board      [BoardSize][BoardSize]int // オセロ盤
	GameID     string                    // ゲームID
	Player     int                       // プレイヤー
	ValidMoves [][2]int                  // 有効な手のリスト
	Winner     int                       // 0: 進行中, 1: 黒の勝利, 2: 白の勝利, 3: 引き分け
}

// 新しいゲームの作成
func NewGame(uuid string) *Othello {
	var o Othello
	mid := BoardSize / 2
	o.Board[mid-1][mid-1] = White
	o.Board[mid-1][mid] = Black
	o.Board[mid][mid-1] = Black
	o.Board[mid][mid] = White
	o.GameID = uuid
	o.Player = Black
	return &o
}

// 盤内にあるか判定
func IsValidPosition(r, c int) bool {
	if 0 <= r && r < BoardSize && 0 <= c && c < BoardSize {
		return true
	}
	return false
}

// 相手の手番を取得
func (o *Othello) getOpponent() (int, error) {
	if o.Player == Black {
		return White, nil
	} else if o.Player == White {
		return Black, nil
	}
	return -1, fmt.Errorf("invalid player value: %d", o.Player)
}

// 手が有効か確認
func (o *Othello) checkMove(r, c int) ([][2]int, bool) {
	// 盤外チェック
	if !IsValidPosition(r, c) {
		return nil, false
	}
	// 空か判定
	if o.Board[r][c] != Empty {
		return nil, false
	}
	// 8 方向のチェック
	var (
		totalFlips [][2]int
		tmpFlips   [BoardSize - 1][2]int
	)
	directions := [8][2]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}
	opponent, err := o.getOpponent()
	if err != nil {
		panic(err)
	}
	var count int
	for _, v := range directions {
		count = 0
		for i := 1; i < BoardSize; i++ {
			nextR, nextL := r+v[0]*i, c+v[1]*i
			// 盤外チェック
			if !IsValidPosition(nextR, nextL) {
				break
			}
			if o.Board[nextR][nextL] == Empty {
				break
			} else if o.Board[nextR][nextL] == o.Player {
				if count > 0 {
					totalFlips = append(totalFlips, tmpFlips[:count]...)
				}
				break
			} else if o.Board[nextR][nextL] == opponent {
				tmpFlips[count][0], tmpFlips[count][1] = nextR, nextL
				count++
			} else {
				panic("Board's num looks strange")
			}
		}
	}
	if len(totalFlips) == 0 {
		return nil, false
	}
	return totalFlips, true
}

// 石のおける箇所を取得
func (o *Othello) updateValidMoves() bool {
	var validMoves [][2]int
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			_, isMove := o.checkMove(i, j)
			if isMove {
				validMoves = append(validMoves, [2]int{i, j})
			}
		}
	}
	o.ValidMoves = validMoves
	return len(o.ValidMoves) != 0
}

// 勝者の判定
func (o *Othello) determineWinner() {
	var whiteNum, blackNum int
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			if o.Board[i][j] == White {
				whiteNum++
			} else if o.Board[i][j] == Black {
				blackNum++
			}
		}
	}
	if whiteNum > blackNum {
		o.Winner = White
	} else if whiteNum < blackNum {
		o.Winner = Black
	} else {
		o.Winner = White + Black
	}
}

// 手番の移動
func (o *Othello) MakeMove(r, c int) error {
	if o.Winner != 0 {
		return fmt.Errorf("game has been already over")
	}
	flips, isMove := o.checkMove(r, c)
	if !isMove {
		return fmt.Errorf("invalid move at (%d, %d)", r, c)
	}
	o.Board[r][c] = o.Player
	for _, v := range flips {
		o.Board[v[0]][v[1]] = o.Player
	}
	opponent, err := o.getOpponent()
	if err != nil {
		panic(err)
	}
	// 手番の交代
	o.Player = opponent
	if !o.updateValidMoves() {
		currentPlayer, _ := o.getOpponent()
		o.Player = currentPlayer
		if !o.updateValidMoves() {
			// 両者とも手がない場合
			o.determineWinner()
		}
	}
	return nil
}

# backend 環境構築
### API サーバーの立ち上げ方
```bash
# docker の build
$ docker compose build
# docker コンテナの up
$ docker compose up -d
# docker コンテナに入る
$ docker compose exec web sh
# ソースコードで必要としている module のインストール
$ go mod tidy
# othello-api を立ち上げる
$ go run cmd/main.go
```
### API の利用方法
1. Postman を利用する
2. curl を利用する
     ```bash
     # ゲームを新しく作成
     $ curl -X POST http://localhost:8080/api/games
     # ゲームの情報を取得
     $ curl -X GET http://localhost:8080/api/games/{ここに uuid を入れる}
     # 作成したゲーム全ての情報を取得
     $ curl -X GET http://localhost:8080/api/games/IDs
     # Row(行)，Col(列)に石を置く
     $ curl -X POST http://localhost:8080/api/games/{ここに uuid を入れる}/moves?Row=4&Col=5
     ```
### テスト方法
`(Row, Col)`に次の座標を入れて9手で黒が勝利するかどうか調べる。
| 番手 | `(Row, Col)` |
| ---- | ------------ |
| 1    | `(4,5)`      |
| 2    | `(5,3)`      |
| 3    | `(4,2)`      |
| 4    | `(3,5)`      |
| 5    | `(6,4)`      |
| 6    | `(5,5)`      |
| 7    | `(4,6)`      |
| 8    | `(5,4)`      |
| 9    | `(2,4)`      |

package controller

// func (h *Handler) TestWebsocket(ctx *gin.Context) {
// 	var (
// 		wsConn, conn, _ = NewWebsocketConn(ctx)
// 		// result          UserGameParam
// 	)

// 	defer wsConn.Close()
// 	defer conn.Close()

// 	go func() {
// 		for {
// 			// log.Println("測試用~")

// 			// fmt.Println("收到玩家端遊戲狀態訊息")
// 			b, _ := json.Marshal(GameParam{
// 				Game: GameModel{
// 					ControlGameStatus: "123",
// 					GameStatus:        "123",
// 				}})
// 			conn.WriteMessage(b)

// 			// ws關閉
// 			if conn.isClose {
// 				return
// 			}

// 			time.Sleep(time.Second * 1)
// 		}
// 	}()

// 	for {
// 		var (
// 			result UserGameParam
// 		)

// 		data, err := conn.ReadMessage()
// 		if err != nil {
// 			// fmt.Println("關閉玩家端遊戲狀態ws")

// 			conn.Close()
// 			return
// 		}

// 		log.Println("收到測試用訊息")

// 		json.Unmarshal(data, &result)

// 		log.Println("訊息: ", result.Game.GameRound, result.Game.GameStatus)
// 	}
// }

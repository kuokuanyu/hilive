package controller

import (
	"encoding/json"
	"hilive/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// GameParam 場次相關參數資訊
type GameParam struct {
	GameIDs    string                 // 遊戲ID(多個，用逗號區隔ID)
	Games      []GameModel            // 場次資訊(多個)
	Game       GameModel              // 場次資訊
	PrizeStaff models.PrizeStaffModel // 中獎資訊
	Reset      bool                   `json:"reset" example:"true"`          // 頁面是否需要重整
	Error      string                 `json:"error" example:"error message"` // 錯誤訊息
}

// GameModel 資料表欄位
type GameModel struct {
	GameID            string  `json:"game_id" example:"game_id"`                                                    // 遊戲ID
	GameID2           string  `json:"game_id_2" example:"game_id_2"`                                                // 遊戲ID2
	GameRound         int64   `json:"game_round" example:"1"`                                                       // 遊戲輪次
	GameAttend        int64   `json:"game_attend" example:"100"`                                                    // 遊戲參加人數
	GameScore         int64   `json:"game_score" example:"100"`                                                     // 分數
	GameScore2        float64 `json:"game_score_2" example:"1.1"`                                                   // 第二分數
	GameCorrect       int64   `json:"game_correct" example:"5"`                                                     // 玩家答對題數
	GameRank          int64   `json:"game_rank" example:"100"`                                                      // 排名
	GameStatus        string  `json:"game_status" example:"open、start、end、close、result、question、calculate"`         // 遊戲狀態
	ControlGameStatus string  `json:"control_game_status" example:"open、start、end、close、result、question、calculate"` // 遙控端遊戲狀態
	GameSecond        int64   `json:"game_second" example:"30"`                                                     // 秒數
	PrizeAmount       int64   `json:"prize_amount" example:"100"`                                                   // 獎品數量
	WinPrizeAmount       int64   `json:"win_prize_amount" example:"100"`                                                   // 獎品數量(勝隊)
	LosePrizeAmount       int64   `json:"lose_prize_amount" example:"100"`                                                   // 獎品數量(敗隊)
	QARound           int64   `json:"qa_round" example:"1"`                                                         // 題目進行題數
	Error             string  `json:"error" example:"error message"`                                                // 錯誤訊息

	QAOption    string `json:"qa_option" example:"A"`      // 遊戲選項
	OriginScore int64  `json:"origin_score" example:"100"` // 原始分數
	AddScore    int64  `json:"add_score" example:"100"`    // 加成分數

	// 賓果遊戲
	BingoAttend    int64   `json:"bingo_attend" example:"100"`   // 完成選號人數
	BingoRound     int64   `json:"bingo_round" example:"1"`      // 賓果遊戲回合數
	BingoPeople    int64   `json:"bingo_people" example:"1"`     // 即將賓果人數
	LastPeople     int64   `json:"last_people" example:"1"`      // 即將賓果人數
	Number         int64   `json:"number" example:"1"`           // 賓果遊戲號碼
	Numbers        []int64 `json:"numbers" example:"1"`          // 賓果遊戲抽球號碼
	GuestNumbers        []int64 `json:"guest_numbers" example:"1"`          // 玩家賓果號碼排序
	UserBingoRound int64   `json:"user_bingo_round" example:"1"` // 完成賓果的回合數
	GoingBingo     bool    `json:"going_bingo" example:"true"`   // 是否即將賓果
	IsBingo        bool    `json:"is_bingo" example:"true"`      // 是否賓果

	// 隊伍資訊
	LeftTeamPlayers     []UserModel // 左方隊伍玩家
	RightTeamPlayers    []UserModel // 右方隊伍玩家
	LeftTeamLeader      UserModel   // 左方隊伍隊長
	RightTeamLeader     UserModel   // 右方隊伍隊長
	LeftTeamMVP         UserModel   // 左方隊伍mvp
	RightTeamMVP        UserModel   // 右方隊伍mvp
	LeftTeamScore       int64       `json:"left_team_score" example:"100"`      // 左方隊伍分數
	RightTeamScore      int64       `json:"right_team_score" example:"100"`     // 右方隊伍分數
	LeftTeamGameAttend  int64       `json:"left_team_game_attend" example:"0"`  // 左方隊伍參加遊戲人數
	RightTeamGameAttend int64       `json:"right_team_game_attend" example:"0"` // 右方隊伍參加遊戲人數
	WinTeam             string      `json:"win_team" example:"left_team"`       // 獲勝隊伍
	IsZero              bool        `json:"is_zero" example:"true"`             // 分數是否為0

	// 扭蛋機遊戲
	IsShake bool `json:"is_shake" example:"true"` // 是否搖晃

	// 投票遊戲
	VoteOptions      []GameVoteOptionModel // 投票選項
	VoteRecords      []GameVoteRecordModel // 投票紀錄
	VoteAvatars []GameVoteAvatarModel // 投票頭像
	UserID           string                `json:"user_id" example:"1"`
	OptionID         string                `json:"option_id" example:"option_id"`    // 選項ID
	Score            int64                 `json:"score" example:"100"`              // 分數
	VoteTimes        int64                 `json:"vote_times" example:"100"`         // 剩餘投票次數
	VoteScore        int64                 `json:"vote_score" example:"100"`         // 投票權重
	VoteMethodPlayer int64                 `json:"vote_method_player" example:"100"` // 玩家投票方式
	// VoteWinningStaffs  []UserModel // 投票遊戲中獎人員

	// qrcode
	QRcode string `json:"qrcode" example:"qrcode"` // qrcode url
}

// @Summary 即時回傳遊戲狀態資訊(玩家端)
// @Tags Websocket
// @version 1.0
// @Accept  json
// @Param activity_id query string true "activity ID"
// @Param game query string true "game" Enums(redpack, ropepack, whack_mole, monopoly, draw_numbers, QA, lottery)
// @param body body GameParam true "game param"
// @Success 200 {array} GameParam
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /ws/v1/game [get]
func (h *Handler) GameWebsocket(ctx *gin.Context) {
	var (
		activityID        = ctx.Query("activity_id")
		game              = ctx.Query("game")
		wsConn, conn, err = NewWebsocketConn(ctx)
		// result            GameParam
	)
	// fmt.Println("開啟回傳即時遊戲狀態資訊(玩家端)ws")
	defer wsConn.Close()
	defer conn.Close()
	// 判斷活動、遊戲資訊是否有效
	if err != nil || activityID == "" || game == "" {
		b, _ := json.Marshal(GameParam{
			Error: "錯誤: 無法辨識活動或遊戲資訊、錯誤的網域請求"})
		conn.WriteMessage(b)
		return
	}

	for {
		var (
			result GameParam
		)

		data, err := conn.ReadMessage()
		if err != nil {
			// fmt.Println("關閉回傳即時遊戲狀態資訊(玩家端)ws")
			conn.Close()
			return
		}

		// 解碼
		json.Unmarshal(data, &result)
		var games = make([]GameModel, 0)              // 紀錄所有場次遊戲狀態資訊
		gameIDs := strings.Split(result.GameIDs, ",") // 多個場次遊戲ID

		if len(gameIDs) > 0 {
			// 取得遊戲狀態資訊
			for _, id := range gameIDs {
				gameModel := h.getGameInfo(id, game) // 遊戲資訊

				games = append(games, GameModel{
					GameID:            gameModel.GameID,
					GameStatus:        gameModel.GameStatus,
					GameAttend:        gameModel.GameAttend,
					GameRound:         gameModel.GameRound,
					ControlGameStatus: gameModel.ControlGameStatus,
				})
			}
		}

		// 回傳所有場次遊戲狀態資訊訊息給前端
		b, _ := json.Marshal(GameParam{
			Games: games,
		})
		conn.WriteMessage(b)
	}
}

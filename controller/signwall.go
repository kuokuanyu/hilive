package controller

import (
	"hilive/modules/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// ShowSignWall 簽到牆(主持端) GET API
func (h *Handler) ShowSignWall(ctx *gin.Context) {
	var (
		path       = ctx.Request.URL.Path
		activityID = ctx.Query("activity_id")
		liffState  = ctx.Query("liff.state") // line裝置，liff url會顯示此參數，ex: liff.state=?activity_id=xxx
		// role       = ctx.Query("role")       // 角色
		routeFile string
	)
	if activityID == "" && liffState == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		return
	}

	// 處理activity_id
	if activityID == "" {
		if liffState != "" {
			if len(liffState) > 13 {
				liffState = liffState[13:] // 不讀取?activity_id字串
			}

			// qrcode都是給玩家掃描用的
			// role = "guest"

			// liff.state參數處理
			// 活動ID(#mst_challenge=xxx)
			activityID = liffState

			if strings.Contains(activityID, "#mst_challenge") {
				activityID = strings.Split(activityID, "#mst_challenge")[0]
			}
		}
	}

	if strings.Contains(path, "general") {
		routeFile = "./hilives/hilive/views/game/general.html"
	} else if strings.Contains(path, "threed") {
		routeFile = "./hilives/hilive/views/game/threed.html"
	} else if strings.Contains(path, "signname") {
		routeFile = "./hilives/hilive/views/game/signname.html"
	}
	// else if strings.Contains(path, "signname") && role == "guest" {
	// 	// 玩家端
	// 	// 判斷呈現內容參數
	// 	// 簽名牆設置資料(redis)
	// 	signnameModel, err := h.getSignnameSetting(true, activityID)
	// 	if err != nil {
	// 		h.executeErrorHTML(ctx, err.Error())
	// 		return
	// 	}

	// 	if signnameModel.SignnameContent == "write" {
	// 		// 讀取簽名牆模板
	// 		log.Println("讀取簽名牆模板")
	// 		routeFile = "./hilives/hilive/views/game/signname.html"
	// 	} else if signnameModel.SignnameContent == "message" {
	// 		// 讀取聊天室模板
	// 		log.Println("讀取聊天室模板")
	// 		routeFile = "./hilives/hilive/views/chatroom/style/default/chatroom.html"
	// 	}
	// }

	h.executeHTML(ctx, routeFile, executeParam{
		Route: route{
			POST: config.CHATROOM_RECORD_API_URL, // 新增聊天紀錄api
		},
	})
}

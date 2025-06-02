package controller

import (
	"fmt"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/parameter"
	"hilive/modules/table"

	"github.com/gin-gonic/gin"
)

// ShowWall 訊息區設置頁面(平台) GET API
func (h *Handler) ShowWall(ctx *gin.Context) {
	var (
		host       = ctx.Request.Host
		path       = ctx.Request.URL.Path
		prefix     = ctx.Param("__prefix")
		activityID = ctx.Query("activity_id")
		panel, _   = h.GetTable(ctx, prefix)
		user       = h.GetLoginUser(ctx.Request, "hilive_session")
		// panelInfo   table.PanelInfo
		formInfo table.FormInfo
		// newFormInfo table.FormInfo
		// editFormInfos []table.FormInfo
		// model, err         = models.DefaultActivityModel().SetConn(h.conn).Find("activity_id", activityID)
		param = parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
			[]string{"id"}, []string{"desc"}).SetPKs(user.UserID, activityID)
		canAdd, canEdit, canDelete bool
		patchURL                   = fmt.Sprintf(config.INTERACT_WALL_API_URL, prefix) // 編輯基本設置
		// sensitivityURL                      = fmt.Sprintf(config.INTERACT_WALL_URL, "message_sensitivity", activityID)
		sensitivityURL, backURL, postURL, putURL, putURL2, deleteURL string
		htmlTmpl                                                     string
		err                                                          error
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}
	if activityID == "" {
		h.executeErrorHTML(ctx, "錯誤: 無法辨識活動資訊，請輸入有效的活動ID")
		return
	}

	if prefix == "message" {
		// tempName = "content_block"
		htmlTmpl = "./hilives/hilive/views/cms/chat_message.html"
	} else if prefix == "topic" {
		// tempName = "content_block"
		htmlTmpl = "./hilives/hilive/views/cms/chat_topic.html"
	} else if prefix == "question" {
		postURL = config.QUESTION_GUEST_API_URL   // 嘉賓
		putURL = config.QUESTION_GUEST_API_URL    // 嘉賓
		deleteURL = config.QUESTION_GUEST_API_URL // 嘉賓
		htmlTmpl = "./hilives/hilive/views/cms/chat_question.html"
	} else if prefix == "danmu" {
		// tempName = "content_block"
		htmlTmpl = "./hilives/hilive/views/cms/chat_NM_danmu.html"
	} else if prefix == "specialdanmu" {
		// tempName = "content_block"
		htmlTmpl = "./hilives/hilive/views/cms/chat_SP_danmu.html"
	} else if prefix == "picture" {
		// tempName = "content"
		// htmlTmpl ="./hilives/hilive/views/cms/chat_message.html"
	} else if prefix == "holdscreen" {
		// tempName = "content_block"
		htmlTmpl = "./hilives/hilive/views/cms/chat_OC_danmu.html"
	} else if prefix == "message_check" {
		putURL = config.CHATROOM_RECORD_API_URL
		putURL2 = config.QUESTION_RECORD_API_URL
		sensitivityURL = fmt.Sprintf(config.INTERACT_WALL_URL, "message_sensitivity", activityID)
		htmlTmpl = "./hilives/hilive/views/cms/message_check.html"
	} else if prefix == "message_sensitivity" {
		backURL = fmt.Sprintf(config.INTERACT_WALL_URL, "message_check", activityID) // 活動訊息審核頁面
		postURL = fmt.Sprintf(config.INTERACT_WALL_API_URL, prefix)
		putURL = fmt.Sprintf(config.INTERACT_WALL_API_URL, prefix)
		deleteURL = fmt.Sprintf(config.INTERACT_WALL_API_URL, prefix)
		sensitivityURL = fmt.Sprintf(config.INTERACT_WALL_URL, "message_sensitivity", activityID)
		htmlTmpl = "./hilives/hilive/views/cms/message_sensitivity.html"
	}

	if prefix == "message" || prefix == "holdscreen" || prefix == "question" ||
		prefix == "specialdanmu" || prefix == "topic"{
		if formInfo, err = panel.GetSettingFormInfo(param, h.services,
			[]string{"user_id", "activity_id"}); err != nil {
			h.executeErrorHTML(ctx, "錯誤: 無法取得表單資訊，請重新整理頁面")
			return
		}
	}

	// 判斷是否有新增編輯刪除權限
	canAdd, canEdit, canDelete = h.checkPermission(user, activityID, path)

	// 提問嘉賓
	if prefix == "question" {
		// 嘉賓資訊
		// if panelInfo, err = panel.GetData(param, h.services, []string{"id"}); err != nil {
		// if panelInfo, err = panel.GetData(param, h.services); err != nil {
		// 	h.executeErrorHTML(ctx, "錯誤: 無法取得資料，請重新整理頁面")
		// 	return
		// }

		// 嘉賓表單
		// newFormInfo = panel.GetNewFormInfo(h.services, param, []string{"activity_id"})

		// 嘉賓編輯表單
		// for _, info := range panelInfo.InfoList {
		// 	var panel2 = h.tableList[prefix]()
		// 	var form table.FormInfo
		// 	editParam := parameter.GetParam(ctx.Request.URL, panel.GetInfo().PageSize,
		// 		[]string{"id"}, []string{"desc"}).SetPKs(fmt.Sprintf("%s", info["id"].Content), activityID)
		// 	if form, err = panel2.GetEditFormInfo(editParam, h.services,
		// 		[]string{"id", "activity_id"}); err != nil {
		// 		h.executeErrorHTML(ctx, "錯誤: 無法取得表單資訊，請重新整理頁面")
		// 		return
		// 	}
		// 	editFormInfos = append(editFormInfos, form)
		// }
	}

	h.executeHTML(ctx, htmlTmpl, executeParam{
		User:            user,
		ActivityID:      activityID,
		SettingFormInfo: formInfo,
		Route: route{
			Sensitivity: sensitivityURL, // 敏感詞頁面
			PATCH:       patchURL,
			POST:        postURL,
			PUT:         putURL,
			PUT2:        putURL2,
			DELETE:      deleteURL,
			Back:        backURL,
		},
		Token:     auth.AddToken(user.UserID),
		CanAdd:    canAdd,
		CanEdit:   canEdit,
		CanDelete: canDelete,
	})
}

// 判斷是否有新增編輯刪除權限
// 判斷用戶權限
// for _, permisssion := range user.Permissions {
// 	if canAdd == false || canEdit == false || canDelete == false {
// 		for _, httppath := range permisssion.HTTPPath {
// 			if httppath == "admin" || httppath == "*" {
// 				canAdd = true
// 				canEdit = true
// 				canDelete = true
// 				break
// 			} else if strings.Contains(path, httppath) {
// 				if strings.Contains(permisssion.HTTPMethod, "POST") {
// 					canAdd = true
// 				}
// 				if strings.Contains(permisssion.HTTPMethod, "PUT") {
// 					canEdit = true
// 				}
// 				if strings.Contains(permisssion.HTTPMethod, "DELETE") {
// 					canDelete = true
// 				}
// 			}
// 		}
// 	} else {
// 		break
// 	}
// }
// // 判斷活動權限
// if canAdd == false || canEdit == false || canDelete == false {
// 	for _, activity := range user.Activitys {
// 		if activity.ActivityID == activityID {
// 			for _, permisssion := range activity.Permissions {
// 				if canAdd == false || canEdit == false || canDelete == false {
// 					for _, httppath := range permisssion.HTTPPath {
// 						if httppath == "admin" || httppath == "*" {
// 							canAdd = true
// 							canEdit = true
// 							canDelete = true
// 							break
// 						} else if strings.Contains(path, httppath) {
// 							if strings.Contains(permisssion.HTTPMethod, "POST") {
// 								canAdd = true
// 							}
// 							if strings.Contains(permisssion.HTTPMethod, "PUT") {
// 								canEdit = true
// 							}
// 							if strings.Contains(permisssion.HTTPMethod, "DELETE") {
// 								canDelete = true
// 							}
// 						}
// 					}
// 				} else {
// 					break
// 				}
// 			}
// 		}
// 	}
// }

// PanelInfo:       panelInfo,
// FormInfo:        newFormInfo,
// FormInfos:       editFormInfos,
// if prefix == "question" || prefix == "danmu" {
// } else {
// 	h.execute(ctx, tempName, htmlTmpl, executeParam{
// 		User:            user,
// 		ActivityID:      activityID,
// 		// PanelInfo:       panelInfo,
// 		// FormInfo:        newFormInfo,
// 		SettingFormInfo: formInfo,
// 		// FormInfos:       editFormInfos,
// 		Route: route{
// 			PATCH:  fmt.Sprintf(config.INTERACT_WALL_API_URL, prefix), // 編輯基本設置
// 			POST:   config.QUESTION_GUEST_API_URL,                     // 新增提問牆嘉賓
// 			PUT:    config.QUESTION_GUEST_API_URL,                     // 編輯提問牆嘉賓
// 			DELETE: config.QUESTION_GUEST_API_URL,                     // 刪除提問牆嘉賓
// 		},
// 		Token:     auth.AddToken( user.UserID),
// 		CanAdd:    canAdd,
// 		CanEdit:   canEdit,
// 		CanDelete: canDelet

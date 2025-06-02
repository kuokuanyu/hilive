package controller

import (
	"fmt"
	"hilive/modules/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// ShowLineRichmenu LINE選單
func (h *Handler) ShowLineRichmenu(ctx *gin.Context) {
	var (
		host                            = ctx.Request.Host
		path                            = ctx.Request.URL.Path
		htmlTmpl, applysignURL, postURL string
		// userID                                       = ctx.Query("user_id")
		// prefix                                       = ctx.Param("__prefix")
		// htmlTmpl, activityURL, applysignURL, backURL string
	)
	if host == config.API_URL {
		h.executeErrorHTML(ctx, "錯誤: 網域請求發生問題")
		return
	}

	if strings.Contains(path, "/activity/search") {
		// prefix = "search"
		htmlTmpl = "./hilives/hilive/views/line/my_activity.html"
		// 活動查詢頁面、用戶報名簽到的活動頁面
		// 因為已經使用line選單，因此報名簽到url使用liff url
		applysignURL = fmt.Sprintf(config.HILIVES_ACTIVITY_ISEXIST_LIFF_URL, "")
	} else if strings.Contains(path, "/activity/create/require") {
		// prefix = "require"
		htmlTmpl = "./hilives/hilive/views/line/contact_form.html"
		postURL = config.ACTIVITY_CREATE_REQUIRE_API_URL
	} else if strings.Contains(path, "/activity/create") {
		// prefix = "create"
		htmlTmpl = "./hilives/hilive/views/line/contact.html"
	} else if strings.Contains(path, "/activity/case") {
		// prefix = "case"
		htmlTmpl = "./hilives/hilive/views/line/social_links.html"
	}

	// if strings.Contains(path, "/activity/search/activity") {
	// 	prefix = "search"
	// 	htmlTmpl = ""
	// 	// 返回、報名簽到頁面
	// 	backURL = fmt.Sprintf(config.LINE_RICHMENU_URL, prefix, userID)
	// 	applysignURL = fmt.Sprintf(config.HTTPS_LINE_LOGIN_REDIRECT_URL, host, "sign", "activity_id", "")
	// } else if strings.Contains(path, "/activity/search/applysign") {
	// 	prefix = "search"
	// 	htmlTmpl = ""
	// 	// 返回、報名簽到頁面
	// 	backURL = fmt.Sprintf(config.LINE_RICHMENU_URL, prefix, userID)
	// 	applysignURL = fmt.Sprintf(config.HTTPS_LINE_LOGIN_REDIRECT_URL, host, "sign", "activity_id", "")
	// } else

	h.executeHTML(ctx, htmlTmpl, executeParam{
		Route: route{
			// Activity:  activityURL,
			ApplySign: applysignURL,
			POST:      postURL,
			// Back:      backURL,
		},
	})
	return
}

package file

import (
	"mime/multipart"
	"strings"

	"hilive/modules/config"
)

// FileUploader 檔案儲存
type FileUploader struct {
	BasePath string // 儲存檔案路徑
}

// GetFileUploader 預設Uploader
func GetFileUploader() Uploader {
	return &FileUploader{
		config.GetStore().Path,
	}
}

// Upload 上傳
func (local *FileUploader) Upload(form *multipart.Form, apiPath, userID, activityID, gameID, game string) error {
	return Upload(func(fileObj *multipart.FileHeader, filename string) (string, error) {

		if strings.Contains(filename, "xlsx") || strings.Contains(filename, "xls") ||
			strings.Contains(filename, "csv") {
			// excel檔案類型: xlsx, xls, csv
			// 判斷是否為excel檔案，是的話存入uploads/excel中
			if err := SaveMultipartFile(fileObj, (*local).BasePath+"/excel/"+filename); err != nil {
				return "", err
			}
		} else {
			uploadPath := (*local).BasePath + "/"
			if apiPath == "/v1/admin/manager" || apiPath == "/v1/user" {
				// 用戶頁面
				uploadPath += userID + "/" + filename
			} else if apiPath == "/v1/line_user" {
				// 報名簽到人員更新姓名頭像頁面
				uploadPath += "line_user" + "/" + filename
			} else if apiPath == "/v1/auth/user" {
				// 補齊資料頁面
				uploadPath += "applysign" + "/" + filename
			} else if apiPath == "/v1/info/introduce" {
				// 活動介紹頁面
				uploadPath += userID + "/" + activityID + "/info/introduce/" + filename
			} else if apiPath == "/v1/info/guest" {
				// 活動嘉賓頁面
				uploadPath += userID + "/" + activityID + "/info/guest/" + filename
			} else if apiPath == "/v1/applysign/apply" {
				// 報名頁面
				uploadPath += userID + "/" + activityID + "/applysign/apply/" + filename
			} else if apiPath == "/v1/applysign/customize" {
				// 自定義欄位設置頁面
				uploadPath += userID + "/" + activityID + "/applysign/customize/" + filename
			} else if apiPath == "/v1/applysign/qrcode" {
				// QRcode自定義欄位頁面
				uploadPath += userID + "/" + activityID + "/applysign/qrcode/" + filename
			} else if apiPath == "/v1/interact/wall/message" {
				// 訊息區頁面
				uploadPath += userID + "/" + activityID + "/interact/wall/message/" + filename
			} else if apiPath == "/v1/interact/wall/topic" {
				// 主題區頁面
				uploadPath += userID + "/" + activityID + "/interact/wall/topic/" + filename
			} else if apiPath == "/v1/interact/wall/question" {
				// 提問區頁面
				uploadPath += userID + "/" + activityID + "/interact/wall/question/" + filename
			} else if apiPath == "/v1/interact/wall/danmu" {
				// 一般彈幕區頁面
				uploadPath += userID + "/" + activityID + "/interact/wall/danmu/" + filename
			} else if apiPath == "/v1/interact/sign/general" {
				// 一般簽到頁面
				uploadPath += userID + "/" + activityID + "/interact/sign/general/" + filename
			} else if apiPath == "/v1/interact/sign/threed" {
				// 3D簽到頁面
				uploadPath += userID + "/" + activityID + "/interact/sign/threed/" + filename
			} else if apiPath == "/v1/interact/sign/signname" {
				// 簽名牆頁面
				uploadPath += userID + "/" + activityID + "/interact/sign/signname/" + filename
			} else if strings.Contains(apiPath, "/v1/interact/sign") {
				// 判斷是否為獎品頁面
				if strings.Contains(apiPath, "/prize/form") {
					game = strings.Split(game, "_prize")[0]
				}

				// 判斷是否為投票選項頁面
				if strings.Contains(apiPath, "/vote/option") {
					game = strings.Split(game, "_option")[0]
				}

				// 遊戲頁面
				uploadPath += userID + "/" + activityID + "/interact/sign/" + game + "/" + gameID + "/" + filename
			} else if strings.Contains(apiPath, "/v1/interact/game") {
				// 判斷是否為獎品頁面
				if strings.Contains(apiPath, "/prize/form") {
					// fmt.Println("獎品頁面: ", strings.Split(game, "_prize")[0])
					game = strings.Split(game, "_prize")[0]
				}

				// 遊戲頁面
				uploadPath += userID + "/" + activityID + "/interact/game/" + game + "/" + gameID + "/" + filename

				// fmt.Println("獎品頁面: ", uploadPath)
			}
			// fmt.Println("上傳圖片路徑: ", uploadPath)

			// 圖片檔案存入uploads中
			if err := SaveMultipartFile(fileObj, uploadPath); err != nil {
				return "", err
			}
		}
		return filename, nil
	}, form, apiPath)
}

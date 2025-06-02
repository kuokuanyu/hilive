package models

import (
	"errors"
	"hilive/modules/db/command"
	"hilive/modules/utils"
	"strconv"
)

// EditQRcodeModel 資料表欄位
type EditQRcodeModel struct {
	ActivityID            string `json:"activity_id"`
	QRcodeLogoPicture     string `json:"qrcode_logo_picture"`
	QRcodeLogoSize        string `json:"qrcode_logo_size"`
	QRcodePicturePoint    string `json:"qrcode_picture_point"`
	QRcodeWhiteDistance   string `json:"qrcode_white_distance"`
	QRcodePointColor      string `json:"qrcode_point_color"`
	QRcodeBackgroundColor string `json:"qrcode_background_color"`
}

// UpdateQRcode 更新QRcode自定義資料
func (a ActivityModel) UpdateQRcode(model EditQRcodeModel) error {
	var (
		fieldValues = command.Value{}
		fields      = []string{
			"qrcode_logo_picture",
			"qrcode_logo_size",
			"qrcode_picture_point",
			"qrcode_white_distance",
			"qrcode_point_color",
			"qrcode_background_color",
		}
		// values = []string{
		// 	model.QRcodeLogoPicture,
		// 	model.QRcodeLogoSize,
		// 	model.QRcodePicturePoint,
		// 	model.QRcodeWhiteDistance,
		// 	model.QRcodePointColor,
		// 	model.QRcodeBackgroundColor,
		// }
	)
	if model.QRcodeLogoSize != "" {
		if size, err := strconv.ParseFloat(model.QRcodeLogoSize, 32); err != nil ||
			size < 0.2 || size > 0.5 {
			return errors.New("錯誤: LOGO尺寸資料發生問題，請輸入有效的LOGO尺寸資料(最小0.2，最大0.5)")
		}
	}
	if model.QRcodePicturePoint != "" {
		if model.QRcodePicturePoint != "open" && model.QRcodePicturePoint != "close" {
			return errors.New("錯誤: LOGO圖片後的點資料發生問題，請輸入有效的LOGO圖片後的點資料")
		}
	}
	if model.QRcodeWhiteDistance != "" {
		if size, err := strconv.Atoi(model.QRcodeWhiteDistance); err != nil ||
			size < 0 || size > 80 {
			return errors.New("錯誤: LOGO留白距離資料發生問題，請輸入有效的LOGO留白距離資料(最小0，最大80)")
		}
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	for _, key := range fields {
		if val, ok := data[key]; ok && val != "" {
			fieldValues[key] = val
		}
	}

	if len(fieldValues) == 0 {
		return nil
	}
	return a.Table(a.Base.TableName).
		Where("activity_id", "=", model.ActivityID).Update(fieldValues)
}

// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }

package controller

import (
	"fmt"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/file"
	"hilive/modules/response"
	"hilive/modules/utils"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// POST 新增 POST API
func (h *Handler) POST(ctx *gin.Context) {
	var (
		host = ctx.Request.Host
		// contentType = ctx.Request.Header.Get("content-type")
		path        = ctx.Request.URL.Path
		prefix      = ctx.Param("__prefix")
		values      = make(map[string][]string)
		activityID  = ctx.Request.FormValue("activity_id")
		gameID      = ctx.Request.FormValue("game_id")
		tokenUserID = ctx.Request.FormValue("user_id")
		token       = ctx.Request.FormValue("token")
		userID      string
		// err           error
	)
	if host != config.API_URL {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: "錯誤: 網域請求發生問題",
		})
		return

	}
	if path == "/v1/admin/manager" || path == "/v1/user" ||
		path == "/v1/activity" || strings.Contains(path, "/v1/interact/game") ||
		path == "/v1/interact/sign/vote/form" {
		// 因為管理員都可以幫忙設置(獎品頁面例外), 所以需要區分user_id參數
		userID = ctx.Request.FormValue("user") // 該活動場次的管理員資料

		// 判斷是否為獎品頁面, 獎品頁面目前沒有user參數
		if strings.Contains(path, "/prize/form") {
			userID = ctx.Request.FormValue("user_id")
		}

		if strings.Contains(path, "/v1/interact/game") {
			// 判斷是否為獎品頁面, 遊戲頁面才要隨機產生game_id
			if !strings.Contains(path, "/prize/form") {
				// 隨機產生game_id
				gameID = utils.UUID(20)

				// 建立遊戲場次, 新增遊戲資料夾
				os.MkdirAll(config.STORE_PATH+"/"+userID+"/"+activityID+"/interact/game/"+prefix+"/"+gameID, os.ModePerm)
			}
		}

		// 投票
		if path == "/v1/interact/sign/vote/form" {
			// 隨機產生game_id
			gameID = utils.UUID(20)

			// 建立遊戲場次, 新增遊戲資料夾
			os.MkdirAll(config.STORE_PATH+"/"+userID+"/"+activityID+"/interact/sign/"+prefix+"/"+gameID, os.ModePerm)
		}
	} else {
		userID = ctx.Request.FormValue("user_id")
	}

	if prefix == "" {
		if strings.Contains(path, "/v1/user") {
			prefix = "user"
		} else if strings.Contains(path, "activity") {
			prefix = "activity"
		} else if strings.Contains(path, "/v1/applysign/users") {
			prefix = "applysign_users" // 匹量匯入自定義報名簽到人員
		} else if strings.Contains(path, "/v1/applysign/user") {
			prefix = "applysign_user" // 單筆匯入自定義報名簽到人員
		} else if strings.Contains(path, "/v1/applysign") {
			prefix = "applysign"
		} else if strings.Contains(path, "/v1/chatroom/record") {
			prefix = "chatroom_record"
		} else if strings.Contains(path, "/v1/import/excel") {
			prefix = "excel"
		}
	} else if strings.Contains(path, "option_lists") {
		// 匹量插入
		prefix += "_option_lists"
	} else if strings.Contains(path, "option_list") {
		prefix += "_option_list"
	} else if strings.Contains(path, "option") {
		prefix += "_option"
	} else if strings.Contains(path, "special_officers") {
		// 匹量插入
		prefix += "_special_officers"
	} else if strings.Contains(path, "special_officer") {
		prefix += "_special_officer"
	} else if strings.Contains(path, "prize") {
		prefix += "_prize"
	} else if strings.Contains(path, "admin") {
		prefix = "admin_" + prefix
	}

	// 處理表單參數
	param, err := ctx.MultipartForm()
	if err != nil {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: fmt.Sprintf("錯誤: 表單參數發生問題，請重新操作，%s", err.Error()),
		})

		// 記錄請求的關鍵資訊
		// log.Println("Request Headers:", ctx.Request.Header)                  // multipart/form-data
		// log.Println("Content-Length:", ctx.Request.ContentLength)            // 不能為0
		// log.Println("Content-Type:", ctx.Request.Header.Get("Content-Type")) // Content-Type: multipart/form-data; boundary=----WebKitFormBoundaryABC123

		// models.DefaultErrorLogModel().SetDbConn(h.dbConn).
		// 	Add(models.EditErrorLogModel{
		// 		UserID: "後端除錯", Message: fmt.Sprintf("Request Headers: %s, Content-Length: %d, Content-Type: %s",
		// 			ctx.Request.Header, ctx.Request.ContentLength, ctx.Request.Header.Get("Content-Type")),
		// 	})

		// // 嘗試讀取部分 Body，確認是否真的有請求內容
		// bodyBytes, _ := io.ReadAll(ctx.Request.Body)
		// log.Println("Request Body (first 500 bytes):", string(bodyBytes[:500]))

		// models.DefaultErrorLogModel().SetDbConn(h.dbConn).
		// 	Add(models.EditErrorLogModel{
		// 		UserID: "後端除錯", Message: fmt.Sprintf("Request Body (first 300 bytes): %s", string(bodyBytes[:300])),
		// 	})

		// ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 10<<20) // 10MB 限制
		// err := ctx.Request.ParseMultipartForm(10 << 20)                              // 10MB
		// if err != nil {
		// 	log.Println("解析表單時發生錯誤(10MB):", err)

		// 	models.DefaultErrorLogModel().SetDbConn(h.dbConn).
		// 		Add(models.EditErrorLogModel{
		// 			UserID: "後端除錯", Message: fmt.Sprintf("解析表單時發生錯誤(10MB): %s", err.Error()),
		// 		})
		// }

		// if ctx.Request.MultipartForm == nil {
		// 	log.Println("Multipart 表單為空(10MB)")

		// 	models.DefaultErrorLogModel().SetDbConn(h.dbConn).
		// 		Add(models.EditErrorLogModel{
		// 			UserID: "後端除錯", Message: "Multipart 表單為空(10MB)",
		// 		})
		// } else {
		// 	fmt.Println("成功解析 Multipart 表單(10MB)")

		// 	models.DefaultErrorLogModel().SetDbConn(h.dbConn).
		// 		Add(models.EditErrorLogModel{
		// 			UserID: "後端除錯", Message: "成功解析 Multipart 表單(10MB)",
		// 		})
		// }

		return
	}

	// 上傳圖片、檔案
	if len(param.File) > 0 {
		if err := file.GetFileEngine(config.FILE_ENGINE).Upload(ctx.Request.MultipartForm, path,
			userID, activityID, gameID, prefix); err != nil {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 上傳檔案發生問題，請重新上傳檔案",
			})
			return

		}
	}

	values = param.Value
	if strings.Contains(path, "/v1/interact/game") ||
		path == "/v1/interact/sign/vote/form" {
		// 建立遊戲場次, 將隨機產生的game_id加入參數中
		values["game_id"] = append(values["game_id"], gameID)
	}

	// 聊天紀錄、舉辦活動需求、新增用戶、新增簽名牆資料不需要token驗證
	if prefix != "chatroom_record" && prefix != "require" && prefix != "user" && prefix != "signname" {
		if !auth.CheckToken(token, tokenUserID) {
			// log.Println("參數: ", activityID, userID, token)
			response.BadRequest(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: Token驗證發生問題，請輸入有效的Token值",
			})
			return

		}
	}

	// 匯入excel檔案處理
	if prefix == "excel" {
		var (
			excel = values["excel"]
			game  = values["game"]
		)

		// 判斷遊戲類型
		if len(game) > 0 && game[0] == "" {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 取得遊戲類型發生問題，請重新操作",
			})
			return

		}

		// 判斷是否上傳excel檔案
		if len(excel) == 0 || (len(excel) > 0 && excel[0] == "") {

			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 讀取excel檔案發生問題(無檔案)，請重新操作",
			})
			return

		}

		//
		file, err := excelize.OpenFile(config.STORE_PATH + "/excel/" + excel[0])
		if err != nil {
			response.Error(ctx, h.dbConn, models.EditErrorLogModel{
				UserID:  userID,
				Method:  ctx.Request.Method,
				Path:    ctx.Request.URL.Path,
				Message: "錯誤: 讀取excel檔案發生問題(開啟excel檔)，請重新操作",
			})
			return

		}
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Println(err)
			}
		}()

		if game[0] == "vote_option_list" { // 投票選項名單
			var (
				users    = make([]models.GameVoteOptionListModel, 0)
				optionID = ctx.Request.FormValue("option_id")
				isLeader bool
			)
			if activityID == "" || gameID == "" || optionID == "" {
				response.Error(ctx, h.dbConn, models.EditErrorLogModel{
					UserID:  userID,
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 參數發生問題(活動、遊戲、選項)，請重新操作",
				})
				return

			}

			// 取得該選項名單資料
			lists, err := models.DefaultGameVoteOptionListModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Find(false, gameID, optionID, "")
			if err != nil {
				response.Error(ctx, h.dbConn, models.EditErrorLogModel{
					UserID:  userID,
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 查詢選項名單資料發生問題，請重新操作",
				})
				return
			}

			// 檢查所有選項名單中是否存在隊長資料
			for _, list := range lists {
				if list.Leader == "leader" {
					// 存在隊長資料
					isLeader = true
					break // 停止
				}
			}

			// excel欄位判斷
			for i := 0; i < 10000; i++ {
				rowIndex := strconv.Itoa(i + 3)
				aa, _ := file.GetCellValue("Sheet1", "A"+rowIndex) // 姓名
				bb, _ := file.GetCellValue("Sheet1", "B"+rowIndex) // 隊長

				// 必填欄位不為空(姓名)，加入陣列
				if aa != "" {
					if bb == "leader" && isLeader {
						// 已有隊長
						response.Error(ctx, h.dbConn, models.EditErrorLogModel{
							UserID:  userID,
							Method:  ctx.Request.Method,
							Path:    ctx.Request.URL.Path,
							Message: "錯誤: 該選項已存在隊長(一個選項隊長只會有一位)，請重新匯入",
						})
					} else if bb == "leader" && !isLeader {
						// 加入隊長資料
						isLeader = true
					}

					users = append(users, models.GameVoteOptionListModel{
						Name:   aa,
						Leader: bb,
					})

				} else {
					// 結束excel判斷
					break
				}
			}

			// 判斷是否還未設置隊長
			if !isLeader && len(users) > 0 {
				// 將第一位人員資料改為隊長
				users[0].Leader = "leader"
			}

			response.OkWithData(ctx, users)
			return

		} else if game[0] == "applysign_users" { // 自定義報名簽到人員
			if activityID == "" {
				response.Error(ctx, h.dbConn, models.EditErrorLogModel{
					UserID:  userID,
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 活動參數發生問題，請重新操作",
				})
				return

			}

			// 取得自定義欄位資料
			customizeModel, err := models.DefaultCustomizeModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				Find(activityID)
			if err != nil || customizeModel.ID == 0 {
				response.Error(ctx, h.dbConn, models.EditErrorLogModel{
					UserID:  userID,
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 查詢自定義欄位資料發生問題，請重新操作",
				})
				return
			}

			// 是否自定義輸入驗證碼
			var customizepassword string
			if customizeModel.CustomizePassword == "open" {
				customizepassword = "true"
			} else {
				customizepassword = "false"
			}

			// 自定義匯入報名簽到人員
			var (
				users   = make([]models.ApplysignUserModel, 0)
				uniques = []string{ // 唯一值
					customizeModel.Ext1Unique, customizeModel.Ext2Unique, customizeModel.Ext3Unique, customizeModel.Ext4Unique,
					customizeModel.Ext5Unique, customizeModel.Ext6Unique, customizeModel.Ext7Unique, customizeModel.Ext8Unique,
					customizeModel.Ext9Unique, customizeModel.Ext10Unique,
					"true", // 驗證碼
				}
				requires = []string{ // 是否必填
					customizeModel.Ext1Required, customizeModel.Ext2Required, customizeModel.Ext3Required, customizeModel.Ext4Required,
					customizeModel.Ext5Required, customizeModel.Ext6Required, customizeModel.Ext7Required, customizeModel.Ext8Required,
					customizeModel.Ext9Required, customizeModel.Ext10Required,
					customizepassword, // 驗證碼
				}
				exts = make([][]string, 11)
			)

			// 取得該活動所有報名簽到人員資料
			applysigns, err := models.DefaultApplysignModel().
				SetConn(h.dbConn, h.redisConn, h.mongoConn).
				FindAll(activityID, "", "", "", 0, 0)
			if err != nil {
				response.Error(ctx, h.dbConn, models.EditErrorLogModel{
					UserID:  userID,
					Method:  ctx.Request.Method,
					Path:    ctx.Request.URL.Path,
					Message: "錯誤: 查詢簽到人員資料發生問題，請重新操作",
				})
				return
			}

			// 取得原有報名簽到資料的唯一值
			for _, applysign := range applysigns {
				applysginExts := []string{
					applysign.Ext1, applysign.Ext2, applysign.Ext3, applysign.Ext4,
					applysign.Ext5, applysign.Ext6, applysign.Ext7, applysign.Ext8,
					applysign.Ext9, applysign.Ext10,
					applysign.ExtPassword, // 驗證碼
				}

				for n, ext := range applysginExts {
					// 判斷自定義欄位是否為唯一值，是的話加入陣列中
					if uniques[n] == "true" {
						exts[n] = append(exts[n], ext)
					}
				}
			}

			// excel欄位判斷
			for i := 0; i < 10000; i++ {

				rowIndex := strconv.Itoa(i + 3)
				aa, _ := file.GetCellValue("Sheet1", "A"+rowIndex) // 姓名
				bb, _ := file.GetCellValue("Sheet1", "B"+rowIndex) // 電話
				cc, _ := file.GetCellValue("Sheet1", "C"+rowIndex) // 信箱
				dd, _ := file.GetCellValue("Sheet1", "D"+rowIndex) // 驗證碼
				ee, _ := file.GetCellValue("Sheet1", "E"+rowIndex) // ex1
				ff, _ := file.GetCellValue("Sheet1", "F"+rowIndex) // ex2
				gg, _ := file.GetCellValue("Sheet1", "G"+rowIndex) // ex3
				hh, _ := file.GetCellValue("Sheet1", "H"+rowIndex) // ex4
				ii, _ := file.GetCellValue("Sheet1", "I"+rowIndex) // ex5
				jj, _ := file.GetCellValue("Sheet1", "J"+rowIndex) // ex6
				kk, _ := file.GetCellValue("Sheet1", "K"+rowIndex) // ex7
				ll, _ := file.GetCellValue("Sheet1", "L"+rowIndex) // ex8
				mm, _ := file.GetCellValue("Sheet1", "M"+rowIndex) // ex9
				nn, _ := file.GetCellValue("Sheet1", "N"+rowIndex) // ex10

				// 必填欄位不為空，加入陣列
				// 判斷辨識ID是否存在陣列中(不能重複)，存在則不加入該筆資料
				if aa != "" {
					// 電話必填或開啟簡訊功能時電話欄位不能為空
					if (customizeModel.ExtPhoneRequired == "true" || customizeModel.PushPhoneMessage == "open") &&
						bb == "" {
						response.Error(ctx, h.dbConn, models.EditErrorLogModel{
							UserID:  userID,
							Method:  ctx.Request.Method,
							Path:    ctx.Request.URL.Path,
							Message: "錯誤: 電話欄位為必填，請輸入有效的資料",
						})
						return
					}

					if bb != "" {
						if len(bb) > 2 {
							if !strings.Contains(bb[:2], "09") || len(bb) != 10 {
								response.Error(ctx, h.dbConn, models.EditErrorLogModel{
									UserID:  userID,
									Method:  ctx.Request.Method,
									Path:    ctx.Request.URL.Path,
									Message: "錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼",
								})
								return
							}
						} else {
							response.Error(ctx, h.dbConn, models.EditErrorLogModel{
								UserID:  userID,
								Method:  ctx.Request.Method,
								Path:    ctx.Request.URL.Path,
								Message: "錯誤: 格式必須為台灣地區的手機號碼，如:09XXXXXXXX，請輸入有效的手機號碼",
							})
							return
						}
					}

					// 信箱必填或開啟郵件功能時信箱欄位不能為空
					if (customizeModel.ExtEmailRequired == "true" || customizeModel.SendMail == "open") &&
						cc == "" {
						response.Error(ctx, h.dbConn, models.EditErrorLogModel{
							UserID:  userID,
							Method:  ctx.Request.Method,
							Path:    ctx.Request.URL.Path,
							Message: "錯誤: 信箱欄位為必填，請輸入有效的資料",
						})
						return
					}

					if cc != "" && !strings.Contains(cc, "@") {
						response.Error(ctx, h.dbConn, models.EditErrorLogModel{
							UserID:  userID,
							Method:  ctx.Request.Method,
							Path:    ctx.Request.URL.Path,
							Message: "錯誤: 電子郵件地址中必須包含「@」，請輸入有效的電子郵件地址",
						})
						return
					}

					// 判斷自定義欄位是否為唯一值
					for n, uniques := range uniques {
						values := []string{ee, ff, gg, hh, ii, jj, kk, ll, mm, nn, dd}

						if n == 10 {
							// 驗證碼欄位驗證判斷
							// 是否自定義輸入驗證碼
							if customizepassword == "true" {
								if values[n] == "" {
									response.Error(ctx, h.dbConn, models.EditErrorLogModel{
										UserID:  userID,
										Method:  ctx.Request.Method,
										Path:    ctx.Request.URL.Path,
										Message: "錯誤: 驗證碼欄位不能為空，請輸入有效資料",
									})
									return
								}

								if utils.InArray(exts[n], values[n]) {
									// 唯一值存在，錯誤
									response.Error(ctx, h.dbConn, models.EditErrorLogModel{
										UserID:  userID,
										Method:  ctx.Request.Method,
										Path:    ctx.Request.URL.Path,
										Message: "錯誤: 驗證碼欄位不能重複，請輸入有效資料",
									})
									return
								} else {
									exts[n] = append(exts[n], values[n])
								}
							} else if customizepassword == "false" {
								// 不自定義輸入驗證碼，驗證碼隨機產生
								dd = utils.RandomNumber(6)

								// 驗證第一次
								if utils.InArray(exts[n], dd) {
									// 唯一值存在，錯誤
									// 再隨機產生一次驗證碼
									dd = utils.RandomNumber(6)

									// 驗證第二次
									if utils.InArray(exts[n], dd) {
										dd = utils.RandomNumber(6)
									}

									// 驗證第三次
									if utils.InArray(exts[n], dd) {
										dd = utils.RandomNumber(6)
									}

									// 驗證第四次
									if utils.InArray(exts[n], dd) {
										dd = utils.RandomNumber(6)
									}

									// 驗證第五次
									if utils.InArray(exts[n], dd) {
										response.Error(ctx, h.dbConn, models.EditErrorLogModel{
											UserID:  userID,
											Method:  ctx.Request.Method,
											Path:    ctx.Request.URL.Path,
											Message: "錯誤: 驗證碼欄位不能重複(隨機產生驗證碼重複)，請再匯入一次excel檔案",
										})
										return
									}

									exts[n] = append(exts[n], dd)
								} else {
									exts[n] = append(exts[n], dd)
								}
							}
						} else {
							// 其它自定義欄位驗證判斷
							if uniques == "true" {
								if values[n] == "" {
									response.Error(ctx, h.dbConn, models.EditErrorLogModel{
										UserID:  userID,
										Method:  ctx.Request.Method,
										Path:    ctx.Request.URL.Path,
										Message: "錯誤: 唯一值欄位不能為空，請輸入有效資料",
									})
									return
								}

								if utils.InArray(exts[n], values[n]) {
									// 唯一值存在，錯誤
									response.Error(ctx, h.dbConn, models.EditErrorLogModel{
										UserID:  userID,
										Method:  ctx.Request.Method,
										Path:    ctx.Request.URL.Path,
										Message: "錯誤: 唯一值欄位不能重複，請輸入有效資料",
									})
									return
								} else {
									exts[n] = append(exts[n], values[n])
								}
							}

							if requires[n] == "true" && values[n] == "" {
								response.Error(ctx, h.dbConn, models.EditErrorLogModel{
									UserID:  userID,
									Method:  ctx.Request.Method,
									Path:    ctx.Request.URL.Path,
									Message: "錯誤: 必填欄位不能為空，請輸入有效資料",
								})
								return
							}
						}

					}

					users = append(users, models.ApplysignUserModel{
						Name:        aa,
						Phone:       bb,
						ExtEmail:    cc,
						ExtPassword: dd,
						Ext1:        ee,
						Ext2:        ff,
						Ext3:        gg,
						Ext4:        hh,
						Ext5:        ii,
						Ext6:        jj,
						Ext7:        kk,
						Ext8:        ll,
						Ext9:        mm,
						Ext10:       nn,
					})

				} else {
					// 結束excel判斷
					break
				}
			}

			response.OkWithData(ctx, users)
			return
		} else if game[0] == "QA" {
			var (
				questions = make([]models.QuestionModel, 20)
				// qa        = make([]string, 80) // 題目設置
				// index     int64
			)

			// 快問快答題目設置
			for i := 0; i < 20; i++ {
				questions[i].Options = make([]string, 0)
				rowIndex := strconv.Itoa(i + 2)
				a, _ := file.GetCellValue("Sheet1", "A"+rowIndex)
				b, _ := file.GetCellValue("Sheet1", "B"+rowIndex)
				c, _ := file.GetCellValue("Sheet1", "C"+rowIndex)
				d, _ := file.GetCellValue("Sheet1", "D"+rowIndex)
				e, _ := file.GetCellValue("Sheet1", "E"+rowIndex)
				f, _ := file.GetCellValue("Sheet1", "F"+rowIndex)
				g, _ := file.GetCellValue("Sheet1", "G"+rowIndex)

				// 某一格欄位為空，停止題目設置
				if a == "" || b == "" || c == "" ||
					d == "" || e == "" || f == "" || g == "" {
					continue
				}
				if f == "A" || f == "a" {
					f = "0"
				} else if f == "B" || f == "b" {
					f = "1"
				} else if f == "C" || f == "c" {
					f = "2"
				} else if f == "D" || f == "d" {
					f = "3"
				} else {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(正確選項只能填寫ABCD)，請重新操作",
					})
					return

				}

				questions[i] = models.QuestionModel{
					Question: a,
					Options:  []string{b, c, d, e},
					Answer:   f,
					Score:    utils.GetInt64(g, 0),
				}

			}

			response.OkWithData(ctx, models.GameModel{Questions: questions})
			return
		} else if game[0] == "introduce" {
			// 活動介紹
			var (
				introduces = make([]models.IntroduceModel, 0)
			)

			// excel欄位判斷
			for i := 0; i < 100; i++ {
				rowIndex := strconv.Itoa(i + 2)
				a, _ := file.GetCellValue("Sheet1", "A"+rowIndex)
				b, _ := file.GetCellValue("Sheet1", "B"+rowIndex)
				c, _ := file.GetCellValue("Sheet1", "C"+rowIndex)

				// 必填欄位為空，停止設置
				if b == "" {
					break
				}

				if utf8.RuneCountInString(a) > 20 {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(介紹標題欄位資料上限為20個字元)，請重新操作",
					})
					return

				}

				if utf8.RuneCountInString(c) > 200 {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(介紹內容欄位資料上限為200個字元)，請重新操作",
					})
					return

				}

				if b != "文字" && b != "圖片" {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(內容類型欄位資料必須為文字、圖片)，請重新操作",
					})
					return

				}

				introduces = append(introduces, models.IntroduceModel{
					Title:         a,
					IntroduceType: b,
					Content:       c,
				})

			}

			response.OkWithData(ctx, introduces)
			return
		} else if game[0] == "schedule" {
			// 活動行程
			var (
				schedules = make([]models.ScheduleModel, 0)
			)

			// excel欄位判斷
			for i := 0; i < 100; i++ {
				rowIndex := strconv.Itoa(i + 2)
				a, _ := file.GetCellValue("Sheet1", "A"+rowIndex)
				b, _ := file.GetCellValue("Sheet1", "B"+rowIndex)
				c, _ := file.GetCellValue("Sheet1", "C"+rowIndex)
				d, _ := file.GetCellValue("Sheet1", "D"+rowIndex)
				e, _ := file.GetCellValue("Sheet1", "E"+rowIndex)

				// 必填欄位為空，停止設置
				if a == "" || b == "" || c == "" || d == "" {
					break
				}

				if utf8.RuneCountInString(a) > 20 {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(行程標題資料上限為20個字元)，請重新操作",
					})
					return
				}

				if utf8.RuneCountInString(e) > 200 {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(行程內容資料上限為200個字元)，請重新操作",
					})
					return

				}

				if !CompareTime(c, d) {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(時間欄位格式必須為XX:XX且結束時間必須大於開始時間)，請重新操作",
					})
					return

				}

				schedules = append(schedules, models.ScheduleModel{
					Title:        a,
					ScheduleDate: b,
					StartTime:    c,
					EndTime:      d,
					Content:      e,
				})

			}

			response.OkWithData(ctx, schedules)
			return
		} else if game[0] == "guest" {
			// 活動嘉賓
			var (
				guests = make([]models.GuestModel, 0)
			)

			// excel欄位判斷
			for i := 0; i < 100; i++ {
				rowIndex := strconv.Itoa(i + 2)
				a, _ := file.GetCellValue("Sheet1", "A"+rowIndex)
				b, _ := file.GetCellValue("Sheet1", "B"+rowIndex)
				c, _ := file.GetCellValue("Sheet1", "C"+rowIndex)
				d, _ := file.GetCellValue("Sheet1", "D"+rowIndex)

				// 必填欄位為空，停止設置
				if b == "" {
					break
				}

				if utf8.RuneCountInString(b) > 20 {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(姓名欄位資料上限為20個字元)，請重新操作",
					})
					return

				}

				if utf8.RuneCountInString(c) > 20 {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(頭銜欄位資料上限為20個字元)，請重新操作",
					})
					return

				}

				if utf8.RuneCountInString(d) > 200 {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(詳細資訊欄位資料上限為200個字元)，請重新操作",
					})
					return

				}

				guests = append(guests, models.GuestModel{
					Avatar:    a,
					Name:      b,
					Introduce: c,
					Detail:    d,
				})

			}

			response.OkWithData(ctx, guests)
			return
		} else if game[0] == "material" {
			// 活動資料
			var (
				materials = make([]models.MaterialModel, 0)
			)

			// excel欄位判斷
			for i := 0; i < 100; i++ {
				rowIndex := strconv.Itoa(i + 2)
				a, _ := file.GetCellValue("Sheet1", "A"+rowIndex)
				b, _ := file.GetCellValue("Sheet1", "B"+rowIndex)
				c, _ := file.GetCellValue("Sheet1", "C"+rowIndex)

				// 必填欄位為空，停止設置
				if a == "" {
					break
				}

				if utf8.RuneCountInString(a) > 20 {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(資料標題欄位資料上限為20個字元)，請重新操作",
					})
					return

				}

				if utf8.RuneCountInString(b) > 200 {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(資料說明欄位資料上限為200個字元)，請重新操作",
					})
					return

				}

				materials = append(materials, models.MaterialModel{
					Title:     a,
					Introduce: b,
					Link:      c,
				})

			}

			response.OkWithData(ctx, materials)
			return
		} else if game[0] == "prize" {
			// 獎品
			var (
				prizes = make([]models.EditPrizeModel, 0)
			)

			// excel欄位判斷
			for i := 0; i < 100; i++ {
				rowIndex := strconv.Itoa(i + 2)
				a, _ := file.GetCellValue("Sheet1", "A"+rowIndex)
				b, _ := file.GetCellValue("Sheet1", "B"+rowIndex)
				c, _ := file.GetCellValue("Sheet1", "C"+rowIndex)
				d, _ := file.GetCellValue("Sheet1", "D"+rowIndex)
				e, _ := file.GetCellValue("Sheet1", "E"+rowIndex)
				f, _ := file.GetCellValue("Sheet1", "F"+rowIndex)
				g, _ := file.GetCellValue("Sheet1", "G"+rowIndex)

				// 必填欄位為空，停止設置
				if b == "" || c == "" || d == "" || e == "" || f == "" || g == "" {
					break
				}

				if utf8.RuneCountInString(b) > 20 {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(獎品名稱欄位)，請重新操作",
					})
					return

				}

				if c != "頭獎" && c != "貳獎" &&
					c != "參獎" && c != "普通獎" &&
					c != "謝謝參與" {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(獎品類型欄位，頭獎、貳獎、參獎、普通獎、謝謝參與)，請重新操作",
					})
					return
				}

				if _, err := strconv.Atoi(d); err != nil {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(獎品價值欄位)，請重新操作",
					})
					return

				}

				if _, err := strconv.Atoi(e); err != nil {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(獎品數量欄位)，請重新操作",
					})
					return

				}

				if f != "現場發放" && f != "宅配郵寄ˇ" &&
					f != "謝謝參與" {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(獎品發送方式欄位，現場發放、宅配郵寄、謝謝參與)，請重新操作",
					})
					return

				}

				if utf8.RuneCountInString(g) > 8 {
					response.Error(ctx, h.dbConn, models.EditErrorLogModel{
						UserID:  userID,
						Method:  ctx.Request.Method,
						Path:    ctx.Request.URL.Path,
						Message: "錯誤: 讀取excel檔案發生問題(獎品密碼欄位)，請重新操作",
					})
					return

				}

				prizes = append(prizes, models.EditPrizeModel{
					PrizePicture:  a,
					PrizeName:     b,
					PrizeType:     c,
					PrizePrice:    d,
					PrizeAmount:   e,
					PrizeMethod:   f,
					PrizePassword: g,
				})

			}

			response.OkWithData(ctx, prizes)
			return
		}
	}

	table, _ := h.GetTable(ctx, prefix)
	if err := table.InsertData(values); err != nil {
		response.Error(ctx, h.dbConn, models.EditErrorLogModel{
			UserID:  userID,
			Method:  ctx.Request.Method,
			Path:    ctx.Request.URL.Path,
			Message: err.Error(),
		})
		return

	}

	response.OkWithURL(ctx, "")
}

// CompareTime 時間比較
func CompareTime(start, end string) bool {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04", "2006-01-02 "+start, time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04", "2006-01-02 "+end, time.Local)
	boolTime := endTime.After(startTime) && startTime.Before(endTime)
	if !boolTime && start != end {
		return false
	}
	return true

}

// @Summary 新增聊天紀錄
// @Tags Chatroom_Record
// @version 1.0
// @Accept  mpfd
// @param user_id formData string true "用戶ID"
// @param activity_id formData string true "活動ID"
// @@@param name formData string true "姓名"
// @@@param avatar formData string false "頭像"
// @param message_type formData string true "訊息類型" Enums(normal-message, normal-barrage, special-barrage, occupy-barrage)
// @param message_style formData string true "訊息風格" Enums(default, happybirthday)
// @param message_price formData integer true "訊息價格"
// @param message_effect formData string true "訊息效果"
// @param message_status formData string true "訊息審核狀態" Enums(yes, no, review)
// @param message formData string false "訊息"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /chatroom/record [post]
func (h *Handler) POSTChatroomRecord(ctx *gin.Context) {
}

// @Summary 匯入excel檔案
// @Tags Excel
// @version 1.0
// @Accept  mpfd
// @param activity_id formData string false "activity_id"
// @param game formData string true "遊戲類型" Enums(QA, introduce, schedule, guest, material, prize, applysign_users, vote_option_list)
// @param excel formData file true "excel檔案"
// @param user_id formData string true "用戶ID"
// @param token formData string true "CSRF Token"
// @Success 200 {array} response.ResponseWithData
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /import/excel [post]
func (h *Handler) POSTImportExcel(ctx *gin.Context) {
}

// qa[index] = a
// qa[index+1] = strings.Join([]string{b, c, d, e}, "&&&")
// qa[index+2] = f
// qa[index+3] = g

// 下一題題目設置的index間隔為4
// index += 4

// gameModel.Questions = append(gameModel.Questions,
// 	models.QuestionModel{Question: qa[0], Options: strings.Split(qa[1], "&&&"), Answer: qa[2], Score: utils.GetInt64(qa[3], 0)},
// 	models.QuestionModel{Question: qa[4], Options: strings.Split(qa[5], "&&&"), Answer: qa[6], Score: utils.GetInt64(qa[7], 0)},
// 	models.QuestionModel{Question: qa[8], Options: strings.Split(qa[9], "&&&"), Answer: qa[10], Score: utils.GetInt64(qa[11], 0)},
// 	models.QuestionModel{Question: qa[12], Options: strings.Split(qa[13], "&&&"), Answer: qa[14], Score: utils.GetInt64(qa[15], 0)},
// 	models.QuestionModel{Question: qa[16], Options: strings.Split(qa[17], "&&&"), Answer: qa[18], Score: utils.GetInt64(qa[19], 0)},
// 	models.QuestionModel{Question: qa[20], Options: strings.Split(qa[21], "&&&"), Answer: qa[22], Score: utils.GetInt64(qa[23], 0)},
// 	models.QuestionModel{Question: qa[24], Options: strings.Split(qa[25], "&&&"), Answer: qa[26], Score: utils.GetInt64(qa[27], 0)},
// 	models.QuestionModel{Question: qa[28], Options: strings.Split(qa[29], "&&&"), Answer: qa[30], Score: utils.GetInt64(qa[31], 0)},
// 	models.QuestionModel{Question: qa[32], Options: strings.Split(qa[33], "&&&"), Answer: qa[34], Score: utils.GetInt64(qa[35], 0)},
// 	models.QuestionModel{Question: qa[36], Options: strings.Split(qa[37], "&&&"), Answer: qa[38], Score: utils.GetInt64(qa[39], 0)},
// 	models.QuestionModel{Question: qa[40], Options: strings.Split(qa[41], "&&&"), Answer: qa[42], Score: utils.GetInt64(qa[43], 0)},
// 	models.QuestionModel{Question: qa[44], Options: strings.Split(qa[45], "&&&"), Answer: qa[46], Score: utils.GetInt64(qa[47], 0)},
// 	models.QuestionModel{Question: qa[48], Options: strings.Split(qa[49], "&&&"), Answer: qa[50], Score: utils.GetInt64(qa[51], 0)},
// 	models.QuestionModel{Question: qa[52], Options: strings.Split(qa[53], "&&&"), Answer: qa[54], Score: utils.GetInt64(qa[55], 0)},
// 	models.QuestionModel{Question: qa[56], Options: strings.Split(qa[57], "&&&"), Answer: qa[58], Score: utils.GetInt64(qa[59], 0)},
// 	models.QuestionModel{Question: qa[60], Options: strings.Split(qa[61], "&&&"), Answer: qa[62], Score: utils.GetInt64(qa[63], 0)},
// 	models.QuestionModel{Question: qa[64], Options: strings.Split(qa[65], "&&&"), Answer: qa[66], Score: utils.GetInt64(qa[67], 0)},
// 	models.QuestionModel{Question: qa[68], Options: strings.Split(qa[69], "&&&"), Answer: qa[70], Score: utils.GetInt64(qa[71], 0)},
// 	models.QuestionModel{Question: qa[72], Options: strings.Split(qa[73], "&&&"), Answer: qa[74], Score: utils.GetInt64(qa[75], 0)},
// 	models.QuestionModel{Question: qa[76], Options: strings.Split(qa[77], "&&&"), Answer: qa[78], Score: utils.GetInt64(qa[79], 0)})

// @@@Summary 新增簽名牆資料
// @@@Tags Sign
// @@@version 1.0
// @@@Accept  mpfd
// @@@param user_id formData string true "用戶ID"
// @@@param activity_id formData string true "活動ID"
// @@@param picture formData string true "標題(上限為20個字元)" minlength(1) maxlength(20)
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/sign/signname [post]
func (h *Handler) POSTSignname(ctx *gin.Context) {
}

// @@@Summary 新增拔河遊戲獎品資料(form-data)
// @@@Tags Tugofwar Prize
// @@@version 1.0
// @@@Accept  mpfd
// @@@param user_id formData string true "用戶ID"
// @@@param activity_id formData string true "活動ID"
// @@@param game_id formData string true "遊戲ID"
// @@@param prize_name formData string true "名稱(上限為20個字元)" minlength(1) maxlength(20)
// @@@param prize_type formData string true "類型" Enums(first, second, third, general)
// @@@param prize_picture formData file false "照片"
// @@@param prize_method formData string true "兌獎方式" Enums(site, mail)
// @@@param prize_password formData string true "兌獎密碼"
// @@@param prize_amount formData integer true "數量"
// @@@param prize_price formData integer true "價值"
// @@@param token formData string true "CSRF Token"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/tugofwar/prize/form [post]
func (h *Handler) POSTTugofwarPrize(ctx *gin.Context) {
}

// ---@Summary 新增遊戲人員資料(form-data)
// ---@Tags Attend Staff
// ---@version 1.0
// ---@Accept  mpfd
// ---@param user_id formData string true "用戶ID"
// ---@param activity_id formData string true "活動ID"
// ---@param game_id formData string true "遊戲ID"
// ---@param name formData string true "姓名"
// ---@param avatar formData string false "頭像"
// ---@param round formData integer true "輪次"
// ---@param status formData string true "狀態" Enums(success, fail)
// ---@param black formData string true "黑名單" Enums(yes, no)
// ---@Success 200 {array} response.Response
// ---@Failure 500 {array} response.ResponseInternalServerError
// ---@Router /staffmanage/attend/form [post]
func (h *Handler) POSTAttend(ctx *gin.Context) {
}

// ---@Summary 新增中獎人員資料(form-data)
// ---@Tags Winning Staff
// ---@version 1.0
// ---@Accept  mpfd
// ---@param user_id formData string true "用戶ID"
// ---@param activity_id formData string true "活動ID"
// ---@param game_id formData string true "遊戲ID"
// ---@param prize_id formData string true "獎品ID"
// ---@param name formData string true "姓名"
// ---@param avatar formData string false "頭像"
// ---@param prize_name formData string true "獎品名稱(上限為20個字元)" minlength(1) maxlength(20)
// ---@param picture formData string false "照片"
// ---@param price formData integer true "價值"
// ---@param method formData string true "兌獎方式" Enums(site, mail, thanks)
// ---@param password formData string true "兌獎密碼(最多八個字元)" minlength(1) maxlength(8)
// ---@param round formData integer true "輪次"
// ---@param win_time formData string true "中獎時間(西元年-月-日 時:分)"
// ---@param status formData string true "兌獎狀態" Enums(yes, no)
// ---@param white formData string true "白名單" Enums(yes, no)
// ---@Success 200 {array} response.Response
// ---@Failure 500 {array} response.ResponseInternalServerError
// ---@Router /staffmanage/winning/form [post]
// func (h *Handler) POSTWinning(ctx *gin.Context) {
// }

// else if contentType == "application/json" {
// 	// fmt.Println("application/json")
// 	if prefix == "black" {
// 		var model models.NewBlackStaffModel
// 		err = ctx.BindJSON(&model)
// 		userID = model.UserID
// 		token = model.Token

// 		values["activity_id"] = []string{model.ActivityID}
// 		values["game_id"] = []string{model.GameID}
// 		values["game"] = []string{model.Game}
// 		values["user_id"] = []string{model.UserID}
// 		values["line_id"] = []string{model.LINEID}
// 		values["reason"] = []string{model.Reason}
// 	} else if strings.Contains(prefix, "_prize") {
// 		// 獎品相關參數
// 		var model models.NewPrizeModel
// 		err = ctx.BindJSON(&model)
// 		userID = model.UserID
// 		token = model.Token

// 		values["activity_id"] = []string{model.ActivityID}
// 		values["game_id"] = []string{model.GameID}
// 		values["prize_name"] = []string{model.PrizeName}
// 		values["prize_type"] = []string{model.PrizeType}
// 		values["prize_picture"] = []string{model.PrizePicture}
// 		values["prize_amount"] = []string{model.PrizeAmount}
// 		values["prize_price"] = []string{model.PrizePrice}
// 		values["prize_method"] = []string{model.PrizeMethod}
// 		values["prize_password"] = []string{model.PrizePassword}
// 	} else if !strings.Contains(prefix, "_prize") {
// 		// 遊戲相關參數
// 		var model models.NewGameModel
// 		err = ctx.BindJSON(&model)
// 		userID = model.UserID
// 		token = model.Token

// 		values["activity_id"] = []string{model.ActivityID}
// 		values["title"] = []string{model.Title}
// 		values["game_type"] = []string{model.GameType}
// 		values["limit_time"] = []string{model.LimitTime}
// 		values["second"] = []string{model.Second}
// 		values["max_people"] = []string{model.MaxPeople}
// 		values["people"] = []string{model.People}
// 		values["max_times"] = []string{model.MaxTimes}
// 		values["allow"] = []string{model.Allow}
// 		values["percent"] = []string{model.Percent}
// 		values["first_prize"] = []string{model.FirstPrize}
// 		values["second_prize"] = []string{model.SecondPrize}
// 		values["third_prize"] = []string{model.ThirdPrize}
// 		values["general_prize"] = []string{model.GeneralPrize}
// 		values["topic"] = []string{model.Topic}
// 		values["skin"] = []string{model.Skin}
// 		values["music"] = []string{model.Music}
// 		values["display_name"] = []string{model.DisplayName}

// 		// 敲敲樂自定義
// 		values["whack_mole_host_background"] = []string{model.WhackMoleHostBackground}
// 		values["whack_mole_guest_background"] = []string{model.WhackMoleGuestBackground}
// 		values["whack_mole_dollar_rat_picture"] = []string{model.WhackMoleDollarRatPicture}
// 		values["whack_mole_redpack_rat_picture"] = []string{model.WhackMoleRedpackRatPicture}
// 		values["whack_mole_bomb_picture"] = []string{model.WhackMoleBombPicture}
// 		values["whack_mole_rat_hole_picture"] = []string{model.WhackMoleRatHolePicture}
// 		values["whack_mole_rock_picture"] = []string{model.WhackMoleRockPicture}
// 		values["whack_mole_rank_picture"] = []string{model.WhackMoleRankPicture}
// 		values["whack_mole_rank_background"] = []string{model.WhackMoleRankBackground}

// 		// 搖號抽獎自定義
// 		values["draw_numbers_background"] = []string{model.DrawNumbersBackground}
// 		values["draw_numbers_title"] = []string{model.DrawNumbersTitle}
// 		values["draw_numbers_gift_inside_picture"] = []string{model.DrawNumbersGiftInsidePicture}
// 		values["draw_numbers_gift_outside_picture"] = []string{model.DrawNumbersGiftOutsidePicture}
// 		values["draw_numbers_prize_left_button"] = []string{model.DrawNumbersPrizeLeftButton}
// 		values["draw_numbers_prize_right_button"] = []string{model.DrawNumbersPrizeRightButton}
// 		values["draw_numbers_prize_leftright_button"] = []string{model.DrawNumbersPrizeLeftrightButton}
// 		values["draw_numbers_addpeople_no_button"] = []string{model.DrawNumbersAddpeopleNoNutton}
// 		values["draw_numbers_addpeople_yes_button"] = []string{model.DrawNumbersAddpeopleYesButton}
// 		values["draw_numbers_people_background"] = []string{model.DrawNumbersPeopleBackground}
// 		values["draw_numbers_add_people"] = []string{model.DrawNumbersAddPeople}
// 		values["draw_numbers_reduce_people"] = []string{model.DrawNumbersReducePeople}
// 		values["draw_numbers_winner_background"] = []string{model.DrawNumbersWinnerBackground}
// 		values["draw_numbers_blackground"] = []string{model.DrawNumbersBlackground}
// 		values["draw_numbers_go_button"] = []string{model.DrawNumbersGoButton}
// 		values["draw_numbers_open_winner_button"] = []string{model.DrawNumbersOpenWinnerButton}
// 		values["draw_numbers_close_winner_button"] = []string{model.DrawNumbersCloseWinnerButton}
// 		values["draw_numbers_current_people"] = []string{model.DrawNumbersCurrentPeople}
// 		// 動圖
// 		values["draw_numbers_gacha_machine"] = []string{model.DrawNumbersGachaMachine}
// 		values["draw_numbers_hood"] = []string{model.DrawNumbersHood}
// 		values["draw_numbers_body"] = []string{model.DrawNumbersBody}
// 		values["draw_numbers_gacha"] = []string{model.DrawNumbersGacha}

// 		// 鑑定師自定義
// 		values["monopoly_screen_again_button"] = []string{model.MonopolyScreenAgainButton}                             // 主持端再來一輪按鈕
// 		values["monopoly_screen_top6_title"] = []string{model.MonopolyScreenTop6Title}                                 // 主持端前六名標題
// 		values["monopoly_screen_end_button"] = []string{model.MonopolyScreenEndButton}                                 // 主持端結束遊戲按鈕
// 		values["monopoly_screen_start_button"] = []string{model.MonopolyScreenStartButton}                             // 主持端開始遊戲按鈕
// 		values["monopoly_screen_gaming_background_png"] = []string{model.MonopolyScreenGamingBackgroundPng}            // 主持端遊戲中背景
// 		values["monopoly_screen_round_countdown"] = []string{model.MonopolyScreenRoundCountdown}                       // 主持端遊戲中輪次倒數
// 		values["monopoly_screen_winner_list"] = []string{model.MonopolyScreenWinnerList}                               // 主持端遊戲和結算中獎列表
// 		values["monopoly_screen_rank_border"] = []string{model.MonopolyScreenRankBorder}                               // 主持端遊戲和結算名次框
// 		values["monopoly_player_carton"] = []string{model.MonopolyPlayerCarton}                                        // 玩家端遊戲中下滑紙箱
// 		values["monopoly_player_any_start_text"] = []string{model.MonopolyPlayerAnyStartText}                          // 玩家端按任意處開始文字
// 		values["monopoly_player_scoreboard"] = []string{model.MonopolyPlayerScoreboard}                                // 玩家端計分和時間板
// 		values["monopoly_player_wait_start_text"] = []string{model.MonopolyPlayerWaitStartText}                        // 玩家端計分和時間板
// 		values["monopoly_player_transparent_background"] = []string{model.MonopolyPlayerTransparentBackground}         // 玩家端開始場景半透明黑底
// 		values["monopoly_player_pile_objects"] = []string{model.MonopolyPlayerPileObjects}                             // 玩家端滑動物件堆
// 		values["monopoly_player_gaming_background"] = []string{model.MonopolyPlayerGamingBackground}                   // 玩家端遊戲中背景
// 		values["monopoly_add_points"] = []string{model.MonopolyAddPoints}                                              // 上滑小標示
// 		values["monopoly_deduct_points"] = []string{model.MonopolyDeductPoints}                                        // 下滑小標示
// 		values["monopoly_player_background_dynamic"] = []string{model.MonopolyPlayerBackgroundDynamic}                 // 玩家端背景
// 		values["monopoly_player_answer_effect"] = []string{model.MonopolyPlayerAnswerEffect}                           // 玩家端遊戲中答對或錯特效
// 		values["monopoly_background_and_gold"] = []string{model.MonopolyBackgroundAndGold}                             // 主持端背景和玩家端金銅條
// 		values["monopoly_screen_redpack_seal"] = []string{model.MonopolyScreenRedpackSeal}                             // 主持端紅包袋封口
// 		values["monopoly_screen_again_button_background"] = []string{model.MonopolyScreenAgainButtonBackground}        // 主持端再來一輪按鈕底
// 		values["monopoly_screen_end_info_skin"] = []string{model.MonopolyScreenEndInfoSkin}                            // 主持端結算3名後玩家資訊木框
// 		values["monopoly_screen_end_npc"] = []string{model.MonopolyScreenEndNpc}                                       // 主持端結算吉祥物
// 		values["monopoly_screen_top_stair"] = []string{model.MonopolyScreenTopStair}                                   // 主持端結算前三名台階
// 		values["monopoly_screen_top_info_skin"] = []string{model.MonopolyScreenTopInfoSkin}                            // 主持端結算前三名資訊框
// 		values["monopoly_screen_top_avatar_skin"] = []string{model.MonopolyScreenTopAvatarSkin}                        // 主持端結算前三名頭像框
// 		values["monopoly_screen_end_background"] = []string{model.MonopolyScreenEndBackground}                         // 主持端結算背景
// 		values["monopoly_screen_start_npc_dialog"] = []string{model.MonopolyScreenStartNpcDialog}                      // 主持端開始畫面人物對話框
// 		values["monopoly_screen_leaderboard"] = []string{model.MonopolyScreenLeaderboard}                              // 主持端遊戲中排行榜
// 		values["monopoly_screen_round_background"] = []string{model.MonopolyScreenRoundBackground}                     // 主持端遊戲中輪次底
// 		values["monopoly_screen_start_end_button_background"] = []string{model.MonopolyScreenStartEndButtonBackground} // 主持端開始和結束按鈕底
// 		values["monopoly_screen_start_background"] = []string{model.MonopolyScreenStartBackground}                     // 主持端開始背景
// 		values["monopoly_screen_start_right_top_decoration"] = []string{model.MonopolyScreenStartRightTopDecoration}   // 主持端開始右上裝飾
// 		values["monopoly_player_tip_arrow"] = []string{model.MonopolyPlayerTipArrow}                                   // 玩家端遊戲中提示箭頭
// 		values["monopoly_player_npc_dialog"] = []string{model.MonopolyPlayerNpcDialog}                                 // 玩家端人物對話框
// 		values["monopoly_player_join_button_background"] = []string{model.MonopolyPlayerJoinButtonBackground}          // 玩家端加入遊戲按鈕底
// 		values["monopoly_player_join_background"] = []string{model.MonopolyPlayerJoinBackground}                       // 玩家端加入遊戲背景
// 		values["monopoly_player_redpack_space"] = []string{model.MonopolyPlayerRedpackSpace}                           // 玩家端紅包袋白底
// 		values["monopoly_player_redpack_seal"] = []string{model.MonopolyPlayerRedpackSeal}                             // 玩家端紅包袋封口
// 		values["monopoly_player_redpack_background"] = []string{model.MonopolyPlayerRedpackBackground}                 // 玩家端紅包袋背景
// 		values["monopoly_player_money_piles"] = []string{model.MonopolyPlayerMoneyPiles}                               // 玩家端鈔票堆
// 		values["monopoly_player_background"] = []string{model.MonopolyPlayerBackground}                                // 玩家端遊戲背景
// 		values["monopoly_player_title"] = []string{model.MonopolyPlayerTitle}                                          // 玩家端遊戲標題
// 		values["monopoly_npc"] = []string{model.MonopolyNpc}                                                           // 代表人物
// 		values["monopoly_button"] = []string{model.MonopolyButton}                                                     // 按鈕
// 		values["monopoly_screen_top_light"] = []string{model.MonopolyScreenTopLight}                                   // 主持端前三名發亮
// 		values["monopoly_screen_end_revolving_light"] = []string{model.MonopolyScreenEndRevolvingLight}                // 主持端結算背景旋轉燈
// 		values["monopoly_screen_end_ribbon"] = []string{model.MonopolyScreenEndRibbon}                                 // 主持端結算彩帶
// 		values["monopoly_player_gaming_redpack"] = []string{model.MonopolyPlayerGamingRedpack}                         // 玩家端遊戲中紅包
// 		values["monopoly_screen_gaming_redpack"] = []string{model.MonopolyScreenGamingRedpack}                         // 主持端遊戲中紅包和玩家端紅包
// 		values["monopoly_screen_top_after_player_info"] = []string{model.MonopolyScreenTopAfterPlayerInfo}             // 主持端結算3名後玩家資訊框
// 		values["monopoly_screen_top_front_player_info"] = []string{model.MonopolyScreenTopFrontPlayerInfo}             // 主持端結算前三名資訊框
// 		values["monopoly_screen_rank"] = []string{model.MonopolyScreenRank}                                            // 主持端遊戲中排行榜
// 		values["monopoly_screen_npc_dialog"] = []string{model.MonopolyScreenNpcDialog}                                 // 主持端對話框
// 		values["monopoly_screen_left_bottom_decoration"] = []string{model.MonopolyScreenLeftBottomDecoration}          // 主持端遊戲中裝飾小物件左下
// 		values["monopoly_player_basket_background"] = []string{model.MonopolyPlayerBasketBackground}                   // 玩家端竹籃背景
// 		values["monopoly_player_gaming_carrots"] = []string{model.MonopolyPlayerGamingCarrots}                         // 玩家端遊戲中紅蘿蔔堆
// 		values["monopoly_button_background"] = []string{model.MonopolyButtonBackground}                                // 按鈕背景
// 		values["monopoly_screen_end_background_dynamic"] = []string{model.MonopolyScreenEndBackgroundDynamic}          // 主持端遊戲中和結算背景
// 		values["monopoly_screen_start_background_dynamic"] = []string{model.MonopolyScreenStartBackgroundDynamic}      // 主持端開始背景
// 		values["monopoly_player_gaming_background_dynamic"] = []string{model.MonopolyPlayerGamingBackgroundDynamic}    // 玩家端遊戲背景
// 		values["monopoly_picking_carrots_and_carrots"] = []string{model.MonopolyPickingCarrotsAndCarrots}              // 主持端遊戲中採蘿蔔和玩家端蘿蔔
// 		values["monopoly_player_top_info"] = []string{model.MonopolyPlayerTopInfo}                                     // 玩家端上方資訊
// 		values["monopoly_player_search_prize_background"] = []string{model.MonopolyPlayerSearchPrizeBackground}        // 玩家端查看獎品背景
// 		values["monopoly_player_food_waste_bin"] = []string{model.MonopolyPlayerFoodWasteBin}                          // 玩家端遊戲中廚餘桶
// 		values["monopoly_screen_end_dynamic"] = []string{model.MonopolyScreenEndDynamic}                               // 主持端結算動圖
// 		values["monopoly_screen_timer"] = []string{model.MonopolyScreenTimer}                                          // 主持端遊戲中計時器
// 		values["monopoly_player_start_gaming_eyecatch"] = []string{model.MonopolyPlayerStartGamingEyecatch}            // 玩家端開始和遊戲過場
// 		values["monopoly_gaming_dynamic_and_fish"] = []string{model.MonopolyGamingDynamicAndFish}                      // 主持端遊戲中動圖和玩家端魚
// 		values["monopoly_screen_gaming_background_jpg"] = []string{model.MonopolyScreenGamingBackgroundJpg}            // 主持端遊戲中背景

// 		// 快問快答自定義
// 		values["qa_mascot"] = []string{model.QAMascot}
// 		values["qa_host_start_background"] = []string{model.QAHostStartBackground}
// 		values["qa_host_game_background"] = []string{model.QAHostGameBackground}
// 		values["qa_host_end_background"] = []string{model.QAhostEndBackground}
// 		values["qa_game_top_1"] = []string{model.QAGameTop1}
// 		values["qa_game_top_2"] = []string{model.QAGameTop2}
// 		values["qa_game_top_3"] = []string{model.QAGameTop3}
// 		values["qa_game_top_4"] = []string{model.QAGameTop4}
// 		values["qa_game_top_5"] = []string{model.QAGameTop5}
// 		values["qa_end_top_1"] = []string{model.QAEndTop1}
// 		values["qa_end_top_2"] = []string{model.QAEndTop2}
// 		values["qa_end_top_3"] = []string{model.QAEndTop3}
// 		values["qa_end_top"] = []string{model.QAEndTop}
// 		values["qa_host_start_game_button"] = []string{model.QAHostStartGameButton}
// 		values["qa_host_pause_countdown_button"] = []string{model.QAHostPauseCountdownButton}
// 		values["qa_host_continue_countdown_button"] = []string{model.QAHostContinueCountdownButton}
// 		values["qa_host_start_answer_button"] = []string{model.QAHostStartAnswerButton}
// 		values["qa_host_see_answer_button"] = []string{model.QAHostSeeAnswerButton}
// 		values["qa_host_next_question_button"] = []string{model.QAHostNextQuestionButton}
// 		values["qa_host_end_game_button"] = []string{model.QAHostEndGameButton}
// 		values["qa_host_again_game_button"] = []string{model.QAHostAgainGameButton}
// 		values["qa_player_start_background"] = []string{model.QAPlayerStartBackground}
// 		values["qa_player_game_background"] = []string{model.QAPlayerGameBackground}
// 		values["qa_player_join_game_button"] = []string{model.QAPlayerJoinGameButton}
// 		values["qa_player_select_answer_button"] = []string{model.QAPlayerSelectAnswerButton}
// 		values["qa_player_confirm_answer_button"] = []string{model.QAPlayerConfirmAnswerButton}
// 		values["qa_player_confirm_status_button"] = []string{model.QAPlayerConfirmStatusButton}

// 		// 搖紅包自定義
// 		values["redpack_screen_again_button"] = []string{model.RedpackScreenAgainButton}
// 		values["redpack_screen_background"] = []string{model.RedpackScreenBackground}
// 		values["redpack_screen_end_button"] = []string{model.RedpackScreenEndButton}
// 		values["redpack_screen_prize_list"] = []string{model.RedpackScreenPrizeList}
// 		values["redpack_screen_prize_redpack"] = []string{model.RedpackScreenPrizeRedpack}
// 		values["redpack_screen_start_button"] = []string{model.RedpackScreenStartButton}
// 		values["redpack_screen_title"] = []string{model.RedpackScreenTitle}
// 		values["redpack_screen_gaming_list"] = []string{model.RedpackScreenGamingList}
// 		values["redpack_screen_gaming_list_background"] = []string{model.RedpackScreenGamingListBackground}
// 		values["redpack_screen_ema"] = []string{model.RedpackScreenEma}
// 		values["redpack_screen_new_list"] = []string{model.RedpackScreenNewList}
// 		values["redpack_screen_lantern1"] = []string{model.RedpackScreenLantern1}
// 		values["redpack_screen_lantern2"] = []string{model.RedpackScreenLantern2}
// 		values["redpack_player_join_button"] = []string{model.RedpackPlayerJoinButton}
// 		values["redpack_player_search_prize_background"] = []string{model.RedpackPlayerSearchPrizeBackground}
// 		values["redpack_player_background"] = []string{model.RedpackPlayerBackground}
// 		values["redpack_player_title"] = []string{model.RedpackPlayerTitle}
// 		values["redpack_player_lantern"] = []string{model.RedpackPlayerLantern}
// 		// 動圖
// 		values["redpack_screen_lucky_bag"] = []string{model.RedpackScreenLuckyBag}
// 		values["redpack_screen_money_piles"] = []string{model.RedpackScreenMoneyPiles}
// 		values["redpack_player_shake"] = []string{model.RedpackPlayerShake}
// 		values["redpack_player_lucky_bag"] = []string{model.RedpackPlayerLuckyBag}
// 		values["redpack_player_money_piles"] = []string{model.RedpackPlayerMoneyPiles}
// 		values["redpack_screen_background_dynamic"] = []string{model.RedpackScreenBackgroundDynamic}
// 		values["redpack_player_background_dynamic"] = []string{model.RedpackPlayerBackgroundDynamic}
// 		values["redpack_player_ready"] = []string{model.RedpackPlayerReady}
// 		// 音樂
// 		values["redpack_bgm_start"] = []string{model.RedpackBgmStart}
// 		values["redpack_bgm_gaming"] = []string{model.RedpackBgmGaming}
// 		values["redpack_bgm_end"] = []string{model.RedpackBgmEnd}

// 		// 套紅包自定義
// 		values["ropepack_screen_prize_list"] = []string{model.RopepackScreenPrizeList}
// 		values["ropepack_screen_again_button"] = []string{model.RopepackScreenAgainButton}
// 		values["ropepack_screen_background"] = []string{model.RopepackScreenBackground}
// 		values["ropepack_screen_decoration"] = []string{model.RopepackScreenDecoration}
// 		values["ropepack_screen_end_button"] = []string{model.RopepackScreenEndButton}
// 		values["ropepack_screen_end_prize_list"] = []string{model.RopepackScreenEndPrizeList}
// 		values["ropepack_screen_prize_redpack"] = []string{model.RopepackScreenPrizeRedpack}
// 		values["ropepack_screen_start_logo"] = []string{model.RopepackScreenStartLogo}
// 		values["ropepack_screen_start_button"] = []string{model.RopepackScreenStartButton}
// 		values["ropepack_screen_prize_skin_red"] = []string{model.RopepackScreenPrizeSkinRed}
// 		values["ropepack_screen_prize_skin_green"] = []string{model.RopepackScreenPrizeSkinGreen}
// 		values["ropepack_player_join_logo"] = []string{model.RopepackPlayerJoinLogo}
// 		values["ropepack_player_join_button"] = []string{model.RopepackPlayerJoinButton}
// 		values["ropepack_player_background"] = []string{model.RopepackPlayerBackground}
// 		values["ropepack_player_ready_redpack1"] = []string{model.RopepackPlayerReadyRedpack1}
// 		values["ropepack_player_ready_redpack2"] = []string{model.RopepackPlayerReadyRedpack2}
// 		values["ropepack_player_ready_background"] = []string{model.RopepackPlayerReadyBackground}
// 		values["ropepack_player_title"] = []string{model.RopepackPlayerTitle}
// 		// 動圖
// 		values["ropepack_screen_background_effect"] = []string{model.RopepackScreenBackgroundEffect}
// 		values["ropepack_player_ropepack_button"] = []string{model.RopepackPlayerRopepackButton}
// 		values["ropepack_player_finger"] = []string{model.RopepackPlayerFinger}
// 		values["ropepack_redpack"] = []string{model.RopepackRedpack}

// 		// 遊戲抽獎自定義
// 		values["lottery_screen_prizer"] = []string{model.LotteryScreenPrizer}
// 		values["lottery_screen_mascot"] = []string{model.LotteryScreenMascot}
// 		values["lottery_screen_background"] = []string{model.LotteryScreenBackground}
// 		values["lottery_screen_prize_notify"] = []string{model.LotteryScreenPrizeNotify}
// 		values["lottery_screen_select_input"] = []string{model.LotteryScreenSelectInput}
// 		values["lottery_screen_close_prize_notify_button"] = []string{model.LotteryScreenClosePrizeNotifyButton}
// 		values["lottery_player_background"] = []string{model.LotteryPlayerBackground}
// 		values["lottery_player_rules"] = []string{model.LotteryPlayerRules}
// 		values["lottery_jiugongge_grid"] = []string{model.LotteryJiugonggeGrid}
// 		values["lottery_jiugongge_start_button"] = []string{model.LotteryJiugonggeStartButton}
// 		values["lottery_turntable_start_button"] = []string{model.LotteryTurntableStartButton}
// 		values["lottery_turntable_roulette"] = []string{model.LotteryTurntableRoulette}
// 		// 動圖
// 		values["lottery_get_prize"] = []string{model.LotteryGetPrize}
// 		values["lottery_jiugongge_border"] = []string{model.LotteryJiugonggeBorder}
// 		values["lottery_jiugongge_title"] = []string{model.LotteryJiugonggeTitle}
// 		values["lottery_turntable_border"] = []string{model.LotteryTurntableBorder}
// 		values["lottery_turntable_title"] = []string{model.LotteryTurntableTitle}

// 		values["qa_1"] = []string{model.QA1}
// 		values["qa_1_options"] = []string{model.QA1Options}
// 		values["qa_1_answer"] = []string{model.QA1Answer}
// 		values["qa_1_score"] = []string{model.QA1Score}

// 		values["qa_2"] = []string{model.QA2}
// 		values["qa_2_options"] = []string{model.QA2Options}
// 		values["qa_2_answer"] = []string{model.QA2Answer}
// 		values["qa_2_score"] = []string{model.QA2Score}

// 		values["qa_3"] = []string{model.QA3}
// 		values["qa_3_options"] = []string{model.QA3Options}
// 		values["qa_3_answer"] = []string{model.QA3Answer}
// 		values["qa_3_score"] = []string{model.QA3Score}

// 		values["qa_4"] = []string{model.QA4}
// 		values["qa_4_options"] = []string{model.QA4Options}
// 		values["qa_4_answer"] = []string{model.QA4Answer}
// 		values["qa_4_score"] = []string{model.QA4Score}

// 		values["qa_5"] = []string{model.QA5}
// 		values["qa_5_options"] = []string{model.QA5Options}
// 		values["qa_5_answer"] = []string{model.QA5Answer}
// 		values["qa_5_score"] = []string{model.QA5Score}

// 		values["qa_6"] = []string{model.QA6}
// 		values["qa_6_options"] = []string{model.QA6Options}
// 		values["qa_6_answer"] = []string{model.QA6Answer}
// 		values["qa_6_score"] = []string{model.QA6Score}

// 		values["qa_7"] = []string{model.QA7}
// 		values["qa_7_options"] = []string{model.QA7Options}
// 		values["qa_7_answer"] = []string{model.QA7Answer}
// 		values["qa_7_score"] = []string{model.QA7Score}

// 		values["qa_8"] = []string{model.QA8}
// 		values["qa_8_options"] = []string{model.QA8Options}
// 		values["qa_8_answer"] = []string{model.QA8Answer}
// 		values["qa_8_score"] = []string{model.QA8Score}

// 		values["qa_9"] = []string{model.QA9}
// 		values["qa_9_options"] = []string{model.QA9Options}
// 		values["qa_9_answer"] = []string{model.QA9Answer}
// 		values["qa_9_score"] = []string{model.QA9Score}

// 		values["qa_10"] = []string{model.QA10}
// 		values["qa_10_options"] = []string{model.QA10Options}
// 		values["qa_10_answer"] = []string{model.QA10Answer}
// 		values["qa_10_score"] = []string{model.QA10Score}

// 		values["qa_11"] = []string{model.QA11}
// 		values["qa_11_options"] = []string{model.QA11Options}
// 		values["qa_11_answer"] = []string{model.QA11Answer}
// 		values["qa_11_score"] = []string{model.QA11Score}

// 		values["qa_12"] = []string{model.QA12}
// 		values["qa_12_options"] = []string{model.QA12Options}
// 		values["qa_12_answer"] = []string{model.QA12Answer}
// 		values["qa_12_score"] = []string{model.QA12Score}

// 		values["qa_13"] = []string{model.QA13}
// 		values["qa_13_options"] = []string{model.QA13Options}
// 		values["qa_13_answer"] = []string{model.QA13Answer}
// 		values["qa_13_score"] = []string{model.QA13Score}

// 		values["qa_14"] = []string{model.QA14}
// 		values["qa_14_options"] = []string{model.QA14Options}
// 		values["qa_14_answer"] = []string{model.QA14Answer}
// 		values["qa_14_score"] = []string{model.QA14Score}

// 		values["qa_15"] = []string{model.QA15}
// 		values["qa_15_options"] = []string{model.QA15Options}
// 		values["qa_15_answer"] = []string{model.QA15Answer}
// 		values["qa_15_score"] = []string{model.QA15Score}

// 		values["qa_16"] = []string{model.QA16}
// 		values["qa_16_options"] = []string{model.QA16Options}
// 		values["qa_16_answer"] = []string{model.QA16Answer}
// 		values["qa_16_score"] = []string{model.QA16Score}

// 		values["qa_17"] = []string{model.QA17}
// 		values["qa_17_options"] = []string{model.QA17Options}
// 		values["qa_17_answer"] = []string{model.QA17Answer}
// 		values["qa_17_score"] = []string{model.QA17Score}

// 		values["qa_18"] = []string{model.QA18}
// 		values["qa_18_options"] = []string{model.QA18Options}
// 		values["qa_18_answer"] = []string{model.QA18Answer}
// 		values["qa_18_score"] = []string{model.QA18Score}

// 		values["qa_19"] = []string{model.QA19}
// 		values["qa_19_options"] = []string{model.QA19Options}
// 		values["qa_19_answer"] = []string{model.QA19Answer}
// 		values["qa_19_score"] = []string{model.QA19Score}

// 		values["qa_20"] = []string{model.QA20}
// 		values["qa_20_options"] = []string{model.QA20Options}
// 		values["qa_20_answer"] = []string{model.QA20Answer}
// 		values["qa_20_score"] = []string{model.QA20Score}

// 		values["total_qa"] = []string{model.TotalQA}
// 		values["qa_second"] = []string{model.QASecond}
// 	}
// 	if err != nil {
// 		if strings.Contains(ctx.GetHeader("Accept"), "json") {
// 			response.Error(ctx, "錯誤: 無法解析JSON參數，請重新操作")
// 			return
// 		} else {
// 			h.executeErrorHTML(ctx, "錯誤: 無法解析JSON參數，請重新操作")
// 			return
// 		}
// 	}
// }

// @@@Summary 新增遊戲抽獎遊戲資料(json)
// @@@Tags Lottery
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.NewGameModel true "Lottery Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/lottery/json [post]
func (h *Handler) POSTJSONLottery(ctx *gin.Context) {
}

// @@@Summary 新增遊戲抽獎獎品資料(json)
// @@@Tags Lottery Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.NewPrizeModel true "Lottery Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/lottery/prize/json [post]
func (h *Handler) POSTJSONLotteryPrize(ctx *gin.Context) {
}

// @@@Summary 新增搖紅包遊戲資料(json)
// @@@Tags Redpack
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.NewGameModel true "Redpack Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/redpack/json [post]
func (h *Handler) POSTJSONRedpack(ctx *gin.Context) {
}

// @@@Summary 新增搖紅包獎品資料(json)
// @@@Tags Redpack Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.NewPrizeModel true "Redpack Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/redpack/prize/json [post]
func (h *Handler) POSTJSONRedpackPrize(ctx *gin.Context) {
}

// @@@Summary 新增套紅包遊戲資料(json)
// @@@Tags Ropepack
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.NewGameModel true "Ropepack Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/ropepack/json [post]
func (h *Handler) POSTJSONRopepack(ctx *gin.Context) {
}

// @@@Summary 新增套紅包獎品資料(json)
// @@@Tags Ropepack Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.NewPrizeModel true "Ropepack Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/ropepack/prize/json [post]
func (h *Handler) POSTJSONRopepackPrize(ctx *gin.Context) {
}

// @@@Summary 新增敲敲樂遊戲資料(json)
// @@@Tags Whack_Mole
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.NewGameModel true "WhackMole Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/whack_mole/json [post]
func (h *Handler) POSTJSONWhackMole(ctx *gin.Context) {
}

// @@@Summary 新增敲敲樂獎品資料(json)
// @@@Tags Whack_Mole Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.NewPrizeModel true "WhackMole Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/whack_mole/prize/json [post]
func (h *Handler) POSTJSONWhackMolePrize(ctx *gin.Context) {
}

// @@@Summary 新增搖號抽獎遊戲資料(json)
// @@@Tags Draw_Numbers
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.NewGameModel true "Draw_Numbers Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/draw_numbers/json [post]
func (h *Handler) POSTJSONDrawNumbers(ctx *gin.Context) {
}

// @@@Summary 新增搖號抽獎獎品資料(json)
// @@@Tags Draw_Numbers Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.NewGameModel true "Draw_Numbers Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/draw_numbers/prize/json [post]
func (h *Handler) POSTJSONDrawNumbersPrize(ctx *gin.Context) {
}

// @@@Summary 新增鑑定師遊戲資料(json)
// @@@Tags Monopoly
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.NewGameModel true "Monopoly Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/monopoly/json [post]
func (h *Handler) POSTJSONMonopoly(ctx *gin.Context) {
}

// @@@Summary 新增鑑定師獎品資料(json)
// @@@Tags Monopoly Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.NewPrizeModel true "Monopoly Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/monopoly/prize/json [post]
func (h *Handler) POSTJSONMonopolyPrize(ctx *gin.Context) {
}

// @@@Summary 新增快問快答遊戲資料(json)
// @@@Tags QA
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.NewGameModel true "QA Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/QA/json [post]
func (h *Handler) POSTJSONQA(ctx *gin.Context) {
}

// @@@Summary 新增快問快答獎品資料(json)
// @@@Tags QA Prize
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.NewPrizeModel true "QA Prize Parameters"
// @@@Success 200 {array} response.Response
// @@@Failure 500 {array} response.ResponseInternalServerError
// @@@Router /interact/game/QA/prize/json [post]
func (h *Handler) POSTJSONQAPrize(ctx *gin.Context) {
}

// @@@Summary 新增黑名單人員資料(json)
// @@@Tags Black Staff
// @@@version 1.0
// @@@Accept  json
// @@@param body body models.NewBlackStaffModel true "Black Staff Parameters"
// @@@Success 200 {array} response.

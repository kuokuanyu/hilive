package controller

// func (h *Handler) Test(ctx *gin.Context) {
// 	// OGhUfh91dEBHZvVVM1Xw
// 	activityID := "H1le0fBZvxuKZ3Wt3Hva"
// 	log.Println("加入10000筆自定義報名簽到用戶，狀態為sign")

// 	var wg sync.WaitGroup
// 	wg.Add(10000) //計數器
// 	for i := 1; i <= 10000; i++ {
// 		go func(i int) {
// 			defer wg.Done()

// 			// 將資料寫入line_users表中
// 			db.Table(config.LINE_USERS_TABLE).Insert(command.Value{
// 				"user_id":   activityID + "_" + utils.UUID(16),
// 				"name":      "...",
// 				"avatar":    config.HTTP_HILIVES_NET_URL + "/admin/uploads/system/img-user-pic.png",
// 				"phone":     "",
// 				"email":     "",
// 				"ext_email": "",
// 				"identify":  utils.UUID(32),
// 				"friend":    "no",
// 				"line":      "晶橙",
// 				"device":    "customize",

// 				"activity_id":  activityID,
// 				"ext_password": "",
// 			})

// 			models.DefaultApplysignModel().SetDbConn(h.dbConn).Add(
// 				false, models.NewApplysignModel{
// 					UserID:     strconv.Itoa(i),
// 					ActivityID: activityID,
// 					Status:     "sign",
// 				})

// 		}(i)
// 	}

// 	h.redisConn.SetCache("unlucky", 0)
// 	h.redisConn.DelCache("bingo_user_" + "n89wB64b1MANPqaDKnws")
// }

// func (h *Handler) Test(ctx *gin.Context) {
// 	var (
// 		wsConn, conn, _ = NewWebsocketConn(ctx)
// 		// a               = ctx.Query("a")
// 	)
// 	defer wsConn.Close()
// 	defer conn.Close()

// 	log.Println("測試中")

// 	// if a == "1" {
// 	// 用來控制 Redis 訂閱開關判斷
// 	context, cancel := context.WithCancel(ctx.Request.Context())
// 	defer cancel() // 確保函式退出時取消訂閱

// 	// 启动 Redis 订阅者
// 	go h.redisConn.Subscribe(context, "channel:notifications", func(channel, message string) {
// 		log.Println("訂閱者收到資料變動的訊息")
// 	})
// 	// }

// 	for {
// 		message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("取消訂閱")

// 			conn.Close()
// 			h.redisConn.Unsubscribe("channel:notifications")
// 			return
// 		}

// 		// log.Println("執行Publish: ", string(message))
// 		if string(message) != "" {
// 			// log.Println("執行Publish: ", string(message))

// 			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 			err = h.redisConn.Publish("channel:notifications", "觸發")
// 			if err != nil {
// 				log.Println("錯誤: 無法發布消息到 Redis 頻道, ", err)
// 			}

// 			// log.Println("修改redis中的資料")
// 			// h.redisConn.SetCache("test", "test.....")
// 		} else {
// 			// log.Println("心跳")
// 		}
// 	}
// }

// func (h *Handler) Test2(ctx *gin.Context) {
// 	activityID := "OGhUfh91dEBHZvVVM1Xw"
// 	log.Println("報名狀態修改為apply，發送訊息")

// 	// 發郵件
// 	for i := 0; i < MaxRetries; i++ {
// 		// 上鎖
// 		ok, _ := h.acquireLock("mail_lock_"+activityID, LockExpiration)
// 		if ok == "OK" {
// 			// 釋放鎖
// 			defer h.releaseLock("mail_lock_" + activityID)

// 			// 遞減郵件數量
// 			models.DefaultActivityModel().SetDbConn(h.dbConn).
// 				SetRedisConn(h.redisConn).DecrMailAmount(true, activityID)

// 			// 釋放鎖
// 			// releaseLock(s.RedisConn, config.MAIL_LOCK_REDIS+model.ActivityID)
// 			break
// 		}

// 		// 鎖被佔用，稍微延遲後重試
// 		time.Sleep(RetryDelay + time.Duration(rand.Intn(100))*time.Millisecond) // 增加隨機延遲防止競爭
// 	}

// 	// 發簡訊
// 	for i := 0; i < MaxRetries; i++ {
// 		// 上鎖
// 		ok, _ := h.acquireLock("message_lock_"+activityID, LockExpiration)
// 		if ok == "OK" {
// 			// 釋放鎖
// 			defer h.releaseLock("message_lock_" + activityID)

// 			// 遞減簡訊數量
// 			models.DefaultActivityModel().SetDbConn(h.dbConn).
// 				SetRedisConn(h.redisConn).DecrMessageAmount(true, activityID)

// 			// 釋放鎖
// 			// releaseLock(s.RedisConn, config.MAIL_LOCK_REDIS+model.ActivityID)
// 			break
// 		}

// 		// 鎖被佔用，稍微延遲後重試
// 		time.Sleep(RetryDelay + time.Duration(rand.Intn(100))*time.Millisecond) // 增加隨機延遲防止競爭
// 	}
// }

// func (h *Handler) Test3(ctx *gin.Context) {
// 	qrcode := ctx.Query("qrcode")
// 	client := &http.Client{Timeout: 20 * time.Second}

// 	log.Println("執行報名簽到處理api，狀態修改為sign")

// 	req, _ := http.NewRequest("GET", fmt.Sprintf("https://dev.hilives.net/applysign?qrcode=%s", qrcode), nil)

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		// fmt.Printf("Worker %d: 請求失敗: %v\n", workerID, err)
// 		return
// 	}
// 	defer resp.Body.Close()

// }

// func (h *Handler) Test4(ctx *gin.Context) {
// 	// log.Println("計算剩餘獎品數")
// 	// amount := 10000

// 	// 計算剩餘獎品數
// 	for i := 0; i < MaxRetries; i++ {
// 		// 上鎖
// 		ok, _ := h.acquireLock("prize_lock_"+"tVk346broSMKMQ6w1fqi", LockExpiration)
// 		if ok == "OK" {
// 			// 釋放鎖
// 			defer h.releaseLock("prize_lock_" + "tVk346broSMKMQ6w1fqi")
// 			// lock.Lock() //佔有資源
// 			h.redisConn.DecrCache("prize")
// 			// amount--
// 			// lock.Unlock() //釋放資源

// 			// 釋放鎖
// 			// releaseLock(s.RedisConn, config.PRIZE_LOCK_REDIS+model.ActivityID)
// 			break
// 		}

// 		// 鎖被佔用，稍微延遲後重試
// 		time.Sleep(RetryDelay + time.Duration(rand.Intn(100))*time.Millisecond) // 增加隨機延遲防止競爭
// 	}

// 	amount, _ := h.redisConn.GetCache("prize")
// 	log.Println("剩餘獎品數: ", amount)
// }

// func (h *Handler) Test(ctx *gin.Context) {
// 	qrcode := fmt.Sprintf(config.HTTPS_APPLYSIGN_QRCODE_URL, config.HILIVES_NET_URL, fmt.Sprintf("activity_id=%s.user_id=%s", "Bfon6SaV6ORhmuDQUioI", "Bfon6SaV6ORhmuDQUioI_yKPwupTQrSS8L9ky"))
// 	// 设置短信参数
// 	params := &openapi.CreateMessageParams{}
// 	params.SetTo("+886932530813")
// 	params.SetFrom(config.PHONE) // 发送者的 Twilio 电话号码
// 	params.SetBody("XXX 您好: 您已報名 XXX 活動，可利用以下連結進行活動簽到(避免資料洩漏，連結勿提供給他人使用): XXX，以下QRcode應用於主持人掃瞄用戶條碼進行活動簽到判斷(避免資料洩漏，連結勿提供給他人使用)： %s")

// 	// 发送短信
// 	_, err := client.Api.CreateMessage(params)
// 	if err != nil {
// 		log.Println("錯誤: 發送簡訊發生問題")
// 	}

// 	return nil
// }

// func (h *Handler) Test(ctx *gin.Context) {
// 	tokenString := ctx.GetHeader("Authorization")

// 	// 驗證
// 	claims := &Claims{}
// 	token, _ := jwt.ParseWithClaims(tokenString,
// 		claims, func(token *jwt.Token) (interface{}, error) {
// 			return jwtKey, nil
// 		})
// 	// log.Println("token.Raw: ", token.Raw)
// 	// log.Println("token.Signature: ", token.Signature)
// 	log.Println("token.Valid: ", token.Valid)
// 	log.Println("claims.Username: ", claims.Username)

// 	h.executeHTML(ctx, "./hilives/hilive/views/test.html", executeParam{})
// }

// @Summary 測試用
// @Tags Test
// @version 1.0
// @Accept  mpfd
// @@@param activity_id query string true "活動ID"
// @@@param user_id query string true "用戶ID"
// @Success 200 {array} response.Response
// @Failure 500 {array} response.ResponseInternalServerError
// @Router /test [get]
// func (h *Handler) Test(ctx *gin.Context) {

// 	// 示例的 Base64 编码字符串
// 	base64Image := ""

// 	// 去掉前缀
// 	base64Image = strings.TrimPrefix(base64Image, "data:image/png;base64,")

// 	// 解码 Base64 数据
// 	imageData, err := base64.StdEncoding.DecodeString(base64Image)
// 	if err != nil {
// 		log.Println("Base64 解码失败: %v", err)
// 	}

// 	// 在遠端建立新檔案(path為遠端的檔案路徑)
// 	ff, err := os.Create("./hilives/hilive/uploads/test.png")
// 	if err != nil {
// 		log.Println("創建檔案失敗")
// 		return
// 	}

// 	defer func() {
// 		if err2 := ff.Close(); err2 != nil {
// 			err = err2
// 		}
// 	}()

// 	ff.Write(imageData)
// 	// 将解码后的数据写入 PNG 文件
// 	_, err = ff.Write(imageData)
// 	if err != nil {
// 		log.Println("寫入文件失败: %v", err)
// 	}

// 	log.Println("PNG圖片成功生成")
// }

// func (h *Handler) Test2(ctx *gin.Context) {
// 	var (
// 		logger = logrus.New()
// 	)

// 	// 設置logger參數
// 	logger.SetLevel(logrus.DebugLevel)
// 	logger.SetFormatter(&logrus.TextFormatter{
// 		TimestampFormat: "2006-01-02 15:04:05",
// 	})

// 	os.Create("/opt/logs/error.log")

// 	// 開啟檔案
// 	src, err := os.OpenFile("/opt/logs/error.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
// 	if err != nil {
// 		log.Println("錯誤: 開啟檔案發生問題")
// 	}

// 	logger.Out = src

// 	logger.Infof("測試寫入虛擬VM的error.log")

// 	h.executeHTML(ctx, "./hilives/hilive/views/test/2.html", executeParam{})
// }

// func (h *Handler) Test3(ctx *gin.Context) {
// 	if err := os.MkdirAll("/opt/hilives/hilive/test123", 0777); err != nil {
// 		log.Println("錯誤: 建立test123的資料夾發生問題")
// 		log.Println(err)
// 	}

// 	// _, err := os.Create("/opt/hilives/hilive/test.txt")
// 	// if err != nil {
// 	// 	log.Println("無法建立txt檔")
// 	// } else {
// 	// 	log.Println("建立txt檔時沒有錯誤")
// 	// }

// 	h.executeHTML(ctx, "./hilives/hilive/views/test/3.html", executeParam{})
// }

// func (h *Handler) Test4(ctx *gin.Context) {
// 	// 在遠端建立新檔案(path為遠端的檔案路徑)
// 	_, err := os.Create("/opt/hilives/hilive/uploads/123.png")
// 	if err != nil {
// 		log.Println("無法建立圖片檔")
// 	} else {
// 		log.Println("建立圖片檔時沒有錯誤")
// 	}

// 	h.executeHTML(ctx, "./hilives/hilive/views/test/4.html", executeParam{})
// }

// applysignLock.Lock()
// defer applysignLock.Unlock()
// var (
// 	activityID = "Bfon6SaV6ORhmuDQUioI"
// 	userID     = "test"
// )

// 新增報名簽到人員資料
// _, err := models.DefaultApplysignModel().SetDbConn(h.dbConn).
// 	Add(
// 		models.NewApplysignModel{
// 			UserID: userID, ActivityID: activityID, Status: "apply",
// 		})
// if err != nil {
// 	response.Error(ctx, err.Error())
// 	return
// }

// 遞增活動人數與活動號碼
// if err = models.DefaultActivityModel().SetDbConn(h.dbConn).
// 	SetRedisConn(h.redisConn).IncrAttend(true, activityID,
// 	userID); err != nil {
// 	response.Error(ctx, err.Error())
//

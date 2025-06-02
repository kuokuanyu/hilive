package controller

// func (h *Handler) SMTP(ctx *gin.Context) {
// 	// 发件人邮箱地址和密码
// 	from := "sales@hilives.net"
// 	password := "kgdl nacn ltkd ycir"

// 	// SMTP 服务器配置
// 	smtpHost := "smtp.gmail.com"
// 	smtpPort := "587"

// 	// 收件人地址
// 	to := "a167829435@gmail.com"
// 	// to := []string{"a167829435@gmail.com"}

// 	// 生成 QR 码
// 	// qrCodeURL := "https://dev.hilives.net/activity/isexist?activity_id=Bfon6SaV6ORhmuDQUioI&sign=open&is_qrocde=true"
// 	// var pngData bytes.Buffer
// 	// qrCode, err := qrcode.New(qrCodeURL, qrcode.Medium)
// 	// if err != nil {
// 	// 	fmt.Println("Error generating QR code:", err)
// 	// 	return
// 	// }
// 	// err = png.Encode(&pngData, qrCode.Image(256))
// 	// if err != nil {
// 	// 	fmt.Println("Error encoding QR code:", err)
// 	// 	return
// 	// }

// 	// 将 QR 码转换为 base64
// 	// encoded := base64.StdEncoding.EncodeToString(pngData.Bytes())

// 	url := "https://dev.hilives.net/applysign?activity_id=Bfon6SaV6ORhmuDQUioI&user_id=U3140fd73cfd35dd992668ab3b6efdae9&sign=open"
// 	// qrCodeURL := url + "&is_qrcode=true"
// 	// qrCodeURL := "https://api.qrserver.com/v1/create-qr-code/?data=https://dev.hilives.net/activity/isexist?activity_id=Bfon6SaV6ORhmuDQUioI&sign=open&is_qrocde=true"

// 	// 生成 QR 码的链接并进行 URL 编码
// 	qrCodeURL := "https://dev.hilives.net/applysign?activity_id=Bfon6SaV6ORhmuDQUioI.user_id=U3140fd73cfd35dd992668ab3b6efdae9.sign=open.is_qrocde=true"

// 	subject := "Subject: 8 活動報名簽到訊息\r\n"
// 	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
// 	body := fmt.Sprintf(`
//         <html>
// 		 	<body>
// 				<p>XXX 您好，</p>
//         		<p>感謝您報名參加我們的 XXX 活動！以下是您的簽到信息：</p>
// 				<p>該活動的驗證碼為gjkdls，可透過驗證碼進行簽到，</p>
//         		<p>或利用以下連結以進行簽到(避免資料洩漏，連結勿提供給他人使用)：</p>
//         		<p><a href="%s">XXX活動簽到連結</a></p>
//         		<p>以下QRcode應用於主持人掃瞄用戶條碼進行活動簽到處理(避免資料洩漏，連結勿提供給他人使用)：</p>
// 				<img src="https://api.qrserver.com/v1/create-qr-code/?size=256x256&data=%s" alt="签到QR码"/>
//         		<p>如有任何問題，請隨時與我們聯繫。</p>
//         		<p>祝您一切順利！</p>
//         		<p>此致，</p>
//         		<p>活動團隊Hilives</p>
//             </body>
//         </html>`, url, qrCodeURL)

// 	// log.Println("body: ", body)
// 	message := []byte(subject + mime + body)

// 	// 设置认证信息
// 	auth := smtp.PlainAuth("", from, password, smtpHost)

// 	// 发送邮件
// 	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
// 	if err != nil {
// 		fmt.Println("錯誤:", err)
// 		return
// 	}
// 	fmt.Println("邮件发送成功!")
// }

// // getClient 從檔案或web取得token值
// func getClient(config *oauth2.Config) *http.Client {
// 	// Load saved token from a file or get it from the web
// 	tokenFile := "./hilives/hilive/token.json"
// 	token, err := tokenFromFile(tokenFile)
// 	if err != nil {
// 		token = getTokenFromWeb(config)
// 		saveToken(tokenFile, token)
// 	}
// 	return config.Client(context.Background(), token)
// }

// // tokenFromFile 從檔案取得token值
// func tokenFromFile(file string) (*oauth2.Token, error) {
// 	// Read the token from a file
// 	f, err := os.ReadFile(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	tok := &oauth2.Token{}
// 	err = json.Unmarshal(f, tok)
// 	return tok, err
// }

// // getTokenFromWeb 從web取得token值
// func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
// 	// Get a token from the web and return the retrieved token
// 	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
// 	fmt.Printf("Go to the following link in your browser and retrieve the authorization code: \n%v\n", authURL)

// 	var authCode string
// 	fmt.Print("Enter the authorization code: ")
// 	if _, err := fmt.Scan(&authCode); err != nil {
// 		log.Println("Unable to read authorization code: %v", err)
// 	}

// 	tok, err := config.Exchange(context.Background(), authCode)
// 	if err != nil {
// 		log.Println("Unable to retrieve token from web: %v", err)
// 	}
// 	return tok
// }

// // saveToken 取得token值後儲存
// func saveToken(path string, token *oauth2.Token) {
// 	// Save the token to a file
// 	fmt.Printf("Saving credential file to: %s\n", path)
// 	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
// 	if err != nil {
// 		log.Println("Unable to cache oauth token: %v", err)
// 	}
// 	defer f.Close()
// 	json.NewEncoder(f).Encode(token)
// }

// func (h *Handler) SendMail(ctx *gin.Context) {
// 	// Read credentials.json file
// 	b, err := os.ReadFile("./hilives/hilive/credentials.json")
// 	if err != nil {
// 		log.Println("Unable to read client secret file: %v", err)
// 	}

// 	// Create a config from the credentials file
// 	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
// 	if err != nil {
// 		log.Println("Unable to parse client secret file to config: %v", err)
// 	}

// 	client := getClient(config)

// 	// Create a new Gmail service
// 	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
// 	if err != nil {
// 		log.Println("Unable to retrieve Gmail client: %v", err)
// 	}

// 	// Email content
// 	qrCodeURL := "https://dev.hilives.net/activity/isexist?activity_id=Bfon6SaV6ORhmuDQUioI&sign=open"
// 	emailTo := "a167829435@gmail.com"
// 	emailSubject := "参与活动并使用二维码签到"
// 	emailBody := fmt.Sprintf(`
// 	To: %s
// 	Subject: %s
// 	Content-Type: text/html; charset="UTF-8"

// 	<p>感谢您报名参加活动！请使用以下二维码进行签到：</p>
// 	<p><a href="%s"><img src="https://chart.googleapis.com/chart?chs=300x300&cht=qr&chl=%s&choe=UTF-8" alt="二维码"></a></p>
// 	<p>或直接点击链接: <a href="%s">%s</a></p>
// `, emailTo, emailSubject, qrCodeURL, qrCodeURL, qrCodeURL, qrCodeURL)

// 	// Base64 encoding of email content
// 	var message gmail.Message
// 	message.Raw = base64.URLEncoding.EncodeToString([]byte(emailBody))

// 	// Send the email
// 	user := "me" // "me" refers to the authorized user
// 	_, err = srv.Users.Messages.Send(user, &message).Do()
// 	if err != nil {
// 		log.Println("Unable to send email: %v", err)
// 	}
// 	fmt.Println("邮件发送成功")
// }

// var (
// 	googleOauthConfig = &oauth2.Config{
// 		ClientID:     "804432621213-vcs0h63r2uslcn9jeuutgsjokorvjl3k.apps.googleusercontent.com",
// 		ClientSecret: "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",
// 		RedirectURL:  "https://dev.hilives.net/gmail/callback",
// 		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
// 			"https://www.googleapis.com/auth/userinfo.email"},
// 		// "https://www.googleapis.com/auth/userinfo.profile" 允许访问用户的基本个人信息，如姓名、头像和公开的 Google 资料信息。
// 		// "https://www.googleapis.com/auth/userinfo.email" 查看用户的电子邮件地址
// 		// "https://www.googleapis.com/auth/gmail.readonly" 允许应用程序以只读方式访问用户的 Gmail 数据。
// 		Endpoint: google.Endpoint,
// 	}

// 	// oauthStateString = "random" // 用于防止CSRF攻击
// )

// func (h *Handler) HandleGoogleMain(ctx *gin.Context) {
// 	var html = `<html><body><a href="/gmail/login">Google 登录</a></body></html>`
// 	fmt.Fprint(ctx.Writer, html)
// }

// func (h *Handler) HandleGoogleLogin(ctx *gin.Context) {

// 	// access_type=online: 适合需要一次性访问用户数据的情况，不需要长期访问。
// 	// access_type=offline: 适合需要长期访问用户数据的情况，应用程序可以在用户不在线时继续访问数据。
// 	url := googleOauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)

// 	http.Redirect(ctx.Writer, ctx.Request, url, http.StatusTemporaryRedirect)
// }

// func (h *Handler) HandleGoogleCallback(ctx *gin.Context) {
// 	state := ctx.Request.FormValue("state")
// 	if state != oauthStateString {
// 		log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
// 		http.Redirect(ctx.Writer, ctx.Request, "/gmail", http.StatusTemporaryRedirect)
// 		return
// 	}

// 	code := ctx.Request.FormValue("code")
// 	// 使用获得的 token 访问用户信息
// 	token, err := googleOauthConfig.Exchange(context.Background(), code)
// 	if err != nil {
// 		log.Printf("code exchange failed: %s\n", err.Error())
// 		http.Redirect(ctx.Writer, ctx.Request, "/gmail", http.StatusTemporaryRedirect)
// 		return
// 	}

// 	// 使用获得的 token 访问用户信息
// 	client := googleOauthConfig.Client(context.Background(), token)
// 	// 获取用户个人信息
// 	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
// 	if err != nil {
// 		log.Println("Unable to retrieve user info: %v", err)
// 		http.Redirect(ctx.Writer, ctx.Request, "/gmail", http.StatusTemporaryRedirect)
// 	}
// 	defer userInfoResp.Body.Close()

// 	userInfo, err := ioutil.ReadAll(userInfoResp.Body)
// 	if err != nil {
// 		log.Println("Unable to read user info response: %v", err)
// 		http.Redirect(ctx.Writer, ctx.Request, "/gmail", http.StatusTemporaryRedirect)
// 	}

// 	fmt.Fprintf(ctx.Writer, "User Info: %s", userInfo)

// }

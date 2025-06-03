package main

import (
	"hilive/engine"
	"hilive/modules/config"

	_ "hilive/docs"

	_ "github.com/go-sql-driver/mysql" // mysql引擎
)
777777777

6

5

4

// @title Hilives API 文檔
// @version 1.0
// @description Hilives 平台

// @contact.name Hilives
// @contact.url

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host apidev.hilives.net
// @BasePath /v1
// @schemes https
func main() {
	// fmt.Println("CPU: ", runtime.NumCPU())
	// runtime.GOMAXPROCS(4)

	// 可以取得內網ip
	// addrs, _ := net.InterfaceAddrs()
	// for _, add := range addrs {
	// 	fmt.Println("add.String(): ", add.String())
	// 	fmt.Println("add.Network(): ", add.Network())
	// }

	// 	host, _ := os.Hostname()
	// addrs, _ := net.LookupIP(host)
	// for _, addr := range addrs {
	//     if ipv4 := addr.To4(); ipv4 != nil {
	//         fmt.Println("IPv4: ", ipv4)
	//     }
	// }

	// t := engine.DefaultEngine().InitRouter()
	// fmt.Println("RouterGroup: ", t.Gin.RouterGroup)
	// fmt.Println("RemoteIPHeaders: ", t.Gin.RemoteIPHeaders)
	// fmt.Println("TrustedPlatform: ", t.Gin.TrustedPlatform)

	cfg := config.Config{
		Databases: config.DatabaseList{
			config.MYSQL_ENGINE: {
				Host:       config.MYSQL_HOST,
				Port:       config.MYSQL_PORT,
				User:       config.MYSQL_USER,
				Pwd:        config.MYSQL_PASSWORD,
				Name:       config.MYSQL_NAME,
				MaxIdleCon: config.MYSQL_MAXIDLECON,
				MaxOpenCon: config.MYSQL_MAXOPENCON,
				Driver:     config.MYSQL_DRIVER,
			},
			// "heroku": {
			// 	Host:       "us-cdbr-east-02.cleardb.com",
			// 	Port:       "3306",
			// 	User:       "be94ad46dfd2c5",
			// 	Pwd:        "0986ac8c",
			// 	Name:       "heroku_340b0d6567ec671",
			// 	MaxIdleCon: 50,
			// 	MaxOpenCon: 150,
			// 	Driver:     "mysql",
			// },
		},
		RedisList: config.RedisList{
			config.REDIS_ENGINE: {
				Host: config.REDIS_HOST,
				Port: config.REDIS_PORT,
			},
		},
		MongoList: config.MongoList{
			config.MONGO_ENGINE: {
				Host: config.MONGO_HOST,
				Port: config.MONGO_PORT,
				User: config.MONGO_USER,
				Pwd:  config.MONGO_PASSWORD,
				Name: config.MONGO_NAME,
			},
		},
		Prefix: config.PREFIX,
		Store: config.Store{
			Path:   config.STORE_PATH,
			Prefix: config.STORE_PREFIX,
		},
	}

	eng := engine.DefaultEngine().
		InitDatabase(cfg).InitMongo(cfg).InitRedis(cfg).
		SetEngine().InitRouter()

	// http連線
	eng.Gin.Run(":80")

	// eng.Gin.Run(":443")

	// 第一種寫法(ssl憑證)-----start
	// listener := autocert.NewListener(
	// 	config.HILIVES_NET_URL,
	// 	config.API_URL,
	// )

	// // listener = netutil.LimitListener(listener, 100000) // 連接數(錯誤)

	// log.Fatal(http.Serve(listener, eng.Gin))
	// 第一種寫法(ssl憑證)-----end

	// 第二種寫法(ssl憑證)-----start
	// certManager := autocert.Manager{
	// 	Prompt: autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist(config.HILIVES_NET_URL,
	// 		config.WWW_HILIVES_NET_URL, config.API_URL), // 網域設置
	// 	Cache: autocert.DirCache("/root/.cache/golang-autocert"), // ssl憑證存放位置
	// }

	// // http package server端參數設置
	// server := &http.Server{
	// 	// ReadTimeout:  5 * time.Second,
	// 	// WriteTimeout: 5 * time.Second,
	// 	// IdleTimeout:  60 * time.Second,
	// 	Addr: ":https", // :443
	// 	TLSConfig: &tls.Config{
	// 		GetCertificate: certManager.GetCertificate},
	// 	Handler: eng.Gin,
	// }

	// server.ListenAndServeTLS("", "")

	// http package預設的Transport
	// MaxIdleConns  所有host的連接池最大連接數(預設100)
	// MaxIdleConnsPerHost  每個host最大的空閒連接數(預設2)。
	// MaxConnsPerHost 每個host最大的連接數，0表示不限制
	// var DefaultTransport RoundTripper = &Transport{
	// 	Proxy: ProxyFromEnvironment,
	// 	DialContext: (&net.Dialer{
	// 	   Timeout:   30 * time.Second,
	// 	   KeepAlive: 30 * time.Second,
	// 	   DualStack: true,
	// 	}).DialContext,
	// 	ForceAttemptHTTP2:     true,
	// 	MaxIdleConns:          100,
	// 	IdleConnTimeout:       90 * time.Second,
	// 	TLSHandshakeTimeout:   10 * time.Second,
	// 	ExpectContinueTimeout: 1 * time.Second,
	//  }

	// http package client端參數設置
	// _ = http.Client{
	// 	Transport: &http.Transport{
	// 		MaxIdleConns:          1000000,
	// 		MaxIdleConnsPerHost:   1000000,
	// 		MaxConnsPerHost:       1000000,
	// 		IdleConnTimeout:       60 * time.Second,
	// 		TLSHandshakeTimeout:   10 * time.Second,
	// 		ExpectContinueTimeout: 1 * time.Second,
	// 	},
	// }
	// 第二種寫法(ssl憑證)-----end

	// 第三種寫法(測試LimitListener最大連接數中)-----start
	// certManager := autocert.Manager{
	// 	Prompt: autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist(config.HILIVES_NET_URL,
	// 		config.WWW_HILIVES_NET_URL, config.API_URL), // 網域設置
	// 	Cache: autocert.DirCache("/root/.cache/golang-autocert"), // ssl憑證存放位置
	// }

	// // http package server端參數設置
	// server := &http.Server{
	// 	// ReadTimeout:  5 * time.Second,
	// 	// WriteTimeout: 5 * time.Second,
	// 	// IdleTimeout:  60 * time.Second,
	// 	Addr: ":443",
	// 	TLSConfig: &tls.Config{
	// 		GetCertificate: certManager.GetCertificate},
	// 	Handler: eng.Gin,
	// }

	// listener, err := net.Listen("tcp", server.Addr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer listener.Close()
	// listener = netutil.LimitListener(listener, 1000000) // 連接數

	// log.Fatal(server.Serve(tls.NewListener(listener, server.TLSConfig)))
	// 第三種寫法(測試LimitListener最大連接數中)-----end
}

// func mainPage(c *fiber.Ctx) error {
// 	fmt.Println("有?")
// 	return c.Render("index", fiber.Map{
// 		"people": 100,
// 		"People": 10,
// 	})
// }

// svelte測試-----start
// template render engine
// engine := html.New("./test/public", ".html")

// app := fiber.New(fiber.Config{
// 	Views: engine, //set as render engine
// })

// app.Static("/test", "./test")
// app.Get("/test", mainPage)
// // app.Listen(":3000", fiber.ListenConfig{})

// // eng.Gin = app

// // listener := autocert.NewListener(
// // 	"dev.hilives.net")

// // eng.Gin.Listen(":3000")
// log.Fatal(app.Listen(":80"))

// svelte測試-----end

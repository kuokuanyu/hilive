package file

import (
	"hilive/modules/utils"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"sync"
)

type Uploader interface {
	Upload(m *multipart.Form, path, userID, activityID, gameID, game string) error
}

type UploaderGenerator func() Uploader

type UploadFun func(*multipart.FileHeader, string) (string, error)

var uploaderList = map[string]UploaderGenerator{
	"hilives": GetFileUploader,
}

// GetFileEngine 透過參數取得Uploader
func GetFileEngine(name string) Uploader {
	if up, ok := uploaderList[name]; ok {
		return up()
	}
	panic("錯誤的檔案引擎名稱")
}

// Upload 儲存上傳的檔案
func Upload(c UploadFun, form *multipart.Form, apiPath string) error {
	var (
		wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
		mu sync.Mutex     // 互斥鎖
	)

	// 多線呈更新圖片資訊
	// wg.Add(len(form.File)) //計數器

	for k := range form.File {
		// log.Println("k(表單欄位名稱): ", k, len(form.File[k]))
		for _, fileObj := range form.File[k] {
			wg.Add(1) // 每個文件都增加計數

			go func(k string, fileObj *multipart.FileHeader) {
				defer wg.Done()
				// log.Println("檔案名稱(用戶原本的檔案名稱): ", fileObj.Filename)

				// 在 Goroutine 內部定義變數，避免競爭條件
				suffix := path.Ext(fileObj.Filename) // 檔案類型
				// fmt.Println("suffix(檔案類型): ", suffix)
				var filename string

				// 存放在遠端上的檔名
				if apiPath == "/v1/admin/manager" || apiPath == "/v1/user" {
					// 用戶頁面
					filename = "user" + suffix
				} else if apiPath == "/v1/line_user" {
					// 報名簽到人員更新姓名頭像頁面
					filename = utils.UUID(8) + suffix
				} else if apiPath == "/v1/applysign/apply" {
					// 報名頁面
					filename = k + suffix
				} else if apiPath == "/v1/applysign/customize" {
					// 自定義欄位設置頁面
					filename = k + suffix
				} else if apiPath == "/v1/applysign/qrcode" {
					// QRcode自定義欄位頁面
					filename = k + suffix
				} else if apiPath == "/v1/interact/wall/message" {
					// 訊息區頁面
					filename = k + suffix
				} else if apiPath == "/v1/interact/wall/topic" {
					// 主題區頁面
					filename = k + suffix
				} else if apiPath == "/v1/interact/wall/question" {
					// 提問區頁面
					filename = k + suffix
				} else if apiPath == "/v1/interact/wall/danmu" {
					// 一般彈幕頁面
					filename = k + suffix
				} else if apiPath == "/v1/interact/sign/general" {
					// 一般簽到頁面
					filename = k + suffix
				} else if apiPath == "/v1/interact/sign/threed" {
					// 3D簽到頁面
					filename = k + suffix
				} else if apiPath == "/v1/interact/sign/signname" {
					// 簽名牆頁面
					filename = k + suffix
				} else if strings.Contains(apiPath, "/v1/interact/sign") {
					// 遊戲頁面
					filename = k + suffix

					// 獎品頁面隨機命名
					if strings.Contains(apiPath, "/prize/form") ||
						strings.Contains(apiPath, "/vote/option") { // 投票選項
						filename = utils.UUID(8) + suffix
					}
				} else if strings.Contains(apiPath, "/v1/interact/game") {
					// 遊戲頁面
					filename = k + suffix

					// 獎品頁面隨機命名
					if strings.Contains(apiPath, "/prize/form") {
						filename = utils.UUID(8) + suffix
					}
				} else {
					// 隨機命名
					filename = utils.UUID(8) + suffix
				}

				// pathStr: 儲存至遠端的檔案名稱
				pathStr, _ := c(fileObj, filename)
				// if err != nil {
				// 	// log.Println("c錯誤是: ", err)
				// 	return err
				// }

				// log.Println("form.Value(原本是空的): ", form.Value[k])
				// log.Println("pathStr(遠端的檔案名稱): ", pathStr)

				mu.Lock() // 加鎖，防止並發寫入 map
				form.Value[k] = append(form.Value[k], pathStr)
				mu.Unlock() // 解鎖

			}(k, fileObj)
		}
	}
	wg.Wait() //等待計數器歸0

	return nil
}

// SaveMultipartFile 儲存檔案在Store設置的路徑中
func SaveMultipartFile(fh *multipart.FileHeader, path string) error {
	// log.Println("檔案名稱(用戶本地端的檔案名稱): ", fh.Filename)
	// log.Println("fh.Header", fh.Header)
	// log.Println("fh.Size", fh.Size)
	// log.Println("fh.Header.Get(fh.Filename)", fh.Header.Get(fh.Filename))

	// fmt.Println("path(遠端的檔案路徑): ", path)

	// 打開上傳的檔案(用戶本地端的檔案名稱)
	f, err := fh.Open()
	if err != nil {
		return err
	}

	if ff, ok := f.(*os.File); ok {
		// log.Println("ok")
		if err := f.Close(); err != nil {
			return err
		}

		// fmt.Println("ff.Name(): ", ff.Name())
		if os.Rename(ff.Name(), path) == nil {
			return nil
		}

		f, err = fh.Open()
		if err != nil {
			return err
		}
	} else {
		// 都是執行這一段
		// log.Println("原本沒有檔案")
	}

	defer func() {
		if err2 := f.Close(); err2 != nil {
			err = err2
		}
	}()

	// err = os.Remove(path)
	// fmt.Println("錯誤: ", err)

	// 在遠端建立新檔案(path為遠端的檔案路徑)
	ff, err := os.Create(path)
	if err != nil {
		return err
	}
	// fmt.Println("ff.Name(): ", ff.Name())

	defer func() {
		if err2 := ff.Close(); err2 != nil {
			err = err2
		}
	}()

	// 將用戶上傳的檔案資料寫入遠端的新檔案中
	_, err = copyZeroAlloc(ff, f)
	return err
}

// copyZeroAlloc 將用戶上傳的檔案資料寫入遠端的新檔案中
func copyZeroAlloc(w io.Writer, r io.Reader) (int64, error) {
	buf := copyBufPool.Get().([]byte)
	n, err := io.CopyBuffer(w, r, buf)
	copyBufPool.Put(buf)
	return n, err
}

var copyBufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 4096)
	},
}

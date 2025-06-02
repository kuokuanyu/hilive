package paginator

import (
	"html/template"
	"math"
	"strconv"

	"hilive/modules/parameter"
)

// Paginator 分頁資訊
type Paginator struct {
	Total         string                   // 總資料數
	URL           string                   // 第n頁的URL(沒pagesize)
	PageSizeList  []string                 // 單頁顯示資料
	PreviousClass string                   // 如果沒上一頁為disabled，否則為空
	PreviousURL   string                   // 前一頁的url參數，例如在第二頁時回傳第一頁的url參數
	Pages         []map[string]string      // 每一分頁的資訊，包刮page、active、issplit、url...
	NextClass     string                   // 如果沒下一頁為disables，否則為空
	NextURL       string                   // 下頁的url參數
	Option        map[string]template.HTML // 對每一個分頁設置HTML
}

// GetPaginator 取得分頁資訊
func GetPaginatorInformation(size int, params parameter.Parameters) Paginator {
	paginator := Paginator{}

	pageInt, _ := strconv.Atoi(params.Page)
	pageSizeInt, _ := strconv.Atoi(params.PageSize)

	totalPage := int(math.Ceil(float64(size) / float64(pageSizeInt)))

	paginator.URL = params.Path + params.GetRouteWithoutPageSize("1")
	paginator.Total = strconv.Itoa(size)

	if pageInt == 1 {
		paginator.PreviousClass = "disabled"
		paginator.PreviousURL = params.Path
	} else {
		paginator.PreviousClass = ""
		paginator.PreviousURL = params.Path + params.GetLastPageRoute()
	}
	if pageInt == totalPage {
		paginator.NextClass = "disabled"
		paginator.NextURL = params.Path
	} else {
		paginator.NextClass = ""
		paginator.NextURL = params.Path + params.GetNextPageRoute()
	}

	paginator.Pages = []map[string]string{}
	if totalPage < 10 {
		var pagesArr []map[string]string
		for i := 1; i < totalPage+1; i++ {
			if i == pageInt {
				pagesArr = append(pagesArr, map[string]string{
					"page":    params.Page,
					"active":  "active",
					"isSplit": "0",
					"url":     params.URL(params.Page),
				})
			} else {
				page := strconv.Itoa(i)
				pagesArr = append(pagesArr, map[string]string{
					"page":    page,
					"active":  "",
					"isSplit": "0",
					"url":     params.URL(page),
				})
			}
		}
		paginator.Pages = pagesArr
	} else {
		var pagesArr []map[string]string
		if pageInt < 6 {
			for i := 1; i < totalPage+1; i++ {

				if i == pageInt {
					pagesArr = append(pagesArr, map[string]string{
						"page":    params.Page,
						"active":  "active",
						"isSplit": "0",
						"url":     params.URL(params.Page),
					})
				} else {
					page := strconv.Itoa(i)
					pagesArr = append(pagesArr, map[string]string{
						"page":    page,
						"active":  "",
						"isSplit": "0",
						"url":     params.URL(page),
					})
				}

				if i == 6 {
					pagesArr = append(pagesArr, map[string]string{
						"page":    "",
						"active":  "",
						"isSplit": "1",
						"url":     params.URL("6"),
					})
					i = totalPage - 1
				}
			}
		} else if pageInt < totalPage-4 {
			for i := 1; i < totalPage+1; i++ {

				if i == pageInt {
					pagesArr = append(pagesArr, map[string]string{
						"page":    params.Page,
						"active":  "active",
						"isSplit": "0",
						"url":     params.URL(params.Page),
					})
				} else {
					page := strconv.Itoa(i)
					pagesArr = append(pagesArr, map[string]string{
						"page":    page,
						"active":  "",
						"isSplit": "0",
						"url":     params.URL(page),
					})
				}

				if i == 2 {
					pagesArr = append(pagesArr, map[string]string{
						"page":    "",
						"active":  "",
						"isSplit": "1",
						"url":     params.URL("2"),
					})
					if pageInt < 7 {
						i = 5
					} else {
						i = pageInt - 2
					}
				}

				if pageInt < 7 {
					if i == pageInt+5 {
						pagesArr = append(pagesArr, map[string]string{
							"page":    "",
							"active":  "",
							"isSplit": "1",
							"url":     params.URL(strconv.Itoa(i)),
						})
						i = totalPage - 1
					}
				} else {
					if i == pageInt+3 {
						pagesArr = append(pagesArr, map[string]string{
							"page":    "",
							"active":  "",
							"isSplit": "1",
							"url":     params.URL(strconv.Itoa(i)),
						})
						i = totalPage - 1
					}
				}
			}
		} else {
			for i := 1; i < totalPage+1; i++ {

				if i == pageInt {
					pagesArr = append(pagesArr, map[string]string{
						"page":    params.Page,
						"active":  "active",
						"isSplit": "0",
						"url":     params.URL(params.Page),
					})
				} else {
					page := strconv.Itoa(i)
					pagesArr = append(pagesArr, map[string]string{
						"page":    page,
						"active":  "",
						"isSplit": "0",
						"url":     params.URL(page),
					})
				}

				if i == 2 {
					pagesArr = append(pagesArr, map[string]string{
						"page":    "",
						"active":  "",
						"isSplit": "1",
						"url":     params.URL("2"),
					})
					i = totalPage - 4
				}
			}
		}
		paginator.Pages = pagesArr
	}
	return paginator
}

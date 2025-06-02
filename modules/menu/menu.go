package menu

// Item menu資料
// type Item struct {
// 	Name         string // menu title 欄位資料
// 	ID           string // menu ID
// 	URL          string // menu URL
// 	SidebarBtnL  string // SidebarBtnL欄位資料
// 	ChildrenList []Item // 放子選單
// }

// Menu 紀錄所有menu資訊
// type Menu struct {
// 	List []Item
// }

// GetMenu 取得menu資料
// func GetMenu(user models.LoginUser, conn db.Connection) *Menu {
// 	var (
// 		menus []map[string]interface{}
// 	)

// 	// TODO: 目前沒有權限相關問題，先拿掉角色權限菜單判斷
// 	// user.GetRoles().GetMenus()
// 	// if user.IsSuperAdmin() {
// 	menus, _ = db.Conn(conn).Table(config.MENU_TABLE).
// 		Where("id", ">", 5).All()
// 	// } else {
// 	// 	var ids []interface{}
// 	// 	for i := 0; i < len(user.MenuIDs); i++ {
// 	// 		ids = append(ids, user.MenuIDs[i])
// 	// 	}
// 	// 	menus, _ = db.Conn(conn).Table(config.MENU_TABLE).
// 	// 		WhereIn("id", ids).All()
// 	// }
// 	menuList := MapToItem(menus, 0)
// 	return &Menu{
// 		List: menuList,
// 	}
// }

// MapToItem map轉換成[]Item
// func MapToItem(menus []map[string]interface{}, parentID int64) []Item {
// 	items := make([]Item, 0)
// 	for j := 0; j < len(menus); j++ {
// 		if parentID == menus[j]["parent_id"].(int64) {
// 			title := menus[j]["title"].(string)
// 			child := Item{
// 				Name:         title,
// 				ID:           strconv.FormatInt(menus[j]["id"].(int64), 10),
// 				URL:          menus[j]["url"].(string),
// 				SidebarBtnL:  menus[j]["sidebarbtnl"].(string),
// 				ChildrenList: MapToItem(menus, menus[j]["id"].(int64)),
// 			}
// 			items = append(items, child)
// 		}
// 	}
// 	return items
// }

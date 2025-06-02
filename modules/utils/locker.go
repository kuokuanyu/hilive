package utils

// // LockLucky 上鎖，抽獎時需要使用，避免用戶併發多次抽獎
// func LockLucky(id string) bool {
// 	return lockLuckyServ(id)
// }

// // UnlockLucky 解鎖，抽獎時需要使用，避免用戶併發多次抽獎
// func UnlockLucky(id string) bool {
// 	return unlockLuckyServ(id)
// }

// // getLuckyLockKey 取得緩存名稱
// func getLuckyLockKey(id string) string {
// 	return fmt.Sprintf("lucky_lock_%s", id)
// }

// // lockLuckyServ 上鎖
// func lockLuckyServ(id string) bool {
// 	key := getLuckyLockKey(id)
// 	cacheObj := datasource.InstanceCache()
// 	rs, _ := cacheObj.Do("SET", key, 1, "EX", 3, "NX")
// 	if rs == "OK" {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// // unlockLuckyServ 解鎖
// func unlockLuckyServ(id string) bool {
// 	key := getLuckyLockKey(id)
// 	cacheObj := datasource.InstanceCache()
// 	rs, _ := cacheObj.Do("DEL", key)
// 	if rs == "OK" {
// 		return true
// 	} else {
// 		return false
// 	}
// }

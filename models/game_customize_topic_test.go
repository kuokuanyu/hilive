package models

import (
	"fmt"
	"log"
	"testing"
)

// 資料表增加遊戲自定義主題欄位
func Test_Customize_Topic(t *testing.T) {
	// ex : ALTER TABLE activity_game_lottery_picture ADD lottery_starrysky_35 varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT '';
	// lottery_遊戲類別_主題_角色_圖片_編號、其他遊戲_主題_角色_圖片_編號

	log.Println("處理資料表")
	var (
		table     = "activity_game_vote_picture" // 資料表，redpack、ropepack、qa、draw_numbers、lottery、monopoly、tugofwar、bingo、whack_mole、vote
		gameUpper = "Vote"                         // 遊戲(大寫)，Redpack、Ropepack、QA、DrawNumbers、Lottery、Monopoly、Tugofwar、Bingo、Whackmole、Vote
		gameLower = "vote"                         // 遊戲(小寫)，redpack、ropepack、qa、draw_numbers、lottery、monopoly、tugofwar、bingo、whackmole、vote

		gameTypeUpper = "" // 遊戲種類，Turntable(輪盤)、Jiugongge(九宮格)
		gameType      = "" // 遊戲種類，turntable(輪盤)、jiugongge(九宮格)

		// 	1. 搖紅包 : 經典主題 (classic)、櫻花主題 (cherry)、christmas
		// 2. 套紅包 : 經典主題 (classic)、兔年主題 (newyear_rabbit)、中秋主題 (moonfestival)、3D(3D)
		// 3. 敲敲樂 : 經典主題 (classic)、萬聖節主題 (halloween)、聖誕節主題 (christmas)
		// 4. 快問快答 : 經典主題 (classic)、電路主題 (electric)、中秋主題 (moonfestival)、龍年(newyear_dragon)
		// 5. 抽獎遊戲 : 經典主題 (classic)、星空主題 (starrysky)
		// 6. 鑑定師 : 經典主題 (classic)、紅包主題(redpack)、兔年主題 (newyear_rabbit)、生魚片主題(sashimi)
		// 7. 搖號抽獎 : 經典主題 (classic) 、 黃金主題 (gold)、龍年(newyear_dragon)、櫻花主題(cherry)
		// 8. 拔河遊戲: 經典主題 (classic)、校園主題 (school)、聖誕節主題 (christmas)
		// 9. 賓果遊戲: 經典主題 (classic)、龍年(newyear_dragon)、櫻花主題(cherry)
		// 10. 投票: 經典主題 (classic)、
		topicUpper = "Classic" // 主題(大寫)
		topicLower = "classic" // 主題(小寫)

		rolesUpper = []string{"H", "G", "C", "H", "G", "C"} // 角色，h、g、c
		roles      = []string{"h", "g", "c", "h", "g", "c"} // 角色，h、g、c

		picStarts = []int{1, 1, 1, 0, 0, 0}   // 開始欄位(靜態主持、靜態玩家、靜態共用、動態主持、動態玩家、動態共用)
		picEnds   = []int{25, 6, 2, 0, 0, 0} // 結束欄位(靜態主持、靜態玩家、靜態共用、動態主持、動態玩家、動態共用)

		startMusic  = false // 是否建立開始音樂欄位
		gamingMusic = false // 是否建立遊戲中音樂欄位
		endMusic    = false // 是否建立結束音樂欄位

		gameQuery     = "ALTER TABLE %s ADD %s_%s_%s_%s_0%d varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT '';"    // 小於十
		gameQuery2    = "ALTER TABLE %s ADD %s_%s_%s_%s_%d varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT '';"     // 大於十
		lotteryQuery  = "ALTER TABLE %s ADD %s_%s_%s_%s_%s_0%d varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT '';" // 小於十
		lotteryQuery2 = "ALTER TABLE %s ADD %s_%s_%s_%s_%s_%d varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT '';"  // 大於十
		musicquery    = "ALTER TABLE %s ADD %s_%s_%s varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' NOT NULL COMMENT '';"           // 音樂
	)

	// 自定義欄位
	for n, start := range picStarts {
		var picType string
		if n <= 2 {
			picType = "pic"
		} else {
			picType = "ani"
		}

		if gameType == "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(gameQuery, table, gameLower, topicLower, roles[n], picType, i))
					}
				}
			}
		} else if gameType != "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(lotteryQuery, table, gameLower, gameType, topicLower, roles[n], picType, i))
					}
				}
			}
		}

		if gameType == "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(gameQuery2, table, gameLower, topicLower, roles[n], picType, i))
				}
			}
		} else if gameType != "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(lotteryQuery2, table, gameLower, gameType, topicLower, roles[n], picType, i))
				}
			}
		}
	}

	// 開始音樂
	if startMusic {
		fmt.Println(fmt.Sprintf(musicquery, table, gameLower, "bgm", "start"))
	}
	// 遊戲中音樂
	if gamingMusic {
		fmt.Println(fmt.Sprintf(musicquery, table, gameLower, "bgm", "gaming"))
	}

	// 結束音樂
	if endMusic {
		fmt.Println(fmt.Sprintf(musicquery, table, gameLower, "bgm", "end"))
	}

	log.Println("處理資料表")

	log.Println("處理model")

	gameQuery = `%s%s%s%s0%d string "json:"%s_%s_%s_%s_0%d" example:"picture""`         // 小於十
	gameQuery2 = `%s%s%s%s%d string "json:"%s_%s_%s_%s_%d" example:"picture""`          // 大於十
	lotteryQuery = `%s%s%s%s%s0%d string "json:"%s_%s_%s_%s_%s_0%d" example:"picture""` // 小於十
	lotteryQuery2 = `%s%s%s%s%s%d string "json:"%s_%s_%s_%s_%s_%d" example:"picture""`  // 大於十
	musicquery = `%s%s%s string "json:"%s_%s_%s" example:"picture""`                    // 音樂

	for n, start := range picStarts {
		var (
			picTypeUpper string
			picType      string
		)
		if n <= 2 {
			picTypeUpper = "Pic"
			picType = "pic"
		} else {
			picTypeUpper = "Ani"
			picType = "ani"
		}

		if gameType == "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(gameQuery,
							gameUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
							gameLower, topicLower, roles[n], picType, i))
					}
				}
			}
		} else if gameType != "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(lotteryQuery,
							gameUpper, gameTypeUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
							gameLower, gameType, topicLower, roles[n], picType, i))
					}
				}
			}
		}

		if gameType == "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(gameQuery2,
						gameUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
						gameLower, topicLower, roles[n], picType, i))
				}
			}
		} else if gameType != "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(lotteryQuery2,
						gameUpper, gameTypeUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
						gameLower, gameType, topicLower, roles[n], picType, i))
				}
			}
		}
	}

	// 開始音樂
	if startMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameUpper, "Bgm", "Start", gameLower, "bgm", "start"))
	}
	// 遊戲中音樂
	if gamingMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameUpper, "Bgm", "Gaming", gameLower, "bgm", "gaming"))
	}

	// 結束音樂
	if endMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameUpper, "Bgm", "End", gameLower, "bgm", "end"))
	}

	log.Println("處理model")

	log.Println("處理Add")

	gameQuery = `"%s_%s_%s_%s_0%d": model.%s%s%s%s0%d,`         // 小於十
	gameQuery2 = `"%s_%s_%s_%s_%d": model.%s%s%s%s%d,`          // 大於十
	lotteryQuery = `"%s_%s_%s_%s_%s_0%d": model.%s%s%s%s%s0%d,` // 小於十
	lotteryQuery2 = `"%s_%s_%s_%s_%s_%d": model.%s%s%s%s%s%d,`  // 大於十
	musicquery = `"%s_%s_%s": model.%s%s%s,`                    // 音樂

	for n, start := range picStarts {
		var (
			picTypeUpper string
			picType      string
		)
		if n <= 2 {
			picTypeUpper = "Pic"
			picType = "pic"
		} else {
			picTypeUpper = "Ani"
			picType = "ani"
		}

		if gameType == "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(gameQuery,
							gameLower, topicLower, roles[n], picType, i,
							gameUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
						))
					}
				}
			}

		} else if gameType != "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(lotteryQuery,
							gameLower, gameType, topicLower, roles[n], picType, i,
							gameUpper, gameTypeUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
						))
					}
				}
			}
		}

		if gameType == "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(gameQuery2,
						gameLower, topicLower, roles[n], picType, i,
						gameUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
					))
				}
			}

		} else if gameType != "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(lotteryQuery2,
						gameLower, gameType, topicLower, roles[n], picType, i,
						gameUpper, gameTypeUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
					))
				}
			}
		}
	}
	// 開始音樂
	if startMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameLower, "bgm", "start", gameUpper, "Bgm", "Start"))
	}
	// 遊戲中音樂
	if gamingMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameLower, "bgm", "gaming", gameUpper, "Bgm", "Gaming"))
	}

	// 結束音樂
	if endMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameLower, "bgm", "end", gameUpper, "Bgm", "End"))
	}

	log.Println("處理Add")

	log.Println("處理Update")

	log.Println("放在fields遊戲的參數裡: ")

	gameQuery = `"%s_%s_%s_%s_0%d",`       // 小於十
	gameQuery2 = `"%s_%s_%s_%s_%d",`       // 大於十
	lotteryQuery = `"%s_%s_%s_%s_%s_0%d",` // 小於十
	lotteryQuery2 = `"%s_%s_%s_%s_%s_%d",` // 大於十
	musicquery = `"%s_%s_%s",`             // 音樂

	for n, start := range picStarts {
		var (
			// picTypeUpper string
			picType string
		)
		if n <= 2 {
			// picTypeUpper = "Pic"
			picType = "pic"
		} else {
			// picTypeUpper = "Ani"
			picType = "ani"
		}

		if gameType == "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(gameQuery,
							gameLower, topicLower, roles[n], picType, i,
						))
					}
				}
			}
		} else if gameType != "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(lotteryQuery,
							gameLower, gameType, topicLower, roles[n], picType, i,
						))
					}
				}
			}
		}

		if gameType == "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(gameQuery2,
						gameLower, topicLower, roles[n], picType, i,
					))
				}
			}
		} else if gameType != "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(lotteryQuery2,
						gameLower, gameType, topicLower, roles[n], picType, i,
					))
				}
			}
		}
	}

	// 開始音樂
	if startMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameLower, "bgm", "start"))
	}
	// 遊戲中音樂
	if gamingMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameLower, "bgm", "gaming"))
	}

	// 結束音樂
	if endMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameLower, "bgm", "end"))
	}

	log.Println("放在values遊戲的參數裡: ")

	gameQuery = `model.%s%s%s%s0%d,`      // 小於十
	gameQuery2 = `model.%s%s%s%s%d,`      // 大於十
	lotteryQuery = `model.%s%s%s%s%s0%d,` // 小於十
	lotteryQuery2 = `model.%s%s%s%s%s%d,` // 大於十
	musicquery = `model.%s%s%s,`          // 音樂

	for n, start := range picStarts {
		var (
			picTypeUpper string
			// picType      string
		)
		if n <= 2 {
			picTypeUpper = "Pic"
			// picType = "pic"
		} else {
			picTypeUpper = "Ani"
			// picType = "ani"
		}

		if gameType == "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(gameQuery,
							gameUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
						))
					}
				}
			}
		} else if gameType != "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(lotteryQuery,
							gameUpper, gameTypeUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
						))
					}
				}
			}
		}

		if gameType == "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(gameQuery2,
						gameUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
					))
				}
			}
		} else if gameType != "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(lotteryQuery2,
						gameUpper, gameTypeUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
					))
				}
			}
		}
	}

	// 開始音樂
	if startMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameUpper, "Bgm", "Start"))
	}
	// 遊戲中音樂
	if gamingMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameUpper, "Bgm", "Gaming"))
	}

	// 結束音樂
	if endMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameUpper, "Bgm", "End"))
	}

	log.Println("處理Update")

	log.Println("處理Find")

	log.Println("Find function裡需要增加的參數: ")

	gameQuery = `"%s.%s_%s_%s_%s_0%d",`       // 小於十
	gameQuery2 = `"%s.%s_%s_%s_%s_%d",`       // 大於十
	lotteryQuery = `"%s.%s_%s_%s_%s_%s_0%d",` // 小於十
	lotteryQuery2 = `"%s.%s_%s_%s_%s_%s_%d",` // 大於十
	musicquery = `"%s.%s_%s_%s",`             // 音樂

	for n, start := range picStarts {
		var (
			// picTypeUpper string
			picType string
		)
		if n <= 2 {
			// picTypeUpper = "Pic"
			picType = "pic"
		} else {
			// picTypeUpper = "Ani"
			picType = "ani"
		}

		if gameType == "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(gameQuery,
							table, gameLower, topicLower, roles[n], picType, i,
						))
					}
				}
			}
		} else if gameType != "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(lotteryQuery,
							table, gameLower, gameType, topicLower, roles[n], picType, i,
						))
					}
				}
			}
		}

		if gameType == "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(gameQuery2,
						table, gameLower, topicLower, roles[n], picType, i,
					))
				}
			}
		} else if gameType != "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(lotteryQuery2,
						table, gameLower, gameType, topicLower, roles[n], picType, i,
					))
				}
			}
		}
	}

	// 開始音樂
	if startMusic {
		fmt.Println(fmt.Sprintf(musicquery, table, gameLower, "bgm", "start"))
	}
	// 遊戲中音樂
	if gamingMusic {
		fmt.Println(fmt.Sprintf(musicquery, table, gameLower, "bgm", "gaming"))
	}

	// 結束音樂
	if endMusic {
		fmt.Println(fmt.Sprintf(musicquery, table, gameLower, "bgm", "end"))
	}

	log.Println("MapToModel function裡需要增加的參數: ")

	gameQuery = `a.%s%s%s%s0%d, _ = m["%s_%s_%s_%s_0%d"].(string)`         // 小於十
	gameQuery2 = `a.%s%s%s%s%d, _ = m["%s_%s_%s_%s_%d"].(string)`          // 大於十
	lotteryQuery = `a.%s%s%s%s%s0%d, _ = m["%s_%s_%s_%s_%s_0%d"].(string)` // 小於十
	lotteryQuery2 = `a.%s%s%s%s%s%d, _ = m["%s_%s_%s_%s_%s_%d"].(string)`  // 大於十
	musicquery = `a.%s%s%s, _ = m["%s_%s_%s"].(string)`                    // 音樂

	for n, start := range picStarts {
		var (
			picTypeUpper string
			picType      string
		)
		if n <= 2 {
			picTypeUpper = "Pic"
			picType = "pic"
		} else {
			picTypeUpper = "Ani"
			picType = "ani"
		}

		if gameType == "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(gameQuery,
							gameUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
							gameLower, topicLower, roles[n], picType, i,
						))
					}
				}
			}
		} else if gameType != "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(lotteryQuery,
							gameUpper, gameTypeUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
							gameLower, gameType, topicLower, roles[n], picType, i,
						))
					}
				}
			}
		}

		if gameType == "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(gameQuery2,
						gameUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
						gameLower, topicLower, roles[n], picType, i,
					))
				}
			}
		} else if gameType != "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(lotteryQuery2,
						gameUpper, gameTypeUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
						gameLower, gameType, topicLower, roles[n], picType, i,
					))
				}
			}
		}
	}

	// 開始音樂
	if startMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameUpper, "Bgm", "Start", gameLower, "bgm", "start"))
	}
	// 遊戲中音樂
	if gamingMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameUpper, "Bgm", "Gaming", gameLower, "bgm", "gaming"))

	}

	// 結束音樂
	if endMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameUpper, "Bgm", "End", gameLower, "bgm", "end"))
	}

	log.Println("MapToModel function append裡需要增加的參數: ")

	gameQuery = `%s%ss = append(%s%ss, a.%s%s%s%s0%d)`      // 小於十
	gameQuery2 = `%s%ss = append(%s%ss, a.%s%s%s%s%d)`      // 大於十
	lotteryQuery = `%s%ss = append(%s%ss, a.%s%s%s%s%s0%d)` // 小於十
	lotteryQuery2 = `%s%ss = append(%s%ss, a.%s%s%s%s%s%d)` // 大於十
	musicquery = `musics = append(musics, a.%s%s%s)`        // 音樂

	for n, start := range picStarts {
		var (
			picTypeUpper string
			picType      string
		)
		if n <= 2 {
			picTypeUpper = "Pic"
			picType = "pic"
		} else {
			picTypeUpper = "Ani"
			picType = "ani"
		}

		if gameType == "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(gameQuery,
							roles[n], picType,
							roles[n], picType,
							gameUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
						))
					}
				}
			}
		} else if gameType != "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(lotteryQuery,
							roles[n], picType,
							roles[n], picType,
							gameUpper, gameTypeUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
						))
					}
				}
			}
		}

		if gameType == "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(gameQuery2,
						roles[n], picType,
						roles[n], picType,
						gameUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
					))
				}
			}
		} else if gameType != "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(lotteryQuery2,
						roles[n], picType,
						roles[n], picType,
						gameUpper, gameTypeUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
					))
				}
			}
		}
	}

	// 開始音樂
	if startMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameUpper, "Bgm", "Start"))
	}
	// 遊戲中音樂
	if gamingMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameUpper, "Bgm", "Gaming"))

	}
	// 結束音樂
	if endMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameUpper, "Bgm", "End"))

	}

	log.Println("處理Find")

	log.Println("處理table套件")

	gameQuery = `"%s/%s/%s_%s_%s_%s_0%d.png",`       // 小於十
	gameQuery2 = `"%s/%s/%s_%s_%s_%s_%d.png",`       // 大於十
	lotteryQuery = `"%s/%s/%s_%s_%s_%s_%s_0%d.png",` // 小於十
	lotteryQuery2 = `"%s/%s/%s_%s_%s_%s_%s_%d.png",` // 大於十
	musicquery = `"%s/%s/bgm/%s.mp3",`               // 音樂

	log.Println("pics裡需要增加的參數: ")
	log.Println("記得要修改圖片為檔案類型(有些是jpg)")

	for n, start := range picStarts {
		var (
			// picTypeUpper string
			picType string
		)
		if n <= 2 {
			// picTypeUpper = "Pic"
			picType = "pic"
		} else {
			// picTypeUpper = "Ani"
			picType = "ani"
		}

		if gameType == "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(gameQuery,
							gameLower, topicLower, gameLower, topicLower, roles[n], picType, i,
						))
					}
				}
			}
		} else if gameType != "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(lotteryQuery,
							gameLower, topicLower, gameLower, gameType, topicLower, roles[n], picType, i,
						))
					}
				}
			}
		}

		if gameType == "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(gameQuery2,
						gameLower, topicLower, gameLower, topicLower, roles[n], picType, i,
					))
				}
			}
		} else if gameType != "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(lotteryQuery2,
						gameLower, topicLower, gameLower, gameType, topicLower, roles[n], picType, i,
					))
				}
			}
		}
	}

	// 開始音樂
	if startMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameLower, "%s", "start"))
	}
	// 遊戲中音樂
	if gamingMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameLower, "%s", "gaming"))
	}

	// 結束音樂
	if endMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameLower, "%s", "end"))
	}

	fmt.Println("fields裡需要增加的參數: ")

	gameQuery = `"%s_%s_%s_%s_0%d",`       // 小於十
	gameQuery2 = `"%s_%s_%s_%s_%d",`       // 大於十
	lotteryQuery = `"%s_%s_%s_%s_%s_0%d",` // 小於十
	lotteryQuery2 = `"%s_%s_%s_%s_%s_%d",` // 大於十
	musicquery = `"%s_%s_%s",`             // 音樂

	for n, start := range picStarts {
		var (
			// picTypeUpper string
			picType string
		)
		if n <= 2 {
			// picTypeUpper = "Pic"
			picType = "pic"
		} else {
			// picTypeUpper = "Ani"
			picType = "ani"
		}

		if gameType == "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(gameQuery,
							gameLower, topicLower, roles[n], picType, i,
						))
					}
				}
			}
		} else if gameType != "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(lotteryQuery,
							gameLower, gameType, topicLower, roles[n], picType, i,
						))
					}
				}
			}
		}

		if gameType == "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(gameQuery2,
						gameLower, topicLower, roles[n], picType, i,
					))
				}
			}
		} else if gameType != "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(lotteryQuery2,
						gameLower, gameType, topicLower, roles[n], picType, i,
					))
				}
			}
		}
	}

	// 開始音樂
	if startMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameLower, "bgm", "start"))
	}
	// 遊戲中音樂
	if gamingMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameLower, "bgm", "gaming"))
	}

	// 結束音樂
	if endMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameLower, "bgm", "end"))
	}

	log.Println("Add.Update function裡需要增加的參數: ")
	gameQuery = `%s%s%s%s0%d: update[],`      // 小於十
	gameQuery2 = `%s%s%s%s%d: update[],`      // 大於十
	lotteryQuery = `%s%s%s%s%s0%d: update[],` // 小於十
	lotteryQuery2 = `%s%s%s%s%s%d: update[],` // 大於十
	musicquery = `%s%s%s: update[],`          // 音樂

	for n, start := range picStarts {
		var (
			picTypeUpper string
			// picType string
		)
		if n <= 2 {
			picTypeUpper = "Pic"
			// picType = "pic"
		} else {
			picTypeUpper = "Ani"
			// picType = "ani"
		}

		if gameType == "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(gameQuery,
							gameUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
						))
					}
				}
			}
		} else if gameType != "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(lotteryQuery,
							gameUpper, gameTypeUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
						))
					}
				}
			}
		}

		if gameType == "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(gameQuery2,
						gameUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
					))
				}
			}
		} else if gameType != "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(lotteryQuery2,
						gameUpper, gameTypeUpper, topicUpper, rolesUpper[n], picTypeUpper, i,
					))
				}
			}
		}
	}

	// 開始音樂
	if startMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameUpper, "Bgm", "Start"))
	}
	// 遊戲中音樂
	if gamingMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameUpper, "Bgm", "Gaming"))
	}

	// 結束音樂
	if endMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameUpper, "Bgm", "End"))
	}

	log.Println("處理table套件")

	log.Println("處理預設值")
	log.Println("記得要修改圖片為檔案類型(有些是jpg)")

	gameQuery = `"%s_%s_%s_%s_0%d": fmt.Sprintf(route, "%s", "%s_%s_%s_%s_0%d.png"),`          // 小於十
	gameQuery2 = `"%s_%s_%s_%s_%d": fmt.Sprintf(route, "%s", "%s_%s_%s_%s_%d.png"),`           // 大於十
	lotteryQuery = `"%s_%s_%s_%s_%s_0%d": fmt.Sprintf(route, "%s", "%s_%s_%s_%s_%s_0%d.png"),` // 小於十
	lotteryQuery2 = `"%s_%s_%s_%s_%s_%d": fmt.Sprintf(route, "%s", "%s_%s_%s_%s_%s_%d.png"),`  // 大於十
	musicquery = `"%s_%s_%s": fmt.Sprintf(route, %s, "bgm/%s.png"),`                           // 音樂

	for n, start := range picStarts {
		var (
			// picTypeUpper string
			picType string
		)
		if n <= 2 {
			// picTypeUpper = "Pic"
			picType = "pic"
		} else {
			// picTypeUpper = "Ani"
			picType = "ani"
		}

		if gameType == "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(gameQuery,
							gameLower, topicLower, roles[n], picType, i,
							topicLower,
							gameLower, topicLower, roles[n], picType, i,
						))
					}
				}
			}
		} else if gameType != "" {
			if start < 10 {
				var end int
				if picEnds[n] < 10 {
					end = picEnds[n]
				} else if picEnds[n] >= 10 {
					end = 9
				}

				if start > 0 {
					for i := start; i <= end; i++ {
						fmt.Println(fmt.Sprintf(lotteryQuery,
							gameLower, gameType, topicLower, roles[n], picType, i,
							topicLower,
							gameLower, gameType, topicLower, roles[n], picType, i,
						))
					}
				}
			}
		}

		if gameType == "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(gameQuery2,
						gameLower, topicLower, roles[n], picType, i,
						topicLower,
						gameLower, topicLower, roles[n], picType, i,
					))
				}
			}
		} else if gameType != "" {
			var (
				end       int
				isExecute bool
			)
			if start >= 10 {
				end = start
				isExecute = true
			} else if picEnds[n] >= 10 {
				end = 10
				isExecute = true
			}

			if isExecute {
				for i := end; i <= picEnds[n]; i++ {
					fmt.Println(fmt.Sprintf(lotteryQuery2,
						gameLower, gameType, topicLower, roles[n], picType, i,
						topicLower,
						gameLower, gameType, topicLower, roles[n], picType, i,
					))
				}
			}
		}
	}

	// 開始音樂
	if startMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameLower, "bgm", "start", "topic", "start"))
	}
	// 遊戲中音樂
	if gamingMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameLower, "bgm", "gaming", "topic", "gaming"))
	}

	// 結束音樂
	if endMusic {
		fmt.Println(fmt.Sprintf(musicquery, gameLower, "bgm", "end", "topic", "end"))
	}

	log.Println("處理預設值")

	// if gameType == "" {
	// } else if gameType != "" {
	// }
}

// // 主持靜態
// gamehpics = `hpics = append(hpics, a.%s%s%s%s0%d)`      // 小於十
// gamehpics2 = `hpics = append(hpics, a.%s%s%s%s%d)`      // 大於十
// lotteryhpics = `hpics = append(hpics, a.%s%s%s%s%s0%d)` // 小於十
// lotteryhpics2 = `hpics = append(hpics, a.%s%s%s%s%s%d)` // 大於十
// // 玩家靜態
// gamegpics = `gpics = append(gpics, a.%s%s%s%s0%d)`      // 小於十
// gamegpics2 = `gpics = append(gpics, a.%s%s%s%s%d)`      // 大於十
// lotterygpics = `gpics = append(gpics, a.%s%s%s%s%s0%d)` // 小於十
// lotterygpics2 = `gpics = append(gpics, a.%s%s%s%s%s%d)` // 大於十
// // 共用靜態
// gamecpics = `cpics = append(cpics, a.%s%s%s%s0%d)`      // 小於十
// gamecpics2 = `cpics = append(cpics, a.%s%s%s%s%d)`      // 大於十
// lotterycpics = `cpics = append(cpics, a.%s%s%s%s%s0%d)` // 小於十
// lotterycpics2 = `cpics = append(cpics, a.%s%s%s%s%s%d)` // 大於十

// // 主持動態
// gamehanis = `hanis = append(hanis, a.%s%s%s%s0%d)`      // 小於十
// gamehanis2 = `hanis = append(hanis, a.%s%s%s%s%d)`      // 大於十
// lotteryhanis = `hanis = append(hanis, a.%s%s%s%s%s0%d)` // 小於十
// lotteryhanis2 = `hanis = append(hanis, a.%s%s%s%s%s%d)` // 大於十
// // 玩家動態
// gameganis = `ganis = append(ganis, a.%s%s%s%s0%d)`      // 小於十
// gameganis2 = `ganis = append(ganis, a.%s%s%s%s%d)`      // 大於十
// lotteryganis = `ganis = append(ganis, a.%s%s%s%s%s0%d)` // 小於十
// lotteryganis2 = `ganis = append(ganis, a.%s%s%s%s%s%d)` // 大於十
// // 共用動態
// gamecanis = `canis = append(canis, a.%s%s%s%s0%d)`      // 小於十
// gamecanis2 = `canis = append(canis, a.%s%s%s%s%d)`      // 大於十
// lotterycanis = `canis = append(canis, a.%s%s%s%s%s0%d)` // 小於十
// lotterycanis2 = `canis = append(canis, a.%s%s%s%s%s%d)` // 大�
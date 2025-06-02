package models

import (
	"errors"
	"fmt"
	"hilive/modules/config"
	"hilive/modules/utils"
	"log"
	"strconv"
	"time"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson"
)

// TableUpdateInfo 紀錄資料表、要更新的欄位及資料
// type TableUpdateInfo struct {
// 	TableName string
// 	Fields    []string
// }

// Update 更新遊戲場次資料
func (a GameModel) Update(isRedis bool, game string, model EditGameModel) error {
	var (
		gameStatus  = "close"
		filter      = bson.M{"game_id": model.GameID} // 過濾條件
		fieldValues = bson.M{}                        // 要更新的欄位資料
		fields      = make([]string, 0)

		// 扭蛋機-----start
		gachaMachinefields = []string{
			// "id",
			// "user_id",
			// "activity_id",
			// "game_id",
			// "game",
			"topic",
			"skin",
			"music",
			// "game_status",
			// "game_order",
			"title",
			"allow",
			"max_times",
			"percent",
			"reflective_switch",
			"gacha_machine_reflection",
			// "game_round",
			// "game_second",
			// "game_attend",
			// "edit_times",

			// 扭蛋遊戲自定義
			"3d_gacha_machine_classic_h_pic_02",
			"3d_gacha_machine_classic_h_pic_03",
			"3d_gacha_machine_classic_h_pic_04",
			"3d_gacha_machine_classic_h_pic_05",
			"3d_gacha_machine_classic_g_pic_01",
			"3d_gacha_machine_classic_g_pic_02",
			"3d_gacha_machine_classic_g_pic_03",
			"3d_gacha_machine_classic_g_pic_04",
			"3d_gacha_machine_classic_g_pic_05",
			"3d_gacha_machine_classic_c_pic_01",

			// 音樂
			"3d_gacha_machine_bgm_gaming", // 遊戲進行中
		}

		// 扭蛋機-----end

		// 賓果-----start
		bingofields = []string{
			// "id",
			// "user_id",
			// "activity_id",
			// "game_id",
			// "game",
			"topic",
			"skin",
			"music",
			// "game_status",
			// "game_order",
			"title",
			"max_people",
			"people",
			"allow",
			"max_number",
			"bingo_line",
			"round_prize",
			// "game_round",
			// "game_second",
			// "game_attend",
			// "bingo_round",
			// "edit_times",

			// 賓果遊戲自定義
			// 音樂
			"bingo_bgm_start",  // 遊戲開始
			"bingo_bgm_gaming", // 遊戲進行中
			"bingo_bgm_end",    // 遊戲結束

			"bingo_classic_h_pic_01",
			"bingo_classic_h_pic_02",
			"bingo_classic_h_pic_03",
			"bingo_classic_h_pic_04",
			"bingo_classic_h_pic_05",
			"bingo_classic_h_pic_06",
			"bingo_classic_h_pic_07",
			"bingo_classic_h_pic_08",
			"bingo_classic_h_pic_09",
			"bingo_classic_h_pic_10",
			"bingo_classic_h_pic_11",
			"bingo_classic_h_pic_12",
			"bingo_classic_h_pic_13",
			"bingo_classic_h_pic_14",
			"bingo_classic_h_pic_15",
			"bingo_classic_h_pic_16",
			"bingo_classic_g_pic_01",
			"bingo_classic_g_pic_02",
			"bingo_classic_g_pic_03",
			"bingo_classic_g_pic_04",
			"bingo_classic_g_pic_05",
			"bingo_classic_g_pic_06",
			"bingo_classic_g_pic_07",
			"bingo_classic_g_pic_08",
			"bingo_classic_c_pic_01",
			"bingo_classic_c_pic_02",
			"bingo_classic_c_pic_03",
			"bingo_classic_c_pic_04",
			"bingo_classic_h_ani_01",
			"bingo_classic_g_ani_01",
			"bingo_classic_c_ani_01",
			"bingo_classic_c_ani_02",

			"bingo_newyear_dragon_h_pic_01",
			"bingo_newyear_dragon_h_pic_02",
			"bingo_newyear_dragon_h_pic_03",
			"bingo_newyear_dragon_h_pic_04",
			"bingo_newyear_dragon_h_pic_05",
			"bingo_newyear_dragon_h_pic_06",
			"bingo_newyear_dragon_h_pic_07",
			"bingo_newyear_dragon_h_pic_08",
			"bingo_newyear_dragon_h_pic_09",
			"bingo_newyear_dragon_h_pic_10",
			"bingo_newyear_dragon_h_pic_11",
			"bingo_newyear_dragon_h_pic_12",
			"bingo_newyear_dragon_h_pic_13",
			"bingo_newyear_dragon_h_pic_14",
			"bingo_newyear_dragon_h_pic_16",
			"bingo_newyear_dragon_h_pic_17",
			"bingo_newyear_dragon_h_pic_18",
			"bingo_newyear_dragon_h_pic_19",
			"bingo_newyear_dragon_h_pic_20",
			"bingo_newyear_dragon_h_pic_21",
			"bingo_newyear_dragon_h_pic_22",
			"bingo_newyear_dragon_g_pic_01",
			"bingo_newyear_dragon_g_pic_02",
			"bingo_newyear_dragon_g_pic_03",
			"bingo_newyear_dragon_g_pic_04",
			"bingo_newyear_dragon_g_pic_05",
			"bingo_newyear_dragon_g_pic_06",
			"bingo_newyear_dragon_g_pic_07",
			"bingo_newyear_dragon_g_pic_08",
			"bingo_newyear_dragon_c_pic_01",
			"bingo_newyear_dragon_c_pic_02",
			"bingo_newyear_dragon_c_pic_03",
			"bingo_newyear_dragon_h_ani_01",
			"bingo_newyear_dragon_g_ani_01",
			"bingo_newyear_dragon_c_ani_01",
			"bingo_newyear_dragon_c_ani_02",
			"bingo_newyear_dragon_c_ani_03",

			"bingo_cherry_h_pic_01",
			"bingo_cherry_h_pic_02",
			"bingo_cherry_h_pic_03",
			"bingo_cherry_h_pic_04",
			"bingo_cherry_h_pic_05",
			"bingo_cherry_h_pic_06",
			"bingo_cherry_h_pic_07",
			"bingo_cherry_h_pic_08",
			"bingo_cherry_h_pic_09",
			"bingo_cherry_h_pic_10",
			"bingo_cherry_h_pic_11",
			"bingo_cherry_h_pic_12",
			"bingo_cherry_h_pic_14",
			"bingo_cherry_h_pic_15",
			"bingo_cherry_h_pic_17",
			"bingo_cherry_h_pic_18",
			"bingo_cherry_h_pic_19",
			"bingo_cherry_g_pic_01",
			"bingo_cherry_g_pic_02",
			"bingo_cherry_g_pic_03",
			"bingo_cherry_g_pic_04",
			"bingo_cherry_g_pic_05",
			"bingo_cherry_g_pic_06",
			"bingo_cherry_c_pic_01",
			"bingo_cherry_c_pic_02",
			"bingo_cherry_c_pic_03",
			"bingo_cherry_c_pic_04",
			"bingo_cherry_h_ani_02",
			"bingo_cherry_h_ani_03",
			"bingo_cherry_g_ani_01",
			"bingo_cherry_g_ani_02",
			"bingo_cherry_c_ani_01",
			"bingo_cherry_c_ani_02",
		}

		// 賓果-----end

		// 搖號抽獎-----start
		drawNumbersfields = []string{
			// "id",
			// "user_id",
			// "activity_id",
			// "game_id",
			// "game",
			"topic",
			"skin",
			"music",
			// "game_status",
			// "game_order",
			"title",
			"limit_time",
			"second",
			"allow",
			"display_name",
			// "game_round",
			"game_second",
			// "game_attend",
			// "edit_times",

			// 搖號抽獎自定義
			// 音樂
			"draw_numbers_bgm_gaming",

			"draw_numbers_classic_h_pic_01",
			"draw_numbers_classic_h_pic_02",
			"draw_numbers_classic_h_pic_03",
			"draw_numbers_classic_h_pic_04",
			"draw_numbers_classic_h_pic_05",
			"draw_numbers_classic_h_pic_06",
			"draw_numbers_classic_h_pic_07",
			"draw_numbers_classic_h_pic_08",
			"draw_numbers_classic_h_pic_09",
			"draw_numbers_classic_h_pic_10",
			"draw_numbers_classic_h_pic_11",
			"draw_numbers_classic_h_pic_12",
			"draw_numbers_classic_h_pic_13",
			"draw_numbers_classic_h_pic_14",
			"draw_numbers_classic_h_pic_15",
			"draw_numbers_classic_h_pic_16",
			"draw_numbers_classic_h_ani_01",

			"draw_numbers_gold_h_pic_01",
			"draw_numbers_gold_h_pic_02",
			"draw_numbers_gold_h_pic_03",
			"draw_numbers_gold_h_pic_04",
			"draw_numbers_gold_h_pic_05",
			"draw_numbers_gold_h_pic_06",
			"draw_numbers_gold_h_pic_07",
			"draw_numbers_gold_h_pic_08",
			"draw_numbers_gold_h_pic_09",
			"draw_numbers_gold_h_pic_10",
			"draw_numbers_gold_h_pic_11",
			"draw_numbers_gold_h_pic_12",
			"draw_numbers_gold_h_pic_13",
			"draw_numbers_gold_h_pic_14",
			"draw_numbers_gold_h_ani_01",
			"draw_numbers_gold_h_ani_02",
			"draw_numbers_gold_h_ani_03",

			"draw_numbers_newyear_dragon_h_pic_01",
			"draw_numbers_newyear_dragon_h_pic_02",
			"draw_numbers_newyear_dragon_h_pic_03",
			"draw_numbers_newyear_dragon_h_pic_04",
			"draw_numbers_newyear_dragon_h_pic_05",
			"draw_numbers_newyear_dragon_h_pic_06",
			"draw_numbers_newyear_dragon_h_pic_07",
			"draw_numbers_newyear_dragon_h_pic_08",
			"draw_numbers_newyear_dragon_h_pic_09",
			"draw_numbers_newyear_dragon_h_pic_10",
			"draw_numbers_newyear_dragon_h_pic_11",
			"draw_numbers_newyear_dragon_h_pic_12",
			"draw_numbers_newyear_dragon_h_pic_13",
			"draw_numbers_newyear_dragon_h_pic_14",
			"draw_numbers_newyear_dragon_h_pic_15",
			"draw_numbers_newyear_dragon_h_pic_16",
			"draw_numbers_newyear_dragon_h_pic_17",
			"draw_numbers_newyear_dragon_h_pic_18",
			"draw_numbers_newyear_dragon_h_pic_19",
			"draw_numbers_newyear_dragon_h_pic_20",
			"draw_numbers_newyear_dragon_h_ani_01",
			"draw_numbers_newyear_dragon_h_ani_02",

			"draw_numbers_cherry_h_pic_01",
			"draw_numbers_cherry_h_pic_02",
			"draw_numbers_cherry_h_pic_03",
			"draw_numbers_cherry_h_pic_04",
			"draw_numbers_cherry_h_pic_05",
			"draw_numbers_cherry_h_pic_06",
			"draw_numbers_cherry_h_pic_07",
			"draw_numbers_cherry_h_pic_08",
			"draw_numbers_cherry_h_pic_09",
			"draw_numbers_cherry_h_pic_10",
			"draw_numbers_cherry_h_pic_11",
			"draw_numbers_cherry_h_pic_12",
			"draw_numbers_cherry_h_pic_13",
			"draw_numbers_cherry_h_pic_14",
			"draw_numbers_cherry_h_pic_15",
			"draw_numbers_cherry_h_pic_16",
			"draw_numbers_cherry_h_pic_17",
			"draw_numbers_cherry_h_ani_01",
			"draw_numbers_cherry_h_ani_02",
			"draw_numbers_cherry_h_ani_03",
			"draw_numbers_cherry_h_ani_04",

			"draw_numbers_3D_space_h_pic_01",
			"draw_numbers_3D_space_h_pic_02",
			"draw_numbers_3D_space_h_pic_03",
			"draw_numbers_3D_space_h_pic_04",
			"draw_numbers_3D_space_h_pic_05",
			"draw_numbers_3D_space_h_pic_06",
			"draw_numbers_3D_space_h_pic_07",
			"draw_numbers_3D_space_h_pic_08",
		}

		// 搖號抽獎-----end

		// 遊戲抽獎-----start
		lotteryfields = []string{
			// "id",
			// "user_id",
			// "activity_id",
			// "game_id",
			// "game",
			"topic",
			"skin",
			"music",
			// "game_status",
			// "game_order",
			"title",
			"allow",
			"max_times",
			"game_type",
			// "game_round",
			// "game_second",
			// "game_attend",
			// "edit_times",

			"lottery_bgm_gaming",

			"lottery_jiugongge_classic_h_pic_01",
			"lottery_jiugongge_classic_h_pic_02",
			"lottery_jiugongge_classic_h_pic_03",
			"lottery_jiugongge_classic_h_pic_04",
			"lottery_jiugongge_classic_g_pic_01",
			"lottery_jiugongge_classic_g_pic_02",
			"lottery_jiugongge_classic_c_pic_01",
			"lottery_jiugongge_classic_c_pic_02",
			"lottery_jiugongge_classic_c_pic_03",
			"lottery_jiugongge_classic_c_pic_04",
			"lottery_jiugongge_classic_c_ani_01",
			"lottery_jiugongge_classic_c_ani_02",
			"lottery_jiugongge_classic_c_ani_03",

			"lottery_turntable_classic_h_pic_01",
			"lottery_turntable_classic_h_pic_02",
			"lottery_turntable_classic_h_pic_03",
			"lottery_turntable_classic_h_pic_04",
			"lottery_turntable_classic_g_pic_01",
			"lottery_turntable_classic_g_pic_02",
			"lottery_turntable_classic_c_pic_01",
			"lottery_turntable_classic_c_pic_02",
			"lottery_turntable_classic_c_pic_03",
			"lottery_turntable_classic_c_pic_04",
			"lottery_turntable_classic_c_pic_05",
			"lottery_turntable_classic_c_pic_06",
			"lottery_turntable_classic_c_ani_01",
			"lottery_turntable_classic_c_ani_02",
			"lottery_turntable_classic_c_ani_03",

			"lottery_jiugongge_starrysky_h_pic_01",
			"lottery_jiugongge_starrysky_h_pic_02",
			"lottery_jiugongge_starrysky_h_pic_03",
			"lottery_jiugongge_starrysky_h_pic_04",
			"lottery_jiugongge_starrysky_h_pic_05",
			"lottery_jiugongge_starrysky_h_pic_06",
			"lottery_jiugongge_starrysky_h_pic_07",
			"lottery_jiugongge_starrysky_g_pic_01",
			"lottery_jiugongge_starrysky_g_pic_02",
			"lottery_jiugongge_starrysky_g_pic_03",
			"lottery_jiugongge_starrysky_g_pic_04",
			"lottery_jiugongge_starrysky_c_pic_01",
			"lottery_jiugongge_starrysky_c_pic_02",
			"lottery_jiugongge_starrysky_c_pic_03",
			"lottery_jiugongge_starrysky_c_pic_04",
			"lottery_jiugongge_starrysky_c_ani_01",
			"lottery_jiugongge_starrysky_c_ani_02",
			"lottery_jiugongge_starrysky_c_ani_03",
			"lottery_jiugongge_starrysky_c_ani_04",
			"lottery_jiugongge_starrysky_c_ani_05",
			"lottery_jiugongge_starrysky_c_ani_06",

			"lottery_turntable_starrysky_h_pic_01",
			"lottery_turntable_starrysky_h_pic_02",
			"lottery_turntable_starrysky_h_pic_03",
			"lottery_turntable_starrysky_h_pic_04",
			"lottery_turntable_starrysky_h_pic_05",
			"lottery_turntable_starrysky_h_pic_06",
			"lottery_turntable_starrysky_h_pic_07",
			"lottery_turntable_starrysky_h_pic_08",
			"lottery_turntable_starrysky_g_pic_01",
			"lottery_turntable_starrysky_g_pic_02",
			"lottery_turntable_starrysky_g_pic_03",
			"lottery_turntable_starrysky_g_pic_04",
			"lottery_turntable_starrysky_g_pic_05",
			"lottery_turntable_starrysky_c_pic_01",
			"lottery_turntable_starrysky_c_pic_02",
			"lottery_turntable_starrysky_c_pic_03",
			"lottery_turntable_starrysky_c_pic_04",
			"lottery_turntable_starrysky_c_pic_05",
			"lottery_turntable_starrysky_c_ani_01",
			"lottery_turntable_starrysky_c_ani_02",
			"lottery_turntable_starrysky_c_ani_03",
			"lottery_turntable_starrysky_c_ani_04",
			"lottery_turntable_starrysky_c_ani_05",
			"lottery_turntable_starrysky_c_ani_06",
			"lottery_turntable_starrysky_c_ani_07",
		}

		// 遊戲抽獎-----end

		// 鑑定師-----start
		monopolyfields = []string{
			// "id",
			// "user_id",
			// "activity_id",
			// "game_id",
			// "game",
			"topic",
			"skin",
			"music",
			// "game_status",
			// "game_order",
			"title",
			"people",
			"max_people",
			"allow",
			"limit_time",
			"second",
			"first_prize",
			"second_prize",
			"third_prize",
			"general_prize",
			// "game_round",
			"game_second",
			// "game_attend",
			// "edit_times",

			"monopoly_bgm_start",
			"monopoly_bgm_gaming",
			"monopoly_bgm_end",

			"monopoly_classic_h_pic_01",
			"monopoly_classic_h_pic_02",
			"monopoly_classic_h_pic_03",
			"monopoly_classic_h_pic_04",
			"monopoly_classic_h_pic_05",
			"monopoly_classic_h_pic_06",
			"monopoly_classic_h_pic_07",
			"monopoly_classic_h_pic_08",
			"monopoly_classic_g_pic_01",
			"monopoly_classic_g_pic_02",
			"monopoly_classic_g_pic_03",
			"monopoly_classic_g_pic_04",
			"monopoly_classic_g_pic_05",
			"monopoly_classic_g_pic_06",
			"monopoly_classic_g_pic_07",
			"monopoly_classic_c_pic_01",
			"monopoly_classic_c_pic_02",
			"monopoly_classic_g_ani_01",
			"monopoly_classic_g_ani_02",
			"monopoly_classic_c_ani_01",

			"monopoly_redpack_h_pic_01",
			"monopoly_redpack_h_pic_02",
			"monopoly_redpack_h_pic_03",
			"monopoly_redpack_h_pic_04",
			"monopoly_redpack_h_pic_05",
			"monopoly_redpack_h_pic_06",
			"monopoly_redpack_h_pic_07",
			"monopoly_redpack_h_pic_08",
			"monopoly_redpack_h_pic_09",
			"monopoly_redpack_h_pic_10",
			"monopoly_redpack_h_pic_11",
			"monopoly_redpack_h_pic_12",
			"monopoly_redpack_h_pic_13",
			"monopoly_redpack_h_pic_14",
			"monopoly_redpack_h_pic_15",
			"monopoly_redpack_h_pic_16",
			"monopoly_redpack_g_pic_01",
			"monopoly_redpack_g_pic_02",
			"monopoly_redpack_g_pic_03",
			"monopoly_redpack_g_pic_04",
			"monopoly_redpack_g_pic_05",
			"monopoly_redpack_g_pic_06",
			"monopoly_redpack_g_pic_07",
			"monopoly_redpack_g_pic_08",
			"monopoly_redpack_g_pic_09",
			"monopoly_redpack_g_pic_10",
			"monopoly_redpack_c_pic_01",
			"monopoly_redpack_c_pic_02",
			"monopoly_redpack_c_pic_03",
			"monopoly_redpack_h_ani_01",
			"monopoly_redpack_h_ani_02",
			"monopoly_redpack_h_ani_03",
			"monopoly_redpack_g_ani_01",
			"monopoly_redpack_g_ani_02",
			"monopoly_redpack_c_ani_01",

			"monopoly_newyear_rabbit_h_pic_01",
			"monopoly_newyear_rabbit_h_pic_02",
			"monopoly_newyear_rabbit_h_pic_03",
			"monopoly_newyear_rabbit_h_pic_04",
			"monopoly_newyear_rabbit_h_pic_05",
			"monopoly_newyear_rabbit_h_pic_06",
			"monopoly_newyear_rabbit_h_pic_07",
			"monopoly_newyear_rabbit_h_pic_08",
			"monopoly_newyear_rabbit_h_pic_09",
			"monopoly_newyear_rabbit_h_pic_10",
			"monopoly_newyear_rabbit_h_pic_11",
			"monopoly_newyear_rabbit_h_pic_12",
			"monopoly_newyear_rabbit_g_pic_01",
			"monopoly_newyear_rabbit_g_pic_02",
			"monopoly_newyear_rabbit_g_pic_03",
			"monopoly_newyear_rabbit_g_pic_04",
			"monopoly_newyear_rabbit_g_pic_05",
			"monopoly_newyear_rabbit_g_pic_06",
			"monopoly_newyear_rabbit_g_pic_07",
			"monopoly_newyear_rabbit_c_pic_01",
			"monopoly_newyear_rabbit_c_pic_02",
			"monopoly_newyear_rabbit_c_pic_03",
			"monopoly_newyear_rabbit_h_ani_01",
			"monopoly_newyear_rabbit_h_ani_02",
			"monopoly_newyear_rabbit_g_ani_01",
			"monopoly_newyear_rabbit_g_ani_02",
			"monopoly_newyear_rabbit_c_ani_01",

			"monopoly_sashimi_h_pic_01",
			"monopoly_sashimi_h_pic_02",
			"monopoly_sashimi_h_pic_03",
			"monopoly_sashimi_h_pic_04",
			"monopoly_sashimi_h_pic_05",
			"monopoly_sashimi_g_pic_01",
			"monopoly_sashimi_g_pic_02",
			"monopoly_sashimi_g_pic_03",
			"monopoly_sashimi_g_pic_04",
			"monopoly_sashimi_g_pic_05",
			"monopoly_sashimi_g_pic_06",
			"monopoly_sashimi_c_pic_01",
			"monopoly_sashimi_c_pic_02",
			"monopoly_sashimi_h_ani_01",
			"monopoly_sashimi_h_ani_02",
			"monopoly_sashimi_g_ani_01",
			"monopoly_sashimi_g_ani_02",
			"monopoly_sashimi_c_ani_01",
		}

		// 鑑定師-----end

		// 快問快答-----start
		qafields = []string{
			// "id",
			// "user_id",
			// "activity_id",
			// "game_id",
			// "game",
			"topic",
			"skin",
			"music",
			// "game_status",
			// "game_order",
			"title",
			"people",
			"max_people",
			"allow",
			"limit_time",
			"second",
			"first_prize",
			"second_prize",
			"third_prize",
			"general_prize",
			// "game_round",
			"game_second",
			// "game_attend",
			// "edit_times",
			// "qa_people",

			"qa_bgm_start",  // 遊戲開始
			"qa_bgm_gaming", // 遊戲進行中
			"qa_bgm_end",    // 遊戲結束

			"qa_classic_h_pic_01",
			"qa_classic_h_pic_02",
			"qa_classic_h_pic_03",
			"qa_classic_h_pic_04",
			"qa_classic_h_pic_05",
			"qa_classic_h_pic_06",
			"qa_classic_h_pic_07",
			"qa_classic_h_pic_08",
			"qa_classic_h_pic_09",
			"qa_classic_h_pic_10",
			"qa_classic_h_pic_11",
			"qa_classic_h_pic_12",
			"qa_classic_h_pic_13",
			"qa_classic_h_pic_14",
			"qa_classic_h_pic_15",
			"qa_classic_h_pic_16",
			"qa_classic_h_pic_17",
			"qa_classic_h_pic_18",
			"qa_classic_h_pic_19",
			"qa_classic_h_pic_20",
			"qa_classic_h_pic_21",
			"qa_classic_h_pic_22",
			"qa_classic_g_pic_01",
			"qa_classic_g_pic_02",
			"qa_classic_g_pic_03",
			"qa_classic_g_pic_04",
			"qa_classic_g_pic_05",
			"qa_classic_c_pic_01",
			"qa_classic_h_ani_01",
			"qa_classic_h_ani_02",
			"qa_classic_g_ani_01",
			"qa_classic_g_ani_02",

			"qa_electric_h_pic_01",
			"qa_electric_h_pic_02",
			"qa_electric_h_pic_03",
			"qa_electric_h_pic_04",
			"qa_electric_h_pic_05",
			"qa_electric_h_pic_06",
			"qa_electric_h_pic_07",
			"qa_electric_h_pic_08",
			"qa_electric_h_pic_09",
			"qa_electric_h_pic_10",
			"qa_electric_h_pic_11",
			"qa_electric_h_pic_12",
			"qa_electric_h_pic_13",
			"qa_electric_h_pic_14",
			"qa_electric_h_pic_15",
			"qa_electric_h_pic_16",
			"qa_electric_h_pic_17",
			"qa_electric_h_pic_18",
			"qa_electric_h_pic_19",
			"qa_electric_h_pic_20",
			"qa_electric_h_pic_21",
			"qa_electric_h_pic_22",
			"qa_electric_h_pic_23",
			"qa_electric_h_pic_24",
			"qa_electric_h_pic_25",
			"qa_electric_h_pic_26",
			"qa_electric_g_pic_01",
			"qa_electric_g_pic_02",
			"qa_electric_g_pic_03",
			"qa_electric_g_pic_04",
			"qa_electric_g_pic_05",
			"qa_electric_g_pic_06",
			"qa_electric_g_pic_07",
			"qa_electric_g_pic_08",
			"qa_electric_g_pic_09",
			"qa_electric_c_pic_01",
			"qa_electric_h_ani_01",
			"qa_electric_h_ani_02",
			"qa_electric_h_ani_03",
			"qa_electric_h_ani_04",
			"qa_electric_h_ani_05",
			"qa_electric_g_ani_01",
			"qa_electric_g_ani_02",
			"qa_electric_c_ani_01",

			"qa_moonfestival_h_pic_01",
			"qa_moonfestival_h_pic_02",
			"qa_moonfestival_h_pic_03",
			"qa_moonfestival_h_pic_04",
			"qa_moonfestival_h_pic_05",
			"qa_moonfestival_h_pic_06",
			"qa_moonfestival_h_pic_07",
			"qa_moonfestival_h_pic_08",
			"qa_moonfestival_h_pic_09",
			"qa_moonfestival_h_pic_10",
			"qa_moonfestival_h_pic_11",
			"qa_moonfestival_h_pic_12",
			"qa_moonfestival_h_pic_13",
			"qa_moonfestival_h_pic_14",
			"qa_moonfestival_h_pic_15",
			"qa_moonfestival_h_pic_16",
			"qa_moonfestival_h_pic_17",
			"qa_moonfestival_h_pic_18",
			"qa_moonfestival_h_pic_19",
			"qa_moonfestival_h_pic_20",
			"qa_moonfestival_h_pic_21",
			"qa_moonfestival_h_pic_22",
			"qa_moonfestival_h_pic_23",
			"qa_moonfestival_h_pic_24",
			"qa_moonfestival_g_pic_01",
			"qa_moonfestival_g_pic_02",
			"qa_moonfestival_g_pic_03",
			"qa_moonfestival_g_pic_04",
			"qa_moonfestival_g_pic_05",
			"qa_moonfestival_c_pic_01",
			"qa_moonfestival_c_pic_02",
			"qa_moonfestival_c_pic_03",
			"qa_moonfestival_h_ani_01",
			"qa_moonfestival_h_ani_02",
			"qa_moonfestival_g_ani_01",
			"qa_moonfestival_g_ani_02",
			"qa_moonfestival_g_ani_03",

			"qa_newyear_dragon_h_pic_01",
			"qa_newyear_dragon_h_pic_02",
			"qa_newyear_dragon_h_pic_03",
			"qa_newyear_dragon_h_pic_04",
			"qa_newyear_dragon_h_pic_05",
			"qa_newyear_dragon_h_pic_06",
			"qa_newyear_dragon_h_pic_07",
			"qa_newyear_dragon_h_pic_08",
			"qa_newyear_dragon_h_pic_09",
			"qa_newyear_dragon_h_pic_10",
			"qa_newyear_dragon_h_pic_11",
			"qa_newyear_dragon_h_pic_12",
			"qa_newyear_dragon_h_pic_13",
			"qa_newyear_dragon_h_pic_14",
			"qa_newyear_dragon_h_pic_15",
			"qa_newyear_dragon_h_pic_16",
			"qa_newyear_dragon_h_pic_17",
			"qa_newyear_dragon_h_pic_18",
			"qa_newyear_dragon_h_pic_19",
			"qa_newyear_dragon_h_pic_20",
			"qa_newyear_dragon_h_pic_21",
			"qa_newyear_dragon_h_pic_22",
			"qa_newyear_dragon_h_pic_23",
			"qa_newyear_dragon_h_pic_24",
			"qa_newyear_dragon_g_pic_01",
			"qa_newyear_dragon_g_pic_02",
			"qa_newyear_dragon_g_pic_03",
			"qa_newyear_dragon_g_pic_04",
			"qa_newyear_dragon_g_pic_05",
			"qa_newyear_dragon_g_pic_06",
			"qa_newyear_dragon_c_pic_01",
			"qa_newyear_dragon_h_ani_01",
			"qa_newyear_dragon_h_ani_02",
			"qa_newyear_dragon_g_ani_01",
			"qa_newyear_dragon_g_ani_02",
			"qa_newyear_dragon_g_ani_03",
			"qa_newyear_dragon_c_ani_01",

			// "qa_1", "qa_1_options", "qa_1_answer", "qa_1_score",
			// "qa_2", "qa_2_options", "qa_2_answer", "qa_2_score",
			// "qa_3", "qa_3_options", "qa_3_answer", "qa_3_score",
			// "qa_4", "qa_4_options", "qa_4_answer", "qa_4_score",
			// "qa_5", "qa_5_options", "qa_5_answer", "qa_5_score",
			// "qa_6", "qa_6_options", "qa_6_answer", "qa_6_score",
			// "qa_7", "qa_7_options", "qa_7_answer", "qa_7_score",
			// "qa_8", "qa_8_options", "qa_8_answer", "qa_8_score",
			// "qa_9", "qa_9_options", "qa_9_answer", "qa_9_score",
			// "qa_10", "qa_10_options", "qa_10_answer", "qa_10_score",
			// "qa_11", "qa_11_options", "qa_11_answer", "qa_11_score",
			// "qa_12", "qa_12_options", "qa_12_answer", "qa_12_score",
			// "qa_13", "qa_13_options", "qa_13_answer", "qa_13_score",
			// "qa_14", "qa_14_options", "qa_14_answer", "qa_14_score",
			// "qa_15", "qa_15_options", "qa_15_answer", "qa_15_score",
			// "qa_16", "qa_16_options", "qa_16_answer", "qa_16_score",
			// "qa_17", "qa_17_options", "qa_17_answer", "qa_17_score",
			// "qa_18", "qa_18_options", "qa_18_answer", "qa_18_score",
			// "qa_19", "qa_19_options", "qa_19_answer", "qa_19_score",
			// "qa_20", "qa_20_options", "qa_20_answer", "qa_20_score",

			"qa_round", "qa_second", "total_qa",
		}

		// 快問快答-----end

		// 搖紅包-----start
		redpackfields = []string{
			// "id",
			// "user_id",
			// "activity_id",
			// "game_id",
			// "game",
			"topic",
			"skin",
			"music",
			// "game_status",
			// "game_order",
			"title",
			"people",
			"max_people",
			"allow",
			"limit_time",
			"second",
			"percent",
			// "game_round",
			"game_second",
			// "game_attend",
			// "edit_times",

			"redpack_bgm_start",
			"redpack_bgm_gaming",
			"redpack_bgm_end",

			"redpack_classic_h_pic_01",
			"redpack_classic_h_pic_02",
			"redpack_classic_h_pic_03",
			"redpack_classic_h_pic_04",
			"redpack_classic_h_pic_05",
			"redpack_classic_h_pic_06",
			"redpack_classic_h_pic_07",
			"redpack_classic_h_pic_08",
			"redpack_classic_h_pic_09",
			"redpack_classic_h_pic_10",
			"redpack_classic_h_pic_11",
			"redpack_classic_h_pic_12",
			"redpack_classic_h_pic_13",
			"redpack_classic_g_pic_01",
			"redpack_classic_g_pic_02",
			"redpack_classic_g_pic_03",
			"redpack_classic_h_ani_01",
			"redpack_classic_h_ani_02",
			"redpack_classic_g_ani_01",
			"redpack_classic_g_ani_02",
			"redpack_classic_g_ani_03",

			"redpack_cherry_h_pic_01",
			"redpack_cherry_h_pic_02",
			"redpack_cherry_h_pic_03",
			"redpack_cherry_h_pic_04",
			"redpack_cherry_h_pic_05",
			"redpack_cherry_h_pic_06",
			"redpack_cherry_h_pic_07",
			"redpack_cherry_g_pic_01",
			"redpack_cherry_g_pic_02",
			"redpack_cherry_h_ani_01",
			"redpack_cherry_h_ani_02",
			"redpack_cherry_g_ani_01",
			"redpack_cherry_g_ani_02",

			"redpack_christmas_h_pic_01",
			"redpack_christmas_h_pic_02",
			"redpack_christmas_h_pic_03",
			"redpack_christmas_h_pic_04",
			"redpack_christmas_h_pic_05",
			"redpack_christmas_h_pic_06",
			"redpack_christmas_h_pic_07",
			"redpack_christmas_h_pic_08",
			"redpack_christmas_h_pic_09",
			"redpack_christmas_h_pic_10",
			"redpack_christmas_h_pic_11",
			"redpack_christmas_h_pic_12",
			"redpack_christmas_h_pic_13",
			"redpack_christmas_g_pic_01",
			"redpack_christmas_g_pic_02",
			"redpack_christmas_g_pic_03",
			"redpack_christmas_g_pic_04",
			"redpack_christmas_c_pic_01",
			"redpack_christmas_c_pic_02",
			"redpack_christmas_c_ani_01",
		}

		// 搖紅包-----end

		// 套紅包-----start
		ropepackfields = []string{
			// "id",
			// "user_id",
			// "activity_id",
			// "game_id",
			// "game",
			"topic",
			"skin",
			"music",
			// "game_status",
			// "game_order",
			"title",
			"people",
			"max_people",
			"allow",
			"limit_time",
			"second",
			"percent",
			// "game_round",
			"game_second",
			// "game_attend",
			// "edit_times",

			"ropepack_bgm_start",
			"ropepack_bgm_gaming",
			"ropepack_bgm_end",

			"ropepack_classic_h_pic_01",
			"ropepack_classic_h_pic_02",
			"ropepack_classic_h_pic_03",
			"ropepack_classic_h_pic_04",
			"ropepack_classic_h_pic_05",
			"ropepack_classic_h_pic_06",
			"ropepack_classic_h_pic_07",
			"ropepack_classic_h_pic_08",
			"ropepack_classic_h_pic_09",
			"ropepack_classic_h_pic_10",
			"ropepack_classic_g_pic_01",
			"ropepack_classic_g_pic_02",
			"ropepack_classic_g_pic_03",
			"ropepack_classic_g_pic_04",
			"ropepack_classic_g_pic_05",
			"ropepack_classic_g_pic_06",
			"ropepack_classic_h_ani_01",
			"ropepack_classic_g_ani_01",
			"ropepack_classic_g_ani_02",
			"ropepack_classic_c_ani_01",

			"ropepack_newyear_rabbit_h_pic_01",
			"ropepack_newyear_rabbit_h_pic_02",
			"ropepack_newyear_rabbit_h_pic_03",
			"ropepack_newyear_rabbit_h_pic_04",
			"ropepack_newyear_rabbit_h_pic_05",
			"ropepack_newyear_rabbit_h_pic_06",
			"ropepack_newyear_rabbit_h_pic_07",
			"ropepack_newyear_rabbit_h_pic_08",
			"ropepack_newyear_rabbit_h_pic_09",
			"ropepack_newyear_rabbit_g_pic_01",
			"ropepack_newyear_rabbit_g_pic_02",
			"ropepack_newyear_rabbit_g_pic_03",
			"ropepack_newyear_rabbit_h_ani_01",
			"ropepack_newyear_rabbit_g_ani_01",
			"ropepack_newyear_rabbit_g_ani_02",
			"ropepack_newyear_rabbit_g_ani_03",
			"ropepack_newyear_rabbit_c_ani_01",
			"ropepack_newyear_rabbit_c_ani_02",

			"ropepack_moonfestival_h_pic_01",
			"ropepack_moonfestival_h_pic_02",
			"ropepack_moonfestival_h_pic_03",
			"ropepack_moonfestival_h_pic_04",
			"ropepack_moonfestival_h_pic_05",
			"ropepack_moonfestival_h_pic_06",
			"ropepack_moonfestival_h_pic_07",
			"ropepack_moonfestival_h_pic_08",
			"ropepack_moonfestival_h_pic_09",
			"ropepack_moonfestival_g_pic_01",
			"ropepack_moonfestival_g_pic_02",
			"ropepack_moonfestival_c_pic_01",
			"ropepack_moonfestival_h_ani_01",
			"ropepack_moonfestival_g_ani_01",
			"ropepack_moonfestival_g_ani_02",
			"ropepack_moonfestival_c_ani_01",
			"ropepack_moonfestival_c_ani_02",

			"ropepack_3D_h_pic_01",
			"ropepack_3D_h_pic_02",
			"ropepack_3D_h_pic_03",
			"ropepack_3D_h_pic_04",
			"ropepack_3D_h_pic_05",
			"ropepack_3D_h_pic_06",
			"ropepack_3D_h_pic_07",
			"ropepack_3D_h_pic_08",
			"ropepack_3D_h_pic_09",
			"ropepack_3D_h_pic_10",
			"ropepack_3D_h_pic_11",
			"ropepack_3D_h_pic_12",
			"ropepack_3D_h_pic_13",
			"ropepack_3D_h_pic_14",
			"ropepack_3D_h_pic_15",
			"ropepack_3D_g_pic_01",
			"ropepack_3D_g_pic_02",
			"ropepack_3D_g_pic_03",
			"ropepack_3D_g_pic_04",
			"ropepack_3D_h_ani_01",
			"ropepack_3D_h_ani_02",
			"ropepack_3D_h_ani_03",
			"ropepack_3D_g_ani_01",
			"ropepack_3D_g_ani_02",
			"ropepack_3D_c_ani_01",
		}

		// 套紅包-----end

		// 拔河遊戲-----start
		tugofwarfields = []string{
			// "id",
			// "user_id",
			// "activity_id",
			// "game_id",
			// "game",
			"topic",
			"skin",
			"music",
			// "game_status",
			// "game_order",
			"title",
			"people",
			"max_people",
			"allow",
			"limit_time",
			"second",
			"allow_choose_team",
			"left_team_name",
			"right_team_name",
			"left_team_picture",
			"right_team_picture",
			"prize",
			"same_people",
			// "game_round",
			"game_second",
			// "game_attend",
			// "edit_times",

			"tugofwar_bgm_start",  // 遊戲開始
			"tugofwar_bgm_gaming", // 遊戲進行中
			"tugofwar_bgm_end",    // 遊戲結束

			"tugofwar_classic_h_pic_01",
			"tugofwar_classic_h_pic_02",
			"tugofwar_classic_h_pic_03",
			"tugofwar_classic_h_pic_04",
			"tugofwar_classic_h_pic_05",
			"tugofwar_classic_h_pic_06",
			"tugofwar_classic_h_pic_07",
			"tugofwar_classic_h_pic_08",
			"tugofwar_classic_h_pic_09",
			"tugofwar_classic_h_pic_10",
			"tugofwar_classic_h_pic_11",
			"tugofwar_classic_h_pic_12",
			"tugofwar_classic_h_pic_13",
			"tugofwar_classic_h_pic_14",
			"tugofwar_classic_h_pic_15",
			"tugofwar_classic_h_pic_16",
			"tugofwar_classic_h_pic_17",
			"tugofwar_classic_h_pic_18",
			"tugofwar_classic_h_pic_19",
			"tugofwar_classic_g_pic_01",
			"tugofwar_classic_g_pic_02",
			"tugofwar_classic_g_pic_03",
			"tugofwar_classic_g_pic_04",
			"tugofwar_classic_g_pic_05",
			"tugofwar_classic_g_pic_06",
			"tugofwar_classic_g_pic_07",
			"tugofwar_classic_g_pic_08",
			"tugofwar_classic_g_pic_09",
			"tugofwar_classic_h_ani_01",
			"tugofwar_classic_h_ani_02",
			"tugofwar_classic_h_ani_03",
			"tugofwar_classic_c_ani_01",

			"tugofwar_school_h_pic_01",
			"tugofwar_school_h_pic_02",
			"tugofwar_school_h_pic_03",
			"tugofwar_school_h_pic_04",
			"tugofwar_school_h_pic_05",
			"tugofwar_school_h_pic_06",
			"tugofwar_school_h_pic_07",
			"tugofwar_school_h_pic_08",
			"tugofwar_school_h_pic_09",
			"tugofwar_school_h_pic_10",
			"tugofwar_school_h_pic_11",
			"tugofwar_school_h_pic_12",
			"tugofwar_school_h_pic_13",
			"tugofwar_school_h_pic_14",
			"tugofwar_school_h_pic_15",
			"tugofwar_school_h_pic_16",
			"tugofwar_school_h_pic_17",
			"tugofwar_school_h_pic_18",
			"tugofwar_school_h_pic_19",
			"tugofwar_school_h_pic_20",
			"tugofwar_school_h_pic_21",
			"tugofwar_school_h_pic_22",
			"tugofwar_school_h_pic_23",
			"tugofwar_school_h_pic_24",
			"tugofwar_school_h_pic_25",
			"tugofwar_school_h_pic_26",
			"tugofwar_school_g_pic_01",
			"tugofwar_school_g_pic_02",
			"tugofwar_school_g_pic_03",
			"tugofwar_school_g_pic_04",
			"tugofwar_school_g_pic_05",
			"tugofwar_school_g_pic_06",
			"tugofwar_school_g_pic_07",
			"tugofwar_school_c_pic_01",
			"tugofwar_school_c_pic_02",
			"tugofwar_school_c_pic_03",
			"tugofwar_school_c_pic_04",
			"tugofwar_school_h_ani_01",
			"tugofwar_school_h_ani_02",
			"tugofwar_school_h_ani_03",
			"tugofwar_school_h_ani_04",
			"tugofwar_school_h_ani_05",
			"tugofwar_school_h_ani_06",
			"tugofwar_school_h_ani_07",

			"tugofwar_christmas_h_pic_01",
			"tugofwar_christmas_h_pic_02",
			"tugofwar_christmas_h_pic_03",
			"tugofwar_christmas_h_pic_04",
			"tugofwar_christmas_h_pic_05",
			"tugofwar_christmas_h_pic_06",
			"tugofwar_christmas_h_pic_07",
			"tugofwar_christmas_h_pic_08",
			"tugofwar_christmas_h_pic_09",
			"tugofwar_christmas_h_pic_10",
			"tugofwar_christmas_h_pic_11",
			"tugofwar_christmas_h_pic_12",
			"tugofwar_christmas_h_pic_13",
			"tugofwar_christmas_h_pic_14",
			"tugofwar_christmas_h_pic_15",
			"tugofwar_christmas_h_pic_16",
			"tugofwar_christmas_h_pic_17",
			"tugofwar_christmas_h_pic_18",
			"tugofwar_christmas_h_pic_19",
			"tugofwar_christmas_h_pic_20",
			"tugofwar_christmas_h_pic_21",
			"tugofwar_christmas_g_pic_01",
			"tugofwar_christmas_g_pic_02",
			"tugofwar_christmas_g_pic_03",
			"tugofwar_christmas_g_pic_04",
			"tugofwar_christmas_g_pic_05",
			"tugofwar_christmas_g_pic_06",
			"tugofwar_christmas_c_pic_01",
			"tugofwar_christmas_c_pic_02",
			"tugofwar_christmas_c_pic_03",
			"tugofwar_christmas_c_pic_04",
			"tugofwar_christmas_h_ani_01",
			"tugofwar_christmas_h_ani_02",
			"tugofwar_christmas_h_ani_03",
			"tugofwar_christmas_c_ani_01",
			"tugofwar_christmas_c_ani_02",
		}

		// 拔河遊戲-----end

		// 投票-----start
		votefields = []string{
			// "id",
			// "user_id",
			// "activity_id",
			// "game_id",
			// "game",
			"topic",
			"skin",
			"music",
			"game_status",
			// "game_order",
			"title",
			"title_switch",
			"vote_screen",
			"vote_times",
			"vote_method",
			"arrangement_guest",
			"vote_method_player",
			"vote_restriction",
			"prize",
			"avatar_shape",
			"auto_play",
			"show_rank",
			"vote_start_time",
			"vote_end_time",
			// "game_round",
			// "game_second",
			// "game_attend",
			// "edit_times",

			"vote_classic_h_pic_01",
			"vote_classic_h_pic_02",
			"vote_classic_h_pic_03",
			"vote_classic_h_pic_04",
			"vote_classic_h_pic_05",
			"vote_classic_h_pic_06",
			"vote_classic_h_pic_07",
			"vote_classic_h_pic_08",
			"vote_classic_h_pic_09",
			"vote_classic_h_pic_10",
			"vote_classic_h_pic_11",
			"vote_classic_h_pic_13",
			"vote_classic_h_pic_14",
			"vote_classic_h_pic_15",
			"vote_classic_h_pic_16",
			"vote_classic_h_pic_17",
			"vote_classic_h_pic_18",
			"vote_classic_h_pic_19",
			"vote_classic_h_pic_20",
			"vote_classic_h_pic_21",
			"vote_classic_h_pic_23",
			"vote_classic_h_pic_24",
			"vote_classic_h_pic_25",
			"vote_classic_h_pic_26",
			"vote_classic_h_pic_27",
			"vote_classic_h_pic_28",
			"vote_classic_h_pic_29",
			"vote_classic_h_pic_30",
			"vote_classic_h_pic_31",
			"vote_classic_h_pic_32",
			"vote_classic_h_pic_33",
			"vote_classic_h_pic_34",
			"vote_classic_h_pic_35",
			"vote_classic_h_pic_36",
			"vote_classic_h_pic_37",
			"vote_classic_g_pic_01",
			"vote_classic_g_pic_02",
			"vote_classic_g_pic_03",
			"vote_classic_g_pic_04",
			"vote_classic_g_pic_05",
			"vote_classic_g_pic_06",
			"vote_classic_g_pic_07",
			"vote_classic_c_pic_01",
			"vote_classic_c_pic_02",
			"vote_classic_c_pic_03",
			"vote_classic_c_pic_04",

			// 音樂
			"vote_bgm_gaming", // 遊戲進行中
		}

		// 投票-----end

		// 敲敲樂-----start
		whackmolefields = []string{
			// "id",
			// "user_id",
			// "activity_id",
			// "game_id",
			// "game",
			"topic",
			"skin",
			"music",
			// "game_status",
			// "game_order",
			"title",
			"people",
			"max_people",
			"allow",
			"limit_time",
			"second",
			"first_prize",
			"second_prize",
			"third_prize",
			"general_prize",
			// "game_round",
			"game_second",
			// "game_attend",
			// "edit_times",

			"whackmole_bgm_start",
			"whackmole_bgm_gaming",
			"whackmole_bgm_end",

			"whackmole_classic_h_pic_01",
			"whackmole_classic_h_pic_02",
			"whackmole_classic_h_pic_03",
			"whackmole_classic_h_pic_04",
			"whackmole_classic_h_pic_05",
			"whackmole_classic_h_pic_06",
			"whackmole_classic_h_pic_07",
			"whackmole_classic_h_pic_08",
			"whackmole_classic_h_pic_09",
			"whackmole_classic_h_pic_10",
			"whackmole_classic_h_pic_11",
			"whackmole_classic_h_pic_12",
			"whackmole_classic_h_pic_13",
			"whackmole_classic_h_pic_14",
			"whackmole_classic_h_pic_15",
			"whackmole_classic_g_pic_01",
			"whackmole_classic_g_pic_02",
			"whackmole_classic_g_pic_03",
			"whackmole_classic_g_pic_04",
			"whackmole_classic_g_pic_05",
			"whackmole_classic_c_pic_01",
			"whackmole_classic_c_pic_02",
			"whackmole_classic_c_pic_03",
			"whackmole_classic_c_pic_04",
			"whackmole_classic_c_pic_05",
			"whackmole_classic_c_pic_06",
			"whackmole_classic_c_pic_07",
			"whackmole_classic_c_pic_08",
			"whackmole_classic_c_ani_01",

			"whackmole_halloween_h_pic_01",
			"whackmole_halloween_h_pic_02",
			"whackmole_halloween_h_pic_03",
			"whackmole_halloween_h_pic_04",
			"whackmole_halloween_h_pic_05",
			"whackmole_halloween_h_pic_06",
			"whackmole_halloween_h_pic_07",
			"whackmole_halloween_h_pic_08",
			"whackmole_halloween_h_pic_09",
			"whackmole_halloween_h_pic_10",
			"whackmole_halloween_h_pic_11",
			"whackmole_halloween_h_pic_12",
			"whackmole_halloween_h_pic_13",
			"whackmole_halloween_h_pic_14",
			"whackmole_halloween_h_pic_15",
			"whackmole_halloween_g_pic_01",
			"whackmole_halloween_g_pic_02",
			"whackmole_halloween_g_pic_03",
			"whackmole_halloween_g_pic_04",
			"whackmole_halloween_g_pic_05",
			"whackmole_halloween_c_pic_01",
			"whackmole_halloween_c_pic_02",
			"whackmole_halloween_c_pic_03",
			"whackmole_halloween_c_pic_04",
			"whackmole_halloween_c_pic_05",
			"whackmole_halloween_c_pic_06",
			"whackmole_halloween_c_pic_07",
			"whackmole_halloween_c_pic_08",
			"whackmole_halloween_c_ani_01",

			"whackmole_christmas_h_pic_01",
			"whackmole_christmas_h_pic_02",
			"whackmole_christmas_h_pic_03",
			"whackmole_christmas_h_pic_04",
			"whackmole_christmas_h_pic_05",
			"whackmole_christmas_h_pic_06",
			"whackmole_christmas_h_pic_07",
			"whackmole_christmas_h_pic_08",
			"whackmole_christmas_h_pic_09",
			"whackmole_christmas_h_pic_10",
			"whackmole_christmas_h_pic_11",
			"whackmole_christmas_h_pic_12",
			"whackmole_christmas_h_pic_13",
			"whackmole_christmas_h_pic_14",
			"whackmole_christmas_g_pic_01",
			"whackmole_christmas_g_pic_02",
			"whackmole_christmas_g_pic_03",
			"whackmole_christmas_g_pic_04",
			"whackmole_christmas_g_pic_05",
			"whackmole_christmas_g_pic_06",
			"whackmole_christmas_g_pic_07",
			"whackmole_christmas_g_pic_08",
			"whackmole_christmas_c_pic_01",
			"whackmole_christmas_c_pic_02",
			"whackmole_christmas_c_pic_03",
			"whackmole_christmas_c_pic_04",
			"whackmole_christmas_c_pic_05",
			"whackmole_christmas_c_pic_06",
			"whackmole_christmas_c_pic_07",
			"whackmole_christmas_c_pic_08",
			"whackmole_christmas_c_ani_01",
			"whackmole_christmas_c_ani_02",
		}

		// 敲敲樂-----end
	)

	if utf8.RuneCountInString(model.Title) > 20 {
		return errors.New("錯誤: 標題上限為20個字元，請輸入有效的標題名稱")
	}

	if model.GameType != "" && game == "lottery" &&
		(model.GameType != "turntable" && model.GameType != "jiugongge") {
		return errors.New("錯誤: 遊戲類型資料發生問題，請輸入有效的遊戲類型")
	}

	if model.LimitTime != "" && (model.LimitTime != "open" && model.LimitTime != "close") {
		return errors.New("錯誤: 是否限時資料發生問題，請輸入有效的資料")
	}

	// 更新秒數時也必須更新game_second的秒數資料
	if model.Second != "" {
		if _, err := strconv.Atoi(model.Second); err != nil {
			return errors.New("錯誤: 限時秒數資料發生問題，請輸入有效的秒數")
		}
		// fieldValues["second"] = utils.GetInt64(model.Second, 0)
		// fieldValues["game_second"] = utils.GetInt64(model.Second, 0)
	}

	// 判斷遊戲人數上限
	if model.MaxPeople != "" && model.People != "" {
		// 判斷遊戲人數上限
		maxPeopleInt, err1 := strconv.Atoi(model.MaxPeople)
		peopleInt, err2 := strconv.Atoi(model.People)
		if err1 != nil || err2 != nil || peopleInt > maxPeopleInt {
			return errors.New("錯誤: 遊戲人數上限資料發生問題，請輸入有效的遊戲人數上限")
		}
		// fieldValues["max_people"] = utils.GetInt64(model.MaxPeople, 0)
		// fieldValues["people"] = utils.GetInt64(model.People, 0)
	}

	if model.MaxTimes != "" {
		if _, err := strconv.Atoi(model.MaxTimes); err != nil {
			return errors.New("錯誤: 遊戲上限次數發生問題，請輸入有效的遊戲次數")
		}
	}

	if model.Allow != "" && (model.Allow != "open" && model.Allow != "close") {
		return errors.New("錯誤: 允許重複搖中資料發生問題，請輸入有效的資料")
	}

	if model.Percent != "" {
		if percentInt, err := strconv.Atoi(model.Percent); err != nil ||
			percentInt > 100 || percentInt < 0 {
			return errors.New("錯誤: 中獎機率必須為0-100，請輸入有效的中獎機率")
		}
	}

	if model.FirstPrize != "" {
		if people, err := strconv.Atoi(model.FirstPrize); err != nil || people > 50 {
			return errors.New("錯誤: 頭獎中獎人數上限資料發生問題，請輸入有效的人數")
		}
	}
	if model.SecondPrize != "" {
		if people, err := strconv.Atoi(model.SecondPrize); err != nil || people > 50 {
			return errors.New("錯誤: 二獎中獎人數上限資料發生問題，請輸入有效的人數")
		}
	}
	if model.ThirdPrize != "" {
		if people, err := strconv.Atoi(model.ThirdPrize); err != nil || people > 100 {
			return errors.New("錯誤: 三獎中獎人數上限資料發生問題，請輸入有效的人數")
		}
	}
	if model.GeneralPrize != "" {
		if people, err := strconv.Atoi(model.GeneralPrize); err != nil || people > 800 {
			return errors.New("錯誤: 普通獎中獎人數上限資料發生問題，請輸入有效的人數")
		}
	}

	// 	1. 搖紅包 : 經典主題 (01_classic)、櫻花主題 (02_cherry)、03_christmas
	// 2. 套紅包 : 經典主題 (01_classic)、兔年主題 (02_newyear_rabbit)、中秋主題 (03_moonfestival)
	// 3. 敲敲樂 : 經典主題 (01_classic)、萬聖節主題 (02_halloween)、聖誕節主題 (03_christmas)
	// 4. 快問快答 : 經典主題 (01_classic)、電路主題 (02_electric)、中秋主題 (03_moonfestival)
	// 5. 抽獎遊戲 : 經典主題 (01_classic)、星空主題 (02_starrysky)
	// 6. 鑑定師 : 經典主題 (01_classic)、紅包主題(02_redpack)、兔年主題 (03_newyear_rabbit)、生魚片主題(04_sashimi)
	// 7. 搖號抽獎 : 經典主題 (01_classic) 、 黃金主題 (02_gold)、 (04_cherry)、 (05_3D_space)
	// 8. 拔河遊戲: 經典主題 (01_classic)、校園主題 (02_school)
	// 9. 賓果戲: 經典主題 (01_classic)
	if model.Topic != "" &&
		(model.Topic != "01_classic" && model.Topic != "02_halloween" &&
			model.Topic != "02_newyear_rabbit" &&
			model.Topic != "02_gold" && model.Topic != "02_electric" &&
			model.Topic != "02_starrysky" && model.Topic != "02_cherry" &&
			model.Topic != "02_redpack" && model.Topic != "03_newyear_rabbit" &&
			model.Topic != "04_sashimi" && model.Topic != "03_moonfestival" &&
			model.Topic != "03_christmas" && model.Topic != "02_school" &&
			model.Topic != "02_newyear_dragon" &&
			model.Topic != "03_newyear_dragon" &&
			model.Topic != "04_newyear_dragon" &&
			model.Topic != "04_cherry" &&
			model.Topic != "03_cherry" &&
			model.Topic != "04_3D" &&
			model.Topic != "05_3D_space") {
		return errors.New("錯誤: 主題資料發生問題，請輸入有效的主題")
	}

	if model.Skin != "" && (model.Skin != "classic" && model.Skin != "customize") {
		return errors.New("錯誤: 樣式資料發生問題，請輸入有效的樣式")
	}

	if model.Music != "" && (model.Music != "classic" && model.Music != "customize") {
		return errors.New("錯誤: 音樂資料發生問題，請輸入有效的音樂")
	}

	if model.DisplayName != "" {
		if model.DisplayName != "open" && model.DisplayName != "close" {
			return errors.New("錯誤: 是否顯示中獎人員姓名頭像資料發生問題，請輸入有效的資料")
		}
	}

	if game == "QA" {
		if model.TotalQA != "" {
			if _, err := strconv.Atoi(model.TotalQA); err != nil {
				return errors.New("錯誤: 總題目數量發生問題，請輸入有效的題目數量")
			}
		}
		if model.QASecond != "" {
			if _, err := strconv.Atoi(model.QASecond); err != nil {
				return errors.New("錯誤: 題目顯示秒數發生問題，請輸入有效的題目顯示秒數")
			}
		}
	}

	if game == "tugofwar" {
		if model.AllowChooseTeam != "" {
			if model.AllowChooseTeam != "open" && model.AllowChooseTeam != "close" {
				return errors.New("錯誤: 允許玩家選擇隊伍資料發生問題，請輸入有效的資料")
			}
		}

		if model.LeftTeamName != "" || model.RightTeamName != "" {
			if utf8.RuneCountInString(model.LeftTeamName) > 20 ||
				utf8.RuneCountInString(model.RightTeamName) > 20 {
				return errors.New("錯誤: 隊伍名稱上限為20個字元，請輸入有效的標題名稱")
			}
		}

		if model.Prize != "" {
			if model.Prize != "uniform" && model.Prize != "all" {
				return errors.New("錯誤: 獎品發放資料發生問題，請輸入有效的資料")
			}
		}
	}

	// 賓果遊戲
	if game == "bingo" {
		if model.BingoLine != "" {
			if line, err := strconv.Atoi(model.BingoLine); err != nil ||
				line < 1 || line > 10 {
				return errors.New("錯誤: 賓果連線數資料發生問題(最多10條線，最少1條線)，請輸入有效的連線數")
			}
		}

		if model.MaxNumber != "" {
			if number, err := strconv.Atoi(model.MaxNumber); err != nil ||
				number < 16 || number > 99 {
				return errors.New("錯誤: 最大號碼資料發生問題(號碼必須大於16且小於100)，請輸入有效的連線數")
			}
		}
	}

	// 投票遊戲
	if game == "vote" {
		if model.VoteScreen != "" {
			if model.VoteScreen != "bar_chart" && model.VoteScreen != "rank" && model.VoteScreen != "detail_information" {
				return errors.New("錯誤: 投票畫面資料發生問題，請輸入有效的資料")
			}
		}

		if model.VoteTimes != "" {
			if _, err := strconv.Atoi(model.VoteTimes); err != nil {
				return errors.New("錯誤: 人員投票次數資料發生問題，請輸入有效的資料")
			}
		}

		if model.VoteMethod != "" {
			if model.VoteMethod != "all_vote" && model.VoteMethod != "single_group" && model.VoteMethod != "all_group" {
				return errors.New("錯誤: 投票模式資料發生問題，請輸入有效的資料")
			}
		}

		if model.VoteMethodPlayer != "" {
			if model.VoteMethodPlayer != "one_vote" && model.VoteMethodPlayer != "free_vote" {
				return errors.New("錯誤: 玩家投票方式資料發生問題，請輸入有效的資料")
			}
		}

		if model.VoteRestriction != "" {
			if model.VoteRestriction != "all_player" && model.VoteRestriction != "special_officer" {
				return errors.New("錯誤: 投票限制資料發生問題，請輸入有效的資料")
			}
		}

		if model.Prize != "" {
			if model.Prize != "uniform" && model.Prize != "all" {
				return errors.New("錯誤: 獎品發放資料發生問題，請輸入有效的資料")
			}
		}

		if model.AvatarShape != "" {
			if model.AvatarShape != "circle" && model.AvatarShape != "square" {
				return errors.New("錯誤: 選項框資料發生問題，請輸入有效的資料")
			}
		}

		// log.Println("update: ", model.VoteStartTime,  model.VoteEndTime )
		if model.VoteStartTime != "" && model.VoteEndTime != "" {
			// var gameStatus string
			// 判斷投票結束時間
			now, _ := time.ParseInLocation("2006-01-02 15:04:05",
				time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05"), time.Local) // 目前時間
			startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", model.VoteStartTime, time.Local)
			endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", model.VoteEndTime, time.Local)

			log.Println("傳遞: ", model.VoteStartTime, model.VoteEndTime)
			log.Println("比較時間: ", now, startTime, endTime)

			if !CompareDatetime(model.VoteStartTime, model.VoteEndTime) {
				return errors.New("錯誤: 投票時間發生問題，結束時間必須大於開始時間，請輸入有效的時間")
			}

			// 比較時間，判斷遊戲狀態
			if now.Before(startTime) { // 現在時間在開始時間之前
				gameStatus = "close" // 關閉
			} else if now.Before(endTime) { // 現在時間在截止時間之前
				gameStatus = "gaming" // 遊戲中
			} else { // 現在時間在截止時間之後
				gameStatus = "calculate" // 結算狀態
			}

			// fieldValues["vote_start_time"] = model.VoteStartTime
			// fieldValues["vote_end_time"] = model.VoteEndTime
			// fieldValues["game_status"] = gameStatus
		}

		if model.AutoPlay != "" {
			if model.AutoPlay != "open" && model.AutoPlay != "close" {
				return errors.New("錯誤: 自動輪播資料發生問題，請輸入有效的資料")
			}
		}
		if model.ShowRank != "" {
			if model.ShowRank != "open" && model.ShowRank != "close" {
				return errors.New("錯誤: 排名展示資料發生問題，請輸入有效的資料")
			}
		}

		if model.TitleSwitch != "open" && model.TitleSwitch != "close" {
			return errors.New("錯誤: 場式名稱開關資料發生問題，請輸入有效的資料")
		}
		if model.ArrangementGuest != "list" && model.ArrangementGuest != "side_by_side" {
			return errors.New("錯誤: 玩家端選項排列方式資料發生問題，請輸入有效的資料")
		}
	}

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	// 手動處理
	// string轉換int64.float64
	// 當data[key]資料不為空時才會將資料寫入data(int64)
	data["game_status"] = gameStatus
	data["game_second"] = utils.GetInt64(model.Second, 0)
	data["second"] = utils.GetInt64(model.Second, 0)
	data["max_people"] = utils.GetInt64(model.MaxPeople, 0)
	data["people"] = utils.GetInt64(model.People, 0)
	data["max_times"] = utils.GetInt64(model.MaxTimes, 0)
	data["percent"] = utils.GetInt64(model.Percent, 0)
	data["first_prize"] = utils.GetInt64(model.FirstPrize, 0)
	data["second_prize"] = utils.GetInt64(model.SecondPrize, 0)
	data["third_prize"] = utils.GetInt64(model.ThirdPrize, 0)
	data["general_prize"] = utils.GetInt64(model.GeneralPrize, 0)
	data["total_qa"] = utils.GetInt64(model.TotalQA, 0)
	data["qa_second"] = utils.GetInt64(model.QASecond, 0)
	data["bingo_line"] = utils.GetInt64(model.BingoLine, 0)
	data["max_number"] = utils.GetInt64(model.MaxNumber, 0)
	data["round_prize"] = utils.GetInt64(model.RoundPrize, 0)
	data["vote_times"] = utils.GetInt64(model.VoteTimes, 0)
	data["gacha_machine_reflection"] = utils.GetFloat64(model.GachaMachineReflection, 0)
	// 編輯時不需要的欄位
	// data["game"] = game
	// data["game_round"] = int64(1)
	// data["game_attend"] = int64(0)
	// data["game_order"] = int64(len(games) + 1)
	// data["edit_times"] = int64(0)
	// data["bingo_round"] = 0
	// data["left_team_game_attend"] = 0
	// data["right_team_game_attend"] = 0

	// 將遊戲欄位資料寫入fields
	if game == "redpack" {
		fields = redpackfields
	} else if game == "ropepack" {
		fields = ropepackfields
	} else if game == "whack_mole" {
		fields = whackmolefields
	} else if game == "lottery" {
		fields = lotteryfields
	} else if game == "monopoly" {
		fields = monopolyfields
	} else if game == "QA" {
		fields = qafields
	} else if game == "draw_numbers" {
		fields = drawNumbersfields
	} else if game == "tugofwar" {
		fields = tugofwarfields
	} else if game == "bingo" {
		fields = bingofields
	} else if game == "3DGachaMachine" {
		fields = gachaMachinefields
	} else if game == "vote" {
		fields = votefields
	}

	// 快問快答
	if game == "QA" {
		if model.QA1 == "" || model.QA1Options == "" ||
			model.QA1Answer == "" || model.QA1Score == "" {
			return errors.New("錯誤: 題目設置最少一題，請重新設置")
		}

		// 更新題目時，回到第一題
		data["qa_round"] = int64(1)

		// 處理qa_n_score值
		for i := 1; i <= 20; i++ {
			// score := data[fmt.Sprintf("qa_%d_score", i)]

			// 將分數轉為int64寫入data中
			data[fmt.Sprintf("qa_%d_score", i)] = utils.GetInt64(data[fmt.Sprintf("qa_%d_score", i)], 0)
		}
	}

	// 過濾後的欄位資料
	filterData := FilterFields(data, fields)

	// 處理要更新的欄位資料，fieldValues
	for key, val := range filterData {
		// 資料不為空時才要更新至mongo，將資料寫入fieldValues
		if val != "" {
			fieldValues[key] = val
		}

		// 判斷是否為自定義欄位(_pic_、_ani_、_bgm_)
		// if strings.Contains(key, "_pic_") ||
		// 	strings.Contains(key, "_ani_") ||
		// 	strings.Contains(key, "_bgm_") {
		// 	if val != "" {
		// 		fieldValues[key] = val
		// 	}
		// } else {
		// 	// 其他欄位，直接寫入fieldValues
		// 	fieldValues[key] = val
		// }
	}

	// 快問快答
	if game == "QA" {
		// 不管題目欄位資訊是否為空，都要更新所有題目資料

		// 題目欄位
		qas := []string{
			"qa_1", "qa_1_options", "qa_1_answer", "qa_1_score",
			"qa_2", "qa_2_options", "qa_2_answer", "qa_2_score",
			"qa_3", "qa_3_options", "qa_3_answer", "qa_3_score",
			"qa_4", "qa_4_options", "qa_4_answer", "qa_4_score",
			"qa_5", "qa_5_options", "qa_5_answer", "qa_5_score",
			"qa_6", "qa_6_options", "qa_6_answer", "qa_6_score",
			"qa_7", "qa_7_options", "qa_7_answer", "qa_7_score",
			"qa_8", "qa_8_options", "qa_8_answer", "qa_8_score",
			"qa_9", "qa_9_options", "qa_9_answer", "qa_9_score",
			"qa_10", "qa_10_options", "qa_10_answer", "qa_10_score",
			"qa_11", "qa_11_options", "qa_11_answer", "qa_11_score",
			"qa_12", "qa_12_options", "qa_12_answer", "qa_12_score",
			"qa_13", "qa_13_options", "qa_13_answer", "qa_13_score",
			"qa_14", "qa_14_options", "qa_14_answer", "qa_14_score",
			"qa_15", "qa_15_options", "qa_15_answer", "qa_15_score",
			"qa_16", "qa_16_options", "qa_16_answer", "qa_16_score",
			"qa_17", "qa_17_options", "qa_17_answer", "qa_17_score",
			"qa_18", "qa_18_options", "qa_18_answer", "qa_18_score",
			"qa_19", "qa_19_options", "qa_19_answer", "qa_19_score",
			"qa_20", "qa_20_options", "qa_20_answer", "qa_20_score",
		}

		// 將所有題目資訊寫入fieldValues
		for _, qa := range qas {
			fieldValues[qa] = data[qa]
		}
	}

	// log.Println("fieldValues: ", fieldValues)

	if len(fieldValues) != 0 {
		// 更新遊戲資料表(mongo，activity_game)
		// 要更新的值
		update := bson.M{
			"$set": fieldValues,
			// "$unset": bson.M{},                // 移除不需要的欄位
			"$inc": bson.M{"edit_times": 1}, // edit 欄位遞增 1
		}

		if _, err := a.MongoConn.UpdateOne(a.TableName, filter, update); err != nil {
			return errors.New("錯誤: 編輯遊戲場次(mongo，activity_game)發生問題")
		}
	}

	if isRedis {
		// gameID := model.GameID

		// 清除遊戲redis資訊(並重新開啟遊戲頁面)
		// Redis 刪除的 key 前綴
		delKeys := []string{
			config.GAME_REDIS,                            // 遊戲設置
			config.GAME_TYPE_REDIS,                       // 遊戲種類
			config.SCORES_REDIS,                          // 分數
			config.SCORES_2_REDIS,                        // 第二分數
			config.CORRECT_REDIS,                         // 答對題數
			config.QA_REDIS,                              // 快問快答題目資訊
			config.QA_RECORD_REDIS,                       // 快問快答答題紀錄
			config.WINNING_STAFFS_REDIS,                  // 中獎人員
			config.NO_WINNING_STAFFS_REDIS,               // 未中獎人員
			config.GAME_TEAM_REDIS,                       // 遊戲隊伍資訊，HASH
			config.GAME_BINGO_NUMBER_REDIS,               // 紀錄抽過的號碼，LIST
			config.GAME_BINGO_USER_REDIS,                 // 賓果中獎人員，ZSET
			config.GAME_BINGO_USER_NUMBER,                // 紀錄玩家的號碼排序，HASH
			config.GAME_BINGO_USER_GOING_BINGO,           // 紀錄玩家是否即將賓果，HASH
			config.GAME_ATTEND_REDIS,                     // 遊戲人數資訊，SET
			config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS,  // 拔河遊戲左隊人數資訊，SET
			config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS, // 拔河遊戲右隊人數資訊，SET
		}

		for _, prefix := range delKeys {
			a.RedisConn.DelCache(prefix + model.GameID)
		}

		// Redis 發布的 channel 前綴
		publishChannels := []string{
			config.CHANNEL_GAME_REDIS,
			config.CHANNEL_GUEST_GAME_STATUS_REDIS,
			config.CHANNEL_GAME_BINGO_NUMBER_REDIS,
			config.CHANNEL_QA_REDIS,
			config.CHANNEL_GAME_TEAM_REDIS,
			config.CHANNEL_GAME_EDIT_TIMES_REDIS,
			config.CHANNEL_WINNING_STAFFS_REDIS,
			config.CHANNEL_GAME_BINGO_USER_NUMBER,
			config.CHANNEL_SCORES_REDIS,
		}

		for _, channel := range publishChannels {
			a.RedisConn.Publish(channel+model.GameID, "修改資料")
		}
	}

	return nil
}

// UpdateVoteEndTime 更新投票結束時間、結算狀態
func (a GameModel) UpdateVoteEndTime(isRedis bool, gameID string) error {
	var (
		now, _ = time.ParseInLocation("2006-01-02 15:04",
			time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04"), time.Local) // 目前時間
		// params      = []interface{}{config.GAME_REDIS + gameID, "vote_end_time", now} // redis參數
	)

	// 更新投票結束時間、結算狀態(mysql)
	// if err := a.Table(a.Base.TableName).
	// 	Where("game_id", "=", gameID).
	// 	Update(command.Value{
	// 		"game_status":   "calculate",
	// 		"vote_end_time": now,
	// 	}); err != nil &&
	// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
	// 	return errors.New("錯誤: 更新投票結束時間發生問題")
	// }

	// 更新遊戲狀態(mongo，activity_game)
	filter := bson.M{"game_id": gameID} // 過濾條件
	// 要更新的值
	update := bson.M{
		"$set": bson.M{"game_status": "calculate",
			"vote_end_time": now},
		// "$unset": bson.M{},                // 移除不需要的欄位
		// "$inc":   bson.M{"edit_times": 1}, // edit 欄位遞增 1
	}

	if _, err := a.MongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
		return errors.New("錯誤: 更新投票結束時間發生問題")
	}

	if isRedis {
		// 從redis取得資料，確定redis中有該場遊戲的設置資料
		a.Find(true, gameID, "vote")

		// 修改redis中的遊戲資訊(投票結束時間)
		a.RedisConn.HashSetCache(config.GAME_REDIS+gameID, "vote_end_time", now)

		// 修改redis中的遊戲資訊(遊戲狀態)
		a.RedisConn.HashSetCache(config.GAME_REDIS+gameID, "game_status", "calculate")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
	}

	return nil
}

// UpdateGameStatus 更新遊戲中的輪次、狀態、秒數...等資料
func (a GameModel) UpdateGameStatus(isRedis bool, gameID, round,
	second, status string) error {
	var (
		fieldValues   = bson.M{} // 更新資料表參數
		incValues     = bson.M{} // 需遞增參數
		fields        = []string{"game_round", "game_second", "game_status"}
		values        = []string{round, second, status}
		params        = []interface{}{config.GAME_REDIS + gameID} // redis參數
		isUpdateRedis bool                                        // 是否需要更新redis
	)

	for i, field := range fields {
		if values[i] != "" {
			// 資料表只更新open、close、order的遊戲狀態、game_round、game_attend資料
			if (field == "game_status" && (values[i] == "open" ||
				values[i] == "close" || values[i] == "order")) ||
				(field == "game_round") {

				if field == "game_round" {
					// game_round不為空，輪次遞增(mysql)
					// fieldValues[field] = fmt.Sprintf("%s + 1", fieldValues[field])

					// game_round不為空，輪次遞增(mongo)
					incValues[field] = 1

					// 遞增redis輪次資訊
					// a.RedisConn.HashIncrCache(config.GAME_REDIS+gameID, "game_round")

					// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
					// a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

					// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
					// a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
				} else {
					fieldValues[field] = values[i]
				}
			}

			// redis資料處理
			if field != "game_round" {
				// 將更新資訊加入redis中(game_status、game_second)
				params = append(params, fields[i], values[i])

				isUpdateRedis = true // 須更新redis
			} else if field == "game_round" {
				// 遞增redis輪次資訊
				a.RedisConn.HashIncrCache(config.GAME_REDIS+gameID, "game_round")

				isUpdateRedis = true // 須更新redis
			}
		}
	}

	// fieldValues裡有資料，需更新資料表
	if len(fieldValues) != 0 || len(incValues) != 0 {
		// mysql
		// if err := a.Table(a.Base.TableName).
		// 	Where("game_id", "=", gameID).
		// 	Update(fieldValues); err != nil &&
		// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
		// 	return err
		// }

		// 更新遊戲狀態(mongo，activity_game)
		filter := bson.M{"game_id": gameID} // 過濾條件
		// 要更新的值
		update := bson.M{
			"$set": fieldValues,
			// "$unset": bson.M{},                // 移除不需要的欄位
			"$inc": incValues, // edit 欄位遞增 1
		}

		if _, err := a.MongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
			return err
		}
	}

	if isRedis && isUpdateRedis {
		// log.Println("推入資料: ", params)

		if len(params) > 1 {
			// 修改redis中的遊戲資訊
			a.RedisConn.HashMultiSetCache(params)
		}

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "遊戲狀態改變")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
	}
	return nil
}

// UpdateQARound 更新快問快答遊戲進行題數資料(資料表、redis)
func (a GameModel) UpdateQARound(isRedis bool, gameID, round string) error {
	var (
		fieldValues = bson.M{} // 更新資料表參數
		incValues   = bson.M{} // 需遞增參數
	// roundValue string
	)

	if round == "1" {
		// 新的一輪，歸零(mysql)
		// roundValue = "1"

		// 新的一輪，歸零(mongo)
		fieldValues["qa_round"] = int64(1)
	} else if round == "+1" {
		// 遞增(mysql)
		// roundValue = "qa_round + 1"

		// 遞增(mongo)
		incValues["qa_round"] = 1
	}

	// mysql
	// if err := a.Table(config.ACTIVITY_GAME_QA_TABLE).
	// 	Where("game_id", "=", gameID).
	// 	// Where(round, "<", "qa_round").
	// 	Update(command.Value{
	// 		"qa_round": roundValue,
	// 	}); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
	// 	return err
	// }

	// 更新遊戲狀態(mongo，activity_game)
	filter := bson.M{"game_id": gameID} // 過濾條件
	// 要更新的值
	update := bson.M{
		"$set": fieldValues,
		// "$unset": bson.M{},                // 移除不需要的欄位
		"$inc": incValues, // edit 欄位遞增 1
	}

	if _, err := a.MongoConn.UpdateOne(config.ACTIVITY_GAME_TABLE, filter, update); err != nil {
		return err
	}

	// 修改redis中的遊戲資訊
	if isRedis {
		if round == "1" {
			// 新的一輪，歸零
			a.RedisConn.HashSetCache(config.GAME_REDIS+gameID,
				"qa_round", round)

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
		} else if round == "+1" {
			// fmt.Println("遞增redis")
			// 遞增
			a.RedisConn.HashIncrCache(config.GAME_REDIS+gameID,
				"qa_round")

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
		}
	}
	return nil
}

// UpdateQAPeople 更新快問快答人數資料(只更新redis資料)
func (a GameModel) UpdateQAPeople(isRedis bool, gameID string, value int64) {
	// if err := a.Table(config.ACTIVITY_GAME_QA_TABLE).
	// 	Where("game_id", "=", gameID).Update(command.Value{
	// 	"qa_people": people,
	// }); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
	// 	return err
	// }

	// 修改redis中的遊戲資訊
	if isRedis {
		if value == 0 {
			a.RedisConn.HashSetCache(config.GAME_REDIS+gameID, "qa_people", 0) // 歸零

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")
		} else if value == 1 {
			a.RedisConn.HashIncrCache(config.GAME_REDIS+gameID, "qa_people") // 遞增

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")
		} else if value == -1 {
			people := a.RedisConn.HashDecrCache(config.GAME_REDIS+gameID, "qa_people") // 遞減
			if people < 0 {
				// 人數不能為負，歸零
				// fmt.Println("負的，歸零")
				a.RedisConn.HashSetCache(config.GAME_REDIS+gameID, "qa_people", 0)
			}

			// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
			a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")
		}
		// a.RedisConn.HashSetCache(config.GAME_REDIS+gameID,
		// 	"qa_people", people)
	}
}

// UpdateShake 更新是否搖晃資料(只更新redis資料)
func (a GameModel) UpdateShake(isRedis bool, gameID string, isShake bool) {
	// 修改redis中的遊戲資訊
	if isRedis {
		a.RedisConn.HashSetCache(config.GAME_REDIS+gameID, "is_shake", strconv.FormatBool(isShake))

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")

		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
		a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+gameID, "修改資料")
	}
}

// IncrAttend 遞增遊戲人數資料
// func (a GameModel) IncrAttend(isRedis bool, gameID string) error {
// 	if err := a.Table(a.Base.TableName).
// 		WhereRaw("`game_id` = ? and `game_attend` < `people`", gameID).
// 		// Where("game_id", "=", gameID).
// 		// Where("game_attend", "<", people).
// 		Update(command.Value{"game_attend": "game_attend + 1"}); err != nil {
// 		if err.Error() == "錯誤: 無更新任何資料，請重新操作" {
// 			return errors.New("錯誤: 遊戲報名人數額滿，無法參加該輪次遊戲。請等待下一輪遊戲開始，謝謝")
// 		}
// 		return err
// 	}

// 	if isRedis {
// 		// 修改redis中的遊戲人數資訊
// 		a.RedisConn.HashIncrCache(config.GAME_REDIS+gameID, "game_attend")

// 		// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// 		a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+gameID, "修改資料")
// 	}
// 	return nil
// }

// 用戶資訊
// userModel, err := DefaultUserModel().SetDbConn(a.DbConn).
// 	SetRedisConn(a.RedisConn).Find(true, true, "",
// 	"users.user_id", model.UserID)
// if err != nil {
// 	return err
// }

// values = []string{
// 	model.Title, model.GameType, model.LimitTime,
// 	model.MaxTimes, model.Allow,
// 	model.Percent, model.FirstPrize, model.SecondPrize, model.ThirdPrize,
// 	model.GeneralPrize, model.Topic, model.Skin, model.Music,
// 	model.DisplayName,

// 	model.AllowChooseTeam,
// 	model.LeftTeamName,
// 	model.LeftTeamPicture,
// 	model.RightTeamName,
// 	model.RightTeamPicture,
// 	model.Prize,

// 	model.MaxNumber,
// 	model.BingoLine,
// 	model.RoundPrize,

// 	model.GachaMachineReflection,
// 	model.ReflectiveSwitch,

// 	model.GameOrder, // 場次排序
// 	model.BoxReflection,
// 	model.SamePeople,

// 	model.VoteScreen,
// 	model.VoteTimes,
// 	model.VoteMethod,
// 	model.VoteMethodPlayer,
// 	model.VoteRestriction,
// 	model.AvatarShape,
// 	model.AutoPlay,
// 	model.ShowRank,
// 	model.TitleSwitch,
// 	model.ArrangementGuest,
// }

// valuesRopepack = []string{
// 	// 套紅包自定義
// 	model.RopepackClassicHPic01,
// 	model.RopepackClassicHPic02,
// 	model.RopepackClassicHPic03,
// 	model.RopepackClassicHPic04,
// 	model.RopepackClassicHPic05,
// 	model.RopepackClassicHPic06,
// 	model.RopepackClassicHPic07,
// 	model.RopepackClassicHPic08,
// 	model.RopepackClassicHPic09,
// 	model.RopepackClassicHPic10,
// 	model.RopepackClassicGPic01,
// 	model.RopepackClassicGPic02,
// 	model.RopepackClassicGPic03,
// 	model.RopepackClassicGPic04,
// 	model.RopepackClassicGPic05,
// 	model.RopepackClassicGPic06,
// 	model.RopepackClassicHAni01,
// 	model.RopepackClassicGAni01,
// 	model.RopepackClassicGAni02,
// 	model.RopepackClassicCAni01,

// 	model.RopepackNewyearRabbitHPic01,
// 	model.RopepackNewyearRabbitHPic02,
// 	model.RopepackNewyearRabbitHPic03,
// 	model.RopepackNewyearRabbitHPic04,
// 	model.RopepackNewyearRabbitHPic05,
// 	model.RopepackNewyearRabbitHPic06,
// 	model.RopepackNewyearRabbitHPic07,
// 	model.RopepackNewyearRabbitHPic08,
// 	model.RopepackNewyearRabbitHPic09,
// 	model.RopepackNewyearRabbitGPic01,
// 	model.RopepackNewyearRabbitGPic02,
// 	model.RopepackNewyearRabbitGPic03,
// 	model.RopepackNewyearRabbitHAni01,
// 	model.RopepackNewyearRabbitGAni01,
// 	model.RopepackNewyearRabbitGAni02,
// 	model.RopepackNewyearRabbitGAni03,
// 	model.RopepackNewyearRabbitCAni01,
// 	model.RopepackNewyearRabbitCAni02,

// 	model.RopepackMoonfestivalHPic01,
// 	model.RopepackMoonfestivalHPic02,
// 	model.RopepackMoonfestivalHPic03,
// 	model.RopepackMoonfestivalHPic04,
// 	model.RopepackMoonfestivalHPic05,
// 	model.RopepackMoonfestivalHPic06,
// 	model.RopepackMoonfestivalHPic07,
// 	model.RopepackMoonfestivalHPic08,
// 	model.RopepackMoonfestivalHPic09,
// 	model.RopepackMoonfestivalGPic01,
// 	model.RopepackMoonfestivalGPic02,
// 	model.RopepackMoonfestivalCPic01,
// 	model.RopepackMoonfestivalHAni01,
// 	model.RopepackMoonfestivalGAni01,
// 	model.RopepackMoonfestivalGAni02,
// 	model.RopepackMoonfestivalCAni01,
// 	model.RopepackMoonfestivalCAni02,

// 	model.Ropepack3DHPic01,
// 	model.Ropepack3DHPic02,
// 	model.Ropepack3DHPic03,
// 	model.Ropepack3DHPic04,
// 	model.Ropepack3DHPic05,
// 	model.Ropepack3DHPic06,
// 	model.Ropepack3DHPic07,
// 	model.Ropepack3DHPic08,
// 	model.Ropepack3DHPic09,
// 	model.Ropepack3DHPic10,
// 	model.Ropepack3DHPic11,
// 	model.Ropepack3DHPic12,
// 	model.Ropepack3DHPic13,
// 	model.Ropepack3DHPic14,
// 	model.Ropepack3DHPic15,
// 	model.Ropepack3DGPic01,
// 	model.Ropepack3DGPic02,
// 	model.Ropepack3DGPic03,
// 	model.Ropepack3DGPic04,
// 	model.Ropepack3DHAni01,
// 	model.Ropepack3DHAni02,
// 	model.Ropepack3DHAni03,
// 	model.Ropepack3DGAni01,
// 	model.Ropepack3DGAni02,
// 	model.Ropepack3DCAni01,

// 	// 音樂
// 	model.RopepackBgmStart,
// 	model.RopepackBgmGaming,
// 	model.RopepackBgmEnd,
// }

// valuesRedpack = []string{
// 	// 搖紅包自定義
// 	model.RedpackClassicHPic01,
// 	model.RedpackClassicHPic02,
// 	model.RedpackClassicHPic03,
// 	model.RedpackClassicHPic04,
// 	model.RedpackClassicHPic05,
// 	model.RedpackClassicHPic06,
// 	model.RedpackClassicHPic07,
// 	model.RedpackClassicHPic08,
// 	model.RedpackClassicHPic09,
// 	model.RedpackClassicHPic10,
// 	model.RedpackClassicHPic11,
// 	model.RedpackClassicHPic12,
// 	model.RedpackClassicHPic13,
// 	model.RedpackClassicGPic01,
// 	model.RedpackClassicGPic02,
// 	model.RedpackClassicGPic03,
// 	model.RedpackClassicHAni01,
// 	model.RedpackClassicHAni02,
// 	model.RedpackClassicGAni01,
// 	model.RedpackClassicGAni02,
// 	model.RedpackClassicGAni03,

// 	model.RedpackCherryHPic01,
// 	model.RedpackCherryHPic02,
// 	model.RedpackCherryHPic03,
// 	model.RedpackCherryHPic04,
// 	model.RedpackCherryHPic05,
// 	model.RedpackCherryHPic06,
// 	model.RedpackCherryHPic07,
// 	model.RedpackCherryGPic01,
// 	model.RedpackCherryGPic02,
// 	model.RedpackCherryHAni01,
// 	model.RedpackCherryHAni02,
// 	model.RedpackCherryGAni01,
// 	model.RedpackCherryGAni02,

// 	model.RedpackChristmasHPic01,
// 	model.RedpackChristmasHPic02,
// 	model.RedpackChristmasHPic03,
// 	model.RedpackChristmasHPic04,
// 	model.RedpackChristmasHPic05,
// 	model.RedpackChristmasHPic06,
// 	model.RedpackChristmasHPic07,
// 	model.RedpackChristmasHPic08,
// 	model.RedpackChristmasHPic09,
// 	model.RedpackChristmasHPic10,
// 	model.RedpackChristmasHPic11,
// 	model.RedpackChristmasHPic12,
// 	model.RedpackChristmasHPic13,
// 	model.RedpackChristmasGPic01,
// 	model.RedpackChristmasGPic02,
// 	model.RedpackChristmasGPic03,
// 	model.RedpackChristmasGPic04,
// 	model.RedpackChristmasCPic01,
// 	model.RedpackChristmasCPic02,
// 	model.RedpackChristmasCAni01,

// 	// 音樂
// 	model.RedpackBgmStart,
// 	model.RedpackBgmGaming,
// 	model.RedpackBgmEnd,
// }

// valuesLottery = []string{
// 	// 遊戲抽獎自定義
// 	model.LotteryJiugonggeClassicHPic01,
// 	model.LotteryJiugonggeClassicHPic02,
// 	model.LotteryJiugonggeClassicHPic03,
// 	model.LotteryJiugonggeClassicHPic04,
// 	model.LotteryJiugonggeClassicGPic01,
// 	model.LotteryJiugonggeClassicGPic02,
// 	model.LotteryJiugonggeClassicCPic01,
// 	model.LotteryJiugonggeClassicCPic02,
// 	model.LotteryJiugonggeClassicCPic03,
// 	model.LotteryJiugonggeClassicCPic04,
// 	model.LotteryJiugonggeClassicCAni01,
// 	model.LotteryJiugonggeClassicCAni02,
// 	model.LotteryJiugonggeClassicCAni03,

// 	model.LotteryTurntableClassicHPic01,
// 	model.LotteryTurntableClassicHPic02,
// 	model.LotteryTurntableClassicHPic03,
// 	model.LotteryTurntableClassicHPic04,
// 	model.LotteryTurntableClassicGPic01,
// 	model.LotteryTurntableClassicGPic02,
// 	model.LotteryTurntableClassicCPic01,
// 	model.LotteryTurntableClassicCPic02,
// 	model.LotteryTurntableClassicCPic03,
// 	model.LotteryTurntableClassicCPic04,
// 	model.LotteryTurntableClassicCPic05,
// 	model.LotteryTurntableClassicCPic06,
// 	model.LotteryTurntableClassicCAni01,
// 	model.LotteryTurntableClassicCAni02,
// 	model.LotteryTurntableClassicCAni03,

// 	model.LotteryJiugonggeStarryskyHPic01,
// 	model.LotteryJiugonggeStarryskyHPic02,
// 	model.LotteryJiugonggeStarryskyHPic03,
// 	model.LotteryJiugonggeStarryskyHPic04,
// 	model.LotteryJiugonggeStarryskyHPic05,
// 	model.LotteryJiugonggeStarryskyHPic06,
// 	model.LotteryJiugonggeStarryskyHPic07,
// 	model.LotteryJiugonggeStarryskyGPic01,
// 	model.LotteryJiugonggeStarryskyGPic02,
// 	model.LotteryJiugonggeStarryskyGPic03,
// 	model.LotteryJiugonggeStarryskyGPic04,
// 	model.LotteryJiugonggeStarryskyCPic01,
// 	model.LotteryJiugonggeStarryskyCPic02,
// 	model.LotteryJiugonggeStarryskyCPic03,
// 	model.LotteryJiugonggeStarryskyCPic04,
// 	model.LotteryJiugonggeStarryskyCAni01,
// 	model.LotteryJiugonggeStarryskyCAni02,
// 	model.LotteryJiugonggeStarryskyCAni03,
// 	model.LotteryJiugonggeStarryskyCAni04,
// 	model.LotteryJiugonggeStarryskyCAni05,
// 	model.LotteryJiugonggeStarryskyCAni06,

// 	model.LotteryTurntableStarryskyHPic01,
// 	model.LotteryTurntableStarryskyHPic02,
// 	model.LotteryTurntableStarryskyHPic03,
// 	model.LotteryTurntableStarryskyHPic04,
// 	model.LotteryTurntableStarryskyHPic05,
// 	model.LotteryTurntableStarryskyHPic06,
// 	model.LotteryTurntableStarryskyHPic07,
// 	model.LotteryTurntableStarryskyHPic08,
// 	model.LotteryTurntableStarryskyGPic01,
// 	model.LotteryTurntableStarryskyGPic02,
// 	model.LotteryTurntableStarryskyGPic03,
// 	model.LotteryTurntableStarryskyGPic04,
// 	model.LotteryTurntableStarryskyGPic05,
// 	model.LotteryTurntableStarryskyCPic01,
// 	model.LotteryTurntableStarryskyCPic02,
// 	model.LotteryTurntableStarryskyCPic03,
// 	model.LotteryTurntableStarryskyCPic04,
// 	model.LotteryTurntableStarryskyCPic05,
// 	model.LotteryTurntableStarryskyCAni01,
// 	model.LotteryTurntableStarryskyCAni02,
// 	model.LotteryTurntableStarryskyCAni03,
// 	model.LotteryTurntableStarryskyCAni04,
// 	model.LotteryTurntableStarryskyCAni05,
// 	model.LotteryTurntableStarryskyCAni06,
// 	model.LotteryTurntableStarryskyCAni07,

// 	// 音樂
// 	model.LotteryBgmGaming,
// }

// valuesDrawNumbers = []string{
// 	// 搖號抽獎自定義
// 	model.DrawNumbersClassicHPic01,
// 	model.DrawNumbersClassicHPic02,
// 	model.DrawNumbersClassicHPic03,
// 	model.DrawNumbersClassicHPic04,
// 	model.DrawNumbersClassicHPic05,
// 	model.DrawNumbersClassicHPic06,
// 	model.DrawNumbersClassicHPic07,
// 	model.DrawNumbersClassicHPic08,
// 	model.DrawNumbersClassicHPic09,
// 	model.DrawNumbersClassicHPic10,
// 	model.DrawNumbersClassicHPic11,
// 	model.DrawNumbersClassicHPic12,
// 	model.DrawNumbersClassicHPic13,
// 	model.DrawNumbersClassicHPic14,
// 	model.DrawNumbersClassicHPic15,
// 	model.DrawNumbersClassicHPic16,
// 	model.DrawNumbersClassicHAni01,

// 	model.DrawNumbersGoldHPic01,
// 	model.DrawNumbersGoldHPic02,
// 	model.DrawNumbersGoldHPic03,
// 	model.DrawNumbersGoldHPic04,
// 	model.DrawNumbersGoldHPic05,
// 	model.DrawNumbersGoldHPic06,
// 	model.DrawNumbersGoldHPic07,
// 	model.DrawNumbersGoldHPic08,
// 	model.DrawNumbersGoldHPic09,
// 	model.DrawNumbersGoldHPic10,
// 	model.DrawNumbersGoldHPic11,
// 	model.DrawNumbersGoldHPic12,
// 	model.DrawNumbersGoldHPic13,
// 	model.DrawNumbersGoldHPic14,
// 	model.DrawNumbersGoldHAni01,
// 	model.DrawNumbersGoldHAni02,
// 	model.DrawNumbersGoldHAni03,

// 	model.DrawNumbersNewyearDragonHPic01,
// 	model.DrawNumbersNewyearDragonHPic02,
// 	model.DrawNumbersNewyearDragonHPic03,
// 	model.DrawNumbersNewyearDragonHPic04,
// 	model.DrawNumbersNewyearDragonHPic05,
// 	model.DrawNumbersNewyearDragonHPic06,
// 	model.DrawNumbersNewyearDragonHPic07,
// 	model.DrawNumbersNewyearDragonHPic08,
// 	model.DrawNumbersNewyearDragonHPic09,
// 	model.DrawNumbersNewyearDragonHPic10,
// 	model.DrawNumbersNewyearDragonHPic11,
// 	model.DrawNumbersNewyearDragonHPic12,
// 	model.DrawNumbersNewyearDragonHPic13,
// 	model.DrawNumbersNewyearDragonHPic14,
// 	model.DrawNumbersNewyearDragonHPic15,
// 	model.DrawNumbersNewyearDragonHPic16,
// 	model.DrawNumbersNewyearDragonHPic17,
// 	model.DrawNumbersNewyearDragonHPic18,
// 	model.DrawNumbersNewyearDragonHPic19,
// 	model.DrawNumbersNewyearDragonHPic20,
// 	model.DrawNumbersNewyearDragonHAni01,
// 	model.DrawNumbersNewyearDragonHAni02,

// 	model.DrawNumbersCherryHPic01,
// 	model.DrawNumbersCherryHPic02,
// 	model.DrawNumbersCherryHPic03,
// 	model.DrawNumbersCherryHPic04,
// 	model.DrawNumbersCherryHPic05,
// 	model.DrawNumbersCherryHPic06,
// 	model.DrawNumbersCherryHPic07,
// 	model.DrawNumbersCherryHPic08,
// 	model.DrawNumbersCherryHPic09,
// 	model.DrawNumbersCherryHPic10,
// 	model.DrawNumbersCherryHPic11,
// 	model.DrawNumbersCherryHPic12,
// 	model.DrawNumbersCherryHPic13,
// 	model.DrawNumbersCherryHPic14,
// 	model.DrawNumbersCherryHPic15,
// 	model.DrawNumbersCherryHPic16,
// 	model.DrawNumbersCherryHPic17,
// 	model.DrawNumbersCherryHAni01,
// 	model.DrawNumbersCherryHAni02,
// 	model.DrawNumbersCherryHAni03,
// 	model.DrawNumbersCherryHAni04,

// 	model.DrawNumbers3DSpaceHPic01,
// 	model.DrawNumbers3DSpaceHPic02,
// 	model.DrawNumbers3DSpaceHPic03,
// 	model.DrawNumbers3DSpaceHPic04,
// 	model.DrawNumbers3DSpaceHPic05,
// 	model.DrawNumbers3DSpaceHPic06,
// 	model.DrawNumbers3DSpaceHPic07,
// 	model.DrawNumbers3DSpaceHPic08,

// 	// 音樂
// 	model.DrawNumbersBgmGaming,
// }

// valuesWhackMole = []string{
// 	// 敲敲樂自定義
// 	model.WhackmoleClassicHPic01,
// 	model.WhackmoleClassicHPic02,
// 	model.WhackmoleClassicHPic03,
// 	model.WhackmoleClassicHPic04,
// 	model.WhackmoleClassicHPic05,
// 	model.WhackmoleClassicHPic06,
// 	model.WhackmoleClassicHPic07,
// 	model.WhackmoleClassicHPic08,
// 	model.WhackmoleClassicHPic09,
// 	model.WhackmoleClassicHPic10,
// 	model.WhackmoleClassicHPic11,
// 	model.WhackmoleClassicHPic12,
// 	model.WhackmoleClassicHPic13,
// 	model.WhackmoleClassicHPic14,
// 	model.WhackmoleClassicHPic15,
// 	model.WhackmoleClassicGPic01,
// 	model.WhackmoleClassicGPic02,
// 	model.WhackmoleClassicGPic03,
// 	model.WhackmoleClassicGPic04,
// 	model.WhackmoleClassicGPic05,
// 	model.WhackmoleClassicCPic01,
// 	model.WhackmoleClassicCPic02,
// 	model.WhackmoleClassicCPic03,
// 	model.WhackmoleClassicCPic04,
// 	model.WhackmoleClassicCPic05,
// 	model.WhackmoleClassicCPic06,
// 	model.WhackmoleClassicCPic07,
// 	model.WhackmoleClassicCPic08,
// 	model.WhackmoleClassicCAni01,

// 	model.WhackmoleHalloweenHPic01,
// 	model.WhackmoleHalloweenHPic02,
// 	model.WhackmoleHalloweenHPic03,
// 	model.WhackmoleHalloweenHPic04,
// 	model.WhackmoleHalloweenHPic05,
// 	model.WhackmoleHalloweenHPic06,
// 	model.WhackmoleHalloweenHPic07,
// 	model.WhackmoleHalloweenHPic08,
// 	model.WhackmoleHalloweenHPic09,
// 	model.WhackmoleHalloweenHPic10,
// 	model.WhackmoleHalloweenHPic11,
// 	model.WhackmoleHalloweenHPic12,
// 	model.WhackmoleHalloweenHPic13,
// 	model.WhackmoleHalloweenHPic14,
// 	model.WhackmoleHalloweenHPic15,
// 	model.WhackmoleHalloweenGPic01,
// 	model.WhackmoleHalloweenGPic02,
// 	model.WhackmoleHalloweenGPic03,
// 	model.WhackmoleHalloweenGPic04,
// 	model.WhackmoleHalloweenGPic05,
// 	model.WhackmoleHalloweenCPic01,
// 	model.WhackmoleHalloweenCPic02,
// 	model.WhackmoleHalloweenCPic03,
// 	model.WhackmoleHalloweenCPic04,
// 	model.WhackmoleHalloweenCPic05,
// 	model.WhackmoleHalloweenCPic06,
// 	model.WhackmoleHalloweenCPic07,
// 	model.WhackmoleHalloweenCPic08,
// 	model.WhackmoleHalloweenCAni01,

// 	model.WhackmoleChristmasHPic01,
// 	model.WhackmoleChristmasHPic02,
// 	model.WhackmoleChristmasHPic03,
// 	model.WhackmoleChristmasHPic04,
// 	model.WhackmoleChristmasHPic05,
// 	model.WhackmoleChristmasHPic06,
// 	model.WhackmoleChristmasHPic07,
// 	model.WhackmoleChristmasHPic08,
// 	model.WhackmoleChristmasHPic09,
// 	model.WhackmoleChristmasHPic10,
// 	model.WhackmoleChristmasHPic11,
// 	model.WhackmoleChristmasHPic12,
// 	model.WhackmoleChristmasHPic13,
// 	model.WhackmoleChristmasHPic14,
// 	model.WhackmoleChristmasGPic01,
// 	model.WhackmoleChristmasGPic02,
// 	model.WhackmoleChristmasGPic03,
// 	model.WhackmoleChristmasGPic04,
// 	model.WhackmoleChristmasGPic05,
// 	model.WhackmoleChristmasGPic06,
// 	model.WhackmoleChristmasGPic07,
// 	model.WhackmoleChristmasGPic08,
// 	model.WhackmoleChristmasCPic01,
// 	model.WhackmoleChristmasCPic02,
// 	model.WhackmoleChristmasCPic03,
// 	model.WhackmoleChristmasCPic04,
// 	model.WhackmoleChristmasCPic05,
// 	model.WhackmoleChristmasCPic06,
// 	model.WhackmoleChristmasCPic07,
// 	model.WhackmoleChristmasCPic08,
// 	model.WhackmoleChristmasCAni01,
// 	model.WhackmoleChristmasCAni02,

// 	// 音樂
// 	model.WhackmoleBgmStart,
// 	model.WhackmoleBgmGaming,
// 	model.WhackmoleBgmEnd,
// }

// valuesMonopoly = []string{
// 	// 鑑定師自定義
// 	model.MonopolyClassicHPic01,
// 	model.MonopolyClassicHPic02,
// 	model.MonopolyClassicHPic03,
// 	model.MonopolyClassicHPic04,
// 	model.MonopolyClassicHPic05,
// 	model.MonopolyClassicHPic06,
// 	model.MonopolyClassicHPic07,
// 	model.MonopolyClassicHPic08,
// 	model.MonopolyClassicGPic01,
// 	model.MonopolyClassicGPic02,
// 	model.MonopolyClassicGPic03,
// 	model.MonopolyClassicGPic04,
// 	model.MonopolyClassicGPic05,
// 	model.MonopolyClassicGPic06,
// 	model.MonopolyClassicGPic07,
// 	model.MonopolyClassicCPic01,
// 	model.MonopolyClassicCPic02,
// 	model.MonopolyClassicGAni01,
// 	model.MonopolyClassicGAni02,
// 	model.MonopolyClassicCAni01,

// 	model.MonopolyRedpackHPic01,
// 	model.MonopolyRedpackHPic02,
// 	model.MonopolyRedpackHPic03,
// 	model.MonopolyRedpackHPic04,
// 	model.MonopolyRedpackHPic05,
// 	model.MonopolyRedpackHPic06,
// 	model.MonopolyRedpackHPic07,
// 	model.MonopolyRedpackHPic08,
// 	model.MonopolyRedpackHPic09,
// 	model.MonopolyRedpackHPic10,
// 	model.MonopolyRedpackHPic11,
// 	model.MonopolyRedpackHPic12,
// 	model.MonopolyRedpackHPic13,
// 	model.MonopolyRedpackHPic14,
// 	model.MonopolyRedpackHPic15,
// 	model.MonopolyRedpackHPic16,
// 	model.MonopolyRedpackGPic01,
// 	model.MonopolyRedpackGPic02,
// 	model.MonopolyRedpackGPic03,
// 	model.MonopolyRedpackGPic04,
// 	model.MonopolyRedpackGPic05,
// 	model.MonopolyRedpackGPic06,
// 	model.MonopolyRedpackGPic07,
// 	model.MonopolyRedpackGPic08,
// 	model.MonopolyRedpackGPic09,
// 	model.MonopolyRedpackGPic10,
// 	model.MonopolyRedpackCPic01,
// 	model.MonopolyRedpackCPic02,
// 	model.MonopolyRedpackCPic03,
// 	model.MonopolyRedpackHAni01,
// 	model.MonopolyRedpackHAni02,
// 	model.MonopolyRedpackHAni03,
// 	model.MonopolyRedpackGAni01,
// 	model.MonopolyRedpackGAni02,
// 	model.MonopolyRedpackCAni01,

// 	model.MonopolyNewyearRabbitHPic01,
// 	model.MonopolyNewyearRabbitHPic02,
// 	model.MonopolyNewyearRabbitHPic03,
// 	model.MonopolyNewyearRabbitHPic04,
// 	model.MonopolyNewyearRabbitHPic05,
// 	model.MonopolyNewyearRabbitHPic06,
// 	model.MonopolyNewyearRabbitHPic07,
// 	model.MonopolyNewyearRabbitHPic08,
// 	model.MonopolyNewyearRabbitHPic09,
// 	model.MonopolyNewyearRabbitHPic10,
// 	model.MonopolyNewyearRabbitHPic11,
// 	model.MonopolyNewyearRabbitHPic12,
// 	model.MonopolyNewyearRabbitGPic01,
// 	model.MonopolyNewyearRabbitGPic02,
// 	model.MonopolyNewyearRabbitGPic03,
// 	model.MonopolyNewyearRabbitGPic04,
// 	model.MonopolyNewyearRabbitGPic05,
// 	model.MonopolyNewyearRabbitGPic06,
// 	model.MonopolyNewyearRabbitGPic07,
// 	model.MonopolyNewyearRabbitCPic01,
// 	model.MonopolyNewyearRabbitCPic02,
// 	model.MonopolyNewyearRabbitCPic03,
// 	model.MonopolyNewyearRabbitHAni01,
// 	model.MonopolyNewyearRabbitHAni02,
// 	model.MonopolyNewyearRabbitGAni01,
// 	model.MonopolyNewyearRabbitGAni02,
// 	model.MonopolyNewyearRabbitCAni01,

// 	model.MonopolySashimiHPic01,
// 	model.MonopolySashimiHPic02,
// 	model.MonopolySashimiHPic03,
// 	model.MonopolySashimiHPic04,
// 	model.MonopolySashimiHPic05,
// 	model.MonopolySashimiGPic01,
// 	model.MonopolySashimiGPic02,
// 	model.MonopolySashimiGPic03,
// 	model.MonopolySashimiGPic04,
// 	model.MonopolySashimiGPic05,
// 	model.MonopolySashimiGPic06,
// 	model.MonopolySashimiCPic01,
// 	model.MonopolySashimiCPic02,
// 	model.MonopolySashimiHAni01,
// 	model.MonopolySashimiHAni02,
// 	model.MonopolySashimiGAni01,
// 	model.MonopolySashimiGAni02,
// 	model.MonopolySashimiCAni01,

// 	// 音樂
// 	model.MonopolyBgmStart,
// 	model.MonopolyBgmGaming,
// 	model.MonopolyBgmEnd,
// }

// valuesQA = []string{
// 	// 快問快答自定義
// 	model.QAClassicHPic01,
// 	model.QAClassicHPic02,
// 	model.QAClassicHPic03,
// 	model.QAClassicHPic04,
// 	model.QAClassicHPic05,
// 	model.QAClassicHPic06,
// 	model.QAClassicHPic07,
// 	model.QAClassicHPic08,
// 	model.QAClassicHPic09,
// 	model.QAClassicHPic10,
// 	model.QAClassicHPic11,
// 	model.QAClassicHPic12,
// 	model.QAClassicHPic13,
// 	model.QAClassicHPic14,
// 	model.QAClassicHPic15,
// 	model.QAClassicHPic16,
// 	model.QAClassicHPic17,
// 	model.QAClassicHPic18,
// 	model.QAClassicHPic19,
// 	model.QAClassicHPic20,
// 	model.QAClassicHPic21,
// 	model.QAClassicHPic22,
// 	model.QAClassicGPic01,
// 	model.QAClassicGPic02,
// 	model.QAClassicGPic03,
// 	model.QAClassicGPic04,
// 	model.QAClassicGPic05,
// 	model.QAClassicCPic01,
// 	model.QAClassicHAni01,
// 	model.QAClassicHAni02,
// 	model.QAClassicGAni01,
// 	model.QAClassicGAni02,

// 	model.QAElectricHPic01,
// 	model.QAElectricHPic02,
// 	model.QAElectricHPic03,
// 	model.QAElectricHPic04,
// 	model.QAElectricHPic05,
// 	model.QAElectricHPic06,
// 	model.QAElectricHPic07,
// 	model.QAElectricHPic08,
// 	model.QAElectricHPic09,
// 	model.QAElectricHPic10,
// 	model.QAElectricHPic11,
// 	model.QAElectricHPic12,
// 	model.QAElectricHPic13,
// 	model.QAElectricHPic14,
// 	model.QAElectricHPic15,
// 	model.QAElectricHPic16,
// 	model.QAElectricHPic17,
// 	model.QAElectricHPic18,
// 	model.QAElectricHPic19,
// 	model.QAElectricHPic20,
// 	model.QAElectricHPic21,
// 	model.QAElectricHPic22,
// 	model.QAElectricHPic23,
// 	model.QAElectricHPic24,
// 	model.QAElectricHPic25,
// 	model.QAElectricHPic26,
// 	model.QAElectricGPic01,
// 	model.QAElectricGPic02,
// 	model.QAElectricGPic03,
// 	model.QAElectricGPic04,
// 	model.QAElectricGPic05,
// 	model.QAElectricGPic06,
// 	model.QAElectricGPic07,
// 	model.QAElectricGPic08,
// 	model.QAElectricGPic09,
// 	model.QAElectricCPic01,
// 	model.QAElectricHAni01,
// 	model.QAElectricHAni02,
// 	model.QAElectricHAni03,
// 	model.QAElectricHAni04,
// 	model.QAElectricHAni05,
// 	model.QAElectricGAni01,
// 	model.QAElectricGAni02,
// 	model.QAElectricCAni01,

// 	model.QAMoonfestivalHPic01,
// 	model.QAMoonfestivalHPic02,
// 	model.QAMoonfestivalHPic03,
// 	model.QAMoonfestivalHPic04,
// 	model.QAMoonfestivalHPic05,
// 	model.QAMoonfestivalHPic06,
// 	model.QAMoonfestivalHPic07,
// 	model.QAMoonfestivalHPic08,
// 	model.QAMoonfestivalHPic09,
// 	model.QAMoonfestivalHPic10,
// 	model.QAMoonfestivalHPic11,
// 	model.QAMoonfestivalHPic12,
// 	model.QAMoonfestivalHPic13,
// 	model.QAMoonfestivalHPic14,
// 	model.QAMoonfestivalHPic15,
// 	model.QAMoonfestivalHPic16,
// 	model.QAMoonfestivalHPic17,
// 	model.QAMoonfestivalHPic18,
// 	model.QAMoonfestivalHPic19,
// 	model.QAMoonfestivalHPic20,
// 	model.QAMoonfestivalHPic21,
// 	model.QAMoonfestivalHPic22,
// 	model.QAMoonfestivalHPic23,
// 	model.QAMoonfestivalHPic24,
// 	model.QAMoonfestivalGPic01,
// 	model.QAMoonfestivalGPic02,
// 	model.QAMoonfestivalGPic03,
// 	model.QAMoonfestivalGPic04,
// 	model.QAMoonfestivalGPic05,
// 	model.QAMoonfestivalCPic01,
// 	model.QAMoonfestivalCPic02,
// 	model.QAMoonfestivalCPic03,
// 	model.QAMoonfestivalHAni01,
// 	model.QAMoonfestivalHAni02,
// 	model.QAMoonfestivalGAni01,
// 	model.QAMoonfestivalGAni02,
// 	model.QAMoonfestivalGAni03,

// 	// 音樂
// 	model.QABgmStart,
// 	model.QABgmGaming,
// 	model.QABgmEnd,
// }

// valuesQA2 = []string{
// 	// 快問快答自定義
// 	model.QANewyearDragonHPic01,
// 	model.QANewyearDragonHPic02,
// 	model.QANewyearDragonHPic03,
// 	model.QANewyearDragonHPic04,
// 	model.QANewyearDragonHPic05,
// 	model.QANewyearDragonHPic06,
// 	model.QANewyearDragonHPic07,
// 	model.QANewyearDragonHPic08,
// 	model.QANewyearDragonHPic09,
// 	model.QANewyearDragonHPic10,
// 	model.QANewyearDragonHPic11,
// 	model.QANewyearDragonHPic12,
// 	model.QANewyearDragonHPic13,
// 	model.QANewyearDragonHPic14,
// 	model.QANewyearDragonHPic15,
// 	model.QANewyearDragonHPic16,
// 	model.QANewyearDragonHPic17,
// 	model.QANewyearDragonHPic18,
// 	model.QANewyearDragonHPic19,
// 	model.QANewyearDragonHPic20,
// 	model.QANewyearDragonHPic21,
// 	model.QANewyearDragonHPic22,
// 	model.QANewyearDragonHPic23,
// 	model.QANewyearDragonHPic24,
// 	model.QANewyearDragonGPic01,
// 	model.QANewyearDragonGPic02,
// 	model.QANewyearDragonGPic03,
// 	model.QANewyearDragonGPic04,
// 	model.QANewyearDragonGPic05,
// 	model.QANewyearDragonGPic06,
// 	model.QANewyearDragonCPic01,
// 	model.QANewyearDragonHAni01,
// 	model.QANewyearDragonHAni02,
// 	model.QANewyearDragonGAni01,
// 	model.QANewyearDragonGAni02,
// 	model.QANewyearDragonGAni03,
// 	model.QANewyearDragonCAni01,
// }

// valuesTugofwar = []string{
// 	// 拔河遊戲自定義
// 	model.TugofwarClassicHPic01,
// 	model.TugofwarClassicHPic02,
// 	model.TugofwarClassicHPic03,
// 	model.TugofwarClassicHPic04,
// 	model.TugofwarClassicHPic05,
// 	model.TugofwarClassicHPic06,
// 	model.TugofwarClassicHPic07,
// 	model.TugofwarClassicHPic08,
// 	model.TugofwarClassicHPic09,
// 	model.TugofwarClassicHPic10,
// 	model.TugofwarClassicHPic11,
// 	model.TugofwarClassicHPic12,
// 	model.TugofwarClassicHPic13,
// 	model.TugofwarClassicHPic14,
// 	model.TugofwarClassicHPic15,
// 	model.TugofwarClassicHPic16,
// 	model.TugofwarClassicHPic17,
// 	model.TugofwarClassicHPic18,
// 	model.TugofwarClassicHPic19,
// 	model.TugofwarClassicGPic01,
// 	model.TugofwarClassicGPic02,
// 	model.TugofwarClassicGPic03,
// 	model.TugofwarClassicGPic04,
// 	model.TugofwarClassicGPic05,
// 	model.TugofwarClassicGPic06,
// 	model.TugofwarClassicGPic07,
// 	model.TugofwarClassicGPic08,
// 	model.TugofwarClassicGPic09,
// 	model.TugofwarClassicHAni01,
// 	model.TugofwarClassicHAni02,
// 	model.TugofwarClassicHAni03,
// 	model.TugofwarClassicCAni01,

// 	model.TugofwarSchoolHPic01,
// 	model.TugofwarSchoolHPic02,
// 	model.TugofwarSchoolHPic03,
// 	model.TugofwarSchoolHPic04,
// 	model.TugofwarSchoolHPic05,
// 	model.TugofwarSchoolHPic06,
// 	model.TugofwarSchoolHPic07,
// 	model.TugofwarSchoolHPic08,
// 	model.TugofwarSchoolHPic09,
// 	model.TugofwarSchoolHPic10,
// 	model.TugofwarSchoolHPic11,
// 	model.TugofwarSchoolHPic12,
// 	model.TugofwarSchoolHPic13,
// 	model.TugofwarSchoolHPic14,
// 	model.TugofwarSchoolHPic15,
// 	model.TugofwarSchoolHPic16,
// 	model.TugofwarSchoolHPic17,
// 	model.TugofwarSchoolHPic18,
// 	model.TugofwarSchoolHPic19,
// 	model.TugofwarSchoolHPic20,
// 	model.TugofwarSchoolHPic21,
// 	model.TugofwarSchoolHPic22,
// 	model.TugofwarSchoolHPic23,
// 	model.TugofwarSchoolHPic24,
// 	model.TugofwarSchoolHPic25,
// 	model.TugofwarSchoolHPic26,
// 	model.TugofwarSchoolGPic01,
// 	model.TugofwarSchoolGPic02,
// 	model.TugofwarSchoolGPic03,
// 	model.TugofwarSchoolGPic04,
// 	model.TugofwarSchoolGPic05,
// 	model.TugofwarSchoolGPic06,
// 	model.TugofwarSchoolGPic07,
// 	model.TugofwarSchoolCPic01,
// 	model.TugofwarSchoolCPic02,
// 	model.TugofwarSchoolCPic03,
// 	model.TugofwarSchoolCPic04,
// 	model.TugofwarSchoolHAni01,
// 	model.TugofwarSchoolHAni02,
// 	model.TugofwarSchoolHAni03,
// 	model.TugofwarSchoolHAni04,
// 	model.TugofwarSchoolHAni05,
// 	model.TugofwarSchoolHAni06,
// 	model.TugofwarSchoolHAni07,

// 	model.TugofwarChristmasHPic01,
// 	model.TugofwarChristmasHPic02,
// 	model.TugofwarChristmasHPic03,
// 	model.TugofwarChristmasHPic04,
// 	model.TugofwarChristmasHPic05,
// 	model.TugofwarChristmasHPic06,
// 	model.TugofwarChristmasHPic07,
// 	model.TugofwarChristmasHPic08,
// 	model.TugofwarChristmasHPic09,
// 	model.TugofwarChristmasHPic10,
// 	model.TugofwarChristmasHPic11,
// 	model.TugofwarChristmasHPic12,
// 	model.TugofwarChristmasHPic13,
// 	model.TugofwarChristmasHPic14,
// 	model.TugofwarChristmasHPic15,
// 	model.TugofwarChristmasHPic16,
// 	model.TugofwarChristmasHPic17,
// 	model.TugofwarChristmasHPic18,
// 	model.TugofwarChristmasHPic19,
// 	model.TugofwarChristmasHPic20,
// 	model.TugofwarChristmasHPic21,
// 	model.TugofwarChristmasGPic01,
// 	model.TugofwarChristmasGPic02,
// 	model.TugofwarChristmasGPic03,
// 	model.TugofwarChristmasGPic04,
// 	model.TugofwarChristmasGPic05,
// 	model.TugofwarChristmasGPic06,
// 	model.TugofwarChristmasCPic01,
// 	model.TugofwarChristmasCPic02,
// 	model.TugofwarChristmasCPic03,
// 	model.TugofwarChristmasCPic04,
// 	model.TugofwarChristmasHAni01,
// 	model.TugofwarChristmasHAni02,
// 	model.TugofwarChristmasHAni03,
// 	model.TugofwarChristmasCAni01,
// 	model.TugofwarChristmasCAni02,

// 	// 音樂
// 	model.TugofwarBgmStart,  // 遊戲開始
// 	model.TugofwarBgmGaming, // 遊戲進行中
// 	model.TugofwarBgmEnd,    // 遊戲結束
// }

// valuesBingo = []string{
// 	// 賓果遊戲自定義
// 	model.BingoClassicHPic01,
// 	model.BingoClassicHPic02,
// 	model.BingoClassicHPic03,
// 	model.BingoClassicHPic04,
// 	model.BingoClassicHPic05,
// 	model.BingoClassicHPic06,
// 	model.BingoClassicHPic07,
// 	model.BingoClassicHPic08,
// 	model.BingoClassicHPic09,
// 	model.BingoClassicHPic10,
// 	model.BingoClassicHPic11,
// 	model.BingoClassicHPic12,
// 	model.BingoClassicHPic13,
// 	model.BingoClassicHPic14,
// 	model.BingoClassicHPic15,
// 	model.BingoClassicHPic16,
// 	model.BingoClassicGPic01,
// 	model.BingoClassicGPic02,
// 	model.BingoClassicGPic03,
// 	model.BingoClassicGPic04,
// 	model.BingoClassicGPic05,
// 	model.BingoClassicGPic06,
// 	model.BingoClassicGPic07,
// 	model.BingoClassicGPic08,
// 	model.BingoClassicCPic01,
// 	model.BingoClassicCPic02,
// 	model.BingoClassicCPic03,
// 	model.BingoClassicCPic04,
// 	// model.BingoClassicCPic05,
// 	model.BingoClassicHAni01,
// 	model.BingoClassicGAni01,
// 	model.BingoClassicCAni01,
// 	model.BingoClassicCAni02,

// 	model.BingoNewyearDragonHPic01,
// 	model.BingoNewyearDragonHPic02,
// 	model.BingoNewyearDragonHPic03,
// 	model.BingoNewyearDragonHPic04,
// 	model.BingoNewyearDragonHPic05,
// 	model.BingoNewyearDragonHPic06,
// 	model.BingoNewyearDragonHPic07,
// 	model.BingoNewyearDragonHPic08,
// 	model.BingoNewyearDragonHPic09,
// 	model.BingoNewyearDragonHPic10,
// 	model.BingoNewyearDragonHPic11,
// 	model.BingoNewyearDragonHPic12,
// 	model.BingoNewyearDragonHPic13,
// 	model.BingoNewyearDragonHPic14,
// 	// model.BingoNewyearDragonHPic15,
// 	model.BingoNewyearDragonHPic16,
// 	model.BingoNewyearDragonHPic17,
// 	model.BingoNewyearDragonHPic18,
// 	model.BingoNewyearDragonHPic19,
// 	model.BingoNewyearDragonHPic20,
// 	model.BingoNewyearDragonHPic21,
// 	model.BingoNewyearDragonHPic22,
// 	model.BingoNewyearDragonGPic01,
// 	model.BingoNewyearDragonGPic02,
// 	model.BingoNewyearDragonGPic03,
// 	model.BingoNewyearDragonGPic04,
// 	model.BingoNewyearDragonGPic05,
// 	model.BingoNewyearDragonGPic06,
// 	model.BingoNewyearDragonGPic07,
// 	model.BingoNewyearDragonGPic08,
// 	model.BingoNewyearDragonCPic01,
// 	model.BingoNewyearDragonCPic02,
// 	model.BingoNewyearDragonCPic03,
// 	model.BingoNewyearDragonHAni01,
// 	model.BingoNewyearDragonGAni01,
// 	model.BingoNewyearDragonCAni01,
// 	model.BingoNewyearDragonCAni02,
// 	model.BingoNewyearDragonCAni03,

// 	model.BingoCherryHPic01,
// 	model.BingoCherryHPic02,
// 	model.BingoCherryHPic03,
// 	model.BingoCherryHPic04,
// 	model.BingoCherryHPic05,
// 	model.BingoCherryHPic06,
// 	model.BingoCherryHPic07,
// 	model.BingoCherryHPic08,
// 	model.BingoCherryHPic09,
// 	model.BingoCherryHPic10,
// 	model.BingoCherryHPic11,
// 	model.BingoCherryHPic12,
// 	// model.BingoCherryHPic13,
// 	model.BingoCherryHPic14,
// 	model.BingoCherryHPic15,
// 	// model.BingoCherryHPic16,
// 	model.BingoCherryHPic17,
// 	model.BingoCherryHPic18,
// 	model.BingoCherryHPic19,
// 	model.BingoCherryGPic01,
// 	model.BingoCherryGPic02,
// 	model.BingoCherryGPic03,
// 	model.BingoCherryGPic04,
// 	model.BingoCherryGPic05,
// 	model.BingoCherryGPic06,
// 	model.BingoCherryCPic01,
// 	model.BingoCherryCPic02,
// 	model.BingoCherryCPic03,
// 	model.BingoCherryCPic04,
// 	// model.BingoCherryHAni01,
// 	model.BingoCherryHAni02,
// 	model.BingoCherryHAni03,
// 	model.BingoCherryGAni01,
// 	model.BingoCherryGAni02,
// 	model.BingoCherryCAni01,
// 	model.BingoCherryCAni02,

// 	// 音樂
// 	model.BingoBgmStart,
// 	model.BingoBgmGaming,
// 	model.BingoBgmEnd,
// }

// values3DGachaMachine = []string{
// 	// 扭蛋機遊戲自定義
// 	model.GachaMachineClassicHPic02,
// 	model.GachaMachineClassicHPic03,
// 	model.GachaMachineClassicHPic04,
// 	model.GachaMachineClassicHPic05,
// 	model.GachaMachineClassicGPic01,
// 	model.GachaMachineClassicGPic02,
// 	model.GachaMachineClassicGPic03,
// 	model.GachaMachineClassicGPic04,
// 	model.GachaMachineClassicGPic05,
// 	model.GachaMachineClassicCPic01,

// 	// 音樂
// 	model.GachaMachineBgmGaming,
// }

// valuesVote = []string{
// 	// 投票遊戲自定義
// 	model.VoteClassicHPic01,
// 	model.VoteClassicHPic02,
// 	model.VoteClassicHPic03,
// 	model.VoteClassicHPic04,
// 	model.VoteClassicHPic05,
// 	model.VoteClassicHPic06,
// 	model.VoteClassicHPic07,
// 	model.VoteClassicHPic08,
// 	model.VoteClassicHPic09,
// 	model.VoteClassicHPic10,
// 	model.VoteClassicHPic11,
// 	// model.VoteClassicHPic12,
// 	model.VoteClassicHPic13,
// 	model.VoteClassicHPic14,
// 	model.VoteClassicHPic15,
// 	model.VoteClassicHPic16,
// 	model.VoteClassicHPic17,
// 	model.VoteClassicHPic18,
// 	model.VoteClassicHPic19,
// 	model.VoteClassicHPic20,
// 	model.VoteClassicHPic21,
// 	// model.VoteClassicHPic22,
// 	model.VoteClassicHPic23,
// 	model.VoteClassicHPic24,
// 	model.VoteClassicHPic25,
// 	model.VoteClassicHPic26,
// 	model.VoteClassicHPic27,
// 	model.VoteClassicHPic28,
// 	model.VoteClassicHPic29,
// 	model.VoteClassicHPic30,
// 	model.VoteClassicHPic31,
// 	model.VoteClassicHPic32,
// 	model.VoteClassicHPic33,
// 	model.VoteClassicHPic34,
// 	model.VoteClassicHPic35,
// 	model.VoteClassicHPic36,
// 	model.VoteClassicHPic37,
// 	model.VoteClassicGPic01,
// 	model.VoteClassicGPic02,
// 	model.VoteClassicGPic03,
// 	model.VoteClassicGPic04,
// 	model.VoteClassicGPic05,
// 	model.VoteClassicGPic06,
// 	model.VoteClassicGPic07,
// 	model.VoteClassicCPic01,
// 	model.VoteClassicCPic02,
// 	model.VoteClassicCPic03,
// 	model.VoteClassicCPic04,

// 	// 音樂
// 	model.VoteBgmGaming,
// }

// 更新遊戲資料表(activity_game)
// for i, value := range values {
// 	if value != "" {
// 		fieldValues[fields[i]] = value
// 	}
// }
// if len(fieldValues) != 0 {
// 	if err := a.Table(a.Base.TableName).
// 		Where("game_id", "=", model.GameID).
// 		Update(fieldValues); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// }

// 更新遊戲資料表(activity_game_ropepack_picture)
// for i, value2 := range valuesRopepack {
// 	if value2 != "" {
// 		fieldValuesRopepack[fieldsRopepack[i]] = value2
// 	}
// }
// if len(fieldValuesRopepack) != 0 {
// 	if err := a.Table(config.ACTIVITY_GAME_ROPEPACK_PICTURE_TABLE).
// 		Where("game_id", "=", model.GameID).
// 		Update(fieldValuesRopepack); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// }

// 更新遊戲資料表(activity_game_redpack_picture)
// for i, value2 := range valuesRedpack {
// 	if value2 != "" {
// 		fieldValuesRedpack[fieldsRedpack[i]] = value2
// 	}
// }
// if len(fieldValuesRedpack) != 0 {
// 	if err := a.Table(config.ACTIVITY_GAME_REDPACK_PICTURE_TABLE).
// 		Where("game_id", "=", model.GameID).
// 		Update(fieldValuesRedpack); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// }

// 更新遊戲資料表(activity_game_lottery_picture)
// for i, value2 := range valuesLottery {
// 	if value2 != "" {
// 		fieldValuesLottery[fieldsLottery[i]] = value2
// 	}
// }
// if len(fieldValuesLottery) != 0 {
// 	if err := a.Table(config.ACTIVITY_GAME_LOTTERY_PICTURE_TABLE).
// 		Where("game_id", "=", model.GameID).
// 		Update(fieldValuesLottery); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// }

// 更新遊戲資料表(activity_game_draw_numbers_picture)
// for i, value2 := range valuesDrawNumbers {
// 	if value2 != "" {
// 		fieldValuesDrawNumbers[fieldsDrawNumbers[i]] = value2
// 	}
// }
// if len(fieldValuesDrawNumbers) != 0 {
// 	if err := a.Table(config.ACTIVITY_GAME_DRAW_NUMBERS_PICTURE_TABLE).
// 		Where("game_id", "=", model.GameID).
// 		Update(fieldValuesDrawNumbers); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// }

// 更新遊戲資料表(activity_game_whack_mole_picture)
// for i, value2 := range valuesWhackMole {
// 	if value2 != "" {
// 		fieldValuesWhackMole[fieldsWhackMole[i]] = value2
// 	}
// }
// if len(fieldValuesWhackMole) != 0 {
// 	if err := a.Table(config.ACTIVITY_GAME_WHACK_MOLE_PICTURE_TABLE).
// 		Where("game_id", "=", model.GameID).
// 		Update(fieldValuesWhackMole); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// }

// 更新遊戲資料表(activity_game_monopoly_picture)
// for i, value2 := range valuesMonopoly {
// 	if value2 != "" {
// 		fieldValuesMonopoly[fieldsMonopoly[i]] = value2
// 	}
// }
// if len(fieldValuesMonopoly) != 0 {
// 	if err := a.Table(config.ACTIVITY_GAME_MONOPOLY_PICTURE_TABLE).
// 		Where("game_id", "=", model.GameID).
// 		Update(fieldValuesMonopoly); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// }

// 更新遊戲資料表(activity_game_qa_picture)
// for i, value2 := range valuesQA {
// 	if value2 != "" {
// 		fieldValuesQA[fieldsQA[i]] = value2
// 	}
// }
// if len(fieldValuesQA) != 0 {
// 	if err := a.Table(config.ACTIVITY_GAME_QA_PICTURE_TABLE_1).
// 		Where("game_id", "=", model.GameID).
// 		Update(fieldValuesQA); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// }

// 更新遊戲資料表(activity_game_qa_picture_2)
// for i, value2 := range valuesQA2 {
// 	if value2 != "" {
// 		fieldValuesQA2[fieldsQA2[i]] = value2
// 	}
// }
// if len(fieldValuesQA2) != 0 {
// 	if err := a.Table(config.ACTIVITY_GAME_QA_PICTURE_TABLE_2).
// 		Where("game_id", "=", model.GameID).
// 		Update(fieldValuesQA2); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// }

// 更新遊戲資料表(activity_game_tugofwar_picture)
// for i, value2 := range valuesTugofwar {
// 	if value2 != "" {
// 		fieldValuesTugofwar[fieldsTugofwar[i]] = value2
// 	}
// }
// if len(fieldValuesTugofwar) != 0 {
// 	if err := a.Table(config.ACTIVITY_GAME_TUGOFWAR_PICTURE_TABLE).
// 		Where("game_id", "=", model.GameID).
// 		Update(fieldValuesTugofwar); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// }

// 更新遊戲資料表(activity_game_bingo_picture)
// for i, value2 := range valuesBingo {
// 	if value2 != "" {
// 		fieldValuesBingo[fieldsBingo[i]] = value2
// 	}
// }
// if len(fieldValuesBingo) != 0 {
// 	if err := a.Table(config.ACTIVITY_GAME_BINGO_PICTURE_TABLE).
// 		Where("game_id", "=", model.GameID).
// 		Update(fieldValuesBingo); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// }

// 更新遊戲資料表(activity_game_3d_gacha_machine_picture)
// for i, value2 := range values3DGachaMachine {
// 	if value2 != "" {
// 		fieldValues3DGachaMachine[fields3DGachaMachine[i]] = value2
// 	}
// }
// if len(fieldValues3DGachaMachine) != 0 {
// 	if err := a.Table(config.ACTIVITY_GAME_GACHAMACHINE_PICTURE_TABLE).
// 		Where("game_id", "=", model.GameID).
// 		Update(fieldValues3DGachaMachine); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// }

// 更新遊戲資料表(activity_game_vote_picture)
// for i, value2 := range valuesVote {
// 	if value2 != "" {
// 		fieldValuesVote[fieldsVote[i]] = value2
// 	}
// }
// if len(fieldValuesVote) != 0 {
// 	if err := a.Table(config.ACTIVITY_GAME_VOTE_PICTURE_TABLE).
// 		Where("game_id", "=", model.GameID).
// 		Update(fieldValuesVote); err != nil &&
// 		err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// }

// gameStatus string
// fieldValuesRopepack       = command.Value{}
// fieldValuesRedpack        = command.Value{}
// fieldValuesLottery        = command.Value{}
// fieldValuesDrawNumbers    = command.Value{}
// fieldValuesWhackMole      = command.Value{}
// fieldValuesMonopoly       = command.Value{}
// fieldValuesQA             = command.Value{}
// fieldValuesQA2            = command.Value{}
// fieldValuesTugofwar       = command.Value{}
// fieldValuesBingo          = command.Value{}
// fieldValues3DGachaMachine = command.Value{}
// fieldValuesVote           = command.Value{}

// if model.Skin != "" && (game == "whack_mole") {
// 	var (
// 		skins = []string{"rat", "ox", "tiger", "rabbit", "dragon", "snake",
// 			"horse", "goat", "monkey", "rooster", "dog", "pig"}
// 		isbool bool
// 	)

// 	for i := range skins {
// 		if model.Skin == skins[i] {
// 			isbool = true
// 			break
// 		}
// 	}
// 	if isbool == false {
// 		return errors.New("錯誤: skin欄位發生問題，請輸入有效的十二生肖")
// 	}
// }

// command.Value{
// 	"qa_1":         model.QA1,
// 	"qa_1_options": model.QA1Options,
// 	"qa_1_answer":  model.QA1Answer,
// 	"qa_1_score":   utils.GetInt64(model.QA1Score, 0),

// 	"qa_2":         model.QA2,
// 	"qa_2_options": model.QA2Options,
// 	"qa_2_answer":  model.QA2Answer,
// 	"qa_2_score":   utils.GetInt64(model.QA2Score, 0),

// 	"qa_3":         model.QA3,
// 	"qa_3_options": model.QA3Options,
// 	"qa_3_answer":  model.QA3Answer,
// 	"qa_3_score":   utils.GetInt64(model.QA3Score, 0),

// 	"qa_4":         model.QA4,
// 	"qa_4_options": model.QA4Options,
// 	"qa_4_answer":  model.QA4Answer,
// 	"qa_4_score":   utils.GetInt64(model.QA4Score, 0),

// 	"qa_5":         model.QA5,
// 	"qa_5_options": model.QA5Options,
// 	"qa_5_answer":  model.QA5Answer,
// 	"qa_5_score":   utils.GetInt64(model.QA5Score, 0),

// 	"qa_6":         model.QA6,
// 	"qa_6_options": model.QA6Options,
// 	"qa_6_answer":  model.QA6Answer,
// 	"qa_6_score":   utils.GetInt64(model.QA6Score, 0),

// 	"qa_7":         model.QA7,
// 	"qa_7_options": model.QA7Options,
// 	"qa_7_answer":  model.QA7Answer,
// 	"qa_7_score":   utils.GetInt64(model.QA7Score, 0),

// 	"qa_8":         model.QA8,
// 	"qa_8_options": model.QA8Options,
// 	"qa_8_answer":  model.QA8Answer,
// 	"qa_8_score":   utils.GetInt64(model.QA8Score, 0),

// 	"qa_9":         model.QA9,
// 	"qa_9_options": model.QA9Options,
// 	"qa_9_answer":  model.QA9Answer,
// 	"qa_9_score":   utils.GetInt64(model.QA9Score, 0),

// 	"qa_10":         model.QA10,
// 	"qa_10_options": model.QA10Options,
// 	"qa_10_answer":  model.QA10Answer,
// 	"qa_10_score":   utils.GetInt64(model.QA10Score, 0),

// 	"qa_11":         model.QA11,
// 	"qa_11_options": model.QA11Options,
// 	"qa_11_answer":  model.QA11Answer,
// 	"qa_11_score":   utils.GetInt64(model.QA11Score, 0),

// 	"qa_12":         model.QA12,
// 	"qa_12_options": model.QA12Options,
// 	"qa_12_answer":  model.QA12Answer,
// 	"qa_12_score":   utils.GetInt64(model.QA12Score, 0),

// 	"qa_13":         model.QA13,
// 	"qa_13_options": model.QA13Options,
// 	"qa_13_answer":  model.QA13Answer,
// 	"qa_13_score":   utils.GetInt64(model.QA13Score, 0),

// 	"qa_14":         model.QA14,
// 	"qa_14_options": model.QA14Options,
// 	"qa_14_answer":  model.QA14Answer,
// 	"qa_14_score":   utils.GetInt64(model.QA14Score, 0),

// 	"qa_15":         model.QA15,
// 	"qa_15_options": model.QA15Options,
// 	"qa_15_answer":  model.QA15Answer,
// 	"qa_15_score":   utils.GetInt64(model.QA15Score, 0),

// 	"qa_16":         model.QA16,
// 	"qa_16_options": model.QA16Options,
// 	"qa_16_answer":  model.QA16Answer,
// 	"qa_16_score":   utils.GetInt64(model.QA16Score, 0),

// 	"qa_17":         model.QA17,
// 	"qa_17_options": model.QA17Options,
// 	"qa_17_answer":  model.QA17Answer,
// 	"qa_17_score":   utils.GetInt64(model.QA17Score, 0),

// 	"qa_18":         model.QA18,
// 	"qa_18_options": model.QA18Options,
// 	"qa_18_answer":  model.QA18Answer,
// 	"qa_18_score":   utils.GetInt64(model.QA18Score, 0),

// 	"qa_19":         model.QA19,
// 	"qa_19_options": model.QA19Options,
// 	"qa_19_answer":  model.QA19Answer,
// 	"qa_19_score":   utils.GetInt64(model.QA19Score, 0),

// 	"qa_20":         model.QA20,
// 	"qa_20_options": model.QA20Options,
// 	"qa_20_answer":  model.QA20Answer,
// 	"qa_20_score":   utils.GetInt64(model.QA20Score, 0),

// 	"total_qa":  model.TotalQA,
// 	"qa_second": model.QASecond,
// 	"qa_round":  1, // 更新題目時，回到第一題
// }

// if isRedis {
// 清除遊戲redis資訊(並重新開啟遊戲頁面)
// a.RedisConn.DelCache(config.GAME_REDIS + model.GameID)                            // 遊戲設置
// a.RedisConn.DelCache(config.GAME_TYPE_REDIS + model.GameID)                       // 遊戲種類
// a.RedisConn.DelCache(config.SCORES_REDIS + model.GameID)                          // 分數
// a.RedisConn.DelCache(config.SCORES_2_REDIS + model.GameID)                        // 第二分數
// a.RedisConn.DelCache(config.CORRECT_REDIS + model.GameID)                         // 答對題數
// a.RedisConn.DelCache(config.QA_REDIS + model.GameID)                              // 快問快答題目資訊
// a.RedisConn.DelCache(config.QA_RECORD_REDIS + model.GameID)                       // 快問快答答題紀錄
// a.RedisConn.DelCache(config.WINNING_STAFFS_REDIS + model.GameID)                  // 中獎人員
// a.RedisConn.DelCache(config.NO_WINNING_STAFFS_REDIS + model.GameID)               // 未中獎人員
// a.RedisConn.DelCache(config.GAME_TEAM_REDIS + model.GameID)                       // 遊戲隊伍資訊，HASH
// a.RedisConn.DelCache(config.GAME_BINGO_NUMBER_REDIS + model.GameID)               // 紀錄抽過的號碼，LIST
// a.RedisConn.DelCache(config.GAME_BINGO_USER_REDIS + model.GameID)                 // 賓果中獎人員，ZSET
// a.RedisConn.DelCache(config.GAME_BINGO_USER_NUMBER + model.GameID)                // 紀錄玩家的號碼排序，HASH
// a.RedisConn.DelCache(config.GAME_BINGO_USER_GOING_BINGO + model.GameID)           // 紀錄玩家是否即將賓果，HASH
// a.RedisConn.DelCache(config.GAME_ATTEND_REDIS + model.GameID)                     // 遊戲人數資訊，SET
// a.RedisConn.DelCache(config.GAME_TUGOFWAR_LEFT_TEAM_ATTEND_REDIS + model.GameID)  // 拔河遊戲左隊人數資訊，SET
// a.RedisConn.DelCache(config.GAME_TUGOFWAR_RIGHT_TEAM_ATTEND_REDIS + model.GameID) // 拔河遊戲右隊人數資訊，SET

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GAME_REDIS+model.GameID, "修改資料")

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GUEST_GAME_STATUS_REDIS+model.GameID, "修改資料")

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GAME_BINGO_NUMBER_REDIS+model.GameID, "修改資料")

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_QA_REDIS+model.GameID, "修改資料")

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GAME_TEAM_REDIS+model.GameID, "修改資料")

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GAME_EDIT_TIMES_REDIS+model.GameID, "修改資料")

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_WINNING_STAFFS_REDIS+model.GameID, "修改資料")

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_GAME_BINGO_USER_NUMBER+model.GameID, "修改資料")

// 當接收到 WebSocket 消息時，將其發布到 Redis 頻道
// a.RedisConn.Publish(config.CHANNEL_SCORES_REDIS+model.GameID, "修改資料")
// }

// if err := a.Table(a.Base.TableName).
// 	Where("game_id", "=", model.GameID).
// 	Update(fieldValues); err != nil &&
// 	err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 	return err
// }

// 其他自定義圖片資料表
// for _, updateInfo := range updateInfos {
// 	for _, field := range updateInfo.Fields {
// 		if val, ok := data[field]; ok && val != "" {
// 			fieldValues[field] = val
// 		}
// 	}

// 檢查每個欄位是否有值並更新
// customizeFieldValues := make(map[string]interface{})

// 如果有需要更新的資料，則執行更新
// if len(customizeFieldValues) > 0 {
// 	// log.Println("更新的資料表: ", updateInfo.TableName)

// 	if err := a.Table(updateInfo.TableName).
// 		Where("game_id", "=", model.GameID).
// 		Update(customizeFieldValues); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 		return err
// 	}
// }
// }

// 快問快答
// if game == "QA" {
// if model.QA1 == "" || model.QA1Options == "" ||
// 	model.QA1Answer == "" || model.QA1Score == "" {
// 	return errors.New("錯誤: 題目設置最少一題，請重新設置")
// }

// 題目設置資訊
// qas := []string{
// 	"total_qa",
// 	"qa_second",
// 	"qa_round",
// }

// 手動處理
// data["qa_round"] = 1 // 更新題目時，回到第一題

// 迴圈處理題目資料並將對應的參數加入 qas 陣列
// for i := 1; i <= 20; i++ {
// 	qaFieldPrefix := fmt.Sprintf("qa_%d", i)
// 	qaOptionsFieldPrefix := fmt.Sprintf("qa_%d_options", i)
// 	qaAnswerFieldPrefix := fmt.Sprintf("qa_%d_answer", i)
// 	qaScoreFieldPrefix := fmt.Sprintf("qa_%d_score", i)

// 	// 把這些動態生成的欄位名稱加入 qas 陣列
// 	qas = append(qas, qaFieldPrefix, qaOptionsFieldPrefix, qaAnswerFieldPrefix, qaScoreFieldPrefix)

// 	// 手動處理qa_1_score資料(轉為int64)
// 	data[fmt.Sprintf("qa_%d_score", i)] = utils.GetInt64(data[fmt.Sprintf("qa_%d_score", i)], 0)

// }

// activity_game_qa資料表
// if err := a.Table(config.ACTIVITY_GAME_QA_TABLE).
// 	Where("game_id", "=", model.GameID).
// 	Update(FilterFields(data, qas)); err != nil && err.Error() != "錯誤: 無更新任何資料，請重新操作" {
// 	return err
// }
// }

// fields = []string{
// 	"title", "game_type", "limit_time",
// 	"max_times", "allow", "percent", "first_prize",
// 	"second_prize", "third_prize", "general_prize", "topic",
// 	"skin", "music", "display_name",

// 	"allow_choose_team",
// 	"left_team_name",
// 	"left_team_picture",
// 	"right_team_name",
// 	"right_team_picture",
// 	"prize",

// 	"max_number",
// 	"bingo_line",
// 	"round_prize",

// 	"gacha_machine_reflection",
// 	"reflective_switch",

// 	// "game_order", // 場次排序

// 	"box_reflection",
// 	"same_people",

// 	"vote_screen",
// 	"vote_times",
// 	"vote_method",
// 	"vote_method_player",
// 	"vote_restriction",
// 	"avatar_shape",
// 	"auto_play",
// 	"show_rank",
// 	"title_switch",
// 	"arrangement_guest",

// 	"qa_1", "qa_1_options", "qa_1_answer", "qa_1_score",
// 	"qa_2", "qa_2_options", "qa_2_answer", "qa_2_score",
// 	"qa_3", "qa_3_options", "qa_3_answer", "qa_3_score",
// 	"qa_4", "qa_4_options", "qa_4_answer", "qa_4_score",
// 	"qa_5", "qa_5_options", "qa_5_answer", "qa_5_score",
// 	"qa_6", "qa_6_options", "qa_6_answer", "qa_6_score",
// 	"qa_7", "qa_7_options", "qa_7_answer", "qa_7_score",
// 	"qa_8", "qa_8_options", "qa_8_answer", "qa_8_score",
// 	"qa_9", "qa_9_options", "qa_9_answer", "qa_9_score",
// 	"qa_10", "qa_10_options", "qa_10_answer", "qa_10_score",
// 	"qa_11", "qa_11_options", "qa_11_answer", "qa_11_score",
// 	"qa_12", "qa_12_options", "qa_12_answer", "qa_12_score",
// 	"qa_13", "qa_13_options", "qa_13_answer", "qa_13_score",
// 	"qa_14", "qa_14_options", "qa_14_answer", "qa_14_score",
// 	"qa_15", "qa_15_options", "qa_15_answer", "qa_15_score",
// 	"qa_16", "qa_16_options", "qa_16_answer", "qa_16_score",
// 	"qa_17", "qa_17_options", "qa_17_answer", "qa_17_score",
// 	"qa_18", "qa_18_options", "qa_18_answer", "qa_18_score",
// 	"qa_19", "qa_19_options", "qa_19_answer", "qa_19_score",
// 	"qa_20", "qa_20_options", "qa_20_answer", "qa_20_score",

// 	"qa_round", "qa_second", "total_qa",
// }

// // 扭蛋機-----start
// fields3DGachaMachine = []string{
// 	// 扭蛋機遊戲自定義
// 	"3d_gacha_machine_classic_h_pic_02",
// 	"3d_gacha_machine_classic_h_pic_03",
// 	"3d_gacha_machine_classic_h_pic_04",
// 	"3d_gacha_machine_classic_h_pic_05",
// 	"3d_gacha_machine_classic_g_pic_01",
// 	"3d_gacha_machine_classic_g_pic_02",
// 	"3d_gacha_machine_classic_g_pic_03",
// 	"3d_gacha_machine_classic_g_pic_04",
// 	"3d_gacha_machine_classic_g_pic_05",
// 	"3d_gacha_machine_classic_c_pic_01",

// 	// 音樂
// 	"3d_gacha_machine_bgm_gaming",
// }

// // 扭蛋機-----end

// // 賓果-----start
// fieldsBingo = []string{
// 	// 賓果遊戲自定義
// 	"bingo_classic_h_pic_01",
// 	"bingo_classic_h_pic_02",
// 	"bingo_classic_h_pic_03",
// 	"bingo_classic_h_pic_04",
// 	"bingo_classic_h_pic_05",
// 	"bingo_classic_h_pic_06",
// 	"bingo_classic_h_pic_07",
// 	"bingo_classic_h_pic_08",
// 	"bingo_classic_h_pic_09",
// 	"bingo_classic_h_pic_10",
// 	"bingo_classic_h_pic_11",
// 	"bingo_classic_h_pic_12",
// 	"bingo_classic_h_pic_13",
// 	"bingo_classic_h_pic_14",
// 	"bingo_classic_h_pic_15",
// 	"bingo_classic_h_pic_16",
// 	"bingo_classic_g_pic_01",
// 	"bingo_classic_g_pic_02",
// 	"bingo_classic_g_pic_03",
// 	"bingo_classic_g_pic_04",
// 	"bingo_classic_g_pic_05",
// 	"bingo_classic_g_pic_06",
// 	"bingo_classic_g_pic_07",
// 	"bingo_classic_g_pic_08",
// 	"bingo_classic_c_pic_01",
// 	"bingo_classic_c_pic_02",
// 	"bingo_classic_c_pic_03",
// 	"bingo_classic_c_pic_04",
// 	"bingo_classic_h_ani_01",
// 	"bingo_classic_g_ani_01",
// 	"bingo_classic_c_ani_01",
// 	"bingo_classic_c_ani_02",

// 	"bingo_newyear_dragon_h_pic_01",
// 	"bingo_newyear_dragon_h_pic_02",
// 	"bingo_newyear_dragon_h_pic_03",
// 	"bingo_newyear_dragon_h_pic_04",
// 	"bingo_newyear_dragon_h_pic_05",
// 	"bingo_newyear_dragon_h_pic_06",
// 	"bingo_newyear_dragon_h_pic_07",
// 	"bingo_newyear_dragon_h_pic_08",
// 	"bingo_newyear_dragon_h_pic_09",
// 	"bingo_newyear_dragon_h_pic_10",
// 	"bingo_newyear_dragon_h_pic_11",
// 	"bingo_newyear_dragon_h_pic_12",
// 	"bingo_newyear_dragon_h_pic_13",
// 	"bingo_newyear_dragon_h_pic_14",
// 	"bingo_newyear_dragon_h_pic_16",
// 	"bingo_newyear_dragon_h_pic_17",
// 	"bingo_newyear_dragon_h_pic_18",
// 	"bingo_newyear_dragon_h_pic_19",
// 	"bingo_newyear_dragon_h_pic_20",
// 	"bingo_newyear_dragon_h_pic_21",
// 	"bingo_newyear_dragon_h_pic_22",
// 	"bingo_newyear_dragon_g_pic_01",
// 	"bingo_newyear_dragon_g_pic_02",
// 	"bingo_newyear_dragon_g_pic_03",
// 	"bingo_newyear_dragon_g_pic_04",
// 	"bingo_newyear_dragon_g_pic_05",
// 	"bingo_newyear_dragon_g_pic_06",
// 	"bingo_newyear_dragon_g_pic_07",
// 	"bingo_newyear_dragon_g_pic_08",
// 	"bingo_newyear_dragon_c_pic_01",
// 	"bingo_newyear_dragon_c_pic_02",
// 	"bingo_newyear_dragon_c_pic_03",
// 	"bingo_newyear_dragon_h_ani_01",
// 	"bingo_newyear_dragon_g_ani_01",
// 	"bingo_newyear_dragon_c_ani_01",
// 	"bingo_newyear_dragon_c_ani_02",
// 	"bingo_newyear_dragon_c_ani_03",

// 	"bingo_cherry_h_pic_01",
// 	"bingo_cherry_h_pic_02",
// 	"bingo_cherry_h_pic_03",
// 	"bingo_cherry_h_pic_04",
// 	"bingo_cherry_h_pic_05",
// 	"bingo_cherry_h_pic_06",
// 	"bingo_cherry_h_pic_07",
// 	"bingo_cherry_h_pic_08",
// 	"bingo_cherry_h_pic_09",
// 	"bingo_cherry_h_pic_10",
// 	"bingo_cherry_h_pic_11",
// 	"bingo_cherry_h_pic_12",
// 	"bingo_cherry_h_pic_14",
// 	"bingo_cherry_h_pic_15",
// 	"bingo_cherry_h_pic_17",
// 	"bingo_cherry_h_pic_18",
// 	"bingo_cherry_h_pic_19",
// 	"bingo_cherry_g_pic_01",
// 	"bingo_cherry_g_pic_02",
// 	"bingo_cherry_g_pic_03",
// 	"bingo_cherry_g_pic_04",
// 	"bingo_cherry_g_pic_05",
// 	"bingo_cherry_g_pic_06",
// 	"bingo_cherry_c_pic_01",
// 	"bingo_cherry_c_pic_02",
// 	"bingo_cherry_c_pic_03",
// 	"bingo_cherry_c_pic_04",
// 	"bingo_cherry_h_ani_02",
// 	"bingo_cherry_h_ani_03",
// 	"bingo_cherry_g_ani_01",
// 	"bingo_cherry_g_ani_02",
// 	"bingo_cherry_c_ani_01",
// 	"bingo_cherry_c_ani_02",

// 	// 音樂
// 	"bingo_bgm_start",
// 	"bingo_bgm_gaming",
// 	"bingo_bgm_end",
// }
// // 賓果-----end

// // 搖號抽獎-----start
// fieldsDrawNumbers = []string{
// 	// 搖號抽獎自定義
// 	"draw_numbers_classic_h_pic_01",
// 	"draw_numbers_classic_h_pic_02",
// 	"draw_numbers_classic_h_pic_03",
// 	"draw_numbers_classic_h_pic_04",
// 	"draw_numbers_classic_h_pic_05",
// 	"draw_numbers_classic_h_pic_06",
// 	"draw_numbers_classic_h_pic_07",
// 	"draw_numbers_classic_h_pic_08",
// 	"draw_numbers_classic_h_pic_09",
// 	"draw_numbers_classic_h_pic_10",
// 	"draw_numbers_classic_h_pic_11",
// 	"draw_numbers_classic_h_pic_12",
// 	"draw_numbers_classic_h_pic_13",
// 	"draw_numbers_classic_h_pic_14",
// 	"draw_numbers_classic_h_pic_15",
// 	"draw_numbers_classic_h_pic_16",
// 	"draw_numbers_classic_h_ani_01",

// 	"draw_numbers_gold_h_pic_01",
// 	"draw_numbers_gold_h_pic_02",
// 	"draw_numbers_gold_h_pic_03",
// 	"draw_numbers_gold_h_pic_04",
// 	"draw_numbers_gold_h_pic_05",
// 	"draw_numbers_gold_h_pic_06",
// 	"draw_numbers_gold_h_pic_07",
// 	"draw_numbers_gold_h_pic_08",
// 	"draw_numbers_gold_h_pic_09",
// 	"draw_numbers_gold_h_pic_10",
// 	"draw_numbers_gold_h_pic_11",
// 	"draw_numbers_gold_h_pic_12",
// 	"draw_numbers_gold_h_pic_13",
// 	"draw_numbers_gold_h_pic_14",
// 	"draw_numbers_gold_h_ani_01",
// 	"draw_numbers_gold_h_ani_02",
// 	"draw_numbers_gold_h_ani_03",

// 	"draw_numbers_newyear_dragon_h_pic_01",
// 	"draw_numbers_newyear_dragon_h_pic_02",
// 	"draw_numbers_newyear_dragon_h_pic_03",
// 	"draw_numbers_newyear_dragon_h_pic_04",
// 	"draw_numbers_newyear_dragon_h_pic_05",
// 	"draw_numbers_newyear_dragon_h_pic_06",
// 	"draw_numbers_newyear_dragon_h_pic_07",
// 	"draw_numbers_newyear_dragon_h_pic_08",
// 	"draw_numbers_newyear_dragon_h_pic_09",
// 	"draw_numbers_newyear_dragon_h_pic_10",
// 	"draw_numbers_newyear_dragon_h_pic_11",
// 	"draw_numbers_newyear_dragon_h_pic_12",
// 	"draw_numbers_newyear_dragon_h_pic_13",
// 	"draw_numbers_newyear_dragon_h_pic_14",
// 	"draw_numbers_newyear_dragon_h_pic_15",
// 	"draw_numbers_newyear_dragon_h_pic_16",
// 	"draw_numbers_newyear_dragon_h_pic_17",
// 	"draw_numbers_newyear_dragon_h_pic_18",
// 	"draw_numbers_newyear_dragon_h_pic_19",
// 	"draw_numbers_newyear_dragon_h_pic_20",
// 	"draw_numbers_newyear_dragon_h_ani_01",
// 	"draw_numbers_newyear_dragon_h_ani_02",

// 	"draw_numbers_cherry_h_pic_01",
// 	"draw_numbers_cherry_h_pic_02",
// 	"draw_numbers_cherry_h_pic_03",
// 	"draw_numbers_cherry_h_pic_04",
// 	"draw_numbers_cherry_h_pic_05",
// 	"draw_numbers_cherry_h_pic_06",
// 	"draw_numbers_cherry_h_pic_07",
// 	"draw_numbers_cherry_h_pic_08",
// 	"draw_numbers_cherry_h_pic_09",
// 	"draw_numbers_cherry_h_pic_10",
// 	"draw_numbers_cherry_h_pic_11",
// 	"draw_numbers_cherry_h_pic_12",
// 	"draw_numbers_cherry_h_pic_13",
// 	"draw_numbers_cherry_h_pic_14",
// 	"draw_numbers_cherry_h_pic_15",
// 	"draw_numbers_cherry_h_pic_16",
// 	"draw_numbers_cherry_h_pic_17",
// 	"draw_numbers_cherry_h_ani_01",
// 	"draw_numbers_cherry_h_ani_02",
// 	"draw_numbers_cherry_h_ani_03",
// 	"draw_numbers_cherry_h_ani_04",

// 	"draw_numbers_3D_space_h_pic_01",
// 	"draw_numbers_3D_space_h_pic_02",
// 	"draw_numbers_3D_space_h_pic_03",
// 	"draw_numbers_3D_space_h_pic_04",
// 	"draw_numbers_3D_space_h_pic_05",
// 	"draw_numbers_3D_space_h_pic_06",
// 	"draw_numbers_3D_space_h_pic_07",
// 	"draw_numbers_3D_space_h_pic_08",

// 	// 音樂
// 	"draw_numbers_bgm_gaming", // 遊戲進行中

// }
// // 搖號抽獎-----end

// // 遊戲抽獎-----start
// fieldsLottery = []string{
// 	// 遊戲抽獎自定義
// 	"lottery_jiugongge_classic_h_pic_01",
// 	"lottery_jiugongge_classic_h_pic_02",
// 	"lottery_jiugongge_classic_h_pic_03",
// 	"lottery_jiugongge_classic_h_pic_04",
// 	"lottery_jiugongge_classic_g_pic_01",
// 	"lottery_jiugongge_classic_g_pic_02",
// 	"lottery_jiugongge_classic_c_pic_01",
// 	"lottery_jiugongge_classic_c_pic_02",
// 	"lottery_jiugongge_classic_c_pic_03",
// 	"lottery_jiugongge_classic_c_pic_04",
// 	"lottery_jiugongge_classic_c_ani_01",
// 	"lottery_jiugongge_classic_c_ani_02",
// 	"lottery_jiugongge_classic_c_ani_03",

// 	"lottery_turntable_classic_h_pic_01",
// 	"lottery_turntable_classic_h_pic_02",
// 	"lottery_turntable_classic_h_pic_03",
// 	"lottery_turntable_classic_h_pic_04",
// 	"lottery_turntable_classic_g_pic_01",
// 	"lottery_turntable_classic_g_pic_02",
// 	"lottery_turntable_classic_c_pic_01",
// 	"lottery_turntable_classic_c_pic_02",
// 	"lottery_turntable_classic_c_pic_03",
// 	"lottery_turntable_classic_c_pic_04",
// 	"lottery_turntable_classic_c_pic_05",
// 	"lottery_turntable_classic_c_pic_06",
// 	"lottery_turntable_classic_c_ani_01",
// 	"lottery_turntable_classic_c_ani_02",
// 	"lottery_turntable_classic_c_ani_03",

// 	"lottery_jiugongge_starrysky_h_pic_01",
// 	"lottery_jiugongge_starrysky_h_pic_02",
// 	"lottery_jiugongge_starrysky_h_pic_03",
// 	"lottery_jiugongge_starrysky_h_pic_04",
// 	"lottery_jiugongge_starrysky_h_pic_05",
// 	"lottery_jiugongge_starrysky_h_pic_06",
// 	"lottery_jiugongge_starrysky_h_pic_07",
// 	"lottery_jiugongge_starrysky_g_pic_01",
// 	"lottery_jiugongge_starrysky_g_pic_02",
// 	"lottery_jiugongge_starrysky_g_pic_03",
// 	"lottery_jiugongge_starrysky_g_pic_04",
// 	"lottery_jiugongge_starrysky_c_pic_01",
// 	"lottery_jiugongge_starrysky_c_pic_02",
// 	"lottery_jiugongge_starrysky_c_pic_03",
// 	"lottery_jiugongge_starrysky_c_pic_04",
// 	"lottery_jiugongge_starrysky_c_ani_01",
// 	"lottery_jiugongge_starrysky_c_ani_02",
// 	"lottery_jiugongge_starrysky_c_ani_03",
// 	"lottery_jiugongge_starrysky_c_ani_04",
// 	"lottery_jiugongge_starrysky_c_ani_05",
// 	"lottery_jiugongge_starrysky_c_ani_06",

// 	"lottery_turntable_starrysky_h_pic_01",
// 	"lottery_turntable_starrysky_h_pic_02",
// 	"lottery_turntable_starrysky_h_pic_03",
// 	"lottery_turntable_starrysky_h_pic_04",
// 	"lottery_turntable_starrysky_h_pic_05",
// 	"lottery_turntable_starrysky_h_pic_06",
// 	"lottery_turntable_starrysky_h_pic_07",
// 	"lottery_turntable_starrysky_h_pic_08",
// 	"lottery_turntable_starrysky_g_pic_01",
// 	"lottery_turntable_starrysky_g_pic_02",
// 	"lottery_turntable_starrysky_g_pic_03",
// 	"lottery_turntable_starrysky_g_pic_04",
// 	"lottery_turntable_starrysky_g_pic_05",
// 	"lottery_turntable_starrysky_c_pic_01",
// 	"lottery_turntable_starrysky_c_pic_02",
// 	"lottery_turntable_starrysky_c_pic_03",
// 	"lottery_turntable_starrysky_c_pic_04",
// 	"lottery_turntable_starrysky_c_pic_05",
// 	"lottery_turntable_starrysky_c_ani_01",
// 	"lottery_turntable_starrysky_c_ani_02",
// 	"lottery_turntable_starrysky_c_ani_03",
// 	"lottery_turntable_starrysky_c_ani_04",
// 	"lottery_turntable_starrysky_c_ani_05",
// 	"lottery_turntable_starrysky_c_ani_06",
// 	"lottery_turntable_starrysky_c_ani_07",

// 	// 音樂
// 	"lottery_bgm_gaming",
// }
// // 遊戲抽獎-----end

// // 鑑定師-----start
// fieldsMonopoly = []string{
// 	// 鑑定師自定義
// 	"monopoly_classic_h_pic_01",
// 	"monopoly_classic_h_pic_02",
// 	"monopoly_classic_h_pic_03",
// 	"monopoly_classic_h_pic_04",
// 	"monopoly_classic_h_pic_05",
// 	"monopoly_classic_h_pic_06",
// 	"monopoly_classic_h_pic_07",
// 	"monopoly_classic_h_pic_08",
// 	"monopoly_classic_g_pic_01",
// 	"monopoly_classic_g_pic_02",
// 	"monopoly_classic_g_pic_03",
// 	"monopoly_classic_g_pic_04",
// 	"monopoly_classic_g_pic_05",
// 	"monopoly_classic_g_pic_06",
// 	"monopoly_classic_g_pic_07",
// 	"monopoly_classic_c_pic_01",
// 	"monopoly_classic_c_pic_02",
// 	"monopoly_classic_g_ani_01",
// 	"monopoly_classic_g_ani_02",
// 	"monopoly_classic_c_ani_01",

// 	"monopoly_redpack_h_pic_01",
// 	"monopoly_redpack_h_pic_02",
// 	"monopoly_redpack_h_pic_03",
// 	"monopoly_redpack_h_pic_04",
// 	"monopoly_redpack_h_pic_05",
// 	"monopoly_redpack_h_pic_06",
// 	"monopoly_redpack_h_pic_07",
// 	"monopoly_redpack_h_pic_08",
// 	"monopoly_redpack_h_pic_09",
// 	"monopoly_redpack_h_pic_10",
// 	"monopoly_redpack_h_pic_11",
// 	"monopoly_redpack_h_pic_12",
// 	"monopoly_redpack_h_pic_13",
// 	"monopoly_redpack_h_pic_14",
// 	"monopoly_redpack_h_pic_15",
// 	"monopoly_redpack_h_pic_16",
// 	"monopoly_redpack_g_pic_01",
// 	"monopoly_redpack_g_pic_02",
// 	"monopoly_redpack_g_pic_03",
// 	"monopoly_redpack_g_pic_04",
// 	"monopoly_redpack_g_pic_05",
// 	"monopoly_redpack_g_pic_06",
// 	"monopoly_redpack_g_pic_07",
// 	"monopoly_redpack_g_pic_08",
// 	"monopoly_redpack_g_pic_09",
// 	"monopoly_redpack_g_pic_10",
// 	"monopoly_redpack_c_pic_01",
// 	"monopoly_redpack_c_pic_02",
// 	"monopoly_redpack_c_pic_03",
// 	"monopoly_redpack_h_ani_01",
// 	"monopoly_redpack_h_ani_02",
// 	"monopoly_redpack_h_ani_03",
// 	"monopoly_redpack_g_ani_01",
// 	"monopoly_redpack_g_ani_02",
// 	"monopoly_redpack_c_ani_01",

// 	"monopoly_newyear_rabbit_h_pic_01",
// 	"monopoly_newyear_rabbit_h_pic_02",
// 	"monopoly_newyear_rabbit_h_pic_03",
// 	"monopoly_newyear_rabbit_h_pic_04",
// 	"monopoly_newyear_rabbit_h_pic_05",
// 	"monopoly_newyear_rabbit_h_pic_06",
// 	"monopoly_newyear_rabbit_h_pic_07",
// 	"monopoly_newyear_rabbit_h_pic_08",
// 	"monopoly_newyear_rabbit_h_pic_09",
// 	"monopoly_newyear_rabbit_h_pic_10",
// 	"monopoly_newyear_rabbit_h_pic_11",
// 	"monopoly_newyear_rabbit_h_pic_12",
// 	"monopoly_newyear_rabbit_g_pic_01",
// 	"monopoly_newyear_rabbit_g_pic_02",
// 	"monopoly_newyear_rabbit_g_pic_03",
// 	"monopoly_newyear_rabbit_g_pic_04",
// 	"monopoly_newyear_rabbit_g_pic_05",
// 	"monopoly_newyear_rabbit_g_pic_06",
// 	"monopoly_newyear_rabbit_g_pic_07",
// 	"monopoly_newyear_rabbit_c_pic_01",
// 	"monopoly_newyear_rabbit_c_pic_02",
// 	"monopoly_newyear_rabbit_c_pic_03",
// 	"monopoly_newyear_rabbit_h_ani_01",
// 	"monopoly_newyear_rabbit_h_ani_02",
// 	"monopoly_newyear_rabbit_g_ani_01",
// 	"monopoly_newyear_rabbit_g_ani_02",
// 	"monopoly_newyear_rabbit_c_ani_01",

// 	"monopoly_sashimi_h_pic_01",
// 	"monopoly_sashimi_h_pic_02",
// 	"monopoly_sashimi_h_pic_03",
// 	"monopoly_sashimi_h_pic_04",
// 	"monopoly_sashimi_h_pic_05",
// 	"monopoly_sashimi_g_pic_01",
// 	"monopoly_sashimi_g_pic_02",
// 	"monopoly_sashimi_g_pic_03",
// 	"monopoly_sashimi_g_pic_04",
// 	"monopoly_sashimi_g_pic_05",
// 	"monopoly_sashimi_g_pic_06",
// 	"monopoly_sashimi_c_pic_01",
// 	"monopoly_sashimi_c_pic_02",
// 	"monopoly_sashimi_h_ani_01",
// 	"monopoly_sashimi_h_ani_02",
// 	"monopoly_sashimi_g_ani_01",
// 	"monopoly_sashimi_g_ani_02",
// 	"monopoly_sashimi_c_ani_01",

// 	// 音樂
// 	"monopoly_bgm_start",
// 	"monopoly_bgm_gaming",
// 	"monopoly_bgm_end",
// }

// // 鑑定師-----end

// // 快問快答-----start

// fieldsQA = []string{
// 	// 快問快答自定義
// 	"qa_classic_h_pic_01",
// 	"qa_classic_h_pic_02",
// 	"qa_classic_h_pic_03",
// 	"qa_classic_h_pic_04",
// 	"qa_classic_h_pic_05",
// 	"qa_classic_h_pic_06",
// 	"qa_classic_h_pic_07",
// 	"qa_classic_h_pic_08",
// 	"qa_classic_h_pic_09",
// 	"qa_classic_h_pic_10",
// 	"qa_classic_h_pic_11",
// 	"qa_classic_h_pic_12",
// 	"qa_classic_h_pic_13",
// 	"qa_classic_h_pic_14",
// 	"qa_classic_h_pic_15",
// 	"qa_classic_h_pic_16",
// 	"qa_classic_h_pic_17",
// 	"qa_classic_h_pic_18",
// 	"qa_classic_h_pic_19",
// 	"qa_classic_h_pic_20",
// 	"qa_classic_h_pic_21",
// 	"qa_classic_h_pic_22",
// 	"qa_classic_g_pic_01",
// 	"qa_classic_g_pic_02",
// 	"qa_classic_g_pic_03",
// 	"qa_classic_g_pic_04",
// 	"qa_classic_g_pic_05",
// 	"qa_classic_c_pic_01",
// 	"qa_classic_h_ani_01",
// 	"qa_classic_h_ani_02",
// 	"qa_classic_g_ani_01",
// 	"qa_classic_g_ani_02",

// 	"qa_electric_h_pic_01",
// 	"qa_electric_h_pic_02",
// 	"qa_electric_h_pic_03",
// 	"qa_electric_h_pic_04",
// 	"qa_electric_h_pic_05",
// 	"qa_electric_h_pic_06",
// 	"qa_electric_h_pic_07",
// 	"qa_electric_h_pic_08",
// 	"qa_electric_h_pic_09",
// 	"qa_electric_h_pic_10",
// 	"qa_electric_h_pic_11",
// 	"qa_electric_h_pic_12",
// 	"qa_electric_h_pic_13",
// 	"qa_electric_h_pic_14",
// 	"qa_electric_h_pic_15",
// 	"qa_electric_h_pic_16",
// 	"qa_electric_h_pic_17",
// 	"qa_electric_h_pic_18",
// 	"qa_electric_h_pic_19",
// 	"qa_electric_h_pic_20",
// 	"qa_electric_h_pic_21",
// 	"qa_electric_h_pic_22",
// 	"qa_electric_h_pic_23",
// 	"qa_electric_h_pic_24",
// 	"qa_electric_h_pic_25",
// 	"qa_electric_h_pic_26",
// 	"qa_electric_g_pic_01",
// 	"qa_electric_g_pic_02",
// 	"qa_electric_g_pic_03",
// 	"qa_electric_g_pic_04",
// 	"qa_electric_g_pic_05",
// 	"qa_electric_g_pic_06",
// 	"qa_electric_g_pic_07",
// 	"qa_electric_g_pic_08",
// 	"qa_electric_g_pic_09",
// 	"qa_electric_c_pic_01",
// 	"qa_electric_h_ani_01",
// 	"qa_electric_h_ani_02",
// 	"qa_electric_h_ani_03",
// 	"qa_electric_h_ani_04",
// 	"qa_electric_h_ani_05",
// 	"qa_electric_g_ani_01",
// 	"qa_electric_g_ani_02",
// 	"qa_electric_c_ani_01",

// 	"qa_moonfestival_h_pic_01",
// 	"qa_moonfestival_h_pic_02",
// 	"qa_moonfestival_h_pic_03",
// 	"qa_moonfestival_h_pic_04",
// 	"qa_moonfestival_h_pic_05",
// 	"qa_moonfestival_h_pic_06",
// 	"qa_moonfestival_h_pic_07",
// 	"qa_moonfestival_h_pic_08",
// 	"qa_moonfestival_h_pic_09",
// 	"qa_moonfestival_h_pic_10",
// 	"qa_moonfestival_h_pic_11",
// 	"qa_moonfestival_h_pic_12",
// 	"qa_moonfestival_h_pic_13",
// 	"qa_moonfestival_h_pic_14",
// 	"qa_moonfestival_h_pic_15",
// 	"qa_moonfestival_h_pic_16",
// 	"qa_moonfestival_h_pic_17",
// 	"qa_moonfestival_h_pic_18",
// 	"qa_moonfestival_h_pic_19",
// 	"qa_moonfestival_h_pic_20",
// 	"qa_moonfestival_h_pic_21",
// 	"qa_moonfestival_h_pic_22",
// 	"qa_moonfestival_h_pic_23",
// 	"qa_moonfestival_h_pic_24",
// 	"qa_moonfestival_g_pic_01",
// 	"qa_moonfestival_g_pic_02",
// 	"qa_moonfestival_g_pic_03",
// 	"qa_moonfestival_g_pic_04",
// 	"qa_moonfestival_g_pic_05",
// 	"qa_moonfestival_c_pic_01",
// 	"qa_moonfestival_c_pic_02",
// 	"qa_moonfestival_c_pic_03",
// 	"qa_moonfestival_h_ani_01",
// 	"qa_moonfestival_h_ani_02",
// 	"qa_moonfestival_g_ani_01",
// 	"qa_moonfestival_g_ani_02",
// 	"qa_moonfestival_g_ani_03",

// 	"qa_newyear_dragon_h_pic_01",
// 	"qa_newyear_dragon_h_pic_02",
// 	"qa_newyear_dragon_h_pic_03",
// 	"qa_newyear_dragon_h_pic_04",
// 	"qa_newyear_dragon_h_pic_05",
// 	"qa_newyear_dragon_h_pic_06",
// 	"qa_newyear_dragon_h_pic_07",
// 	"qa_newyear_dragon_h_pic_08",
// 	"qa_newyear_dragon_h_pic_09",
// 	"qa_newyear_dragon_h_pic_10",
// 	"qa_newyear_dragon_h_pic_11",
// 	"qa_newyear_dragon_h_pic_12",
// 	"qa_newyear_dragon_h_pic_13",
// 	"qa_newyear_dragon_h_pic_14",
// 	"qa_newyear_dragon_h_pic_15",
// 	"qa_newyear_dragon_h_pic_16",
// 	"qa_newyear_dragon_h_pic_17",
// 	"qa_newyear_dragon_h_pic_18",
// 	"qa_newyear_dragon_h_pic_19",
// 	"qa_newyear_dragon_h_pic_20",
// 	"qa_newyear_dragon_h_pic_21",
// 	"qa_newyear_dragon_h_pic_22",
// 	"qa_newyear_dragon_h_pic_23",
// 	"qa_newyear_dragon_h_pic_24",
// 	"qa_newyear_dragon_g_pic_01",
// 	"qa_newyear_dragon_g_pic_02",
// 	"qa_newyear_dragon_g_pic_03",
// 	"qa_newyear_dragon_g_pic_04",
// 	"qa_newyear_dragon_g_pic_05",
// 	"qa_newyear_dragon_g_pic_06",
// 	"qa_newyear_dragon_c_pic_01",
// 	"qa_newyear_dragon_h_ani_01",
// 	"qa_newyear_dragon_h_ani_02",
// 	"qa_newyear_dragon_g_ani_01",
// 	"qa_newyear_dragon_g_ani_02",
// 	"qa_newyear_dragon_g_ani_03",
// 	"qa_newyear_dragon_c_ani_01",

// 	// 音樂
// 	"qa_bgm_start",
// 	"qa_bgm_gaming",
// 	"qa_bgm_end",
// }
// // 快問快答-----end

// // 搖紅包-----start
// fieldsRedpack = []string{
// 	// 搖紅包自定義
// 	"redpack_classic_h_pic_01",
// 	"redpack_classic_h_pic_02",
// 	"redpack_classic_h_pic_03",
// 	"redpack_classic_h_pic_04",
// 	"redpack_classic_h_pic_05",
// 	"redpack_classic_h_pic_06",
// 	"redpack_classic_h_pic_07",
// 	"redpack_classic_h_pic_08",
// 	"redpack_classic_h_pic_09",
// 	"redpack_classic_h_pic_10",
// 	"redpack_classic_h_pic_11",
// 	"redpack_classic_h_pic_12",
// 	"redpack_classic_h_pic_13",
// 	"redpack_classic_g_pic_01",
// 	"redpack_classic_g_pic_02",
// 	"redpack_classic_g_pic_03",
// 	"redpack_classic_h_ani_01",
// 	"redpack_classic_h_ani_02",
// 	"redpack_classic_g_ani_01",
// 	"redpack_classic_g_ani_02",
// 	"redpack_classic_g_ani_03",

// 	"redpack_cherry_h_pic_01",
// 	"redpack_cherry_h_pic_02",
// 	"redpack_cherry_h_pic_03",
// 	"redpack_cherry_h_pic_04",
// 	"redpack_cherry_h_pic_05",
// 	"redpack_cherry_h_pic_06",
// 	"redpack_cherry_h_pic_07",
// 	"redpack_cherry_g_pic_01",
// 	"redpack_cherry_g_pic_02",
// 	"redpack_cherry_h_ani_01",
// 	"redpack_cherry_h_ani_02",
// 	"redpack_cherry_g_ani_01",
// 	"redpack_cherry_g_ani_02",

// 	"redpack_christmas_h_pic_01",
// 	"redpack_christmas_h_pic_02",
// 	"redpack_christmas_h_pic_03",
// 	"redpack_christmas_h_pic_04",
// 	"redpack_christmas_h_pic_05",
// 	"redpack_christmas_h_pic_06",
// 	"redpack_christmas_h_pic_07",
// 	"redpack_christmas_h_pic_08",
// 	"redpack_christmas_h_pic_09",
// 	"redpack_christmas_h_pic_10",
// 	"redpack_christmas_h_pic_11",
// 	"redpack_christmas_h_pic_12",
// 	"redpack_christmas_h_pic_13",
// 	"redpack_christmas_g_pic_01",
// 	"redpack_christmas_g_pic_02",
// 	"redpack_christmas_g_pic_03",
// 	"redpack_christmas_g_pic_04",
// 	"redpack_christmas_c_pic_01",
// 	"redpack_christmas_c_pic_02",
// 	"redpack_christmas_c_ani_01",

// 	// 音樂
// 	"redpack_bgm_start",
// 	"redpack_bgm_gaming",
// 	"redpack_bgm_end",
// }
// // 搖紅包-----end

// // 套紅包-----start
// fieldsRopepack = []string{
// 	// 套紅包自定義
// 	"ropepack_classic_h_pic_01",
// 	"ropepack_classic_h_pic_02",
// 	"ropepack_classic_h_pic_03",
// 	"ropepack_classic_h_pic_04",
// 	"ropepack_classic_h_pic_05",
// 	"ropepack_classic_h_pic_06",
// 	"ropepack_classic_h_pic_07",
// 	"ropepack_classic_h_pic_08",
// 	"ropepack_classic_h_pic_09",
// 	"ropepack_classic_h_pic_10",
// 	"ropepack_classic_g_pic_01",
// 	"ropepack_classic_g_pic_02",
// 	"ropepack_classic_g_pic_03",
// 	"ropepack_classic_g_pic_04",
// 	"ropepack_classic_g_pic_05",
// 	"ropepack_classic_g_pic_06",
// 	"ropepack_classic_h_ani_01",
// 	"ropepack_classic_g_ani_01",
// 	"ropepack_classic_g_ani_02",
// 	"ropepack_classic_c_ani_01",

// 	"ropepack_newyear_rabbit_h_pic_01",
// 	"ropepack_newyear_rabbit_h_pic_02",
// 	"ropepack_newyear_rabbit_h_pic_03",
// 	"ropepack_newyear_rabbit_h_pic_04",
// 	"ropepack_newyear_rabbit_h_pic_05",
// 	"ropepack_newyear_rabbit_h_pic_06",
// 	"ropepack_newyear_rabbit_h_pic_07",
// 	"ropepack_newyear_rabbit_h_pic_08",
// 	"ropepack_newyear_rabbit_h_pic_09",
// 	"ropepack_newyear_rabbit_g_pic_01",
// 	"ropepack_newyear_rabbit_g_pic_02",
// 	"ropepack_newyear_rabbit_g_pic_03",
// 	"ropepack_newyear_rabbit_h_ani_01",
// 	"ropepack_newyear_rabbit_g_ani_01",
// 	"ropepack_newyear_rabbit_g_ani_02",
// 	"ropepack_newyear_rabbit_g_ani_03",
// 	"ropepack_newyear_rabbit_c_ani_01",
// 	"ropepack_newyear_rabbit_c_ani_02",

// 	"ropepack_moonfestival_h_pic_01",
// 	"ropepack_moonfestival_h_pic_02",
// 	"ropepack_moonfestival_h_pic_03",
// 	"ropepack_moonfestival_h_pic_04",
// 	"ropepack_moonfestival_h_pic_05",
// 	"ropepack_moonfestival_h_pic_06",
// 	"ropepack_moonfestival_h_pic_07",
// 	"ropepack_moonfestival_h_pic_08",
// 	"ropepack_moonfestival_h_pic_09",
// 	"ropepack_moonfestival_g_pic_01",
// 	"ropepack_moonfestival_g_pic_02",
// 	"ropepack_moonfestival_c_pic_01",
// 	"ropepack_moonfestival_h_ani_01",
// 	"ropepack_moonfestival_g_ani_01",
// 	"ropepack_moonfestival_g_ani_02",
// 	"ropepack_moonfestival_c_ani_01",
// 	"ropepack_moonfestival_c_ani_02",

// 	"ropepack_3D_h_pic_01",
// 	"ropepack_3D_h_pic_02",
// 	"ropepack_3D_h_pic_03",
// 	"ropepack_3D_h_pic_04",
// 	"ropepack_3D_h_pic_05",
// 	"ropepack_3D_h_pic_06",
// 	"ropepack_3D_h_pic_07",
// 	"ropepack_3D_h_pic_08",
// 	"ropepack_3D_h_pic_09",
// 	"ropepack_3D_h_pic_10",
// 	"ropepack_3D_h_pic_11",
// 	"ropepack_3D_h_pic_12",
// 	"ropepack_3D_h_pic_13",
// 	"ropepack_3D_h_pic_14",
// 	"ropepack_3D_h_pic_15",
// 	"ropepack_3D_g_pic_01",
// 	"ropepack_3D_g_pic_02",
// 	"ropepack_3D_g_pic_03",
// 	"ropepack_3D_g_pic_04",
// 	"ropepack_3D_h_ani_01",
// 	"ropepack_3D_h_ani_02",
// 	"ropepack_3D_h_ani_03",
// 	"ropepack_3D_g_ani_01",
// 	"ropepack_3D_g_ani_02",
// 	"ropepack_3D_c_ani_01",

// 	// 音樂
// 	"ropepack_bgm_start",
// 	"ropepack_bgm_gaming",
// 	"ropepack_bgm_end",
// }
// // 套紅包-----end

// // 拔河遊戲-----start
// fieldsTugofwar = []string{
// 	// 拔河遊戲自定義
// 	"tugofwar_classic_h_pic_01",
// 	"tugofwar_classic_h_pic_02",
// 	"tugofwar_classic_h_pic_03",
// 	"tugofwar_classic_h_pic_04",
// 	"tugofwar_classic_h_pic_05",
// 	"tugofwar_classic_h_pic_06",
// 	"tugofwar_classic_h_pic_07",
// 	"tugofwar_classic_h_pic_08",
// 	"tugofwar_classic_h_pic_09",
// 	"tugofwar_classic_h_pic_10",
// 	"tugofwar_classic_h_pic_11",
// 	"tugofwar_classic_h_pic_12",
// 	"tugofwar_classic_h_pic_13",
// 	"tugofwar_classic_h_pic_14",
// 	"tugofwar_classic_h_pic_15",
// 	"tugofwar_classic_h_pic_16",
// 	"tugofwar_classic_h_pic_17",
// 	"tugofwar_classic_h_pic_18",
// 	"tugofwar_classic_h_pic_19",
// 	"tugofwar_classic_g_pic_01",
// 	"tugofwar_classic_g_pic_02",
// 	"tugofwar_classic_g_pic_03",
// 	"tugofwar_classic_g_pic_04",
// 	"tugofwar_classic_g_pic_05",
// 	"tugofwar_classic_g_pic_06",
// 	"tugofwar_classic_g_pic_07",
// 	"tugofwar_classic_g_pic_08",
// 	"tugofwar_classic_g_pic_09",
// 	"tugofwar_classic_h_ani_01",
// 	"tugofwar_classic_h_ani_02",
// 	"tugofwar_classic_h_ani_03",
// 	"tugofwar_classic_c_ani_01",

// 	"tugofwar_school_h_pic_01",
// 	"tugofwar_school_h_pic_02",
// 	"tugofwar_school_h_pic_03",
// 	"tugofwar_school_h_pic_04",
// 	"tugofwar_school_h_pic_05",
// 	"tugofwar_school_h_pic_06",
// 	"tugofwar_school_h_pic_07",
// 	"tugofwar_school_h_pic_08",
// 	"tugofwar_school_h_pic_09",
// 	"tugofwar_school_h_pic_10",
// 	"tugofwar_school_h_pic_11",
// 	"tugofwar_school_h_pic_12",
// 	"tugofwar_school_h_pic_13",
// 	"tugofwar_school_h_pic_14",
// 	"tugofwar_school_h_pic_15",
// 	"tugofwar_school_h_pic_16",
// 	"tugofwar_school_h_pic_17",
// 	"tugofwar_school_h_pic_18",
// 	"tugofwar_school_h_pic_19",
// 	"tugofwar_school_h_pic_20",
// 	"tugofwar_school_h_pic_21",
// 	"tugofwar_school_h_pic_22",
// 	"tugofwar_school_h_pic_23",
// 	"tugofwar_school_h_pic_24",
// 	"tugofwar_school_h_pic_25",
// 	"tugofwar_school_h_pic_26",
// 	"tugofwar_school_g_pic_01",
// 	"tugofwar_school_g_pic_02",
// 	"tugofwar_school_g_pic_03",
// 	"tugofwar_school_g_pic_04",
// 	"tugofwar_school_g_pic_05",
// 	"tugofwar_school_g_pic_06",
// 	"tugofwar_school_g_pic_07",
// 	"tugofwar_school_c_pic_01",
// 	"tugofwar_school_c_pic_02",
// 	"tugofwar_school_c_pic_03",
// 	"tugofwar_school_c_pic_04",
// 	"tugofwar_school_h_ani_01",
// 	"tugofwar_school_h_ani_02",
// 	"tugofwar_school_h_ani_03",
// 	"tugofwar_school_h_ani_04",
// 	"tugofwar_school_h_ani_05",
// 	"tugofwar_school_h_ani_06",
// 	"tugofwar_school_h_ani_07",

// 	"tugofwar_christmas_h_pic_01",
// 	"tugofwar_christmas_h_pic_02",
// 	"tugofwar_christmas_h_pic_03",
// 	"tugofwar_christmas_h_pic_04",
// 	"tugofwar_christmas_h_pic_05",
// 	"tugofwar_christmas_h_pic_06",
// 	"tugofwar_christmas_h_pic_07",
// 	"tugofwar_christmas_h_pic_08",
// 	"tugofwar_christmas_h_pic_09",
// 	"tugofwar_christmas_h_pic_10",
// 	"tugofwar_christmas_h_pic_11",
// 	"tugofwar_christmas_h_pic_12",
// 	"tugofwar_christmas_h_pic_13",
// 	"tugofwar_christmas_h_pic_14",
// 	"tugofwar_christmas_h_pic_15",
// 	"tugofwar_christmas_h_pic_16",
// 	"tugofwar_christmas_h_pic_17",
// 	"tugofwar_christmas_h_pic_18",
// 	"tugofwar_christmas_h_pic_19",
// 	"tugofwar_christmas_h_pic_20",
// 	"tugofwar_christmas_h_pic_21",
// 	"tugofwar_christmas_g_pic_01",
// 	"tugofwar_christmas_g_pic_02",
// 	"tugofwar_christmas_g_pic_03",
// 	"tugofwar_christmas_g_pic_04",
// 	"tugofwar_christmas_g_pic_05",
// 	"tugofwar_christmas_g_pic_06",
// 	"tugofwar_christmas_c_pic_01",
// 	"tugofwar_christmas_c_pic_02",
// 	"tugofwar_christmas_c_pic_03",
// 	"tugofwar_christmas_c_pic_04",
// 	"tugofwar_christmas_h_ani_01",
// 	"tugofwar_christmas_h_ani_02",
// 	"tugofwar_christmas_h_ani_03",
// 	"tugofwar_christmas_c_ani_01",
// 	"tugofwar_christmas_c_ani_02",

// 	// 音樂
// 	"tugofwar_bgm_start",
// 	"tugofwar_bgm_gaming",
// 	"tugofwar_bgm_end",
// }
// // 拔河遊戲-----end

// // 投票-----start
// fieldsVote = []string{
// 	// 投票遊戲自定義
// 	"vote_classic_h_pic_01",
// 	"vote_classic_h_pic_02",
// 	"vote_classic_h_pic_03",
// 	"vote_classic_h_pic_04",
// 	"vote_classic_h_pic_05",
// 	"vote_classic_h_pic_06",
// 	"vote_classic_h_pic_07",
// 	"vote_classic_h_pic_08",
// 	"vote_classic_h_pic_09",
// 	"vote_classic_h_pic_10",
// 	"vote_classic_h_pic_11",
// 	"vote_classic_h_pic_13",
// 	"vote_classic_h_pic_14",
// 	"vote_classic_h_pic_15",
// 	"vote_classic_h_pic_16",
// 	"vote_classic_h_pic_17",
// 	"vote_classic_h_pic_18",
// 	"vote_classic_h_pic_19",
// 	"vote_classic_h_pic_20",
// 	"vote_classic_h_pic_21",
// 	"vote_classic_h_pic_23",
// 	"vote_classic_h_pic_24",
// 	"vote_classic_h_pic_25",
// 	"vote_classic_h_pic_26",
// 	"vote_classic_h_pic_27",
// 	"vote_classic_h_pic_28",
// 	"vote_classic_h_pic_29",
// 	"vote_classic_h_pic_30",
// 	"vote_classic_h_pic_31",
// 	"vote_classic_h_pic_32",
// 	"vote_classic_h_pic_33",
// 	"vote_classic_h_pic_34",
// 	"vote_classic_h_pic_35",
// 	"vote_classic_h_pic_36",
// 	"vote_classic_h_pic_37",
// 	"vote_classic_g_pic_01",
// 	"vote_classic_g_pic_02",
// 	"vote_classic_g_pic_03",
// 	"vote_classic_g_pic_04",
// 	"vote_classic_g_pic_05",
// 	"vote_classic_g_pic_06",
// 	"vote_classic_g_pic_07",
// 	"vote_classic_c_pic_01",
// 	"vote_classic_c_pic_02",
// 	"vote_classic_c_pic_03",
// 	"vote_classic_c_pic_04",

// 	// 音樂
// 	"vote_bgm_gaming",
// }
// // 投票-----end

// // 敲敲樂-----start
// fieldsWhackMole = []string{
// 	// 敲敲樂自定義
// 	"whackmole_classic_h_pic_01",
// 	"whackmole_classic_h_pic_02",
// 	"whackmole_classic_h_pic_03",
// 	"whackmole_classic_h_pic_04",
// 	"whackmole_classic_h_pic_05",
// 	"whackmole_classic_h_pic_06",
// 	"whackmole_classic_h_pic_07",
// 	"whackmole_classic_h_pic_08",
// 	"whackmole_classic_h_pic_09",
// 	"whackmole_classic_h_pic_10",
// 	"whackmole_classic_h_pic_11",
// 	"whackmole_classic_h_pic_12",
// 	"whackmole_classic_h_pic_13",
// 	"whackmole_classic_h_pic_14",
// 	"whackmole_classic_h_pic_15",
// 	"whackmole_classic_g_pic_01",
// 	"whackmole_classic_g_pic_02",
// 	"whackmole_classic_g_pic_03",
// 	"whackmole_classic_g_pic_04",
// 	"whackmole_classic_g_pic_05",
// 	"whackmole_classic_c_pic_01",
// 	"whackmole_classic_c_pic_02",
// 	"whackmole_classic_c_pic_03",
// 	"whackmole_classic_c_pic_04",
// 	"whackmole_classic_c_pic_05",
// 	"whackmole_classic_c_pic_06",
// 	"whackmole_classic_c_pic_07",
// 	"whackmole_classic_c_pic_08",
// 	"whackmole_classic_c_ani_01",

// 	"whackmole_halloween_h_pic_01",
// 	"whackmole_halloween_h_pic_02",
// 	"whackmole_halloween_h_pic_03",
// 	"whackmole_halloween_h_pic_04",
// 	"whackmole_halloween_h_pic_05",
// 	"whackmole_halloween_h_pic_06",
// 	"whackmole_halloween_h_pic_07",
// 	"whackmole_halloween_h_pic_08",
// 	"whackmole_halloween_h_pic_09",
// 	"whackmole_halloween_h_pic_10",
// 	"whackmole_halloween_h_pic_11",
// 	"whackmole_halloween_h_pic_12",
// 	"whackmole_halloween_h_pic_13",
// 	"whackmole_halloween_h_pic_14",
// 	"whackmole_halloween_h_pic_15",
// 	"whackmole_halloween_g_pic_01",
// 	"whackmole_halloween_g_pic_02",
// 	"whackmole_halloween_g_pic_03",
// 	"whackmole_halloween_g_pic_04",
// 	"whackmole_halloween_g_pic_05",
// 	"whackmole_halloween_c_pic_01",
// 	"whackmole_halloween_c_pic_02",
// 	"whackmole_halloween_c_pic_03",
// 	"whackmole_halloween_c_pic_04",
// 	"whackmole_halloween_c_pic_05",
// 	"whackmole_halloween_c_pic_06",
// 	"whackmole_halloween_c_pic_07",
// 	"whackmole_halloween_c_pic_08",
// 	"whackmole_halloween_c_ani_01",

// 	"whackmole_christmas_h_pic_01",
// 	"whackmole_christmas_h_pic_02",
// 	"whackmole_christmas_h_pic_03",
// 	"whackmole_christmas_h_pic_04",
// 	"whackmole_christmas_h_pic_05",
// 	"whackmole_christmas_h_pic_06",
// 	"whackmole_christmas_h_pic_07",
// 	"whackmole_christmas_h_pic_08",
// 	"whackmole_christmas_h_pic_09",
// 	"whackmole_christmas_h_pic_10",
// 	"whackmole_christmas_h_pic_11",
// 	"whackmole_christmas_h_pic_12",
// 	"whackmole_christmas_h_pic_13",
// 	"whackmole_christmas_h_pic_14",
// 	"whackmole_christmas_g_pic_01",
// 	"whackmole_christmas_g_pic_02",
// 	"whackmole_christmas_g_pic_03",
// 	"whackmole_christmas_g_pic_04",
// 	"whackmole_christmas_g_pic_05",
// 	"whackmole_christmas_g_pic_06",
// 	"whackmole_christmas_g_pic_07",
// 	"whackmole_christmas_g_pic_08",
// 	"whackmole_christmas_c_pic_01",
// 	"whackmole_christmas_c_pic_02",
// 	"whackmole_christmas_c_pic_03",
// 	"whackmole_christmas_c_pic_04",
// 	"whackmole_christmas_c_pic_05",
// 	"whackmole_christmas_c_pic_06",
// 	"whackmole_christmas_c_pic_07",
// 	"whackmole_christmas_c_pic_08",
// 	"whackmole_christmas_c_ani_01",
// 	"whackmole_christmas_c_ani_02",:

// 	// 音樂
// 	"whackmole_bgm_start",
// 	"whackmole_bgm_gaming",
// 	"whackmole_bgm_end",
// }

// // 敲敲樂-----end

// updateInfos = []TableUpdateInfo{
// 	{TableName: "", Fields: fieldsRopepack},
// 	{TableName: "", Fields: fieldsRedpack},
// 	{TableName: "", Fields: fieldsLottery},
// 	{TableName: "", Fields: fieldsDrawNumbers},
// 	{TableName: "", Fields: fieldsWhackMole},
// 	{TableName: "", Fields: fieldsMonopoly},
// 	{TableName: "", Fields: fieldsQA},
// 	{TableName: "", Fields: fieldsTugofwar},
// 	{TableName: "", Fields: fieldsBingo},
// 	{TableName: "", Fields: fields3DGachaMachine},
// 	{TableName: "", Fields: fieldsVote},
// }

// utils.SetInt64IfNotEmpty(data, "max_times", model.MaxTimes)
// utils.SetInt64IfNotEmpty(data, "percent", model.Percent)
// utils.SetInt64IfNotEmpty(data, "first_prize", model.FirstPrize)
// utils.SetInt64IfNotEmpty(data, "second_prize", model.SecondPrize)
// utils.SetInt64IfNotEmpty(data, "third_prize", model.ThirdPrize)
// utils.SetInt64IfNotEmpty(data, "general_prize", model.GeneralPrize)
// utils.SetInt64IfNotEmpty(data, "total_qa", model.TotalQA)
// utils.SetInt64IfNotEmpty(data, "qa_second", model.QASecond)
// utils.SetInt64IfNotEmpty(data, "bingo_line", model.BingoLine)
// utils.SetInt64IfNotEmpty(data, "max_number", model.MaxNumber)
// utils.SetInt64IfNotEmpty(data, "round_prize", model.RoundPrize)
// utils.SetInt64IfNotEmpty(data, "vote_times", model.VoteTimes)
// utils.SetFloat64IfNotEmpty(data, "gacha_machine_reflection", model.GachaMachineReflection)

// if val, ok := data[key]; ok && val != "" {
// 	fieldValues[key] = val
// }

// 所有遊戲自定義圖片欄位資料
// for _, updateInfo := range updateInfos {
// 	for _, field := range updateInfo.Fields {
// 		if val, ok := data[field]; ok && val != "" {
// 			// 資料不為空時，將資料加入fieldValues
// 			fieldValues[field] 
package models

import (
	"errors"
	"fmt"
	"hilive/modules/config"
	"hilive/modules/utils"
	"strconv"
	"time"
	"unicode/utf8"
)

// Add 新增遊戲場次資料
func (a GameModel) Add(isRedis bool, game, gameid string, model EditGameModel) error {
	var (
		gameStatus = "close"
		now, _     = time.ParseInLocation("2006-01-02 15:04:05",
			time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04"), time.Local) // 目前時間

		fields = make([]string, 0)

		// 扭蛋機-----start
		gachaMachinefields = []string{
			"id",
			"user_id",
			"activity_id",
			"game_id",
			"game",
			"topic",
			"skin",
			"music",
			"game_status",
			"game_order",
			"title",
			"allow",
			"max_times",
			"percent",
			"reflective_switch",
			"gacha_machine_reflection",
			"game_round",
			"game_second",
			"game_attend",
			"edit_times",

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
			"id",
			"user_id",
			"activity_id",
			"game_id",
			"game",
			"topic",
			"skin",
			"music",
			"game_status",
			"game_order",
			"title",
			"max_people",
			"people",
			"allow",
			"max_number",
			"bingo_line",
			"round_prize",
			"game_round",
			"game_second",
			"game_attend",
			// "bingo_round",
			"edit_times",

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
			"id",
			"user_id",
			"activity_id",
			"game_id",
			"game",
			"topic",
			"skin",
			"music",
			"game_status",
			"game_order",
			"title",
			"limit_time",
			"second",
			"allow",
			"display_name",
			"game_round",
			"game_second",
			"game_attend",
			"edit_times",

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
			"id",
			"user_id",
			"activity_id",
			"game_id",
			"game",
			"topic",
			"skin",
			"music",
			"game_status",
			"game_order",
			"title",
			"allow",
			"max_times",
			"game_type",
			"game_round",
			"game_second",
			"game_attend",
			"edit_times",

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
			"id",
			"user_id",
			"activity_id",
			"game_id",
			"game",
			"topic",
			"skin",
			"music",
			"game_status",
			"game_order",
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
			"game_round",
			"game_second",
			"game_attend",
			"edit_times",

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
			"id",
			"user_id",
			"activity_id",
			"game_id",
			"game",
			"topic",
			"skin",
			"music",
			"game_status",
			"game_order",
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
			"game_round",
			"game_second",
			"game_attend",
			"edit_times",
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

			"qa_round", "qa_second", "total_qa",
		}

		// 快問快答-----end

		// 搖紅包-----start
		redpackfields = []string{
			"id",
			"user_id",
			"activity_id",
			"game_id",
			"game",
			"topic",
			"skin",
			"music",
			"game_status",
			"game_order",
			"title",
			"people",
			"max_people",
			"allow",
			"limit_time",
			"second",
			"percent",
			"game_round",
			"game_second",
			"game_attend",
			"edit_times",

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
			"id",
			"user_id",
			"activity_id",
			"game_id",
			"game",
			"topic",
			"skin",
			"music",
			"game_status",
			"game_order",
			"title",
			"people",
			"max_people",
			"allow",
			"limit_time",
			"second",
			"percent",
			"game_round",
			"game_second",
			"game_attend",
			"edit_times",

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
			"id",
			"user_id",
			"activity_id",
			"game_id",
			"game",
			"topic",
			"skin",
			"music",
			"game_status",
			"game_order",
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
			"game_round",
			"game_second",
			"game_attend",
			"edit_times",

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
			"id",
			"user_id",
			"activity_id",
			"game_id",
			"game",
			"topic",
			"skin",
			"music",
			"game_status",
			"game_order",
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
			"game_round",
			"game_second",
			"game_attend",
			"edit_times",

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
			"id",
			"user_id",
			"activity_id",
			"game_id",
			"game",
			"topic",
			"skin",
			"music",
			"game_status",
			"game_order",
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
			"game_round",
			"game_second",
			"game_attend",
			"edit_times",

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

	if game != "redpack" && game != "ropepack" && game != "whack_mole" &&
		game != "lottery" && game != "monopoly" && game != "QA" &&
		game != "draw_numbers" && game != "tugofwar" &&
		game != "bingo" && game != "3DGachaMachine" &&
		game != "vote" {
		return errors.New("錯誤: 遊戲種類發生問題，請輸入有效的遊戲種類")
	}

	if model.Title == "" || utf8.RuneCountInString(model.Title) > 20 {
		return errors.New("錯誤: 標題上限為20個字元，請輸入有效的標題名稱")
	}

	if game == "lottery" {
		if model.GameType != "turntable" && model.GameType != "jiugongge" {
			return errors.New("錯誤: 遊戲類型資料發生問題，請輸入有效的遊戲類型")
		}
	}

	if model.LimitTime != "" {
		if model.LimitTime != "open" && model.LimitTime != "close" {
			return errors.New("錯誤: 是否限時資料發生問題，請輸入有效的資料")
		}
	}

	if model.Second != "" {
		if _, err := strconv.Atoi(model.Second); err != nil {
			return errors.New("錯誤: 限時秒數資料發生問題，請輸入有效的秒數")
		}
	}

	if model.MaxPeople != "" && model.People != "" {
		// 判斷遊戲人數上限
		maxPeopleInt, err1 := strconv.Atoi(model.MaxPeople)
		peopleInt, err2 := strconv.Atoi(model.People)
		if err1 != nil || err2 != nil || peopleInt > maxPeopleInt {
			return errors.New("錯誤: 遊戲人數上限資料發生問題，請輸入有效的遊戲人數上限")
		}
	}

	if model.MaxTimes != "" {
		if _, err := strconv.Atoi(model.MaxTimes); err != nil {
			return errors.New("錯誤: 遊戲上限次數發生問題，請輸入有效的遊戲次數")
		}
	}

	if model.Allow != "" {
		if model.Allow != "open" && model.Allow != "close" {
			return errors.New("錯誤: 允許重複搖中資料發生問題，請輸入有效的資料")
		}
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
	// 9. 賓果遊戲: 經典主題 (01_classic)
	if model.Topic != "01_classic" && model.Topic != "02_halloween" &&
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
		model.Topic != "05_3D_space" {
		return errors.New("錯誤: 主題資料發生問題，請輸入有效的主題")
	}

	if model.Skin != "classic" && model.Skin != "customize" {
		return errors.New("錯誤: 樣式資料發生問題，請輸入有效的樣式")
	}

	if model.Music != "classic" && model.Music != "customize" {
		return errors.New("錯誤: 音樂資料發生問題，請輸入有效的音樂")
	}

	if model.DisplayName != "" {
		if model.DisplayName != "open" && model.DisplayName != "close" {
			return errors.New("錯誤: 是否顯示中獎人員姓名頭像資料發生問題，請輸入有效的資料")
		}
	}

	if game == "3DGachaMachine" {
		if model.ReflectiveSwitch != "open" && model.ReflectiveSwitch != "close" {
			return errors.New("錯誤: 扭蛋盒的反光資料發生問題，請輸入有效的資料")
		}
	}

	if game == "QA" {
		if _, err := strconv.Atoi(model.TotalQA); err != nil {
			return errors.New("錯誤: 總題目數量發生問題，請輸入有效的題目數量")
		}
		if _, err := strconv.Atoi(model.QASecond); err != nil {
			return errors.New("錯誤: 題目顯示秒數發生問題，請輸入有效的題目顯示秒數")
		}
	}

	// 拔河遊戲
	if game == "tugofwar" {
		if model.AllowChooseTeam != "open" && model.AllowChooseTeam != "close" {
			return errors.New("錯誤: 允許玩家選擇隊伍資料發生問題，請輸入有效的資料")
		}

		if model.LeftTeamName == "" || model.LeftTeamPicture == "" ||
			model.RightTeamName == "" || model.RightTeamPicture == "" {
			return errors.New("錯誤: 隊伍名稱.照片資料發生問題，請輸入有效的資料")
		}

		if utf8.RuneCountInString(model.LeftTeamName) > 20 ||
			utf8.RuneCountInString(model.RightTeamName) > 20 {
			return errors.New("錯誤: 隊伍名稱上限為20個字元，請輸入有效的標題名稱")
		}

		if model.Prize != "uniform" && model.Prize != "all" {
			return errors.New("錯誤: 獎品發放資料發生問題，請輸入有效的資料")
		}
	}

	// 賓果遊戲
	if game == "bingo" {
		if line, err := strconv.Atoi(model.BingoLine); err != nil ||
			line < 1 || line > 10 {
			return errors.New("錯誤: 賓果連線數資料發生問題(最多10條線，最少1條線)，請輸入有效的連線數")
		}

		if number, err := strconv.Atoi(model.MaxNumber); err != nil ||
			number < 16 || number > 99 {
			return errors.New("錯誤: 最大號碼資料發生問題(號碼必須大於16且小於100)，請輸入有效的連線數")
		}

		if _, err := strconv.Atoi(model.RoundPrize); err != nil {
			return errors.New("錯誤: 每輪發獎人數資料發生問題，請輸入有效的資料")
		}
	}

	// 投票遊戲
	if game == "vote" {
		if model.VoteScreen != "bar_chart" && model.VoteScreen != "rank" && model.VoteScreen != "detail_information" {
			return errors.New("錯誤: 投票畫面資料發生問題，請輸入有效的資料")
		}

		if _, err := strconv.Atoi(model.VoteTimes); err != nil {
			return errors.New("錯誤: 人員投票次數資料發生問題，請輸入有效的資料")
		}

		if model.VoteMethod != "all_vote" && model.VoteMethod != "single_group" && model.VoteMethod != "all_group" {
			return errors.New("錯誤: 投票模式資料發生問題，請輸入有效的資料")
		}

		if model.VoteMethodPlayer != "one_vote" && model.VoteMethodPlayer != "free_vote" {
			return errors.New("錯誤: 玩家投票方式資料發生問題，請輸入有效的資料")
		}

		if model.VoteRestriction != "all_player" && model.VoteRestriction != "special_officer" {
			return errors.New("錯誤: 投票限制資料發生問題，請輸入有效的資料")
		}

		if model.Prize != "uniform" && model.Prize != "all" {
			return errors.New("錯誤: 獎品發放資料發生問題，請輸入有效的資料")
		}

		if model.AvatarShape != "circle" && model.AvatarShape != "square" {
			return errors.New("錯誤: 選項框資料發生問題，請輸入有效的資料")
		}

		// 判斷投票結束時間
		startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", model.VoteStartTime, time.Local)
		endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", model.VoteEndTime, time.Local)

		if !CompareDatetime(model.VoteStartTime, model.VoteEndTime) {
			return errors.New("錯誤: 投票時間發生問題，結束時間必須大於開始時間，請輸入有效的時間")
		}

		// log.Println("比較時間: ", now, startTime, endTime)
		// 比較時間，判斷遊戲狀態
		if now.Before(startTime) { // 現在時間在開始時間之前
			gameStatus = "close" // 關閉
		} else if now.Before(endTime) { // 現在時間在截止時間之前
			gameStatus = "gaming" // 遊戲中
		} else { // 現在時間在截止時間之後
			gameStatus = "calculate" // 結算狀態
		}

		if model.AutoPlay != "open" && model.AutoPlay != "close" {
			return errors.New("錯誤: 自動輪播資料發生問題，請輸入有效的資料")
		}
		if model.ShowRank != "open" && model.ShowRank != "close" {
			return errors.New("錯誤: 排名展示資料發生問題，請輸入有效的資料")
		}
		if model.TitleSwitch != "open" && model.TitleSwitch != "close" {
			return errors.New("錯誤: 場式名稱開關資料發生問題，請輸入有效的資料")
		}
		if model.ArrangementGuest != "list" && model.ArrangementGuest != "side_by_side" {
			return errors.New("錯誤: 玩家端選項排列方式資料發生問題，請輸入有效的資料")
		}
	}

	// 取得該活動遊戲場次數量(用於game_order，mongo)
	games, err := a.FindAll(model.ActivityID, "")
	if err != nil {
		return errors.New("錯誤: 查詢遊戲場次(activity_game)發生問題")
	}

	// 取得該活動遊戲場次數量(用於game_order，mysql)
	// games, err := a.Table(a.TableName).
	// 	Where("activity_id", "=", model.ActivityID).All()
	// if err != nil {
	// 	return errors.New("錯誤: 查詢遊戲場次(activity_game)發生問題")
	// }

	// 將 struct 轉換為 map[string]interface{} 格式
	data := utils.StructToMap(model)

	// 手動處理
	data["game"] = game
	data["game_status"] = gameStatus
	data["game_round"] = int64(1)
	data["game_second"] = utils.GetInt64(model.Second, 0)
	data["game_attend"] = int64(0)
	data["game_order"] = int64(len(games) + 1)
	data["edit_times"] = int64(0)
	data["qa_round"] = int64(1)
	// data["bingo_round"] = 0
	// data["left_team_game_attend"] = 0
	// data["right_team_game_attend"] = 0
	// string轉換int64.float64
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

	if game == "QA" {
		if model.QA1 == "" || model.QA1Options == "" ||
			model.QA1Answer == "" || model.QA1Score == "" {
			return errors.New("錯誤: 題目設置最少一題，請重新設置")
		}

		// 處理qa_n_score值
		for i := 1; i <= 20; i++ {
			// score := data[fmt.Sprintf("qa_%d_score", i)]

			// 將分數轉為int64寫入data中
			data[fmt.Sprintf("qa_%d_score", i)] = utils.GetInt64(data[fmt.Sprintf("qa_%d_score", i)], 0)
		}
	}

	// 寫入mongo，activity_game資料表
	// 將id資料手動寫入data
	// 取得mongo中的id資料(遞增處理)
	mongoID, _ := a.MongoConn.GetNextSequence(a.TableName)
	data["id"] = mongoID

	if _, err = a.MongoConn.InsertOne(a.TableName, FilterFields(data, fields)); err != nil {
		return errors.New("錯誤: 新增遊戲場次(mongo，activity_game)發生問題")
	}

	// 自定義處理
	if game == "redpack" {

	} else if game == "ropepack" {

	} else if game == "whack_mole" {

	} else if game == "lottery" {

	} else if game == "monopoly" {

	} else if game == "QA" {

	} else if game == "draw_numbers" {

	} else if game == "tugofwar" {

		// 建立拔河遊戲時，預設兩個獎品(勝方敗方)
		if err := DefaultPrizeModel().
			SetConn(a.DbConn, a.RedisConn, a.MongoConn).
			Add(true, "tugofwar", utils.UUID(20), EditPrizeModel{
				ActivityID:    model.ActivityID,
				GameID:        gameid,
				PrizeName:     "勝利隊伍",
				PrizeType:     "first",
				PrizePicture:  config.UPLOAD_SYSTEM_URL + "img-prize-pic.png",
				PrizeAmount:   "0",
				PrizePrice:    "0",
				PrizeMethod:   "site",
				PrizePassword: "win",
				TeamType:      "win",
			}); err != nil {
			return errors.New("錯誤: 新增勝方遊戲獎品發生問題")
		}

		if err := DefaultPrizeModel().
			SetConn(a.DbConn, a.RedisConn, a.MongoConn).
			Add(true, "tugofwar", utils.UUID(20), EditPrizeModel{
				ActivityID:    model.ActivityID,
				GameID:        gameid,
				PrizeName:     "落敗隊伍",
				PrizeType:     "thanks",
				PrizePicture:  config.UPLOAD_SYSTEM_URL + "img-prize-pic.png",
				PrizeAmount:   "0",
				PrizePrice:    "0",
				PrizeMethod:   "site",
				PrizePassword: "lose",
				TeamType:      "lose",
			}); err != nil {
			return errors.New("錯誤: 新增敗方遊戲獎品發生問題")
		}

	} else if game == "bingo" {

		// 建立賓果遊戲時，預設一個獎品
		if err := DefaultPrizeModel().
			SetConn(a.DbConn, a.RedisConn, a.MongoConn).
			Add(true, "bingo", utils.UUID(20), EditPrizeModel{
				ActivityID:    model.ActivityID,
				GameID:        gameid,
				PrizeName:     "賓果遊戲獎品",
				PrizeType:     "first",
				PrizePicture:  config.UPLOAD_SYSTEM_URL + "img-prize-pic.png",
				PrizeAmount:   "0",
				PrizePrice:    "0",
				PrizeMethod:   "site",
				PrizePassword: "bingo",
				TeamType:      "",
			}); err != nil {
			return errors.New("錯誤: 新增賓果遊戲獎品發生問題")
		}

	} else if game == "3DGachaMachine" {

	} else if game == "vote" {

	}

	return nil
}

// command.Value{
// 	"user_id":       model.UserID,
// 	"activity_id":   model.ActivityID,
// 	"game_id":       gameid,
// 	"game":          game,
// 	"title":         model.Title,
// 	"game_type":     model.GameType,
// 	"limit_time":    model.LimitTime,
// 	"second":        model.Second,
// 	"max_people":    model.MaxPeople,
// 	"people":        model.People,
// 	"max_times":     model.MaxTimes,
// 	"allow":         model.Allow,
// 	"percent":       model.Percent,
// 	"first_prize":   model.FirstPrize,
// 	"second_prize":  model.SecondPrize,
// 	"third_prize":   model.ThirdPrize,
// 	"general_prize": model.GeneralPrize,
// 	"topic":         model.Topic,
// 	"skin":          model.Skin,
// 	"music":         model.Music,
// 	"display_name":  model.DisplayName,
// 	"game_round":    1,
// 	"game_second":   model.Second,
// 	"game_status":   gameStatus,
// 	"game_attend":   0,

// 	"allow_choose_team":      model.AllowChooseTeam,
// 	"left_team_name":         model.LeftTeamName,
// 	"left_team_picture":      model.LeftTeamPicture,
// 	"right_team_name":        model.RightTeamName,
// 	"right_team_picture":     model.RightTeamPicture,
// 	"left_team_game_attend":  0,
// 	"right_team_game_attend": 0,
// 	"prize":                  model.Prize,

// 	"max_number":  model.MaxNumber,
// 	"bingo_line":  model.BingoLine,
// 	"round_prize": model.RoundPrize,
// 	"bingo_round": 0,

// 	"gacha_machine_reflection": model.GachaMachineReflection,
// 	"reflective_switch":        model.ReflectiveSwitch,

// 	"edit_times": 0,

// 	"game_order": len(games) + 1,

// 	"box_reflection": model.BoxReflection,
// 	"same_people":    model.SamePeople,

// 	"vote_screen":        model.VoteScreen,
// 	"vote_times":         model.VoteTimes,
// 	"vote_method":        model.VoteMethod,
// 	"vote_method_player": model.VoteMethodPlayer,
// 	"vote_restriction":   model.VoteRestriction,
// 	"avatar_shape":       model.AvatarShape,
// 	"vote_start_time":    model.VoteStartTime,
// 	"vote_end_time":      model.VoteEndTime,
// 	"auto_play":          model.AutoPlay,
// 	"show_rank":          model.ShowRank,
// 	"title_switch":       model.TitleSwitch,
// 	"arrangement_guest":  model.ArrangementGuest,
// }

// activity_game_ropepack_picture資料表
// if _, err := a.Table(config.ACTIVITY_GAME_ROPEPACK_PICTURE_TABLE).Insert(command.Value{
// 	"activity_id": model.ActivityID,
// 	"game_id":     gameid,

// 	// 套紅包自定義
// 	"ropepack_classic_h_pic_01": model.RopepackClassicHPic01,
// 	"ropepack_classic_h_pic_02": model.RopepackClassicHPic02,
// 	"ropepack_classic_h_pic_03": model.RopepackClassicHPic03,
// 	"ropepack_classic_h_pic_04": model.RopepackClassicHPic04,
// 	"ropepack_classic_h_pic_05": model.RopepackClassicHPic05,
// 	"ropepack_classic_h_pic_06": model.RopepackClassicHPic06,
// 	"ropepack_classic_h_pic_07": model.RopepackClassicHPic07,
// 	"ropepack_classic_h_pic_08": model.RopepackClassicHPic08,
// 	"ropepack_classic_h_pic_09": model.RopepackClassicHPic09,
// 	"ropepack_classic_h_pic_10": model.RopepackClassicHPic10,
// 	"ropepack_classic_g_pic_01": model.RopepackClassicGPic01,
// 	"ropepack_classic_g_pic_02": model.RopepackClassicGPic02,
// 	"ropepack_classic_g_pic_03": model.RopepackClassicGPic03,
// 	"ropepack_classic_g_pic_04": model.RopepackClassicGPic04,
// 	"ropepack_classic_g_pic_05": model.RopepackClassicGPic05,
// 	"ropepack_classic_g_pic_06": model.RopepackClassicGPic06,
// 	"ropepack_classic_h_ani_01": model.RopepackClassicHAni01,
// 	"ropepack_classic_g_ani_01": model.RopepackClassicGAni01,
// 	"ropepack_classic_g_ani_02": model.RopepackClassicGAni02,
// 	"ropepack_classic_c_ani_01": model.RopepackClassicCAni01,

// 	"ropepack_newyear_rabbit_h_pic_01": model.RopepackNewyearRabbitHPic01,
// 	"ropepack_newyear_rabbit_h_pic_02": model.RopepackNewyearRabbitHPic02,
// 	"ropepack_newyear_rabbit_h_pic_03": model.RopepackNewyearRabbitHPic03,
// 	"ropepack_newyear_rabbit_h_pic_04": model.RopepackNewyearRabbitHPic04,
// 	"ropepack_newyear_rabbit_h_pic_05": model.RopepackNewyearRabbitHPic05,
// 	"ropepack_newyear_rabbit_h_pic_06": model.RopepackNewyearRabbitHPic06,
// 	"ropepack_newyear_rabbit_h_pic_07": model.RopepackNewyearRabbitHPic07,
// 	"ropepack_newyear_rabbit_h_pic_08": model.RopepackNewyearRabbitHPic08,
// 	"ropepack_newyear_rabbit_h_pic_09": model.RopepackNewyearRabbitHPic09,
// 	"ropepack_newyear_rabbit_g_pic_01": model.RopepackNewyearRabbitGPic01,
// 	"ropepack_newyear_rabbit_g_pic_02": model.RopepackNewyearRabbitGPic02,
// 	"ropepack_newyear_rabbit_g_pic_03": model.RopepackNewyearRabbitGPic03,
// 	"ropepack_newyear_rabbit_h_ani_01": model.RopepackNewyearRabbitHAni01,
// 	"ropepack_newyear_rabbit_g_ani_01": model.RopepackNewyearRabbitGAni01,
// 	"ropepack_newyear_rabbit_g_ani_02": model.RopepackNewyearRabbitGAni02,
// 	"ropepack_newyear_rabbit_g_ani_03": model.RopepackNewyearRabbitGAni03,
// 	"ropepack_newyear_rabbit_c_ani_01": model.RopepackNewyearRabbitCAni01,
// 	"ropepack_newyear_rabbit_c_ani_02": model.RopepackNewyearRabbitCAni02,

// 	"ropepack_moonfestival_h_pic_01": model.RopepackMoonfestivalHPic01,
// 	"ropepack_moonfestival_h_pic_02": model.RopepackMoonfestivalHPic02,
// 	"ropepack_moonfestival_h_pic_03": model.RopepackMoonfestivalHPic03,
// 	"ropepack_moonfestival_h_pic_04": model.RopepackMoonfestivalHPic04,
// 	"ropepack_moonfestival_h_pic_05": model.RopepackMoonfestivalHPic05,
// 	"ropepack_moonfestival_h_pic_06": model.RopepackMoonfestivalHPic06,
// 	"ropepack_moonfestival_h_pic_07": model.RopepackMoonfestivalHPic07,
// 	"ropepack_moonfestival_h_pic_08": model.RopepackMoonfestivalHPic08,
// 	"ropepack_moonfestival_h_pic_09": model.RopepackMoonfestivalHPic09,
// 	"ropepack_moonfestival_g_pic_01": model.RopepackMoonfestivalGPic01,
// 	"ropepack_moonfestival_g_pic_02": model.RopepackMoonfestivalGPic02,
// 	"ropepack_moonfestival_c_pic_01": model.RopepackMoonfestivalCPic01,
// 	"ropepack_moonfestival_h_ani_01": model.RopepackMoonfestivalHAni01,
// 	"ropepack_moonfestival_g_ani_01": model.RopepackMoonfestivalGAni01,
// 	"ropepack_moonfestival_g_ani_02": model.RopepackMoonfestivalGAni02,
// 	"ropepack_moonfestival_c_ani_01": model.RopepackMoonfestivalCAni01,
// 	"ropepack_moonfestival_c_ani_02": model.RopepackMoonfestivalCAni02,

// 	"ropepack_3D_h_pic_01": model.Ropepack3DHPic01,
// 	"ropepack_3D_h_pic_02": model.Ropepack3DHPic02,
// 	"ropepack_3D_h_pic_03": model.Ropepack3DHPic03,
// 	"ropepack_3D_h_pic_04": model.Ropepack3DHPic04,
// 	"ropepack_3D_h_pic_05": model.Ropepack3DHPic05,
// 	"ropepack_3D_h_pic_06": model.Ropepack3DHPic06,
// 	"ropepack_3D_h_pic_07": model.Ropepack3DHPic07,
// 	"ropepack_3D_h_pic_08": model.Ropepack3DHPic08,
// 	"ropepack_3D_h_pic_09": model.Ropepack3DHPic09,
// 	"ropepack_3D_h_pic_10": model.Ropepack3DHPic10,
// 	"ropepack_3D_h_pic_11": model.Ropepack3DHPic11,
// 	"ropepack_3D_h_pic_12": model.Ropepack3DHPic12,
// 	"ropepack_3D_h_pic_13": model.Ropepack3DHPic13,
// 	"ropepack_3D_h_pic_14": model.Ropepack3DHPic14,
// 	"ropepack_3D_h_pic_15": model.Ropepack3DHPic15,
// 	"ropepack_3D_g_pic_01": model.Ropepack3DGPic01,
// 	"ropepack_3D_g_pic_02": model.Ropepack3DGPic02,
// 	"ropepack_3D_g_pic_03": model.Ropepack3DGPic03,
// 	"ropepack_3D_g_pic_04": model.Ropepack3DGPic04,
// 	"ropepack_3D_h_ani_01": model.Ropepack3DHAni01,
// 	"ropepack_3D_h_ani_02": model.Ropepack3DHAni02,
// 	"ropepack_3D_h_ani_03": model.Ropepack3DHAni03,
// 	"ropepack_3D_g_ani_01": model.Ropepack3DGAni01,
// 	"ropepack_3D_g_ani_02": model.Ropepack3DGAni02,
// 	"ropepack_3D_c_ani_01": model.Ropepack3DCAni01,

// 	// 音樂
// 	"ropepack_bgm_start":  model.RopepackBgmStart,
// 	"ropepack_bgm_gaming": model.RopepackBgmGaming,
// 	"ropepack_bgm_end":    model.RopepackBgmEnd,
// }); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_ropepack_picture)發生問題")
// }

// activity_game_monopoly_picture資料表
// if _, err := a.Table(config.ACTIVITY_GAME_MONOPOLY_PICTURE_TABLE).Insert(command.Value{
// 	"activity_id": model.ActivityID,
// 	"game_id":     gameid,

// 	// 鑑定師自定義
// 	"monopoly_classic_h_pic_01": model.MonopolyClassicHPic01,
// 	"monopoly_classic_h_pic_02": model.MonopolyClassicHPic02,
// 	"monopoly_classic_h_pic_03": model.MonopolyClassicHPic03,
// 	"monopoly_classic_h_pic_04": model.MonopolyClassicHPic04,
// 	"monopoly_classic_h_pic_05": model.MonopolyClassicHPic05,
// 	"monopoly_classic_h_pic_06": model.MonopolyClassicHPic06,
// 	"monopoly_classic_h_pic_07": model.MonopolyClassicHPic07,
// 	"monopoly_classic_h_pic_08": model.MonopolyClassicHPic08,
// 	"monopoly_classic_g_pic_01": model.MonopolyClassicGPic01,
// 	"monopoly_classic_g_pic_02": model.MonopolyClassicGPic02,
// 	"monopoly_classic_g_pic_03": model.MonopolyClassicGPic03,
// 	"monopoly_classic_g_pic_04": model.MonopolyClassicGPic04,
// 	"monopoly_classic_g_pic_05": model.MonopolyClassicGPic05,
// 	"monopoly_classic_g_pic_06": model.MonopolyClassicGPic06,
// 	"monopoly_classic_g_pic_07": model.MonopolyClassicGPic07,
// 	"monopoly_classic_c_pic_01": model.MonopolyClassicCPic01,
// 	"monopoly_classic_c_pic_02": model.MonopolyClassicCPic02,
// 	"monopoly_classic_g_ani_01": model.MonopolyClassicGAni01,
// 	"monopoly_classic_g_ani_02": model.MonopolyClassicGAni02,
// 	"monopoly_classic_c_ani_01": model.MonopolyClassicCAni01,

// 	"monopoly_redpack_h_pic_01": model.MonopolyRedpackHPic01,
// 	"monopoly_redpack_h_pic_02": model.MonopolyRedpackHPic02,
// 	"monopoly_redpack_h_pic_03": model.MonopolyRedpackHPic03,
// 	"monopoly_redpack_h_pic_04": model.MonopolyRedpackHPic04,
// 	"monopoly_redpack_h_pic_05": model.MonopolyRedpackHPic05,
// 	"monopoly_redpack_h_pic_06": model.MonopolyRedpackHPic06,
// 	"monopoly_redpack_h_pic_07": model.MonopolyRedpackHPic07,
// 	"monopoly_redpack_h_pic_08": model.MonopolyRedpackHPic08,
// 	"monopoly_redpack_h_pic_09": model.MonopolyRedpackHPic09,
// 	"monopoly_redpack_h_pic_10": model.MonopolyRedpackHPic10,
// 	"monopoly_redpack_h_pic_11": model.MonopolyRedpackHPic11,
// 	"monopoly_redpack_h_pic_12": model.MonopolyRedpackHPic12,
// 	"monopoly_redpack_h_pic_13": model.MonopolyRedpackHPic13,
// 	"monopoly_redpack_h_pic_14": model.MonopolyRedpackHPic14,
// 	"monopoly_redpack_h_pic_15": model.MonopolyRedpackHPic15,
// 	"monopoly_redpack_h_pic_16": model.MonopolyRedpackHPic16,
// 	"monopoly_redpack_g_pic_01": model.MonopolyRedpackGPic01,
// 	"monopoly_redpack_g_pic_02": model.MonopolyRedpackGPic02,
// 	"monopoly_redpack_g_pic_03": model.MonopolyRedpackGPic03,
// 	"monopoly_redpack_g_pic_04": model.MonopolyRedpackGPic04,
// 	"monopoly_redpack_g_pic_05": model.MonopolyRedpackGPic05,
// 	"monopoly_redpack_g_pic_06": model.MonopolyRedpackGPic06,
// 	"monopoly_redpack_g_pic_07": model.MonopolyRedpackGPic07,
// 	"monopoly_redpack_g_pic_08": model.MonopolyRedpackGPic08,
// 	"monopoly_redpack_g_pic_09": model.MonopolyRedpackGPic09,
// 	"monopoly_redpack_g_pic_10": model.MonopolyRedpackGPic10,
// 	"monopoly_redpack_c_pic_01": model.MonopolyRedpackCPic01,
// 	"monopoly_redpack_c_pic_02": model.MonopolyRedpackCPic02,
// 	"monopoly_redpack_c_pic_03": model.MonopolyRedpackCPic03,
// 	"monopoly_redpack_h_ani_01": model.MonopolyRedpackHAni01,
// 	"monopoly_redpack_h_ani_02": model.MonopolyRedpackHAni02,
// 	"monopoly_redpack_h_ani_03": model.MonopolyRedpackHAni03,
// 	"monopoly_redpack_g_ani_01": model.MonopolyRedpackGAni01,
// 	"monopoly_redpack_g_ani_02": model.MonopolyRedpackGAni02,
// 	"monopoly_redpack_c_ani_01": model.MonopolyRedpackCAni01,

// 	"monopoly_newyear_rabbit_h_pic_01": model.MonopolyNewyearRabbitHPic01,
// 	"monopoly_newyear_rabbit_h_pic_02": model.MonopolyNewyearRabbitHPic02,
// 	"monopoly_newyear_rabbit_h_pic_03": model.MonopolyNewyearRabbitHPic03,
// 	"monopoly_newyear_rabbit_h_pic_04": model.MonopolyNewyearRabbitHPic04,
// 	"monopoly_newyear_rabbit_h_pic_05": model.MonopolyNewyearRabbitHPic05,
// 	"monopoly_newyear_rabbit_h_pic_06": model.MonopolyNewyearRabbitHPic06,
// 	"monopoly_newyear_rabbit_h_pic_07": model.MonopolyNewyearRabbitHPic07,
// 	"monopoly_newyear_rabbit_h_pic_08": model.MonopolyNewyearRabbitHPic08,
// 	"monopoly_newyear_rabbit_h_pic_09": model.MonopolyNewyearRabbitHPic09,
// 	"monopoly_newyear_rabbit_h_pic_10": model.MonopolyNewyearRabbitHPic10,
// 	"monopoly_newyear_rabbit_h_pic_11": model.MonopolyNewyearRabbitHPic11,
// 	"monopoly_newyear_rabbit_h_pic_12": model.MonopolyNewyearRabbitHPic12,
// 	"monopoly_newyear_rabbit_g_pic_01": model.MonopolyNewyearRabbitGPic01,
// 	"monopoly_newyear_rabbit_g_pic_02": model.MonopolyNewyearRabbitGPic02,
// 	"monopoly_newyear_rabbit_g_pic_03": model.MonopolyNewyearRabbitGPic03,
// 	"monopoly_newyear_rabbit_g_pic_04": model.MonopolyNewyearRabbitGPic04,
// 	"monopoly_newyear_rabbit_g_pic_05": model.MonopolyNewyearRabbitGPic05,
// 	"monopoly_newyear_rabbit_g_pic_06": model.MonopolyNewyearRabbitGPic06,
// 	"monopoly_newyear_rabbit_g_pic_07": model.MonopolyNewyearRabbitGPic07,
// 	"monopoly_newyear_rabbit_c_pic_01": model.MonopolyNewyearRabbitCPic01,
// 	"monopoly_newyear_rabbit_c_pic_02": model.MonopolyNewyearRabbitCPic02,
// 	"monopoly_newyear_rabbit_c_pic_03": model.MonopolyNewyearRabbitCPic03,
// 	"monopoly_newyear_rabbit_h_ani_01": model.MonopolyNewyearRabbitHAni01,
// 	"monopoly_newyear_rabbit_h_ani_02": model.MonopolyNewyearRabbitHAni02,
// 	"monopoly_newyear_rabbit_g_ani_01": model.MonopolyNewyearRabbitGAni01,
// 	"monopoly_newyear_rabbit_g_ani_02": model.MonopolyNewyearRabbitGAni02,
// 	"monopoly_newyear_rabbit_c_ani_01": model.MonopolyNewyearRabbitCAni01,

// 	"monopoly_sashimi_h_pic_01": model.MonopolySashimiHPic01,
// 	"monopoly_sashimi_h_pic_02": model.MonopolySashimiHPic02,
// 	"monopoly_sashimi_h_pic_03": model.MonopolySashimiHPic03,
// 	"monopoly_sashimi_h_pic_04": model.MonopolySashimiHPic04,
// 	"monopoly_sashimi_h_pic_05": model.MonopolySashimiHPic05,
// 	"monopoly_sashimi_g_pic_01": model.MonopolySashimiGPic01,
// 	"monopoly_sashimi_g_pic_02": model.MonopolySashimiGPic02,
// 	"monopoly_sashimi_g_pic_03": model.MonopolySashimiGPic03,
// 	"monopoly_sashimi_g_pic_04": model.MonopolySashimiGPic04,
// 	"monopoly_sashimi_g_pic_05": model.MonopolySashimiGPic05,
// 	"monopoly_sashimi_g_pic_06": model.MonopolySashimiGPic06,
// 	"monopoly_sashimi_c_pic_01": model.MonopolySashimiCPic01,
// 	"monopoly_sashimi_c_pic_02": model.MonopolySashimiCPic02,
// 	"monopoly_sashimi_h_ani_01": model.MonopolySashimiHAni01,
// 	"monopoly_sashimi_h_ani_02": model.MonopolySashimiHAni02,
// 	"monopoly_sashimi_g_ani_01": model.MonopolySashimiGAni01,
// 	"monopoly_sashimi_g_ani_02": model.MonopolySashimiGAni02,
// 	"monopoly_sashimi_c_ani_01": model.MonopolySashimiCAni01,

// 	// 音樂
// 	"monopoly_bgm_start":  model.MonopolyBgmStart,
// 	"monopoly_bgm_gaming": model.MonopolyBgmGaming,
// 	"monopoly_bgm_end":    model.MonopolyBgmEnd,
// }); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_monopoly_picture)發生問題")
// }

// activity_game_draw_numbers_picture資料表
// if _, err := a.Table(config.ACTIVITY_GAME_DRAW_NUMBERS_PICTURE_TABLE).Insert(command.Value{
// 	"activity_id": model.ActivityID,
// 	"game_id":     gameid,

// 	// 搖號抽獎自定義
// 	"draw_numbers_classic_h_pic_01": model.DrawNumbersClassicHPic01,
// 	"draw_numbers_classic_h_pic_02": model.DrawNumbersClassicHPic02,
// 	"draw_numbers_classic_h_pic_03": model.DrawNumbersClassicHPic03,
// 	"draw_numbers_classic_h_pic_04": model.DrawNumbersClassicHPic04,
// 	"draw_numbers_classic_h_pic_05": model.DrawNumbersClassicHPic05,
// 	"draw_numbers_classic_h_pic_06": model.DrawNumbersClassicHPic06,
// 	"draw_numbers_classic_h_pic_07": model.DrawNumbersClassicHPic07,
// 	"draw_numbers_classic_h_pic_08": model.DrawNumbersClassicHPic08,
// 	"draw_numbers_classic_h_pic_09": model.DrawNumbersClassicHPic09,
// 	"draw_numbers_classic_h_pic_10": model.DrawNumbersClassicHPic10,
// 	"draw_numbers_classic_h_pic_11": model.DrawNumbersClassicHPic11,
// 	"draw_numbers_classic_h_pic_12": model.DrawNumbersClassicHPic12,
// 	"draw_numbers_classic_h_pic_13": model.DrawNumbersClassicHPic13,
// 	"draw_numbers_classic_h_pic_14": model.DrawNumbersClassicHPic14,
// 	"draw_numbers_classic_h_pic_15": model.DrawNumbersClassicHPic15,
// 	"draw_numbers_classic_h_pic_16": model.DrawNumbersClassicHPic16,
// 	"draw_numbers_classic_h_ani_01": model.DrawNumbersClassicHAni01,

// 	"draw_numbers_gold_h_pic_01": model.DrawNumbersGoldHPic01,
// 	"draw_numbers_gold_h_pic_02": model.DrawNumbersGoldHPic02,
// 	"draw_numbers_gold_h_pic_03": model.DrawNumbersGoldHPic03,
// 	"draw_numbers_gold_h_pic_04": model.DrawNumbersGoldHPic04,
// 	"draw_numbers_gold_h_pic_05": model.DrawNumbersGoldHPic05,
// 	"draw_numbers_gold_h_pic_06": model.DrawNumbersGoldHPic06,
// 	"draw_numbers_gold_h_pic_07": model.DrawNumbersGoldHPic07,
// 	"draw_numbers_gold_h_pic_08": model.DrawNumbersGoldHPic08,
// 	"draw_numbers_gold_h_pic_09": model.DrawNumbersGoldHPic09,
// 	"draw_numbers_gold_h_pic_10": model.DrawNumbersGoldHPic10,
// 	"draw_numbers_gold_h_pic_11": model.DrawNumbersGoldHPic11,
// 	"draw_numbers_gold_h_pic_12": model.DrawNumbersGoldHPic12,
// 	"draw_numbers_gold_h_pic_13": model.DrawNumbersGoldHPic13,
// 	"draw_numbers_gold_h_pic_14": model.DrawNumbersGoldHPic14,
// 	"draw_numbers_gold_h_ani_01": model.DrawNumbersGoldHAni01,
// 	"draw_numbers_gold_h_ani_02": model.DrawNumbersGoldHAni02,
// 	"draw_numbers_gold_h_ani_03": model.DrawNumbersGoldHAni03,

// 	"draw_numbers_newyear_dragon_h_pic_01": model.DrawNumbersNewyearDragonHPic01,
// 	"draw_numbers_newyear_dragon_h_pic_02": model.DrawNumbersNewyearDragonHPic02,
// 	"draw_numbers_newyear_dragon_h_pic_03": model.DrawNumbersNewyearDragonHPic03,
// 	"draw_numbers_newyear_dragon_h_pic_04": model.DrawNumbersNewyearDragonHPic04,
// 	"draw_numbers_newyear_dragon_h_pic_05": model.DrawNumbersNewyearDragonHPic05,
// 	"draw_numbers_newyear_dragon_h_pic_06": model.DrawNumbersNewyearDragonHPic06,
// 	"draw_numbers_newyear_dragon_h_pic_07": model.DrawNumbersNewyearDragonHPic07,
// 	"draw_numbers_newyear_dragon_h_pic_08": model.DrawNumbersNewyearDragonHPic08,
// 	"draw_numbers_newyear_dragon_h_pic_09": model.DrawNumbersNewyearDragonHPic09,
// 	"draw_numbers_newyear_dragon_h_pic_10": model.DrawNumbersNewyearDragonHPic10,
// 	"draw_numbers_newyear_dragon_h_pic_11": model.DrawNumbersNewyearDragonHPic11,
// 	"draw_numbers_newyear_dragon_h_pic_12": model.DrawNumbersNewyearDragonHPic12,
// 	"draw_numbers_newyear_dragon_h_pic_13": model.DrawNumbersNewyearDragonHPic13,
// 	"draw_numbers_newyear_dragon_h_pic_14": model.DrawNumbersNewyearDragonHPic14,
// 	"draw_numbers_newyear_dragon_h_pic_15": model.DrawNumbersNewyearDragonHPic15,
// 	"draw_numbers_newyear_dragon_h_pic_16": model.DrawNumbersNewyearDragonHPic16,
// 	"draw_numbers_newyear_dragon_h_pic_17": model.DrawNumbersNewyearDragonHPic17,
// 	"draw_numbers_newyear_dragon_h_pic_18": model.DrawNumbersNewyearDragonHPic18,
// 	"draw_numbers_newyear_dragon_h_pic_19": model.DrawNumbersNewyearDragonHPic19,
// 	"draw_numbers_newyear_dragon_h_pic_20": model.DrawNumbersNewyearDragonHPic20,
// 	"draw_numbers_newyear_dragon_h_ani_01": model.DrawNumbersNewyearDragonHAni01,
// 	"draw_numbers_newyear_dragon_h_ani_02": model.DrawNumbersNewyearDragonHAni02,

// 	"draw_numbers_cherry_h_pic_01": model.DrawNumbersCherryHPic01,
// 	"draw_numbers_cherry_h_pic_02": model.DrawNumbersCherryHPic02,
// 	"draw_numbers_cherry_h_pic_03": model.DrawNumbersCherryHPic03,
// 	"draw_numbers_cherry_h_pic_04": model.DrawNumbersCherryHPic04,
// 	"draw_numbers_cherry_h_pic_05": model.DrawNumbersCherryHPic05,
// 	"draw_numbers_cherry_h_pic_06": model.DrawNumbersCherryHPic06,
// 	"draw_numbers_cherry_h_pic_07": model.DrawNumbersCherryHPic07,
// 	"draw_numbers_cherry_h_pic_08": model.DrawNumbersCherryHPic08,
// 	"draw_numbers_cherry_h_pic_09": model.DrawNumbersCherryHPic09,
// 	"draw_numbers_cherry_h_pic_10": model.DrawNumbersCherryHPic10,
// 	"draw_numbers_cherry_h_pic_11": model.DrawNumbersCherryHPic11,
// 	"draw_numbers_cherry_h_pic_12": model.DrawNumbersCherryHPic12,
// 	"draw_numbers_cherry_h_pic_13": model.DrawNumbersCherryHPic13,
// 	"draw_numbers_cherry_h_pic_14": model.DrawNumbersCherryHPic14,
// 	"draw_numbers_cherry_h_pic_15": model.DrawNumbersCherryHPic15,
// 	"draw_numbers_cherry_h_pic_16": model.DrawNumbersCherryHPic16,
// 	"draw_numbers_cherry_h_pic_17": model.DrawNumbersCherryHPic17,
// 	"draw_numbers_cherry_h_ani_01": model.DrawNumbersCherryHAni01,
// 	"draw_numbers_cherry_h_ani_02": model.DrawNumbersCherryHAni02,
// 	"draw_numbers_cherry_h_ani_03": model.DrawNumbersCherryHAni03,
// 	"draw_numbers_cherry_h_ani_04": model.DrawNumbersCherryHAni04,

// 	"draw_numbers_3D_space_h_pic_01": model.DrawNumbers3DSpaceHPic01,
// 	"draw_numbers_3D_space_h_pic_02": model.DrawNumbers3DSpaceHPic02,
// 	"draw_numbers_3D_space_h_pic_03": model.DrawNumbers3DSpaceHPic03,
// 	"draw_numbers_3D_space_h_pic_04": model.DrawNumbers3DSpaceHPic04,
// 	"draw_numbers_3D_space_h_pic_05": model.DrawNumbers3DSpaceHPic05,
// 	"draw_numbers_3D_space_h_pic_06": model.DrawNumbers3DSpaceHPic06,
// 	"draw_numbers_3D_space_h_pic_07": model.DrawNumbers3DSpaceHPic07,
// 	"draw_numbers_3D_space_h_pic_08": model.DrawNumbers3DSpaceHPic08,

// 	// 音樂
// 	"draw_numbers_bgm_gaming": model.DrawNumbersBgmGaming,
// }); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_draw_numbers_picture)發生問題")
// }

// activity_game_tugofwar_picture資料表
// if _, err := a.Table(config.ACTIVITY_GAME_TUGOFWAR_PICTURE_TABLE).Insert(command.Value{
// 	"activity_id": model.ActivityID,
// 	"game_id":     gameid,

// 	// 拔河遊戲自定義
// 	"tugofwar_classic_h_pic_01": model.TugofwarClassicHPic01,
// 	"tugofwar_classic_h_pic_02": model.TugofwarClassicHPic02,
// 	"tugofwar_classic_h_pic_03": model.TugofwarClassicHPic03,
// 	"tugofwar_classic_h_pic_04": model.TugofwarClassicHPic04,
// 	"tugofwar_classic_h_pic_05": model.TugofwarClassicHPic05,
// 	"tugofwar_classic_h_pic_06": model.TugofwarClassicHPic06,
// 	"tugofwar_classic_h_pic_07": model.TugofwarClassicHPic07,
// 	"tugofwar_classic_h_pic_08": model.TugofwarClassicHPic08,
// 	"tugofwar_classic_h_pic_09": model.TugofwarClassicHPic09,
// 	"tugofwar_classic_h_pic_10": model.TugofwarClassicHPic10,
// 	"tugofwar_classic_h_pic_11": model.TugofwarClassicHPic11,
// 	"tugofwar_classic_h_pic_12": model.TugofwarClassicHPic12,
// 	"tugofwar_classic_h_pic_13": model.TugofwarClassicHPic13,
// 	"tugofwar_classic_h_pic_14": model.TugofwarClassicHPic14,
// 	"tugofwar_classic_h_pic_15": model.TugofwarClassicHPic15,
// 	"tugofwar_classic_h_pic_16": model.TugofwarClassicHPic16,
// 	"tugofwar_classic_h_pic_17": model.TugofwarClassicHPic17,
// 	"tugofwar_classic_h_pic_18": model.TugofwarClassicHPic18,
// 	"tugofwar_classic_h_pic_19": model.TugofwarClassicHPic19,
// 	"tugofwar_classic_g_pic_01": model.TugofwarClassicGPic01,
// 	"tugofwar_classic_g_pic_02": model.TugofwarClassicGPic02,
// 	"tugofwar_classic_g_pic_03": model.TugofwarClassicGPic03,
// 	"tugofwar_classic_g_pic_04": model.TugofwarClassicGPic04,
// 	"tugofwar_classic_g_pic_05": model.TugofwarClassicGPic05,
// 	"tugofwar_classic_g_pic_06": model.TugofwarClassicGPic06,
// 	"tugofwar_classic_g_pic_07": model.TugofwarClassicGPic07,
// 	"tugofwar_classic_g_pic_08": model.TugofwarClassicGPic08,
// 	"tugofwar_classic_g_pic_09": model.TugofwarClassicGPic09,
// 	"tugofwar_classic_h_ani_01": model.TugofwarClassicHAni01,
// 	"tugofwar_classic_h_ani_02": model.TugofwarClassicHAni02,
// 	"tugofwar_classic_h_ani_03": model.TugofwarClassicHAni03,
// 	"tugofwar_classic_c_ani_01": model.TugofwarClassicCAni01,

// 	"tugofwar_school_h_pic_01": model.TugofwarSchoolHPic01,
// 	"tugofwar_school_h_pic_02": model.TugofwarSchoolHPic02,
// 	"tugofwar_school_h_pic_03": model.TugofwarSchoolHPic03,
// 	"tugofwar_school_h_pic_04": model.TugofwarSchoolHPic04,
// 	"tugofwar_school_h_pic_05": model.TugofwarSchoolHPic05,
// 	"tugofwar_school_h_pic_06": model.TugofwarSchoolHPic06,
// 	"tugofwar_school_h_pic_07": model.TugofwarSchoolHPic07,
// 	"tugofwar_school_h_pic_08": model.TugofwarSchoolHPic08,
// 	"tugofwar_school_h_pic_09": model.TugofwarSchoolHPic09,
// 	"tugofwar_school_h_pic_10": model.TugofwarSchoolHPic10,
// 	"tugofwar_school_h_pic_11": model.TugofwarSchoolHPic11,
// 	"tugofwar_school_h_pic_12": model.TugofwarSchoolHPic12,
// 	"tugofwar_school_h_pic_13": model.TugofwarSchoolHPic13,
// 	"tugofwar_school_h_pic_14": model.TugofwarSchoolHPic14,
// 	"tugofwar_school_h_pic_15": model.TugofwarSchoolHPic15,
// 	"tugofwar_school_h_pic_16": model.TugofwarSchoolHPic16,
// 	"tugofwar_school_h_pic_17": model.TugofwarSchoolHPic17,
// 	"tugofwar_school_h_pic_18": model.TugofwarSchoolHPic18,
// 	"tugofwar_school_h_pic_19": model.TugofwarSchoolHPic19,
// 	"tugofwar_school_h_pic_20": model.TugofwarSchoolHPic20,
// 	"tugofwar_school_h_pic_21": model.TugofwarSchoolHPic21,
// 	"tugofwar_school_h_pic_22": model.TugofwarSchoolHPic22,
// 	"tugofwar_school_h_pic_23": model.TugofwarSchoolHPic23,
// 	"tugofwar_school_h_pic_24": model.TugofwarSchoolHPic24,
// 	"tugofwar_school_h_pic_25": model.TugofwarSchoolHPic25,
// 	"tugofwar_school_h_pic_26": model.TugofwarSchoolHPic26,
// 	"tugofwar_school_g_pic_01": model.TugofwarSchoolGPic01,
// 	"tugofwar_school_g_pic_02": model.TugofwarSchoolGPic02,
// 	"tugofwar_school_g_pic_03": model.TugofwarSchoolGPic03,
// 	"tugofwar_school_g_pic_04": model.TugofwarSchoolGPic04,
// 	"tugofwar_school_g_pic_05": model.TugofwarSchoolGPic05,
// 	"tugofwar_school_g_pic_06": model.TugofwarSchoolGPic06,
// 	"tugofwar_school_g_pic_07": model.TugofwarSchoolGPic07,
// 	"tugofwar_school_c_pic_01": model.TugofwarSchoolCPic01,
// 	"tugofwar_school_c_pic_02": model.TugofwarSchoolCPic02,
// 	"tugofwar_school_c_pic_03": model.TugofwarSchoolCPic03,
// 	"tugofwar_school_c_pic_04": model.TugofwarSchoolCPic04,
// 	"tugofwar_school_h_ani_01": model.TugofwarSchoolHAni01,
// 	"tugofwar_school_h_ani_02": model.TugofwarSchoolHAni02,
// 	"tugofwar_school_h_ani_03": model.TugofwarSchoolHAni03,
// 	"tugofwar_school_h_ani_04": model.TugofwarSchoolHAni04,
// 	"tugofwar_school_h_ani_05": model.TugofwarSchoolHAni05,
// 	"tugofwar_school_h_ani_06": model.TugofwarSchoolHAni06,
// 	"tugofwar_school_h_ani_07": model.TugofwarSchoolHAni07,

// 	"tugofwar_christmas_h_pic_01": model.TugofwarChristmasHPic01,
// 	"tugofwar_christmas_h_pic_02": model.TugofwarChristmasHPic02,
// 	"tugofwar_christmas_h_pic_03": model.TugofwarChristmasHPic03,
// 	"tugofwar_christmas_h_pic_04": model.TugofwarChristmasHPic04,
// 	"tugofwar_christmas_h_pic_05": model.TugofwarChristmasHPic05,
// 	"tugofwar_christmas_h_pic_06": model.TugofwarChristmasHPic06,
// 	"tugofwar_christmas_h_pic_07": model.TugofwarChristmasHPic07,
// 	"tugofwar_christmas_h_pic_08": model.TugofwarChristmasHPic08,
// 	"tugofwar_christmas_h_pic_09": model.TugofwarChristmasHPic09,
// 	"tugofwar_christmas_h_pic_10": model.TugofwarChristmasHPic10,
// 	"tugofwar_christmas_h_pic_11": model.TugofwarChristmasHPic11,
// 	"tugofwar_christmas_h_pic_12": model.TugofwarChristmasHPic12,
// 	"tugofwar_christmas_h_pic_13": model.TugofwarChristmasHPic13,
// 	"tugofwar_christmas_h_pic_14": model.TugofwarChristmasHPic14,
// 	"tugofwar_christmas_h_pic_15": model.TugofwarChristmasHPic15,
// 	"tugofwar_christmas_h_pic_16": model.TugofwarChristmasHPic16,
// 	"tugofwar_christmas_h_pic_17": model.TugofwarChristmasHPic17,
// 	"tugofwar_christmas_h_pic_18": model.TugofwarChristmasHPic18,
// 	"tugofwar_christmas_h_pic_19": model.TugofwarChristmasHPic19,
// 	"tugofwar_christmas_h_pic_20": model.TugofwarChristmasHPic20,
// 	"tugofwar_christmas_h_pic_21": model.TugofwarChristmasHPic21,
// 	"tugofwar_christmas_g_pic_01": model.TugofwarChristmasGPic01,
// 	"tugofwar_christmas_g_pic_02": model.TugofwarChristmasGPic02,
// 	"tugofwar_christmas_g_pic_03": model.TugofwarChristmasGPic03,
// 	"tugofwar_christmas_g_pic_04": model.TugofwarChristmasGPic04,
// 	"tugofwar_christmas_g_pic_05": model.TugofwarChristmasGPic05,
// 	"tugofwar_christmas_g_pic_06": model.TugofwarChristmasGPic06,
// 	"tugofwar_christmas_c_pic_01": model.TugofwarChristmasCPic01,
// 	"tugofwar_christmas_c_pic_02": model.TugofwarChristmasCPic02,
// 	"tugofwar_christmas_c_pic_03": model.TugofwarChristmasCPic03,
// 	"tugofwar_christmas_c_pic_04": model.TugofwarChristmasCPic04,
// 	"tugofwar_christmas_h_ani_01": model.TugofwarChristmasHAni01,
// 	"tugofwar_christmas_h_ani_02": model.TugofwarChristmasHAni02,
// 	"tugofwar_christmas_h_ani_03": model.TugofwarChristmasHAni03,
// 	"tugofwar_christmas_c_ani_01": model.TugofwarChristmasCAni01,
// 	"tugofwar_christmas_c_ani_02": model.TugofwarChristmasCAni02,

// 	// 音樂
// 	"tugofwar_bgm_start":  model.TugofwarBgmStart,  // 遊戲開始
// 	"tugofwar_bgm_gaming": model.TugofwarBgmGaming, // 遊戲進行中
// 	"tugofwar_bgm_end":    model.TugofwarBgmEnd,    // 遊戲結束
// }); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_tugofwar_picture)發生問題")
// }

// activity_game_qa_picture資料表
// if _, err := a.Table(config.ACTIVITY_GAME_QA_PICTURE_TABLE_1).Insert(command.Value{
// 	"activity_id": model.ActivityID,
// 	"game_id":     gameid,

// 	// 快問快答自定義
// 	"qa_classic_h_pic_01": model.QAClassicHPic01,
// 	"qa_classic_h_pic_02": model.QAClassicHPic02,
// 	"qa_classic_h_pic_03": model.QAClassicHPic03,
// 	"qa_classic_h_pic_04": model.QAClassicHPic04,
// 	"qa_classic_h_pic_05": model.QAClassicHPic05,
// 	"qa_classic_h_pic_06": model.QAClassicHPic06,
// 	"qa_classic_h_pic_07": model.QAClassicHPic07,
// 	"qa_classic_h_pic_08": model.QAClassicHPic08,
// 	"qa_classic_h_pic_09": model.QAClassicHPic09,
// 	"qa_classic_h_pic_10": model.QAClassicHPic10,
// 	"qa_classic_h_pic_11": model.QAClassicHPic11,
// 	"qa_classic_h_pic_12": model.QAClassicHPic12,
// 	"qa_classic_h_pic_13": model.QAClassicHPic13,
// 	"qa_classic_h_pic_14": model.QAClassicHPic14,
// 	"qa_classic_h_pic_15": model.QAClassicHPic15,
// 	"qa_classic_h_pic_16": model.QAClassicHPic16,
// 	"qa_classic_h_pic_17": model.QAClassicHPic17,
// 	"qa_classic_h_pic_18": model.QAClassicHPic18,
// 	"qa_classic_h_pic_19": model.QAClassicHPic19,
// 	"qa_classic_h_pic_20": model.QAClassicHPic20,
// 	"qa_classic_h_pic_21": model.QAClassicHPic21,
// 	"qa_classic_h_pic_22": model.QAClassicHPic22,
// 	"qa_classic_g_pic_01": model.QAClassicGPic01,
// 	"qa_classic_g_pic_02": model.QAClassicGPic02,
// 	"qa_classic_g_pic_03": model.QAClassicGPic03,
// 	"qa_classic_g_pic_04": model.QAClassicGPic04,
// 	"qa_classic_g_pic_05": model.QAClassicGPic05,
// 	"qa_classic_c_pic_01": model.QAClassicCPic01,
// 	"qa_classic_h_ani_01": model.QAClassicHAni01,
// 	"qa_classic_h_ani_02": model.QAClassicHAni02,
// 	"qa_classic_g_ani_01": model.QAClassicGAni01,
// 	"qa_classic_g_ani_02": model.QAClassicGAni02,

// 	"qa_electric_h_pic_01": model.QAElectricHPic01,
// 	"qa_electric_h_pic_02": model.QAElectricHPic02,
// 	"qa_electric_h_pic_03": model.QAElectricHPic03,
// 	"qa_electric_h_pic_04": model.QAElectricHPic04,
// 	"qa_electric_h_pic_05": model.QAElectricHPic05,
// 	"qa_electric_h_pic_06": model.QAElectricHPic06,
// 	"qa_electric_h_pic_07": model.QAElectricHPic07,
// 	"qa_electric_h_pic_08": model.QAElectricHPic08,
// 	"qa_electric_h_pic_09": model.QAElectricHPic09,
// 	"qa_electric_h_pic_10": model.QAElectricHPic10,
// 	"qa_electric_h_pic_11": model.QAElectricHPic11,
// 	"qa_electric_h_pic_12": model.QAElectricHPic12,
// 	"qa_electric_h_pic_13": model.QAElectricHPic13,
// 	"qa_electric_h_pic_14": model.QAElectricHPic14,
// 	"qa_electric_h_pic_15": model.QAElectricHPic15,
// 	"qa_electric_h_pic_16": model.QAElectricHPic16,
// 	"qa_electric_h_pic_17": model.QAElectricHPic17,
// 	"qa_electric_h_pic_18": model.QAElectricHPic18,
// 	"qa_electric_h_pic_19": model.QAElectricHPic19,
// 	"qa_electric_h_pic_20": model.QAElectricHPic20,
// 	"qa_electric_h_pic_21": model.QAElectricHPic21,
// 	"qa_electric_h_pic_22": model.QAElectricHPic22,
// 	"qa_electric_h_pic_23": model.QAElectricHPic23,
// 	"qa_electric_h_pic_24": model.QAElectricHPic24,
// 	"qa_electric_h_pic_25": model.QAElectricHPic25,
// 	"qa_electric_h_pic_26": model.QAElectricHPic26,
// 	"qa_electric_g_pic_01": model.QAElectricGPic01,
// 	"qa_electric_g_pic_02": model.QAElectricGPic02,
// 	"qa_electric_g_pic_03": model.QAElectricGPic03,
// 	"qa_electric_g_pic_04": model.QAElectricGPic04,
// 	"qa_electric_g_pic_05": model.QAElectricGPic05,
// 	"qa_electric_g_pic_06": model.QAElectricGPic06,
// 	"qa_electric_g_pic_07": model.QAElectricGPic07,
// 	"qa_electric_g_pic_08": model.QAElectricGPic08,
// 	"qa_electric_g_pic_09": model.QAElectricGPic09,
// 	"qa_electric_c_pic_01": model.QAElectricCPic01,
// 	"qa_electric_h_ani_01": model.QAElectricHAni01,
// 	"qa_electric_h_ani_02": model.QAElectricHAni02,
// 	"qa_electric_h_ani_03": model.QAElectricHAni03,
// 	"qa_electric_h_ani_04": model.QAElectricHAni04,
// 	"qa_electric_h_ani_05": model.QAElectricHAni05,
// 	"qa_electric_g_ani_01": model.QAElectricGAni01,
// 	"qa_electric_g_ani_02": model.QAElectricGAni02,
// 	"qa_electric_c_ani_01": model.QAElectricCAni01,

// 	"qa_moonfestival_h_pic_01": model.QAMoonfestivalHPic01,
// 	"qa_moonfestival_h_pic_02": model.QAMoonfestivalHPic02,
// 	"qa_moonfestival_h_pic_03": model.QAMoonfestivalHPic03,
// 	"qa_moonfestival_h_pic_04": model.QAMoonfestivalHPic04,
// 	"qa_moonfestival_h_pic_05": model.QAMoonfestivalHPic05,
// 	"qa_moonfestival_h_pic_06": model.QAMoonfestivalHPic06,
// 	"qa_moonfestival_h_pic_07": model.QAMoonfestivalHPic07,
// 	"qa_moonfestival_h_pic_08": model.QAMoonfestivalHPic08,
// 	"qa_moonfestival_h_pic_09": model.QAMoonfestivalHPic09,
// 	"qa_moonfestival_h_pic_10": model.QAMoonfestivalHPic10,
// 	"qa_moonfestival_h_pic_11": model.QAMoonfestivalHPic11,
// 	"qa_moonfestival_h_pic_12": model.QAMoonfestivalHPic12,
// 	"qa_moonfestival_h_pic_13": model.QAMoonfestivalHPic13,
// 	"qa_moonfestival_h_pic_14": model.QAMoonfestivalHPic14,
// 	"qa_moonfestival_h_pic_15": model.QAMoonfestivalHPic15,
// 	"qa_moonfestival_h_pic_16": model.QAMoonfestivalHPic16,
// 	"qa_moonfestival_h_pic_17": model.QAMoonfestivalHPic17,
// 	"qa_moonfestival_h_pic_18": model.QAMoonfestivalHPic18,
// 	"qa_moonfestival_h_pic_19": model.QAMoonfestivalHPic19,
// 	"qa_moonfestival_h_pic_20": model.QAMoonfestivalHPic20,
// 	"qa_moonfestival_h_pic_21": model.QAMoonfestivalHPic21,
// 	"qa_moonfestival_h_pic_22": model.QAMoonfestivalHPic22,
// 	"qa_moonfestival_h_pic_23": model.QAMoonfestivalHPic23,
// 	"qa_moonfestival_h_pic_24": model.QAMoonfestivalHPic24,
// 	"qa_moonfestival_g_pic_01": model.QAMoonfestivalGPic01,
// 	"qa_moonfestival_g_pic_02": model.QAMoonfestivalGPic02,
// 	"qa_moonfestival_g_pic_03": model.QAMoonfestivalGPic03,
// 	"qa_moonfestival_g_pic_04": model.QAMoonfestivalGPic04,
// 	"qa_moonfestival_g_pic_05": model.QAMoonfestivalGPic05,
// 	"qa_moonfestival_c_pic_01": model.QAMoonfestivalCPic01,
// 	"qa_moonfestival_c_pic_02": model.QAMoonfestivalCPic02,
// 	"qa_moonfestival_c_pic_03": model.QAMoonfestivalCPic03,
// 	"qa_moonfestival_h_ani_01": model.QAMoonfestivalHAni01,
// 	"qa_moonfestival_h_ani_02": model.QAMoonfestivalHAni02,
// 	"qa_moonfestival_g_ani_01": model.QAMoonfestivalGAni01,
// 	"qa_moonfestival_g_ani_02": model.QAMoonfestivalGAni02,
// 	"qa_moonfestival_g_ani_03": model.QAMoonfestivalGAni03,

// 	// 音樂
// 	"qa_bgm_start":  model.QABgmStart,  // 遊戲開始
// 	"qa_bgm_gaming": model.QABgmGaming, // 遊戲進行中
// 	"qa_bgm_end":    model.QABgmEnd,    // 遊戲結束
// }); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_qa_picture)發生問題")
// }

// activity_game_qa_picture_2資料表
// if _, err := a.Table(config.ACTIVITY_GAME_QA_PICTURE_TABLE_2).Insert(command.Value{
// 	"activity_id": model.ActivityID,
// 	"game_id":     gameid,

// 	"qa_newyear_dragon_h_pic_01": model.QANewyearDragonHPic01,
// 	"qa_newyear_dragon_h_pic_02": model.QANewyearDragonHPic02,
// 	"qa_newyear_dragon_h_pic_03": model.QANewyearDragonHPic03,
// 	"qa_newyear_dragon_h_pic_04": model.QANewyearDragonHPic04,
// 	"qa_newyear_dragon_h_pic_05": model.QANewyearDragonHPic05,
// 	"qa_newyear_dragon_h_pic_06": model.QANewyearDragonHPic06,
// 	"qa_newyear_dragon_h_pic_07": model.QANewyearDragonHPic07,
// 	"qa_newyear_dragon_h_pic_08": model.QANewyearDragonHPic08,
// 	"qa_newyear_dragon_h_pic_09": model.QANewyearDragonHPic09,
// 	"qa_newyear_dragon_h_pic_10": model.QANewyearDragonHPic10,
// 	"qa_newyear_dragon_h_pic_11": model.QANewyearDragonHPic11,
// 	"qa_newyear_dragon_h_pic_12": model.QANewyearDragonHPic12,
// 	"qa_newyear_dragon_h_pic_13": model.QANewyearDragonHPic13,
// 	"qa_newyear_dragon_h_pic_14": model.QANewyearDragonHPic14,
// 	"qa_newyear_dragon_h_pic_15": model.QANewyearDragonHPic15,
// 	"qa_newyear_dragon_h_pic_16": model.QANewyearDragonHPic16,
// 	"qa_newyear_dragon_h_pic_17": model.QANewyearDragonHPic17,
// 	"qa_newyear_dragon_h_pic_18": model.QANewyearDragonHPic18,
// 	"qa_newyear_dragon_h_pic_19": model.QANewyearDragonHPic19,
// 	"qa_newyear_dragon_h_pic_20": model.QANewyearDragonHPic20,
// 	"qa_newyear_dragon_h_pic_21": model.QANewyearDragonHPic21,
// 	"qa_newyear_dragon_h_pic_22": model.QANewyearDragonHPic22,
// 	"qa_newyear_dragon_h_pic_23": model.QANewyearDragonHPic23,
// 	"qa_newyear_dragon_h_pic_24": model.QANewyearDragonHPic24,
// 	"qa_newyear_dragon_g_pic_01": model.QANewyearDragonGPic01,
// 	"qa_newyear_dragon_g_pic_02": model.QANewyearDragonGPic02,
// 	"qa_newyear_dragon_g_pic_03": model.QANewyearDragonGPic03,
// 	"qa_newyear_dragon_g_pic_04": model.QANewyearDragonGPic04,
// 	"qa_newyear_dragon_g_pic_05": model.QANewyearDragonGPic05,
// 	"qa_newyear_dragon_g_pic_06": model.QANewyearDragonGPic06,
// 	"qa_newyear_dragon_c_pic_01": model.QANewyearDragonCPic01,
// 	"qa_newyear_dragon_h_ani_01": model.QANewyearDragonHAni01,
// 	"qa_newyear_dragon_h_ani_02": model.QANewyearDragonHAni02,
// 	"qa_newyear_dragon_g_ani_01": model.QANewyearDragonGAni01,
// 	"qa_newyear_dragon_g_ani_02": model.QANewyearDragonGAni02,
// 	"qa_newyear_dragon_g_ani_03": model.QANewyearDragonGAni03,
// 	"qa_newyear_dragon_c_ani_01": model.QANewyearDragonCAni01,
// }); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_qa_picture_2)發生問題")
// }

// activity_game_whack_mole_picture資料表
// if _, err := a.Table(config.ACTIVITY_GAME_WHACK_MOLE_PICTURE_TABLE).Insert(command.Value{
// 	"activity_id": model.ActivityID,
// 	"game_id":     gameid,

// 	// 敲敲樂自定義
// 	"whackmole_classic_h_pic_01": model.WhackmoleClassicHPic01,
// 	"whackmole_classic_h_pic_02": model.WhackmoleClassicHPic02,
// 	"whackmole_classic_h_pic_03": model.WhackmoleClassicHPic03,
// 	"whackmole_classic_h_pic_04": model.WhackmoleClassicHPic04,
// 	"whackmole_classic_h_pic_05": model.WhackmoleClassicHPic05,
// 	"whackmole_classic_h_pic_06": model.WhackmoleClassicHPic06,
// 	"whackmole_classic_h_pic_07": model.WhackmoleClassicHPic07,
// 	"whackmole_classic_h_pic_08": model.WhackmoleClassicHPic08,
// 	"whackmole_classic_h_pic_09": model.WhackmoleClassicHPic09,
// 	"whackmole_classic_h_pic_10": model.WhackmoleClassicHPic10,
// 	"whackmole_classic_h_pic_11": model.WhackmoleClassicHPic11,
// 	"whackmole_classic_h_pic_12": model.WhackmoleClassicHPic12,
// 	"whackmole_classic_h_pic_13": model.WhackmoleClassicHPic13,
// 	"whackmole_classic_h_pic_14": model.WhackmoleClassicHPic14,
// 	"whackmole_classic_h_pic_15": model.WhackmoleClassicHPic15,
// 	"whackmole_classic_g_pic_01": model.WhackmoleClassicGPic01,
// 	"whackmole_classic_g_pic_02": model.WhackmoleClassicGPic02,
// 	"whackmole_classic_g_pic_03": model.WhackmoleClassicGPic03,
// 	"whackmole_classic_g_pic_04": model.WhackmoleClassicGPic04,
// 	"whackmole_classic_g_pic_05": model.WhackmoleClassicGPic05,
// 	"whackmole_classic_c_pic_01": model.WhackmoleClassicCPic01,
// 	"whackmole_classic_c_pic_02": model.WhackmoleClassicCPic02,
// 	"whackmole_classic_c_pic_03": model.WhackmoleClassicCPic03,
// 	"whackmole_classic_c_pic_04": model.WhackmoleClassicCPic04,
// 	"whackmole_classic_c_pic_05": model.WhackmoleClassicCPic05,
// 	"whackmole_classic_c_pic_06": model.WhackmoleClassicCPic06,
// 	"whackmole_classic_c_pic_07": model.WhackmoleClassicCPic07,
// 	"whackmole_classic_c_pic_08": model.WhackmoleClassicCPic08,
// 	"whackmole_classic_c_ani_01": model.WhackmoleClassicCAni01,

// 	"whackmole_halloween_h_pic_01": model.WhackmoleHalloweenHPic01,
// 	"whackmole_halloween_h_pic_02": model.WhackmoleHalloweenHPic02,
// 	"whackmole_halloween_h_pic_03": model.WhackmoleHalloweenHPic03,
// 	"whackmole_halloween_h_pic_04": model.WhackmoleHalloweenHPic04,
// 	"whackmole_halloween_h_pic_05": model.WhackmoleHalloweenHPic05,
// 	"whackmole_halloween_h_pic_06": model.WhackmoleHalloweenHPic06,
// 	"whackmole_halloween_h_pic_07": model.WhackmoleHalloweenHPic07,
// 	"whackmole_halloween_h_pic_08": model.WhackmoleHalloweenHPic08,
// 	"whackmole_halloween_h_pic_09": model.WhackmoleHalloweenHPic09,
// 	"whackmole_halloween_h_pic_10": model.WhackmoleHalloweenHPic10,
// 	"whackmole_halloween_h_pic_11": model.WhackmoleHalloweenHPic11,
// 	"whackmole_halloween_h_pic_12": model.WhackmoleHalloweenHPic12,
// 	"whackmole_halloween_h_pic_13": model.WhackmoleHalloweenHPic13,
// 	"whackmole_halloween_h_pic_14": model.WhackmoleHalloweenHPic14,
// 	"whackmole_halloween_h_pic_15": model.WhackmoleHalloweenHPic15,
// 	"whackmole_halloween_g_pic_01": model.WhackmoleHalloweenGPic01,
// 	"whackmole_halloween_g_pic_02": model.WhackmoleHalloweenGPic02,
// 	"whackmole_halloween_g_pic_03": model.WhackmoleHalloweenGPic03,
// 	"whackmole_halloween_g_pic_04": model.WhackmoleHalloweenGPic04,
// 	"whackmole_halloween_g_pic_05": model.WhackmoleHalloweenGPic05,
// 	"whackmole_halloween_c_pic_01": model.WhackmoleHalloweenCPic01,
// 	"whackmole_halloween_c_pic_02": model.WhackmoleHalloweenCPic02,
// 	"whackmole_halloween_c_pic_03": model.WhackmoleHalloweenCPic03,
// 	"whackmole_halloween_c_pic_04": model.WhackmoleHalloweenCPic04,
// 	"whackmole_halloween_c_pic_05": model.WhackmoleHalloweenCPic05,
// 	"whackmole_halloween_c_pic_06": model.WhackmoleHalloweenCPic06,
// 	"whackmole_halloween_c_pic_07": model.WhackmoleHalloweenCPic07,
// 	"whackmole_halloween_c_pic_08": model.WhackmoleHalloweenCPic08,
// 	"whackmole_halloween_c_ani_01": model.WhackmoleHalloweenCAni01,

// 	"whackmole_christmas_h_pic_01": model.WhackmoleChristmasHPic01,
// 	"whackmole_christmas_h_pic_02": model.WhackmoleChristmasHPic02,
// 	"whackmole_christmas_h_pic_03": model.WhackmoleChristmasHPic03,
// 	"whackmole_christmas_h_pic_04": model.WhackmoleChristmasHPic04,
// 	"whackmole_christmas_h_pic_05": model.WhackmoleChristmasHPic05,
// 	"whackmole_christmas_h_pic_06": model.WhackmoleChristmasHPic06,
// 	"whackmole_christmas_h_pic_07": model.WhackmoleChristmasHPic07,
// 	"whackmole_christmas_h_pic_08": model.WhackmoleChristmasHPic08,
// 	"whackmole_christmas_h_pic_09": model.WhackmoleChristmasHPic09,
// 	"whackmole_christmas_h_pic_10": model.WhackmoleChristmasHPic10,
// 	"whackmole_christmas_h_pic_11": model.WhackmoleChristmasHPic11,
// 	"whackmole_christmas_h_pic_12": model.WhackmoleChristmasHPic12,
// 	"whackmole_christmas_h_pic_13": model.WhackmoleChristmasHPic13,
// 	"whackmole_christmas_h_pic_14": model.WhackmoleChristmasHPic14,
// 	"whackmole_christmas_g_pic_01": model.WhackmoleChristmasGPic01,
// 	"whackmole_christmas_g_pic_02": model.WhackmoleChristmasGPic02,
// 	"whackmole_christmas_g_pic_03": model.WhackmoleChristmasGPic03,
// 	"whackmole_christmas_g_pic_04": model.WhackmoleChristmasGPic04,
// 	"whackmole_christmas_g_pic_05": model.WhackmoleChristmasGPic05,
// 	"whackmole_christmas_g_pic_06": model.WhackmoleChristmasGPic06,
// 	"whackmole_christmas_g_pic_07": model.WhackmoleChristmasGPic07,
// 	"whackmole_christmas_g_pic_08": model.WhackmoleChristmasGPic08,
// 	"whackmole_christmas_c_pic_01": model.WhackmoleChristmasCPic01,
// 	"whackmole_christmas_c_pic_02": model.WhackmoleChristmasCPic02,
// 	"whackmole_christmas_c_pic_03": model.WhackmoleChristmasCPic03,
// 	"whackmole_christmas_c_pic_04": model.WhackmoleChristmasCPic04,
// 	"whackmole_christmas_c_pic_05": model.WhackmoleChristmasCPic05,
// 	"whackmole_christmas_c_pic_06": model.WhackmoleChristmasCPic06,
// 	"whackmole_christmas_c_pic_07": model.WhackmoleChristmasCPic07,
// 	"whackmole_christmas_c_pic_08": model.WhackmoleChristmasCPic08,
// 	"whackmole_christmas_c_ani_01": model.WhackmoleChristmasCAni01,
// 	"whackmole_christmas_c_ani_02": model.WhackmoleChristmasCAni02,

// 	// 音樂
// 	"whackmole_bgm_start":  model.WhackmoleBgmStart,
// 	"whackmole_bgm_gaming": model.WhackmoleBgmGaming,
// 	"whackmole_bgm_end":    model.WhackmoleBgmEnd,
// }); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_whack_mole_picture)發生問題")
// }

// activity_game_bingo_picture資料表
// if _, err := a.Table(config.ACTIVITY_GAME_BINGO_PICTURE_TABLE).Insert(command.Value{
// 	"activity_id": model.ActivityID,
// 	"game_id":     gameid,

// 	// 賓果遊戲自定義
// 	"bingo_classic_h_pic_01": model.BingoClassicHPic01,
// 	"bingo_classic_h_pic_02": model.BingoClassicHPic02,
// 	"bingo_classic_h_pic_03": model.BingoClassicHPic03,
// 	"bingo_classic_h_pic_04": model.BingoClassicHPic04,
// 	"bingo_classic_h_pic_05": model.BingoClassicHPic05,
// 	"bingo_classic_h_pic_06": model.BingoClassicHPic06,
// 	"bingo_classic_h_pic_07": model.BingoClassicHPic07,
// 	"bingo_classic_h_pic_08": model.BingoClassicHPic08,
// 	"bingo_classic_h_pic_09": model.BingoClassicHPic09,
// 	"bingo_classic_h_pic_10": model.BingoClassicHPic10,
// 	"bingo_classic_h_pic_11": model.BingoClassicHPic11,
// 	"bingo_classic_h_pic_12": model.BingoClassicHPic12,
// 	"bingo_classic_h_pic_13": model.BingoClassicHPic13,
// 	"bingo_classic_h_pic_14": model.BingoClassicHPic14,
// 	"bingo_classic_h_pic_15": model.BingoClassicHPic15,
// 	"bingo_classic_h_pic_16": model.BingoClassicHPic16,
// 	"bingo_classic_g_pic_01": model.BingoClassicGPic01,
// 	"bingo_classic_g_pic_02": model.BingoClassicGPic02,
// 	"bingo_classic_g_pic_03": model.BingoClassicGPic03,
// 	"bingo_classic_g_pic_04": model.BingoClassicGPic04,
// 	"bingo_classic_g_pic_05": model.BingoClassicGPic05,
// 	"bingo_classic_g_pic_06": model.BingoClassicGPic06,
// 	"bingo_classic_g_pic_07": model.BingoClassicGPic07,
// 	"bingo_classic_g_pic_08": model.BingoClassicGPic08,
// 	"bingo_classic_c_pic_01": model.BingoClassicCPic01,
// 	"bingo_classic_c_pic_02": model.BingoClassicCPic02,
// 	"bingo_classic_c_pic_03": model.BingoClassicCPic03,
// 	"bingo_classic_c_pic_04": model.BingoClassicCPic04,
// 	// "bingo_classic_c_pic_05": model.BingoClassicCPic05,
// 	"bingo_classic_h_ani_01": model.BingoClassicHAni01,
// 	"bingo_classic_g_ani_01": model.BingoClassicGAni01,
// 	"bingo_classic_c_ani_01": model.BingoClassicCAni01,
// 	"bingo_classic_c_ani_02": model.BingoClassicCAni02,

// 	"bingo_newyear_dragon_h_pic_01": model.BingoNewyearDragonHPic01,
// 	"bingo_newyear_dragon_h_pic_02": model.BingoNewyearDragonHPic02,
// 	"bingo_newyear_dragon_h_pic_03": model.BingoNewyearDragonHPic03,
// 	"bingo_newyear_dragon_h_pic_04": model.BingoNewyearDragonHPic04,
// 	"bingo_newyear_dragon_h_pic_05": model.BingoNewyearDragonHPic05,
// 	"bingo_newyear_dragon_h_pic_06": model.BingoNewyearDragonHPic06,
// 	"bingo_newyear_dragon_h_pic_07": model.BingoNewyearDragonHPic07,
// 	"bingo_newyear_dragon_h_pic_08": model.BingoNewyearDragonHPic08,
// 	"bingo_newyear_dragon_h_pic_09": model.BingoNewyearDragonHPic09,
// 	"bingo_newyear_dragon_h_pic_10": model.BingoNewyearDragonHPic10,
// 	"bingo_newyear_dragon_h_pic_11": model.BingoNewyearDragonHPic11,
// 	"bingo_newyear_dragon_h_pic_12": model.BingoNewyearDragonHPic12,
// 	"bingo_newyear_dragon_h_pic_13": model.BingoNewyearDragonHPic13,
// 	"bingo_newyear_dragon_h_pic_14": model.BingoNewyearDragonHPic14,
// 	// "bingo_newyear_dragon_h_pic_15": model.BingoNewyearDragonHPic15,
// 	"bingo_newyear_dragon_h_pic_16": model.BingoNewyearDragonHPic16,
// 	"bingo_newyear_dragon_h_pic_17": model.BingoNewyearDragonHPic17,
// 	"bingo_newyear_dragon_h_pic_18": model.BingoNewyearDragonHPic18,
// 	"bingo_newyear_dragon_h_pic_19": model.BingoNewyearDragonHPic19,
// 	"bingo_newyear_dragon_h_pic_20": model.BingoNewyearDragonHPic20,
// 	"bingo_newyear_dragon_h_pic_21": model.BingoNewyearDragonHPic21,
// 	"bingo_newyear_dragon_h_pic_22": model.BingoNewyearDragonHPic22,
// 	"bingo_newyear_dragon_g_pic_01": model.BingoNewyearDragonGPic01,
// 	"bingo_newyear_dragon_g_pic_02": model.BingoNewyearDragonGPic02,
// 	"bingo_newyear_dragon_g_pic_03": model.BingoNewyearDragonGPic03,
// 	"bingo_newyear_dragon_g_pic_04": model.BingoNewyearDragonGPic04,
// 	"bingo_newyear_dragon_g_pic_05": model.BingoNewyearDragonGPic05,
// 	"bingo_newyear_dragon_g_pic_06": model.BingoNewyearDragonGPic06,
// 	"bingo_newyear_dragon_g_pic_07": model.BingoNewyearDragonGPic07,
// 	"bingo_newyear_dragon_g_pic_08": model.BingoNewyearDragonGPic08,
// 	"bingo_newyear_dragon_c_pic_01": model.BingoNewyearDragonCPic01,
// 	"bingo_newyear_dragon_c_pic_02": model.BingoNewyearDragonCPic02,
// 	"bingo_newyear_dragon_c_pic_03": model.BingoNewyearDragonCPic03,
// 	"bingo_newyear_dragon_h_ani_01": model.BingoNewyearDragonHAni01,
// 	"bingo_newyear_dragon_g_ani_01": model.BingoNewyearDragonGAni01,
// 	"bingo_newyear_dragon_c_ani_01": model.BingoNewyearDragonCAni01,
// 	"bingo_newyear_dragon_c_ani_02": model.BingoNewyearDragonCAni02,
// 	"bingo_newyear_dragon_c_ani_03": model.BingoNewyearDragonCAni03,

// 	"bingo_cherry_h_pic_01": model.BingoCherryHPic01,
// 	"bingo_cherry_h_pic_02": model.BingoCherryHPic02,
// 	"bingo_cherry_h_pic_03": model.BingoCherryHPic03,
// 	"bingo_cherry_h_pic_04": model.BingoCherryHPic04,
// 	"bingo_cherry_h_pic_05": model.BingoCherryHPic05,
// 	"bingo_cherry_h_pic_06": model.BingoCherryHPic06,
// 	"bingo_cherry_h_pic_07": model.BingoCherryHPic07,
// 	"bingo_cherry_h_pic_08": model.BingoCherryHPic08,
// 	"bingo_cherry_h_pic_09": model.BingoCherryHPic09,
// 	"bingo_cherry_h_pic_10": model.BingoCherryHPic10,
// 	"bingo_cherry_h_pic_11": model.BingoCherryHPic11,
// 	"bingo_cherry_h_pic_12": model.BingoCherryHPic12,
// 	// "bingo_cherry_h_pic_13": model.BingoCherryHPic13,
// 	"bingo_cherry_h_pic_14": model.BingoCherryHPic14,
// 	"bingo_cherry_h_pic_15": model.BingoCherryHPic15,
// 	// "bingo_cherry_h_pic_16": model.BingoCherryHPic16,
// 	"bingo_cherry_h_pic_17": model.BingoCherryHPic17,
// 	"bingo_cherry_h_pic_18": model.BingoCherryHPic18,
// 	"bingo_cherry_h_pic_19": model.BingoCherryHPic19,
// 	"bingo_cherry_g_pic_01": model.BingoCherryGPic01,
// 	"bingo_cherry_g_pic_02": model.BingoCherryGPic02,
// 	"bingo_cherry_g_pic_03": model.BingoCherryGPic03,
// 	"bingo_cherry_g_pic_04": model.BingoCherryGPic04,
// 	"bingo_cherry_g_pic_05": model.BingoCherryGPic05,
// 	"bingo_cherry_g_pic_06": model.BingoCherryGPic06,
// 	"bingo_cherry_c_pic_01": model.BingoCherryCPic01,
// 	"bingo_cherry_c_pic_02": model.BingoCherryCPic02,
// 	"bingo_cherry_c_pic_03": model.BingoCherryCPic03,
// 	"bingo_cherry_c_pic_04": model.BingoCherryCPic04,
// 	// "bingo_cherry_h_ani_01": model.BingoCherryHAni01,
// 	"bingo_cherry_h_ani_02": model.BingoCherryHAni02,
// 	"bingo_cherry_h_ani_03": model.BingoCherryHAni03,
// 	"bingo_cherry_g_ani_01": model.BingoCherryGAni01,
// 	"bingo_cherry_g_ani_02": model.BingoCherryGAni02,
// 	"bingo_cherry_c_ani_01": model.BingoCherryCAni01,
// 	"bingo_cherry_c_ani_02": model.BingoCherryCAni02,

// 	// 音樂
// 	"bingo_bgm_start":  model.BingoBgmStart,
// 	"bingo_bgm_gaming": model.BingoBgmGaming,
// 	"bingo_bgm_end":    model.BingoBgmEnd,
// }); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_bingo_picture)發生問題")
// }

// activity_game_redpack_picture資料表
// if _, err := a.Table(config.ACTIVITY_GAME_REDPACK_PICTURE_TABLE).Insert(command.Value{
// 	"activity_id": model.ActivityID,
// 	"game_id":     gameid,

// 	// 搖紅包自定義
// 	"redpack_classic_h_pic_01": model.RedpackClassicHPic01,
// 	"redpack_classic_h_pic_02": model.RedpackClassicHPic02,
// 	"redpack_classic_h_pic_03": model.RedpackClassicHPic03,
// 	"redpack_classic_h_pic_04": model.RedpackClassicHPic04,
// 	"redpack_classic_h_pic_05": model.RedpackClassicHPic05,
// 	"redpack_classic_h_pic_06": model.RedpackClassicHPic06,
// 	"redpack_classic_h_pic_07": model.RedpackClassicHPic07,
// 	"redpack_classic_h_pic_08": model.RedpackClassicHPic08,
// 	"redpack_classic_h_pic_09": model.RedpackClassicHPic09,
// 	"redpack_classic_h_pic_10": model.RedpackClassicHPic10,
// 	"redpack_classic_h_pic_11": model.RedpackClassicHPic11,
// 	"redpack_classic_h_pic_12": model.RedpackClassicHPic12,
// 	"redpack_classic_h_pic_13": model.RedpackClassicHPic13,
// 	"redpack_classic_g_pic_01": model.RedpackClassicGPic01,
// 	"redpack_classic_g_pic_02": model.RedpackClassicGPic02,
// 	"redpack_classic_g_pic_03": model.RedpackClassicGPic03,
// 	"redpack_classic_h_ani_01": model.RedpackClassicHAni01,
// 	"redpack_classic_h_ani_02": model.RedpackClassicHAni02,
// 	"redpack_classic_g_ani_01": model.RedpackClassicGAni01,
// 	"redpack_classic_g_ani_02": model.RedpackClassicGAni02,
// 	"redpack_classic_g_ani_03": model.RedpackClassicGAni03,

// 	"redpack_cherry_h_pic_01": model.RedpackCherryHPic01,
// 	"redpack_cherry_h_pic_02": model.RedpackCherryHPic02,
// 	"redpack_cherry_h_pic_03": model.RedpackCherryHPic03,
// 	"redpack_cherry_h_pic_04": model.RedpackCherryHPic04,
// 	"redpack_cherry_h_pic_05": model.RedpackCherryHPic05,
// 	"redpack_cherry_h_pic_06": model.RedpackCherryHPic06,
// 	"redpack_cherry_h_pic_07": model.RedpackCherryHPic07,
// 	"redpack_cherry_g_pic_01": model.RedpackCherryGPic01,
// 	"redpack_cherry_g_pic_02": model.RedpackCherryGPic02,
// 	"redpack_cherry_h_ani_01": model.RedpackCherryHAni01,
// 	"redpack_cherry_h_ani_02": model.RedpackCherryHAni02,
// 	"redpack_cherry_g_ani_01": model.RedpackCherryGAni01,
// 	"redpack_cherry_g_ani_02": model.RedpackCherryGAni02,

// 	"redpack_christmas_h_pic_01": model.RedpackChristmasHPic01,
// 	"redpack_christmas_h_pic_02": model.RedpackChristmasHPic02,
// 	"redpack_christmas_h_pic_03": model.RedpackChristmasHPic03,
// 	"redpack_christmas_h_pic_04": model.RedpackChristmasHPic04,
// 	"redpack_christmas_h_pic_05": model.RedpackChristmasHPic05,
// 	"redpack_christmas_h_pic_06": model.RedpackChristmasHPic06,
// 	"redpack_christmas_h_pic_07": model.RedpackChristmasHPic07,
// 	"redpack_christmas_h_pic_08": model.RedpackChristmasHPic08,
// 	"redpack_christmas_h_pic_09": model.RedpackChristmasHPic09,
// 	"redpack_christmas_h_pic_10": model.RedpackChristmasHPic10,
// 	"redpack_christmas_h_pic_11": model.RedpackChristmasHPic11,
// 	"redpack_christmas_h_pic_12": model.RedpackChristmasHPic12,
// 	"redpack_christmas_h_pic_13": model.RedpackChristmasHPic13,
// 	"redpack_christmas_g_pic_01": model.RedpackChristmasGPic01,
// 	"redpack_christmas_g_pic_02": model.RedpackChristmasGPic02,
// 	"redpack_christmas_g_pic_03": model.RedpackChristmasGPic03,
// 	"redpack_christmas_g_pic_04": model.RedpackChristmasGPic04,
// 	"redpack_christmas_c_pic_01": model.RedpackChristmasCPic01,
// 	"redpack_christmas_c_pic_02": model.RedpackChristmasCPic02,
// 	"redpack_christmas_c_ani_01": model.RedpackChristmasCAni01,

// 	// 音樂
// 	"redpack_bgm_start":  model.RedpackBgmStart,
// 	"redpack_bgm_gaming": model.RedpackBgmGaming,
// 	"redpack_bgm_end":    model.RedpackBgmEnd,
// }); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_redpack_picture)發生問題")
// }

// activity_game_lottery_picture資料表
// if _, err := a.Table(config.ACTIVITY_GAME_LOTTERY_PICTURE_TABLE).Insert(command.Value{
// 	"activity_id": model.ActivityID,
// 	"game_id":     gameid,

// 	// 遊戲抽獎自定義
// 	"lottery_jiugongge_classic_h_pic_01": model.LotteryJiugonggeClassicHPic01,
// 	"lottery_jiugongge_classic_h_pic_02": model.LotteryJiugonggeClassicHPic02,
// 	"lottery_jiugongge_classic_h_pic_03": model.LotteryJiugonggeClassicHPic03,
// 	"lottery_jiugongge_classic_h_pic_04": model.LotteryJiugonggeClassicHPic04,
// 	"lottery_jiugongge_classic_g_pic_01": model.LotteryJiugonggeClassicGPic01,
// 	"lottery_jiugongge_classic_g_pic_02": model.LotteryJiugonggeClassicGPic02,
// 	"lottery_jiugongge_classic_c_pic_01": model.LotteryJiugonggeClassicCPic01,
// 	"lottery_jiugongge_classic_c_pic_02": model.LotteryJiugonggeClassicCPic02,
// 	"lottery_jiugongge_classic_c_pic_03": model.LotteryJiugonggeClassicCPic03,
// 	"lottery_jiugongge_classic_c_pic_04": model.LotteryJiugonggeClassicCPic04,
// 	"lottery_jiugongge_classic_c_ani_01": model.LotteryJiugonggeClassicCAni01,
// 	"lottery_jiugongge_classic_c_ani_02": model.LotteryJiugonggeClassicCAni02,
// 	"lottery_jiugongge_classic_c_ani_03": model.LotteryJiugonggeClassicCAni03,

// 	"lottery_turntable_classic_h_pic_01": model.LotteryTurntableClassicHPic01,
// 	"lottery_turntable_classic_h_pic_02": model.LotteryTurntableClassicHPic02,
// 	"lottery_turntable_classic_h_pic_03": model.LotteryTurntableClassicHPic03,
// 	"lottery_turntable_classic_h_pic_04": model.LotteryTurntableClassicHPic04,
// 	"lottery_turntable_classic_g_pic_01": model.LotteryTurntableClassicGPic01,
// 	"lottery_turntable_classic_g_pic_02": model.LotteryTurntableClassicGPic02,
// 	"lottery_turntable_classic_c_pic_01": model.LotteryTurntableClassicCPic01,
// 	"lottery_turntable_classic_c_pic_02": model.LotteryTurntableClassicCPic02,
// 	"lottery_turntable_classic_c_pic_03": model.LotteryTurntableClassicCPic03,
// 	"lottery_turntable_classic_c_pic_04": model.LotteryTurntableClassicCPic04,
// 	"lottery_turntable_classic_c_pic_05": model.LotteryTurntableClassicCPic05,
// 	"lottery_turntable_classic_c_pic_06": model.LotteryTurntableClassicCPic06,
// 	"lottery_turntable_classic_c_ani_01": model.LotteryTurntableClassicCAni01,
// 	"lottery_turntable_classic_c_ani_02": model.LotteryTurntableClassicCAni02,
// 	"lottery_turntable_classic_c_ani_03": model.LotteryTurntableClassicCAni03,

// 	"lottery_jiugongge_starrysky_h_pic_01": model.LotteryJiugonggeStarryskyHPic01,
// 	"lottery_jiugongge_starrysky_h_pic_02": model.LotteryJiugonggeStarryskyHPic02,
// 	"lottery_jiugongge_starrysky_h_pic_03": model.LotteryJiugonggeStarryskyHPic03,
// 	"lottery_jiugongge_starrysky_h_pic_04": model.LotteryJiugonggeStarryskyHPic04,
// 	"lottery_jiugongge_starrysky_h_pic_05": model.LotteryJiugonggeStarryskyHPic05,
// 	"lottery_jiugongge_starrysky_h_pic_06": model.LotteryJiugonggeStarryskyHPic06,
// 	"lottery_jiugongge_starrysky_h_pic_07": model.LotteryJiugonggeStarryskyHPic07,
// 	"lottery_jiugongge_starrysky_g_pic_01": model.LotteryJiugonggeStarryskyGPic01,
// 	"lottery_jiugongge_starrysky_g_pic_02": model.LotteryJiugonggeStarryskyGPic02,
// 	"lottery_jiugongge_starrysky_g_pic_03": model.LotteryJiugonggeStarryskyGPic03,
// 	"lottery_jiugongge_starrysky_g_pic_04": model.LotteryJiugonggeStarryskyGPic04,
// 	"lottery_jiugongge_starrysky_c_pic_01": model.LotteryJiugonggeStarryskyCPic01,
// 	"lottery_jiugongge_starrysky_c_pic_02": model.LotteryJiugonggeStarryskyCPic02,
// 	"lottery_jiugongge_starrysky_c_pic_03": model.LotteryJiugonggeStarryskyCPic03,
// 	"lottery_jiugongge_starrysky_c_pic_04": model.LotteryJiugonggeStarryskyCPic04,
// 	"lottery_jiugongge_starrysky_c_ani_01": model.LotteryJiugonggeStarryskyCAni01,
// 	"lottery_jiugongge_starrysky_c_ani_02": model.LotteryJiugonggeStarryskyCAni02,
// 	"lottery_jiugongge_starrysky_c_ani_03": model.LotteryJiugonggeStarryskyCAni03,
// 	"lottery_jiugongge_starrysky_c_ani_04": model.LotteryJiugonggeStarryskyCAni04,
// 	"lottery_jiugongge_starrysky_c_ani_05": model.LotteryJiugonggeStarryskyCAni05,
// 	"lottery_jiugongge_starrysky_c_ani_06": model.LotteryJiugonggeStarryskyCAni06,

// 	"lottery_turntable_starrysky_h_pic_01": model.LotteryTurntableStarryskyHPic01,
// 	"lottery_turntable_starrysky_h_pic_02": model.LotteryTurntableStarryskyHPic02,
// 	"lottery_turntable_starrysky_h_pic_03": model.LotteryTurntableStarryskyHPic03,
// 	"lottery_turntable_starrysky_h_pic_04": model.LotteryTurntableStarryskyHPic04,
// 	"lottery_turntable_starrysky_h_pic_05": model.LotteryTurntableStarryskyHPic05,
// 	"lottery_turntable_starrysky_h_pic_06": model.LotteryTurntableStarryskyHPic06,
// 	"lottery_turntable_starrysky_h_pic_07": model.LotteryTurntableStarryskyHPic07,
// 	"lottery_turntable_starrysky_h_pic_08": model.LotteryTurntableStarryskyHPic08,
// 	"lottery_turntable_starrysky_g_pic_01": model.LotteryTurntableStarryskyGPic01,
// 	"lottery_turntable_starrysky_g_pic_02": model.LotteryTurntableStarryskyGPic02,
// 	"lottery_turntable_starrysky_g_pic_03": model.LotteryTurntableStarryskyGPic03,
// 	"lottery_turntable_starrysky_g_pic_04": model.LotteryTurntableStarryskyGPic04,
// 	"lottery_turntable_starrysky_g_pic_05": model.LotteryTurntableStarryskyGPic05,
// 	"lottery_turntable_starrysky_c_pic_01": model.LotteryTurntableStarryskyCPic01,
// 	"lottery_turntable_starrysky_c_pic_02": model.LotteryTurntableStarryskyCPic02,
// 	"lottery_turntable_starrysky_c_pic_03": model.LotteryTurntableStarryskyCPic03,
// 	"lottery_turntable_starrysky_c_pic_04": model.LotteryTurntableStarryskyCPic04,
// 	"lottery_turntable_starrysky_c_pic_05": model.LotteryTurntableStarryskyCPic05,
// 	"lottery_turntable_starrysky_c_ani_01": model.LotteryTurntableStarryskyCAni01,
// 	"lottery_turntable_starrysky_c_ani_02": model.LotteryTurntableStarryskyCAni02,
// 	"lottery_turntable_starrysky_c_ani_03": model.LotteryTurntableStarryskyCAni03,
// 	"lottery_turntable_starrysky_c_ani_04": model.LotteryTurntableStarryskyCAni04,
// 	"lottery_turntable_starrysky_c_ani_05": model.LotteryTurntableStarryskyCAni05,
// 	"lottery_turntable_starrysky_c_ani_06": model.LotteryTurntableStarryskyCAni06,
// 	"lottery_turntable_starrysky_c_ani_07": model.LotteryTurntableStarryskyCAni07,

// 	// 音樂
// 	"lottery_bgm_gaming": model.LotteryBgmGaming,
// }); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_lottery_picture)發生問題")
// }

// activity_game_3d_gacha_machine_picture資料表
// if _, err := a.Table(config.ACTIVITY_GAME_GACHAMACHINE_PICTURE_TABLE).Insert(command.Value{
// 	"activity_id": model.ActivityID,
// 	"game_id":     gameid,

// 	// 扭蛋機自定義
// 	"3d_gacha_machine_classic_h_pic_02": model.GachaMachineClassicHPic02,
// 	"3d_gacha_machine_classic_h_pic_03": model.GachaMachineClassicHPic03,
// 	"3d_gacha_machine_classic_h_pic_04": model.GachaMachineClassicHPic04,
// 	"3d_gacha_machine_classic_h_pic_05": model.GachaMachineClassicHPic05,
// 	"3d_gacha_machine_classic_g_pic_01": model.GachaMachineClassicGPic01,
// 	"3d_gacha_machine_classic_g_pic_02": model.GachaMachineClassicGPic02,
// 	"3d_gacha_machine_classic_g_pic_03": model.GachaMachineClassicGPic03,
// 	"3d_gacha_machine_classic_g_pic_04": model.GachaMachineClassicGPic04,
// 	"3d_gacha_machine_classic_g_pic_05": model.GachaMachineClassicGPic05,
// 	"3d_gacha_machine_classic_c_pic_01": model.GachaMachineClassicCPic01,

// 	// 音樂
// 	"3d_gacha_machine_bgm_gaming": model.GachaMachineBgmGaming,
// }); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_3d_gacha_machine_picture)發生問題")
// }

// // activity_game_vote_picture資料表
// if _, err := a.Table(config.ACTIVITY_GAME_VOTE_PICTURE_TABLE).Insert(command.Value{
// 	"activity_id": model.ActivityID,
// 	"game_id":     gameid,

// 	// 投票自定義
// 	"vote_classic_h_pic_01": model.VoteClassicHPic01,
// 	"vote_classic_h_pic_02": model.VoteClassicHPic02,
// 	"vote_classic_h_pic_03": model.VoteClassicHPic03,
// 	"vote_classic_h_pic_04": model.VoteClassicHPic04,
// 	"vote_classic_h_pic_05": model.VoteClassicHPic05,
// 	"vote_classic_h_pic_06": model.VoteClassicHPic06,
// 	"vote_classic_h_pic_07": model.VoteClassicHPic07,
// 	"vote_classic_h_pic_08": model.VoteClassicHPic08,
// 	"vote_classic_h_pic_09": model.VoteClassicHPic09,
// 	"vote_classic_h_pic_10": model.VoteClassicHPic10,
// 	"vote_classic_h_pic_11": model.VoteClassicHPic11,
// 	// "vote_classic_h_pic_12": model.VoteClassicHPic12,
// 	"vote_classic_h_pic_13": model.VoteClassicHPic13,
// 	"vote_classic_h_pic_14": model.VoteClassicHPic14,
// 	"vote_classic_h_pic_15": model.VoteClassicHPic15,
// 	"vote_classic_h_pic_16": model.VoteClassicHPic16,
// 	"vote_classic_h_pic_17": model.VoteClassicHPic17,
// 	"vote_classic_h_pic_18": model.VoteClassicHPic18,
// 	"vote_classic_h_pic_19": model.VoteClassicHPic19,
// 	"vote_classic_h_pic_20": model.VoteClassicHPic20,
// 	"vote_classic_h_pic_21": model.VoteClassicHPic21,
// 	// "vote_classic_h_pic_22": model.VoteClassicHPic22,
// 	"vote_classic_h_pic_23": model.VoteClassicHPic23,
// 	"vote_classic_h_pic_24": model.VoteClassicHPic24,
// 	"vote_classic_h_pic_25": model.VoteClassicHPic25,
// 	"vote_classic_h_pic_26": model.VoteClassicHPic26,
// 	"vote_classic_h_pic_27": model.VoteClassicHPic27,
// 	"vote_classic_h_pic_28": model.VoteClassicHPic28,
// 	"vote_classic_h_pic_29": model.VoteClassicHPic29,
// 	"vote_classic_h_pic_30": model.VoteClassicHPic30,
// 	"vote_classic_h_pic_31": model.VoteClassicHPic31,
// 	"vote_classic_h_pic_32": model.VoteClassicHPic32,
// 	"vote_classic_h_pic_33": model.VoteClassicHPic33,
// 	"vote_classic_h_pic_34": model.VoteClassicHPic34,
// 	"vote_classic_h_pic_35": model.VoteClassicHPic35,
// 	"vote_classic_h_pic_36": model.VoteClassicHPic36,
// 	"vote_classic_h_pic_37": model.VoteClassicHPic37,
// 	"vote_classic_g_pic_01": model.VoteClassicGPic01,
// 	"vote_classic_g_pic_02": model.VoteClassicGPic02,
// 	"vote_classic_g_pic_03": model.VoteClassicGPic03,
// 	"vote_classic_g_pic_04": model.VoteClassicGPic04,
// 	"vote_classic_g_pic_05": model.VoteClassicGPic05,
// 	"vote_classic_g_pic_06": model.VoteClassicGPic06,
// 	"vote_classic_g_pic_07": model.VoteClassicGPic07,
// 	"vote_classic_c_pic_01": model.VoteClassicCPic01,
// 	"vote_classic_c_pic_02": model.VoteClassicCPic02,
// 	"vote_classic_c_pic_03": model.VoteClassicCPic03,
// 	"vote_classic_c_pic_04": model.VoteClassicCPic04,

// 	// 音樂
// 	"vote_bgm_gaming": model.VoteBgmGaming,
// }); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_3d_gacha_machine_picture)發生問題")
// }

// 建立拔河遊戲時，預設兩個獎品(勝方敗方)
// if game == "tugofwar" {
// 	// var amount string
// 	// if model.Prize == "all" {
// 	// 	// 全部發獎(獎品數預設為遊戲人數)
// 	// 	amount = model.People
// 	// } else if model.Prize == "uniform" {
// 	// 	// 統一發獎(獎品數預設1)
// 	// 	amount = "1"
// 	// }

// 	if err := DefaultPrizeModel().SetDbConn(a.DbConn).SetRedisConn(a.RedisConn).
// 		Add(true, "tugofwar", utils.UUID(20), NewPrizeModel{
// 			ActivityID:    model.ActivityID,
// 			GameID:        gameid,
// 			PrizeName:     "勝利隊伍",
// 			PrizeType:     "first",
// 			PrizePicture:  config.UPLOAD_SYSTEM_URL + "img-prize-pic.png",
// 			PrizeAmount:   "0",
// 			PrizePrice:    "0",
// 			PrizeMethod:   "site",
// 			PrizePassword: "win",
// 			TeamType:      "win",
// 		}); err != nil {
// 		return errors.New("錯誤: 新增勝方遊戲獎品發生問題")
// 	}

// 	if err := DefaultPrizeModel().SetDbConn(a.DbConn).SetRedisConn(a.RedisConn).
// 		Add(true, "tugofwar", utils.UUID(20), NewPrizeModel{
// 			ActivityID:    model.ActivityID,
// 			GameID:        gameid,
// 			PrizeName:     "落敗隊伍",
// 			PrizeType:     "thanks",
// 			PrizePicture:  config.UPLOAD_SYSTEM_URL + "img-prize-pic.png",
// 			PrizeAmount:   "0",
// 			PrizePrice:    "0",
// 			PrizeMethod:   "site",
// 			PrizePassword: "lose",
// 			TeamType:      "lose",
// 		}); err != nil {
// 		return errors.New("錯誤: 新增敗方遊戲獎品發生問題")
// 	}
// }

// 建立賓果遊戲時，預設一個獎品
// if game == "bingo" {
// 	if err := DefaultPrizeModel().SetDbConn(a.DbConn).SetRedisConn(a.RedisConn).
// 		Add(true, "bingo", utils.UUID(20), NewPrizeModel{
// 			ActivityID:    model.ActivityID,
// 			GameID:        gameid,
// 			PrizeName:     "賓果遊戲獎品",
// 			PrizeType:     "first",
// 			PrizePicture:  config.UPLOAD_SYSTEM_URL + "img-prize-pic.png",
// 			PrizeAmount:   "0",
// 			PrizePrice:    "0",
// 			PrizeMethod:   "site",
// 			PrizePassword: "bingo",
// 			TeamType:      "",
// 		}); err != nil {
// 		return errors.New("錯誤: 新增賓果遊戲獎品發生問題")
// 	}
// }

// if game == "QA" {
// 	if model.QA1 == "" || model.QA1Options == "" ||
// 		model.QA1Answer == "" || model.QA1Score == "" {
// 		return errors.New("錯誤: 題目設置最少一題，請重新設置")
// 	}

// 	// activity_game_qa資料表
// 	if _, err := a.Table(config.ACTIVITY_GAME_QA_TABLE).Insert(command.Value{
// 		"activity_id":  model.ActivityID,
// 		"game_id":      gameid,
// 		"qa_1":         model.QA1,
// 		"qa_1_options": model.QA1Options,
// 		"qa_1_answer":  model.QA1Answer,
// 		"qa_1_score":   utils.GetInt64(model.QA1Score, 0),

// 		"qa_2":         model.QA2,
// 		"qa_2_options": model.QA2Options,
// 		"qa_2_answer":  model.QA2Answer,
// 		"qa_2_score":   utils.GetInt64(model.QA2Score, 0),

// 		"qa_3":         model.QA3,
// 		"qa_3_options": model.QA3Options,
// 		"qa_3_answer":  model.QA3Answer,
// 		"qa_3_score":   utils.GetInt64(model.QA3Score, 0),

// 		"qa_4":         model.QA4,
// 		"qa_4_options": model.QA4Options,
// 		"qa_4_answer":  model.QA4Answer,
// 		"qa_4_score":   utils.GetInt64(model.QA4Score, 0),

// 		"qa_5":         model.QA5,
// 		"qa_5_options": model.QA5Options,
// 		"qa_5_answer":  model.QA5Answer,
// 		"qa_5_score":   utils.GetInt64(model.QA5Score, 0),

// 		"qa_6":         model.QA6,
// 		"qa_6_options": model.QA6Options,
// 		"qa_6_answer":  model.QA6Answer,
// 		"qa_6_score":   utils.GetInt64(model.QA6Score, 0),

// 		"qa_7":         model.QA7,
// 		"qa_7_options": model.QA7Options,
// 		"qa_7_answer":  model.QA7Answer,
// 		"qa_7_score":   utils.GetInt64(model.QA7Score, 0),

// 		"qa_8":         model.QA8,
// 		"qa_8_options": model.QA8Options,
// 		"qa_8_answer":  model.QA8Answer,
// 		"qa_8_score":   utils.GetInt64(model.QA8Score, 0),

// 		"qa_9":         model.QA9,
// 		"qa_9_options": model.QA9Options,
// 		"qa_9_answer":  model.QA9Answer,
// 		"qa_9_score":   utils.GetInt64(model.QA9Score, 0),

// 		"qa_10":         model.QA10,
// 		"qa_10_options": model.QA10Options,
// 		"qa_10_answer":  model.QA10Answer,
// 		"qa_10_score":   utils.GetInt64(model.QA10Score, 0),

// 		"qa_11":         model.QA11,
// 		"qa_11_options": model.QA11Options,
// 		"qa_11_answer":  model.QA11Answer,
// 		"qa_11_score":   utils.GetInt64(model.QA11Score, 0),

// 		"qa_12":         model.QA12,
// 		"qa_12_options": model.QA12Options,
// 		"qa_12_answer":  model.QA12Answer,
// 		"qa_12_score":   utils.GetInt64(model.QA12Score, 0),

// 		"qa_13":         model.QA13,
// 		"qa_13_options": model.QA13Options,
// 		"qa_13_answer":  model.QA13Answer,
// 		"qa_13_score":   utils.GetInt64(model.QA13Score, 0),

// 		"qa_14":         model.QA14,
// 		"qa_14_options": model.QA14Options,
// 		"qa_14_answer":  model.QA14Answer,
// 		"qa_14_score":   utils.GetInt64(model.QA14Score, 0),

// 		"qa_15":         model.QA15,
// 		"qa_15_options": model.QA15Options,
// 		"qa_15_answer":  model.QA15Answer,
// 		"qa_15_score":   utils.GetInt64(model.QA15Score, 0),

// 		"qa_16":         model.QA16,
// 		"qa_16_options": model.QA16Options,
// 		"qa_16_answer":  model.QA16Answer,
// 		"qa_16_score":   utils.GetInt64(model.QA16Score, 0),

// 		"qa_17":         model.QA17,
// 		"qa_17_options": model.QA17Options,
// 		"qa_17_answer":  model.QA17Answer,
// 		"qa_17_score":   utils.GetInt64(model.QA17Score, 0),

// 		"qa_18":         model.QA18,
// 		"qa_18_options": model.QA18Options,
// 		"qa_18_answer":  model.QA18Answer,
// 		"qa_18_score":   utils.GetInt64(model.QA18Score, 0),

// 		"qa_19":         model.QA19,
// 		"qa_19_options": model.QA19Options,
// 		"qa_19_answer":  model.QA19Answer,
// 		"qa_19_score":   utils.GetInt64(model.QA19Score, 0),

// 		"qa_20":         model.QA20,
// 		"qa_20_options": model.QA20Options,
// 		"qa_20_answer":  model.QA20Answer,
// 		"qa_20_score":   utils.GetInt64(model.QA20Score, 0),

// 		"total_qa":  model.TotalQA,
// 		"qa_second": model.QASecond,
// 		"qa_round":  1,
// 		"qa_people": 0,
// 	}); err != nil {
// 		return errors.New("錯誤: 新增遊戲快問快答資料發生問題")
// 	}
// }

// 用戶資訊
// userModel, err := DefaultUserModel().SetDbConn(a.DbConn).
// 	SetRedisConn(a.RedisConn).Find(true, true, "",
// 	"users.user_id", model.UserID)
// if err != nil {
// 	return err
// }
// if people, err := strconv.Atoi(model.People); err != nil || people > 1000 {
// 	return errors.New("錯誤: 遊戲上限人數為1000人，請輸入有效的遊戲人數")
// }

// if game == "whack_mole" {
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

// if model.VoteStartTime == "" {
// 	// 將 now 格式化為字串
// 	nowStr := now.Format("2006-01-02 15:04")
// 	model.VoteStartTime = nowStr
// }
// if model.VoteEndTime == "" {
// 	// 將 now 格式化為字串
// 	nowStr := now.Format("2006-01-02 15:04")
// 	model.VoteEndTime = nowStr
// }

// gachaMachineCustomizefields = []string{
// 	"activity_id",
// 	"game_id",

// 	// 扭蛋遊戲自定義
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
// 	"3d_gacha_machine_bgm_gaming", // 遊戲進行中
// }

// bingoCustomizefields = []string{
// 	"activity_id",
// 	"game_id",

// 	// 賓果遊戲自定義
// 	// 音樂
// 	"bingo_bgm_start",  // 遊戲開始
// 	"bingo_bgm_gaming", // 遊戲進行中
// 	"bingo_bgm_end",    // 遊戲結束

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
// }

// drawNumbersCustomizefields = []string{
// 	"activity_id",
// 	"game_id",

// 	// 搖號抽獎自定義
// 	// 音樂
// 	"draw_numbers_bgm_gaming",

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
// }

// lotteryCustomizefields = []string{
// 	"activity_id",
// 	"game_id",

// 	"lottery_bgm_gaming",

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
// }

// monopolyCustomizefields = []string{
// 	"activity_id",
// 	"game_id",

// 	"monopoly_bgm_start",
// 	"monopoly_bgm_gaming",
// 	"monopoly_bgm_end",

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
// }

// qaCustomizefields = []string{
// 	"activity_id",
// 	"game_id",

// 	"qa_bgm_start",  // 遊戲開始
// 	"qa_bgm_gaming", // 遊戲進行中
// 	"qa_bgm_end",    // 遊戲結束

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
// }

// qa2Customizefields = []string{
// 	"activity_id",
// 	"game_id",

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
// }

// redpackCustomizefields = []string{
// 	"activity_id",
// 	"game_id",

// 	"redpack_bgm_start",
// 	"redpack_bgm_gaming",
// 	"redpack_bgm_end",

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
// }

// ropepackCustomizefields = []string{
// 	"activity_id",
// 	"game_id",

// 	"ropepack_bgm_start",
// 	"ropepack_bgm_gaming",
// 	"ropepack_bgm_end",

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
// }

// tugofwarCustomizefields = []string{
// 	"activity_id",
// 	"game_id",

// 	"tugofwar_bgm_start",  // 遊戲開始
// 	"tugofwar_bgm_gaming", // 遊戲進行中
// 	"tugofwar_bgm_end",    // 遊戲結束

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
// }

// voteCustomizefields = []string{
// 	"activity_id",
// 	"game_id",

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
// 	"vote_bgm_gaming", // 遊戲進行中
// }

// whackmoleCustomizefields = []string{
// 	"activity_id",
// 	"game_id",

// 	"whackmole_bgm_start",
// 	"whackmole_bgm_gaming",
// 	"whackmole_bgm_end",

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
// 	"whackmole_christmas_c_ani_02",
// }

// if model.Second != "" {
// 	if _, err := strconv.Atoi(model.Second); err != nil {
// 		return errors.New("錯誤: 限時秒數資料發生問題，請輸入有效的秒數")
// 	}
// }

// if model.MaxPeople != "" && model.People != "" {
// 判斷遊戲人數上限
// maxPeopleInt, err1 := strconv.Atoi(model.MaxPeople)
// peopleInt, err2 := strconv.Atoi(model.People)
// if err1 != nil || err2 != nil || peopleInt > maxPeopleInt {
// 	return errors.New("錯誤: 遊戲人數上限資料發生問題，請輸入有效的遊戲人數上限")
// }
// }

// if model.MaxTimes != "" {
// 	if _, err := strconv.Atoi(model.MaxTimes); err != nil {
// 		return errors.New("錯誤: 遊戲上限次數發生問題，請輸入有效的遊戲次數")
// 	}
// }

// if model.Percent != "" {
// 	if percentInt, err := strconv.Atoi(model.Percent); err != nil ||
// 		percentInt > 100 || percentInt < 0 {
// 		return errors.New("錯誤: 中獎機率必須為0-100，請輸入有效的中獎機率")
// 	}
// }

// if model.FirstPrize != "" {
// 	if people, err := strconv.Atoi(model.FirstPrize); err != nil || people > 50 {
// 		return errors.New("錯誤: 頭獎中獎人數上限資料發生問題，請輸入有效的人數")
// 	}
// }

// if model.SecondPrize != "" {
// 	if people, err := strconv.Atoi(model.SecondPrize); err != nil || people > 50 {
// 		return errors.New("錯誤: 二獎中獎人數上限資料發生問題，請輸入有效的人數")
// 	}
// }

// if model.ThirdPrize != "" {
// 	if people, err := strconv.Atoi(model.ThirdPrize); err != nil || people > 100 {
// 		return errors.New("錯誤: 三獎中獎人數上限資料發生問題，請輸入有效的人數")
// 	}
// }

// if model.GeneralPrize != "" {
// 	if people, err := strconv.Atoi(model.GeneralPrize); err != nil || people > 800 {
// 		return errors.New("錯誤: 普通獎中獎人數上限資料發生問題，請輸入有效的人數")
// 	}
// }

// if line, err := strconv.Atoi(model.BingoLine); err != nil ||
// 	line < 1 || line > 10 {
// 	return errors.New("錯誤: 賓果連線數資料發生問題(最多10條線，最少1條線)，請輸入有效的連線數")
// }

// if number, err := strconv.Atoi(model.MaxNumber); err != nil ||
// 	number < 16 || number > 99 {
// 	return errors.New("錯誤: 最大號碼資料發生問題(號碼必須大於16且小於100)，請輸入有效的連線數")
// }

// if _, err := strconv.Atoi(model.RoundPrize); err != nil {
// 	return errors.New("錯誤: 每輪發獎人數資料發生問題，請輸入有效的資料")
// }

// if _, err := strconv.Atoi(model.VoteTimes); err != nil {
// 	return errors.New("錯誤: 人員投票次數資料發生問題，請輸入有效的資料")
// }

// if _, err := a.Table(config.ACTIVITY_GAME_REDPACK_PICTURE_TABLE).
// 	Insert(FilterFields(data, redpackCustomizefields)); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_redpack_picture)發生問題")
// }

// if _, err := a.Table(config.ACTIVITY_GAME_ROPEPACK_PICTURE_TABLE).
// 	Insert(FilterFields(data, ropepackCustomizefields)); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_ropepack_picture)發生問題")
// }

// if _, err := a.Table(config.ACTIVITY_GAME_WHACK_MOLE_PICTURE_TABLE).
// 	Insert(FilterFields(data, whackmoleCustomizefields)); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_whack_mole_picture)發生問題")
// }

// if _, err := a.Table(config.ACTIVITY_GAME_LOTTERY_PICTURE_TABLE).
// 	Insert(FilterFields(data, lotteryCustomizefields)); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_lottery_picture)發生問題")
// }

// if _, err := a.Table(config.ACTIVITY_GAME_MONOPOLY_PICTURE_TABLE).
// 	Insert(FilterFields(data, monopolyCustomizefields)); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_monopoly_picture)發生問題")
// }

// if _, err := a.Table(config.ACTIVITY_GAME_QA_PICTURE_TABLE_1).
// 	Insert(FilterFields(data, qaCustomizefields)); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_qa_picture)發生問題")
// }

// if _, err := a.Table(config.ACTIVITY_GAME_QA_PICTURE_TABLE_2).
// 	Insert(FilterFields(data, qa2Customizefields)); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_qa_picture_2)發生問題")
// }

// if _, err := a.Table(config.ACTIVITY_GAME_DRAW_NUMBERS_PICTURE_TABLE).
// 	Insert(FilterFields(data, drawNumbersCustomizefields)); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_draw_numbers_picture)發生問題")
// }

// if _, err := a.Table(config.ACTIVITY_GAME_TUGOFWAR_PICTURE_TABLE).
// 	Insert(FilterFields(data, tugofwarCustomizefields)); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_tugofwar_picture)發生問題")
// }

// if _, err := a.Table(config.ACTIVITY_GAME_BINGO_PICTURE_TABLE).
// 	Insert(FilterFields(data, bingoCustomizefields)); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_bingo_picture)發生問題")
// }

// if _, err := a.Table(config.ACTIVITY_GAME_GACHAMACHINE_PICTURE_TABLE).
// 	Insert(FilterFields(data, gachaMachineCustomizefields)); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_3d_gacha_machine_picture)發生問題")
// }

// if _, err := a.Table(config.ACTIVITY_GAME_VOTE_PICTURE_TABLE).
// 	Insert(FilterFields(data, voteCustomizefields)); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game_vote_picture)發生問題")
// }

// activity_game資料表(mysql)
// if _, err := a.Table(a.TableName).
// 	Insert(FilterFields(data, fields)); err != nil {
// 	return errors.New("錯誤: 新增遊戲場次(activity_game)發生問題")
// }

// if model.QA1 == "" || model.QA1Options == "" ||
// 	model.QA1Answer == "" || model.QA1Score == "" {
// 	return errors.New("錯誤: 題目設置最少一題，請重新設置")
// }

// // 題目設置資訊
// qas := []string{
// 	"activity_id",
// 	"game_id",
// 	"total_qa",
// 	"qa_second",
// }

// total, err := strconv.Atoi(model.TotalQA)
// if err != nil {
// 	return errors.New("錯誤: 題目資料數發生問題，請重新設置")
// }

// // 迴圈處理題目資料並將對應的參數加入 qas 陣列
// for i := 1; i <= total; i++ {
// 	qaFieldPrefix := fmt.Sprintf("qa_%d", i)
// 	qaOptionsFieldPrefix := fmt.Sprintf("qa_%d_options", i)
// 	qaAnswerFieldPrefix := fmt.Sprintf("qa_%d_answer", i)
// 	qaScoreFieldPrefix := fmt.Sprintf("qa_%d_score", i)

// 	// 把這些動態生成的欄位名稱加入 qas 陣列
// 	qas = append(qas, qaFieldPrefix, qaOptionsFieldPrefix, qaAnswerFieldPrefix, qaScoreFieldPrefix)
// }

// // log.Println("快問快答qa參數: ", qas)

// // activity_game_qa資料表(題目設置資訊)
// if _, err := a.Table(config.ACTIVITY_GAME_QA_TABLE).
// 	Insert(FilterFields(data, qas)); err != nil {
// 	return errors.New("錯誤: 新增遊戲快問快答資料發生問題")
// }

// 題目設置資訊
// qas := []string{
// 	"activity_id",
// 	"game_id",
// 	"total_qa",
// 	"qa_second",
// }

// total, err := strconv.Atoi(model.TotalQA)
// if err != nil {
// 	return errors.New("錯誤: 題目資料數發生問題，請重新設置")
// }

// 迴圈處理題目資料並將對應的參數加入 qas 陣列
// for i := 1; i <= total; i++ {
// 	qaFieldPrefix := fmt.Sprintf("qa_%d", i)
// 	qaOptionsFieldPrefix := fmt.Sprintf("qa_%d_options", i)
// 	qaAnswerFieldPrefix := fmt.Sprintf("qa_%d_answer", i)
// 	qaScoreFieldPrefix := fmt.Sprintf("qa_%d_score", i)

// 	// 把這些動態生成的欄位名稱加入 qas 陣列
// 	qas = append(qas, qaFieldPrefix, qaOptionsFieldPrefix, qaAnswerFieldPrefix, qaScoreFieldPrefix)
// }

// log.Println("快問快答qa參數: ", qas)

// activity_game_qa資料表(題目設置資訊)
// if _, err := a.Table(config.ACTIVITY_GAME_QA_TABLE).
// 	Insert(FilterFields(data, qas)); err != nil {
// 	return errors.New("錯誤: 新增遊戲快問快答資料發生問題")
// }

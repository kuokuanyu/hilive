package models

import (
	"hilive/modules/cache"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/mongo"
)

// GameModel 資料表欄位
type GameModel struct {
	Base         `json:"-"`
	ID           int64  `json:"id"`
	UserID       string `json:"user_id"`
	ActivityID   string `json:"activity_id" example:"activity_id"`
	GameID       string `json:"game_id" example:"game_id"`
	Game         string `json:"game" example:"game name"`
	Title        string `json:"title" example:"game title"`
	GameType     string `json:"game_type" example:"game type"`
	LimitTime    string `json:"limit_time" example:"open、close"`
	Second       int64  `json:"second" example:"30"`
	MaxPeople    int64  `json:"max_people" example:"100"`
	People       int64  `json:"people" example:"100"`
	MaxTimes     int64  `json:"max_times" example:"10"`
	Allow        string `json:"allow" example:"open、close"`
	Percent      int64  `json:"percent" example:"0-100"`
	FirstPrize   int64  `json:"first_prize" example:"50(上限為50人)"`
	SecondPrize  int64  `json:"second_prize" example:"50(上限為50人)"`
	ThirdPrize   int64  `json:"third_prize" example:"100(上限為100人)"`
	GeneralPrize int64  `json:"general_prize" example:"800(上限為800人)"`

	// 主題相關
	Topic              string             `json:"topic" example:"01_classic"`
	Skin               string             `json:"skin" example:"classic"`
	Music              string             `json:"music" example:"classic"`
	CustomizeSceneData map[string]interface{} `json:"customize_scene_data" bson:"customize_scene_data"` // 這裡會包含所有畫面設定資料

	GameRound         int64  `json:"game_round" example:"1"`
	GameSecond        int64  `json:"game_second" example:"30"`
	GameStatus        string `json:"game_status" example:"open、start、end、close、calculate"`
	ControlGameStatus string `json:"control_game_status" example:"open、start、end、close、calculate"`
	GameAttend        int64  `json:"game_attend" example:"0"`
	DisplayName       string `json:"display_name" example:"open、close"`   // 是否顯示中獎人員姓名頭像
	GameOrder         int64  `json:"game_order" example:"1"`              // 遊戲場次排序
	BoxReflection     string `json:"box_reflection" example:"open、close"` // 扭蛋機遊戲開關盒反射
	SamePeople        string `json:"same_people" example:"open、close"`    // 拔河遊戲人數是否一致

	// 拔河遊戲
	AllowChooseTeam     string `json:"allow_choose_team" example:"open、close"` // 允許玩家選擇隊伍
	LeftTeamName        string `json:"left_team_name" example:"name"`          // 左邊隊伍名稱
	LeftTeamPicture     string `json:"left_team_picture" example:"picture"`    // 左邊隊伍照片
	RightTeamName       string `json:"right_team_name" example:"name"`         // 右邊隊伍名稱
	RightTeamPicture    string `json:"right_team_picture" example:"picture"`   // 右邊隊伍照片
	LeftTeamGameAttend  int64  `json:"left_team_game_attend" example:"0"`      // 左方隊伍參加遊戲人數
	RightTeamGameAttend int64  `json:"right_team_game_attend" example:"0"`     // 右方隊伍參加遊戲人數
	Prize               string `json:"prize" example:"uniform、all"`            // 獎品發放

	// 賓果遊戲
	MaxNumber  int64 `json:"max_number" example:"0"`  // 最大號碼
	BingoLine  int64 `json:"bingo_line" example:"0"`  // 賓果連線數
	RoundPrize int64 `json:"round_prize" example:"0"` // 每輪發獎人數
	BingoRound int64 `json:"bingo_round" example:"0"` // 賓果遊戲進行回合數

	// 隊伍資訊
	LeftTeamPlayers  []string `json:"left_team_players" example:"A"`  // 左方隊伍玩家
	RightTeamPlayers []string `json:"right_team_players" example:"A"` // 右方隊伍玩家
	LeftTeamLeader   string   `json:"left_team_leader" example:"A"`   // 左方隊伍隊長
	RightTeamLeader  string   `json:"right_team_leader" example:"A"`  // 右方隊伍隊長
	WinTeam          string   `json:"win_team" example:"A"`           // 獲勝隊伍
	// LeftTeamMVP      UserModel `json:"left_team_mvp" example:"A"`      // 左方隊伍mvp
	// RightTeamMVP     UserModel `json:"right_team_mvp" example:"A"`     // 右方隊伍mvp

	// 扭蛋機遊戲
	GachaMachineReflection float64 `json:"gacha_machine_reflection" example:"open、close"` // 球的反射度
	IsShake                string  `json:"is_shake" example:"true、false"`                 // 是否搖晃
	ReflectiveSwitch       string  `json:"reflective_switch" example:"open、close"`        // 反射開關

	// 投票遊戲
	VoteScreen       string `json:"vote_screen" example:"bar_chart、rank、detail_information"` // 投票畫面(長條圖顯示、排名顯示、詳細資訊顯示)
	VoteTimes        int64  `json:"vote_times" example:"0"`                                  // 人員投票次數
	VoteMethod       string `json:"vote_method" example:"all_vote、single_group、all_group"`   // 投票模式(全選項投票)
	VoteMethodPlayer string `json:"vote_method_player" example:"one_vote、free_vote"`         // 玩家投票方式(一個選項一票、自由投票)
	VoteRestriction  string `json:"vote_restriction" example:"all_player、special_officer"`   // 投票限制(所有人員都能投票、特殊人員才能投票)
	AvatarShape      string `json:"avatar_shape" example:"circle、square"`                    // 選項框是圓形還是方形
	VoteStartTime    string `json:"vote_start_time" example:""`                              // 投票開始時間
	VoteEndTime      string `json:"vote_end_time" example:""`                                // 投票結束時間
	AutoPlay         string `json:"auto_play" example:"open、close"`                          // 自動輪播
	ShowRank         string `json:"show_rank" example:"open、close"`                          // 排名展示
	TitleSwitch      string `json:"title_switch" example:"open、close"`                       // 場次名稱
	ArrangementGuest string `json:"arrangement_guest" example:"list、side_by_side"`           // 玩家端選項排列方式

	// 編輯次數
	EditTimes int64 `json:"edit_times" example:"0"` // 編輯次數

	// 基本設置
	LotteryGameAllow      string `json:"lottery_game_allow" example:"open、close"`          // 遊戲抽獎遊戲是否允許重複中獎
	RedpackGameAllow      string `json:"redpack_game_allow" example:"open、close"`          // 搖紅包遊戲是否允許重複中獎
	RopepackGameAllow     string `json:"ropepack_game_allow" example:"open、close"`         // 套紅包遊戲是否允許重複中獎
	WhackMoleGameAllow    string `json:"whack_mole_game_allow" example:"open、close"`       // 敲敲樂遊戲是否允許重複中獎
	MonopolyGameAllow     string `json:"monopoly_game_allow" example:"open、close"`         // 鑑定師遊戲是否允許重複中獎
	QAGameAllow           string `json:"qa_game_allow" example:"open、close"`               // 快問快答遊戲是否允許重複中獎
	DrawNumbersGameAllow  string `json:"draw_numbers_game_allow" example:"open、close"`     // 搖號抽獎遊戲是否允許重複中獎
	TugofwarGameAllow     string `json:"tugofwar_game_allow" example:"open、close"`         // 拔河遊戲是否允許重複中獎
	BingoGameAllow        string `json:"bingo_game_allow" example:"open、close"`            // 賓果遊戲是否允許重複中獎
	GachaMachineGameAllow string `json:"3d_gacha_machine_game_allow" example:"open、close"` // 扭蛋機遊戲是否允許重複中獎
	VoteGameAllow         string `json:"vote_game_allow" example:"open、close"`             // 投票遊戲是否允許重複中獎
	AllGameAllow          string `json:"all_game_allow" example:"open、close"`              // 所有遊戲是否允許重複中獎

	// 敲敲樂自定義
	WhackmoleClassicHPic01 string `json:"whackmole_classic_h_pic_01" example:"picture"`
	WhackmoleClassicHPic02 string `json:"whackmole_classic_h_pic_02" example:"picture"`
	WhackmoleClassicHPic03 string `json:"whackmole_classic_h_pic_03" example:"picture"`
	WhackmoleClassicHPic04 string `json:"whackmole_classic_h_pic_04" example:"picture"`
	WhackmoleClassicHPic05 string `json:"whackmole_classic_h_pic_05" example:"picture"`
	WhackmoleClassicHPic06 string `json:"whackmole_classic_h_pic_06" example:"picture"`
	WhackmoleClassicHPic07 string `json:"whackmole_classic_h_pic_07" example:"picture"`
	WhackmoleClassicHPic08 string `json:"whackmole_classic_h_pic_08" example:"picture"`
	WhackmoleClassicHPic09 string `json:"whackmole_classic_h_pic_09" example:"picture"`
	WhackmoleClassicHPic10 string `json:"whackmole_classic_h_pic_10" example:"picture"`
	WhackmoleClassicHPic11 string `json:"whackmole_classic_h_pic_11" example:"picture"`
	WhackmoleClassicHPic12 string `json:"whackmole_classic_h_pic_12" example:"picture"`
	WhackmoleClassicHPic13 string `json:"whackmole_classic_h_pic_13" example:"picture"`
	WhackmoleClassicHPic14 string `json:"whackmole_classic_h_pic_14" example:"picture"`
	WhackmoleClassicHPic15 string `json:"whackmole_classic_h_pic_15" example:"picture"`
	WhackmoleClassicGPic01 string `json:"whackmole_classic_g_pic_01" example:"picture"`
	WhackmoleClassicGPic02 string `json:"whackmole_classic_g_pic_02" example:"picture"`
	WhackmoleClassicGPic03 string `json:"whackmole_classic_g_pic_03" example:"picture"`
	WhackmoleClassicGPic04 string `json:"whackmole_classic_g_pic_04" example:"picture"`
	WhackmoleClassicGPic05 string `json:"whackmole_classic_g_pic_05" example:"picture"`
	WhackmoleClassicCPic01 string `json:"whackmole_classic_c_pic_01" example:"picture"`
	WhackmoleClassicCPic02 string `json:"whackmole_classic_c_pic_02" example:"picture"`
	WhackmoleClassicCPic03 string `json:"whackmole_classic_c_pic_03" example:"picture"`
	WhackmoleClassicCPic04 string `json:"whackmole_classic_c_pic_04" example:"picture"`
	WhackmoleClassicCPic05 string `json:"whackmole_classic_c_pic_05" example:"picture"`
	WhackmoleClassicCPic06 string `json:"whackmole_classic_c_pic_06" example:"picture"`
	WhackmoleClassicCPic07 string `json:"whackmole_classic_c_pic_07" example:"picture"`
	WhackmoleClassicCPic08 string `json:"whackmole_classic_c_pic_08" example:"picture"`
	WhackmoleClassicCAni01 string `json:"whackmole_classic_c_ani_01" example:"picture"`

	WhackmoleHalloweenHPic01 string `json:"whackmole_halloween_h_pic_01" example:"picture"`
	WhackmoleHalloweenHPic02 string `json:"whackmole_halloween_h_pic_02" example:"picture"`
	WhackmoleHalloweenHPic03 string `json:"whackmole_halloween_h_pic_03" example:"picture"`
	WhackmoleHalloweenHPic04 string `json:"whackmole_halloween_h_pic_04" example:"picture"`
	WhackmoleHalloweenHPic05 string `json:"whackmole_halloween_h_pic_05" example:"picture"`
	WhackmoleHalloweenHPic06 string `json:"whackmole_halloween_h_pic_06" example:"picture"`
	WhackmoleHalloweenHPic07 string `json:"whackmole_halloween_h_pic_07" example:"picture"`
	WhackmoleHalloweenHPic08 string `json:"whackmole_halloween_h_pic_08" example:"picture"`
	WhackmoleHalloweenHPic09 string `json:"whackmole_halloween_h_pic_09" example:"picture"`
	WhackmoleHalloweenHPic10 string `json:"whackmole_halloween_h_pic_10" example:"picture"`
	WhackmoleHalloweenHPic11 string `json:"whackmole_halloween_h_pic_11" example:"picture"`
	WhackmoleHalloweenHPic12 string `json:"whackmole_halloween_h_pic_12" example:"picture"`
	WhackmoleHalloweenHPic13 string `json:"whackmole_halloween_h_pic_13" example:"picture"`
	WhackmoleHalloweenHPic14 string `json:"whackmole_halloween_h_pic_14" example:"picture"`
	WhackmoleHalloweenHPic15 string `json:"whackmole_halloween_h_pic_15" example:"picture"`
	WhackmoleHalloweenGPic01 string `json:"whackmole_halloween_g_pic_01" example:"picture"`
	WhackmoleHalloweenGPic02 string `json:"whackmole_halloween_g_pic_02" example:"picture"`
	WhackmoleHalloweenGPic03 string `json:"whackmole_halloween_g_pic_03" example:"picture"`
	WhackmoleHalloweenGPic04 string `json:"whackmole_halloween_g_pic_04" example:"picture"`
	WhackmoleHalloweenGPic05 string `json:"whackmole_halloween_g_pic_05" example:"picture"`
	WhackmoleHalloweenCPic01 string `json:"whackmole_halloween_c_pic_01" example:"picture"`
	WhackmoleHalloweenCPic02 string `json:"whackmole_halloween_c_pic_02" example:"picture"`
	WhackmoleHalloweenCPic03 string `json:"whackmole_halloween_c_pic_03" example:"picture"`
	WhackmoleHalloweenCPic04 string `json:"whackmole_halloween_c_pic_04" example:"picture"`
	WhackmoleHalloweenCPic05 string `json:"whackmole_halloween_c_pic_05" example:"picture"`
	WhackmoleHalloweenCPic06 string `json:"whackmole_halloween_c_pic_06" example:"picture"`
	WhackmoleHalloweenCPic07 string `json:"whackmole_halloween_c_pic_07" example:"picture"`
	WhackmoleHalloweenCPic08 string `json:"whackmole_halloween_c_pic_08" example:"picture"`
	WhackmoleHalloweenCAni01 string `json:"whackmole_halloween_c_ani_01" example:"picture"`

	WhackmoleChristmasHPic01 string `json:"whackmole_christmas_h_pic_01" example:"picture"`
	WhackmoleChristmasHPic02 string `json:"whackmole_christmas_h_pic_02" example:"picture"`
	WhackmoleChristmasHPic03 string `json:"whackmole_christmas_h_pic_03" example:"picture"`
	WhackmoleChristmasHPic04 string `json:"whackmole_christmas_h_pic_04" example:"picture"`
	WhackmoleChristmasHPic05 string `json:"whackmole_christmas_h_pic_05" example:"picture"`
	WhackmoleChristmasHPic06 string `json:"whackmole_christmas_h_pic_06" example:"picture"`
	WhackmoleChristmasHPic07 string `json:"whackmole_christmas_h_pic_07" example:"picture"`
	WhackmoleChristmasHPic08 string `json:"whackmole_christmas_h_pic_08" example:"picture"`
	WhackmoleChristmasHPic09 string `json:"whackmole_christmas_h_pic_09" example:"picture"`
	WhackmoleChristmasHPic10 string `json:"whackmole_christmas_h_pic_10" example:"picture"`
	WhackmoleChristmasHPic11 string `json:"whackmole_christmas_h_pic_11" example:"picture"`
	WhackmoleChristmasHPic12 string `json:"whackmole_christmas_h_pic_12" example:"picture"`
	WhackmoleChristmasHPic13 string `json:"whackmole_christmas_h_pic_13" example:"picture"`
	WhackmoleChristmasHPic14 string `json:"whackmole_christmas_h_pic_14" example:"picture"`
	WhackmoleChristmasGPic01 string `json:"whackmole_christmas_g_pic_01" example:"picture"`
	WhackmoleChristmasGPic02 string `json:"whackmole_christmas_g_pic_02" example:"picture"`
	WhackmoleChristmasGPic03 string `json:"whackmole_christmas_g_pic_03" example:"picture"`
	WhackmoleChristmasGPic04 string `json:"whackmole_christmas_g_pic_04" example:"picture"`
	WhackmoleChristmasGPic05 string `json:"whackmole_christmas_g_pic_05" example:"picture"`
	WhackmoleChristmasGPic06 string `json:"whackmole_christmas_g_pic_06" example:"picture"`
	WhackmoleChristmasGPic07 string `json:"whackmole_christmas_g_pic_07" example:"picture"`
	WhackmoleChristmasGPic08 string `json:"whackmole_christmas_g_pic_08" example:"picture"`
	WhackmoleChristmasCPic01 string `json:"whackmole_christmas_c_pic_01" example:"picture"`
	WhackmoleChristmasCPic02 string `json:"whackmole_christmas_c_pic_02" example:"picture"`
	WhackmoleChristmasCPic03 string `json:"whackmole_christmas_c_pic_03" example:"picture"`
	WhackmoleChristmasCPic04 string `json:"whackmole_christmas_c_pic_04" example:"picture"`
	WhackmoleChristmasCPic05 string `json:"whackmole_christmas_c_pic_05" example:"picture"`
	WhackmoleChristmasCPic06 string `json:"whackmole_christmas_c_pic_06" example:"picture"`
	WhackmoleChristmasCPic07 string `json:"whackmole_christmas_c_pic_07" example:"picture"`
	WhackmoleChristmasCPic08 string `json:"whackmole_christmas_c_pic_08" example:"picture"`
	WhackmoleChristmasCAni01 string `json:"whackmole_christmas_c_ani_01" example:"picture"`
	WhackmoleChristmasCAni02 string `json:"whackmole_christmas_c_ani_02" example:"picture"`

	// 敲敲樂音樂
	WhackmoleBgmStart  string `json:"whackmole_bgm_start" example:"picture"`  // 遊戲開始
	WhackmoleBgmGaming string `json:"whackmole_bgm_gaming" example:"picture"` // 遊戲進行中
	WhackmoleBgmEnd    string `json:"whackmole_bgm_end" example:"picture"`    // 遊戲結束

	// 搖號抽獎自定義
	DrawNumbersClassicHPic01 string `json:"draw_numbers_classic_h_pic_01" example:"picture"`
	DrawNumbersClassicHPic02 string `json:"draw_numbers_classic_h_pic_02" example:"picture"`
	DrawNumbersClassicHPic03 string `json:"draw_numbers_classic_h_pic_03" example:"picture"`
	DrawNumbersClassicHPic04 string `json:"draw_numbers_classic_h_pic_04" example:"picture"`
	DrawNumbersClassicHPic05 string `json:"draw_numbers_classic_h_pic_05" example:"picture"`
	DrawNumbersClassicHPic06 string `json:"draw_numbers_classic_h_pic_06" example:"picture"`
	DrawNumbersClassicHPic07 string `json:"draw_numbers_classic_h_pic_07" example:"picture"`
	DrawNumbersClassicHPic08 string `json:"draw_numbers_classic_h_pic_08" example:"picture"`
	DrawNumbersClassicHPic09 string `json:"draw_numbers_classic_h_pic_09" example:"picture"`
	DrawNumbersClassicHPic10 string `json:"draw_numbers_classic_h_pic_10" example:"picture"`
	DrawNumbersClassicHPic11 string `json:"draw_numbers_classic_h_pic_11" example:"picture"`
	DrawNumbersClassicHPic12 string `json:"draw_numbers_classic_h_pic_12" example:"picture"`
	DrawNumbersClassicHPic13 string `json:"draw_numbers_classic_h_pic_13" example:"picture"`
	DrawNumbersClassicHPic14 string `json:"draw_numbers_classic_h_pic_14" example:"picture"`
	DrawNumbersClassicHPic15 string `json:"draw_numbers_classic_h_pic_15" example:"picture"`
	DrawNumbersClassicHPic16 string `json:"draw_numbers_classic_h_pic_16" example:"picture"`
	DrawNumbersClassicHAni01 string `json:"draw_numbers_classic_h_ani_01" example:"picture"`

	DrawNumbersGoldHPic01 string `json:"draw_numbers_gold_h_pic_01" example:"picture"`
	DrawNumbersGoldHPic02 string `json:"draw_numbers_gold_h_pic_02" example:"picture"`
	DrawNumbersGoldHPic03 string `json:"draw_numbers_gold_h_pic_03" example:"picture"`
	DrawNumbersGoldHPic04 string `json:"draw_numbers_gold_h_pic_04" example:"picture"`
	DrawNumbersGoldHPic05 string `json:"draw_numbers_gold_h_pic_05" example:"picture"`
	DrawNumbersGoldHPic06 string `json:"draw_numbers_gold_h_pic_06" example:"picture"`
	DrawNumbersGoldHPic07 string `json:"draw_numbers_gold_h_pic_07" example:"picture"`
	DrawNumbersGoldHPic08 string `json:"draw_numbers_gold_h_pic_08" example:"picture"`
	DrawNumbersGoldHPic09 string `json:"draw_numbers_gold_h_pic_09" example:"picture"`
	DrawNumbersGoldHPic10 string `json:"draw_numbers_gold_h_pic_10" example:"picture"`
	DrawNumbersGoldHPic11 string `json:"draw_numbers_gold_h_pic_11" example:"picture"`
	DrawNumbersGoldHPic12 string `json:"draw_numbers_gold_h_pic_12" example:"picture"`
	DrawNumbersGoldHPic13 string `json:"draw_numbers_gold_h_pic_13" example:"picture"`
	DrawNumbersGoldHPic14 string `json:"draw_numbers_gold_h_pic_14" example:"picture"`
	DrawNumbersGoldHAni01 string `json:"draw_numbers_gold_h_ani_01" example:"picture"`
	DrawNumbersGoldHAni02 string `json:"draw_numbers_gold_h_ani_02" example:"picture"`
	DrawNumbersGoldHAni03 string `json:"draw_numbers_gold_h_ani_03" example:"picture"`

	DrawNumbersNewyearDragonHPic01 string `json:"draw_numbers_newyear_dragon_h_pic_01" example:"picture"`
	DrawNumbersNewyearDragonHPic02 string `json:"draw_numbers_newyear_dragon_h_pic_02" example:"picture"`
	DrawNumbersNewyearDragonHPic03 string `json:"draw_numbers_newyear_dragon_h_pic_03" example:"picture"`
	DrawNumbersNewyearDragonHPic04 string `json:"draw_numbers_newyear_dragon_h_pic_04" example:"picture"`
	DrawNumbersNewyearDragonHPic05 string `json:"draw_numbers_newyear_dragon_h_pic_05" example:"picture"`
	DrawNumbersNewyearDragonHPic06 string `json:"draw_numbers_newyear_dragon_h_pic_06" example:"picture"`
	DrawNumbersNewyearDragonHPic07 string `json:"draw_numbers_newyear_dragon_h_pic_07" example:"picture"`
	DrawNumbersNewyearDragonHPic08 string `json:"draw_numbers_newyear_dragon_h_pic_08" example:"picture"`
	DrawNumbersNewyearDragonHPic09 string `json:"draw_numbers_newyear_dragon_h_pic_09" example:"picture"`
	DrawNumbersNewyearDragonHPic10 string `json:"draw_numbers_newyear_dragon_h_pic_10" example:"picture"`
	DrawNumbersNewyearDragonHPic11 string `json:"draw_numbers_newyear_dragon_h_pic_11" example:"picture"`
	DrawNumbersNewyearDragonHPic12 string `json:"draw_numbers_newyear_dragon_h_pic_12" example:"picture"`
	DrawNumbersNewyearDragonHPic13 string `json:"draw_numbers_newyear_dragon_h_pic_13" example:"picture"`
	DrawNumbersNewyearDragonHPic14 string `json:"draw_numbers_newyear_dragon_h_pic_14" example:"picture"`
	DrawNumbersNewyearDragonHPic15 string `json:"draw_numbers_newyear_dragon_h_pic_15" example:"picture"`
	DrawNumbersNewyearDragonHPic16 string `json:"draw_numbers_newyear_dragon_h_pic_16" example:"picture"`
	DrawNumbersNewyearDragonHPic17 string `json:"draw_numbers_newyear_dragon_h_pic_17" example:"picture"`
	DrawNumbersNewyearDragonHPic18 string `json:"draw_numbers_newyear_dragon_h_pic_18" example:"picture"`
	DrawNumbersNewyearDragonHPic19 string `json:"draw_numbers_newyear_dragon_h_pic_19" example:"picture"`
	DrawNumbersNewyearDragonHPic20 string `json:"draw_numbers_newyear_dragon_h_pic_20" example:"picture"`
	DrawNumbersNewyearDragonHAni01 string `json:"draw_numbers_newyear_dragon_h_ani_01" example:"picture"`
	DrawNumbersNewyearDragonHAni02 string `json:"draw_numbers_newyear_dragon_h_ani_02" example:"picture"`

	DrawNumbersCherryHPic01 string `json:"draw_numbers_cherry_h_pic_01" example:"picture"`
	DrawNumbersCherryHPic02 string `json:"draw_numbers_cherry_h_pic_02" example:"picture"`
	DrawNumbersCherryHPic03 string `json:"draw_numbers_cherry_h_pic_03" example:"picture"`
	DrawNumbersCherryHPic04 string `json:"draw_numbers_cherry_h_pic_04" example:"picture"`
	DrawNumbersCherryHPic05 string `json:"draw_numbers_cherry_h_pic_05" example:"picture"`
	DrawNumbersCherryHPic06 string `json:"draw_numbers_cherry_h_pic_06" example:"picture"`
	DrawNumbersCherryHPic07 string `json:"draw_numbers_cherry_h_pic_07" example:"picture"`
	DrawNumbersCherryHPic08 string `json:"draw_numbers_cherry_h_pic_08" example:"picture"`
	DrawNumbersCherryHPic09 string `json:"draw_numbers_cherry_h_pic_09" example:"picture"`
	DrawNumbersCherryHPic10 string `json:"draw_numbers_cherry_h_pic_10" example:"picture"`
	DrawNumbersCherryHPic11 string `json:"draw_numbers_cherry_h_pic_11" example:"picture"`
	DrawNumbersCherryHPic12 string `json:"draw_numbers_cherry_h_pic_12" example:"picture"`
	DrawNumbersCherryHPic13 string `json:"draw_numbers_cherry_h_pic_13" example:"picture"`
	DrawNumbersCherryHPic14 string `json:"draw_numbers_cherry_h_pic_14" example:"picture"`
	DrawNumbersCherryHPic15 string `json:"draw_numbers_cherry_h_pic_15" example:"picture"`
	DrawNumbersCherryHPic16 string `json:"draw_numbers_cherry_h_pic_16" example:"picture"`
	DrawNumbersCherryHPic17 string `json:"draw_numbers_cherry_h_pic_17" example:"picture"`
	DrawNumbersCherryHAni01 string `json:"draw_numbers_cherry_h_ani_01" example:"picture"`
	DrawNumbersCherryHAni02 string `json:"draw_numbers_cherry_h_ani_02" example:"picture"`
	DrawNumbersCherryHAni03 string `json:"draw_numbers_cherry_h_ani_03" example:"picture"`
	DrawNumbersCherryHAni04 string `json:"draw_numbers_cherry_h_ani_04" example:"picture"`

	// 太空主題
	DrawNumbers3DSpaceHPic01 string `json:"draw_numbers_3D_space_h_pic_01" example:"picture"`
	DrawNumbers3DSpaceHPic02 string `json:"draw_numbers_3D_space_h_pic_02" example:"picture"`
	DrawNumbers3DSpaceHPic03 string `json:"draw_numbers_3D_space_h_pic_03" example:"picture"`
	DrawNumbers3DSpaceHPic04 string `json:"draw_numbers_3D_space_h_pic_04" example:"picture"`
	DrawNumbers3DSpaceHPic05 string `json:"draw_numbers_3D_space_h_pic_05" example:"picture"`
	DrawNumbers3DSpaceHPic06 string `json:"draw_numbers_3D_space_h_pic_06" example:"picture"`
	DrawNumbers3DSpaceHPic07 string `json:"draw_numbers_3D_space_h_pic_07" example:"picture"`
	DrawNumbers3DSpaceHPic08 string `json:"draw_numbers_3D_space_h_pic_08" example:"picture"`

	// 音樂
	DrawNumbersBgmGaming string `json:"draw_numbers_bgm_gaming" example:"picture"` // 遊戲進行中

	// 快問快答自定義
	QAClassicHPic01 string `json:"qa_classic_h_pic_01" example:"picture"`
	QAClassicHPic02 string `json:"qa_classic_h_pic_02" example:"picture"`
	QAClassicHPic03 string `json:"qa_classic_h_pic_03" example:"picture"`
	QAClassicHPic04 string `json:"qa_classic_h_pic_04" example:"picture"`
	QAClassicHPic05 string `json:"qa_classic_h_pic_05" example:"picture"`
	QAClassicHPic06 string `json:"qa_classic_h_pic_06" example:"picture"`
	QAClassicHPic07 string `json:"qa_classic_h_pic_07" example:"picture"`
	QAClassicHPic08 string `json:"qa_classic_h_pic_08" example:"picture"`
	QAClassicHPic09 string `json:"qa_classic_h_pic_09" example:"picture"`
	QAClassicHPic10 string `json:"qa_classic_h_pic_10" example:"picture"`
	QAClassicHPic11 string `json:"qa_classic_h_pic_11" example:"picture"`
	QAClassicHPic12 string `json:"qa_classic_h_pic_12" example:"picture"`
	QAClassicHPic13 string `json:"qa_classic_h_pic_13" example:"picture"`
	QAClassicHPic14 string `json:"qa_classic_h_pic_14" example:"picture"`
	QAClassicHPic15 string `json:"qa_classic_h_pic_15" example:"picture"`
	QAClassicHPic16 string `json:"qa_classic_h_pic_16" example:"picture"`
	QAClassicHPic17 string `json:"qa_classic_h_pic_17" example:"picture"`
	QAClassicHPic18 string `json:"qa_classic_h_pic_18" example:"picture"`
	QAClassicHPic19 string `json:"qa_classic_h_pic_19" example:"picture"`
	QAClassicHPic20 string `json:"qa_classic_h_pic_20" example:"picture"`
	QAClassicHPic21 string `json:"qa_classic_h_pic_21" example:"picture"`
	QAClassicHPic22 string `json:"qa_classic_h_pic_22" example:"picture"`
	QAClassicGPic01 string `json:"qa_classic_g_pic_01" example:"picture"`
	QAClassicGPic02 string `json:"qa_classic_g_pic_02" example:"picture"`
	QAClassicGPic03 string `json:"qa_classic_g_pic_03" example:"picture"`
	QAClassicGPic04 string `json:"qa_classic_g_pic_04" example:"picture"`
	QAClassicGPic05 string `json:"qa_classic_g_pic_05" example:"picture"`
	QAClassicCPic01 string `json:"qa_classic_c_pic_01" example:"picture"`
	QAClassicHAni01 string `json:"qa_classic_h_ani_01" example:"picture"`
	QAClassicHAni02 string `json:"qa_classic_h_ani_02" example:"picture"`
	QAClassicGAni01 string `json:"qa_classic_g_ani_01" example:"picture"`
	QAClassicGAni02 string `json:"qa_classic_g_ani_02" example:"picture"`

	QAElectricHPic01 string `json:"qa_electric_h_pic_01" example:"picture"`
	QAElectricHPic02 string `json:"qa_electric_h_pic_02" example:"picture"`
	QAElectricHPic03 string `json:"qa_electric_h_pic_03" example:"picture"`
	QAElectricHPic04 string `json:"qa_electric_h_pic_04" example:"picture"`
	QAElectricHPic05 string `json:"qa_electric_h_pic_05" example:"picture"`
	QAElectricHPic06 string `json:"qa_electric_h_pic_06" example:"picture"`
	QAElectricHPic07 string `json:"qa_electric_h_pic_07" example:"picture"`
	QAElectricHPic08 string `json:"qa_electric_h_pic_08" example:"picture"`
	QAElectricHPic09 string `json:"qa_electric_h_pic_09" example:"picture"`
	QAElectricHPic10 string `json:"qa_electric_h_pic_10" example:"picture"`
	QAElectricHPic11 string `json:"qa_electric_h_pic_11" example:"picture"`
	QAElectricHPic12 string `json:"qa_electric_h_pic_12" example:"picture"`
	QAElectricHPic13 string `json:"qa_electric_h_pic_13" example:"picture"`
	QAElectricHPic14 string `json:"qa_electric_h_pic_14" example:"picture"`
	QAElectricHPic15 string `json:"qa_electric_h_pic_15" example:"picture"`
	QAElectricHPic16 string `json:"qa_electric_h_pic_16" example:"picture"`
	QAElectricHPic17 string `json:"qa_electric_h_pic_17" example:"picture"`
	QAElectricHPic18 string `json:"qa_electric_h_pic_18" example:"picture"`
	QAElectricHPic19 string `json:"qa_electric_h_pic_19" example:"picture"`
	QAElectricHPic20 string `json:"qa_electric_h_pic_20" example:"picture"`
	QAElectricHPic21 string `json:"qa_electric_h_pic_21" example:"picture"`
	QAElectricHPic22 string `json:"qa_electric_h_pic_22" example:"picture"`
	QAElectricHPic23 string `json:"qa_electric_h_pic_23" example:"picture"`
	QAElectricHPic24 string `json:"qa_electric_h_pic_24" example:"picture"`
	QAElectricHPic25 string `json:"qa_electric_h_pic_25" example:"picture"`
	QAElectricHPic26 string `json:"qa_electric_h_pic_26" example:"picture"`
	QAElectricGPic01 string `json:"qa_electric_g_pic_01" example:"picture"`
	QAElectricGPic02 string `json:"qa_electric_g_pic_02" example:"picture"`
	QAElectricGPic03 string `json:"qa_electric_g_pic_03" example:"picture"`
	QAElectricGPic04 string `json:"qa_electric_g_pic_04" example:"picture"`
	QAElectricGPic05 string `json:"qa_electric_g_pic_05" example:"picture"`
	QAElectricGPic06 string `json:"qa_electric_g_pic_06" example:"picture"`
	QAElectricGPic07 string `json:"qa_electric_g_pic_07" example:"picture"`
	QAElectricGPic08 string `json:"qa_electric_g_pic_08" example:"picture"`
	QAElectricGPic09 string `json:"qa_electric_g_pic_09" example:"picture"`
	QAElectricCPic01 string `json:"qa_electric_c_pic_01" example:"picture"`
	QAElectricHAni01 string `json:"qa_electric_h_ani_01" example:"picture"`
	QAElectricHAni02 string `json:"qa_electric_h_ani_02" example:"picture"`
	QAElectricHAni03 string `json:"qa_electric_h_ani_03" example:"picture"`
	QAElectricHAni04 string `json:"qa_electric_h_ani_04" example:"picture"`
	QAElectricHAni05 string `json:"qa_electric_h_ani_05" example:"picture"`
	QAElectricGAni01 string `json:"qa_electric_g_ani_01" example:"picture"`
	QAElectricGAni02 string `json:"qa_electric_g_ani_02" example:"picture"`
	QAElectricCAni01 string `json:"qa_electric_c_ani_01" example:"picture"`

	QAMoonfestivalHPic01 string `json:"qa_moonfestival_h_pic_01" example:"picture"`
	QAMoonfestivalHPic02 string `json:"qa_moonfestival_h_pic_02" example:"picture"`
	QAMoonfestivalHPic03 string `json:"qa_moonfestival_h_pic_03" example:"picture"`
	QAMoonfestivalHPic04 string `json:"qa_moonfestival_h_pic_04" example:"picture"`
	QAMoonfestivalHPic05 string `json:"qa_moonfestival_h_pic_05" example:"picture"`
	QAMoonfestivalHPic06 string `json:"qa_moonfestival_h_pic_06" example:"picture"`
	QAMoonfestivalHPic07 string `json:"qa_moonfestival_h_pic_07" example:"picture"`
	QAMoonfestivalHPic08 string `json:"qa_moonfestival_h_pic_08" example:"picture"`
	QAMoonfestivalHPic09 string `json:"qa_moonfestival_h_pic_09" example:"picture"`
	QAMoonfestivalHPic10 string `json:"qa_moonfestival_h_pic_10" example:"picture"`
	QAMoonfestivalHPic11 string `json:"qa_moonfestival_h_pic_11" example:"picture"`
	QAMoonfestivalHPic12 string `json:"qa_moonfestival_h_pic_12" example:"picture"`
	QAMoonfestivalHPic13 string `json:"qa_moonfestival_h_pic_13" example:"picture"`
	QAMoonfestivalHPic14 string `json:"qa_moonfestival_h_pic_14" example:"picture"`
	QAMoonfestivalHPic15 string `json:"qa_moonfestival_h_pic_15" example:"picture"`
	QAMoonfestivalHPic16 string `json:"qa_moonfestival_h_pic_16" example:"picture"`
	QAMoonfestivalHPic17 string `json:"qa_moonfestival_h_pic_17" example:"picture"`
	QAMoonfestivalHPic18 string `json:"qa_moonfestival_h_pic_18" example:"picture"`
	QAMoonfestivalHPic19 string `json:"qa_moonfestival_h_pic_19" example:"picture"`
	QAMoonfestivalHPic20 string `json:"qa_moonfestival_h_pic_20" example:"picture"`
	QAMoonfestivalHPic21 string `json:"qa_moonfestival_h_pic_21" example:"picture"`
	QAMoonfestivalHPic22 string `json:"qa_moonfestival_h_pic_22" example:"picture"`
	QAMoonfestivalHPic23 string `json:"qa_moonfestival_h_pic_23" example:"picture"`
	QAMoonfestivalHPic24 string `json:"qa_moonfestival_h_pic_24" example:"picture"`
	QAMoonfestivalGPic01 string `json:"qa_moonfestival_g_pic_01" example:"picture"`
	QAMoonfestivalGPic02 string `json:"qa_moonfestival_g_pic_02" example:"picture"`
	QAMoonfestivalGPic03 string `json:"qa_moonfestival_g_pic_03" example:"picture"`
	QAMoonfestivalGPic04 string `json:"qa_moonfestival_g_pic_04" example:"picture"`
	QAMoonfestivalGPic05 string `json:"qa_moonfestival_g_pic_05" example:"picture"`
	QAMoonfestivalCPic01 string `json:"qa_moonfestival_c_pic_01" example:"picture"`
	QAMoonfestivalCPic02 string `json:"qa_moonfestival_c_pic_02" example:"picture"`
	QAMoonfestivalCPic03 string `json:"qa_moonfestival_c_pic_03" example:"picture"`
	QAMoonfestivalHAni01 string `json:"qa_moonfestival_h_ani_01" example:"picture"`
	QAMoonfestivalHAni02 string `json:"qa_moonfestival_h_ani_02" example:"picture"`
	QAMoonfestivalGAni01 string `json:"qa_moonfestival_g_ani_01" example:"picture"`
	QAMoonfestivalGAni02 string `json:"qa_moonfestival_g_ani_02" example:"picture"`
	QAMoonfestivalGAni03 string `json:"qa_moonfestival_g_ani_03" example:"picture"`

	QANewyearDragonHPic01 string `json:"qa_newyear_dragon_h_pic_01" example:"picture"`
	QANewyearDragonHPic02 string `json:"qa_newyear_dragon_h_pic_02" example:"picture"`
	QANewyearDragonHPic03 string `json:"qa_newyear_dragon_h_pic_03" example:"picture"`
	QANewyearDragonHPic04 string `json:"qa_newyear_dragon_h_pic_04" example:"picture"`
	QANewyearDragonHPic05 string `json:"qa_newyear_dragon_h_pic_05" example:"picture"`
	QANewyearDragonHPic06 string `json:"qa_newyear_dragon_h_pic_06" example:"picture"`
	QANewyearDragonHPic07 string `json:"qa_newyear_dragon_h_pic_07" example:"picture"`
	QANewyearDragonHPic08 string `json:"qa_newyear_dragon_h_pic_08" example:"picture"`
	QANewyearDragonHPic09 string `json:"qa_newyear_dragon_h_pic_09" example:"picture"`
	QANewyearDragonHPic10 string `json:"qa_newyear_dragon_h_pic_10" example:"picture"`
	QANewyearDragonHPic11 string `json:"qa_newyear_dragon_h_pic_11" example:"picture"`
	QANewyearDragonHPic12 string `json:"qa_newyear_dragon_h_pic_12" example:"picture"`
	QANewyearDragonHPic13 string `json:"qa_newyear_dragon_h_pic_13" example:"picture"`
	QANewyearDragonHPic14 string `json:"qa_newyear_dragon_h_pic_14" example:"picture"`
	QANewyearDragonHPic15 string `json:"qa_newyear_dragon_h_pic_15" example:"picture"`
	QANewyearDragonHPic16 string `json:"qa_newyear_dragon_h_pic_16" example:"picture"`
	QANewyearDragonHPic17 string `json:"qa_newyear_dragon_h_pic_17" example:"picture"`
	QANewyearDragonHPic18 string `json:"qa_newyear_dragon_h_pic_18" example:"picture"`
	QANewyearDragonHPic19 string `json:"qa_newyear_dragon_h_pic_19" example:"picture"`
	QANewyearDragonHPic20 string `json:"qa_newyear_dragon_h_pic_20" example:"picture"`
	QANewyearDragonHPic21 string `json:"qa_newyear_dragon_h_pic_21" example:"picture"`
	QANewyearDragonHPic22 string `json:"qa_newyear_dragon_h_pic_22" example:"picture"`
	QANewyearDragonHPic23 string `json:"qa_newyear_dragon_h_pic_23" example:"picture"`
	QANewyearDragonHPic24 string `json:"qa_newyear_dragon_h_pic_24" example:"picture"`
	QANewyearDragonGPic01 string `json:"qa_newyear_dragon_g_pic_01" example:"picture"`
	QANewyearDragonGPic02 string `json:"qa_newyear_dragon_g_pic_02" example:"picture"`
	QANewyearDragonGPic03 string `json:"qa_newyear_dragon_g_pic_03" example:"picture"`
	QANewyearDragonGPic04 string `json:"qa_newyear_dragon_g_pic_04" example:"picture"`
	QANewyearDragonGPic05 string `json:"qa_newyear_dragon_g_pic_05" example:"picture"`
	QANewyearDragonGPic06 string `json:"qa_newyear_dragon_g_pic_06" example:"picture"`
	QANewyearDragonCPic01 string `json:"qa_newyear_dragon_c_pic_01" example:"picture"`
	QANewyearDragonHAni01 string `json:"qa_newyear_dragon_h_ani_01" example:"picture"`
	QANewyearDragonHAni02 string `json:"qa_newyear_dragon_h_ani_02" example:"picture"`
	QANewyearDragonGAni01 string `json:"qa_newyear_dragon_g_ani_01" example:"picture"`
	QANewyearDragonGAni02 string `json:"qa_newyear_dragon_g_ani_02" example:"picture"`
	QANewyearDragonGAni03 string `json:"qa_newyear_dragon_g_ani_03" example:"picture"`
	QANewyearDragonCAni01 string `json:"qa_newyear_dragon_c_ani_01" example:"picture"`

	// 音樂
	QABgmStart  string `json:"qa_bgm_start" example:"picture"`  // 遊戲開始
	QABgmGaming string `json:"qa_bgm_gaming" example:"picture"` // 遊戲進行中
	QABgmEnd    string `json:"qa_bgm_end" example:"picture"`    // 遊戲結束

	// 搖紅包自定義
	RedpackClassicHPic01 string `json:"redpack_classic_h_pic_01" example:"picture"`
	RedpackClassicHPic02 string `json:"redpack_classic_h_pic_02" example:"picture"`
	RedpackClassicHPic03 string `json:"redpack_classic_h_pic_03" example:"picture"`
	RedpackClassicHPic04 string `json:"redpack_classic_h_pic_04" example:"picture"`
	RedpackClassicHPic05 string `json:"redpack_classic_h_pic_05" example:"picture"`
	RedpackClassicHPic06 string `json:"redpack_classic_h_pic_06" example:"picture"`
	RedpackClassicHPic07 string `json:"redpack_classic_h_pic_07" example:"picture"`
	RedpackClassicHPic08 string `json:"redpack_classic_h_pic_08" example:"picture"`
	RedpackClassicHPic09 string `json:"redpack_classic_h_pic_09" example:"picture"`
	RedpackClassicHPic10 string `json:"redpack_classic_h_pic_10" example:"picture"`
	RedpackClassicHPic11 string `json:"redpack_classic_h_pic_11" example:"picture"`
	RedpackClassicHPic12 string `json:"redpack_classic_h_pic_12" example:"picture"`
	RedpackClassicHPic13 string `json:"redpack_classic_h_pic_13" example:"picture"`
	RedpackClassicGPic01 string `json:"redpack_classic_g_pic_01" example:"picture"`
	RedpackClassicGPic02 string `json:"redpack_classic_g_pic_02" example:"picture"`
	RedpackClassicGPic03 string `json:"redpack_classic_g_pic_03" example:"picture"`
	RedpackClassicHAni01 string `json:"redpack_classic_h_ani_01" example:"picture"`
	RedpackClassicHAni02 string `json:"redpack_classic_h_ani_02" example:"picture"`
	RedpackClassicGAni01 string `json:"redpack_classic_g_ani_01" example:"picture"`
	RedpackClassicGAni02 string `json:"redpack_classic_g_ani_02" example:"picture"`
	RedpackClassicGAni03 string `json:"redpack_classic_g_ani_03" example:"picture"`

	RedpackCherryHPic01 string `json:"redpack_cherry_h_pic_01" example:"picture"`
	RedpackCherryHPic02 string `json:"redpack_cherry_h_pic_02" example:"picture"`
	RedpackCherryHPic03 string `json:"redpack_cherry_h_pic_03" example:"picture"`
	RedpackCherryHPic04 string `json:"redpack_cherry_h_pic_04" example:"picture"`
	RedpackCherryHPic05 string `json:"redpack_cherry_h_pic_05" example:"picture"`
	RedpackCherryHPic06 string `json:"redpack_cherry_h_pic_06" example:"picture"`
	RedpackCherryHPic07 string `json:"redpack_cherry_h_pic_07" example:"picture"`
	RedpackCherryGPic01 string `json:"redpack_cherry_g_pic_01" example:"picture"`
	RedpackCherryGPic02 string `json:"redpack_cherry_g_pic_02" example:"picture"`
	RedpackCherryHAni01 string `json:"redpack_cherry_h_ani_01" example:"picture"`
	RedpackCherryHAni02 string `json:"redpack_cherry_h_ani_02" example:"picture"`
	RedpackCherryGAni01 string `json:"redpack_cherry_g_ani_01" example:"picture"`
	RedpackCherryGAni02 string `json:"redpack_cherry_g_ani_02" example:"picture"`

	RedpackChristmasHPic01 string `json:"redpack_christmas_h_pic_01" example:"picture"`
	RedpackChristmasHPic02 string `json:"redpack_christmas_h_pic_02" example:"picture"`
	RedpackChristmasHPic03 string `json:"redpack_christmas_h_pic_03" example:"picture"`
	RedpackChristmasHPic04 string `json:"redpack_christmas_h_pic_04" example:"picture"`
	RedpackChristmasHPic05 string `json:"redpack_christmas_h_pic_05" example:"picture"`
	RedpackChristmasHPic06 string `json:"redpack_christmas_h_pic_06" example:"picture"`
	RedpackChristmasHPic07 string `json:"redpack_christmas_h_pic_07" example:"picture"`
	RedpackChristmasHPic08 string `json:"redpack_christmas_h_pic_08" example:"picture"`
	RedpackChristmasHPic09 string `json:"redpack_christmas_h_pic_09" example:"picture"`
	RedpackChristmasHPic10 string `json:"redpack_christmas_h_pic_10" example:"picture"`
	RedpackChristmasHPic11 string `json:"redpack_christmas_h_pic_11" example:"picture"`
	RedpackChristmasHPic12 string `json:"redpack_christmas_h_pic_12" example:"picture"`
	RedpackChristmasHPic13 string `json:"redpack_christmas_h_pic_13" example:"picture"`
	RedpackChristmasGPic01 string `json:"redpack_christmas_g_pic_01" example:"picture"`
	RedpackChristmasGPic02 string `json:"redpack_christmas_g_pic_02" example:"picture"`
	RedpackChristmasGPic03 string `json:"redpack_christmas_g_pic_03" example:"picture"`
	RedpackChristmasGPic04 string `json:"redpack_christmas_g_pic_04" example:"picture"`
	RedpackChristmasCPic01 string `json:"redpack_christmas_c_pic_01" example:"picture"`
	RedpackChristmasCPic02 string `json:"redpack_christmas_c_pic_02" example:"picture"`
	RedpackChristmasCAni01 string `json:"redpack_christmas_c_ani_01" example:"picture"`

	// 音樂
	RedpackBgmStart  string `json:"redpack_bgm_start" example:"picture"`  // 遊戲開始
	RedpackBgmGaming string `json:"redpack_bgm_gaming" example:"picture"` // 遊戲進行中
	RedpackBgmEnd    string `json:"redpack_bgm_end" example:"picture"`    // 遊戲結束

	// 套紅包自定義
	RopepackClassicHPic01 string `json:"ropepack_classic_h_pic_01" example:"picture"`
	RopepackClassicHPic02 string `json:"ropepack_classic_h_pic_02" example:"picture"`
	RopepackClassicHPic03 string `json:"ropepack_classic_h_pic_03" example:"picture"`
	RopepackClassicHPic04 string `json:"ropepack_classic_h_pic_04" example:"picture"`
	RopepackClassicHPic05 string `json:"ropepack_classic_h_pic_05" example:"picture"`
	RopepackClassicHPic06 string `json:"ropepack_classic_h_pic_06" example:"picture"`
	RopepackClassicHPic07 string `json:"ropepack_classic_h_pic_07" example:"picture"`
	RopepackClassicHPic08 string `json:"ropepack_classic_h_pic_08" example:"picture"`
	RopepackClassicHPic09 string `json:"ropepack_classic_h_pic_09" example:"picture"`
	RopepackClassicHPic10 string `json:"ropepack_classic_h_pic_10" example:"picture"`
	RopepackClassicGPic01 string `json:"ropepack_classic_g_pic_01" example:"picture"`
	RopepackClassicGPic02 string `json:"ropepack_classic_g_pic_02" example:"picture"`
	RopepackClassicGPic03 string `json:"ropepack_classic_g_pic_03" example:"picture"`
	RopepackClassicGPic04 string `json:"ropepack_classic_g_pic_04" example:"picture"`
	RopepackClassicGPic05 string `json:"ropepack_classic_g_pic_05" example:"picture"`
	RopepackClassicGPic06 string `json:"ropepack_classic_g_pic_06" example:"picture"`
	RopepackClassicHAni01 string `json:"ropepack_classic_h_ani_01" example:"picture"`
	RopepackClassicGAni01 string `json:"ropepack_classic_g_ani_01" example:"picture"`
	RopepackClassicGAni02 string `json:"ropepack_classic_g_ani_02" example:"picture"`
	RopepackClassicCAni01 string `json:"ropepack_classic_c_ani_01" example:"picture"`

	RopepackNewyearRabbitHPic01 string `json:"ropepack_newyear_rabbit_h_pic_01" example:"picture"`
	RopepackNewyearRabbitHPic02 string `json:"ropepack_newyear_rabbit_h_pic_02" example:"picture"`
	RopepackNewyearRabbitHPic03 string `json:"ropepack_newyear_rabbit_h_pic_03" example:"picture"`
	RopepackNewyearRabbitHPic04 string `json:"ropepack_newyear_rabbit_h_pic_04" example:"picture"`
	RopepackNewyearRabbitHPic05 string `json:"ropepack_newyear_rabbit_h_pic_05" example:"picture"`
	RopepackNewyearRabbitHPic06 string `json:"ropepack_newyear_rabbit_h_pic_06" example:"picture"`
	RopepackNewyearRabbitHPic07 string `json:"ropepack_newyear_rabbit_h_pic_07" example:"picture"`
	RopepackNewyearRabbitHPic08 string `json:"ropepack_newyear_rabbit_h_pic_08" example:"picture"`
	RopepackNewyearRabbitHPic09 string `json:"ropepack_newyear_rabbit_h_pic_09" example:"picture"`
	RopepackNewyearRabbitGPic01 string `json:"ropepack_newyear_rabbit_g_pic_01" example:"picture"`
	RopepackNewyearRabbitGPic02 string `json:"ropepack_newyear_rabbit_g_pic_02" example:"picture"`
	RopepackNewyearRabbitGPic03 string `json:"ropepack_newyear_rabbit_g_pic_03" example:"picture"`
	RopepackNewyearRabbitHAni01 string `json:"ropepack_newyear_rabbit_h_ani_01" example:"picture"`
	RopepackNewyearRabbitGAni01 string `json:"ropepack_newyear_rabbit_g_ani_01" example:"picture"`
	RopepackNewyearRabbitGAni02 string `json:"ropepack_newyear_rabbit_g_ani_02" example:"picture"`
	RopepackNewyearRabbitGAni03 string `json:"ropepack_newyear_rabbit_g_ani_03" example:"picture"`
	RopepackNewyearRabbitCAni01 string `json:"ropepack_newyear_rabbit_c_ani_01" example:"picture"`
	RopepackNewyearRabbitCAni02 string `json:"ropepack_newyear_rabbit_c_ani_02" example:"picture"`

	RopepackMoonfestivalHPic01 string `json:"ropepack_moonfestival_h_pic_01" example:"picture"`
	RopepackMoonfestivalHPic02 string `json:"ropepack_moonfestival_h_pic_02" example:"picture"`
	RopepackMoonfestivalHPic03 string `json:"ropepack_moonfestival_h_pic_03" example:"picture"`
	RopepackMoonfestivalHPic04 string `json:"ropepack_moonfestival_h_pic_04" example:"picture"`
	RopepackMoonfestivalHPic05 string `json:"ropepack_moonfestival_h_pic_05" example:"picture"`
	RopepackMoonfestivalHPic06 string `json:"ropepack_moonfestival_h_pic_06" example:"picture"`
	RopepackMoonfestivalHPic07 string `json:"ropepack_moonfestival_h_pic_07" example:"picture"`
	RopepackMoonfestivalHPic08 string `json:"ropepack_moonfestival_h_pic_08" example:"picture"`
	RopepackMoonfestivalHPic09 string `json:"ropepack_moonfestival_h_pic_09" example:"picture"`
	RopepackMoonfestivalGPic01 string `json:"ropepack_moonfestival_g_pic_01" example:"picture"`
	RopepackMoonfestivalGPic02 string `json:"ropepack_moonfestival_g_pic_02" example:"picture"`
	RopepackMoonfestivalCPic01 string `json:"ropepack_moonfestival_c_pic_01" example:"picture"`
	RopepackMoonfestivalHAni01 string `json:"ropepack_moonfestival_h_ani_01" example:"picture"`
	RopepackMoonfestivalGAni01 string `json:"ropepack_moonfestival_g_ani_01" example:"picture"`
	RopepackMoonfestivalGAni02 string `json:"ropepack_moonfestival_g_ani_02" example:"picture"`
	RopepackMoonfestivalCAni01 string `json:"ropepack_moonfestival_c_ani_01" example:"picture"`
	RopepackMoonfestivalCAni02 string `json:"ropepack_moonfestival_c_ani_02" example:"picture"`

	Ropepack3DHPic01 string `json:"ropepack_3D_h_pic_01" example:"picture"`
	Ropepack3DHPic02 string `json:"ropepack_3D_h_pic_02" example:"picture"`
	Ropepack3DHPic03 string `json:"ropepack_3D_h_pic_03" example:"picture"`
	Ropepack3DHPic04 string `json:"ropepack_3D_h_pic_04" example:"picture"`
	Ropepack3DHPic05 string `json:"ropepack_3D_h_pic_05" example:"picture"`
	Ropepack3DHPic06 string `json:"ropepack_3D_h_pic_06" example:"picture"`
	Ropepack3DHPic07 string `json:"ropepack_3D_h_pic_07" example:"picture"`
	Ropepack3DHPic08 string `json:"ropepack_3D_h_pic_08" example:"picture"`
	Ropepack3DHPic09 string `json:"ropepack_3D_h_pic_09" example:"picture"`
	Ropepack3DHPic10 string `json:"ropepack_3D_h_pic_10" example:"picture"`
	Ropepack3DHPic11 string `json:"ropepack_3D_h_pic_11" example:"picture"`
	Ropepack3DHPic12 string `json:"ropepack_3D_h_pic_12" example:"picture"`
	Ropepack3DHPic13 string `json:"ropepack_3D_h_pic_13" example:"picture"`
	Ropepack3DHPic14 string `json:"ropepack_3D_h_pic_14" example:"picture"`
	Ropepack3DHPic15 string `json:"ropepack_3D_h_pic_15" example:"picture"`
	Ropepack3DGPic01 string `json:"ropepack_3D_g_pic_01" example:"picture"`
	Ropepack3DGPic02 string `json:"ropepack_3D_g_pic_02" example:"picture"`
	Ropepack3DGPic03 string `json:"ropepack_3D_g_pic_03" example:"picture"`
	Ropepack3DGPic04 string `json:"ropepack_3D_g_pic_04" example:"picture"`
	Ropepack3DHAni01 string `json:"ropepack_3D_h_ani_01" example:"picture"`
	Ropepack3DHAni02 string `json:"ropepack_3D_h_ani_02" example:"picture"`
	Ropepack3DHAni03 string `json:"ropepack_3D_h_ani_03" example:"picture"`
	Ropepack3DGAni01 string `json:"ropepack_3D_g_ani_01" example:"picture"`
	Ropepack3DGAni02 string `json:"ropepack_3D_g_ani_02" example:"picture"`
	Ropepack3DCAni01 string `json:"ropepack_3D_c_ani_01" example:"picture"`

	// 音樂
	RopepackBgmStart  string `json:"ropepack_bgm_start" example:"picture"`  // 遊戲開始
	RopepackBgmGaming string `json:"ropepack_bgm_gaming" example:"picture"` // 遊戲進行中
	RopepackBgmEnd    string `json:"ropepack_bgm_end" example:"picture"`    // 遊戲結束

	// 遊戲抽獎自定義
	LotteryJiugonggeClassicHPic01 string `json:"lottery_jiugongge_classic_h_pic_01" example:"picture"`
	LotteryJiugonggeClassicHPic02 string `json:"lottery_jiugongge_classic_h_pic_02" example:"picture"`
	LotteryJiugonggeClassicHPic03 string `json:"lottery_jiugongge_classic_h_pic_03" example:"picture"`
	LotteryJiugonggeClassicHPic04 string `json:"lottery_jiugongge_classic_h_pic_04" example:"picture"`
	LotteryJiugonggeClassicGPic01 string `json:"lottery_jiugongge_classic_g_pic_01" example:"picture"`
	LotteryJiugonggeClassicGPic02 string `json:"lottery_jiugongge_classic_g_pic_02" example:"picture"`
	LotteryJiugonggeClassicCPic01 string `json:"lottery_jiugongge_classic_c_pic_01" example:"picture"`
	LotteryJiugonggeClassicCPic02 string `json:"lottery_jiugongge_classic_c_pic_02" example:"picture"`
	LotteryJiugonggeClassicCPic03 string `json:"lottery_jiugongge_classic_c_pic_03" example:"picture"`
	LotteryJiugonggeClassicCPic04 string `json:"lottery_jiugongge_classic_c_pic_04" example:"picture"`
	LotteryJiugonggeClassicCAni01 string `json:"lottery_jiugongge_classic_c_ani_01" example:"picture"`
	LotteryJiugonggeClassicCAni02 string `json:"lottery_jiugongge_classic_c_ani_02" example:"picture"`
	LotteryJiugonggeClassicCAni03 string `json:"lottery_jiugongge_classic_c_ani_03" example:"picture"`

	LotteryTurntableClassicHPic01 string `json:"lottery_turntable_classic_h_pic_01" example:"picture"`
	LotteryTurntableClassicHPic02 string `json:"lottery_turntable_classic_h_pic_02" example:"picture"`
	LotteryTurntableClassicHPic03 string `json:"lottery_turntable_classic_h_pic_03" example:"picture"`
	LotteryTurntableClassicHPic04 string `json:"lottery_turntable_classic_h_pic_04" example:"picture"`
	LotteryTurntableClassicGPic01 string `json:"lottery_turntable_classic_g_pic_01" example:"picture"`
	LotteryTurntableClassicGPic02 string `json:"lottery_turntable_classic_g_pic_02" example:"picture"`
	LotteryTurntableClassicCPic01 string `json:"lottery_turntable_classic_c_pic_01" example:"picture"`
	LotteryTurntableClassicCPic02 string `json:"lottery_turntable_classic_c_pic_02" example:"picture"`
	LotteryTurntableClassicCPic03 string `json:"lottery_turntable_classic_c_pic_03" example:"picture"`
	LotteryTurntableClassicCPic04 string `json:"lottery_turntable_classic_c_pic_04" example:"picture"`
	LotteryTurntableClassicCPic05 string `json:"lottery_turntable_classic_c_pic_05" example:"picture"`
	LotteryTurntableClassicCPic06 string `json:"lottery_turntable_classic_c_pic_06" example:"picture"`
	LotteryTurntableClassicCAni01 string `json:"lottery_turntable_classic_c_ani_01" example:"picture"`
	LotteryTurntableClassicCAni02 string `json:"lottery_turntable_classic_c_ani_02" example:"picture"`
	LotteryTurntableClassicCAni03 string `json:"lottery_turntable_classic_c_ani_03" example:"picture"`

	LotteryJiugonggeStarryskyHPic01 string `json:"lottery_jiugongge_starrysky_h_pic_01" example:"picture"`
	LotteryJiugonggeStarryskyHPic02 string `json:"lottery_jiugongge_starrysky_h_pic_02" example:"picture"`
	LotteryJiugonggeStarryskyHPic03 string `json:"lottery_jiugongge_starrysky_h_pic_03" example:"picture"`
	LotteryJiugonggeStarryskyHPic04 string `json:"lottery_jiugongge_starrysky_h_pic_04" example:"picture"`
	LotteryJiugonggeStarryskyHPic05 string `json:"lottery_jiugongge_starrysky_h_pic_05" example:"picture"`
	LotteryJiugonggeStarryskyHPic06 string `json:"lottery_jiugongge_starrysky_h_pic_06" example:"picture"`
	LotteryJiugonggeStarryskyHPic07 string `json:"lottery_jiugongge_starrysky_h_pic_07" example:"picture"`
	LotteryJiugonggeStarryskyGPic01 string `json:"lottery_jiugongge_starrysky_g_pic_01" example:"picture"`
	LotteryJiugonggeStarryskyGPic02 string `json:"lottery_jiugongge_starrysky_g_pic_02" example:"picture"`
	LotteryJiugonggeStarryskyGPic03 string `json:"lottery_jiugongge_starrysky_g_pic_03" example:"picture"`
	LotteryJiugonggeStarryskyGPic04 string `json:"lottery_jiugongge_starrysky_g_pic_04" example:"picture"`
	LotteryJiugonggeStarryskyCPic01 string `json:"lottery_jiugongge_starrysky_c_pic_01" example:"picture"`
	LotteryJiugonggeStarryskyCPic02 string `json:"lottery_jiugongge_starrysky_c_pic_02" example:"picture"`
	LotteryJiugonggeStarryskyCPic03 string `json:"lottery_jiugongge_starrysky_c_pic_03" example:"picture"`
	LotteryJiugonggeStarryskyCPic04 string `json:"lottery_jiugongge_starrysky_c_pic_04" example:"picture"`
	LotteryJiugonggeStarryskyCAni01 string `json:"lottery_jiugongge_starrysky_c_ani_01" example:"picture"`
	LotteryJiugonggeStarryskyCAni02 string `json:"lottery_jiugongge_starrysky_c_ani_02" example:"picture"`
	LotteryJiugonggeStarryskyCAni03 string `json:"lottery_jiugongge_starrysky_c_ani_03" example:"picture"`
	LotteryJiugonggeStarryskyCAni04 string `json:"lottery_jiugongge_starrysky_c_ani_04" example:"picture"`
	LotteryJiugonggeStarryskyCAni05 string `json:"lottery_jiugongge_starrysky_c_ani_05" example:"picture"`
	LotteryJiugonggeStarryskyCAni06 string `json:"lottery_jiugongge_starrysky_c_ani_06" example:"picture"`

	LotteryTurntableStarryskyHPic01 string `json:"lottery_turntable_starrysky_h_pic_01" example:"picture"`
	LotteryTurntableStarryskyHPic02 string `json:"lottery_turntable_starrysky_h_pic_02" example:"picture"`
	LotteryTurntableStarryskyHPic03 string `json:"lottery_turntable_starrysky_h_pic_03" example:"picture"`
	LotteryTurntableStarryskyHPic04 string `json:"lottery_turntable_starrysky_h_pic_04" example:"picture"`
	LotteryTurntableStarryskyHPic05 string `json:"lottery_turntable_starrysky_h_pic_05" example:"picture"`
	LotteryTurntableStarryskyHPic06 string `json:"lottery_turntable_starrysky_h_pic_06" example:"picture"`
	LotteryTurntableStarryskyHPic07 string `json:"lottery_turntable_starrysky_h_pic_07" example:"picture"`
	LotteryTurntableStarryskyHPic08 string `json:"lottery_turntable_starrysky_h_pic_08" example:"picture"`
	LotteryTurntableStarryskyGPic01 string `json:"lottery_turntable_starrysky_g_pic_01" example:"picture"`
	LotteryTurntableStarryskyGPic02 string `json:"lottery_turntable_starrysky_g_pic_02" example:"picture"`
	LotteryTurntableStarryskyGPic03 string `json:"lottery_turntable_starrysky_g_pic_03" example:"picture"`
	LotteryTurntableStarryskyGPic04 string `json:"lottery_turntable_starrysky_g_pic_04" example:"picture"`
	LotteryTurntableStarryskyGPic05 string `json:"lottery_turntable_starrysky_g_pic_05" example:"picture"`
	LotteryTurntableStarryskyCPic01 string `json:"lottery_turntable_starrysky_c_pic_01" example:"picture"`
	LotteryTurntableStarryskyCPic02 string `json:"lottery_turntable_starrysky_c_pic_02" example:"picture"`
	LotteryTurntableStarryskyCPic03 string `json:"lottery_turntable_starrysky_c_pic_03" example:"picture"`
	LotteryTurntableStarryskyCPic04 string `json:"lottery_turntable_starrysky_c_pic_04" example:"picture"`
	LotteryTurntableStarryskyCPic05 string `json:"lottery_turntable_starrysky_c_pic_05" example:"picture"`
	LotteryTurntableStarryskyCAni01 string `json:"lottery_turntable_starrysky_c_ani_01" example:"picture"`
	LotteryTurntableStarryskyCAni02 string `json:"lottery_turntable_starrysky_c_ani_02" example:"picture"`
	LotteryTurntableStarryskyCAni03 string `json:"lottery_turntable_starrysky_c_ani_03" example:"picture"`
	LotteryTurntableStarryskyCAni04 string `json:"lottery_turntable_starrysky_c_ani_04" example:"picture"`
	LotteryTurntableStarryskyCAni05 string `json:"lottery_turntable_starrysky_c_ani_05" example:"picture"`
	LotteryTurntableStarryskyCAni06 string `json:"lottery_turntable_starrysky_c_ani_06" example:"picture"`
	LotteryTurntableStarryskyCAni07 string `json:"lottery_turntable_starrysky_c_ani_07" example:"picture"`

	// 音樂
	LotteryBgmGaming string `json:"lottery_bgm_gaming" example:"picture"` // 遊戲進行中

	// 鑑定師自定義
	MonopolyClassicHPic01 string `json:"monopoly_classic_h_pic_01" example:"picture"`
	MonopolyClassicHPic02 string `json:"monopoly_classic_h_pic_02" example:"picture"`
	MonopolyClassicHPic03 string `json:"monopoly_classic_h_pic_03" example:"picture"`
	MonopolyClassicHPic04 string `json:"monopoly_classic_h_pic_04" example:"picture"`
	MonopolyClassicHPic05 string `json:"monopoly_classic_h_pic_05" example:"picture"`
	MonopolyClassicHPic06 string `json:"monopoly_classic_h_pic_06" example:"picture"`
	MonopolyClassicHPic07 string `json:"monopoly_classic_h_pic_07" example:"picture"`
	MonopolyClassicHPic08 string `json:"monopoly_classic_h_pic_08" example:"picture"`
	MonopolyClassicGPic01 string `json:"monopoly_classic_g_pic_01" example:"picture"`
	MonopolyClassicGPic02 string `json:"monopoly_classic_g_pic_02" example:"picture"`
	MonopolyClassicGPic03 string `json:"monopoly_classic_g_pic_03" example:"picture"`
	MonopolyClassicGPic04 string `json:"monopoly_classic_g_pic_04" example:"picture"`
	MonopolyClassicGPic05 string `json:"monopoly_classic_g_pic_05" example:"picture"`
	MonopolyClassicGPic06 string `json:"monopoly_classic_g_pic_06" example:"picture"`
	MonopolyClassicGPic07 string `json:"monopoly_classic_g_pic_07" example:"picture"`
	MonopolyClassicCPic01 string `json:"monopoly_classic_c_pic_01" example:"picture"`
	MonopolyClassicCPic02 string `json:"monopoly_classic_c_pic_02" example:"picture"`
	MonopolyClassicGAni01 string `json:"monopoly_classic_g_ani_01" example:"picture"`
	MonopolyClassicGAni02 string `json:"monopoly_classic_g_ani_02" example:"picture"`
	MonopolyClassicCAni01 string `json:"monopoly_classic_c_ani_01" example:"picture"`

	MonopolyRedpackHPic01 string `json:"monopoly_redpack_h_pic_01" example:"picture"`
	MonopolyRedpackHPic02 string `json:"monopoly_redpack_h_pic_02" example:"picture"`
	MonopolyRedpackHPic03 string `json:"monopoly_redpack_h_pic_03" example:"picture"`
	MonopolyRedpackHPic04 string `json:"monopoly_redpack_h_pic_04" example:"picture"`
	MonopolyRedpackHPic05 string `json:"monopoly_redpack_h_pic_05" example:"picture"`
	MonopolyRedpackHPic06 string `json:"monopoly_redpack_h_pic_06" example:"picture"`
	MonopolyRedpackHPic07 string `json:"monopoly_redpack_h_pic_07" example:"picture"`
	MonopolyRedpackHPic08 string `json:"monopoly_redpack_h_pic_08" example:"picture"`
	MonopolyRedpackHPic09 string `json:"monopoly_redpack_h_pic_09" example:"picture"`
	MonopolyRedpackHPic10 string `json:"monopoly_redpack_h_pic_10" example:"picture"`
	MonopolyRedpackHPic11 string `json:"monopoly_redpack_h_pic_11" example:"picture"`
	MonopolyRedpackHPic12 string `json:"monopoly_redpack_h_pic_12" example:"picture"`
	MonopolyRedpackHPic13 string `json:"monopoly_redpack_h_pic_13" example:"picture"`
	MonopolyRedpackHPic14 string `json:"monopoly_redpack_h_pic_14" example:"picture"`
	MonopolyRedpackHPic15 string `json:"monopoly_redpack_h_pic_15" example:"picture"`
	MonopolyRedpackHPic16 string `json:"monopoly_redpack_h_pic_16" example:"picture"`
	MonopolyRedpackGPic01 string `json:"monopoly_redpack_g_pic_01" example:"picture"`
	MonopolyRedpackGPic02 string `json:"monopoly_redpack_g_pic_02" example:"picture"`
	MonopolyRedpackGPic03 string `json:"monopoly_redpack_g_pic_03" example:"picture"`
	MonopolyRedpackGPic04 string `json:"monopoly_redpack_g_pic_04" example:"picture"`
	MonopolyRedpackGPic05 string `json:"monopoly_redpack_g_pic_05" example:"picture"`
	MonopolyRedpackGPic06 string `json:"monopoly_redpack_g_pic_06" example:"picture"`
	MonopolyRedpackGPic07 string `json:"monopoly_redpack_g_pic_07" example:"picture"`
	MonopolyRedpackGPic08 string `json:"monopoly_redpack_g_pic_08" example:"picture"`
	MonopolyRedpackGPic09 string `json:"monopoly_redpack_g_pic_09" example:"picture"`
	MonopolyRedpackGPic10 string `json:"monopoly_redpack_g_pic_10" example:"picture"`
	MonopolyRedpackCPic01 string `json:"monopoly_redpack_c_pic_01" example:"picture"`
	MonopolyRedpackCPic02 string `json:"monopoly_redpack_c_pic_02" example:"picture"`
	MonopolyRedpackCPic03 string `json:"monopoly_redpack_c_pic_03" example:"picture"`
	MonopolyRedpackHAni01 string `json:"monopoly_redpack_h_ani_01" example:"picture"`
	MonopolyRedpackHAni02 string `json:"monopoly_redpack_h_ani_02" example:"picture"`
	MonopolyRedpackHAni03 string `json:"monopoly_redpack_h_ani_03" example:"picture"`
	MonopolyRedpackGAni01 string `json:"monopoly_redpack_g_ani_01" example:"picture"`
	MonopolyRedpackGAni02 string `json:"monopoly_redpack_g_ani_02" example:"picture"`
	MonopolyRedpackCAni01 string `json:"monopoly_redpack_c_ani_01" example:"picture"`

	MonopolyNewyearRabbitHPic01 string `json:"monopoly_newyear_rabbit_h_pic_01" example:"picture"`
	MonopolyNewyearRabbitHPic02 string `json:"monopoly_newyear_rabbit_h_pic_02" example:"picture"`
	MonopolyNewyearRabbitHPic03 string `json:"monopoly_newyear_rabbit_h_pic_03" example:"picture"`
	MonopolyNewyearRabbitHPic04 string `json:"monopoly_newyear_rabbit_h_pic_04" example:"picture"`
	MonopolyNewyearRabbitHPic05 string `json:"monopoly_newyear_rabbit_h_pic_05" example:"picture"`
	MonopolyNewyearRabbitHPic06 string `json:"monopoly_newyear_rabbit_h_pic_06" example:"picture"`
	MonopolyNewyearRabbitHPic07 string `json:"monopoly_newyear_rabbit_h_pic_07" example:"picture"`
	MonopolyNewyearRabbitHPic08 string `json:"monopoly_newyear_rabbit_h_pic_08" example:"picture"`
	MonopolyNewyearRabbitHPic09 string `json:"monopoly_newyear_rabbit_h_pic_09" example:"picture"`
	MonopolyNewyearRabbitHPic10 string `json:"monopoly_newyear_rabbit_h_pic_10" example:"picture"`
	MonopolyNewyearRabbitHPic11 string `json:"monopoly_newyear_rabbit_h_pic_11" example:"picture"`
	MonopolyNewyearRabbitHPic12 string `json:"monopoly_newyear_rabbit_h_pic_12" example:"picture"`
	MonopolyNewyearRabbitGPic01 string `json:"monopoly_newyear_rabbit_g_pic_01" example:"picture"`
	MonopolyNewyearRabbitGPic02 string `json:"monopoly_newyear_rabbit_g_pic_02" example:"picture"`
	MonopolyNewyearRabbitGPic03 string `json:"monopoly_newyear_rabbit_g_pic_03" example:"picture"`
	MonopolyNewyearRabbitGPic04 string `json:"monopoly_newyear_rabbit_g_pic_04" example:"picture"`
	MonopolyNewyearRabbitGPic05 string `json:"monopoly_newyear_rabbit_g_pic_05" example:"picture"`
	MonopolyNewyearRabbitGPic06 string `json:"monopoly_newyear_rabbit_g_pic_06" example:"picture"`
	MonopolyNewyearRabbitGPic07 string `json:"monopoly_newyear_rabbit_g_pic_07" example:"picture"`
	MonopolyNewyearRabbitCPic01 string `json:"monopoly_newyear_rabbit_c_pic_01" example:"picture"`
	MonopolyNewyearRabbitCPic02 string `json:"monopoly_newyear_rabbit_c_pic_02" example:"picture"`
	MonopolyNewyearRabbitCPic03 string `json:"monopoly_newyear_rabbit_c_pic_03" example:"picture"`
	MonopolyNewyearRabbitHAni01 string `json:"monopoly_newyear_rabbit_h_ani_01" example:"picture"`
	MonopolyNewyearRabbitHAni02 string `json:"monopoly_newyear_rabbit_h_ani_02" example:"picture"`
	MonopolyNewyearRabbitGAni01 string `json:"monopoly_newyear_rabbit_g_ani_01" example:"picture"`
	MonopolyNewyearRabbitGAni02 string `json:"monopoly_newyear_rabbit_g_ani_02" example:"picture"`
	MonopolyNewyearRabbitCAni01 string `json:"monopoly_newyear_rabbit_c_ani_01" example:"picture"`

	MonopolySashimiHPic01 string `json:"monopoly_sashimi_h_pic_01" example:"picture"`
	MonopolySashimiHPic02 string `json:"monopoly_sashimi_h_pic_02" example:"picture"`
	MonopolySashimiHPic03 string `json:"monopoly_sashimi_h_pic_03" example:"picture"`
	MonopolySashimiHPic04 string `json:"monopoly_sashimi_h_pic_04" example:"picture"`
	MonopolySashimiHPic05 string `json:"monopoly_sashimi_h_pic_05" example:"picture"`
	MonopolySashimiGPic01 string `json:"monopoly_sashimi_g_pic_01" example:"picture"`
	MonopolySashimiGPic02 string `json:"monopoly_sashimi_g_pic_02" example:"picture"`
	MonopolySashimiGPic03 string `json:"monopoly_sashimi_g_pic_03" example:"picture"`
	MonopolySashimiGPic04 string `json:"monopoly_sashimi_g_pic_04" example:"picture"`
	MonopolySashimiGPic05 string `json:"monopoly_sashimi_g_pic_05" example:"picture"`
	MonopolySashimiGPic06 string `json:"monopoly_sashimi_g_pic_06" example:"picture"`
	MonopolySashimiCPic01 string `json:"monopoly_sashimi_c_pic_01" example:"picture"`
	MonopolySashimiCPic02 string `json:"monopoly_sashimi_c_pic_02" example:"picture"`
	MonopolySashimiHAni01 string `json:"monopoly_sashimi_h_ani_01" example:"picture"`
	MonopolySashimiHAni02 string `json:"monopoly_sashimi_h_ani_02" example:"picture"`
	MonopolySashimiGAni01 string `json:"monopoly_sashimi_g_ani_01" example:"picture"`
	MonopolySashimiGAni02 string `json:"monopoly_sashimi_g_ani_02" example:"picture"`
	MonopolySashimiCAni01 string `json:"monopoly_sashimi_c_ani_01" example:"picture"`

	// 音樂
	MonopolyBgmStart  string `json:"monopoly_bgm_start" example:"picture"`  // 遊戲開始
	MonopolyBgmGaming string `json:"monopoly_bgm_gaming" example:"picture"` // 遊戲進行中
	MonopolyBgmEnd    string `json:"monopoly_bgm_end" example:"picture"`    // 遊戲結束

	// 拔河遊戲自定義
	TugofwarClassicHPic01 string `json:"tugofwar_classic_h_pic_01" example:"picture"`
	TugofwarClassicHPic02 string `json:"tugofwar_classic_h_pic_02" example:"picture"`
	TugofwarClassicHPic03 string `json:"tugofwar_classic_h_pic_03" example:"picture"`
	TugofwarClassicHPic04 string `json:"tugofwar_classic_h_pic_04" example:"picture"`
	TugofwarClassicHPic05 string `json:"tugofwar_classic_h_pic_05" example:"picture"`
	TugofwarClassicHPic06 string `json:"tugofwar_classic_h_pic_06" example:"picture"`
	TugofwarClassicHPic07 string `json:"tugofwar_classic_h_pic_07" example:"picture"`
	TugofwarClassicHPic08 string `json:"tugofwar_classic_h_pic_08" example:"picture"`
	TugofwarClassicHPic09 string `json:"tugofwar_classic_h_pic_09" example:"picture"`
	TugofwarClassicHPic10 string `json:"tugofwar_classic_h_pic_10" example:"picture"`
	TugofwarClassicHPic11 string `json:"tugofwar_classic_h_pic_11" example:"picture"`
	TugofwarClassicHPic12 string `json:"tugofwar_classic_h_pic_12" example:"picture"`
	TugofwarClassicHPic13 string `json:"tugofwar_classic_h_pic_13" example:"picture"`
	TugofwarClassicHPic14 string `json:"tugofwar_classic_h_pic_14" example:"picture"`
	TugofwarClassicHPic15 string `json:"tugofwar_classic_h_pic_15" example:"picture"`
	TugofwarClassicHPic16 string `json:"tugofwar_classic_h_pic_16" example:"picture"`
	TugofwarClassicHPic17 string `json:"tugofwar_classic_h_pic_17" example:"picture"`
	TugofwarClassicHPic18 string `json:"tugofwar_classic_h_pic_18" example:"picture"`
	TugofwarClassicHPic19 string `json:"tugofwar_classic_h_pic_19" example:"picture"`
	TugofwarClassicGPic01 string `json:"tugofwar_classic_g_pic_01" example:"picture"`
	TugofwarClassicGPic02 string `json:"tugofwar_classic_g_pic_02" example:"picture"`
	TugofwarClassicGPic03 string `json:"tugofwar_classic_g_pic_03" example:"picture"`
	TugofwarClassicGPic04 string `json:"tugofwar_classic_g_pic_04" example:"picture"`
	TugofwarClassicGPic05 string `json:"tugofwar_classic_g_pic_05" example:"picture"`
	TugofwarClassicGPic06 string `json:"tugofwar_classic_g_pic_06" example:"picture"`
	TugofwarClassicGPic07 string `json:"tugofwar_classic_g_pic_07" example:"picture"`
	TugofwarClassicGPic08 string `json:"tugofwar_classic_g_pic_08" example:"picture"`
	TugofwarClassicGPic09 string `json:"tugofwar_classic_g_pic_09" example:"picture"`
	TugofwarClassicHAni01 string `json:"tugofwar_classic_h_ani_01" example:"picture"`
	TugofwarClassicHAni02 string `json:"tugofwar_classic_h_ani_02" example:"picture"`
	TugofwarClassicHAni03 string `json:"tugofwar_classic_h_ani_03" example:"picture"`
	TugofwarClassicCAni01 string `json:"tugofwar_classic_c_ani_01" example:"picture"`

	TugofwarSchoolHPic01 string `json:"tugofwar_school_h_pic_01" example:"picture"`
	TugofwarSchoolHPic02 string `json:"tugofwar_school_h_pic_02" example:"picture"`
	TugofwarSchoolHPic03 string `json:"tugofwar_school_h_pic_03" example:"picture"`
	TugofwarSchoolHPic04 string `json:"tugofwar_school_h_pic_04" example:"picture"`
	TugofwarSchoolHPic05 string `json:"tugofwar_school_h_pic_05" example:"picture"`
	TugofwarSchoolHPic06 string `json:"tugofwar_school_h_pic_06" example:"picture"`
	TugofwarSchoolHPic07 string `json:"tugofwar_school_h_pic_07" example:"picture"`
	TugofwarSchoolHPic08 string `json:"tugofwar_school_h_pic_08" example:"picture"`
	TugofwarSchoolHPic09 string `json:"tugofwar_school_h_pic_09" example:"picture"`
	TugofwarSchoolHPic10 string `json:"tugofwar_school_h_pic_10" example:"picture"`
	TugofwarSchoolHPic11 string `json:"tugofwar_school_h_pic_11" example:"picture"`
	TugofwarSchoolHPic12 string `json:"tugofwar_school_h_pic_12" example:"picture"`
	TugofwarSchoolHPic13 string `json:"tugofwar_school_h_pic_13" example:"picture"`
	TugofwarSchoolHPic14 string `json:"tugofwar_school_h_pic_14" example:"picture"`
	TugofwarSchoolHPic15 string `json:"tugofwar_school_h_pic_15" example:"picture"`
	TugofwarSchoolHPic16 string `json:"tugofwar_school_h_pic_16" example:"picture"`
	TugofwarSchoolHPic17 string `json:"tugofwar_school_h_pic_17" example:"picture"`
	TugofwarSchoolHPic18 string `json:"tugofwar_school_h_pic_18" example:"picture"`
	TugofwarSchoolHPic19 string `json:"tugofwar_school_h_pic_19" example:"picture"`
	TugofwarSchoolHPic20 string `json:"tugofwar_school_h_pic_20" example:"picture"`
	TugofwarSchoolHPic21 string `json:"tugofwar_school_h_pic_21" example:"picture"`
	TugofwarSchoolHPic22 string `json:"tugofwar_school_h_pic_22" example:"picture"`
	TugofwarSchoolHPic23 string `json:"tugofwar_school_h_pic_23" example:"picture"`
	TugofwarSchoolHPic24 string `json:"tugofwar_school_h_pic_24" example:"picture"`
	TugofwarSchoolHPic25 string `json:"tugofwar_school_h_pic_25" example:"picture"`
	TugofwarSchoolHPic26 string `json:"tugofwar_school_h_pic_26" example:"picture"`
	TugofwarSchoolGPic01 string `json:"tugofwar_school_g_pic_01" example:"picture"`
	TugofwarSchoolGPic02 string `json:"tugofwar_school_g_pic_02" example:"picture"`
	TugofwarSchoolGPic03 string `json:"tugofwar_school_g_pic_03" example:"picture"`
	TugofwarSchoolGPic04 string `json:"tugofwar_school_g_pic_04" example:"picture"`
	TugofwarSchoolGPic05 string `json:"tugofwar_school_g_pic_05" example:"picture"`
	TugofwarSchoolGPic06 string `json:"tugofwar_school_g_pic_06" example:"picture"`
	TugofwarSchoolGPic07 string `json:"tugofwar_school_g_pic_07" example:"picture"`
	TugofwarSchoolCPic01 string `json:"tugofwar_school_c_pic_01" example:"picture"`
	TugofwarSchoolCPic02 string `json:"tugofwar_school_c_pic_02" example:"picture"`
	TugofwarSchoolCPic03 string `json:"tugofwar_school_c_pic_03" example:"picture"`
	TugofwarSchoolCPic04 string `json:"tugofwar_school_c_pic_04" example:"picture"`
	TugofwarSchoolHAni01 string `json:"tugofwar_school_h_ani_01" example:"picture"`
	TugofwarSchoolHAni02 string `json:"tugofwar_school_h_ani_02" example:"picture"`
	TugofwarSchoolHAni03 string `json:"tugofwar_school_h_ani_03" example:"picture"`
	TugofwarSchoolHAni04 string `json:"tugofwar_school_h_ani_04" example:"picture"`
	TugofwarSchoolHAni05 string `json:"tugofwar_school_h_ani_05" example:"picture"`
	TugofwarSchoolHAni06 string `json:"tugofwar_school_h_ani_06" example:"picture"`
	TugofwarSchoolHAni07 string `json:"tugofwar_school_h_ani_07" example:"picture"`

	TugofwarChristmasHPic01 string `json:"tugofwar_christmas_h_pic_01" example:"picture"`
	TugofwarChristmasHPic02 string `json:"tugofwar_christmas_h_pic_02" example:"picture"`
	TugofwarChristmasHPic03 string `json:"tugofwar_christmas_h_pic_03" example:"picture"`
	TugofwarChristmasHPic04 string `json:"tugofwar_christmas_h_pic_04" example:"picture"`
	TugofwarChristmasHPic05 string `json:"tugofwar_christmas_h_pic_05" example:"picture"`
	TugofwarChristmasHPic06 string `json:"tugofwar_christmas_h_pic_06" example:"picture"`
	TugofwarChristmasHPic07 string `json:"tugofwar_christmas_h_pic_07" example:"picture"`
	TugofwarChristmasHPic08 string `json:"tugofwar_christmas_h_pic_08" example:"picture"`
	TugofwarChristmasHPic09 string `json:"tugofwar_christmas_h_pic_09" example:"picture"`
	TugofwarChristmasHPic10 string `json:"tugofwar_christmas_h_pic_10" example:"picture"`
	TugofwarChristmasHPic11 string `json:"tugofwar_christmas_h_pic_11" example:"picture"`
	TugofwarChristmasHPic12 string `json:"tugofwar_christmas_h_pic_12" example:"picture"`
	TugofwarChristmasHPic13 string `json:"tugofwar_christmas_h_pic_13" example:"picture"`
	TugofwarChristmasHPic14 string `json:"tugofwar_christmas_h_pic_14" example:"picture"`
	TugofwarChristmasHPic15 string `json:"tugofwar_christmas_h_pic_15" example:"picture"`
	TugofwarChristmasHPic16 string `json:"tugofwar_christmas_h_pic_16" example:"picture"`
	TugofwarChristmasHPic17 string `json:"tugofwar_christmas_h_pic_17" example:"picture"`
	TugofwarChristmasHPic18 string `json:"tugofwar_christmas_h_pic_18" example:"picture"`
	TugofwarChristmasHPic19 string `json:"tugofwar_christmas_h_pic_19" example:"picture"`
	TugofwarChristmasHPic20 string `json:"tugofwar_christmas_h_pic_20" example:"picture"`
	TugofwarChristmasHPic21 string `json:"tugofwar_christmas_h_pic_21" example:"picture"`
	TugofwarChristmasGPic01 string `json:"tugofwar_christmas_g_pic_01" example:"picture"`
	TugofwarChristmasGPic02 string `json:"tugofwar_christmas_g_pic_02" example:"picture"`
	TugofwarChristmasGPic03 string `json:"tugofwar_christmas_g_pic_03" example:"picture"`
	TugofwarChristmasGPic04 string `json:"tugofwar_christmas_g_pic_04" example:"picture"`
	TugofwarChristmasGPic05 string `json:"tugofwar_christmas_g_pic_05" example:"picture"`
	TugofwarChristmasGPic06 string `json:"tugofwar_christmas_g_pic_06" example:"picture"`
	TugofwarChristmasCPic01 string `json:"tugofwar_christmas_c_pic_01" example:"picture"`
	TugofwarChristmasCPic02 string `json:"tugofwar_christmas_c_pic_02" example:"picture"`
	TugofwarChristmasCPic03 string `json:"tugofwar_christmas_c_pic_03" example:"picture"`
	TugofwarChristmasCPic04 string `json:"tugofwar_christmas_c_pic_04" example:"picture"`
	TugofwarChristmasHAni01 string `json:"tugofwar_christmas_h_ani_01" example:"picture"`
	TugofwarChristmasHAni02 string `json:"tugofwar_christmas_h_ani_02" example:"picture"`
	TugofwarChristmasHAni03 string `json:"tugofwar_christmas_h_ani_03" example:"picture"`
	TugofwarChristmasCAni01 string `json:"tugofwar_christmas_c_ani_01" example:"picture"`
	TugofwarChristmasCAni02 string `json:"tugofwar_christmas_c_ani_02" example:"picture"`

	// 音樂
	TugofwarBgmStart  string `json:"tugofwar_bgm_start" example:"picture"`  // 遊戲開始
	TugofwarBgmGaming string `json:"tugofwar_bgm_gaming" example:"picture"` // 遊戲進行中
	TugofwarBgmEnd    string `json:"tugofwar_bgm_end" example:"picture"`    // 遊戲結束

	// 賓果遊戲自定義
	BingoClassicHPic01 string `json:"bingo_classic_h_pic_01" example:"picture"`
	BingoClassicHPic02 string `json:"bingo_classic_h_pic_02" example:"picture"`
	BingoClassicHPic03 string `json:"bingo_classic_h_pic_03" example:"picture"`
	BingoClassicHPic04 string `json:"bingo_classic_h_pic_04" example:"picture"`
	BingoClassicHPic05 string `json:"bingo_classic_h_pic_05" example:"picture"`
	BingoClassicHPic06 string `json:"bingo_classic_h_pic_06" example:"picture"`
	BingoClassicHPic07 string `json:"bingo_classic_h_pic_07" example:"picture"`
	BingoClassicHPic08 string `json:"bingo_classic_h_pic_08" example:"picture"`
	BingoClassicHPic09 string `json:"bingo_classic_h_pic_09" example:"picture"`
	BingoClassicHPic10 string `json:"bingo_classic_h_pic_10" example:"picture"`
	BingoClassicHPic11 string `json:"bingo_classic_h_pic_11" example:"picture"`
	BingoClassicHPic12 string `json:"bingo_classic_h_pic_12" example:"picture"`
	BingoClassicHPic13 string `json:"bingo_classic_h_pic_13" example:"picture"`
	BingoClassicHPic14 string `json:"bingo_classic_h_pic_14" example:"picture"`
	BingoClassicHPic15 string `json:"bingo_classic_h_pic_15" example:"picture"`
	BingoClassicHPic16 string `json:"bingo_classic_h_pic_16" example:"picture"`
	BingoClassicGPic01 string `json:"bingo_classic_g_pic_01" example:"picture"`
	BingoClassicGPic02 string `json:"bingo_classic_g_pic_02" example:"picture"`
	BingoClassicGPic03 string `json:"bingo_classic_g_pic_03" example:"picture"`
	BingoClassicGPic04 string `json:"bingo_classic_g_pic_04" example:"picture"`
	BingoClassicGPic05 string `json:"bingo_classic_g_pic_05" example:"picture"`
	BingoClassicGPic06 string `json:"bingo_classic_g_pic_06" example:"picture"`
	BingoClassicGPic07 string `json:"bingo_classic_g_pic_07" example:"picture"`
	BingoClassicGPic08 string `json:"bingo_classic_g_pic_08" example:"picture"`
	BingoClassicCPic01 string `json:"bingo_classic_c_pic_01" example:"picture"`
	BingoClassicCPic02 string `json:"bingo_classic_c_pic_02" example:"picture"`
	BingoClassicCPic03 string `json:"bingo_classic_c_pic_03" example:"picture"`
	BingoClassicCPic04 string `json:"bingo_classic_c_pic_04" example:"picture"`
	// BingoClassicCPic05 string `json:"bingo_classic_c_pic_05" example:"picture"`
	BingoClassicHAni01 string `json:"bingo_classic_h_ani_01" example:"picture"`
	BingoClassicGAni01 string `json:"bingo_classic_g_ani_01" example:"picture"`
	BingoClassicCAni01 string `json:"bingo_classic_c_ani_01" example:"picture"`
	BingoClassicCAni02 string `json:"bingo_classic_c_ani_02" example:"picture"`

	BingoNewyearDragonHPic01 string `json:"bingo_newyear_dragon_h_pic_01" example:"picture"`
	BingoNewyearDragonHPic02 string `json:"bingo_newyear_dragon_h_pic_02" example:"picture"`
	BingoNewyearDragonHPic03 string `json:"bingo_newyear_dragon_h_pic_03" example:"picture"`
	BingoNewyearDragonHPic04 string `json:"bingo_newyear_dragon_h_pic_04" example:"picture"`
	BingoNewyearDragonHPic05 string `json:"bingo_newyear_dragon_h_pic_05" example:"picture"`
	BingoNewyearDragonHPic06 string `json:"bingo_newyear_dragon_h_pic_06" example:"picture"`
	BingoNewyearDragonHPic07 string `json:"bingo_newyear_dragon_h_pic_07" example:"picture"`
	BingoNewyearDragonHPic08 string `json:"bingo_newyear_dragon_h_pic_08" example:"picture"`
	BingoNewyearDragonHPic09 string `json:"bingo_newyear_dragon_h_pic_09" example:"picture"`
	BingoNewyearDragonHPic10 string `json:"bingo_newyear_dragon_h_pic_10" example:"picture"`
	BingoNewyearDragonHPic11 string `json:"bingo_newyear_dragon_h_pic_11" example:"picture"`
	BingoNewyearDragonHPic12 string `json:"bingo_newyear_dragon_h_pic_12" example:"picture"`
	BingoNewyearDragonHPic13 string `json:"bingo_newyear_dragon_h_pic_13" example:"picture"`
	BingoNewyearDragonHPic14 string `json:"bingo_newyear_dragon_h_pic_14" example:"picture"`
	// BingoNewyearDragonHPic15 string `json:"bingo_newyear_dragon_h_pic_15" example:"picture"`
	BingoNewyearDragonHPic16 string `json:"bingo_newyear_dragon_h_pic_16" example:"picture"`
	BingoNewyearDragonHPic17 string `json:"bingo_newyear_dragon_h_pic_17" example:"picture"`
	BingoNewyearDragonHPic18 string `json:"bingo_newyear_dragon_h_pic_18" example:"picture"`
	BingoNewyearDragonHPic19 string `json:"bingo_newyear_dragon_h_pic_19" example:"picture"`
	BingoNewyearDragonHPic20 string `json:"bingo_newyear_dragon_h_pic_20" example:"picture"`
	BingoNewyearDragonHPic21 string `json:"bingo_newyear_dragon_h_pic_21" example:"picture"`
	BingoNewyearDragonHPic22 string `json:"bingo_newyear_dragon_h_pic_22" example:"picture"`
	BingoNewyearDragonGPic01 string `json:"bingo_newyear_dragon_g_pic_01" example:"picture"`
	BingoNewyearDragonGPic02 string `json:"bingo_newyear_dragon_g_pic_02" example:"picture"`
	BingoNewyearDragonGPic03 string `json:"bingo_newyear_dragon_g_pic_03" example:"picture"`
	BingoNewyearDragonGPic04 string `json:"bingo_newyear_dragon_g_pic_04" example:"picture"`
	BingoNewyearDragonGPic05 string `json:"bingo_newyear_dragon_g_pic_05" example:"picture"`
	BingoNewyearDragonGPic06 string `json:"bingo_newyear_dragon_g_pic_06" example:"picture"`
	BingoNewyearDragonGPic07 string `json:"bingo_newyear_dragon_g_pic_07" example:"picture"`
	BingoNewyearDragonGPic08 string `json:"bingo_newyear_dragon_g_pic_08" example:"picture"`
	BingoNewyearDragonCPic01 string `json:"bingo_newyear_dragon_c_pic_01" example:"picture"`
	BingoNewyearDragonCPic02 string `json:"bingo_newyear_dragon_c_pic_02" example:"picture"`
	BingoNewyearDragonCPic03 string `json:"bingo_newyear_dragon_c_pic_03" example:"picture"`
	BingoNewyearDragonHAni01 string `json:"bingo_newyear_dragon_h_ani_01" example:"picture"`
	BingoNewyearDragonGAni01 string `json:"bingo_newyear_dragon_g_ani_01" example:"picture"`
	BingoNewyearDragonCAni01 string `json:"bingo_newyear_dragon_c_ani_01" example:"picture"`
	BingoNewyearDragonCAni02 string `json:"bingo_newyear_dragon_c_ani_02" example:"picture"`
	BingoNewyearDragonCAni03 string `json:"bingo_newyear_dragon_c_ani_03" example:"picture"`

	BingoCherryHPic01 string `json:"bingo_cherry_h_pic_01" example:"picture"`
	BingoCherryHPic02 string `json:"bingo_cherry_h_pic_02" example:"picture"`
	BingoCherryHPic03 string `json:"bingo_cherry_h_pic_03" example:"picture"`
	BingoCherryHPic04 string `json:"bingo_cherry_h_pic_04" example:"picture"`
	BingoCherryHPic05 string `json:"bingo_cherry_h_pic_05" example:"picture"`
	BingoCherryHPic06 string `json:"bingo_cherry_h_pic_06" example:"picture"`
	BingoCherryHPic07 string `json:"bingo_cherry_h_pic_07" example:"picture"`
	BingoCherryHPic08 string `json:"bingo_cherry_h_pic_08" example:"picture"`
	BingoCherryHPic09 string `json:"bingo_cherry_h_pic_09" example:"picture"`
	BingoCherryHPic10 string `json:"bingo_cherry_h_pic_10" example:"picture"`
	BingoCherryHPic11 string `json:"bingo_cherry_h_pic_11" example:"picture"`
	BingoCherryHPic12 string `json:"bingo_cherry_h_pic_12" example:"picture"`
	// BingoCherryHPic13 string `json:"bingo_cherry_h_pic_13" example:"picture"`
	BingoCherryHPic14 string `json:"bingo_cherry_h_pic_14" example:"picture"`
	BingoCherryHPic15 string `json:"bingo_cherry_h_pic_15" example:"picture"`
	// BingoCherryHPic16 string `json:"bingo_cherry_h_pic_16" example:"picture"`
	BingoCherryHPic17 string `json:"bingo_cherry_h_pic_17" example:"picture"`
	BingoCherryHPic18 string `json:"bingo_cherry_h_pic_18" example:"picture"`
	BingoCherryHPic19 string `json:"bingo_cherry_h_pic_19" example:"picture"`
	BingoCherryGPic01 string `json:"bingo_cherry_g_pic_01" example:"picture"`
	BingoCherryGPic02 string `json:"bingo_cherry_g_pic_02" example:"picture"`
	BingoCherryGPic03 string `json:"bingo_cherry_g_pic_03" example:"picture"`
	BingoCherryGPic04 string `json:"bingo_cherry_g_pic_04" example:"picture"`
	BingoCherryGPic05 string `json:"bingo_cherry_g_pic_05" example:"picture"`
	BingoCherryGPic06 string `json:"bingo_cherry_g_pic_06" example:"picture"`
	BingoCherryCPic01 string `json:"bingo_cherry_c_pic_01" example:"picture"`
	BingoCherryCPic02 string `json:"bingo_cherry_c_pic_02" example:"picture"`
	BingoCherryCPic03 string `json:"bingo_cherry_c_pic_03" example:"picture"`
	BingoCherryCPic04 string `json:"bingo_cherry_c_pic_04" example:"picture"`
	// BingoCherryHAni01 string `json:"bingo_cherry_h_ani_01" example:"picture"`
	BingoCherryHAni02 string `json:"bingo_cherry_h_ani_02" example:"picture"`
	BingoCherryHAni03 string `json:"bingo_cherry_h_ani_03" example:"picture"`
	BingoCherryGAni01 string `json:"bingo_cherry_g_ani_01" example:"picture"`
	BingoCherryGAni02 string `json:"bingo_cherry_g_ani_02" example:"picture"`
	BingoCherryCAni01 string `json:"bingo_cherry_c_ani_01" example:"picture"`
	BingoCherryCAni02 string `json:"bingo_cherry_c_ani_02" example:"picture"`

	// 音樂
	BingoBgmStart  string `json:"bingo_bgm_start" example:"picture"`  // 遊戲開始
	BingoBgmGaming string `json:"bingo_bgm_gaming" example:"picture"` // 遊戲進行中
	BingoBgmEnd    string `json:"bingo_bgm_end" example:"picture"`    // 遊戲結束

	// 扭蛋機自定義
	GachaMachineClassicHPic02 string `json:"3d_gacha_machine_classic_h_pic_02" example:"picture"`
	GachaMachineClassicHPic03 string `json:"3d_gacha_machine_classic_h_pic_03" example:"picture"`
	GachaMachineClassicHPic04 string `json:"3d_gacha_machine_classic_h_pic_04" example:"picture"`
	GachaMachineClassicHPic05 string `json:"3d_gacha_machine_classic_h_pic_05" example:"picture"`
	GachaMachineClassicGPic01 string `json:"3d_gacha_machine_classic_g_pic_01" example:"picture"`
	GachaMachineClassicGPic02 string `json:"3d_gacha_machine_classic_g_pic_02" example:"picture"`
	GachaMachineClassicGPic03 string `json:"3d_gacha_machine_classic_g_pic_03" example:"picture"`
	GachaMachineClassicGPic04 string `json:"3d_gacha_machine_classic_g_pic_04" example:"picture"`
	GachaMachineClassicGPic05 string `json:"3d_gacha_machine_classic_g_pic_05" example:"picture"`
	GachaMachineClassicCPic01 string `json:"3d_gacha_machine_classic_c_pic_01" example:"picture"`

	// 音樂
	GachaMachineBgmGaming string `json:"3d_gacha_machine_bgm_gaming" example:"picture"`

	// 投票自定義
	VoteClassicHPic01 string `json:"vote_classic_h_pic_01" example:"picture"`
	VoteClassicHPic02 string `json:"vote_classic_h_pic_02" example:"picture"`
	VoteClassicHPic03 string `json:"vote_classic_h_pic_03" example:"picture"`
	VoteClassicHPic04 string `json:"vote_classic_h_pic_04" example:"picture"`
	VoteClassicHPic05 string `json:"vote_classic_h_pic_05" example:"picture"`
	VoteClassicHPic06 string `json:"vote_classic_h_pic_06" example:"picture"`
	VoteClassicHPic07 string `json:"vote_classic_h_pic_07" example:"picture"`
	VoteClassicHPic08 string `json:"vote_classic_h_pic_08" example:"picture"`
	VoteClassicHPic09 string `json:"vote_classic_h_pic_09" example:"picture"`
	VoteClassicHPic10 string `json:"vote_classic_h_pic_10" example:"picture"`
	VoteClassicHPic11 string `json:"vote_classic_h_pic_11" example:"picture"`
	// VoteClassicHPic12 string `json:"vote_classic_h_pic_12" example:"picture"`
	VoteClassicHPic13 string `json:"vote_classic_h_pic_13" example:"picture"`
	VoteClassicHPic14 string `json:"vote_classic_h_pic_14" example:"picture"`
	VoteClassicHPic15 string `json:"vote_classic_h_pic_15" example:"picture"`
	VoteClassicHPic16 string `json:"vote_classic_h_pic_16" example:"picture"`
	VoteClassicHPic17 string `json:"vote_classic_h_pic_17" example:"picture"`
	VoteClassicHPic18 string `json:"vote_classic_h_pic_18" example:"picture"`
	VoteClassicHPic19 string `json:"vote_classic_h_pic_19" example:"picture"`
	VoteClassicHPic20 string `json:"vote_classic_h_pic_20" example:"picture"`
	VoteClassicHPic21 string `json:"vote_classic_h_pic_21" example:"picture"`
	// VoteClassicHPic22 string `json:"vote_classic_h_pic_22" example:"picture"`
	VoteClassicHPic23 string `json:"vote_classic_h_pic_23" example:"picture"`
	VoteClassicHPic24 string `json:"vote_classic_h_pic_24" example:"picture"`
	VoteClassicHPic25 string `json:"vote_classic_h_pic_25" example:"picture"`
	VoteClassicHPic26 string `json:"vote_classic_h_pic_26" example:"picture"`
	VoteClassicHPic27 string `json:"vote_classic_h_pic_27" example:"picture"`
	VoteClassicHPic28 string `json:"vote_classic_h_pic_28" example:"picture"`
	VoteClassicHPic29 string `json:"vote_classic_h_pic_29" example:"picture"`
	VoteClassicHPic30 string `json:"vote_classic_h_pic_30" example:"picture"`
	VoteClassicHPic31 string `json:"vote_classic_h_pic_31" example:"picture"`
	VoteClassicHPic32 string `json:"vote_classic_h_pic_32" example:"picture"`
	VoteClassicHPic33 string `json:"vote_classic_h_pic_33" example:"picture"`
	VoteClassicHPic34 string `json:"vote_classic_h_pic_34" example:"picture"`
	VoteClassicHPic35 string `json:"vote_classic_h_pic_35" example:"picture"`
	VoteClassicHPic36 string `json:"vote_classic_h_pic_36" example:"picture"`
	VoteClassicHPic37 string `json:"vote_classic_h_pic_37" example:"picture"`
	VoteClassicGPic01 string `json:"vote_classic_g_pic_01" example:"picture"`
	VoteClassicGPic02 string `json:"vote_classic_g_pic_02" example:"picture"`
	VoteClassicGPic03 string `json:"vote_classic_g_pic_03" example:"picture"`
	VoteClassicGPic04 string `json:"vote_classic_g_pic_04" example:"picture"`
	VoteClassicGPic05 string `json:"vote_classic_g_pic_05" example:"picture"`
	VoteClassicGPic06 string `json:"vote_classic_g_pic_06" example:"picture"`
	VoteClassicGPic07 string `json:"vote_classic_g_pic_07" example:"picture"`
	VoteClassicCPic01 string `json:"vote_classic_c_pic_01" example:"picture"`
	VoteClassicCPic02 string `json:"vote_classic_c_pic_02" example:"picture"`
	VoteClassicCPic03 string `json:"vote_classic_c_pic_03" example:"picture"`
	VoteClassicCPic04 string `json:"vote_classic_c_pic_04" example:"picture"`
	// 音樂
	VoteBgmGaming string `json:"vote_bgm_gaming" example:"picture"`

	// 自定義圖片陣列參數
	CustomizeHostPictures      []string `json:"customize_host_pictures" example:"pictures"`      // 主持靜態圖片
	CustomizeGuestPictures     []string `json:"customize_guest_pictures" example:"pictures"`     // 玩家靜態圖片
	CustomizeCommonPictures    []string `json:"customize_common_pictures" example:"pictures"`    // 共用靜態圖片
	CustomizeHostAnipictures   []string `json:"customize_host_anipictures" example:"pictures"`   // 主持動態圖片
	CustomizeGuestAnipictures  []string `json:"customize_guest_anipictures" example:"pictures"`  // 玩家動態圖片
	CustomizeCommonAnipictures []string `json:"customize_common_anipictures" example:"pictures"` // 共用動態圖片
	CustomizeMusics            []string `json:"customize_musics" example:"pictures"`             // 音樂

	// 活動資訊(join activity)
	// UserID string `json:"user_id"`
	// Device string `json:"device"`

	// 用戶
	// MaxActivityPeople int64 `json:"max_activity_people"`
	// MaxGamePeople     int64 `json:"max_game_people"`

	// 快問快答
	Questions   []QuestionModel `json:"questions"`  // 所有題目資訊
	TotalQA     int64           `json:"total_qa"`   // 總題數
	QASecond    int64           `json:"qa_second"`  // 題目顯示秒數
	QARound     int64           `json:"qa_round"`   // 題目進行題數
	QAPeople    int64           `json:"qa_people" ` // 快問快答人數
	QA1         string          `json:"qa_1"`
	QA1Options  []string        `json:"qa_1_options"`
	QA1Answer   string          `json:"qa_1_answer"`
	QA1Score    int64           `json:"qa_1_score"`
	QA2         string          `json:"qa_2"`
	QA2Options  []string        `json:"qa_2_options"`
	QA2Answer   string          `json:"qa_2_answer"`
	QA2Score    int64           `json:"qa_2_score"`
	QA3         string          `json:"qa_3"`
	QA3Options  []string        `json:"qa_3_options"`
	QA3Answer   string          `json:"qa_3_answer"`
	QA3Score    int64           `json:"qa_3_score"`
	QA4         string          `json:"qa_4"`
	QA4Options  []string        `json:"qa_4_options"`
	QA4Answer   string          `json:"qa_4_answer"`
	QA4Score    int64           `json:"qa_4_score"`
	QA5         string          `json:"qa_5"`
	QA5Options  []string        `json:"qa_5_options"`
	QA5Answer   string          `json:"qa_5_answer"`
	QA5Score    int64           `json:"qa_5_score"`
	QA6         string          `json:"qa_6"`
	QA6Options  []string        `json:"qa_6_options"`
	QA6Answer   string          `json:"qa_6_answer"`
	QA6Score    int64           `json:"qa_6_score"`
	QA7         string          `json:"qa_7"`
	QA7Options  []string        `json:"qa_7_options"`
	QA7Answer   string          `json:"qa_7_answer"`
	QA7Score    int64           `json:"qa_7_score"`
	QA8         string          `json:"qa_8"`
	QA8Options  []string        `json:"qa_8_options"`
	QA8Answer   string          `json:"qa_8_answer"`
	QA8Score    int64           `json:"qa_8_score"`
	QA9         string          `json:"qa_9"`
	QA9Options  []string        `json:"qa_9_options"`
	QA9Answer   string          `json:"qa_9_answer"`
	QA9Score    int64           `json:"qa_9_score"`
	QA10        string          `json:"qa_10"`
	QA10Options []string        `json:"qa_10_options"`
	QA10Answer  string          `json:"qa_10_answer"`
	QA10Score   int64           `json:"qa_10_score"`
	QA11        string          `json:"qa_11"`
	QA11Options []string        `json:"qa_11_options"`
	QA11Answer  string          `json:"qa_11_answer"`
	QA11Score   int64           `json:"qa_11_score"`
	QA12        string          `json:"qa_12"`
	QA12Options []string        `json:"qa_12_options"`
	QA12Answer  string          `json:"qa_12_answer"`
	QA12Score   int64           `json:"qa_12_score"`
	QA13        string          `json:"qa_13"`
	QA13Options []string        `json:"qa_13_options"`
	QA13Answer  string          `json:"qa_13_answer"`
	QA13Score   int64           `json:"qa_13_score"`
	QA14        string          `json:"qa_14"`
	QA14Options []string        `json:"qa_14_options"`
	QA14Answer  string          `json:"qa_14_answer"`
	QA14Score   int64           `json:"qa_14_score"`
	QA15        string          `json:"qa_15"`
	QA15Options []string        `json:"qa_15_options"`
	QA15Answer  string          `json:"qa_15_answer"`
	QA15Score   int64           `json:"qa_15_score"`
	QA16        string          `json:"qa_16"`
	QA16Options []string        `json:"qa_16_options"`
	QA16Answer  string          `json:"qa_16_answer"`
	QA16Score   int64           `json:"qa_16_score"`
	QA17        string          `json:"qa_17"`
	QA17Options []string        `json:"qa_17_options"`
	QA17Answer  string          `json:"qa_17_answer"`
	QA17Score   int64           `json:"qa_17_score"`
	QA18        string          `json:"qa_18"`
	QA18Options []string        `json:"qa_18_options"`
	QA18Answer  string          `json:"qa_18_answer"`
	QA18Score   int64           `json:"qa_18_score"`
	QA19        string          `json:"qa_19"`
	QA19Options []string        `json:"qa_19_options"`
	QA19Answer  string          `json:"qa_19_answer"`
	QA19Score   int64           `json:"qa_19_score"`
	QA20        string          `json:"qa_20"`
	QA20Options []string        `json:"qa_20_options"`
	QA20Answer  string          `json:"qa_20_answer"`
	QA20Score   int64           `json:"qa_20_score"`
}

// 快問快答題目資訊
type QuestionModel struct {
	Question string `json:"question" example:"qa"`
	// Picture  string   `json:"picture" example:"picture"`
	Options []string `json:"options" example:"0&&&1&&&2&&&3"`
	Answer  string   `json:"answer" example:"0"`
	Score   int64    `json:"score" example:"10"`
}

// EditGameModel 資料表欄位
type EditGameModel struct {
	UserID     string `json:"user_id" example:"user_id"`
	ActivityID string `json:"activity_id" example:"activity_id"`
	GameID     string `json:"game_id" example:"game_id"`
	// Game         string `json:"game" example:"game name"`
	Title         string `json:"title" example:"game title"`
	GameType      string `json:"game_type" example:"game type"`
	LimitTime     string `json:"limit_time" example:"open、close"`
	Second        string `json:"second" example:"30"`
	MaxPeople     string `json:"max_people" example:"100(依照用戶權限判斷上限)"`
	People        string `json:"people" example:"100(依照max_people資料判斷上限)"`
	MaxTimes      string `json:"max_times" example:"10"`
	Allow         string `json:"allow" example:"open、close"`
	Percent       string `json:"percent" example:"0-100"`
	FirstPrize    string `json:"first_prize" example:"50(上限為50人)"`
	SecondPrize   string `json:"second_prize" example:"50(上限為50人)"`
	ThirdPrize    string `json:"third_prize" example:"100(上限為100人)"`
	GeneralPrize  string `json:"general_prize" example:"800(上限為800人)"`
	Topic         string `json:"topic" example:"01_classic"`
	Skin          string `json:"skin" example:"classic"`
	Music         string `json:"music" example:"classic"`
	DisplayName   string `json:"display_name" example:"open、close"`   // 是否顯示中獎人員姓名頭像
	// GameOrder     int64 `json:"game_order" example:"1"`              // 遊戲場次排序
	BoxReflection string `json:"box_reflection" example:"open、close"` // 扭蛋機遊戲開關盒反射
	SamePeople    string `json:"same_people" example:"open、close"`    // 拔河遊戲人數是否一致

	// 拔河遊戲
	AllowChooseTeam  string `json:"allow_choose_team" example:"open、close"` // 允許玩家選擇隊伍
	LeftTeamName     string `json:"left_team_name" example:"name"`          // 左邊隊伍名稱
	LeftTeamPicture  string `json:"left_team_picture" example:"picture"`    // 左邊隊伍照片
	RightTeamName    string `json:"right_team_name" example:"name"`         // 右邊隊伍名稱
	RightTeamPicture string `json:"right_team_picture" example:"picture"`   // 右邊隊伍照片
	Prize            string `json:"prize" example:"uniform、all"`            // 獎品發放

	// 賓果遊戲
	MaxNumber  string `json:"max_number" example:"0"`  // 最大號碼
	BingoLine  string `json:"bingo_line" example:"0"`  // 賓果連線數
	RoundPrize string `json:"round_prize" example:"0"` // 每輪發獎人數

	// 扭蛋機遊戲
	GachaMachineReflection string `json:"gacha_machine_reflection" example:"open、close"` // 球的反射度
	ReflectiveSwitch       string `json:"reflective_switch" example:"open、close"`        // 反射開關

	// 投票遊戲
	VoteScreen       string `json:"vote_screen" example:"bar_chart、rank、detail_information"` // 投票畫面(長條圖顯示、排名顯示、詳細資訊顯示)
	VoteTimes        string `json:"vote_times" example:"0"`                                  // 人員投票次數
	VoteMethod       string `json:"vote_method" example:"all_vote、single_group、all_group"`   // 投票模式(全選項投票)
	VoteMethodPlayer string `json:"vote_method_player" example:"one_vote、free_vote"`         // 玩家投票方式(一個選項一票、自由投票)
	VoteRestriction  string `json:"vote_restriction" example:"all_player、special_officer"`   // 投票限制(所有人員都能投票、特殊人員才能投票)
	AvatarShape      string `json:"avatar_shape" example:"circle、square"`                    // 選項框是圓形還是方形
	VoteStartTime    string `json:"vote_start_time" example:""`                              // 投票開始時間
	VoteEndTime      string `json:"vote_end_time" example:""`                                // 投票結束時間
	AutoPlay         string `json:"auto_play" example:"open、close"`                          // 自動輪播
	ShowRank         string `json:"show_rank" example:"open、close"`                          // 排名展示
	TitleSwitch      string `json:"title_switch" example:"open、close"`                       // 場次名稱
	ArrangementGuest string `json:"arrangement_guest" example:"list、side_by_side"`           // 玩家端選項排列方式

	// 敲敲樂自定義
	WhackmoleClassicHPic01 string `json:"whackmole_classic_h_pic_01" example:"picture"`
	WhackmoleClassicHPic02 string `json:"whackmole_classic_h_pic_02" example:"picture"`
	WhackmoleClassicHPic03 string `json:"whackmole_classic_h_pic_03" example:"picture"`
	WhackmoleClassicHPic04 string `json:"whackmole_classic_h_pic_04" example:"picture"`
	WhackmoleClassicHPic05 string `json:"whackmole_classic_h_pic_05" example:"picture"`
	WhackmoleClassicHPic06 string `json:"whackmole_classic_h_pic_06" example:"picture"`
	WhackmoleClassicHPic07 string `json:"whackmole_classic_h_pic_07" example:"picture"`
	WhackmoleClassicHPic08 string `json:"whackmole_classic_h_pic_08" example:"picture"`
	WhackmoleClassicHPic09 string `json:"whackmole_classic_h_pic_09" example:"picture"`
	WhackmoleClassicHPic10 string `json:"whackmole_classic_h_pic_10" example:"picture"`
	WhackmoleClassicHPic11 string `json:"whackmole_classic_h_pic_11" example:"picture"`
	WhackmoleClassicHPic12 string `json:"whackmole_classic_h_pic_12" example:"picture"`
	WhackmoleClassicHPic13 string `json:"whackmole_classic_h_pic_13" example:"picture"`
	WhackmoleClassicHPic14 string `json:"whackmole_classic_h_pic_14" example:"picture"`
	WhackmoleClassicHPic15 string `json:"whackmole_classic_h_pic_15" example:"picture"`
	WhackmoleClassicGPic01 string `json:"whackmole_classic_g_pic_01" example:"picture"`
	WhackmoleClassicGPic02 string `json:"whackmole_classic_g_pic_02" example:"picture"`
	WhackmoleClassicGPic03 string `json:"whackmole_classic_g_pic_03" example:"picture"`
	WhackmoleClassicGPic04 string `json:"whackmole_classic_g_pic_04" example:"picture"`
	WhackmoleClassicGPic05 string `json:"whackmole_classic_g_pic_05" example:"picture"`
	WhackmoleClassicCPic01 string `json:"whackmole_classic_c_pic_01" example:"picture"`
	WhackmoleClassicCPic02 string `json:"whackmole_classic_c_pic_02" example:"picture"`
	WhackmoleClassicCPic03 string `json:"whackmole_classic_c_pic_03" example:"picture"`
	WhackmoleClassicCPic04 string `json:"whackmole_classic_c_pic_04" example:"picture"`
	WhackmoleClassicCPic05 string `json:"whackmole_classic_c_pic_05" example:"picture"`
	WhackmoleClassicCPic06 string `json:"whackmole_classic_c_pic_06" example:"picture"`
	WhackmoleClassicCPic07 string `json:"whackmole_classic_c_pic_07" example:"picture"`
	WhackmoleClassicCPic08 string `json:"whackmole_classic_c_pic_08" example:"picture"`
	WhackmoleClassicCAni01 string `json:"whackmole_classic_c_ani_01" example:"picture"`

	WhackmoleHalloweenHPic01 string `json:"whackmole_halloween_h_pic_01" example:"picture"`
	WhackmoleHalloweenHPic02 string `json:"whackmole_halloween_h_pic_02" example:"picture"`
	WhackmoleHalloweenHPic03 string `json:"whackmole_halloween_h_pic_03" example:"picture"`
	WhackmoleHalloweenHPic04 string `json:"whackmole_halloween_h_pic_04" example:"picture"`
	WhackmoleHalloweenHPic05 string `json:"whackmole_halloween_h_pic_05" example:"picture"`
	WhackmoleHalloweenHPic06 string `json:"whackmole_halloween_h_pic_06" example:"picture"`
	WhackmoleHalloweenHPic07 string `json:"whackmole_halloween_h_pic_07" example:"picture"`
	WhackmoleHalloweenHPic08 string `json:"whackmole_halloween_h_pic_08" example:"picture"`
	WhackmoleHalloweenHPic09 string `json:"whackmole_halloween_h_pic_09" example:"picture"`
	WhackmoleHalloweenHPic10 string `json:"whackmole_halloween_h_pic_10" example:"picture"`
	WhackmoleHalloweenHPic11 string `json:"whackmole_halloween_h_pic_11" example:"picture"`
	WhackmoleHalloweenHPic12 string `json:"whackmole_halloween_h_pic_12" example:"picture"`
	WhackmoleHalloweenHPic13 string `json:"whackmole_halloween_h_pic_13" example:"picture"`
	WhackmoleHalloweenHPic14 string `json:"whackmole_halloween_h_pic_14" example:"picture"`
	WhackmoleHalloweenHPic15 string `json:"whackmole_halloween_h_pic_15" example:"picture"`
	WhackmoleHalloweenGPic01 string `json:"whackmole_halloween_g_pic_01" example:"picture"`
	WhackmoleHalloweenGPic02 string `json:"whackmole_halloween_g_pic_02" example:"picture"`
	WhackmoleHalloweenGPic03 string `json:"whackmole_halloween_g_pic_03" example:"picture"`
	WhackmoleHalloweenGPic04 string `json:"whackmole_halloween_g_pic_04" example:"picture"`
	WhackmoleHalloweenGPic05 string `json:"whackmole_halloween_g_pic_05" example:"picture"`
	WhackmoleHalloweenCPic01 string `json:"whackmole_halloween_c_pic_01" example:"picture"`
	WhackmoleHalloweenCPic02 string `json:"whackmole_halloween_c_pic_02" example:"picture"`
	WhackmoleHalloweenCPic03 string `json:"whackmole_halloween_c_pic_03" example:"picture"`
	WhackmoleHalloweenCPic04 string `json:"whackmole_halloween_c_pic_04" example:"picture"`
	WhackmoleHalloweenCPic05 string `json:"whackmole_halloween_c_pic_05" example:"picture"`
	WhackmoleHalloweenCPic06 string `json:"whackmole_halloween_c_pic_06" example:"picture"`
	WhackmoleHalloweenCPic07 string `json:"whackmole_halloween_c_pic_07" example:"picture"`
	WhackmoleHalloweenCPic08 string `json:"whackmole_halloween_c_pic_08" example:"picture"`
	WhackmoleHalloweenCAni01 string `json:"whackmole_halloween_c_ani_01" example:"picture"`

	WhackmoleChristmasHPic01 string `json:"whackmole_christmas_h_pic_01" example:"picture"`
	WhackmoleChristmasHPic02 string `json:"whackmole_christmas_h_pic_02" example:"picture"`
	WhackmoleChristmasHPic03 string `json:"whackmole_christmas_h_pic_03" example:"picture"`
	WhackmoleChristmasHPic04 string `json:"whackmole_christmas_h_pic_04" example:"picture"`
	WhackmoleChristmasHPic05 string `json:"whackmole_christmas_h_pic_05" example:"picture"`
	WhackmoleChristmasHPic06 string `json:"whackmole_christmas_h_pic_06" example:"picture"`
	WhackmoleChristmasHPic07 string `json:"whackmole_christmas_h_pic_07" example:"picture"`
	WhackmoleChristmasHPic08 string `json:"whackmole_christmas_h_pic_08" example:"picture"`
	WhackmoleChristmasHPic09 string `json:"whackmole_christmas_h_pic_09" example:"picture"`
	WhackmoleChristmasHPic10 string `json:"whackmole_christmas_h_pic_10" example:"picture"`
	WhackmoleChristmasHPic11 string `json:"whackmole_christmas_h_pic_11" example:"picture"`
	WhackmoleChristmasHPic12 string `json:"whackmole_christmas_h_pic_12" example:"picture"`
	WhackmoleChristmasHPic13 string `json:"whackmole_christmas_h_pic_13" example:"picture"`
	WhackmoleChristmasHPic14 string `json:"whackmole_christmas_h_pic_14" example:"picture"`
	WhackmoleChristmasGPic01 string `json:"whackmole_christmas_g_pic_01" example:"picture"`
	WhackmoleChristmasGPic02 string `json:"whackmole_christmas_g_pic_02" example:"picture"`
	WhackmoleChristmasGPic03 string `json:"whackmole_christmas_g_pic_03" example:"picture"`
	WhackmoleChristmasGPic04 string `json:"whackmole_christmas_g_pic_04" example:"picture"`
	WhackmoleChristmasGPic05 string `json:"whackmole_christmas_g_pic_05" example:"picture"`
	WhackmoleChristmasGPic06 string `json:"whackmole_christmas_g_pic_06" example:"picture"`
	WhackmoleChristmasGPic07 string `json:"whackmole_christmas_g_pic_07" example:"picture"`
	WhackmoleChristmasGPic08 string `json:"whackmole_christmas_g_pic_08" example:"picture"`
	WhackmoleChristmasCPic01 string `json:"whackmole_christmas_c_pic_01" example:"picture"`
	WhackmoleChristmasCPic02 string `json:"whackmole_christmas_c_pic_02" example:"picture"`
	WhackmoleChristmasCPic03 string `json:"whackmole_christmas_c_pic_03" example:"picture"`
	WhackmoleChristmasCPic04 string `json:"whackmole_christmas_c_pic_04" example:"picture"`
	WhackmoleChristmasCPic05 string `json:"whackmole_christmas_c_pic_05" example:"picture"`
	WhackmoleChristmasCPic06 string `json:"whackmole_christmas_c_pic_06" example:"picture"`
	WhackmoleChristmasCPic07 string `json:"whackmole_christmas_c_pic_07" example:"picture"`
	WhackmoleChristmasCPic08 string `json:"whackmole_christmas_c_pic_08" example:"picture"`
	WhackmoleChristmasCAni01 string `json:"whackmole_christmas_c_ani_01" example:"picture"`
	WhackmoleChristmasCAni02 string `json:"whackmole_christmas_c_ani_02" example:"picture"`

	// 敲敲樂音樂
	WhackmoleBgmStart  string `json:"whackmole_bgm_start" example:"picture"`  // 遊戲開始
	WhackmoleBgmGaming string `json:"whackmole_bgm_gaming" example:"picture"` // 遊戲進行中
	WhackmoleBgmEnd    string `json:"whackmole_bgm_end" example:"picture"`    // 遊戲結束

	// 搖號抽獎自定義
	DrawNumbersClassicHPic01 string `json:"draw_numbers_classic_h_pic_01" example:"picture"`
	DrawNumbersClassicHPic02 string `json:"draw_numbers_classic_h_pic_02" example:"picture"`
	DrawNumbersClassicHPic03 string `json:"draw_numbers_classic_h_pic_03" example:"picture"`
	DrawNumbersClassicHPic04 string `json:"draw_numbers_classic_h_pic_04" example:"picture"`
	DrawNumbersClassicHPic05 string `json:"draw_numbers_classic_h_pic_05" example:"picture"`
	DrawNumbersClassicHPic06 string `json:"draw_numbers_classic_h_pic_06" example:"picture"`
	DrawNumbersClassicHPic07 string `json:"draw_numbers_classic_h_pic_07" example:"picture"`
	DrawNumbersClassicHPic08 string `json:"draw_numbers_classic_h_pic_08" example:"picture"`
	DrawNumbersClassicHPic09 string `json:"draw_numbers_classic_h_pic_09" example:"picture"`
	DrawNumbersClassicHPic10 string `json:"draw_numbers_classic_h_pic_10" example:"picture"`
	DrawNumbersClassicHPic11 string `json:"draw_numbers_classic_h_pic_11" example:"picture"`
	DrawNumbersClassicHPic12 string `json:"draw_numbers_classic_h_pic_12" example:"picture"`
	DrawNumbersClassicHPic13 string `json:"draw_numbers_classic_h_pic_13" example:"picture"`
	DrawNumbersClassicHPic14 string `json:"draw_numbers_classic_h_pic_14" example:"picture"`
	DrawNumbersClassicHPic15 string `json:"draw_numbers_classic_h_pic_15" example:"picture"`
	DrawNumbersClassicHPic16 string `json:"draw_numbers_classic_h_pic_16" example:"picture"`
	DrawNumbersClassicHAni01 string `json:"draw_numbers_classic_h_ani_01" example:"picture"`

	DrawNumbersGoldHPic01 string `json:"draw_numbers_gold_h_pic_01" example:"picture"`
	DrawNumbersGoldHPic02 string `json:"draw_numbers_gold_h_pic_02" example:"picture"`
	DrawNumbersGoldHPic03 string `json:"draw_numbers_gold_h_pic_03" example:"picture"`
	DrawNumbersGoldHPic04 string `json:"draw_numbers_gold_h_pic_04" example:"picture"`
	DrawNumbersGoldHPic05 string `json:"draw_numbers_gold_h_pic_05" example:"picture"`
	DrawNumbersGoldHPic06 string `json:"draw_numbers_gold_h_pic_06" example:"picture"`
	DrawNumbersGoldHPic07 string `json:"draw_numbers_gold_h_pic_07" example:"picture"`
	DrawNumbersGoldHPic08 string `json:"draw_numbers_gold_h_pic_08" example:"picture"`
	DrawNumbersGoldHPic09 string `json:"draw_numbers_gold_h_pic_09" example:"picture"`
	DrawNumbersGoldHPic10 string `json:"draw_numbers_gold_h_pic_10" example:"picture"`
	DrawNumbersGoldHPic11 string `json:"draw_numbers_gold_h_pic_11" example:"picture"`
	DrawNumbersGoldHPic12 string `json:"draw_numbers_gold_h_pic_12" example:"picture"`
	DrawNumbersGoldHPic13 string `json:"draw_numbers_gold_h_pic_13" example:"picture"`
	DrawNumbersGoldHPic14 string `json:"draw_numbers_gold_h_pic_14" example:"picture"`
	DrawNumbersGoldHAni01 string `json:"draw_numbers_gold_h_ani_01" example:"picture"`
	DrawNumbersGoldHAni02 string `json:"draw_numbers_gold_h_ani_02" example:"picture"`
	DrawNumbersGoldHAni03 string `json:"draw_numbers_gold_h_ani_03" example:"picture"`

	DrawNumbersNewyearDragonHPic01 string `json:"draw_numbers_newyear_dragon_h_pic_01" example:"picture"`
	DrawNumbersNewyearDragonHPic02 string `json:"draw_numbers_newyear_dragon_h_pic_02" example:"picture"`
	DrawNumbersNewyearDragonHPic03 string `json:"draw_numbers_newyear_dragon_h_pic_03" example:"picture"`
	DrawNumbersNewyearDragonHPic04 string `json:"draw_numbers_newyear_dragon_h_pic_04" example:"picture"`
	DrawNumbersNewyearDragonHPic05 string `json:"draw_numbers_newyear_dragon_h_pic_05" example:"picture"`
	DrawNumbersNewyearDragonHPic06 string `json:"draw_numbers_newyear_dragon_h_pic_06" example:"picture"`
	DrawNumbersNewyearDragonHPic07 string `json:"draw_numbers_newyear_dragon_h_pic_07" example:"picture"`
	DrawNumbersNewyearDragonHPic08 string `json:"draw_numbers_newyear_dragon_h_pic_08" example:"picture"`
	DrawNumbersNewyearDragonHPic09 string `json:"draw_numbers_newyear_dragon_h_pic_09" example:"picture"`
	DrawNumbersNewyearDragonHPic10 string `json:"draw_numbers_newyear_dragon_h_pic_10" example:"picture"`
	DrawNumbersNewyearDragonHPic11 string `json:"draw_numbers_newyear_dragon_h_pic_11" example:"picture"`
	DrawNumbersNewyearDragonHPic12 string `json:"draw_numbers_newyear_dragon_h_pic_12" example:"picture"`
	DrawNumbersNewyearDragonHPic13 string `json:"draw_numbers_newyear_dragon_h_pic_13" example:"picture"`
	DrawNumbersNewyearDragonHPic14 string `json:"draw_numbers_newyear_dragon_h_pic_14" example:"picture"`
	DrawNumbersNewyearDragonHPic15 string `json:"draw_numbers_newyear_dragon_h_pic_15" example:"picture"`
	DrawNumbersNewyearDragonHPic16 string `json:"draw_numbers_newyear_dragon_h_pic_16" example:"picture"`
	DrawNumbersNewyearDragonHPic17 string `json:"draw_numbers_newyear_dragon_h_pic_17" example:"picture"`
	DrawNumbersNewyearDragonHPic18 string `json:"draw_numbers_newyear_dragon_h_pic_18" example:"picture"`
	DrawNumbersNewyearDragonHPic19 string `json:"draw_numbers_newyear_dragon_h_pic_19" example:"picture"`
	DrawNumbersNewyearDragonHPic20 string `json:"draw_numbers_newyear_dragon_h_pic_20" example:"picture"`
	DrawNumbersNewyearDragonHAni01 string `json:"draw_numbers_newyear_dragon_h_ani_01" example:"picture"`
	DrawNumbersNewyearDragonHAni02 string `json:"draw_numbers_newyear_dragon_h_ani_02" example:"picture"`

	DrawNumbersCherryHPic01 string `json:"draw_numbers_cherry_h_pic_01" example:"picture"`
	DrawNumbersCherryHPic02 string `json:"draw_numbers_cherry_h_pic_02" example:"picture"`
	DrawNumbersCherryHPic03 string `json:"draw_numbers_cherry_h_pic_03" example:"picture"`
	DrawNumbersCherryHPic04 string `json:"draw_numbers_cherry_h_pic_04" example:"picture"`
	DrawNumbersCherryHPic05 string `json:"draw_numbers_cherry_h_pic_05" example:"picture"`
	DrawNumbersCherryHPic06 string `json:"draw_numbers_cherry_h_pic_06" example:"picture"`
	DrawNumbersCherryHPic07 string `json:"draw_numbers_cherry_h_pic_07" example:"picture"`
	DrawNumbersCherryHPic08 string `json:"draw_numbers_cherry_h_pic_08" example:"picture"`
	DrawNumbersCherryHPic09 string `json:"draw_numbers_cherry_h_pic_09" example:"picture"`
	DrawNumbersCherryHPic10 string `json:"draw_numbers_cherry_h_pic_10" example:"picture"`
	DrawNumbersCherryHPic11 string `json:"draw_numbers_cherry_h_pic_11" example:"picture"`
	DrawNumbersCherryHPic12 string `json:"draw_numbers_cherry_h_pic_12" example:"picture"`
	DrawNumbersCherryHPic13 string `json:"draw_numbers_cherry_h_pic_13" example:"picture"`
	DrawNumbersCherryHPic14 string `json:"draw_numbers_cherry_h_pic_14" example:"picture"`
	DrawNumbersCherryHPic15 string `json:"draw_numbers_cherry_h_pic_15" example:"picture"`
	DrawNumbersCherryHPic16 string `json:"draw_numbers_cherry_h_pic_16" example:"picture"`
	DrawNumbersCherryHPic17 string `json:"draw_numbers_cherry_h_pic_17" example:"picture"`
	DrawNumbersCherryHAni01 string `json:"draw_numbers_cherry_h_ani_01" example:"picture"`
	DrawNumbersCherryHAni02 string `json:"draw_numbers_cherry_h_ani_02" example:"picture"`
	DrawNumbersCherryHAni03 string `json:"draw_numbers_cherry_h_ani_03" example:"picture"`
	DrawNumbersCherryHAni04 string `json:"draw_numbers_cherry_h_ani_04" example:"picture"`

	// 太空主題
	DrawNumbers3DSpaceHPic01 string `json:"draw_numbers_3D_space_h_pic_01" example:"picture"`
	DrawNumbers3DSpaceHPic02 string `json:"draw_numbers_3D_space_h_pic_02" example:"picture"`
	DrawNumbers3DSpaceHPic03 string `json:"draw_numbers_3D_space_h_pic_03" example:"picture"`
	DrawNumbers3DSpaceHPic04 string `json:"draw_numbers_3D_space_h_pic_04" example:"picture"`
	DrawNumbers3DSpaceHPic05 string `json:"draw_numbers_3D_space_h_pic_05" example:"picture"`
	DrawNumbers3DSpaceHPic06 string `json:"draw_numbers_3D_space_h_pic_06" example:"picture"`
	DrawNumbers3DSpaceHPic07 string `json:"draw_numbers_3D_space_h_pic_07" example:"picture"`
	DrawNumbers3DSpaceHPic08 string `json:"draw_numbers_3D_space_h_pic_08" example:"picture"`

	// 音樂
	DrawNumbersBgmGaming string `json:"draw_numbers_bgm_gaming" example:"picture"` // 遊戲進行中

	// 快問快答自定義
	QAClassicHPic01 string `json:"qa_classic_h_pic_01" example:"picture"`
	QAClassicHPic02 string `json:"qa_classic_h_pic_02" example:"picture"`
	QAClassicHPic03 string `json:"qa_classic_h_pic_03" example:"picture"`
	QAClassicHPic04 string `json:"qa_classic_h_pic_04" example:"picture"`
	QAClassicHPic05 string `json:"qa_classic_h_pic_05" example:"picture"`
	QAClassicHPic06 string `json:"qa_classic_h_pic_06" example:"picture"`
	QAClassicHPic07 string `json:"qa_classic_h_pic_07" example:"picture"`
	QAClassicHPic08 string `json:"qa_classic_h_pic_08" example:"picture"`
	QAClassicHPic09 string `json:"qa_classic_h_pic_09" example:"picture"`
	QAClassicHPic10 string `json:"qa_classic_h_pic_10" example:"picture"`
	QAClassicHPic11 string `json:"qa_classic_h_pic_11" example:"picture"`
	QAClassicHPic12 string `json:"qa_classic_h_pic_12" example:"picture"`
	QAClassicHPic13 string `json:"qa_classic_h_pic_13" example:"picture"`
	QAClassicHPic14 string `json:"qa_classic_h_pic_14" example:"picture"`
	QAClassicHPic15 string `json:"qa_classic_h_pic_15" example:"picture"`
	QAClassicHPic16 string `json:"qa_classic_h_pic_16" example:"picture"`
	QAClassicHPic17 string `json:"qa_classic_h_pic_17" example:"picture"`
	QAClassicHPic18 string `json:"qa_classic_h_pic_18" example:"picture"`
	QAClassicHPic19 string `json:"qa_classic_h_pic_19" example:"picture"`
	QAClassicHPic20 string `json:"qa_classic_h_pic_20" example:"picture"`
	QAClassicHPic21 string `json:"qa_classic_h_pic_21" example:"picture"`
	QAClassicHPic22 string `json:"qa_classic_h_pic_22" example:"picture"`
	QAClassicGPic01 string `json:"qa_classic_g_pic_01" example:"picture"`
	QAClassicGPic02 string `json:"qa_classic_g_pic_02" example:"picture"`
	QAClassicGPic03 string `json:"qa_classic_g_pic_03" example:"picture"`
	QAClassicGPic04 string `json:"qa_classic_g_pic_04" example:"picture"`
	QAClassicGPic05 string `json:"qa_classic_g_pic_05" example:"picture"`
	QAClassicCPic01 string `json:"qa_classic_c_pic_01" example:"picture"`
	QAClassicHAni01 string `json:"qa_classic_h_ani_01" example:"picture"`
	QAClassicHAni02 string `json:"qa_classic_h_ani_02" example:"picture"`
	QAClassicGAni01 string `json:"qa_classic_g_ani_01" example:"picture"`
	QAClassicGAni02 string `json:"qa_classic_g_ani_02" example:"picture"`

	QAElectricHPic01 string `json:"qa_electric_h_pic_01" example:"picture"`
	QAElectricHPic02 string `json:"qa_electric_h_pic_02" example:"picture"`
	QAElectricHPic03 string `json:"qa_electric_h_pic_03" example:"picture"`
	QAElectricHPic04 string `json:"qa_electric_h_pic_04" example:"picture"`
	QAElectricHPic05 string `json:"qa_electric_h_pic_05" example:"picture"`
	QAElectricHPic06 string `json:"qa_electric_h_pic_06" example:"picture"`
	QAElectricHPic07 string `json:"qa_electric_h_pic_07" example:"picture"`
	QAElectricHPic08 string `json:"qa_electric_h_pic_08" example:"picture"`
	QAElectricHPic09 string `json:"qa_electric_h_pic_09" example:"picture"`
	QAElectricHPic10 string `json:"qa_electric_h_pic_10" example:"picture"`
	QAElectricHPic11 string `json:"qa_electric_h_pic_11" example:"picture"`
	QAElectricHPic12 string `json:"qa_electric_h_pic_12" example:"picture"`
	QAElectricHPic13 string `json:"qa_electric_h_pic_13" example:"picture"`
	QAElectricHPic14 string `json:"qa_electric_h_pic_14" example:"picture"`
	QAElectricHPic15 string `json:"qa_electric_h_pic_15" example:"picture"`
	QAElectricHPic16 string `json:"qa_electric_h_pic_16" example:"picture"`
	QAElectricHPic17 string `json:"qa_electric_h_pic_17" example:"picture"`
	QAElectricHPic18 string `json:"qa_electric_h_pic_18" example:"picture"`
	QAElectricHPic19 string `json:"qa_electric_h_pic_19" example:"picture"`
	QAElectricHPic20 string `json:"qa_electric_h_pic_20" example:"picture"`
	QAElectricHPic21 string `json:"qa_electric_h_pic_21" example:"picture"`
	QAElectricHPic22 string `json:"qa_electric_h_pic_22" example:"picture"`
	QAElectricHPic23 string `json:"qa_electric_h_pic_23" example:"picture"`
	QAElectricHPic24 string `json:"qa_electric_h_pic_24" example:"picture"`
	QAElectricHPic25 string `json:"qa_electric_h_pic_25" example:"picture"`
	QAElectricHPic26 string `json:"qa_electric_h_pic_26" example:"picture"`
	QAElectricGPic01 string `json:"qa_electric_g_pic_01" example:"picture"`
	QAElectricGPic02 string `json:"qa_electric_g_pic_02" example:"picture"`
	QAElectricGPic03 string `json:"qa_electric_g_pic_03" example:"picture"`
	QAElectricGPic04 string `json:"qa_electric_g_pic_04" example:"picture"`
	QAElectricGPic05 string `json:"qa_electric_g_pic_05" example:"picture"`
	QAElectricGPic06 string `json:"qa_electric_g_pic_06" example:"picture"`
	QAElectricGPic07 string `json:"qa_electric_g_pic_07" example:"picture"`
	QAElectricGPic08 string `json:"qa_electric_g_pic_08" example:"picture"`
	QAElectricGPic09 string `json:"qa_electric_g_pic_09" example:"picture"`
	QAElectricCPic01 string `json:"qa_electric_c_pic_01" example:"picture"`
	QAElectricHAni01 string `json:"qa_electric_h_ani_01" example:"picture"`
	QAElectricHAni02 string `json:"qa_electric_h_ani_02" example:"picture"`
	QAElectricHAni03 string `json:"qa_electric_h_ani_03" example:"picture"`
	QAElectricHAni04 string `json:"qa_electric_h_ani_04" example:"picture"`
	QAElectricHAni05 string `json:"qa_electric_h_ani_05" example:"picture"`
	QAElectricGAni01 string `json:"qa_electric_g_ani_01" example:"picture"`
	QAElectricGAni02 string `json:"qa_electric_g_ani_02" example:"picture"`
	QAElectricCAni01 string `json:"qa_electric_c_ani_01" example:"picture"`

	QAMoonfestivalHPic01 string `json:"qa_moonfestival_h_pic_01" example:"picture"`
	QAMoonfestivalHPic02 string `json:"qa_moonfestival_h_pic_02" example:"picture"`
	QAMoonfestivalHPic03 string `json:"qa_moonfestival_h_pic_03" example:"picture"`
	QAMoonfestivalHPic04 string `json:"qa_moonfestival_h_pic_04" example:"picture"`
	QAMoonfestivalHPic05 string `json:"qa_moonfestival_h_pic_05" example:"picture"`
	QAMoonfestivalHPic06 string `json:"qa_moonfestival_h_pic_06" example:"picture"`
	QAMoonfestivalHPic07 string `json:"qa_moonfestival_h_pic_07" example:"picture"`
	QAMoonfestivalHPic08 string `json:"qa_moonfestival_h_pic_08" example:"picture"`
	QAMoonfestivalHPic09 string `json:"qa_moonfestival_h_pic_09" example:"picture"`
	QAMoonfestivalHPic10 string `json:"qa_moonfestival_h_pic_10" example:"picture"`
	QAMoonfestivalHPic11 string `json:"qa_moonfestival_h_pic_11" example:"picture"`
	QAMoonfestivalHPic12 string `json:"qa_moonfestival_h_pic_12" example:"picture"`
	QAMoonfestivalHPic13 string `json:"qa_moonfestival_h_pic_13" example:"picture"`
	QAMoonfestivalHPic14 string `json:"qa_moonfestival_h_pic_14" example:"picture"`
	QAMoonfestivalHPic15 string `json:"qa_moonfestival_h_pic_15" example:"picture"`
	QAMoonfestivalHPic16 string `json:"qa_moonfestival_h_pic_16" example:"picture"`
	QAMoonfestivalHPic17 string `json:"qa_moonfestival_h_pic_17" example:"picture"`
	QAMoonfestivalHPic18 string `json:"qa_moonfestival_h_pic_18" example:"picture"`
	QAMoonfestivalHPic19 string `json:"qa_moonfestival_h_pic_19" example:"picture"`
	QAMoonfestivalHPic20 string `json:"qa_moonfestival_h_pic_20" example:"picture"`
	QAMoonfestivalHPic21 string `json:"qa_moonfestival_h_pic_21" example:"picture"`
	QAMoonfestivalHPic22 string `json:"qa_moonfestival_h_pic_22" example:"picture"`
	QAMoonfestivalHPic23 string `json:"qa_moonfestival_h_pic_23" example:"picture"`
	QAMoonfestivalHPic24 string `json:"qa_moonfestival_h_pic_24" example:"picture"`
	QAMoonfestivalGPic01 string `json:"qa_moonfestival_g_pic_01" example:"picture"`
	QAMoonfestivalGPic02 string `json:"qa_moonfestival_g_pic_02" example:"picture"`
	QAMoonfestivalGPic03 string `json:"qa_moonfestival_g_pic_03" example:"picture"`
	QAMoonfestivalGPic04 string `json:"qa_moonfestival_g_pic_04" example:"picture"`
	QAMoonfestivalGPic05 string `json:"qa_moonfestival_g_pic_05" example:"picture"`
	QAMoonfestivalCPic01 string `json:"qa_moonfestival_c_pic_01" example:"picture"`
	QAMoonfestivalCPic02 string `json:"qa_moonfestival_c_pic_02" example:"picture"`
	QAMoonfestivalCPic03 string `json:"qa_moonfestival_c_pic_03" example:"picture"`
	QAMoonfestivalHAni01 string `json:"qa_moonfestival_h_ani_01" example:"picture"`
	QAMoonfestivalHAni02 string `json:"qa_moonfestival_h_ani_02" example:"picture"`
	QAMoonfestivalGAni01 string `json:"qa_moonfestival_g_ani_01" example:"picture"`
	QAMoonfestivalGAni02 string `json:"qa_moonfestival_g_ani_02" example:"picture"`
	QAMoonfestivalGAni03 string `json:"qa_moonfestival_g_ani_03" example:"picture"`

	QANewyearDragonHPic01 string `json:"qa_newyear_dragon_h_pic_01" example:"picture"`
	QANewyearDragonHPic02 string `json:"qa_newyear_dragon_h_pic_02" example:"picture"`
	QANewyearDragonHPic03 string `json:"qa_newyear_dragon_h_pic_03" example:"picture"`
	QANewyearDragonHPic04 string `json:"qa_newyear_dragon_h_pic_04" example:"picture"`
	QANewyearDragonHPic05 string `json:"qa_newyear_dragon_h_pic_05" example:"picture"`
	QANewyearDragonHPic06 string `json:"qa_newyear_dragon_h_pic_06" example:"picture"`
	QANewyearDragonHPic07 string `json:"qa_newyear_dragon_h_pic_07" example:"picture"`
	QANewyearDragonHPic08 string `json:"qa_newyear_dragon_h_pic_08" example:"picture"`
	QANewyearDragonHPic09 string `json:"qa_newyear_dragon_h_pic_09" example:"picture"`
	QANewyearDragonHPic10 string `json:"qa_newyear_dragon_h_pic_10" example:"picture"`
	QANewyearDragonHPic11 string `json:"qa_newyear_dragon_h_pic_11" example:"picture"`
	QANewyearDragonHPic12 string `json:"qa_newyear_dragon_h_pic_12" example:"picture"`
	QANewyearDragonHPic13 string `json:"qa_newyear_dragon_h_pic_13" example:"picture"`
	QANewyearDragonHPic14 string `json:"qa_newyear_dragon_h_pic_14" example:"picture"`
	QANewyearDragonHPic15 string `json:"qa_newyear_dragon_h_pic_15" example:"picture"`
	QANewyearDragonHPic16 string `json:"qa_newyear_dragon_h_pic_16" example:"picture"`
	QANewyearDragonHPic17 string `json:"qa_newyear_dragon_h_pic_17" example:"picture"`
	QANewyearDragonHPic18 string `json:"qa_newyear_dragon_h_pic_18" example:"picture"`
	QANewyearDragonHPic19 string `json:"qa_newyear_dragon_h_pic_19" example:"picture"`
	QANewyearDragonHPic20 string `json:"qa_newyear_dragon_h_pic_20" example:"picture"`
	QANewyearDragonHPic21 string `json:"qa_newyear_dragon_h_pic_21" example:"picture"`
	QANewyearDragonHPic22 string `json:"qa_newyear_dragon_h_pic_22" example:"picture"`
	QANewyearDragonHPic23 string `json:"qa_newyear_dragon_h_pic_23" example:"picture"`
	QANewyearDragonHPic24 string `json:"qa_newyear_dragon_h_pic_24" example:"picture"`
	QANewyearDragonGPic01 string `json:"qa_newyear_dragon_g_pic_01" example:"picture"`
	QANewyearDragonGPic02 string `json:"qa_newyear_dragon_g_pic_02" example:"picture"`
	QANewyearDragonGPic03 string `json:"qa_newyear_dragon_g_pic_03" example:"picture"`
	QANewyearDragonGPic04 string `json:"qa_newyear_dragon_g_pic_04" example:"picture"`
	QANewyearDragonGPic05 string `json:"qa_newyear_dragon_g_pic_05" example:"picture"`
	QANewyearDragonGPic06 string `json:"qa_newyear_dragon_g_pic_06" example:"picture"`
	QANewyearDragonCPic01 string `json:"qa_newyear_dragon_c_pic_01" example:"picture"`
	QANewyearDragonHAni01 string `json:"qa_newyear_dragon_h_ani_01" example:"picture"`
	QANewyearDragonHAni02 string `json:"qa_newyear_dragon_h_ani_02" example:"picture"`
	QANewyearDragonGAni01 string `json:"qa_newyear_dragon_g_ani_01" example:"picture"`
	QANewyearDragonGAni02 string `json:"qa_newyear_dragon_g_ani_02" example:"picture"`
	QANewyearDragonGAni03 string `json:"qa_newyear_dragon_g_ani_03" example:"picture"`
	QANewyearDragonCAni01 string `json:"qa_newyear_dragon_c_ani_01" example:"picture"`

	// 音樂
	QABgmStart  string `json:"qa_bgm_start" example:"picture"`  // 遊戲開始
	QABgmGaming string `json:"qa_bgm_gaming" example:"picture"` // 遊戲進行中
	QABgmEnd    string `json:"qa_bgm_end" example:"picture"`    // 遊戲結束

	// 搖紅包自定義
	RedpackClassicHPic01 string `json:"redpack_classic_h_pic_01" example:"picture"`
	RedpackClassicHPic02 string `json:"redpack_classic_h_pic_02" example:"picture"`
	RedpackClassicHPic03 string `json:"redpack_classic_h_pic_03" example:"picture"`
	RedpackClassicHPic04 string `json:"redpack_classic_h_pic_04" example:"picture"`
	RedpackClassicHPic05 string `json:"redpack_classic_h_pic_05" example:"picture"`
	RedpackClassicHPic06 string `json:"redpack_classic_h_pic_06" example:"picture"`
	RedpackClassicHPic07 string `json:"redpack_classic_h_pic_07" example:"picture"`
	RedpackClassicHPic08 string `json:"redpack_classic_h_pic_08" example:"picture"`
	RedpackClassicHPic09 string `json:"redpack_classic_h_pic_09" example:"picture"`
	RedpackClassicHPic10 string `json:"redpack_classic_h_pic_10" example:"picture"`
	RedpackClassicHPic11 string `json:"redpack_classic_h_pic_11" example:"picture"`
	RedpackClassicHPic12 string `json:"redpack_classic_h_pic_12" example:"picture"`
	RedpackClassicHPic13 string `json:"redpack_classic_h_pic_13" example:"picture"`
	RedpackClassicGPic01 string `json:"redpack_classic_g_pic_01" example:"picture"`
	RedpackClassicGPic02 string `json:"redpack_classic_g_pic_02" example:"picture"`
	RedpackClassicGPic03 string `json:"redpack_classic_g_pic_03" example:"picture"`
	RedpackClassicHAni01 string `json:"redpack_classic_h_ani_01" example:"picture"`
	RedpackClassicHAni02 string `json:"redpack_classic_h_ani_02" example:"picture"`
	RedpackClassicGAni01 string `json:"redpack_classic_g_ani_01" example:"picture"`
	RedpackClassicGAni02 string `json:"redpack_classic_g_ani_02" example:"picture"`
	RedpackClassicGAni03 string `json:"redpack_classic_g_ani_03" example:"picture"`

	RedpackCherryHPic01 string `json:"redpack_cherry_h_pic_01" example:"picture"`
	RedpackCherryHPic02 string `json:"redpack_cherry_h_pic_02" example:"picture"`
	RedpackCherryHPic03 string `json:"redpack_cherry_h_pic_03" example:"picture"`
	RedpackCherryHPic04 string `json:"redpack_cherry_h_pic_04" example:"picture"`
	RedpackCherryHPic05 string `json:"redpack_cherry_h_pic_05" example:"picture"`
	RedpackCherryHPic06 string `json:"redpack_cherry_h_pic_06" example:"picture"`
	RedpackCherryHPic07 string `json:"redpack_cherry_h_pic_07" example:"picture"`
	RedpackCherryGPic01 string `json:"redpack_cherry_g_pic_01" example:"picture"`
	RedpackCherryGPic02 string `json:"redpack_cherry_g_pic_02" example:"picture"`
	RedpackCherryHAni01 string `json:"redpack_cherry_h_ani_01" example:"picture"`
	RedpackCherryHAni02 string `json:"redpack_cherry_h_ani_02" example:"picture"`
	RedpackCherryGAni01 string `json:"redpack_cherry_g_ani_01" example:"picture"`
	RedpackCherryGAni02 string `json:"redpack_cherry_g_ani_02" example:"picture"`

	RedpackChristmasHPic01 string `json:"redpack_christmas_h_pic_01" example:"picture"`
	RedpackChristmasHPic02 string `json:"redpack_christmas_h_pic_02" example:"picture"`
	RedpackChristmasHPic03 string `json:"redpack_christmas_h_pic_03" example:"picture"`
	RedpackChristmasHPic04 string `json:"redpack_christmas_h_pic_04" example:"picture"`
	RedpackChristmasHPic05 string `json:"redpack_christmas_h_pic_05" example:"picture"`
	RedpackChristmasHPic06 string `json:"redpack_christmas_h_pic_06" example:"picture"`
	RedpackChristmasHPic07 string `json:"redpack_christmas_h_pic_07" example:"picture"`
	RedpackChristmasHPic08 string `json:"redpack_christmas_h_pic_08" example:"picture"`
	RedpackChristmasHPic09 string `json:"redpack_christmas_h_pic_09" example:"picture"`
	RedpackChristmasHPic10 string `json:"redpack_christmas_h_pic_10" example:"picture"`
	RedpackChristmasHPic11 string `json:"redpack_christmas_h_pic_11" example:"picture"`
	RedpackChristmasHPic12 string `json:"redpack_christmas_h_pic_12" example:"picture"`
	RedpackChristmasHPic13 string `json:"redpack_christmas_h_pic_13" example:"picture"`
	RedpackChristmasGPic01 string `json:"redpack_christmas_g_pic_01" example:"picture"`
	RedpackChristmasGPic02 string `json:"redpack_christmas_g_pic_02" example:"picture"`
	RedpackChristmasGPic03 string `json:"redpack_christmas_g_pic_03" example:"picture"`
	RedpackChristmasGPic04 string `json:"redpack_christmas_g_pic_04" example:"picture"`
	RedpackChristmasCPic01 string `json:"redpack_christmas_c_pic_01" example:"picture"`
	RedpackChristmasCPic02 string `json:"redpack_christmas_c_pic_02" example:"picture"`
	RedpackChristmasCAni01 string `json:"redpack_christmas_c_ani_01" example:"picture"`

	// 音樂
	RedpackBgmStart  string `json:"redpack_bgm_start" example:"picture"`  // 遊戲開始
	RedpackBgmGaming string `json:"redpack_bgm_gaming" example:"picture"` // 遊戲進行中
	RedpackBgmEnd    string `json:"redpack_bgm_end" example:"picture"`    // 遊戲結束

	// 套紅包自定義
	RopepackClassicHPic01 string `json:"ropepack_classic_h_pic_01" example:"picture"`
	RopepackClassicHPic02 string `json:"ropepack_classic_h_pic_02" example:"picture"`
	RopepackClassicHPic03 string `json:"ropepack_classic_h_pic_03" example:"picture"`
	RopepackClassicHPic04 string `json:"ropepack_classic_h_pic_04" example:"picture"`
	RopepackClassicHPic05 string `json:"ropepack_classic_h_pic_05" example:"picture"`
	RopepackClassicHPic06 string `json:"ropepack_classic_h_pic_06" example:"picture"`
	RopepackClassicHPic07 string `json:"ropepack_classic_h_pic_07" example:"picture"`
	RopepackClassicHPic08 string `json:"ropepack_classic_h_pic_08" example:"picture"`
	RopepackClassicHPic09 string `json:"ropepack_classic_h_pic_09" example:"picture"`
	RopepackClassicHPic10 string `json:"ropepack_classic_h_pic_10" example:"picture"`
	RopepackClassicGPic01 string `json:"ropepack_classic_g_pic_01" example:"picture"`
	RopepackClassicGPic02 string `json:"ropepack_classic_g_pic_02" example:"picture"`
	RopepackClassicGPic03 string `json:"ropepack_classic_g_pic_03" example:"picture"`
	RopepackClassicGPic04 string `json:"ropepack_classic_g_pic_04" example:"picture"`
	RopepackClassicGPic05 string `json:"ropepack_classic_g_pic_05" example:"picture"`
	RopepackClassicGPic06 string `json:"ropepack_classic_g_pic_06" example:"picture"`
	RopepackClassicHAni01 string `json:"ropepack_classic_h_ani_01" example:"picture"`
	RopepackClassicGAni01 string `json:"ropepack_classic_g_ani_01" example:"picture"`
	RopepackClassicGAni02 string `json:"ropepack_classic_g_ani_02" example:"picture"`
	RopepackClassicCAni01 string `json:"ropepack_classic_c_ani_01" example:"picture"`

	RopepackNewyearRabbitHPic01 string `json:"ropepack_newyear_rabbit_h_pic_01" example:"picture"`
	RopepackNewyearRabbitHPic02 string `json:"ropepack_newyear_rabbit_h_pic_02" example:"picture"`
	RopepackNewyearRabbitHPic03 string `json:"ropepack_newyear_rabbit_h_pic_03" example:"picture"`
	RopepackNewyearRabbitHPic04 string `json:"ropepack_newyear_rabbit_h_pic_04" example:"picture"`
	RopepackNewyearRabbitHPic05 string `json:"ropepack_newyear_rabbit_h_pic_05" example:"picture"`
	RopepackNewyearRabbitHPic06 string `json:"ropepack_newyear_rabbit_h_pic_06" example:"picture"`
	RopepackNewyearRabbitHPic07 string `json:"ropepack_newyear_rabbit_h_pic_07" example:"picture"`
	RopepackNewyearRabbitHPic08 string `json:"ropepack_newyear_rabbit_h_pic_08" example:"picture"`
	RopepackNewyearRabbitHPic09 string `json:"ropepack_newyear_rabbit_h_pic_09" example:"picture"`
	RopepackNewyearRabbitGPic01 string `json:"ropepack_newyear_rabbit_g_pic_01" example:"picture"`
	RopepackNewyearRabbitGPic02 string `json:"ropepack_newyear_rabbit_g_pic_02" example:"picture"`
	RopepackNewyearRabbitGPic03 string `json:"ropepack_newyear_rabbit_g_pic_03" example:"picture"`
	RopepackNewyearRabbitHAni01 string `json:"ropepack_newyear_rabbit_h_ani_01" example:"picture"`
	RopepackNewyearRabbitGAni01 string `json:"ropepack_newyear_rabbit_g_ani_01" example:"picture"`
	RopepackNewyearRabbitGAni02 string `json:"ropepack_newyear_rabbit_g_ani_02" example:"picture"`
	RopepackNewyearRabbitGAni03 string `json:"ropepack_newyear_rabbit_g_ani_03" example:"picture"`
	RopepackNewyearRabbitCAni01 string `json:"ropepack_newyear_rabbit_c_ani_01" example:"picture"`
	RopepackNewyearRabbitCAni02 string `json:"ropepack_newyear_rabbit_c_ani_02" example:"picture"`

	RopepackMoonfestivalHPic01 string `json:"ropepack_moonfestival_h_pic_01" example:"picture"`
	RopepackMoonfestivalHPic02 string `json:"ropepack_moonfestival_h_pic_02" example:"picture"`
	RopepackMoonfestivalHPic03 string `json:"ropepack_moonfestival_h_pic_03" example:"picture"`
	RopepackMoonfestivalHPic04 string `json:"ropepack_moonfestival_h_pic_04" example:"picture"`
	RopepackMoonfestivalHPic05 string `json:"ropepack_moonfestival_h_pic_05" example:"picture"`
	RopepackMoonfestivalHPic06 string `json:"ropepack_moonfestival_h_pic_06" example:"picture"`
	RopepackMoonfestivalHPic07 string `json:"ropepack_moonfestival_h_pic_07" example:"picture"`
	RopepackMoonfestivalHPic08 string `json:"ropepack_moonfestival_h_pic_08" example:"picture"`
	RopepackMoonfestivalHPic09 string `json:"ropepack_moonfestival_h_pic_09" example:"picture"`
	RopepackMoonfestivalGPic01 string `json:"ropepack_moonfestival_g_pic_01" example:"picture"`
	RopepackMoonfestivalGPic02 string `json:"ropepack_moonfestival_g_pic_02" example:"picture"`
	RopepackMoonfestivalCPic01 string `json:"ropepack_moonfestival_c_pic_01" example:"picture"`
	RopepackMoonfestivalHAni01 string `json:"ropepack_moonfestival_h_ani_01" example:"picture"`
	RopepackMoonfestivalGAni01 string `json:"ropepack_moonfestival_g_ani_01" example:"picture"`
	RopepackMoonfestivalGAni02 string `json:"ropepack_moonfestival_g_ani_02" example:"picture"`
	RopepackMoonfestivalCAni01 string `json:"ropepack_moonfestival_c_ani_01" example:"picture"`
	RopepackMoonfestivalCAni02 string `json:"ropepack_moonfestival_c_ani_02" example:"picture"`

	Ropepack3DHPic01 string `json:"ropepack_3D_h_pic_01" example:"picture"`
	Ropepack3DHPic02 string `json:"ropepack_3D_h_pic_02" example:"picture"`
	Ropepack3DHPic03 string `json:"ropepack_3D_h_pic_03" example:"picture"`
	Ropepack3DHPic04 string `json:"ropepack_3D_h_pic_04" example:"picture"`
	Ropepack3DHPic05 string `json:"ropepack_3D_h_pic_05" example:"picture"`
	Ropepack3DHPic06 string `json:"ropepack_3D_h_pic_06" example:"picture"`
	Ropepack3DHPic07 string `json:"ropepack_3D_h_pic_07" example:"picture"`
	Ropepack3DHPic08 string `json:"ropepack_3D_h_pic_08" example:"picture"`
	Ropepack3DHPic09 string `json:"ropepack_3D_h_pic_09" example:"picture"`
	Ropepack3DHPic10 string `json:"ropepack_3D_h_pic_10" example:"picture"`
	Ropepack3DHPic11 string `json:"ropepack_3D_h_pic_11" example:"picture"`
	Ropepack3DHPic12 string `json:"ropepack_3D_h_pic_12" example:"picture"`
	Ropepack3DHPic13 string `json:"ropepack_3D_h_pic_13" example:"picture"`
	Ropepack3DHPic14 string `json:"ropepack_3D_h_pic_14" example:"picture"`
	Ropepack3DHPic15 string `json:"ropepack_3D_h_pic_15" example:"picture"`
	Ropepack3DGPic01 string `json:"ropepack_3D_g_pic_01" example:"picture"`
	Ropepack3DGPic02 string `json:"ropepack_3D_g_pic_02" example:"picture"`
	Ropepack3DGPic03 string `json:"ropepack_3D_g_pic_03" example:"picture"`
	Ropepack3DGPic04 string `json:"ropepack_3D_g_pic_04" example:"picture"`
	Ropepack3DHAni01 string `json:"ropepack_3D_h_ani_01" example:"picture"`
	Ropepack3DHAni02 string `json:"ropepack_3D_h_ani_02" example:"picture"`
	Ropepack3DHAni03 string `json:"ropepack_3D_h_ani_03" example:"picture"`
	Ropepack3DGAni01 string `json:"ropepack_3D_g_ani_01" example:"picture"`
	Ropepack3DGAni02 string `json:"ropepack_3D_g_ani_02" example:"picture"`
	Ropepack3DCAni01 string `json:"ropepack_3D_c_ani_01" example:"picture"`

	// 音樂
	RopepackBgmStart  string `json:"ropepack_bgm_start" example:"picture"`  // 遊戲開始
	RopepackBgmGaming string `json:"ropepack_bgm_gaming" example:"picture"` // 遊戲進行中
	RopepackBgmEnd    string `json:"ropepack_bgm_end" example:"picture"`    // 遊戲結束

	// 遊戲抽獎自定義
	LotteryJiugonggeClassicHPic01 string `json:"lottery_jiugongge_classic_h_pic_01" example:"picture"`
	LotteryJiugonggeClassicHPic02 string `json:"lottery_jiugongge_classic_h_pic_02" example:"picture"`
	LotteryJiugonggeClassicHPic03 string `json:"lottery_jiugongge_classic_h_pic_03" example:"picture"`
	LotteryJiugonggeClassicHPic04 string `json:"lottery_jiugongge_classic_h_pic_04" example:"picture"`
	LotteryJiugonggeClassicGPic01 string `json:"lottery_jiugongge_classic_g_pic_01" example:"picture"`
	LotteryJiugonggeClassicGPic02 string `json:"lottery_jiugongge_classic_g_pic_02" example:"picture"`
	LotteryJiugonggeClassicCPic01 string `json:"lottery_jiugongge_classic_c_pic_01" example:"picture"`
	LotteryJiugonggeClassicCPic02 string `json:"lottery_jiugongge_classic_c_pic_02" example:"picture"`
	LotteryJiugonggeClassicCPic03 string `json:"lottery_jiugongge_classic_c_pic_03" example:"picture"`
	LotteryJiugonggeClassicCPic04 string `json:"lottery_jiugongge_classic_c_pic_04" example:"picture"`
	LotteryJiugonggeClassicCAni01 string `json:"lottery_jiugongge_classic_c_ani_01" example:"picture"`
	LotteryJiugonggeClassicCAni02 string `json:"lottery_jiugongge_classic_c_ani_02" example:"picture"`
	LotteryJiugonggeClassicCAni03 string `json:"lottery_jiugongge_classic_c_ani_03" example:"picture"`

	LotteryTurntableClassicHPic01 string `json:"lottery_turntable_classic_h_pic_01" example:"picture"`
	LotteryTurntableClassicHPic02 string `json:"lottery_turntable_classic_h_pic_02" example:"picture"`
	LotteryTurntableClassicHPic03 string `json:"lottery_turntable_classic_h_pic_03" example:"picture"`
	LotteryTurntableClassicHPic04 string `json:"lottery_turntable_classic_h_pic_04" example:"picture"`
	LotteryTurntableClassicGPic01 string `json:"lottery_turntable_classic_g_pic_01" example:"picture"`
	LotteryTurntableClassicGPic02 string `json:"lottery_turntable_classic_g_pic_02" example:"picture"`
	LotteryTurntableClassicCPic01 string `json:"lottery_turntable_classic_c_pic_01" example:"picture"`
	LotteryTurntableClassicCPic02 string `json:"lottery_turntable_classic_c_pic_02" example:"picture"`
	LotteryTurntableClassicCPic03 string `json:"lottery_turntable_classic_c_pic_03" example:"picture"`
	LotteryTurntableClassicCPic04 string `json:"lottery_turntable_classic_c_pic_04" example:"picture"`
	LotteryTurntableClassicCPic05 string `json:"lottery_turntable_classic_c_pic_05" example:"picture"`
	LotteryTurntableClassicCPic06 string `json:"lottery_turntable_classic_c_pic_06" example:"picture"`
	LotteryTurntableClassicCAni01 string `json:"lottery_turntable_classic_c_ani_01" example:"picture"`
	LotteryTurntableClassicCAni02 string `json:"lottery_turntable_classic_c_ani_02" example:"picture"`
	LotteryTurntableClassicCAni03 string `json:"lottery_turntable_classic_c_ani_03" example:"picture"`

	LotteryJiugonggeStarryskyHPic01 string `json:"lottery_jiugongge_starrysky_h_pic_01" example:"picture"`
	LotteryJiugonggeStarryskyHPic02 string `json:"lottery_jiugongge_starrysky_h_pic_02" example:"picture"`
	LotteryJiugonggeStarryskyHPic03 string `json:"lottery_jiugongge_starrysky_h_pic_03" example:"picture"`
	LotteryJiugonggeStarryskyHPic04 string `json:"lottery_jiugongge_starrysky_h_pic_04" example:"picture"`
	LotteryJiugonggeStarryskyHPic05 string `json:"lottery_jiugongge_starrysky_h_pic_05" example:"picture"`
	LotteryJiugonggeStarryskyHPic06 string `json:"lottery_jiugongge_starrysky_h_pic_06" example:"picture"`
	LotteryJiugonggeStarryskyHPic07 string `json:"lottery_jiugongge_starrysky_h_pic_07" example:"picture"`
	LotteryJiugonggeStarryskyGPic01 string `json:"lottery_jiugongge_starrysky_g_pic_01" example:"picture"`
	LotteryJiugonggeStarryskyGPic02 string `json:"lottery_jiugongge_starrysky_g_pic_02" example:"picture"`
	LotteryJiugonggeStarryskyGPic03 string `json:"lottery_jiugongge_starrysky_g_pic_03" example:"picture"`
	LotteryJiugonggeStarryskyGPic04 string `json:"lottery_jiugongge_starrysky_g_pic_04" example:"picture"`
	LotteryJiugonggeStarryskyCPic01 string `json:"lottery_jiugongge_starrysky_c_pic_01" example:"picture"`
	LotteryJiugonggeStarryskyCPic02 string `json:"lottery_jiugongge_starrysky_c_pic_02" example:"picture"`
	LotteryJiugonggeStarryskyCPic03 string `json:"lottery_jiugongge_starrysky_c_pic_03" example:"picture"`
	LotteryJiugonggeStarryskyCPic04 string `json:"lottery_jiugongge_starrysky_c_pic_04" example:"picture"`
	LotteryJiugonggeStarryskyCAni01 string `json:"lottery_jiugongge_starrysky_c_ani_01" example:"picture"`
	LotteryJiugonggeStarryskyCAni02 string `json:"lottery_jiugongge_starrysky_c_ani_02" example:"picture"`
	LotteryJiugonggeStarryskyCAni03 string `json:"lottery_jiugongge_starrysky_c_ani_03" example:"picture"`
	LotteryJiugonggeStarryskyCAni04 string `json:"lottery_jiugongge_starrysky_c_ani_04" example:"picture"`
	LotteryJiugonggeStarryskyCAni05 string `json:"lottery_jiugongge_starrysky_c_ani_05" example:"picture"`
	LotteryJiugonggeStarryskyCAni06 string `json:"lottery_jiugongge_starrysky_c_ani_06" example:"picture"`

	LotteryTurntableStarryskyHPic01 string `json:"lottery_turntable_starrysky_h_pic_01" example:"picture"`
	LotteryTurntableStarryskyHPic02 string `json:"lottery_turntable_starrysky_h_pic_02" example:"picture"`
	LotteryTurntableStarryskyHPic03 string `json:"lottery_turntable_starrysky_h_pic_03" example:"picture"`
	LotteryTurntableStarryskyHPic04 string `json:"lottery_turntable_starrysky_h_pic_04" example:"picture"`
	LotteryTurntableStarryskyHPic05 string `json:"lottery_turntable_starrysky_h_pic_05" example:"picture"`
	LotteryTurntableStarryskyHPic06 string `json:"lottery_turntable_starrysky_h_pic_06" example:"picture"`
	LotteryTurntableStarryskyHPic07 string `json:"lottery_turntable_starrysky_h_pic_07" example:"picture"`
	LotteryTurntableStarryskyHPic08 string `json:"lottery_turntable_starrysky_h_pic_08" example:"picture"`
	LotteryTurntableStarryskyGPic01 string `json:"lottery_turntable_starrysky_g_pic_01" example:"picture"`
	LotteryTurntableStarryskyGPic02 string `json:"lottery_turntable_starrysky_g_pic_02" example:"picture"`
	LotteryTurntableStarryskyGPic03 string `json:"lottery_turntable_starrysky_g_pic_03" example:"picture"`
	LotteryTurntableStarryskyGPic04 string `json:"lottery_turntable_starrysky_g_pic_04" example:"picture"`
	LotteryTurntableStarryskyGPic05 string `json:"lottery_turntable_starrysky_g_pic_05" example:"picture"`
	LotteryTurntableStarryskyCPic01 string `json:"lottery_turntable_starrysky_c_pic_01" example:"picture"`
	LotteryTurntableStarryskyCPic02 string `json:"lottery_turntable_starrysky_c_pic_02" example:"picture"`
	LotteryTurntableStarryskyCPic03 string `json:"lottery_turntable_starrysky_c_pic_03" example:"picture"`
	LotteryTurntableStarryskyCPic04 string `json:"lottery_turntable_starrysky_c_pic_04" example:"picture"`
	LotteryTurntableStarryskyCPic05 string `json:"lottery_turntable_starrysky_c_pic_05" example:"picture"`
	LotteryTurntableStarryskyCAni01 string `json:"lottery_turntable_starrysky_c_ani_01" example:"picture"`
	LotteryTurntableStarryskyCAni02 string `json:"lottery_turntable_starrysky_c_ani_02" example:"picture"`
	LotteryTurntableStarryskyCAni03 string `json:"lottery_turntable_starrysky_c_ani_03" example:"picture"`
	LotteryTurntableStarryskyCAni04 string `json:"lottery_turntable_starrysky_c_ani_04" example:"picture"`
	LotteryTurntableStarryskyCAni05 string `json:"lottery_turntable_starrysky_c_ani_05" example:"picture"`
	LotteryTurntableStarryskyCAni06 string `json:"lottery_turntable_starrysky_c_ani_06" example:"picture"`
	LotteryTurntableStarryskyCAni07 string `json:"lottery_turntable_starrysky_c_ani_07" example:"picture"`

	// 音樂
	LotteryBgmGaming string `json:"lottery_bgm_gaming" example:"picture"` // 遊戲進行中

	// 鑑定師自定義
	MonopolyClassicHPic01 string `json:"monopoly_classic_h_pic_01" example:"picture"`
	MonopolyClassicHPic02 string `json:"monopoly_classic_h_pic_02" example:"picture"`
	MonopolyClassicHPic03 string `json:"monopoly_classic_h_pic_03" example:"picture"`
	MonopolyClassicHPic04 string `json:"monopoly_classic_h_pic_04" example:"picture"`
	MonopolyClassicHPic05 string `json:"monopoly_classic_h_pic_05" example:"picture"`
	MonopolyClassicHPic06 string `json:"monopoly_classic_h_pic_06" example:"picture"`
	MonopolyClassicHPic07 string `json:"monopoly_classic_h_pic_07" example:"picture"`
	MonopolyClassicHPic08 string `json:"monopoly_classic_h_pic_08" example:"picture"`
	MonopolyClassicGPic01 string `json:"monopoly_classic_g_pic_01" example:"picture"`
	MonopolyClassicGPic02 string `json:"monopoly_classic_g_pic_02" example:"picture"`
	MonopolyClassicGPic03 string `json:"monopoly_classic_g_pic_03" example:"picture"`
	MonopolyClassicGPic04 string `json:"monopoly_classic_g_pic_04" example:"picture"`
	MonopolyClassicGPic05 string `json:"monopoly_classic_g_pic_05" example:"picture"`
	MonopolyClassicGPic06 string `json:"monopoly_classic_g_pic_06" example:"picture"`
	MonopolyClassicGPic07 string `json:"monopoly_classic_g_pic_07" example:"picture"`
	MonopolyClassicCPic01 string `json:"monopoly_classic_c_pic_01" example:"picture"`
	MonopolyClassicCPic02 string `json:"monopoly_classic_c_pic_02" example:"picture"`
	MonopolyClassicGAni01 string `json:"monopoly_classic_g_ani_01" example:"picture"`
	MonopolyClassicGAni02 string `json:"monopoly_classic_g_ani_02" example:"picture"`
	MonopolyClassicCAni01 string `json:"monopoly_classic_c_ani_01" example:"picture"`

	MonopolyRedpackHPic01 string `json:"monopoly_redpack_h_pic_01" example:"picture"`
	MonopolyRedpackHPic02 string `json:"monopoly_redpack_h_pic_02" example:"picture"`
	MonopolyRedpackHPic03 string `json:"monopoly_redpack_h_pic_03" example:"picture"`
	MonopolyRedpackHPic04 string `json:"monopoly_redpack_h_pic_04" example:"picture"`
	MonopolyRedpackHPic05 string `json:"monopoly_redpack_h_pic_05" example:"picture"`
	MonopolyRedpackHPic06 string `json:"monopoly_redpack_h_pic_06" example:"picture"`
	MonopolyRedpackHPic07 string `json:"monopoly_redpack_h_pic_07" example:"picture"`
	MonopolyRedpackHPic08 string `json:"monopoly_redpack_h_pic_08" example:"picture"`
	MonopolyRedpackHPic09 string `json:"monopoly_redpack_h_pic_09" example:"picture"`
	MonopolyRedpackHPic10 string `json:"monopoly_redpack_h_pic_10" example:"picture"`
	MonopolyRedpackHPic11 string `json:"monopoly_redpack_h_pic_11" example:"picture"`
	MonopolyRedpackHPic12 string `json:"monopoly_redpack_h_pic_12" example:"picture"`
	MonopolyRedpackHPic13 string `json:"monopoly_redpack_h_pic_13" example:"picture"`
	MonopolyRedpackHPic14 string `json:"monopoly_redpack_h_pic_14" example:"picture"`
	MonopolyRedpackHPic15 string `json:"monopoly_redpack_h_pic_15" example:"picture"`
	MonopolyRedpackHPic16 string `json:"monopoly_redpack_h_pic_16" example:"picture"`
	MonopolyRedpackGPic01 string `json:"monopoly_redpack_g_pic_01" example:"picture"`
	MonopolyRedpackGPic02 string `json:"monopoly_redpack_g_pic_02" example:"picture"`
	MonopolyRedpackGPic03 string `json:"monopoly_redpack_g_pic_03" example:"picture"`
	MonopolyRedpackGPic04 string `json:"monopoly_redpack_g_pic_04" example:"picture"`
	MonopolyRedpackGPic05 string `json:"monopoly_redpack_g_pic_05" example:"picture"`
	MonopolyRedpackGPic06 string `json:"monopoly_redpack_g_pic_06" example:"picture"`
	MonopolyRedpackGPic07 string `json:"monopoly_redpack_g_pic_07" example:"picture"`
	MonopolyRedpackGPic08 string `json:"monopoly_redpack_g_pic_08" example:"picture"`
	MonopolyRedpackGPic09 string `json:"monopoly_redpack_g_pic_09" example:"picture"`
	MonopolyRedpackGPic10 string `json:"monopoly_redpack_g_pic_10" example:"picture"`
	MonopolyRedpackCPic01 string `json:"monopoly_redpack_c_pic_01" example:"picture"`
	MonopolyRedpackCPic02 string `json:"monopoly_redpack_c_pic_02" example:"picture"`
	MonopolyRedpackCPic03 string `json:"monopoly_redpack_c_pic_03" example:"picture"`
	MonopolyRedpackHAni01 string `json:"monopoly_redpack_h_ani_01" example:"picture"`
	MonopolyRedpackHAni02 string `json:"monopoly_redpack_h_ani_02" example:"picture"`
	MonopolyRedpackHAni03 string `json:"monopoly_redpack_h_ani_03" example:"picture"`
	MonopolyRedpackGAni01 string `json:"monopoly_redpack_g_ani_01" example:"picture"`
	MonopolyRedpackGAni02 string `json:"monopoly_redpack_g_ani_02" example:"picture"`
	MonopolyRedpackCAni01 string `json:"monopoly_redpack_c_ani_01" example:"picture"`

	MonopolyNewyearRabbitHPic01 string `json:"monopoly_newyear_rabbit_h_pic_01" example:"picture"`
	MonopolyNewyearRabbitHPic02 string `json:"monopoly_newyear_rabbit_h_pic_02" example:"picture"`
	MonopolyNewyearRabbitHPic03 string `json:"monopoly_newyear_rabbit_h_pic_03" example:"picture"`
	MonopolyNewyearRabbitHPic04 string `json:"monopoly_newyear_rabbit_h_pic_04" example:"picture"`
	MonopolyNewyearRabbitHPic05 string `json:"monopoly_newyear_rabbit_h_pic_05" example:"picture"`
	MonopolyNewyearRabbitHPic06 string `json:"monopoly_newyear_rabbit_h_pic_06" example:"picture"`
	MonopolyNewyearRabbitHPic07 string `json:"monopoly_newyear_rabbit_h_pic_07" example:"picture"`
	MonopolyNewyearRabbitHPic08 string `json:"monopoly_newyear_rabbit_h_pic_08" example:"picture"`
	MonopolyNewyearRabbitHPic09 string `json:"monopoly_newyear_rabbit_h_pic_09" example:"picture"`
	MonopolyNewyearRabbitHPic10 string `json:"monopoly_newyear_rabbit_h_pic_10" example:"picture"`
	MonopolyNewyearRabbitHPic11 string `json:"monopoly_newyear_rabbit_h_pic_11" example:"picture"`
	MonopolyNewyearRabbitHPic12 string `json:"monopoly_newyear_rabbit_h_pic_12" example:"picture"`
	MonopolyNewyearRabbitGPic01 string `json:"monopoly_newyear_rabbit_g_pic_01" example:"picture"`
	MonopolyNewyearRabbitGPic02 string `json:"monopoly_newyear_rabbit_g_pic_02" example:"picture"`
	MonopolyNewyearRabbitGPic03 string `json:"monopoly_newyear_rabbit_g_pic_03" example:"picture"`
	MonopolyNewyearRabbitGPic04 string `json:"monopoly_newyear_rabbit_g_pic_04" example:"picture"`
	MonopolyNewyearRabbitGPic05 string `json:"monopoly_newyear_rabbit_g_pic_05" example:"picture"`
	MonopolyNewyearRabbitGPic06 string `json:"monopoly_newyear_rabbit_g_pic_06" example:"picture"`
	MonopolyNewyearRabbitGPic07 string `json:"monopoly_newyear_rabbit_g_pic_07" example:"picture"`
	MonopolyNewyearRabbitCPic01 string `json:"monopoly_newyear_rabbit_c_pic_01" example:"picture"`
	MonopolyNewyearRabbitCPic02 string `json:"monopoly_newyear_rabbit_c_pic_02" example:"picture"`
	MonopolyNewyearRabbitCPic03 string `json:"monopoly_newyear_rabbit_c_pic_03" example:"picture"`
	MonopolyNewyearRabbitHAni01 string `json:"monopoly_newyear_rabbit_h_ani_01" example:"picture"`
	MonopolyNewyearRabbitHAni02 string `json:"monopoly_newyear_rabbit_h_ani_02" example:"picture"`
	MonopolyNewyearRabbitGAni01 string `json:"monopoly_newyear_rabbit_g_ani_01" example:"picture"`
	MonopolyNewyearRabbitGAni02 string `json:"monopoly_newyear_rabbit_g_ani_02" example:"picture"`
	MonopolyNewyearRabbitCAni01 string `json:"monopoly_newyear_rabbit_c_ani_01" example:"picture"`

	MonopolySashimiHPic01 string `json:"monopoly_sashimi_h_pic_01" example:"picture"`
	MonopolySashimiHPic02 string `json:"monopoly_sashimi_h_pic_02" example:"picture"`
	MonopolySashimiHPic03 string `json:"monopoly_sashimi_h_pic_03" example:"picture"`
	MonopolySashimiHPic04 string `json:"monopoly_sashimi_h_pic_04" example:"picture"`
	MonopolySashimiHPic05 string `json:"monopoly_sashimi_h_pic_05" example:"picture"`
	MonopolySashimiGPic01 string `json:"monopoly_sashimi_g_pic_01" example:"picture"`
	MonopolySashimiGPic02 string `json:"monopoly_sashimi_g_pic_02" example:"picture"`
	MonopolySashimiGPic03 string `json:"monopoly_sashimi_g_pic_03" example:"picture"`
	MonopolySashimiGPic04 string `json:"monopoly_sashimi_g_pic_04" example:"picture"`
	MonopolySashimiGPic05 string `json:"monopoly_sashimi_g_pic_05" example:"picture"`
	MonopolySashimiGPic06 string `json:"monopoly_sashimi_g_pic_06" example:"picture"`
	MonopolySashimiCPic01 string `json:"monopoly_sashimi_c_pic_01" example:"picture"`
	MonopolySashimiCPic02 string `json:"monopoly_sashimi_c_pic_02" example:"picture"`
	MonopolySashimiHAni01 string `json:"monopoly_sashimi_h_ani_01" example:"picture"`
	MonopolySashimiHAni02 string `json:"monopoly_sashimi_h_ani_02" example:"picture"`
	MonopolySashimiGAni01 string `json:"monopoly_sashimi_g_ani_01" example:"picture"`
	MonopolySashimiGAni02 string `json:"monopoly_sashimi_g_ani_02" example:"picture"`
	MonopolySashimiCAni01 string `json:"monopoly_sashimi_c_ani_01" example:"picture"`

	// 音樂
	MonopolyBgmStart  string `json:"monopoly_bgm_start" example:"picture"`  // 遊戲開始
	MonopolyBgmGaming string `json:"monopoly_bgm_gaming" example:"picture"` // 遊戲進行中
	MonopolyBgmEnd    string `json:"monopoly_bgm_end" example:"picture"`    // 遊戲結束

	// 拔河遊戲自定義
	TugofwarClassicHPic01 string `json:"tugofwar_classic_h_pic_01" example:"picture"`
	TugofwarClassicHPic02 string `json:"tugofwar_classic_h_pic_02" example:"picture"`
	TugofwarClassicHPic03 string `json:"tugofwar_classic_h_pic_03" example:"picture"`
	TugofwarClassicHPic04 string `json:"tugofwar_classic_h_pic_04" example:"picture"`
	TugofwarClassicHPic05 string `json:"tugofwar_classic_h_pic_05" example:"picture"`
	TugofwarClassicHPic06 string `json:"tugofwar_classic_h_pic_06" example:"picture"`
	TugofwarClassicHPic07 string `json:"tugofwar_classic_h_pic_07" example:"picture"`
	TugofwarClassicHPic08 string `json:"tugofwar_classic_h_pic_08" example:"picture"`
	TugofwarClassicHPic09 string `json:"tugofwar_classic_h_pic_09" example:"picture"`
	TugofwarClassicHPic10 string `json:"tugofwar_classic_h_pic_10" example:"picture"`
	TugofwarClassicHPic11 string `json:"tugofwar_classic_h_pic_11" example:"picture"`
	TugofwarClassicHPic12 string `json:"tugofwar_classic_h_pic_12" example:"picture"`
	TugofwarClassicHPic13 string `json:"tugofwar_classic_h_pic_13" example:"picture"`
	TugofwarClassicHPic14 string `json:"tugofwar_classic_h_pic_14" example:"picture"`
	TugofwarClassicHPic15 string `json:"tugofwar_classic_h_pic_15" example:"picture"`
	TugofwarClassicHPic16 string `json:"tugofwar_classic_h_pic_16" example:"picture"`
	TugofwarClassicHPic17 string `json:"tugofwar_classic_h_pic_17" example:"picture"`
	TugofwarClassicHPic18 string `json:"tugofwar_classic_h_pic_18" example:"picture"`
	TugofwarClassicHPic19 string `json:"tugofwar_classic_h_pic_19" example:"picture"`
	TugofwarClassicGPic01 string `json:"tugofwar_classic_g_pic_01" example:"picture"`
	TugofwarClassicGPic02 string `json:"tugofwar_classic_g_pic_02" example:"picture"`
	TugofwarClassicGPic03 string `json:"tugofwar_classic_g_pic_03" example:"picture"`
	TugofwarClassicGPic04 string `json:"tugofwar_classic_g_pic_04" example:"picture"`
	TugofwarClassicGPic05 string `json:"tugofwar_classic_g_pic_05" example:"picture"`
	TugofwarClassicGPic06 string `json:"tugofwar_classic_g_pic_06" example:"picture"`
	TugofwarClassicGPic07 string `json:"tugofwar_classic_g_pic_07" example:"picture"`
	TugofwarClassicGPic08 string `json:"tugofwar_classic_g_pic_08" example:"picture"`
	TugofwarClassicGPic09 string `json:"tugofwar_classic_g_pic_09" example:"picture"`
	TugofwarClassicHAni01 string `json:"tugofwar_classic_h_ani_01" example:"picture"`
	TugofwarClassicHAni02 string `json:"tugofwar_classic_h_ani_02" example:"picture"`
	TugofwarClassicHAni03 string `json:"tugofwar_classic_h_ani_03" example:"picture"`
	TugofwarClassicCAni01 string `json:"tugofwar_classic_c_ani_01" example:"picture"`

	TugofwarSchoolHPic01 string `json:"tugofwar_school_h_pic_01" example:"picture"`
	TugofwarSchoolHPic02 string `json:"tugofwar_school_h_pic_02" example:"picture"`
	TugofwarSchoolHPic03 string `json:"tugofwar_school_h_pic_03" example:"picture"`
	TugofwarSchoolHPic04 string `json:"tugofwar_school_h_pic_04" example:"picture"`
	TugofwarSchoolHPic05 string `json:"tugofwar_school_h_pic_05" example:"picture"`
	TugofwarSchoolHPic06 string `json:"tugofwar_school_h_pic_06" example:"picture"`
	TugofwarSchoolHPic07 string `json:"tugofwar_school_h_pic_07" example:"picture"`
	TugofwarSchoolHPic08 string `json:"tugofwar_school_h_pic_08" example:"picture"`
	TugofwarSchoolHPic09 string `json:"tugofwar_school_h_pic_09" example:"picture"`
	TugofwarSchoolHPic10 string `json:"tugofwar_school_h_pic_10" example:"picture"`
	TugofwarSchoolHPic11 string `json:"tugofwar_school_h_pic_11" example:"picture"`
	TugofwarSchoolHPic12 string `json:"tugofwar_school_h_pic_12" example:"picture"`
	TugofwarSchoolHPic13 string `json:"tugofwar_school_h_pic_13" example:"picture"`
	TugofwarSchoolHPic14 string `json:"tugofwar_school_h_pic_14" example:"picture"`
	TugofwarSchoolHPic15 string `json:"tugofwar_school_h_pic_15" example:"picture"`
	TugofwarSchoolHPic16 string `json:"tugofwar_school_h_pic_16" example:"picture"`
	TugofwarSchoolHPic17 string `json:"tugofwar_school_h_pic_17" example:"picture"`
	TugofwarSchoolHPic18 string `json:"tugofwar_school_h_pic_18" example:"picture"`
	TugofwarSchoolHPic19 string `json:"tugofwar_school_h_pic_19" example:"picture"`
	TugofwarSchoolHPic20 string `json:"tugofwar_school_h_pic_20" example:"picture"`
	TugofwarSchoolHPic21 string `json:"tugofwar_school_h_pic_21" example:"picture"`
	TugofwarSchoolHPic22 string `json:"tugofwar_school_h_pic_22" example:"picture"`
	TugofwarSchoolHPic23 string `json:"tugofwar_school_h_pic_23" example:"picture"`
	TugofwarSchoolHPic24 string `json:"tugofwar_school_h_pic_24" example:"picture"`
	TugofwarSchoolHPic25 string `json:"tugofwar_school_h_pic_25" example:"picture"`
	TugofwarSchoolHPic26 string `json:"tugofwar_school_h_pic_26" example:"picture"`
	TugofwarSchoolGPic01 string `json:"tugofwar_school_g_pic_01" example:"picture"`
	TugofwarSchoolGPic02 string `json:"tugofwar_school_g_pic_02" example:"picture"`
	TugofwarSchoolGPic03 string `json:"tugofwar_school_g_pic_03" example:"picture"`
	TugofwarSchoolGPic04 string `json:"tugofwar_school_g_pic_04" example:"picture"`
	TugofwarSchoolGPic05 string `json:"tugofwar_school_g_pic_05" example:"picture"`
	TugofwarSchoolGPic06 string `json:"tugofwar_school_g_pic_06" example:"picture"`
	TugofwarSchoolGPic07 string `json:"tugofwar_school_g_pic_07" example:"picture"`
	TugofwarSchoolCPic01 string `json:"tugofwar_school_c_pic_01" example:"picture"`
	TugofwarSchoolCPic02 string `json:"tugofwar_school_c_pic_02" example:"picture"`
	TugofwarSchoolCPic03 string `json:"tugofwar_school_c_pic_03" example:"picture"`
	TugofwarSchoolCPic04 string `json:"tugofwar_school_c_pic_04" example:"picture"`
	TugofwarSchoolHAni01 string `json:"tugofwar_school_h_ani_01" example:"picture"`
	TugofwarSchoolHAni02 string `json:"tugofwar_school_h_ani_02" example:"picture"`
	TugofwarSchoolHAni03 string `json:"tugofwar_school_h_ani_03" example:"picture"`
	TugofwarSchoolHAni04 string `json:"tugofwar_school_h_ani_04" example:"picture"`
	TugofwarSchoolHAni05 string `json:"tugofwar_school_h_ani_05" example:"picture"`
	TugofwarSchoolHAni06 string `json:"tugofwar_school_h_ani_06" example:"picture"`
	TugofwarSchoolHAni07 string `json:"tugofwar_school_h_ani_07" example:"picture"`

	TugofwarChristmasHPic01 string `json:"tugofwar_christmas_h_pic_01" example:"picture"`
	TugofwarChristmasHPic02 string `json:"tugofwar_christmas_h_pic_02" example:"picture"`
	TugofwarChristmasHPic03 string `json:"tugofwar_christmas_h_pic_03" example:"picture"`
	TugofwarChristmasHPic04 string `json:"tugofwar_christmas_h_pic_04" example:"picture"`
	TugofwarChristmasHPic05 string `json:"tugofwar_christmas_h_pic_05" example:"picture"`
	TugofwarChristmasHPic06 string `json:"tugofwar_christmas_h_pic_06" example:"picture"`
	TugofwarChristmasHPic07 string `json:"tugofwar_christmas_h_pic_07" example:"picture"`
	TugofwarChristmasHPic08 string `json:"tugofwar_christmas_h_pic_08" example:"picture"`
	TugofwarChristmasHPic09 string `json:"tugofwar_christmas_h_pic_09" example:"picture"`
	TugofwarChristmasHPic10 string `json:"tugofwar_christmas_h_pic_10" example:"picture"`
	TugofwarChristmasHPic11 string `json:"tugofwar_christmas_h_pic_11" example:"picture"`
	TugofwarChristmasHPic12 string `json:"tugofwar_christmas_h_pic_12" example:"picture"`
	TugofwarChristmasHPic13 string `json:"tugofwar_christmas_h_pic_13" example:"picture"`
	TugofwarChristmasHPic14 string `json:"tugofwar_christmas_h_pic_14" example:"picture"`
	TugofwarChristmasHPic15 string `json:"tugofwar_christmas_h_pic_15" example:"picture"`
	TugofwarChristmasHPic16 string `json:"tugofwar_christmas_h_pic_16" example:"picture"`
	TugofwarChristmasHPic17 string `json:"tugofwar_christmas_h_pic_17" example:"picture"`
	TugofwarChristmasHPic18 string `json:"tugofwar_christmas_h_pic_18" example:"picture"`
	TugofwarChristmasHPic19 string `json:"tugofwar_christmas_h_pic_19" example:"picture"`
	TugofwarChristmasHPic20 string `json:"tugofwar_christmas_h_pic_20" example:"picture"`
	TugofwarChristmasHPic21 string `json:"tugofwar_christmas_h_pic_21" example:"picture"`
	TugofwarChristmasGPic01 string `json:"tugofwar_christmas_g_pic_01" example:"picture"`
	TugofwarChristmasGPic02 string `json:"tugofwar_christmas_g_pic_02" example:"picture"`
	TugofwarChristmasGPic03 string `json:"tugofwar_christmas_g_pic_03" example:"picture"`
	TugofwarChristmasGPic04 string `json:"tugofwar_christmas_g_pic_04" example:"picture"`
	TugofwarChristmasGPic05 string `json:"tugofwar_christmas_g_pic_05" example:"picture"`
	TugofwarChristmasGPic06 string `json:"tugofwar_christmas_g_pic_06" example:"picture"`
	TugofwarChristmasCPic01 string `json:"tugofwar_christmas_c_pic_01" example:"picture"`
	TugofwarChristmasCPic02 string `json:"tugofwar_christmas_c_pic_02" example:"picture"`
	TugofwarChristmasCPic03 string `json:"tugofwar_christmas_c_pic_03" example:"picture"`
	TugofwarChristmasCPic04 string `json:"tugofwar_christmas_c_pic_04" example:"picture"`
	TugofwarChristmasHAni01 string `json:"tugofwar_christmas_h_ani_01" example:"picture"`
	TugofwarChristmasHAni02 string `json:"tugofwar_christmas_h_ani_02" example:"picture"`
	TugofwarChristmasHAni03 string `json:"tugofwar_christmas_h_ani_03" example:"picture"`
	TugofwarChristmasCAni01 string `json:"tugofwar_christmas_c_ani_01" example:"picture"`
	TugofwarChristmasCAni02 string `json:"tugofwar_christmas_c_ani_02" example:"picture"`

	// 音樂
	TugofwarBgmStart  string `json:"tugofwar_bgm_start" example:"picture"`  // 遊戲開始
	TugofwarBgmGaming string `json:"tugofwar_bgm_gaming" example:"picture"` // 遊戲進行中
	TugofwarBgmEnd    string `json:"tugofwar_bgm_end" example:"picture"`    // 遊戲結束

	// 賓果遊戲自定義
	BingoClassicHPic01 string `json:"bingo_classic_h_pic_01" example:"picture"`
	BingoClassicHPic02 string `json:"bingo_classic_h_pic_02" example:"picture"`
	BingoClassicHPic03 string `json:"bingo_classic_h_pic_03" example:"picture"`
	BingoClassicHPic04 string `json:"bingo_classic_h_pic_04" example:"picture"`
	BingoClassicHPic05 string `json:"bingo_classic_h_pic_05" example:"picture"`
	BingoClassicHPic06 string `json:"bingo_classic_h_pic_06" example:"picture"`
	BingoClassicHPic07 string `json:"bingo_classic_h_pic_07" example:"picture"`
	BingoClassicHPic08 string `json:"bingo_classic_h_pic_08" example:"picture"`
	BingoClassicHPic09 string `json:"bingo_classic_h_pic_09" example:"picture"`
	BingoClassicHPic10 string `json:"bingo_classic_h_pic_10" example:"picture"`
	BingoClassicHPic11 string `json:"bingo_classic_h_pic_11" example:"picture"`
	BingoClassicHPic12 string `json:"bingo_classic_h_pic_12" example:"picture"`
	BingoClassicHPic13 string `json:"bingo_classic_h_pic_13" example:"picture"`
	BingoClassicHPic14 string `json:"bingo_classic_h_pic_14" example:"picture"`
	BingoClassicHPic15 string `json:"bingo_classic_h_pic_15" example:"picture"`
	BingoClassicHPic16 string `json:"bingo_classic_h_pic_16" example:"picture"`
	BingoClassicGPic01 string `json:"bingo_classic_g_pic_01" example:"picture"`
	BingoClassicGPic02 string `json:"bingo_classic_g_pic_02" example:"picture"`
	BingoClassicGPic03 string `json:"bingo_classic_g_pic_03" example:"picture"`
	BingoClassicGPic04 string `json:"bingo_classic_g_pic_04" example:"picture"`
	BingoClassicGPic05 string `json:"bingo_classic_g_pic_05" example:"picture"`
	BingoClassicGPic06 string `json:"bingo_classic_g_pic_06" example:"picture"`
	BingoClassicGPic07 string `json:"bingo_classic_g_pic_07" example:"picture"`
	BingoClassicGPic08 string `json:"bingo_classic_g_pic_08" example:"picture"`
	BingoClassicCPic01 string `json:"bingo_classic_c_pic_01" example:"picture"`
	BingoClassicCPic02 string `json:"bingo_classic_c_pic_02" example:"picture"`
	BingoClassicCPic03 string `json:"bingo_classic_c_pic_03" example:"picture"`
	BingoClassicCPic04 string `json:"bingo_classic_c_pic_04" example:"picture"`
	// BingoClassicCPic05 string `json:"bingo_classic_c_pic_05" example:"picture"`
	BingoClassicHAni01 string `json:"bingo_classic_h_ani_01" example:"picture"`
	BingoClassicGAni01 string `json:"bingo_classic_g_ani_01" example:"picture"`
	BingoClassicCAni01 string `json:"bingo_classic_c_ani_01" example:"picture"`
	BingoClassicCAni02 string `json:"bingo_classic_c_ani_02" example:"picture"`

	BingoNewyearDragonHPic01 string `json:"bingo_newyear_dragon_h_pic_01" example:"picture"`
	BingoNewyearDragonHPic02 string `json:"bingo_newyear_dragon_h_pic_02" example:"picture"`
	BingoNewyearDragonHPic03 string `json:"bingo_newyear_dragon_h_pic_03" example:"picture"`
	BingoNewyearDragonHPic04 string `json:"bingo_newyear_dragon_h_pic_04" example:"picture"`
	BingoNewyearDragonHPic05 string `json:"bingo_newyear_dragon_h_pic_05" example:"picture"`
	BingoNewyearDragonHPic06 string `json:"bingo_newyear_dragon_h_pic_06" example:"picture"`
	BingoNewyearDragonHPic07 string `json:"bingo_newyear_dragon_h_pic_07" example:"picture"`
	BingoNewyearDragonHPic08 string `json:"bingo_newyear_dragon_h_pic_08" example:"picture"`
	BingoNewyearDragonHPic09 string `json:"bingo_newyear_dragon_h_pic_09" example:"picture"`
	BingoNewyearDragonHPic10 string `json:"bingo_newyear_dragon_h_pic_10" example:"picture"`
	BingoNewyearDragonHPic11 string `json:"bingo_newyear_dragon_h_pic_11" example:"picture"`
	BingoNewyearDragonHPic12 string `json:"bingo_newyear_dragon_h_pic_12" example:"picture"`
	BingoNewyearDragonHPic13 string `json:"bingo_newyear_dragon_h_pic_13" example:"picture"`
	BingoNewyearDragonHPic14 string `json:"bingo_newyear_dragon_h_pic_14" example:"picture"`
	// BingoNewyearDragonHPic15 string `json:"bingo_newyear_dragon_h_pic_15" example:"picture"`
	BingoNewyearDragonHPic16 string `json:"bingo_newyear_dragon_h_pic_16" example:"picture"`
	BingoNewyearDragonHPic17 string `json:"bingo_newyear_dragon_h_pic_17" example:"picture"`
	BingoNewyearDragonHPic18 string `json:"bingo_newyear_dragon_h_pic_18" example:"picture"`
	BingoNewyearDragonHPic19 string `json:"bingo_newyear_dragon_h_pic_19" example:"picture"`
	BingoNewyearDragonHPic20 string `json:"bingo_newyear_dragon_h_pic_20" example:"picture"`
	BingoNewyearDragonHPic21 string `json:"bingo_newyear_dragon_h_pic_21" example:"picture"`
	BingoNewyearDragonHPic22 string `json:"bingo_newyear_dragon_h_pic_22" example:"picture"`
	BingoNewyearDragonGPic01 string `json:"bingo_newyear_dragon_g_pic_01" example:"picture"`
	BingoNewyearDragonGPic02 string `json:"bingo_newyear_dragon_g_pic_02" example:"picture"`
	BingoNewyearDragonGPic03 string `json:"bingo_newyear_dragon_g_pic_03" example:"picture"`
	BingoNewyearDragonGPic04 string `json:"bingo_newyear_dragon_g_pic_04" example:"picture"`
	BingoNewyearDragonGPic05 string `json:"bingo_newyear_dragon_g_pic_05" example:"picture"`
	BingoNewyearDragonGPic06 string `json:"bingo_newyear_dragon_g_pic_06" example:"picture"`
	BingoNewyearDragonGPic07 string `json:"bingo_newyear_dragon_g_pic_07" example:"picture"`
	BingoNewyearDragonGPic08 string `json:"bingo_newyear_dragon_g_pic_08" example:"picture"`
	BingoNewyearDragonCPic01 string `json:"bingo_newyear_dragon_c_pic_01" example:"picture"`
	BingoNewyearDragonCPic02 string `json:"bingo_newyear_dragon_c_pic_02" example:"picture"`
	BingoNewyearDragonCPic03 string `json:"bingo_newyear_dragon_c_pic_03" example:"picture"`
	BingoNewyearDragonHAni01 string `json:"bingo_newyear_dragon_h_ani_01" example:"picture"`
	BingoNewyearDragonGAni01 string `json:"bingo_newyear_dragon_g_ani_01" example:"picture"`
	BingoNewyearDragonCAni01 string `json:"bingo_newyear_dragon_c_ani_01" example:"picture"`
	BingoNewyearDragonCAni02 string `json:"bingo_newyear_dragon_c_ani_02" example:"picture"`
	BingoNewyearDragonCAni03 string `json:"bingo_newyear_dragon_c_ani_03" example:"picture"`

	BingoCherryHPic01 string `json:"bingo_cherry_h_pic_01" example:"picture"`
	BingoCherryHPic02 string `json:"bingo_cherry_h_pic_02" example:"picture"`
	BingoCherryHPic03 string `json:"bingo_cherry_h_pic_03" example:"picture"`
	BingoCherryHPic04 string `json:"bingo_cherry_h_pic_04" example:"picture"`
	BingoCherryHPic05 string `json:"bingo_cherry_h_pic_05" example:"picture"`
	BingoCherryHPic06 string `json:"bingo_cherry_h_pic_06" example:"picture"`
	BingoCherryHPic07 string `json:"bingo_cherry_h_pic_07" example:"picture"`
	BingoCherryHPic08 string `json:"bingo_cherry_h_pic_08" example:"picture"`
	BingoCherryHPic09 string `json:"bingo_cherry_h_pic_09" example:"picture"`
	BingoCherryHPic10 string `json:"bingo_cherry_h_pic_10" example:"picture"`
	BingoCherryHPic11 string `json:"bingo_cherry_h_pic_11" example:"picture"`
	BingoCherryHPic12 string `json:"bingo_cherry_h_pic_12" example:"picture"`
	// BingoCherryHPic13 string `json:"bingo_cherry_h_pic_13" example:"picture"`
	BingoCherryHPic14 string `json:"bingo_cherry_h_pic_14" example:"picture"`
	BingoCherryHPic15 string `json:"bingo_cherry_h_pic_15" example:"picture"`
	// BingoCherryHPic16 string `json:"bingo_cherry_h_pic_16" example:"picture"`
	BingoCherryHPic17 string `json:"bingo_cherry_h_pic_17" example:"picture"`
	BingoCherryHPic18 string `json:"bingo_cherry_h_pic_18" example:"picture"`
	BingoCherryHPic19 string `json:"bingo_cherry_h_pic_19" example:"picture"`
	BingoCherryGPic01 string `json:"bingo_cherry_g_pic_01" example:"picture"`
	BingoCherryGPic02 string `json:"bingo_cherry_g_pic_02" example:"picture"`
	BingoCherryGPic03 string `json:"bingo_cherry_g_pic_03" example:"picture"`
	BingoCherryGPic04 string `json:"bingo_cherry_g_pic_04" example:"picture"`
	BingoCherryGPic05 string `json:"bingo_cherry_g_pic_05" example:"picture"`
	BingoCherryGPic06 string `json:"bingo_cherry_g_pic_06" example:"picture"`
	BingoCherryCPic01 string `json:"bingo_cherry_c_pic_01" example:"picture"`
	BingoCherryCPic02 string `json:"bingo_cherry_c_pic_02" example:"picture"`
	BingoCherryCPic03 string `json:"bingo_cherry_c_pic_03" example:"picture"`
	BingoCherryCPic04 string `json:"bingo_cherry_c_pic_04" example:"picture"`
	// BingoCherryHAni01 string `json:"bingo_cherry_h_ani_01" example:"picture"`
	BingoCherryHAni02 string `json:"bingo_cherry_h_ani_02" example:"picture"`
	BingoCherryHAni03 string `json:"bingo_cherry_h_ani_03" example:"picture"`
	BingoCherryGAni01 string `json:"bingo_cherry_g_ani_01" example:"picture"`
	BingoCherryGAni02 string `json:"bingo_cherry_g_ani_02" example:"picture"`
	BingoCherryCAni01 string `json:"bingo_cherry_c_ani_01" example:"picture"`
	BingoCherryCAni02 string `json:"bingo_cherry_c_ani_02" example:"picture"`

	// 音樂
	BingoBgmStart  string `json:"bingo_bgm_start" example:"picture"`  // 遊戲開始
	BingoBgmGaming string `json:"bingo_bgm_gaming" example:"picture"` // 遊戲進行中
	BingoBgmEnd    string `json:"bingo_bgm_end" example:"picture"`    // 遊戲結束

	// 扭蛋機自定義
	GachaMachineClassicHPic02 string `json:"3d_gacha_machine_classic_h_pic_02" example:"picture"`
	GachaMachineClassicHPic03 string `json:"3d_gacha_machine_classic_h_pic_03" example:"picture"`
	GachaMachineClassicHPic04 string `json:"3d_gacha_machine_classic_h_pic_04" example:"picture"`
	GachaMachineClassicHPic05 string `json:"3d_gacha_machine_classic_h_pic_05" example:"picture"`
	GachaMachineClassicGPic01 string `json:"3d_gacha_machine_classic_g_pic_01" example:"picture"`
	GachaMachineClassicGPic02 string `json:"3d_gacha_machine_classic_g_pic_02" example:"picture"`
	GachaMachineClassicGPic03 string `json:"3d_gacha_machine_classic_g_pic_03" example:"picture"`
	GachaMachineClassicGPic04 string `json:"3d_gacha_machine_classic_g_pic_04" example:"picture"`
	GachaMachineClassicGPic05 string `json:"3d_gacha_machine_classic_g_pic_05" example:"picture"`
	GachaMachineClassicCPic01 string `json:"3d_gacha_machine_classic_c_pic_01" example:"picture"`

	// 音樂
	GachaMachineBgmGaming string `json:"3d_gacha_machine_bgm_gaming" example:"picture"`

	// 投票自定義
	VoteClassicHPic01 string `json:"vote_classic_h_pic_01" example:"picture"`
	VoteClassicHPic02 string `json:"vote_classic_h_pic_02" example:"picture"`
	VoteClassicHPic03 string `json:"vote_classic_h_pic_03" example:"picture"`
	VoteClassicHPic04 string `json:"vote_classic_h_pic_04" example:"picture"`
	VoteClassicHPic05 string `json:"vote_classic_h_pic_05" example:"picture"`
	VoteClassicHPic06 string `json:"vote_classic_h_pic_06" example:"picture"`
	VoteClassicHPic07 string `json:"vote_classic_h_pic_07" example:"picture"`
	VoteClassicHPic08 string `json:"vote_classic_h_pic_08" example:"picture"`
	VoteClassicHPic09 string `json:"vote_classic_h_pic_09" example:"picture"`
	VoteClassicHPic10 string `json:"vote_classic_h_pic_10" example:"picture"`
	VoteClassicHPic11 string `json:"vote_classic_h_pic_11" example:"picture"`
	// VoteClassicHPic12 string `json:"vote_classic_h_pic_12" example:"picture"`
	VoteClassicHPic13 string `json:"vote_classic_h_pic_13" example:"picture"`
	VoteClassicHPic14 string `json:"vote_classic_h_pic_14" example:"picture"`
	VoteClassicHPic15 string `json:"vote_classic_h_pic_15" example:"picture"`
	VoteClassicHPic16 string `json:"vote_classic_h_pic_16" example:"picture"`
	VoteClassicHPic17 string `json:"vote_classic_h_pic_17" example:"picture"`
	VoteClassicHPic18 string `json:"vote_classic_h_pic_18" example:"picture"`
	VoteClassicHPic19 string `json:"vote_classic_h_pic_19" example:"picture"`
	VoteClassicHPic20 string `json:"vote_classic_h_pic_20" example:"picture"`
	VoteClassicHPic21 string `json:"vote_classic_h_pic_21" example:"picture"`
	// VoteClassicHPic22 string `json:"vote_classic_h_pic_22" example:"picture"`
	VoteClassicHPic23 string `json:"vote_classic_h_pic_23" example:"picture"`
	VoteClassicHPic24 string `json:"vote_classic_h_pic_24" example:"picture"`
	VoteClassicHPic25 string `json:"vote_classic_h_pic_25" example:"picture"`
	VoteClassicHPic26 string `json:"vote_classic_h_pic_26" example:"picture"`
	VoteClassicHPic27 string `json:"vote_classic_h_pic_27" example:"picture"`
	VoteClassicHPic28 string `json:"vote_classic_h_pic_28" example:"picture"`
	VoteClassicHPic29 string `json:"vote_classic_h_pic_29" example:"picture"`
	VoteClassicHPic30 string `json:"vote_classic_h_pic_30" example:"picture"`
	VoteClassicHPic31 string `json:"vote_classic_h_pic_31" example:"picture"`
	VoteClassicHPic32 string `json:"vote_classic_h_pic_32" example:"picture"`
	VoteClassicHPic33 string `json:"vote_classic_h_pic_33" example:"picture"`
	VoteClassicHPic34 string `json:"vote_classic_h_pic_34" example:"picture"`
	VoteClassicHPic35 string `json:"vote_classic_h_pic_35" example:"picture"`
	VoteClassicHPic36 string `json:"vote_classic_h_pic_36" example:"picture"`
	VoteClassicHPic37 string `json:"vote_classic_h_pic_37" example:"picture"`
	VoteClassicGPic01 string `json:"vote_classic_g_pic_01" example:"picture"`
	VoteClassicGPic02 string `json:"vote_classic_g_pic_02" example:"picture"`
	VoteClassicGPic03 string `json:"vote_classic_g_pic_03" example:"picture"`
	VoteClassicGPic04 string `json:"vote_classic_g_pic_04" example:"picture"`
	VoteClassicGPic05 string `json:"vote_classic_g_pic_05" example:"picture"`
	VoteClassicGPic06 string `json:"vote_classic_g_pic_06" example:"picture"`
	VoteClassicGPic07 string `json:"vote_classic_g_pic_07" example:"picture"`
	VoteClassicCPic01 string `json:"vote_classic_c_pic_01" example:"picture"`
	VoteClassicCPic02 string `json:"vote_classic_c_pic_02" example:"picture"`
	VoteClassicCPic03 string `json:"vote_classic_c_pic_03" example:"picture"`
	VoteClassicCPic04 string `json:"vote_classic_c_pic_04" example:"picture"`
	// 音樂
	VoteBgmGaming string `json:"vote_bgm_gaming" example:"picture"`

	// 快問快答
	QA1         string `json:"qa_1" example:"qa"`
	QA1Options  string `json:"qa_1_options" example:"A&&&B&&&C&&&D"`
	QA1Answer   string `json:"qa_1_answer" example:"1"`
	QA1Score    string `json:"qa_1_score" example:"10"`
	QA2         string `json:"qa_2" example:"qa"`
	QA2Options  string `json:"qa_2_options" example:"A&&&B&&&C&&&D"`
	QA2Answer   string `json:"qa_2_answer" example:"1"`
	QA2Score    string `json:"qa_2_score" example:"10"`
	QA3         string `json:"qa_3" example:"qa"`
	QA3Options  string `json:"qa_3_options" example:"A&&&B&&&C&&&D"`
	QA3Answer   string `json:"qa_3_answer" example:"1"`
	QA3Score    string `json:"qa_3_score" example:"10"`
	QA4         string `json:"qa_4" example:"qa"`
	QA4Options  string `json:"qa_4_options" example:"A&&&B&&&C&&&D"`
	QA4Answer   string `json:"qa_4_answer" example:"1"`
	QA4Score    string `json:"qa_4_score" example:"10"`
	QA5         string `json:"qa_5" example:"qa"`
	QA5Options  string `json:"qa_5_options" example:"A&&&B&&&C&&&D"`
	QA5Answer   string `json:"qa_5_answer" example:"1"`
	QA5Score    string `json:"qa_5_score" example:"10"`
	QA6         string `json:"qa_6" example:"qa"`
	QA6Options  string `json:"qa_6_options" example:"A&&&B&&&C&&&D"`
	QA6Answer   string `json:"qa_6_answer" example:"1"`
	QA6Score    string `json:"qa_6_score" example:"10"`
	QA7         string `json:"qa_7" example:"qa"`
	QA7Options  string `json:"qa_7_options" example:"A&&&B&&&C&&&D"`
	QA7Answer   string `json:"qa_7_answer" example:"1"`
	QA7Score    string `json:"qa_7_score" example:"10"`
	QA8         string `json:"qa_8" example:"qa"`
	QA8Options  string `json:"qa_8_options" example:"A&&&B&&&C&&&D"`
	QA8Answer   string `json:"qa_8_answer" example:"1"`
	QA8Score    string `json:"qa_8_score" example:"10"`
	QA9         string `json:"qa_9" example:"qa"`
	QA9Options  string `json:"qa_9_options" example:"A&&&B&&&C&&&D"`
	QA9Answer   string `json:"qa_9_answer" example:"1"`
	QA9Score    string `json:"qa_9_score" example:"10"`
	QA10        string `json:"qa_10" example:"qa"`
	QA10Options string `json:"qa_10_options" example:"A&&&B&&&C&&&D"`
	QA10Answer  string `json:"qa_10_answer" example:"1"`
	QA10Score   string `json:"qa_10_score" example:"10"`
	QA11        string `json:"qa_11" example:"qa"`
	QA11Options string `json:"qa_11_options" example:"A&&&B&&&C&&&D"`
	QA11Answer  string `json:"qa_11_answer" example:"1"`
	QA11Score   string `json:"qa_11_score" example:"10"`
	QA12        string `json:"qa_12" example:"qa"`
	QA12Options string `json:"qa_12_options" example:"A&&&B&&&C&&&D"`
	QA12Answer  string `json:"qa_12_answer" example:"1"`
	QA12Score   string `json:"qa_12_score" example:"10"`
	QA13        string `json:"qa_13" example:"qa"`
	QA13Options string `json:"qa_13_options" example:"A&&&B&&&C&&&D"`
	QA13Answer  string `json:"qa_13_answer" example:"1"`
	QA13Score   string `json:"qa_13_score" example:"10"`
	QA14        string `json:"qa_14" example:"qa"`
	QA14Options string `json:"qa_14_options" example:"A&&&B&&&C&&&D"`
	QA14Answer  string `json:"qa_14_answer" example:"1"`
	QA14Score   string `json:"qa_14_score" example:"10"`
	QA15        string `json:"qa_15" example:"qa"`
	QA15Options string `json:"qa_15_options" example:"A&&&B&&&C&&&D"`
	QA15Answer  string `json:"qa_15_answer" example:"1"`
	QA15Score   string `json:"qa_15_score" example:"10"`
	QA16        string `json:"qa_16" example:"qa"`
	QA16Options string `json:"qa_16_options" example:"A&&&B&&&C&&&D"`
	QA16Answer  string `json:"qa_16_answer" example:"1"`
	QA16Score   string `json:"qa_16_score" example:"10"`
	QA17        string `json:"qa_17" example:"qa"`
	QA17Options string `json:"qa_17_options" example:"A&&&B&&&C&&&D"`
	QA17Answer  string `json:"qa_17_answer" example:"1"`
	QA17Score   string `json:"qa_17_score" example:"10"`
	QA18        string `json:"qa_18" example:"qa"`
	QA18Options string `json:"qa_18_options" example:"A&&&B&&&C&&&D"`
	QA18Answer  string `json:"qa_18_answer" example:"1"`
	QA18Score   string `json:"qa_18_score" example:"10"`
	QA19        string `json:"qa_19" example:"qa"`
	QA19Options string `json:"qa_19_options" example:"A&&&B&&&C&&&D"`
	QA19Answer  string `json:"qa_19_answer" example:"1"`
	QA19Score   string `json:"qa_19_score" example:"10"`
	QA20        string `json:"qa_20" example:"qa"`
	QA20Options string `json:"qa_20_options" example:"A&&&B&&&C&&&D"`
	QA20Answer  string `json:"qa_20_answer" example:"1"`
	QA20Score   string `json:"qa_20_score" example:"10"`
	TotalQA     string `json:"total_qa" example:"1"`   // 總題數
	QASecond    string `json:"qa_second" example:"30"` // 題目顯示秒數

	// Token string `json:"token" example:"token"`
}

// DefaultGameModel 預設GameModel
func DefaultGameModel() GameModel {
	return GameModel{Base: Base{TableName: config.ACTIVITY_GAME_TABLE}}
}

// SetConn 設定connection(mysql.redis.mongo統一設置)
func (a GameModel) SetConn(dbconn db.Connection,
	cacheconn cache.Connection, mongoconn mongo.Connection) GameModel {
	a.DbConn = dbconn
	a.RedisConn = cacheconn
	a.MongoConn = mongoconn
	return a
}

// SetDbConn 設定connection
// func (a GameModel) SetDbConn(conn db.Connection) GameModel {
// 	a.DbConn = conn
// 	return a
// }

// // SetRedisConn 設定connection
// func (a GameModel) SetRedisConn(conn cache.Connection) GameModel {
// 	a.RedisConn = conn
// 	return a
// }

// // SetMongoConn 設定connection
// func (a GameModel) SetMongoConn(conn mongo.Connection) GameModel {
// 	a.MongoConn = conn
// 	return a
// }

// GameRound    string `json:"game_round" example:"1"`
// GameSecond   string `json:"game_second" example:"30"`
// GameStatus   string `json:"game_status" example:"open、start、end、close"`
// GameAttend   string `json:"game_attend" example:"0"`
// LeftTeamGameAttend  string `json:"left_team_game_attend" example:"0"`      // 左方隊伍參加遊戲人數
// RightTeamGameAttend string `json:"right_team_game_attend" example:"0"`     // 右方隊伍參加遊戲人數
// BingoRound string `json:"bingo_round" example:"0"` // 賓果遊戲進行回合數

// NewGameModel 資料表欄位
// type NewGameModel struct {
// 	UserID     string `json:"user_id" example:"user_id"`
// 	ActivityID string `json:"activity_id" example:"activity_id"`
// 	// GameID       string `json:"game_id" example:"game_id"`
// 	// Game         string `json:"game" example:"game name"`
// 	Title         string `json:"title" example:"game title"`
// 	GameType      string `json:"game_type" example:"game type"`
// 	LimitTime     string `json:"limit_time" example:"open、close"`
// 	Second        string `json:"second" example:"30"`
// 	MaxPeople     string `json:"max_people" example:"100(依照用戶權限判斷上限)"`
// 	People        string `json:"people" example:"100(依照max_people資料判斷上限)"`
// 	MaxTimes      string `json:"max_times" example:"10"`
// 	Allow         string `json:"allow" example:"open、close"`
// 	Percent       string `json:"percent" example:"0-100"`
// 	FirstPrize    string `json:"first_prize" example:"50(上限為50人)"`
// 	SecondPrize   string `json:"second_prize" example:"50(上限為50人)"`
// 	ThirdPrize    string `json:"third_prize" example:"100(上限為100人)"`
// 	GeneralPrize  string `json:"general_prize" example:"800(上限為800人)"`
// 	Topic         string `json:"topic" example:"01_classic"`
// 	Skin          string `json:"skin" example:"classic"`
// 	Music         string `json:"music" example:"classic"`
// 	DisplayName   string `json:"display_name" example:"open、close"`   // 是否顯示中獎人員姓名頭像
// 	GameOrder     int64  `json:"game_order" example:"1"`              // 遊戲場次排序
// 	BoxReflection string `json:"box_reflection" example:"open、close"` // 扭蛋機遊戲開關盒反射
// 	SamePeople    string `json:"same_people" example:"open、close"`    // 拔河遊戲人數是否一致
// 	// GameRound    string `json:"game_round" example:"1"`
// 	// GameSecond   string `json:"game_second" example:"30"`
// 	// GameStatus   string `json:"game_status" example:"open、start、end、close"`
// 	// GameAttend   string `json:"game_attend" example:"0"`
// 	// LeftTeamGameAttend  string `json:"left_team_game_attend" example:"0"`      // 左方隊伍參加遊戲人數
// 	// RightTeamGameAttend string `json:"right_team_game_attend" example:"0"`     // 右方隊伍參加遊戲人數

// 	// 拔河遊戲
// 	AllowChooseTeam  string `json:"allow_choose_team" example:"open、close"` // 允許玩家選擇隊伍
// 	LeftTeamName     string `json:"left_team_name" example:"name"`          // 左邊隊伍名稱
// 	LeftTeamPicture  string `json:"left_team_picture" example:"picture"`    // 左邊隊伍照片
// 	RightTeamName    string `json:"right_team_name" example:"name"`         // 右邊隊伍名稱
// 	RightTeamPicture string `json:"right_team_picture" example:"picture"`   // 右邊隊伍照片
// 	Prize            string `json:"prize" example:"uniform、all"`            // 獎品發放

// 	// 賓果遊戲
// 	MaxNumber  string `json:"max_number" example:"0"`  // 最大號碼
// 	BingoLine  string `json:"bingo_line" example:"0"`  // 賓果連線數
// 	RoundPrize string `json:"round_prize" example:"0"` // 每輪發獎人數
// 	// BingoRound string `json:"bingo_round" example:"0"` // 賓果遊戲進行回合數

// 	// 扭蛋機遊戲
// 	GachaMachineReflection string `json:"gacha_machine_reflection" example:"open、close"` // 球的反射度
// 	ReflectiveSwitch       string `json:"reflective_switch" example:"open、close"`        // 反射開關

// 	// 投票遊戲
// 	VoteScreen       string `json:"vote_screen" example:"bar_chart、rank、detail_information"` // 投票畫面(長條圖顯示、排名顯示、詳細資訊顯示)
// 	VoteTimes        string `json:"vote_times" example:"0"`                                  // 人員投票次數
// 	VoteMethod       string `json:"vote_method" example:"all_vote、single_group、all_group"`   // 投票模式(全選項投票)
// 	VoteMethodPlayer string `json:"vote_method_player" example:"one_vote、free_vote"`         // 玩家投票方式(一個選項一票、自由投票)
// 	VoteRestriction  string `json:"vote_restriction" example:"all_player、special_officer"`   // 投票限制(所有人員都能投票、特殊人員才能投票)
// 	AvatarShape      string `json:"avatar_shape" example:"circle、square"`                    // 選項框是圓形還是方形
// 	VoteStartTime    string `json:"vote_start_time" example:""`                              // 投票開始時間
// 	VoteEndTime      string `json:"vote_end_time" example:""`                                // 投票結束時間
// 	AutoPlay         string `json:"auto_play" example:"open、close"`                          // 自動輪播
// 	ShowRank         string `json:"show_rank" example:"open、close"`                          // 排名展示
// 	TitleSwitch      string `json:"title_switch" example:"open、close"`                       // 場次名稱
// 	ArrangementGuest string `json:"arrangement_guest" example:"list、side_by_side"`           // 玩家端選項排列方式

// 	// 敲敲樂自定義
// 	WhackmoleClassicHPic01 string `json:"whackmole_classic_h_pic_01" example:"picture"`
// 	WhackmoleClassicHPic02 string `json:"whackmole_classic_h_pic_02" example:"picture"`
// 	WhackmoleClassicHPic03 string `json:"whackmole_classic_h_pic_03" example:"picture"`
// 	WhackmoleClassicHPic04 string `json:"whackmole_classic_h_pic_04" example:"picture"`
// 	WhackmoleClassicHPic05 string `json:"whackmole_classic_h_pic_05" example:"picture"`
// 	WhackmoleClassicHPic06 string `json:"whackmole_classic_h_pic_06" example:"picture"`
// 	WhackmoleClassicHPic07 string `json:"whackmole_classic_h_pic_07" example:"picture"`
// 	WhackmoleClassicHPic08 string `json:"whackmole_classic_h_pic_08" example:"picture"`
// 	WhackmoleClassicHPic09 string `json:"whackmole_classic_h_pic_09" example:"picture"`
// 	WhackmoleClassicHPic10 string `json:"whackmole_classic_h_pic_10" example:"picture"`
// 	WhackmoleClassicHPic11 string `json:"whackmole_classic_h_pic_11" example:"picture"`
// 	WhackmoleClassicHPic12 string `json:"whackmole_classic_h_pic_12" example:"picture"`
// 	WhackmoleClassicHPic13 string `json:"whackmole_classic_h_pic_13" example:"picture"`
// 	WhackmoleClassicHPic14 string `json:"whackmole_classic_h_pic_14" example:"picture"`
// 	WhackmoleClassicHPic15 string `json:"whackmole_classic_h_pic_15" example:"picture"`
// 	WhackmoleClassicGPic01 string `json:"whackmole_classic_g_pic_01" example:"picture"`
// 	WhackmoleClassicGPic02 string `json:"whackmole_classic_g_pic_02" example:"picture"`
// 	WhackmoleClassicGPic03 string `json:"whackmole_classic_g_pic_03" example:"picture"`
// 	WhackmoleClassicGPic04 string `json:"whackmole_classic_g_pic_04" example:"picture"`
// 	WhackmoleClassicGPic05 string `json:"whackmole_classic_g_pic_05" example:"picture"`
// 	WhackmoleClassicCPic01 string `json:"whackmole_classic_c_pic_01" example:"picture"`
// 	WhackmoleClassicCPic02 string `json:"whackmole_classic_c_pic_02" example:"picture"`
// 	WhackmoleClassicCPic03 string `json:"whackmole_classic_c_pic_03" example:"picture"`
// 	WhackmoleClassicCPic04 string `json:"whackmole_classic_c_pic_04" example:"picture"`
// 	WhackmoleClassicCPic05 string `json:"whackmole_classic_c_pic_05" example:"picture"`
// 	WhackmoleClassicCPic06 string `json:"whackmole_classic_c_pic_06" example:"picture"`
// 	WhackmoleClassicCPic07 string `json:"whackmole_classic_c_pic_07" example:"picture"`
// 	WhackmoleClassicCPic08 string `json:"whackmole_classic_c_pic_08" example:"picture"`
// 	WhackmoleClassicCAni01 string `json:"whackmole_classic_c_ani_01" example:"picture"`

// 	WhackmoleHalloweenHPic01 string `json:"whackmole_halloween_h_pic_01" example:"picture"`
// 	WhackmoleHalloweenHPic02 string `json:"whackmole_halloween_h_pic_02" example:"picture"`
// 	WhackmoleHalloweenHPic03 string `json:"whackmole_halloween_h_pic_03" example:"picture"`
// 	WhackmoleHalloweenHPic04 string `json:"whackmole_halloween_h_pic_04" example:"picture"`
// 	WhackmoleHalloweenHPic05 string `json:"whackmole_halloween_h_pic_05" example:"picture"`
// 	WhackmoleHalloweenHPic06 string `json:"whackmole_halloween_h_pic_06" example:"picture"`
// 	WhackmoleHalloweenHPic07 string `json:"whackmole_halloween_h_pic_07" example:"picture"`
// 	WhackmoleHalloweenHPic08 string `json:"whackmole_halloween_h_pic_08" example:"picture"`
// 	WhackmoleHalloweenHPic09 string `json:"whackmole_halloween_h_pic_09" example:"picture"`
// 	WhackmoleHalloweenHPic10 string `json:"whackmole_halloween_h_pic_10" example:"picture"`
// 	WhackmoleHalloweenHPic11 string `json:"whackmole_halloween_h_pic_11" example:"picture"`
// 	WhackmoleHalloweenHPic12 string `json:"whackmole_halloween_h_pic_12" example:"picture"`
// 	WhackmoleHalloweenHPic13 string `json:"whackmole_halloween_h_pic_13" example:"picture"`
// 	WhackmoleHalloweenHPic14 string `json:"whackmole_halloween_h_pic_14" example:"picture"`
// 	WhackmoleHalloweenHPic15 string `json:"whackmole_halloween_h_pic_15" example:"picture"`
// 	WhackmoleHalloweenGPic01 string `json:"whackmole_halloween_g_pic_01" example:"picture"`
// 	WhackmoleHalloweenGPic02 string `json:"whackmole_halloween_g_pic_02" example:"picture"`
// 	WhackmoleHalloweenGPic03 string `json:"whackmole_halloween_g_pic_03" example:"picture"`
// 	WhackmoleHalloweenGPic04 string `json:"whackmole_halloween_g_pic_04" example:"picture"`
// 	WhackmoleHalloweenGPic05 string `json:"whackmole_halloween_g_pic_05" example:"picture"`
// 	WhackmoleHalloweenCPic01 string `json:"whackmole_halloween_c_pic_01" example:"picture"`
// 	WhackmoleHalloweenCPic02 string `json:"whackmole_halloween_c_pic_02" example:"picture"`
// 	WhackmoleHalloweenCPic03 string `json:"whackmole_halloween_c_pic_03" example:"picture"`
// 	WhackmoleHalloweenCPic04 string `json:"whackmole_halloween_c_pic_04" example:"picture"`
// 	WhackmoleHalloweenCPic05 string `json:"whackmole_halloween_c_pic_05" example:"picture"`
// 	WhackmoleHalloweenCPic06 string `json:"whackmole_halloween_c_pic_06" example:"picture"`
// 	WhackmoleHalloweenCPic07 string `json:"whackmole_halloween_c_pic_07" example:"picture"`
// 	WhackmoleHalloweenCPic08 string `json:"whackmole_halloween_c_pic_08" example:"picture"`
// 	WhackmoleHalloweenCAni01 string `json:"whackmole_halloween_c_ani_01" example:"picture"`

// 	WhackmoleChristmasHPic01 string `json:"whackmole_christmas_h_pic_01" example:"picture"`
// 	WhackmoleChristmasHPic02 string `json:"whackmole_christmas_h_pic_02" example:"picture"`
// 	WhackmoleChristmasHPic03 string `json:"whackmole_christmas_h_pic_03" example:"picture"`
// 	WhackmoleChristmasHPic04 string `json:"whackmole_christmas_h_pic_04" example:"picture"`
// 	WhackmoleChristmasHPic05 string `json:"whackmole_christmas_h_pic_05" example:"picture"`
// 	WhackmoleChristmasHPic06 string `json:"whackmole_christmas_h_pic_06" example:"picture"`
// 	WhackmoleChristmasHPic07 string `json:"whackmole_christmas_h_pic_07" example:"picture"`
// 	WhackmoleChristmasHPic08 string `json:"whackmole_christmas_h_pic_08" example:"picture"`
// 	WhackmoleChristmasHPic09 string `json:"whackmole_christmas_h_pic_09" example:"picture"`
// 	WhackmoleChristmasHPic10 string `json:"whackmole_christmas_h_pic_10" example:"picture"`
// 	WhackmoleChristmasHPic11 string `json:"whackmole_christmas_h_pic_11" example:"picture"`
// 	WhackmoleChristmasHPic12 string `json:"whackmole_christmas_h_pic_12" example:"picture"`
// 	WhackmoleChristmasHPic13 string `json:"whackmole_christmas_h_pic_13" example:"picture"`
// 	WhackmoleChristmasHPic14 string `json:"whackmole_christmas_h_pic_14" example:"picture"`
// 	WhackmoleChristmasGPic01 string `json:"whackmole_christmas_g_pic_01" example:"picture"`
// 	WhackmoleChristmasGPic02 string `json:"whackmole_christmas_g_pic_02" example:"picture"`
// 	WhackmoleChristmasGPic03 string `json:"whackmole_christmas_g_pic_03" example:"picture"`
// 	WhackmoleChristmasGPic04 string `json:"whackmole_christmas_g_pic_04" example:"picture"`
// 	WhackmoleChristmasGPic05 string `json:"whackmole_christmas_g_pic_05" example:"picture"`
// 	WhackmoleChristmasGPic06 string `json:"whackmole_christmas_g_pic_06" example:"picture"`
// 	WhackmoleChristmasGPic07 string `json:"whackmole_christmas_g_pic_07" example:"picture"`
// 	WhackmoleChristmasGPic08 string `json:"whackmole_christmas_g_pic_08" example:"picture"`
// 	WhackmoleChristmasCPic01 string `json:"whackmole_christmas_c_pic_01" example:"picture"`
// 	WhackmoleChristmasCPic02 string `json:"whackmole_christmas_c_pic_02" example:"picture"`
// 	WhackmoleChristmasCPic03 string `json:"whackmole_christmas_c_pic_03" example:"picture"`
// 	WhackmoleChristmasCPic04 string `json:"whackmole_christmas_c_pic_04" example:"picture"`
// 	WhackmoleChristmasCPic05 string `json:"whackmole_christmas_c_pic_05" example:"picture"`
// 	WhackmoleChristmasCPic06 string `json:"whackmole_christmas_c_pic_06" example:"picture"`
// 	WhackmoleChristmasCPic07 string `json:"whackmole_christmas_c_pic_07" example:"picture"`
// 	WhackmoleChristmasCPic08 string `json:"whackmole_christmas_c_pic_08" example:"picture"`
// 	WhackmoleChristmasCAni01 string `json:"whackmole_christmas_c_ani_01" example:"picture"`
// 	WhackmoleChristmasCAni02 string `json:"whackmole_christmas_c_ani_02" example:"picture"`

// 	// 敲敲樂音樂
// 	WhackmoleBgmStart  string `json:"whackmole_bgm_start" example:"picture"`  // 遊戲開始
// 	WhackmoleBgmGaming string `json:"whackmole_bgm_gaming" example:"picture"` // 遊戲進行中
// 	WhackmoleBgmEnd    string `json:"whackmole_bgm_end" example:"picture"`    // 遊戲結束

// 	// 搖號抽獎自定義
// 	DrawNumbersClassicHPic01 string `json:"draw_numbers_classic_h_pic_01" example:"picture"`
// 	DrawNumbersClassicHPic02 string `json:"draw_numbers_classic_h_pic_02" example:"picture"`
// 	DrawNumbersClassicHPic03 string `json:"draw_numbers_classic_h_pic_03" example:"picture"`
// 	DrawNumbersClassicHPic04 string `json:"draw_numbers_classic_h_pic_04" example:"picture"`
// 	DrawNumbersClassicHPic05 string `json:"draw_numbers_classic_h_pic_05" example:"picture"`
// 	DrawNumbersClassicHPic06 string `json:"draw_numbers_classic_h_pic_06" example:"picture"`
// 	DrawNumbersClassicHPic07 string `json:"draw_numbers_classic_h_pic_07" example:"picture"`
// 	DrawNumbersClassicHPic08 string `json:"draw_numbers_classic_h_pic_08" example:"picture"`
// 	DrawNumbersClassicHPic09 string `json:"draw_numbers_classic_h_pic_09" example:"picture"`
// 	DrawNumbersClassicHPic10 string `json:"draw_numbers_classic_h_pic_10" example:"picture"`
// 	DrawNumbersClassicHPic11 string `json:"draw_numbers_classic_h_pic_11" example:"picture"`
// 	DrawNumbersClassicHPic12 string `json:"draw_numbers_classic_h_pic_12" example:"picture"`
// 	DrawNumbersClassicHPic13 string `json:"draw_numbers_classic_h_pic_13" example:"picture"`
// 	DrawNumbersClassicHPic14 string `json:"draw_numbers_classic_h_pic_14" example:"picture"`
// 	DrawNumbersClassicHPic15 string `json:"draw_numbers_classic_h_pic_15" example:"picture"`
// 	DrawNumbersClassicHPic16 string `json:"draw_numbers_classic_h_pic_16" example:"picture"`
// 	DrawNumbersClassicHAni01 string `json:"draw_numbers_classic_h_ani_01" example:"picture"`

// 	DrawNumbersGoldHPic01 string `json:"draw_numbers_gold_h_pic_01" example:"picture"`
// 	DrawNumbersGoldHPic02 string `json:"draw_numbers_gold_h_pic_02" example:"picture"`
// 	DrawNumbersGoldHPic03 string `json:"draw_numbers_gold_h_pic_03" example:"picture"`
// 	DrawNumbersGoldHPic04 string `json:"draw_numbers_gold_h_pic_04" example:"picture"`
// 	DrawNumbersGoldHPic05 string `json:"draw_numbers_gold_h_pic_05" example:"picture"`
// 	DrawNumbersGoldHPic06 string `json:"draw_numbers_gold_h_pic_06" example:"picture"`
// 	DrawNumbersGoldHPic07 string `json:"draw_numbers_gold_h_pic_07" example:"picture"`
// 	DrawNumbersGoldHPic08 string `json:"draw_numbers_gold_h_pic_08" example:"picture"`
// 	DrawNumbersGoldHPic09 string `json:"draw_numbers_gold_h_pic_09" example:"picture"`
// 	DrawNumbersGoldHPic10 string `json:"draw_numbers_gold_h_pic_10" example:"picture"`
// 	DrawNumbersGoldHPic11 string `json:"draw_numbers_gold_h_pic_11" example:"picture"`
// 	DrawNumbersGoldHPic12 string `json:"draw_numbers_gold_h_pic_12" example:"picture"`
// 	DrawNumbersGoldHPic13 string `json:"draw_numbers_gold_h_pic_13" example:"picture"`
// 	DrawNumbersGoldHPic14 string `json:"draw_numbers_gold_h_pic_14" example:"picture"`
// 	DrawNumbersGoldHAni01 string `json:"draw_numbers_gold_h_ani_01" example:"picture"`
// 	DrawNumbersGoldHAni02 string `json:"draw_numbers_gold_h_ani_02" example:"picture"`
// 	DrawNumbersGoldHAni03 string `json:"draw_numbers_gold_h_ani_03" example:"picture"`

// 	DrawNumbersNewyearDragonHPic01 string `json:"draw_numbers_newyear_dragon_h_pic_01" example:"picture"`
// 	DrawNumbersNewyearDragonHPic02 string `json:"draw_numbers_newyear_dragon_h_pic_02" example:"picture"`
// 	DrawNumbersNewyearDragonHPic03 string `json:"draw_numbers_newyear_dragon_h_pic_03" example:"picture"`
// 	DrawNumbersNewyearDragonHPic04 string `json:"draw_numbers_newyear_dragon_h_pic_04" example:"picture"`
// 	DrawNumbersNewyearDragonHPic05 string `json:"draw_numbers_newyear_dragon_h_pic_05" example:"picture"`
// 	DrawNumbersNewyearDragonHPic06 string `json:"draw_numbers_newyear_dragon_h_pic_06" example:"picture"`
// 	DrawNumbersNewyearDragonHPic07 string `json:"draw_numbers_newyear_dragon_h_pic_07" example:"picture"`
// 	DrawNumbersNewyearDragonHPic08 string `json:"draw_numbers_newyear_dragon_h_pic_08" example:"picture"`
// 	DrawNumbersNewyearDragonHPic09 string `json:"draw_numbers_newyear_dragon_h_pic_09" example:"picture"`
// 	DrawNumbersNewyearDragonHPic10 string `json:"draw_numbers_newyear_dragon_h_pic_10" example:"picture"`
// 	DrawNumbersNewyearDragonHPic11 string `json:"draw_numbers_newyear_dragon_h_pic_11" example:"picture"`
// 	DrawNumbersNewyearDragonHPic12 string `json:"draw_numbers_newyear_dragon_h_pic_12" example:"picture"`
// 	DrawNumbersNewyearDragonHPic13 string `json:"draw_numbers_newyear_dragon_h_pic_13" example:"picture"`
// 	DrawNumbersNewyearDragonHPic14 string `json:"draw_numbers_newyear_dragon_h_pic_14" example:"picture"`
// 	DrawNumbersNewyearDragonHPic15 string `json:"draw_numbers_newyear_dragon_h_pic_15" example:"picture"`
// 	DrawNumbersNewyearDragonHPic16 string `json:"draw_numbers_newyear_dragon_h_pic_16" example:"picture"`
// 	DrawNumbersNewyearDragonHPic17 string `json:"draw_numbers_newyear_dragon_h_pic_17" example:"picture"`
// 	DrawNumbersNewyearDragonHPic18 string `json:"draw_numbers_newyear_dragon_h_pic_18" example:"picture"`
// 	DrawNumbersNewyearDragonHPic19 string `json:"draw_numbers_newyear_dragon_h_pic_19" example:"picture"`
// 	DrawNumbersNewyearDragonHPic20 string `json:"draw_numbers_newyear_dragon_h_pic_20" example:"picture"`
// 	DrawNumbersNewyearDragonHAni01 string `json:"draw_numbers_newyear_dragon_h_ani_01" example:"picture"`
// 	DrawNumbersNewyearDragonHAni02 string `json:"draw_numbers_newyear_dragon_h_ani_02" example:"picture"`

// 	DrawNumbersCherryHPic01 string `json:"draw_numbers_cherry_h_pic_01" example:"picture"`
// 	DrawNumbersCherryHPic02 string `json:"draw_numbers_cherry_h_pic_02" example:"picture"`
// 	DrawNumbersCherryHPic03 string `json:"draw_numbers_cherry_h_pic_03" example:"picture"`
// 	DrawNumbersCherryHPic04 string `json:"draw_numbers_cherry_h_pic_04" example:"picture"`
// 	DrawNumbersCherryHPic05 string `json:"draw_numbers_cherry_h_pic_05" example:"picture"`
// 	DrawNumbersCherryHPic06 string `json:"draw_numbers_cherry_h_pic_06" example:"picture"`
// 	DrawNumbersCherryHPic07 string `json:"draw_numbers_cherry_h_pic_07" example:"picture"`
// 	DrawNumbersCherryHPic08 string `json:"draw_numbers_cherry_h_pic_08" example:"picture"`
// 	DrawNumbersCherryHPic09 string `json:"draw_numbers_cherry_h_pic_09" example:"picture"`
// 	DrawNumbersCherryHPic10 string `json:"draw_numbers_cherry_h_pic_10" example:"picture"`
// 	DrawNumbersCherryHPic11 string `json:"draw_numbers_cherry_h_pic_11" example:"picture"`
// 	DrawNumbersCherryHPic12 string `json:"draw_numbers_cherry_h_pic_12" example:"picture"`
// 	DrawNumbersCherryHPic13 string `json:"draw_numbers_cherry_h_pic_13" example:"picture"`
// 	DrawNumbersCherryHPic14 string `json:"draw_numbers_cherry_h_pic_14" example:"picture"`
// 	DrawNumbersCherryHPic15 string `json:"draw_numbers_cherry_h_pic_15" example:"picture"`
// 	DrawNumbersCherryHPic16 string `json:"draw_numbers_cherry_h_pic_16" example:"picture"`
// 	DrawNumbersCherryHPic17 string `json:"draw_numbers_cherry_h_pic_17" example:"picture"`
// 	DrawNumbersCherryHAni01 string `json:"draw_numbers_cherry_h_ani_01" example:"picture"`
// 	DrawNumbersCherryHAni02 string `json:"draw_numbers_cherry_h_ani_02" example:"picture"`
// 	DrawNumbersCherryHAni03 string `json:"draw_numbers_cherry_h_ani_03" example:"picture"`
// 	DrawNumbersCherryHAni04 string `json:"draw_numbers_cherry_h_ani_04" example:"picture"`

// 	// 太空主題
// 	DrawNumbers3DSpaceHPic01 string `json:"draw_numbers_3D_space_h_pic_01" example:"picture"`
// 	DrawNumbers3DSpaceHPic02 string `json:"draw_numbers_3D_space_h_pic_02" example:"picture"`
// 	DrawNumbers3DSpaceHPic03 string `json:"draw_numbers_3D_space_h_pic_03" example:"picture"`
// 	DrawNumbers3DSpaceHPic04 string `json:"draw_numbers_3D_space_h_pic_04" example:"picture"`
// 	DrawNumbers3DSpaceHPic05 string `json:"draw_numbers_3D_space_h_pic_05" example:"picture"`
// 	DrawNumbers3DSpaceHPic06 string `json:"draw_numbers_3D_space_h_pic_06" example:"picture"`
// 	DrawNumbers3DSpaceHPic07 string `json:"draw_numbers_3D_space_h_pic_07" example:"picture"`
// 	DrawNumbers3DSpaceHPic08 string `json:"draw_numbers_3D_space_h_pic_08" example:"picture"`

// 	// 音樂
// 	DrawNumbersBgmGaming string `json:"draw_numbers_bgm_gaming" example:"picture"` // 遊戲進行中

// 	// 快問快答自定義
// 	QAClassicHPic01 string `json:"qa_classic_h_pic_01" example:"picture"`
// 	QAClassicHPic02 string `json:"qa_classic_h_pic_02" example:"picture"`
// 	QAClassicHPic03 string `json:"qa_classic_h_pic_03" example:"picture"`
// 	QAClassicHPic04 string `json:"qa_classic_h_pic_04" example:"picture"`
// 	QAClassicHPic05 string `json:"qa_classic_h_pic_05" example:"picture"`
// 	QAClassicHPic06 string `json:"qa_classic_h_pic_06" example:"picture"`
// 	QAClassicHPic07 string `json:"qa_classic_h_pic_07" example:"picture"`
// 	QAClassicHPic08 string `json:"qa_classic_h_pic_08" example:"picture"`
// 	QAClassicHPic09 string `json:"qa_classic_h_pic_09" example:"picture"`
// 	QAClassicHPic10 string `json:"qa_classic_h_pic_10" example:"picture"`
// 	QAClassicHPic11 string `json:"qa_classic_h_pic_11" example:"picture"`
// 	QAClassicHPic12 string `json:"qa_classic_h_pic_12" example:"picture"`
// 	QAClassicHPic13 string `json:"qa_classic_h_pic_13" example:"picture"`
// 	QAClassicHPic14 string `json:"qa_classic_h_pic_14" example:"picture"`
// 	QAClassicHPic15 string `json:"qa_classic_h_pic_15" example:"picture"`
// 	QAClassicHPic16 string `json:"qa_classic_h_pic_16" example:"picture"`
// 	QAClassicHPic17 string `json:"qa_classic_h_pic_17" example:"picture"`
// 	QAClassicHPic18 string `json:"qa_classic_h_pic_18" example:"picture"`
// 	QAClassicHPic19 string `json:"qa_classic_h_pic_19" example:"picture"`
// 	QAClassicHPic20 string `json:"qa_classic_h_pic_20" example:"picture"`
// 	QAClassicHPic21 string `json:"qa_classic_h_pic_21" example:"picture"`
// 	QAClassicHPic22 string `json:"qa_classic_h_pic_22" example:"picture"`
// 	QAClassicGPic01 string `json:"qa_classic_g_pic_01" example:"picture"`
// 	QAClassicGPic02 string `json:"qa_classic_g_pic_02" example:"picture"`
// 	QAClassicGPic03 string `json:"qa_classic_g_pic_03" example:"picture"`
// 	QAClassicGPic04 string `json:"qa_classic_g_pic_04" example:"picture"`
// 	QAClassicGPic05 string `json:"qa_classic_g_pic_05" example:"picture"`
// 	QAClassicCPic01 string `json:"qa_classic_c_pic_01" example:"picture"`
// 	QAClassicHAni01 string `json:"qa_classic_h_ani_01" example:"picture"`
// 	QAClassicHAni02 string `json:"qa_classic_h_ani_02" example:"picture"`
// 	QAClassicGAni01 string `json:"qa_classic_g_ani_01" example:"picture"`
// 	QAClassicGAni02 string `json:"qa_classic_g_ani_02" example:"picture"`

// 	QAElectricHPic01 string `json:"qa_electric_h_pic_01" example:"picture"`
// 	QAElectricHPic02 string `json:"qa_electric_h_pic_02" example:"picture"`
// 	QAElectricHPic03 string `json:"qa_electric_h_pic_03" example:"picture"`
// 	QAElectricHPic04 string `json:"qa_electric_h_pic_04" example:"picture"`
// 	QAElectricHPic05 string `json:"qa_electric_h_pic_05" example:"picture"`
// 	QAElectricHPic06 string `json:"qa_electric_h_pic_06" example:"picture"`
// 	QAElectricHPic07 string `json:"qa_electric_h_pic_07" example:"picture"`
// 	QAElectricHPic08 string `json:"qa_electric_h_pic_08" example:"picture"`
// 	QAElectricHPic09 string `json:"qa_electric_h_pic_09" example:"picture"`
// 	QAElectricHPic10 string `json:"qa_electric_h_pic_10" example:"picture"`
// 	QAElectricHPic11 string `json:"qa_electric_h_pic_11" example:"picture"`
// 	QAElectricHPic12 string `json:"qa_electric_h_pic_12" example:"picture"`
// 	QAElectricHPic13 string `json:"qa_electric_h_pic_13" example:"picture"`
// 	QAElectricHPic14 string `json:"qa_electric_h_pic_14" example:"picture"`
// 	QAElectricHPic15 string `json:"qa_electric_h_pic_15" example:"picture"`
// 	QAElectricHPic16 string `json:"qa_electric_h_pic_16" example:"picture"`
// 	QAElectricHPic17 string `json:"qa_electric_h_pic_17" example:"picture"`
// 	QAElectricHPic18 string `json:"qa_electric_h_pic_18" example:"picture"`
// 	QAElectricHPic19 string `json:"qa_electric_h_pic_19" example:"picture"`
// 	QAElectricHPic20 string `json:"qa_electric_h_pic_20" example:"picture"`
// 	QAElectricHPic21 string `json:"qa_electric_h_pic_21" example:"picture"`
// 	QAElectricHPic22 string `json:"qa_electric_h_pic_22" example:"picture"`
// 	QAElectricHPic23 string `json:"qa_electric_h_pic_23" example:"picture"`
// 	QAElectricHPic24 string `json:"qa_electric_h_pic_24" example:"picture"`
// 	QAElectricHPic25 string `json:"qa_electric_h_pic_25" example:"picture"`
// 	QAElectricHPic26 string `json:"qa_electric_h_pic_26" example:"picture"`
// 	QAElectricGPic01 string `json:"qa_electric_g_pic_01" example:"picture"`
// 	QAElectricGPic02 string `json:"qa_electric_g_pic_02" example:"picture"`
// 	QAElectricGPic03 string `json:"qa_electric_g_pic_03" example:"picture"`
// 	QAElectricGPic04 string `json:"qa_electric_g_pic_04" example:"picture"`
// 	QAElectricGPic05 string `json:"qa_electric_g_pic_05" example:"picture"`
// 	QAElectricGPic06 string `json:"qa_electric_g_pic_06" example:"picture"`
// 	QAElectricGPic07 string `json:"qa_electric_g_pic_07" example:"picture"`
// 	QAElectricGPic08 string `json:"qa_electric_g_pic_08" example:"picture"`
// 	QAElectricGPic09 string `json:"qa_electric_g_pic_09" example:"picture"`
// 	QAElectricCPic01 string `json:"qa_electric_c_pic_01" example:"picture"`
// 	QAElectricHAni01 string `json:"qa_electric_h_ani_01" example:"picture"`
// 	QAElectricHAni02 string `json:"qa_electric_h_ani_02" example:"picture"`
// 	QAElectricHAni03 string `json:"qa_electric_h_ani_03" example:"picture"`
// 	QAElectricHAni04 string `json:"qa_electric_h_ani_04" example:"picture"`
// 	QAElectricHAni05 string `json:"qa_electric_h_ani_05" example:"picture"`
// 	QAElectricGAni01 string `json:"qa_electric_g_ani_01" example:"picture"`
// 	QAElectricGAni02 string `json:"qa_electric_g_ani_02" example:"picture"`
// 	QAElectricCAni01 string `json:"qa_electric_c_ani_01" example:"picture"`

// 	QAMoonfestivalHPic01 string `json:"qa_moonfestival_h_pic_01" example:"picture"`
// 	QAMoonfestivalHPic02 string `json:"qa_moonfestival_h_pic_02" example:"picture"`
// 	QAMoonfestivalHPic03 string `json:"qa_moonfestival_h_pic_03" example:"picture"`
// 	QAMoonfestivalHPic04 string `json:"qa_moonfestival_h_pic_04" example:"picture"`
// 	QAMoonfestivalHPic05 string `json:"qa_moonfestival_h_pic_05" example:"picture"`
// 	QAMoonfestivalHPic06 string `json:"qa_moonfestival_h_pic_06" example:"picture"`
// 	QAMoonfestivalHPic07 string `json:"qa_moonfestival_h_pic_07" example:"picture"`
// 	QAMoonfestivalHPic08 string `json:"qa_moonfestival_h_pic_08" example:"picture"`
// 	QAMoonfestivalHPic09 string `json:"qa_moonfestival_h_pic_09" example:"picture"`
// 	QAMoonfestivalHPic10 string `json:"qa_moonfestival_h_pic_10" example:"picture"`
// 	QAMoonfestivalHPic11 string `json:"qa_moonfestival_h_pic_11" example:"picture"`
// 	QAMoonfestivalHPic12 string `json:"qa_moonfestival_h_pic_12" example:"picture"`
// 	QAMoonfestivalHPic13 string `json:"qa_moonfestival_h_pic_13" example:"picture"`
// 	QAMoonfestivalHPic14 string `json:"qa_moonfestival_h_pic_14" example:"picture"`
// 	QAMoonfestivalHPic15 string `json:"qa_moonfestival_h_pic_15" example:"picture"`
// 	QAMoonfestivalHPic16 string `json:"qa_moonfestival_h_pic_16" example:"picture"`
// 	QAMoonfestivalHPic17 string `json:"qa_moonfestival_h_pic_17" example:"picture"`
// 	QAMoonfestivalHPic18 string `json:"qa_moonfestival_h_pic_18" example:"picture"`
// 	QAMoonfestivalHPic19 string `json:"qa_moonfestival_h_pic_19" example:"picture"`
// 	QAMoonfestivalHPic20 string `json:"qa_moonfestival_h_pic_20" example:"picture"`
// 	QAMoonfestivalHPic21 string `json:"qa_moonfestival_h_pic_21" example:"picture"`
// 	QAMoonfestivalHPic22 string `json:"qa_moonfestival_h_pic_22" example:"picture"`
// 	QAMoonfestivalHPic23 string `json:"qa_moonfestival_h_pic_23" example:"picture"`
// 	QAMoonfestivalHPic24 string `json:"qa_moonfestival_h_pic_24" example:"picture"`
// 	QAMoonfestivalGPic01 string `json:"qa_moonfestival_g_pic_01" example:"picture"`
// 	QAMoonfestivalGPic02 string `json:"qa_moonfestival_g_pic_02" example:"picture"`
// 	QAMoonfestivalGPic03 string `json:"qa_moonfestival_g_pic_03" example:"picture"`
// 	QAMoonfestivalGPic04 string `json:"qa_moonfestival_g_pic_04" example:"picture"`
// 	QAMoonfestivalGPic05 string `json:"qa_moonfestival_g_pic_05" example:"picture"`
// 	QAMoonfestivalCPic01 string `json:"qa_moonfestival_c_pic_01" example:"picture"`
// 	QAMoonfestivalCPic02 string `json:"qa_moonfestival_c_pic_02" example:"picture"`
// 	QAMoonfestivalCPic03 string `json:"qa_moonfestival_c_pic_03" example:"picture"`
// 	QAMoonfestivalHAni01 string `json:"qa_moonfestival_h_ani_01" example:"picture"`
// 	QAMoonfestivalHAni02 string `json:"qa_moonfestival_h_ani_02" example:"picture"`
// 	QAMoonfestivalGAni01 string `json:"qa_moonfestival_g_ani_01" example:"picture"`
// 	QAMoonfestivalGAni02 string `json:"qa_moonfestival_g_ani_02" example:"picture"`
// 	QAMoonfestivalGAni03 string `json:"qa_moonfestival_g_ani_03" example:"picture"`

// 	QANewyearDragonHPic01 string `json:"qa_newyear_dragon_h_pic_01" example:"picture"`
// 	QANewyearDragonHPic02 string `json:"qa_newyear_dragon_h_pic_02" example:"picture"`
// 	QANewyearDragonHPic03 string `json:"qa_newyear_dragon_h_pic_03" example:"picture"`
// 	QANewyearDragonHPic04 string `json:"qa_newyear_dragon_h_pic_04" example:"picture"`
// 	QANewyearDragonHPic05 string `json:"qa_newyear_dragon_h_pic_05" example:"picture"`
// 	QANewyearDragonHPic06 string `json:"qa_newyear_dragon_h_pic_06" example:"picture"`
// 	QANewyearDragonHPic07 string `json:"qa_newyear_dragon_h_pic_07" example:"picture"`
// 	QANewyearDragonHPic08 string `json:"qa_newyear_dragon_h_pic_08" example:"picture"`
// 	QANewyearDragonHPic09 string `json:"qa_newyear_dragon_h_pic_09" example:"picture"`
// 	QANewyearDragonHPic10 string `json:"qa_newyear_dragon_h_pic_10" example:"picture"`
// 	QANewyearDragonHPic11 string `json:"qa_newyear_dragon_h_pic_11" example:"picture"`
// 	QANewyearDragonHPic12 string `json:"qa_newyear_dragon_h_pic_12" example:"picture"`
// 	QANewyearDragonHPic13 string `json:"qa_newyear_dragon_h_pic_13" example:"picture"`
// 	QANewyearDragonHPic14 string `json:"qa_newyear_dragon_h_pic_14" example:"picture"`
// 	QANewyearDragonHPic15 string `json:"qa_newyear_dragon_h_pic_15" example:"picture"`
// 	QANewyearDragonHPic16 string `json:"qa_newyear_dragon_h_pic_16" example:"picture"`
// 	QANewyearDragonHPic17 string `json:"qa_newyear_dragon_h_pic_17" example:"picture"`
// 	QANewyearDragonHPic18 string `json:"qa_newyear_dragon_h_pic_18" example:"picture"`
// 	QANewyearDragonHPic19 string `json:"qa_newyear_dragon_h_pic_19" example:"picture"`
// 	QANewyearDragonHPic20 string `json:"qa_newyear_dragon_h_pic_20" example:"picture"`
// 	QANewyearDragonHPic21 string `json:"qa_newyear_dragon_h_pic_21" example:"picture"`
// 	QANewyearDragonHPic22 string `json:"qa_newyear_dragon_h_pic_22" example:"picture"`
// 	QANewyearDragonHPic23 string `json:"qa_newyear_dragon_h_pic_23" example:"picture"`
// 	QANewyearDragonHPic24 string `json:"qa_newyear_dragon_h_pic_24" example:"picture"`
// 	QANewyearDragonGPic01 string `json:"qa_newyear_dragon_g_pic_01" example:"picture"`
// 	QANewyearDragonGPic02 string `json:"qa_newyear_dragon_g_pic_02" example:"picture"`
// 	QANewyearDragonGPic03 string `json:"qa_newyear_dragon_g_pic_03" example:"picture"`
// 	QANewyearDragonGPic04 string `json:"qa_newyear_dragon_g_pic_04" example:"picture"`
// 	QANewyearDragonGPic05 string `json:"qa_newyear_dragon_g_pic_05" example:"picture"`
// 	QANewyearDragonGPic06 string `json:"qa_newyear_dragon_g_pic_06" example:"picture"`
// 	QANewyearDragonCPic01 string `json:"qa_newyear_dragon_c_pic_01" example:"picture"`
// 	QANewyearDragonHAni01 string `json:"qa_newyear_dragon_h_ani_01" example:"picture"`
// 	QANewyearDragonHAni02 string `json:"qa_newyear_dragon_h_ani_02" example:"picture"`
// 	QANewyearDragonGAni01 string `json:"qa_newyear_dragon_g_ani_01" example:"picture"`
// 	QANewyearDragonGAni02 string `json:"qa_newyear_dragon_g_ani_02" example:"picture"`
// 	QANewyearDragonGAni03 string `json:"qa_newyear_dragon_g_ani_03" example:"picture"`
// 	QANewyearDragonCAni01 string `json:"qa_newyear_dragon_c_ani_01" example:"picture"`

// 	// 音樂
// 	QABgmStart  string `json:"qa_bgm_start" example:"picture"`  // 遊戲開始
// 	QABgmGaming string `json:"qa_bgm_gaming" example:"picture"` // 遊戲進行中
// 	QABgmEnd    string `json:"qa_bgm_end" example:"picture"`    // 遊戲結束

// 	// 搖紅包自定義
// 	RedpackClassicHPic01 string `json:"redpack_classic_h_pic_01" example:"picture"`
// 	RedpackClassicHPic02 string `json:"redpack_classic_h_pic_02" example:"picture"`
// 	RedpackClassicHPic03 string `json:"redpack_classic_h_pic_03" example:"picture"`
// 	RedpackClassicHPic04 string `json:"redpack_classic_h_pic_04" example:"picture"`
// 	RedpackClassicHPic05 string `json:"redpack_classic_h_pic_05" example:"picture"`
// 	RedpackClassicHPic06 string `json:"redpack_classic_h_pic_06" example:"picture"`
// 	RedpackClassicHPic07 string `json:"redpack_classic_h_pic_07" example:"picture"`
// 	RedpackClassicHPic08 string `json:"redpack_classic_h_pic_08" example:"picture"`
// 	RedpackClassicHPic09 string `json:"redpack_classic_h_pic_09" example:"picture"`
// 	RedpackClassicHPic10 string `json:"redpack_classic_h_pic_10" example:"picture"`
// 	RedpackClassicHPic11 string `json:"redpack_classic_h_pic_11" example:"picture"`
// 	RedpackClassicHPic12 string `json:"redpack_classic_h_pic_12" example:"picture"`
// 	RedpackClassicHPic13 string `json:"redpack_classic_h_pic_13" example:"picture"`
// 	RedpackClassicGPic01 string `json:"redpack_classic_g_pic_01" example:"picture"`
// 	RedpackClassicGPic02 string `json:"redpack_classic_g_pic_02" example:"picture"`
// 	RedpackClassicGPic03 string `json:"redpack_classic_g_pic_03" example:"picture"`
// 	RedpackClassicHAni01 string `json:"redpack_classic_h_ani_01" example:"picture"`
// 	RedpackClassicHAni02 string `json:"redpack_classic_h_ani_02" example:"picture"`
// 	RedpackClassicGAni01 string `json:"redpack_classic_g_ani_01" example:"picture"`
// 	RedpackClassicGAni02 string `json:"redpack_classic_g_ani_02" example:"picture"`
// 	RedpackClassicGAni03 string `json:"redpack_classic_g_ani_03" example:"picture"`

// 	RedpackCherryHPic01 string `json:"redpack_cherry_h_pic_01" example:"picture"`
// 	RedpackCherryHPic02 string `json:"redpack_cherry_h_pic_02" example:"picture"`
// 	RedpackCherryHPic03 string `json:"redpack_cherry_h_pic_03" example:"picture"`
// 	RedpackCherryHPic04 string `json:"redpack_cherry_h_pic_04" example:"picture"`
// 	RedpackCherryHPic05 string `json:"redpack_cherry_h_pic_05" example:"picture"`
// 	RedpackCherryHPic06 string `json:"redpack_cherry_h_pic_06" example:"picture"`
// 	RedpackCherryHPic07 string `json:"redpack_cherry_h_pic_07" example:"picture"`
// 	RedpackCherryGPic01 string `json:"redpack_cherry_g_pic_01" example:"picture"`
// 	RedpackCherryGPic02 string `json:"redpack_cherry_g_pic_02" example:"picture"`
// 	RedpackCherryHAni01 string `json:"redpack_cherry_h_ani_01" example:"picture"`
// 	RedpackCherryHAni02 string `json:"redpack_cherry_h_ani_02" example:"picture"`
// 	RedpackCherryGAni01 string `json:"redpack_cherry_g_ani_01" example:"picture"`
// 	RedpackCherryGAni02 string `json:"redpack_cherry_g_ani_02" example:"picture"`

// 	RedpackChristmasHPic01 string `json:"redpack_christmas_h_pic_01" example:"picture"`
// 	RedpackChristmasHPic02 string `json:"redpack_christmas_h_pic_02" example:"picture"`
// 	RedpackChristmasHPic03 string `json:"redpack_christmas_h_pic_03" example:"picture"`
// 	RedpackChristmasHPic04 string `json:"redpack_christmas_h_pic_04" example:"picture"`
// 	RedpackChristmasHPic05 string `json:"redpack_christmas_h_pic_05" example:"picture"`
// 	RedpackChristmasHPic06 string `json:"redpack_christmas_h_pic_06" example:"picture"`
// 	RedpackChristmasHPic07 string `json:"redpack_christmas_h_pic_07" example:"picture"`
// 	RedpackChristmasHPic08 string `json:"redpack_christmas_h_pic_08" example:"picture"`
// 	RedpackChristmasHPic09 string `json:"redpack_christmas_h_pic_09" example:"picture"`
// 	RedpackChristmasHPic10 string `json:"redpack_christmas_h_pic_10" example:"picture"`
// 	RedpackChristmasHPic11 string `json:"redpack_christmas_h_pic_11" example:"picture"`
// 	RedpackChristmasHPic12 string `json:"redpack_christmas_h_pic_12" example:"picture"`
// 	RedpackChristmasHPic13 string `json:"redpack_christmas_h_pic_13" example:"picture"`
// 	RedpackChristmasGPic01 string `json:"redpack_christmas_g_pic_01" example:"picture"`
// 	RedpackChristmasGPic02 string `json:"redpack_christmas_g_pic_02" example:"picture"`
// 	RedpackChristmasGPic03 string `json:"redpack_christmas_g_pic_03" example:"picture"`
// 	RedpackChristmasGPic04 string `json:"redpack_christmas_g_pic_04" example:"picture"`
// 	RedpackChristmasCPic01 string `json:"redpack_christmas_c_pic_01" example:"picture"`
// 	RedpackChristmasCPic02 string `json:"redpack_christmas_c_pic_02" example:"picture"`
// 	RedpackChristmasCAni01 string `json:"redpack_christmas_c_ani_01" example:"picture"`

// 	// 音樂
// 	RedpackBgmStart  string `json:"redpack_bgm_start" example:"picture"`  // 遊戲開始
// 	RedpackBgmGaming string `json:"redpack_bgm_gaming" example:"picture"` // 遊戲進行中
// 	RedpackBgmEnd    string `json:"redpack_bgm_end" example:"picture"`    // 遊戲結束

// 	// 套紅包自定義
// 	RopepackClassicHPic01 string `json:"ropepack_classic_h_pic_01" example:"picture"`
// 	RopepackClassicHPic02 string `json:"ropepack_classic_h_pic_02" example:"picture"`
// 	RopepackClassicHPic03 string `json:"ropepack_classic_h_pic_03" example:"picture"`
// 	RopepackClassicHPic04 string `json:"ropepack_classic_h_pic_04" example:"picture"`
// 	RopepackClassicHPic05 string `json:"ropepack_classic_h_pic_05" example:"picture"`
// 	RopepackClassicHPic06 string `json:"ropepack_classic_h_pic_06" example:"picture"`
// 	RopepackClassicHPic07 string `json:"ropepack_classic_h_pic_07" example:"picture"`
// 	RopepackClassicHPic08 string `json:"ropepack_classic_h_pic_08" example:"picture"`
// 	RopepackClassicHPic09 string `json:"ropepack_classic_h_pic_09" example:"picture"`
// 	RopepackClassicHPic10 string `json:"ropepack_classic_h_pic_10" example:"picture"`
// 	RopepackClassicGPic01 string `json:"ropepack_classic_g_pic_01" example:"picture"`
// 	RopepackClassicGPic02 string `json:"ropepack_classic_g_pic_02" example:"picture"`
// 	RopepackClassicGPic03 string `json:"ropepack_classic_g_pic_03" example:"picture"`
// 	RopepackClassicGPic04 string `json:"ropepack_classic_g_pic_04" example:"picture"`
// 	RopepackClassicGPic05 string `json:"ropepack_classic_g_pic_05" example:"picture"`
// 	RopepackClassicGPic06 string `json:"ropepack_classic_g_pic_06" example:"picture"`
// 	RopepackClassicHAni01 string `json:"ropepack_classic_h_ani_01" example:"picture"`
// 	RopepackClassicGAni01 string `json:"ropepack_classic_g_ani_01" example:"picture"`
// 	RopepackClassicGAni02 string `json:"ropepack_classic_g_ani_02" example:"picture"`
// 	RopepackClassicCAni01 string `json:"ropepack_classic_c_ani_01" example:"picture"`

// 	RopepackNewyearRabbitHPic01 string `json:"ropepack_newyear_rabbit_h_pic_01" example:"picture"`
// 	RopepackNewyearRabbitHPic02 string `json:"ropepack_newyear_rabbit_h_pic_02" example:"picture"`
// 	RopepackNewyearRabbitHPic03 string `json:"ropepack_newyear_rabbit_h_pic_03" example:"picture"`
// 	RopepackNewyearRabbitHPic04 string `json:"ropepack_newyear_rabbit_h_pic_04" example:"picture"`
// 	RopepackNewyearRabbitHPic05 string `json:"ropepack_newyear_rabbit_h_pic_05" example:"picture"`
// 	RopepackNewyearRabbitHPic06 string `json:"ropepack_newyear_rabbit_h_pic_06" example:"picture"`
// 	RopepackNewyearRabbitHPic07 string `json:"ropepack_newyear_rabbit_h_pic_07" example:"picture"`
// 	RopepackNewyearRabbitHPic08 string `json:"ropepack_newyear_rabbit_h_pic_08" example:"picture"`
// 	RopepackNewyearRabbitHPic09 string `json:"ropepack_newyear_rabbit_h_pic_09" example:"picture"`
// 	RopepackNewyearRabbitGPic01 string `json:"ropepack_newyear_rabbit_g_pic_01" example:"picture"`
// 	RopepackNewyearRabbitGPic02 string `json:"ropepack_newyear_rabbit_g_pic_02" example:"picture"`
// 	RopepackNewyearRabbitGPic03 string `json:"ropepack_newyear_rabbit_g_pic_03" example:"picture"`
// 	RopepackNewyearRabbitHAni01 string `json:"ropepack_newyear_rabbit_h_ani_01" example:"picture"`
// 	RopepackNewyearRabbitGAni01 string `json:"ropepack_newyear_rabbit_g_ani_01" example:"picture"`
// 	RopepackNewyearRabbitGAni02 string `json:"ropepack_newyear_rabbit_g_ani_02" example:"picture"`
// 	RopepackNewyearRabbitGAni03 string `json:"ropepack_newyear_rabbit_g_ani_03" example:"picture"`
// 	RopepackNewyearRabbitCAni01 string `json:"ropepack_newyear_rabbit_c_ani_01" example:"picture"`
// 	RopepackNewyearRabbitCAni02 string `json:"ropepack_newyear_rabbit_c_ani_02" example:"picture"`

// 	RopepackMoonfestivalHPic01 string `json:"ropepack_moonfestival_h_pic_01" example:"picture"`
// 	RopepackMoonfestivalHPic02 string `json:"ropepack_moonfestival_h_pic_02" example:"picture"`
// 	RopepackMoonfestivalHPic03 string `json:"ropepack_moonfestival_h_pic_03" example:"picture"`
// 	RopepackMoonfestivalHPic04 string `json:"ropepack_moonfestival_h_pic_04" example:"picture"`
// 	RopepackMoonfestivalHPic05 string `json:"ropepack_moonfestival_h_pic_05" example:"picture"`
// 	RopepackMoonfestivalHPic06 string `json:"ropepack_moonfestival_h_pic_06" example:"picture"`
// 	RopepackMoonfestivalHPic07 string `json:"ropepack_moonfestival_h_pic_07" example:"picture"`
// 	RopepackMoonfestivalHPic08 string `json:"ropepack_moonfestival_h_pic_08" example:"picture"`
// 	RopepackMoonfestivalHPic09 string `json:"ropepack_moonfestival_h_pic_09" example:"picture"`
// 	RopepackMoonfestivalGPic01 string `json:"ropepack_moonfestival_g_pic_01" example:"picture"`
// 	RopepackMoonfestivalGPic02 string `json:"ropepack_moonfestival_g_pic_02" example:"picture"`
// 	RopepackMoonfestivalCPic01 string `json:"ropepack_moonfestival_c_pic_01" example:"picture"`
// 	RopepackMoonfestivalHAni01 string `json:"ropepack_moonfestival_h_ani_01" example:"picture"`
// 	RopepackMoonfestivalGAni01 string `json:"ropepack_moonfestival_g_ani_01" example:"picture"`
// 	RopepackMoonfestivalGAni02 string `json:"ropepack_moonfestival_g_ani_02" example:"picture"`
// 	RopepackMoonfestivalCAni01 string `json:"ropepack_moonfestival_c_ani_01" example:"picture"`
// 	RopepackMoonfestivalCAni02 string `json:"ropepack_moonfestival_c_ani_02" example:"picture"`

// 	Ropepack3DHPic01 string `json:"ropepack_3D_h_pic_01" example:"picture"`
// 	Ropepack3DHPic02 string `json:"ropepack_3D_h_pic_02" example:"picture"`
// 	Ropepack3DHPic03 string `json:"ropepack_3D_h_pic_03" example:"picture"`
// 	Ropepack3DHPic04 string `json:"ropepack_3D_h_pic_04" example:"picture"`
// 	Ropepack3DHPic05 string `json:"ropepack_3D_h_pic_05" example:"picture"`
// 	Ropepack3DHPic06 string `json:"ropepack_3D_h_pic_06" example:"picture"`
// 	Ropepack3DHPic07 string `json:"ropepack_3D_h_pic_07" example:"picture"`
// 	Ropepack3DHPic08 string `json:"ropepack_3D_h_pic_08" example:"picture"`
// 	Ropepack3DHPic09 string `json:"ropepack_3D_h_pic_09" example:"picture"`
// 	Ropepack3DHPic10 string `json:"ropepack_3D_h_pic_10" example:"picture"`
// 	Ropepack3DHPic11 string `json:"ropepack_3D_h_pic_11" example:"picture"`
// 	Ropepack3DHPic12 string `json:"ropepack_3D_h_pic_12" example:"picture"`
// 	Ropepack3DHPic13 string `json:"ropepack_3D_h_pic_13" example:"picture"`
// 	Ropepack3DHPic14 string `json:"ropepack_3D_h_pic_14" example:"picture"`
// 	Ropepack3DHPic15 string `json:"ropepack_3D_h_pic_15" example:"picture"`
// 	Ropepack3DGPic01 string `json:"ropepack_3D_g_pic_01" example:"picture"`
// 	Ropepack3DGPic02 string `json:"ropepack_3D_g_pic_02" example:"picture"`
// 	Ropepack3DGPic03 string `json:"ropepack_3D_g_pic_03" example:"picture"`
// 	Ropepack3DGPic04 string `json:"ropepack_3D_g_pic_04" example:"picture"`
// 	Ropepack3DHAni01 string `json:"ropepack_3D_h_ani_01" example:"picture"`
// 	Ropepack3DHAni02 string `json:"ropepack_3D_h_ani_02" example:"picture"`
// 	Ropepack3DHAni03 string `json:"ropepack_3D_h_ani_03" example:"picture"`
// 	Ropepack3DGAni01 string `json:"ropepack_3D_g_ani_01" example:"picture"`
// 	Ropepack3DGAni02 string `json:"ropepack_3D_g_ani_02" example:"picture"`
// 	Ropepack3DCAni01 string `json:"ropepack_3D_c_ani_01" example:"picture"`

// 	// 音樂
// 	RopepackBgmStart  string `json:"ropepack_bgm_start" example:"picture"`  // 遊戲開始
// 	RopepackBgmGaming string `json:"ropepack_bgm_gaming" example:"picture"` // 遊戲進行中
// 	RopepackBgmEnd    string `json:"ropepack_bgm_end" example:"picture"`    // 遊戲結束

// 	// 遊戲抽獎自定義
// 	LotteryJiugonggeClassicHPic01 string `json:"lottery_jiugongge_classic_h_pic_01" example:"picture"`
// 	LotteryJiugonggeClassicHPic02 string `json:"lottery_jiugongge_classic_h_pic_02" example:"picture"`
// 	LotteryJiugonggeClassicHPic03 string `json:"lottery_jiugongge_classic_h_pic_03" example:"picture"`
// 	LotteryJiugonggeClassicHPic04 string `json:"lottery_jiugongge_classic_h_pic_04" example:"picture"`
// 	LotteryJiugonggeClassicGPic01 string `json:"lottery_jiugongge_classic_g_pic_01" example:"picture"`
// 	LotteryJiugonggeClassicGPic02 string `json:"lottery_jiugongge_classic_g_pic_02" example:"picture"`
// 	LotteryJiugonggeClassicCPic01 string `json:"lottery_jiugongge_classic_c_pic_01" example:"picture"`
// 	LotteryJiugonggeClassicCPic02 string `json:"lottery_jiugongge_classic_c_pic_02" example:"picture"`
// 	LotteryJiugonggeClassicCPic03 string `json:"lottery_jiugongge_classic_c_pic_03" example:"picture"`
// 	LotteryJiugonggeClassicCPic04 string `json:"lottery_jiugongge_classic_c_pic_04" example:"picture"`
// 	LotteryJiugonggeClassicCAni01 string `json:"lottery_jiugongge_classic_c_ani_01" example:"picture"`
// 	LotteryJiugonggeClassicCAni02 string `json:"lottery_jiugongge_classic_c_ani_02" example:"picture"`
// 	LotteryJiugonggeClassicCAni03 string `json:"lottery_jiugongge_classic_c_ani_03" example:"picture"`

// 	LotteryTurntableClassicHPic01 string `json:"lottery_turntable_classic_h_pic_01" example:"picture"`
// 	LotteryTurntableClassicHPic02 string `json:"lottery_turntable_classic_h_pic_02" example:"picture"`
// 	LotteryTurntableClassicHPic03 string `json:"lottery_turntable_classic_h_pic_03" example:"picture"`
// 	LotteryTurntableClassicHPic04 string `json:"lottery_turntable_classic_h_pic_04" example:"picture"`
// 	LotteryTurntableClassicGPic01 string `json:"lottery_turntable_classic_g_pic_01" example:"picture"`
// 	LotteryTurntableClassicGPic02 string `json:"lottery_turntable_classic_g_pic_02" example:"picture"`
// 	LotteryTurntableClassicCPic01 string `json:"lottery_turntable_classic_c_pic_01" example:"picture"`
// 	LotteryTurntableClassicCPic02 string `json:"lottery_turntable_classic_c_pic_02" example:"picture"`
// 	LotteryTurntableClassicCPic03 string `json:"lottery_turntable_classic_c_pic_03" example:"picture"`
// 	LotteryTurntableClassicCPic04 string `json:"lottery_turntable_classic_c_pic_04" example:"picture"`
// 	LotteryTurntableClassicCPic05 string `json:"lottery_turntable_classic_c_pic_05" example:"picture"`
// 	LotteryTurntableClassicCPic06 string `json:"lottery_turntable_classic_c_pic_06" example:"picture"`
// 	LotteryTurntableClassicCAni01 string `json:"lottery_turntable_classic_c_ani_01" example:"picture"`
// 	LotteryTurntableClassicCAni02 string `json:"lottery_turntable_classic_c_ani_02" example:"picture"`
// 	LotteryTurntableClassicCAni03 string `json:"lottery_turntable_classic_c_ani_03" example:"picture"`

// 	LotteryJiugonggeStarryskyHPic01 string `json:"lottery_jiugongge_starrysky_h_pic_01" example:"picture"`
// 	LotteryJiugonggeStarryskyHPic02 string `json:"lottery_jiugongge_starrysky_h_pic_02" example:"picture"`
// 	LotteryJiugonggeStarryskyHPic03 string `json:"lottery_jiugongge_starrysky_h_pic_03" example:"picture"`
// 	LotteryJiugonggeStarryskyHPic04 string `json:"lottery_jiugongge_starrysky_h_pic_04" example:"picture"`
// 	LotteryJiugonggeStarryskyHPic05 string `json:"lottery_jiugongge_starrysky_h_pic_05" example:"picture"`
// 	LotteryJiugonggeStarryskyHPic06 string `json:"lottery_jiugongge_starrysky_h_pic_06" example:"picture"`
// 	LotteryJiugonggeStarryskyHPic07 string `json:"lottery_jiugongge_starrysky_h_pic_07" example:"picture"`
// 	LotteryJiugonggeStarryskyGPic01 string `json:"lottery_jiugongge_starrysky_g_pic_01" example:"picture"`
// 	LotteryJiugonggeStarryskyGPic02 string `json:"lottery_jiugongge_starrysky_g_pic_02" example:"picture"`
// 	LotteryJiugonggeStarryskyGPic03 string `json:"lottery_jiugongge_starrysky_g_pic_03" example:"picture"`
// 	LotteryJiugonggeStarryskyGPic04 string `json:"lottery_jiugongge_starrysky_g_pic_04" example:"picture"`
// 	LotteryJiugonggeStarryskyCPic01 string `json:"lottery_jiugongge_starrysky_c_pic_01" example:"picture"`
// 	LotteryJiugonggeStarryskyCPic02 string `json:"lottery_jiugongge_starrysky_c_pic_02" example:"picture"`
// 	LotteryJiugonggeStarryskyCPic03 string `json:"lottery_jiugongge_starrysky_c_pic_03" example:"picture"`
// 	LotteryJiugonggeStarryskyCPic04 string `json:"lottery_jiugongge_starrysky_c_pic_04" example:"picture"`
// 	LotteryJiugonggeStarryskyCAni01 string `json:"lottery_jiugongge_starrysky_c_ani_01" example:"picture"`
// 	LotteryJiugonggeStarryskyCAni02 string `json:"lottery_jiugongge_starrysky_c_ani_02" example:"picture"`
// 	LotteryJiugonggeStarryskyCAni03 string `json:"lottery_jiugongge_starrysky_c_ani_03" example:"picture"`
// 	LotteryJiugonggeStarryskyCAni04 string `json:"lottery_jiugongge_starrysky_c_ani_04" example:"picture"`
// 	LotteryJiugonggeStarryskyCAni05 string `json:"lottery_jiugongge_starrysky_c_ani_05" example:"picture"`
// 	LotteryJiugonggeStarryskyCAni06 string `json:"lottery_jiugongge_starrysky_c_ani_06" example:"picture"`

// 	LotteryTurntableStarryskyHPic01 string `json:"lottery_turntable_starrysky_h_pic_01" example:"picture"`
// 	LotteryTurntableStarryskyHPic02 string `json:"lottery_turntable_starrysky_h_pic_02" example:"picture"`
// 	LotteryTurntableStarryskyHPic03 string `json:"lottery_turntable_starrysky_h_pic_03" example:"picture"`
// 	LotteryTurntableStarryskyHPic04 string `json:"lottery_turntable_starrysky_h_pic_04" example:"picture"`
// 	LotteryTurntableStarryskyHPic05 string `json:"lottery_turntable_starrysky_h_pic_05" example:"picture"`
// 	LotteryTurntableStarryskyHPic06 string `json:"lottery_turntable_starrysky_h_pic_06" example:"picture"`
// 	LotteryTurntableStarryskyHPic07 string `json:"lottery_turntable_starrysky_h_pic_07" example:"picture"`
// 	LotteryTurntableStarryskyHPic08 string `json:"lottery_turntable_starrysky_h_pic_08" example:"picture"`
// 	LotteryTurntableStarryskyGPic01 string `json:"lottery_turntable_starrysky_g_pic_01" example:"picture"`
// 	LotteryTurntableStarryskyGPic02 string `json:"lottery_turntable_starrysky_g_pic_02" example:"picture"`
// 	LotteryTurntableStarryskyGPic03 string `json:"lottery_turntable_starrysky_g_pic_03" example:"picture"`
// 	LotteryTurntableStarryskyGPic04 string `json:"lottery_turntable_starrysky_g_pic_04" example:"picture"`
// 	LotteryTurntableStarryskyGPic05 string `json:"lottery_turntable_starrysky_g_pic_05" example:"picture"`
// 	LotteryTurntableStarryskyCPic01 string `json:"lottery_turntable_starrysky_c_pic_01" example:"picture"`
// 	LotteryTurntableStarryskyCPic02 string `json:"lottery_turntable_starrysky_c_pic_02" example:"picture"`
// 	LotteryTurntableStarryskyCPic03 string `json:"lottery_turntable_starrysky_c_pic_03" example:"picture"`
// 	LotteryTurntableStarryskyCPic04 string `json:"lottery_turntable_starrysky_c_pic_04" example:"picture"`
// 	LotteryTurntableStarryskyCPic05 string `json:"lottery_turntable_starrysky_c_pic_05" example:"picture"`
// 	LotteryTurntableStarryskyCAni01 string `json:"lottery_turntable_starrysky_c_ani_01" example:"picture"`
// 	LotteryTurntableStarryskyCAni02 string `json:"lottery_turntable_starrysky_c_ani_02" example:"picture"`
// 	LotteryTurntableStarryskyCAni03 string `json:"lottery_turntable_starrysky_c_ani_03" example:"picture"`
// 	LotteryTurntableStarryskyCAni04 string `json:"lottery_turntable_starrysky_c_ani_04" example:"picture"`
// 	LotteryTurntableStarryskyCAni05 string `json:"lottery_turntable_starrysky_c_ani_05" example:"picture"`
// 	LotteryTurntableStarryskyCAni06 string `json:"lottery_turntable_starrysky_c_ani_06" example:"picture"`
// 	LotteryTurntableStarryskyCAni07 string `json:"lottery_turntable_starrysky_c_ani_07" example:"picture"`

// 	// 音樂
// 	LotteryBgmGaming string `json:"lottery_bgm_gaming" example:"picture"` // 遊戲進行中

// 	// 鑑定師自定義
// 	MonopolyClassicHPic01 string `json:"monopoly_classic_h_pic_01" example:"picture"`
// 	MonopolyClassicHPic02 string `json:"monopoly_classic_h_pic_02" example:"picture"`
// 	MonopolyClassicHPic03 string `json:"monopoly_classic_h_pic_03" example:"picture"`
// 	MonopolyClassicHPic04 string `json:"monopoly_classic_h_pic_04" example:"picture"`
// 	MonopolyClassicHPic05 string `json:"monopoly_classic_h_pic_05" example:"picture"`
// 	MonopolyClassicHPic06 string `json:"monopoly_classic_h_pic_06" example:"picture"`
// 	MonopolyClassicHPic07 string `json:"monopoly_classic_h_pic_07" example:"picture"`
// 	MonopolyClassicHPic08 string `json:"monopoly_classic_h_pic_08" example:"picture"`
// 	MonopolyClassicGPic01 string `json:"monopoly_classic_g_pic_01" example:"picture"`
// 	MonopolyClassicGPic02 string `json:"monopoly_classic_g_pic_02" example:"picture"`
// 	MonopolyClassicGPic03 string `json:"monopoly_classic_g_pic_03" example:"picture"`
// 	MonopolyClassicGPic04 string `json:"monopoly_classic_g_pic_04" example:"picture"`
// 	MonopolyClassicGPic05 string `json:"monopoly_classic_g_pic_05" example:"picture"`
// 	MonopolyClassicGPic06 string `json:"monopoly_classic_g_pic_06" example:"picture"`
// 	MonopolyClassicGPic07 string `json:"monopoly_classic_g_pic_07" example:"picture"`
// 	MonopolyClassicCPic01 string `json:"monopoly_classic_c_pic_01" example:"picture"`
// 	MonopolyClassicCPic02 string `json:"monopoly_classic_c_pic_02" example:"picture"`
// 	MonopolyClassicGAni01 string `json:"monopoly_classic_g_ani_01" example:"picture"`
// 	MonopolyClassicGAni02 string `json:"monopoly_classic_g_ani_02" example:"picture"`
// 	MonopolyClassicCAni01 string `json:"monopoly_classic_c_ani_01" example:"picture"`

// 	MonopolyRedpackHPic01 string `json:"monopoly_redpack_h_pic_01" example:"picture"`
// 	MonopolyRedpackHPic02 string `json:"monopoly_redpack_h_pic_02" example:"picture"`
// 	MonopolyRedpackHPic03 string `json:"monopoly_redpack_h_pic_03" example:"picture"`
// 	MonopolyRedpackHPic04 string `json:"monopoly_redpack_h_pic_04" example:"picture"`
// 	MonopolyRedpackHPic05 string `json:"monopoly_redpack_h_pic_05" example:"picture"`
// 	MonopolyRedpackHPic06 string `json:"monopoly_redpack_h_pic_06" example:"picture"`
// 	MonopolyRedpackHPic07 string `json:"monopoly_redpack_h_pic_07" example:"picture"`
// 	MonopolyRedpackHPic08 string `json:"monopoly_redpack_h_pic_08" example:"picture"`
// 	MonopolyRedpackHPic09 string `json:"monopoly_redpack_h_pic_09" example:"picture"`
// 	MonopolyRedpackHPic10 string `json:"monopoly_redpack_h_pic_10" example:"picture"`
// 	MonopolyRedpackHPic11 string `json:"monopoly_redpack_h_pic_11" example:"picture"`
// 	MonopolyRedpackHPic12 string `json:"monopoly_redpack_h_pic_12" example:"picture"`
// 	MonopolyRedpackHPic13 string `json:"monopoly_redpack_h_pic_13" example:"picture"`
// 	MonopolyRedpackHPic14 string `json:"monopoly_redpack_h_pic_14" example:"picture"`
// 	MonopolyRedpackHPic15 string `json:"monopoly_redpack_h_pic_15" example:"picture"`
// 	MonopolyRedpackHPic16 string `json:"monopoly_redpack_h_pic_16" example:"picture"`
// 	MonopolyRedpackGPic01 string `json:"monopoly_redpack_g_pic_01" example:"picture"`
// 	MonopolyRedpackGPic02 string `json:"monopoly_redpack_g_pic_02" example:"picture"`
// 	MonopolyRedpackGPic03 string `json:"monopoly_redpack_g_pic_03" example:"picture"`
// 	MonopolyRedpackGPic04 string `json:"monopoly_redpack_g_pic_04" example:"picture"`
// 	MonopolyRedpackGPic05 string `json:"monopoly_redpack_g_pic_05" example:"picture"`
// 	MonopolyRedpackGPic06 string `json:"monopoly_redpack_g_pic_06" example:"picture"`
// 	MonopolyRedpackGPic07 string `json:"monopoly_redpack_g_pic_07" example:"picture"`
// 	MonopolyRedpackGPic08 string `json:"monopoly_redpack_g_pic_08" example:"picture"`
// 	MonopolyRedpackGPic09 string `json:"monopoly_redpack_g_pic_09" example:"picture"`
// 	MonopolyRedpackGPic10 string `json:"monopoly_redpack_g_pic_10" example:"picture"`
// 	MonopolyRedpackCPic01 string `json:"monopoly_redpack_c_pic_01" example:"picture"`
// 	MonopolyRedpackCPic02 string `json:"monopoly_redpack_c_pic_02" example:"picture"`
// 	MonopolyRedpackCPic03 string `json:"monopoly_redpack_c_pic_03" example:"picture"`
// 	MonopolyRedpackHAni01 string `json:"monopoly_redpack_h_ani_01" example:"picture"`
// 	MonopolyRedpackHAni02 string `json:"monopoly_redpack_h_ani_02" example:"picture"`
// 	MonopolyRedpackHAni03 string `json:"monopoly_redpack_h_ani_03" example:"picture"`
// 	MonopolyRedpackGAni01 string `json:"monopoly_redpack_g_ani_01" example:"picture"`
// 	MonopolyRedpackGAni02 string `json:"monopoly_redpack_g_ani_02" example:"picture"`
// 	MonopolyRedpackCAni01 string `json:"monopoly_redpack_c_ani_01" example:"picture"`

// 	MonopolyNewyearRabbitHPic01 string `json:"monopoly_newyear_rabbit_h_pic_01" example:"picture"`
// 	MonopolyNewyearRabbitHPic02 string `json:"monopoly_newyear_rabbit_h_pic_02" example:"picture"`
// 	MonopolyNewyearRabbitHPic03 string `json:"monopoly_newyear_rabbit_h_pic_03" example:"picture"`
// 	MonopolyNewyearRabbitHPic04 string `json:"monopoly_newyear_rabbit_h_pic_04" example:"picture"`
// 	MonopolyNewyearRabbitHPic05 string `json:"monopoly_newyear_rabbit_h_pic_05" example:"picture"`
// 	MonopolyNewyearRabbitHPic06 string `json:"monopoly_newyear_rabbit_h_pic_06" example:"picture"`
// 	MonopolyNewyearRabbitHPic07 string `json:"monopoly_newyear_rabbit_h_pic_07" example:"picture"`
// 	MonopolyNewyearRabbitHPic08 string `json:"monopoly_newyear_rabbit_h_pic_08" example:"picture"`
// 	MonopolyNewyearRabbitHPic09 string `json:"monopoly_newyear_rabbit_h_pic_09" example:"picture"`
// 	MonopolyNewyearRabbitHPic10 string `json:"monopoly_newyear_rabbit_h_pic_10" example:"picture"`
// 	MonopolyNewyearRabbitHPic11 string `json:"monopoly_newyear_rabbit_h_pic_11" example:"picture"`
// 	MonopolyNewyearRabbitHPic12 string `json:"monopoly_newyear_rabbit_h_pic_12" example:"picture"`
// 	MonopolyNewyearRabbitGPic01 string `json:"monopoly_newyear_rabbit_g_pic_01" example:"picture"`
// 	MonopolyNewyearRabbitGPic02 string `json:"monopoly_newyear_rabbit_g_pic_02" example:"picture"`
// 	MonopolyNewyearRabbitGPic03 string `json:"monopoly_newyear_rabbit_g_pic_03" example:"picture"`
// 	MonopolyNewyearRabbitGPic04 string `json:"monopoly_newyear_rabbit_g_pic_04" example:"picture"`
// 	MonopolyNewyearRabbitGPic05 string `json:"monopoly_newyear_rabbit_g_pic_05" example:"picture"`
// 	MonopolyNewyearRabbitGPic06 string `json:"monopoly_newyear_rabbit_g_pic_06" example:"picture"`
// 	MonopolyNewyearRabbitGPic07 string `json:"monopoly_newyear_rabbit_g_pic_07" example:"picture"`
// 	MonopolyNewyearRabbitCPic01 string `json:"monopoly_newyear_rabbit_c_pic_01" example:"picture"`
// 	MonopolyNewyearRabbitCPic02 string `json:"monopoly_newyear_rabbit_c_pic_02" example:"picture"`
// 	MonopolyNewyearRabbitCPic03 string `json:"monopoly_newyear_rabbit_c_pic_03" example:"picture"`
// 	MonopolyNewyearRabbitHAni01 string `json:"monopoly_newyear_rabbit_h_ani_01" example:"picture"`
// 	MonopolyNewyearRabbitHAni02 string `json:"monopoly_newyear_rabbit_h_ani_02" example:"picture"`
// 	MonopolyNewyearRabbitGAni01 string `json:"monopoly_newyear_rabbit_g_ani_01" example:"picture"`
// 	MonopolyNewyearRabbitGAni02 string `json:"monopoly_newyear_rabbit_g_ani_02" example:"picture"`
// 	MonopolyNewyearRabbitCAni01 string `json:"monopoly_newyear_rabbit_c_ani_01" example:"picture"`

// 	MonopolySashimiHPic01 string `json:"monopoly_sashimi_h_pic_01" example:"picture"`
// 	MonopolySashimiHPic02 string `json:"monopoly_sashimi_h_pic_02" example:"picture"`
// 	MonopolySashimiHPic03 string `json:"monopoly_sashimi_h_pic_03" example:"picture"`
// 	MonopolySashimiHPic04 string `json:"monopoly_sashimi_h_pic_04" example:"picture"`
// 	MonopolySashimiHPic05 string `json:"monopoly_sashimi_h_pic_05" example:"picture"`
// 	MonopolySashimiGPic01 string `json:"monopoly_sashimi_g_pic_01" example:"picture"`
// 	MonopolySashimiGPic02 string `json:"monopoly_sashimi_g_pic_02" example:"picture"`
// 	MonopolySashimiGPic03 string `json:"monopoly_sashimi_g_pic_03" example:"picture"`
// 	MonopolySashimiGPic04 string `json:"monopoly_sashimi_g_pic_04" example:"picture"`
// 	MonopolySashimiGPic05 string `json:"monopoly_sashimi_g_pic_05" example:"picture"`
// 	MonopolySashimiGPic06 string `json:"monopoly_sashimi_g_pic_06" example:"picture"`
// 	MonopolySashimiCPic01 string `json:"monopoly_sashimi_c_pic_01" example:"picture"`
// 	MonopolySashimiCPic02 string `json:"monopoly_sashimi_c_pic_02" example:"picture"`
// 	MonopolySashimiHAni01 string `json:"monopoly_sashimi_h_ani_01" example:"picture"`
// 	MonopolySashimiHAni02 string `json:"monopoly_sashimi_h_ani_02" example:"picture"`
// 	MonopolySashimiGAni01 string `json:"monopoly_sashimi_g_ani_01" example:"picture"`
// 	MonopolySashimiGAni02 string `json:"monopoly_sashimi_g_ani_02" example:"picture"`
// 	MonopolySashimiCAni01 string `json:"monopoly_sashimi_c_ani_01" example:"picture"`

// 	// 音樂
// 	MonopolyBgmStart  string `json:"monopoly_bgm_start" example:"picture"`  // 遊戲開始
// 	MonopolyBgmGaming string `json:"monopoly_bgm_gaming" example:"picture"` // 遊戲進行中
// 	MonopolyBgmEnd    string `json:"monopoly_bgm_end" example:"picture"`    // 遊戲結束

// 	// 拔河遊戲自定義
// 	TugofwarClassicHPic01 string `json:"tugofwar_classic_h_pic_01" example:"picture"`
// 	TugofwarClassicHPic02 string `json:"tugofwar_classic_h_pic_02" example:"picture"`
// 	TugofwarClassicHPic03 string `json:"tugofwar_classic_h_pic_03" example:"picture"`
// 	TugofwarClassicHPic04 string `json:"tugofwar_classic_h_pic_04" example:"picture"`
// 	TugofwarClassicHPic05 string `json:"tugofwar_classic_h_pic_05" example:"picture"`
// 	TugofwarClassicHPic06 string `json:"tugofwar_classic_h_pic_06" example:"picture"`
// 	TugofwarClassicHPic07 string `json:"tugofwar_classic_h_pic_07" example:"picture"`
// 	TugofwarClassicHPic08 string `json:"tugofwar_classic_h_pic_08" example:"picture"`
// 	TugofwarClassicHPic09 string `json:"tugofwar_classic_h_pic_09" example:"picture"`
// 	TugofwarClassicHPic10 string `json:"tugofwar_classic_h_pic_10" example:"picture"`
// 	TugofwarClassicHPic11 string `json:"tugofwar_classic_h_pic_11" example:"picture"`
// 	TugofwarClassicHPic12 string `json:"tugofwar_classic_h_pic_12" example:"picture"`
// 	TugofwarClassicHPic13 string `json:"tugofwar_classic_h_pic_13" example:"picture"`
// 	TugofwarClassicHPic14 string `json:"tugofwar_classic_h_pic_14" example:"picture"`
// 	TugofwarClassicHPic15 string `json:"tugofwar_classic_h_pic_15" example:"picture"`
// 	TugofwarClassicHPic16 string `json:"tugofwar_classic_h_pic_16" example:"picture"`
// 	TugofwarClassicHPic17 string `json:"tugofwar_classic_h_pic_17" example:"picture"`
// 	TugofwarClassicHPic18 string `json:"tugofwar_classic_h_pic_18" example:"picture"`
// 	TugofwarClassicHPic19 string `json:"tugofwar_classic_h_pic_19" example:"picture"`
// 	TugofwarClassicGPic01 string `json:"tugofwar_classic_g_pic_01" example:"picture"`
// 	TugofwarClassicGPic02 string `json:"tugofwar_classic_g_pic_02" example:"picture"`
// 	TugofwarClassicGPic03 string `json:"tugofwar_classic_g_pic_03" example:"picture"`
// 	TugofwarClassicGPic04 string `json:"tugofwar_classic_g_pic_04" example:"picture"`
// 	TugofwarClassicGPic05 string `json:"tugofwar_classic_g_pic_05" example:"picture"`
// 	TugofwarClassicGPic06 string `json:"tugofwar_classic_g_pic_06" example:"picture"`
// 	TugofwarClassicGPic07 string `json:"tugofwar_classic_g_pic_07" example:"picture"`
// 	TugofwarClassicGPic08 string `json:"tugofwar_classic_g_pic_08" example:"picture"`
// 	TugofwarClassicGPic09 string `json:"tugofwar_classic_g_pic_09" example:"picture"`
// 	TugofwarClassicHAni01 string `json:"tugofwar_classic_h_ani_01" example:"picture"`
// 	TugofwarClassicHAni02 string `json:"tugofwar_classic_h_ani_02" example:"picture"`
// 	TugofwarClassicHAni03 string `json:"tugofwar_classic_h_ani_03" example:"picture"`
// 	TugofwarClassicCAni01 string `json:"tugofwar_classic_c_ani_01" example:"picture"`

// 	TugofwarSchoolHPic01 string `json:"tugofwar_school_h_pic_01" example:"picture"`
// 	TugofwarSchoolHPic02 string `json:"tugofwar_school_h_pic_02" example:"picture"`
// 	TugofwarSchoolHPic03 string `json:"tugofwar_school_h_pic_03" example:"picture"`
// 	TugofwarSchoolHPic04 string `json:"tugofwar_school_h_pic_04" example:"picture"`
// 	TugofwarSchoolHPic05 string `json:"tugofwar_school_h_pic_05" example:"picture"`
// 	TugofwarSchoolHPic06 string `json:"tugofwar_school_h_pic_06" example:"picture"`
// 	TugofwarSchoolHPic07 string `json:"tugofwar_school_h_pic_07" example:"picture"`
// 	TugofwarSchoolHPic08 string `json:"tugofwar_school_h_pic_08" example:"picture"`
// 	TugofwarSchoolHPic09 string `json:"tugofwar_school_h_pic_09" example:"picture"`
// 	TugofwarSchoolHPic10 string `json:"tugofwar_school_h_pic_10" example:"picture"`
// 	TugofwarSchoolHPic11 string `json:"tugofwar_school_h_pic_11" example:"picture"`
// 	TugofwarSchoolHPic12 string `json:"tugofwar_school_h_pic_12" example:"picture"`
// 	TugofwarSchoolHPic13 string `json:"tugofwar_school_h_pic_13" example:"picture"`
// 	TugofwarSchoolHPic14 string `json:"tugofwar_school_h_pic_14" example:"picture"`
// 	TugofwarSchoolHPic15 string `json:"tugofwar_school_h_pic_15" example:"picture"`
// 	TugofwarSchoolHPic16 string `json:"tugofwar_school_h_pic_16" example:"picture"`
// 	TugofwarSchoolHPic17 string `json:"tugofwar_school_h_pic_17" example:"picture"`
// 	TugofwarSchoolHPic18 string `json:"tugofwar_school_h_pic_18" example:"picture"`
// 	TugofwarSchoolHPic19 string `json:"tugofwar_school_h_pic_19" example:"picture"`
// 	TugofwarSchoolHPic20 string `json:"tugofwar_school_h_pic_20" example:"picture"`
// 	TugofwarSchoolHPic21 string `json:"tugofwar_school_h_pic_21" example:"picture"`
// 	TugofwarSchoolHPic22 string `json:"tugofwar_school_h_pic_22" example:"picture"`
// 	TugofwarSchoolHPic23 string `json:"tugofwar_school_h_pic_23" example:"picture"`
// 	TugofwarSchoolHPic24 string `json:"tugofwar_school_h_pic_24" example:"picture"`
// 	TugofwarSchoolHPic25 string `json:"tugofwar_school_h_pic_25" example:"picture"`
// 	TugofwarSchoolHPic26 string `json:"tugofwar_school_h_pic_26" example:"picture"`
// 	TugofwarSchoolGPic01 string `json:"tugofwar_school_g_pic_01" example:"picture"`
// 	TugofwarSchoolGPic02 string `json:"tugofwar_school_g_pic_02" example:"picture"`
// 	TugofwarSchoolGPic03 string `json:"tugofwar_school_g_pic_03" example:"picture"`
// 	TugofwarSchoolGPic04 string `json:"tugofwar_school_g_pic_04" example:"picture"`
// 	TugofwarSchoolGPic05 string `json:"tugofwar_school_g_pic_05" example:"picture"`
// 	TugofwarSchoolGPic06 string `json:"tugofwar_school_g_pic_06" example:"picture"`
// 	TugofwarSchoolGPic07 string `json:"tugofwar_school_g_pic_07" example:"picture"`
// 	TugofwarSchoolCPic01 string `json:"tugofwar_school_c_pic_01" example:"picture"`
// 	TugofwarSchoolCPic02 string `json:"tugofwar_school_c_pic_02" example:"picture"`
// 	TugofwarSchoolCPic03 string `json:"tugofwar_school_c_pic_03" example:"picture"`
// 	TugofwarSchoolCPic04 string `json:"tugofwar_school_c_pic_04" example:"picture"`
// 	TugofwarSchoolHAni01 string `json:"tugofwar_school_h_ani_01" example:"picture"`
// 	TugofwarSchoolHAni02 string `json:"tugofwar_school_h_ani_02" example:"picture"`
// 	TugofwarSchoolHAni03 string `json:"tugofwar_school_h_ani_03" example:"picture"`
// 	TugofwarSchoolHAni04 string `json:"tugofwar_school_h_ani_04" example:"picture"`
// 	TugofwarSchoolHAni05 string `json:"tugofwar_school_h_ani_05" example:"picture"`
// 	TugofwarSchoolHAni06 string `json:"tugofwar_school_h_ani_06" example:"picture"`
// 	TugofwarSchoolHAni07 string `json:"tugofwar_school_h_ani_07" example:"picture"`

// 	TugofwarChristmasHPic01 string `json:"tugofwar_christmas_h_pic_01" example:"picture"`
// 	TugofwarChristmasHPic02 string `json:"tugofwar_christmas_h_pic_02" example:"picture"`
// 	TugofwarChristmasHPic03 string `json:"tugofwar_christmas_h_pic_03" example:"picture"`
// 	TugofwarChristmasHPic04 string `json:"tugofwar_christmas_h_pic_04" example:"picture"`
// 	TugofwarChristmasHPic05 string `json:"tugofwar_christmas_h_pic_05" example:"picture"`
// 	TugofwarChristmasHPic06 string `json:"tugofwar_christmas_h_pic_06" example:"picture"`
// 	TugofwarChristmasHPic07 string `json:"tugofwar_christmas_h_pic_07" example:"picture"`
// 	TugofwarChristmasHPic08 string `json:"tugofwar_christmas_h_pic_08" example:"picture"`
// 	TugofwarChristmasHPic09 string `json:"tugofwar_christmas_h_pic_09" example:"picture"`
// 	TugofwarChristmasHPic10 string `json:"tugofwar_christmas_h_pic_10" example:"picture"`
// 	TugofwarChristmasHPic11 string `json:"tugofwar_christmas_h_pic_11" example:"picture"`
// 	TugofwarChristmasHPic12 string `json:"tugofwar_christmas_h_pic_12" example:"picture"`
// 	TugofwarChristmasHPic13 string `json:"tugofwar_christmas_h_pic_13" example:"picture"`
// 	TugofwarChristmasHPic14 string `json:"tugofwar_christmas_h_pic_14" example:"picture"`
// 	TugofwarChristmasHPic15 string `json:"tugofwar_christmas_h_pic_15" example:"picture"`
// 	TugofwarChristmasHPic16 string `json:"tugofwar_christmas_h_pic_16" example:"picture"`
// 	TugofwarChristmasHPic17 string `json:"tugofwar_christmas_h_pic_17" example:"picture"`
// 	TugofwarChristmasHPic18 string `json:"tugofwar_christmas_h_pic_18" example:"picture"`
// 	TugofwarChristmasHPic19 string `json:"tugofwar_christmas_h_pic_19" example:"picture"`
// 	TugofwarChristmasHPic20 string `json:"tugofwar_christmas_h_pic_20" example:"picture"`
// 	TugofwarChristmasHPic21 string `json:"tugofwar_christmas_h_pic_21" example:"picture"`
// 	TugofwarChristmasGPic01 string `json:"tugofwar_christmas_g_pic_01" example:"picture"`
// 	TugofwarChristmasGPic02 string `json:"tugofwar_christmas_g_pic_02" example:"picture"`
// 	TugofwarChristmasGPic03 string `json:"tugofwar_christmas_g_pic_03" example:"picture"`
// 	TugofwarChristmasGPic04 string `json:"tugofwar_christmas_g_pic_04" example:"picture"`
// 	TugofwarChristmasGPic05 string `json:"tugofwar_christmas_g_pic_05" example:"picture"`
// 	TugofwarChristmasGPic06 string `json:"tugofwar_christmas_g_pic_06" example:"picture"`
// 	TugofwarChristmasCPic01 string `json:"tugofwar_christmas_c_pic_01" example:"picture"`
// 	TugofwarChristmasCPic02 string `json:"tugofwar_christmas_c_pic_02" example:"picture"`
// 	TugofwarChristmasCPic03 string `json:"tugofwar_christmas_c_pic_03" example:"picture"`
// 	TugofwarChristmasCPic04 string `json:"tugofwar_christmas_c_pic_04" example:"picture"`
// 	TugofwarChristmasHAni01 string `json:"tugofwar_christmas_h_ani_01" example:"picture"`
// 	TugofwarChristmasHAni02 string `json:"tugofwar_christmas_h_ani_02" example:"picture"`
// 	TugofwarChristmasHAni03 string `json:"tugofwar_christmas_h_ani_03" example:"picture"`
// 	TugofwarChristmasCAni01 string `json:"tugofwar_christmas_c_ani_01" example:"picture"`
// 	TugofwarChristmasCAni02 string `json:"tugofwar_christmas_c_ani_02" example:"picture"`

// 	// 音樂
// 	TugofwarBgmStart  string `json:"tugofwar_bgm_start" example:"picture"`  // 遊戲開始
// 	TugofwarBgmGaming string `json:"tugofwar_bgm_gaming" example:"picture"` // 遊戲進行中
// 	TugofwarBgmEnd    string `json:"tugofwar_bgm_end" example:"picture"`    // 遊戲結束

// 	// 賓果遊戲自定義
// 	BingoClassicHPic01 string `json:"bingo_classic_h_pic_01" example:"picture"`
// 	BingoClassicHPic02 string `json:"bingo_classic_h_pic_02" example:"picture"`
// 	BingoClassicHPic03 string `json:"bingo_classic_h_pic_03" example:"picture"`
// 	BingoClassicHPic04 string `json:"bingo_classic_h_pic_04" example:"picture"`
// 	BingoClassicHPic05 string `json:"bingo_classic_h_pic_05" example:"picture"`
// 	BingoClassicHPic06 string `json:"bingo_classic_h_pic_06" example:"picture"`
// 	BingoClassicHPic07 string `json:"bingo_classic_h_pic_07" example:"picture"`
// 	BingoClassicHPic08 string `json:"bingo_classic_h_pic_08" example:"picture"`
// 	BingoClassicHPic09 string `json:"bingo_classic_h_pic_09" example:"picture"`
// 	BingoClassicHPic10 string `json:"bingo_classic_h_pic_10" example:"picture"`
// 	BingoClassicHPic11 string `json:"bingo_classic_h_pic_11" example:"picture"`
// 	BingoClassicHPic12 string `json:"bingo_classic_h_pic_12" example:"picture"`
// 	BingoClassicHPic13 string `json:"bingo_classic_h_pic_13" example:"picture"`
// 	BingoClassicHPic14 string `json:"bingo_classic_h_pic_14" example:"picture"`
// 	BingoClassicHPic15 string `json:"bingo_classic_h_pic_15" example:"picture"`
// 	BingoClassicHPic16 string `json:"bingo_classic_h_pic_16" example:"picture"`
// 	BingoClassicGPic01 string `json:"bingo_classic_g_pic_01" example:"picture"`
// 	BingoClassicGPic02 string `json:"bingo_classic_g_pic_02" example:"picture"`
// 	BingoClassicGPic03 string `json:"bingo_classic_g_pic_03" example:"picture"`
// 	BingoClassicGPic04 string `json:"bingo_classic_g_pic_04" example:"picture"`
// 	BingoClassicGPic05 string `json:"bingo_classic_g_pic_05" example:"picture"`
// 	BingoClassicGPic06 string `json:"bingo_classic_g_pic_06" example:"picture"`
// 	BingoClassicGPic07 string `json:"bingo_classic_g_pic_07" example:"picture"`
// 	BingoClassicGPic08 string `json:"bingo_classic_g_pic_08" example:"picture"`
// 	BingoClassicCPic01 string `json:"bingo_classic_c_pic_01" example:"picture"`
// 	BingoClassicCPic02 string `json:"bingo_classic_c_pic_02" example:"picture"`
// 	BingoClassicCPic03 string `json:"bingo_classic_c_pic_03" example:"picture"`
// 	BingoClassicCPic04 string `json:"bingo_classic_c_pic_04" example:"picture"`
// 	// BingoClassicCPic05 string `json:"bingo_classic_c_pic_05" example:"picture"`
// 	BingoClassicHAni01 string `json:"bingo_classic_h_ani_01" example:"picture"`
// 	BingoClassicGAni01 string `json:"bingo_classic_g_ani_01" example:"picture"`
// 	BingoClassicCAni01 string `json:"bingo_classic_c_ani_01" example:"picture"`
// 	BingoClassicCAni02 string `json:"bingo_classic_c_ani_02" example:"picture"`

// 	BingoNewyearDragonHPic01 string `json:"bingo_newyear_dragon_h_pic_01" example:"picture"`
// 	BingoNewyearDragonHPic02 string `json:"bingo_newyear_dragon_h_pic_02" example:"picture"`
// 	BingoNewyearDragonHPic03 string `json:"bingo_newyear_dragon_h_pic_03" example:"picture"`
// 	BingoNewyearDragonHPic04 string `json:"bingo_newyear_dragon_h_pic_04" example:"picture"`
// 	BingoNewyearDragonHPic05 string `json:"bingo_newyear_dragon_h_pic_05" example:"picture"`
// 	BingoNewyearDragonHPic06 string `json:"bingo_newyear_dragon_h_pic_06" example:"picture"`
// 	BingoNewyearDragonHPic07 string `json:"bingo_newyear_dragon_h_pic_07" example:"picture"`
// 	BingoNewyearDragonHPic08 string `json:"bingo_newyear_dragon_h_pic_08" example:"picture"`
// 	BingoNewyearDragonHPic09 string `json:"bingo_newyear_dragon_h_pic_09" example:"picture"`
// 	BingoNewyearDragonHPic10 string `json:"bingo_newyear_dragon_h_pic_10" example:"picture"`
// 	BingoNewyearDragonHPic11 string `json:"bingo_newyear_dragon_h_pic_11" example:"picture"`
// 	BingoNewyearDragonHPic12 string `json:"bingo_newyear_dragon_h_pic_12" example:"picture"`
// 	BingoNewyearDragonHPic13 string `json:"bingo_newyear_dragon_h_pic_13" example:"picture"`
// 	BingoNewyearDragonHPic14 string `json:"bingo_newyear_dragon_h_pic_14" example:"picture"`
// 	// BingoNewyearDragonHPic15 string `json:"bingo_newyear_dragon_h_pic_15" example:"picture"`
// 	BingoNewyearDragonHPic16 string `json:"bingo_newyear_dragon_h_pic_16" example:"picture"`
// 	BingoNewyearDragonHPic17 string `json:"bingo_newyear_dragon_h_pic_17" example:"picture"`
// 	BingoNewyearDragonHPic18 string `json:"bingo_newyear_dragon_h_pic_18" example:"picture"`
// 	BingoNewyearDragonHPic19 string `json:"bingo_newyear_dragon_h_pic_19" example:"picture"`
// 	BingoNewyearDragonHPic20 string `json:"bingo_newyear_dragon_h_pic_20" example:"picture"`
// 	BingoNewyearDragonHPic21 string `json:"bingo_newyear_dragon_h_pic_21" example:"picture"`
// 	BingoNewyearDragonHPic22 string `json:"bingo_newyear_dragon_h_pic_22" example:"picture"`
// 	BingoNewyearDragonGPic01 string `json:"bingo_newyear_dragon_g_pic_01" example:"picture"`
// 	BingoNewyearDragonGPic02 string `json:"bingo_newyear_dragon_g_pic_02" example:"picture"`
// 	BingoNewyearDragonGPic03 string `json:"bingo_newyear_dragon_g_pic_03" example:"picture"`
// 	BingoNewyearDragonGPic04 string `json:"bingo_newyear_dragon_g_pic_04" example:"picture"`
// 	BingoNewyearDragonGPic05 string `json:"bingo_newyear_dragon_g_pic_05" example:"picture"`
// 	BingoNewyearDragonGPic06 string `json:"bingo_newyear_dragon_g_pic_06" example:"picture"`
// 	BingoNewyearDragonGPic07 string `json:"bingo_newyear_dragon_g_pic_07" example:"picture"`
// 	BingoNewyearDragonGPic08 string `json:"bingo_newyear_dragon_g_pic_08" example:"picture"`
// 	BingoNewyearDragonCPic01 string `json:"bingo_newyear_dragon_c_pic_01" example:"picture"`
// 	BingoNewyearDragonCPic02 string `json:"bingo_newyear_dragon_c_pic_02" example:"picture"`
// 	BingoNewyearDragonCPic03 string `json:"bingo_newyear_dragon_c_pic_03" example:"picture"`
// 	BingoNewyearDragonHAni01 string `json:"bingo_newyear_dragon_h_ani_01" example:"picture"`
// 	BingoNewyearDragonGAni01 string `json:"bingo_newyear_dragon_g_ani_01" example:"picture"`
// 	BingoNewyearDragonCAni01 string `json:"bingo_newyear_dragon_c_ani_01" example:"picture"`
// 	BingoNewyearDragonCAni02 string `json:"bingo_newyear_dragon_c_ani_02" example:"picture"`
// 	BingoNewyearDragonCAni03 string `json:"bingo_newyear_dragon_c_ani_03" example:"picture"`

// 	BingoCherryHPic01 string `json:"bingo_cherry_h_pic_01" example:"picture"`
// 	BingoCherryHPic02 string `json:"bingo_cherry_h_pic_02" example:"picture"`
// 	BingoCherryHPic03 string `json:"bingo_cherry_h_pic_03" example:"picture"`
// 	BingoCherryHPic04 string `json:"bingo_cherry_h_pic_04" example:"picture"`
// 	BingoCherryHPic05 string `json:"bingo_cherry_h_pic_05" example:"picture"`
// 	BingoCherryHPic06 string `json:"bingo_cherry_h_pic_06" example:"picture"`
// 	BingoCherryHPic07 string `json:"bingo_cherry_h_pic_07" example:"picture"`
// 	BingoCherryHPic08 string `json:"bingo_cherry_h_pic_08" example:"picture"`
// 	BingoCherryHPic09 string `json:"bingo_cherry_h_pic_09" example:"picture"`
// 	BingoCherryHPic10 string `json:"bingo_cherry_h_pic_10" example:"picture"`
// 	BingoCherryHPic11 string `json:"bingo_cherry_h_pic_11" example:"picture"`
// 	BingoCherryHPic12 string `json:"bingo_cherry_h_pic_12" example:"picture"`
// 	// BingoCherryHPic13 string `json:"bingo_cherry_h_pic_13" example:"picture"`
// 	BingoCherryHPic14 string `json:"bingo_cherry_h_pic_14" example:"picture"`
// 	BingoCherryHPic15 string `json:"bingo_cherry_h_pic_15" example:"picture"`
// 	// BingoCherryHPic16 string `json:"bingo_cherry_h_pic_16" example:"picture"`
// 	BingoCherryHPic17 string `json:"bingo_cherry_h_pic_17" example:"picture"`
// 	BingoCherryHPic18 string `json:"bingo_cherry_h_pic_18" example:"picture"`
// 	BingoCherryHPic19 string `json:"bingo_cherry_h_pic_19" example:"picture"`
// 	BingoCherryGPic01 string `json:"bingo_cherry_g_pic_01" example:"picture"`
// 	BingoCherryGPic02 string `json:"bingo_cherry_g_pic_02" example:"picture"`
// 	BingoCherryGPic03 string `json:"bingo_cherry_g_pic_03" example:"picture"`
// 	BingoCherryGPic04 string `json:"bingo_cherry_g_pic_04" example:"picture"`
// 	BingoCherryGPic05 string `json:"bingo_cherry_g_pic_05" example:"picture"`
// 	BingoCherryGPic06 string `json:"bingo_cherry_g_pic_06" example:"picture"`
// 	BingoCherryCPic01 string `json:"bingo_cherry_c_pic_01" example:"picture"`
// 	BingoCherryCPic02 string `json:"bingo_cherry_c_pic_02" example:"picture"`
// 	BingoCherryCPic03 string `json:"bingo_cherry_c_pic_03" example:"picture"`
// 	BingoCherryCPic04 string `json:"bingo_cherry_c_pic_04" example:"picture"`
// 	// BingoCherryHAni01 string `json:"bingo_cherry_h_ani_01" example:"picture"`
// 	BingoCherryHAni02 string `json:"bingo_cherry_h_ani_02" example:"picture"`
// 	BingoCherryHAni03 string `json:"bingo_cherry_h_ani_03" example:"picture"`
// 	BingoCherryGAni01 string `json:"bingo_cherry_g_ani_01" example:"picture"`
// 	BingoCherryGAni02 string `json:"bingo_cherry_g_ani_02" example:"picture"`
// 	BingoCherryCAni01 string `json:"bingo_cherry_c_ani_01" example:"picture"`
// 	BingoCherryCAni02 string `json:"bingo_cherry_c_ani_02" example:"picture"`

// 	// 音樂
// 	BingoBgmStart  string `json:"bingo_bgm_start" example:"picture"`  // 遊戲開始
// 	BingoBgmGaming string `json:"bingo_bgm_gaming" example:"picture"` // 遊戲進行中
// 	BingoBgmEnd    string `json:"bingo_bgm_end" example:"picture"`    // 遊戲結束

// 	// 扭蛋機自定義
// 	GachaMachineClassicHPic02 string `json:"3d_gacha_machine_classic_h_pic_02" example:"picture"`
// 	GachaMachineClassicHPic03 string `json:"3d_gacha_machine_classic_h_pic_03" example:"picture"`
// 	GachaMachineClassicHPic04 string `json:"3d_gacha_machine_classic_h_pic_04" example:"picture"`
// 	GachaMachineClassicHPic05 string `json:"3d_gacha_machine_classic_h_pic_05" example:"picture"`
// 	GachaMachineClassicGPic01 string `json:"3d_gacha_machine_classic_g_pic_01" example:"picture"`
// 	GachaMachineClassicGPic02 string `json:"3d_gacha_machine_classic_g_pic_02" example:"picture"`
// 	GachaMachineClassicGPic03 string `json:"3d_gacha_machine_classic_g_pic_03" example:"picture"`
// 	GachaMachineClassicGPic04 string `json:"3d_gacha_machine_classic_g_pic_04" example:"picture"`
// 	GachaMachineClassicGPic05 string `json:"3d_gacha_machine_classic_g_pic_05" example:"picture"`
// 	GachaMachineClassicCPic01 string `json:"3d_gacha_machine_classic_c_pic_01" example:"picture"`

// 	// 音樂
// 	GachaMachineBgmGaming string `json:"3d_gacha_machine_bgm_gaming" example:"picture"`

// 	// 投票自定義
// 	VoteClassicHPic01 string `json:"vote_classic_h_pic_01" example:"picture"`
// 	VoteClassicHPic02 string `json:"vote_classic_h_pic_02" example:"picture"`
// 	VoteClassicHPic03 string `json:"vote_classic_h_pic_03" example:"picture"`
// 	VoteClassicHPic04 string `json:"vote_classic_h_pic_04" example:"picture"`
// 	VoteClassicHPic05 string `json:"vote_classic_h_pic_05" example:"picture"`
// 	VoteClassicHPic06 string `json:"vote_classic_h_pic_06" example:"picture"`
// 	VoteClassicHPic07 string `json:"vote_classic_h_pic_07" example:"picture"`
// 	VoteClassicHPic08 string `json:"vote_classic_h_pic_08" example:"picture"`
// 	VoteClassicHPic09 string `json:"vote_classic_h_pic_09" example:"picture"`
// 	VoteClassicHPic10 string `json:"vote_classic_h_pic_10" example:"picture"`
// 	VoteClassicHPic11 string `json:"vote_classic_h_pic_11" example:"picture"`
// 	// VoteClassicHPic12 string `json:"vote_classic_h_pic_12" example:"picture"`
// 	VoteClassicHPic13 string `json:"vote_classic_h_pic_13" example:"picture"`
// 	VoteClassicHPic14 string `json:"vote_classic_h_pic_14" example:"picture"`
// 	VoteClassicHPic15 string `json:"vote_classic_h_pic_15" example:"picture"`
// 	VoteClassicHPic16 string `json:"vote_classic_h_pic_16" example:"picture"`
// 	VoteClassicHPic17 string `json:"vote_classic_h_pic_17" example:"picture"`
// 	VoteClassicHPic18 string `json:"vote_classic_h_pic_18" example:"picture"`
// 	VoteClassicHPic19 string `json:"vote_classic_h_pic_19" example:"picture"`
// 	VoteClassicHPic20 string `json:"vote_classic_h_pic_20" example:"picture"`
// 	VoteClassicHPic21 string `json:"vote_classic_h_pic_21" example:"picture"`
// 	// VoteClassicHPic22 string `json:"vote_classic_h_pic_22" example:"picture"`
// 	VoteClassicHPic23 string `json:"vote_classic_h_pic_23" example:"picture"`
// 	VoteClassicHPic24 string `json:"vote_classic_h_pic_24" example:"picture"`
// 	VoteClassicHPic25 string `json:"vote_classic_h_pic_25" example:"picture"`
// 	VoteClassicHPic26 string `json:"vote_classic_h_pic_26" example:"picture"`
// 	VoteClassicHPic27 string `json:"vote_classic_h_pic_27" example:"picture"`
// 	VoteClassicHPic28 string `json:"vote_classic_h_pic_28" example:"picture"`
// 	VoteClassicHPic29 string `json:"vote_classic_h_pic_29" example:"picture"`
// 	VoteClassicHPic30 string `json:"vote_classic_h_pic_30" example:"picture"`
// 	VoteClassicHPic31 string `json:"vote_classic_h_pic_31" example:"picture"`
// 	VoteClassicHPic32 string `json:"vote_classic_h_pic_32" example:"picture"`
// 	VoteClassicHPic33 string `json:"vote_classic_h_pic_33" example:"picture"`
// 	VoteClassicHPic34 string `json:"vote_classic_h_pic_34" example:"picture"`
// 	VoteClassicHPic35 string `json:"vote_classic_h_pic_35" example:"picture"`
// 	VoteClassicHPic36 string `json:"vote_classic_h_pic_36" example:"picture"`
// 	VoteClassicHPic37 string `json:"vote_classic_h_pic_37" example:"picture"`
// 	VoteClassicGPic01 string `json:"vote_classic_g_pic_01" example:"picture"`
// 	VoteClassicGPic02 string `json:"vote_classic_g_pic_02" example:"picture"`
// 	VoteClassicGPic03 string `json:"vote_classic_g_pic_03" example:"picture"`
// 	VoteClassicGPic04 string `json:"vote_classic_g_pic_04" example:"picture"`
// 	VoteClassicGPic05 string `json:"vote_classic_g_pic_05" example:"picture"`
// 	VoteClassicGPic06 string `json:"vote_classic_g_pic_06" example:"picture"`
// 	VoteClassicGPic07 string `json:"vote_classic_g_pic_07" example:"picture"`
// 	VoteClassicCPic01 string `json:"vote_classic_c_pic_01" example:"picture"`
// 	VoteClassicCPic02 string `json:"vote_classic_c_pic_02" example:"picture"`
// 	VoteClassicCPic03 string `json:"vote_classic_c_pic_03" example:"picture"`
// 	VoteClassicCPic04 string `json:"vote_classic_c_pic_04" example:"picture"`
// 	// 音樂
// 	VoteBgmGaming string `json:"vote_bgm_gaming" example:"picture"`

// 	// 快問快答
// 	QA1 string `json:"qa_1" example:"qa"`
// 	// QA1Picture  string `json:"qa_1_picture" example:"picture"`
// 	QA1Options string `json:"qa_1_options" example:"A&&&B&&&C&&&D"`
// 	QA1Answer  string `json:"qa_1_answer" example:"1"`
// 	QA1Score   string `json:"qa_1_score" example:"10"`
// 	QA2        string `json:"qa_2" example:"qa"`
// 	// QA2Picture  string `json:"qa_2_picture" example:"picture"`
// 	QA2Options string `json:"qa_2_options" example:"A&&&B&&&C&&&D"`
// 	QA2Answer  string `json:"qa_2_pic_answer" example:"1"`
// 	QA2Score   string `json:"qa_2_score" example:"10"`
// 	QA3        string `json:"qa_3" example:"qa"`
// 	// QA3Picture  string `json:"qa_3_picture" example:"picture"`
// 	QA3Options string `json:"qa_3_options" example:"A&&&B&&&C&&&D"`
// 	QA3Answer  string `json:"qa_3_pic_answer" example:"1"`
// 	QA3Score   string `json:"qa_3_score" example:"10"`
// 	QA4        string `json:"qa_4" example:"qa"`
// 	// QA4Picture  string `json:"qa_4_picture" example:"picture"`
// 	QA4Options string `json:"qa_4_options" example:"A&&&B&&&C&&&D"`
// 	QA4Answer  string `json:"qa_4_pic_answer" example:"1"`
// 	QA4Score   string `json:"qa_4_score" example:"10"`
// 	QA5        string `json:"qa_5" example:"qa"`
// 	// QA5Picture  string `json:"qa_5_picture" example:"picture"`
// 	QA5Options string `json:"qa_5_options" example:"A&&&B&&&C&&&D"`
// 	QA5Answer  string `json:"qa_5_pic_answer" example:"1"`
// 	QA5Score   string `json:"qa_5_score" example:"10"`
// 	QA6        string `json:"qa_6" example:"qa"`
// 	// QA6Picture  string `json:"qa_6_picture" example:"picture"`
// 	QA6Options string `json:"qa_6_options" example:"A&&&B&&&C&&&D"`
// 	QA6Answer  string `json:"qa_6_pic_answer" example:"1"`
// 	QA6Score   string `json:"qa_6_score" example:"10"`
// 	QA7        string `json:"qa_7" example:"qa"`
// 	// QA7Picture  string `json:"qa_7_picture" example:"picture"`
// 	QA7Options string `json:"qa_7_options" example:"A&&&B&&&C&&&D"`
// 	QA7Answer  string `json:"qa_7_pic_answer" example:"1"`
// 	QA7Score   string `json:"qa_7_score" example:"10"`
// 	QA8        string `json:"qa_8" example:"qa"`
// 	// QA8Picture  string `json:"qa_8_picture" example:"picture"`
// 	QA8Options string `json:"qa_8_options" example:"A&&&B&&&C&&&D"`
// 	QA8Answer  string `json:"qa_8_pic_answer" example:"1"`
// 	QA8Score   string `json:"qa_8_score" example:"10"`
// 	QA9        string `json:"qa_9" example:"qa"`
// 	// QA9Picture  string `json:"qa_9_picture" example:"picture"`
// 	QA9Options string `json:"qa_9_options" example:"A&&&B&&&C&&&D"`
// 	QA9Answer  string `json:"qa_9_pic_answer" example:"1"`
// 	QA9Score   string `json:"qa_9_score" example:"10"`
// 	QA10       string `json:"qa_10" example:"qa"`
// 	// QA10Picture string `json:"qa_10_picture" example:"picture"`
// 	QA10Options string `json:"qa_10_options" example:"A&&&B&&&C&&&D"`
// 	QA10Answer  string `json:"qa_10_pic_answer" example:"1"`
// 	QA10Score   string `json:"qa_10_score" example:"10"`
// 	QA11        string `json:"qa_11" example:"qa"`
// 	// QA11Picture string `json:"qa_11_picture" example:"picture"`
// 	QA11Options string `json:"qa_11_options" example:"A&&&B&&&C&&&D"`
// 	QA11Answer  string `json:"qa_11_pic_answer" example:"1"`
// 	QA11Score   string `json:"qa_11_score" example:"10"`
// 	QA12        string `json:"qa_12" example:"qa"`
// 	// QA12Picture string `json:"qa_12_picture" example:"picture"`
// 	QA12Options string `json:"qa_12_options" example:"A&&&B&&&C&&&D"`
// 	QA12Answer  string `json:"qa_12_pic_answer" example:"1"`
// 	QA12Score   string `json:"qa_12_score" example:"10"`
// 	QA13        string `json:"qa_13" example:"qa"`
// 	// QA13Picture string `json:"qa_13_picture" example:"picture"`
// 	QA13Options string `json:"qa_13_options" example:"A&&&B&&&C&&&D"`
// 	QA13Answer  string `json:"qa_13_pic_answer" example:"1"`
// 	QA13Score   string `json:"qa_13_score" example:"10"`
// 	QA14        string `json:"qa_14" example:"qa"`
// 	// QA14Picture string `json:"qa_14_picture" example:"picture"`
// 	QA14Options string `json:"qa_14_options" example:"A&&&B&&&C&&&D"`
// 	QA14Answer  string `json:"qa_14_pic_answer" example:"1"`
// 	QA14Score   string `json:"qa_14_score" example:"10"`
// 	QA15        string `json:"qa_15" example:"qa"`
// 	// QA15Picture string `json:"qa_15_picture" example:"picture"`
// 	QA15Options string `json:"qa_15_options" example:"A&&&B&&&C&&&D"`
// 	QA15Answer  string `json:"qa_15_pic_answer" example:"1"`
// 	QA15Score   string `json:"qa_15_score" example:"10"`
// 	QA16        string `json:"qa_16" example:"qa"`
// 	// QA16Picture string `json:"qa_16_picture" example:"picture"`
// 	QA16Options string `json:"qa_16_options" example:"A&&&B&&&C&&&D"`
// 	QA16Answer  string `json:"qa_16_pic_answer" example:"1"`
// 	QA16Score   string `json:"qa_16_score" example:"10"`
// 	QA17        string `json:"qa_17" example:"qa"`
// 	// QA17Picture string `json:"qa_17_picture" example:"picture"`
// 	QA17Options string `json:"qa_17_options" example:"A&&&B&&&C&&&D"`
// 	QA17Answer  string `json:"qa_17_pic_answer" example:"1"`
// 	QA17Score   string `json:"qa_17_score" example:"10"`
// 	QA18        string `json:"qa_18" example:"qa"`
// 	// QA18Picture string `json:"qa_18_picture" example:"picture"`
// 	QA18Options string `json:"qa_18_options" example:"A&&&B&&&C&&&D"`
// 	QA18Answer  string `json:"qa_18_pic_answer" example:"1"`
// 	QA18Score   string `json:"qa_18_score" example:"10"`
// 	QA19        string `json:"qa_19" example:"qa"`
// 	// QA19Picture string `json:"qa_19_picture" example:"picture"`
// 	QA19Options string `json:"qa_19_options" example:"A&&&B&&&C&&&D"`
// 	QA19Answer  string `json:"qa_19_pic_answer" example:"1"`
// 	QA19Score   string `json:"qa_19_score" example:"10"`
// 	QA20        string `json:"qa_20" example:"qa"`
// 	// QA20Picture string `json:"qa_20_picture" example:"picture"`
// 	QA20Options string `json:"qa_20_options" example:"A&&&B&&&C&&&D"`
// 	QA20Answer  string `json:"qa_20_pic_answer" example:"1"`
// 	QA20Score   string `json:"qa_20_score" example:"10"`
// 	TotalQA     string `json:"total_qa" example:"1"`   // 總題數
// 	QASecond    string `json:"qa_second" example:"30"` // 題目顯示秒數

// 	Token string `json:"token" example:"token"`
// }

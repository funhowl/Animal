package game

import (
	"Animal/common"
	"Animal/user"
)

var BeginingChess = [63]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

//16个棋子 象>狮>虎>豹>狼>狗>猫>鼠  9 8 7 6 5 4 3 2
//地型   陆地 水  陷阱 兽穴  0  100 200 300

func init() {
	//地形 水
	BeginingChess[12] = 100
	BeginingChess[13] = 100
	BeginingChess[14] = 100

	BeginingChess[21] = 100
	BeginingChess[22] = 100
	BeginingChess[23] = 100

	BeginingChess[39] = 100
	BeginingChess[40] = 100
	BeginingChess[41] = 100

	BeginingChess[48] = 100
	BeginingChess[49] = 100
	BeginingChess[50] = 100

	//地形 陷阱
	BeginingChess[18] = 200
	BeginingChess[28] = 200
	BeginingChess[36] = 200

	BeginingChess[26] = 200
	BeginingChess[34] = 200
	BeginingChess[44] = 200

	//地形 兽穴
	BeginingChess[27] = 300
	BeginingChess[35] = 300

	//棋子
	BeginingChess[0] = 7
	BeginingChess[2] = 9
	BeginingChess[10] = 3
	BeginingChess[20] = 5
	BeginingChess[38] = 6
	BeginingChess[46] = 4
	BeginingChess[54] = 8
	BeginingChess[56] = 2

	BeginingChess[62] = 17
	BeginingChess[60] = 19
	BeginingChess[52] = 13
	BeginingChess[42] = 15
	BeginingChess[24] = 16
	BeginingChess[16] = 14
	BeginingChess[8] = 18
	BeginingChess[6] = 12

}

type gameMsg struct {
	Key     float64
	GameMsg *common.GameRow
}

func startgame(roomS *room) {
	var (
		left  string
		right string
		trun  int // 0代表left  1代表right
	)
	left = roomS.Left
	right = roomS.Right

	trun = common.IntnBytime(1)

	SendGameMsg(left, trun, BeginingChess)
	SendGameMsg(right, trun, BeginingChess)
}

func SendGameMsg(accid string, trun int, chess [63]int) {
	var (
		urow *common.UserRow
		err  error
		ans  gameMsg
	)
	if urow, err = user.GetRowById(accid); err != nil {
		panic(err.Error())
	}
	urow.GameMsg.Chess = chess
	urow.GameMsg.Trun = trun
	user.UserRowData[accid] = urow

	ans.Key = common.S2Cgamemsg
	ans.GameMsg = urow.GameMsg

	if _, err = user.Send(accid, ans); err != nil {
		panic(err.Error())
	}
}

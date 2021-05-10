package game

import (
	"Animal/common"
	"Animal/user"
)

var BeginingChess = [63]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

//16个棋子 象>狮>虎>豹>狼>狗>猫>鼠  9 8 7 6 5 4 3 2
//地型   陆地 水  陷阱 兽穴  0  100 200 300

func init() {

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
	left = roomS.left
	right = roomS.right

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
	urow.GameMsg.Trun = trun
	user.UserRowData[accid] = urow

	ans.Key = common.C2Sgamemsg
	ans.GameMsg = urow.GameMsg

	if _, err = user.Send(accid, ans); err != nil {
		panic(err.Error())
	}
}

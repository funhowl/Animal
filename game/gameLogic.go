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

	user.RegisterEvent(common.C2Sgamechange, new(roomessage).gamechange) //  棋子移动指令
}

func (*roomessage) gamechange(conn *user.Connection, param map[string]interface{}) {

	var (
		err  error
		urow *common.UserRow
	)

	if urow, err = user.GetRowById(conn.Accid); err != nil {
		panic(err.Error())
	}

	turn := urow.GameMsg.Trun
	if !(turn == 0 && urow.GameMsg.Left == conn.Accid) || !(turn == 1 && urow.GameMsg.Right == conn.Accid) {
		conn.ResultMsg(common.GameTurnErr, common.GameTurnErrMsg)
		return
	}
	start := int(param["Start"].(float64)) // 起始位置
	move := int(param["Move"].(float64))   //移动指令 1左 2右 3上 4下
	chess := urow.GameMsg.Chess
	value := chess[start]                  // 起始值
	animal := valueByChessid(chess, start) // 棋子的值  象>狮>虎>豹>狼>狗>猫>鼠  9 8 7 6 5 4 3 2
	if move > 4 || move < 1 || animal == 0 {
		conn.ResultMsg(common.MessageParamsErr, common.MessageParamsErrMsg)
		return
	}
	isLeft := (value % 100) < 10 // 是否是左边玩家
	if (isLeft && urow.GameMsg.Left == conn.Accid) || (!isLeft && urow.GameMsg.Right == conn.Accid) {
		conn.ResultMsg(common.GameTurnErr, common.GameTurnErrMsg) //不是自己的棋子
		return
	}
	var end int    // 找到结束的位置
	if move == 1 { //  左
		if start%9 == 8 { //当在最后一行的时候不能向左移动
			conn.ResultMsg(common.MessageParamsErr, common.MessageParamsErrMsg)
			return
		}
		end = findEndofMove(chess, start, 1)
	}
	if move == 2 { //  右
		if start%9 == 0 { //当在第一行的时候不能向右移动
			conn.ResultMsg(common.MessageParamsErr, common.MessageParamsErrMsg)
			return
		}
		end = findEndofMove(chess, start, -1)
	}
	if move == 3 { //  上
		if start < 9 {
			conn.ResultMsg(common.MessageParamsErr, common.MessageParamsErrMsg)
			return
		}
		end = findEndofMove(chess, start, -9)
	}
	if move == 4 { //  下
		if start > 53 {
			conn.ResultMsg(common.MessageParamsErr, common.MessageParamsErrMsg)
			return
		}
		end = findEndofMove(chess, start, 9)
	}
	if valueByChessid(chess, start) < valueByChessid(chess, end) {
		conn.ResultMsg(common.MessageParamsErr, common.MessageParamsErrMsg)
		return
	}
	chess[start] = BeginingChess[start]
	chess[end] = value + BeginingChess[end]

	turn++
	if turn > 1 {
		turn = 0
	}

	SendGameMsg(urow.GameMsg.Left, turn, chess)
	SendGameMsg(urow.GameMsg.Right, turn, chess)

}

func findEndofMove(chess [63]int, start int, up int) int {
	end := start + up
	animal := valueByChessid(chess, start) // 棋子的值  象>狮>虎>豹>狼>狗>猫>鼠  9 8 7 6 5 4 3 2
	if checkByChessid(chess, end) == 1 {   //是水
		if animal != 2 { // 老鼠可以进水
			return end
		}
		if animal == 8 || animal == 7 { // 狮虎可以跳水但前提是水里没老鼠
			for i := 1; i < 5; i++ { // 最多三格水 start后面第四格必是陆地
				if checkByChessid(chess, start+up*i) == 1 {
					if valueByChessid(chess, start+up*i) != 0 {
						return 0
					}
				} else {
					end = start + up*i
					break
				}
			}
		} else {
			return 0
		}
	}
	if (end == 27 || end == 35) || (animal != 1) { //兽穴只有对面的才能进
		return 0
	}
	return end
}

func checkByChessid(chess [63]int, chessid int) int { // 0是陆地  1是水 2是陷阱 3是兽穴
	check := chess[chessid] / 100
	return check // 0是陆地  1是水 2是陷阱 3是兽穴
}

func valueByChessid(chess [63]int, chessid int) int { // 棋子的值  象>狮>虎>豹>狼>狗>猫>鼠  9 8 7 6 5 4 3 2

	//在兽穴中战斗力为1
	if chessid == 18 || chessid == 28 || chessid == 36 { //左边兽穴
		if chess[chessid]%100 > 10 {
			return 1
		}
	}
	if chessid == 26 || chessid == 34 || chessid == 44 { //右边边兽穴
		if chess[chessid]%100 < 10 {
			return 1
		}
	}
	value := chess[chessid] % 10
	return value // 棋子的值  象>狮>虎>豹>狼>狗>猫>鼠  9 8 7 6 5 4 3 2
}

func startgame(roomS *room) {
	var (
		left  string
		right string
		turn  int // 0代表left  1代表right
	)
	left = roomS.Left
	right = roomS.Right

	turn = common.IntnBytime(1)

	SendGameMsg(left, turn, BeginingChess)
	SendGameMsg(right, turn, BeginingChess)
}

type gameMsg struct {
	Key     float64
	GameMsg *common.GameRow
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

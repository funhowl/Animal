package common

const (
	USER_NAME = "root"
	PASS_WORD = "206206A@b"
	HOST      = "localhost"
	PORT      = "3306"
	DATABASE  = "chess"
	CHARSET   = "utf8"

	// C2S client to server   S2C server to client
	C2Slogin = 1000 // 登录请求
	S2Clogin = 1001

	C2Sroomessage = 1002 // 房间信息请求
	S2Croomessage = 1003

	C2Sjoinroom = 1004 // 加入房间请求
	S2Cjoinroom = 1005

	C2Sgamechange = 1006 // 游戏棋子移动请求
	// S2Cgamechange = 1007

	C2Sgamemsg = 1008 // 游戏信息请求
	S2Cgamemsg = 1009

	// C2Sgamemsg = 1010
	S2Cgamewin = 1011 // 游戏胜利通知

	ClientConflict    = 9000 //客户端冲突，将原客户端踢下线
	ClientConflictMsg = "账号异地登录"
	LoginErr          = 9001
	LoginErrMsg       = "登录异常"

	RoomEnoughErr    = 9002
	RoomEnoughErrMsg = "房间已满"

	RoomRepeatErr    = 9003
	RoomRepeatErrMsg = "已经在房间里"

	MessageParamsErr    = 9004
	MessageParamsErrMsg = "参数错误"

	GameTurnErr    = 9005
	GameTurnErrMsg = "不是你的回合"

	GameChessErr    = 9006
	GameChessErrMsg = "不是你的棋子"
)

type Userinfo struct { //个人信息模块
	Name string
	Icon string
}

type TaskRow struct {
	Point  int         // 游戏积分 每次游戏结束获得100积分  再写个排行榜每天6点前十名发奖励
	Record string      // 游戏记录
	Uct    interface{} //嵌套扩展模块
}

//https://baike.baidu.com/item/%E6%96%97%E5%85%BD%E6%A3%8B/896091?fr=aladdin
type GameRow struct {
	Roomid int //房间号
	Left   string
	Right  string
	Chess  [63]int // 7 * 9 63个格子的
	Trun   int     // 0代表left  1代表right
}

type UserRow struct { //总模块
	Accid    string
	Gold     int // 金币
	Userinfo *Userinfo
	Tasks    *TaskRow
	GameMsg  *GameRow // game 模块相关的信息 只保存 不更新入数据库
}

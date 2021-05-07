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

	ClientConflict    = 9000 //客户端冲突，将原客户端踢下线
	ClientConflictMsg = "账号异地登录"
	LoginErr          = 9001
	LoginErrMsg       = "登录异常"

	RoomEnoughErr    = 9002
	RoomEnoughErrMsg = "房间已满"
)

type Userinfo struct {
	Name string
	Icon string
}

type TaskRow struct {
	Point int
	Str   string
	Uct   interface{}
}

type GameRow struct {
	Roomid int //房间号
}

type UserRow struct {
	Accid    string
	Gold     int
	Userinfo *Userinfo
	Tasks    *TaskRow
	GameMsg  *GameRow // game 模块相关的信息 只保存 不更新入数据库
}

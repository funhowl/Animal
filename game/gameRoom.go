package game

import (
	"Animal/common"
	"Animal/user"
	"reflect"
)

type room struct {
	roomid int
	// count int
	left  string // room 里面最多两个玩家 就直接用两个字段
	right string
}

var roomList map[int]room
var maxroomid int

type roomessage struct {
	Key    float64
	Data   map[int]room
	Sucess bool
}

type joinroom struct {
	Key  float64
	Data map[int]room
}

func init() {

	roomList = make(map[int]room)
	roomid := 5000

	for i := 0; i < 4; i++ { // 先 生成6个空房间
		var roomNew room
		roomNew.roomid = roomid + i + 1
		roomList[roomNew.roomid] = roomNew
		if roomNew.roomid > maxroomid {
			maxroomid = roomNew.roomid
		}
	}

	user.RegisterEvent(common.C2Sroomessage, new(roomessage).roomessage) // 房间信息请求
	user.RegisterEvent(common.C2Sjoinroom, new(joinroom).joinroom)       //  加入房间请求
}

func (*roomessage) roomessage(conn *user.Connection, param map[string]interface{}) {

	var (
		err error
		ans joinroom
	)

	ans.Key = common.S2Croomessage
	ans.Data = roomList

	if err = conn.Send(ans); err != nil {
		panic(err.Error())
	}
}

func (*joinroom) joinroom(conn *user.Connection, param map[string]interface{}) {
	var (
		err error
		ans roomessage
	)

	ans.Key = common.S2Cjoinroom

	roomid := param["roomid"].(int)
	roomS := roomList[roomid]
	if reflect.ValueOf(roomS.left).IsZero() { // 没人
		ans.Sucess = true
		roomS.left = conn.Accid
		if roomid == maxroomid && roomid < 5008 { // 由于本人前端水平 就只开8个房间
			for i := 0; i < 2; i++ { // 生成2个空房间
				var roomNew room
				roomNew.roomid = roomid + i + 1
				roomList[roomNew.roomid] = roomNew
				if roomNew.roomid > maxroomid {
					maxroomid = roomNew.roomid
				}
			}
		}
	} else if reflect.ValueOf(roomS.right).IsZero() { // 有一个
		ans.Sucess = true
		roomS.right = conn.Accid
		// start game
	} else { // 人满
		ans.Sucess = false
	}

	roomList[roomid] = roomS
	ans.Data = roomList

	if err = conn.Send(ans); err != nil {
		panic(err.Error())
	}
}

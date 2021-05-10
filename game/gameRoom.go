package game

import (
	"Animal/common"
	"Animal/user"
	"fmt"
	"reflect"
)

type room struct {
	roomid int
	// count int
	left  string // room 里面最多两个玩家 就直接用两个字段
	right string
}

var roomList map[int]*room
var maxroomid int

type roomessage struct {
	Key  float64
	Data interface{}
}

type joinroom struct {
	Key  float64
	Data interface{}
}

func init() {

	roomList = make(map[int]*room)
	roomid := 5000

	for i := 0; i < 4; i++ { // 先 生成4个空房间
		var roomNew room
		roomNew.roomid = roomid + i + 1
		roomNew.left = "wo shi "
		roomList[roomNew.roomid] = &roomNew
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
	fmt.Print("ans:", ans)

	if err = conn.Send(ans); err != nil {
		panic(err.Error())
	}
}

func (*joinroom) joinroom(conn *user.Connection, param map[string]interface{}) {
	var (
		err        error
		ans        roomessage
		urow       *common.UserRow
		istartgame bool
	)

	if urow, err = user.GetRowById(conn.Accid); err != nil {
		panic(err.Error())
	}
	if urow.GameMsg.Roomid > 5000 {
		conn.ResultMsg(common.RoomRepeatErr, common.RoomRepeatErrMsg)
		return
	}

	ans.Key = common.S2Cjoinroom

	roomid := param["roomid"].(int)
	roomS := roomList[roomid]
	if reflect.ValueOf(roomS.left).IsZero() { // 没人
		roomS.left = conn.Accid

		if roomid == maxroomid && roomid < 5008 { // 由于本人前端水平 就只开8个房间
			for i := 0; i < 2; i++ { // 生成2个空房间
				var roomNew *room
				roomNew.roomid = roomid + i + 1
				roomList[roomNew.roomid] = roomNew
				if roomNew.roomid > maxroomid {
					maxroomid = roomNew.roomid
				}
			}
		}
	} else if reflect.ValueOf(roomS.right).IsZero() { // 有一个

		roomS.right = conn.Accid
		// start game
		istartgame = true
		// left := roomS.left

	} else { // 人满
		conn.ResultMsg(common.RoomEnoughErr, common.RoomEnoughErrMsg)
		return
	}

	roomList[roomid] = roomS
	ans.Data = roomList

	// uplist := *urow
	// uplist.GameMsg.Roomid = roomid
	// *urow = uplist

	urow.GameMsg.Roomid = roomid
	urow.GameMsg.Left = roomS.left
	urow.GameMsg.Right = roomS.right
	user.UserRowData[conn.Accid] = urow

	if err = conn.Send(ans); err != nil {
		panic(err.Error())
	}
	if istartgame {
		startgame(roomS)
	}
}

//离开房间暂时就不写了

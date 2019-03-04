package golang

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
	"masterlab_socket/lib"
)

type Root struct {
	Code int             `json:"code"`
	Msg  string           `json:"msg"`
	Data interface{}     `json:"data"`
	Pages int             `json:"pages"`
}

type ListType struct {
	Mine   map[string]string      `json:"mine"`
	Friend []FriendType            `json:"friend"`
	Group  []map[string]string     `json:"group"`
}

type FriendType struct {
	Groupname string               `json:"groupname"`
	Online    int                        `json:"online"`
	Id        int                        `json:"id"`
	List      []map[string]string        `json:"list"`
}

type MemberType struct {
	Owner   map[string]string    `json:"owner"`
	Members int                        `json:"members"`
	List    []map[string]string        `json:"list"`
}

type SysMsgType struct {
	Id         int64                    `json:"id"`
	Content    string                   `json:"content"`
	Username    string                   `json:"username"`
	Uid        int64                     `json:"uid"`
	From       int64                     `json:"from"`
	From_group int64                     `json:"from_group"`
	Type       int64                      `json:"type"`
	Href       string                     `json:"href"`
	Read       int64                      `json:"read"`
	Remark     string                     `json:"remark"`
	Time       string                      `json:"time"`
	Status     int64                       `json:"status"`
	User       map[string]string           `json:"user"`
}

func GetUserRow(db *sql.DB, sid string) map[string]string {

	sql_str := `select id,nick,status ,sign, avatar,token  from user where sid=?`
	var id, nick, status, sign, avatar, token string
	record := make(map[string]string)
	err := db.QueryRow(sql_str, sid).Scan(&id, &nick, &status, &sign, &avatar, &token)
	if err != nil {
		fmt.Println("getUserRow err:", err.Error())
		return record
	}
	record["id"] = id
	record["username"] = nick
	record["sign"] = sign
	record["status"] = status
	record["sid"] = sid
	record["avatar"] = avatar
	record["token"] = token

	return record
}

func getMyContacts(db *sql.DB, uid int) []map[string]string {

	sql_str := "SELECT  u.id,u.nick as nick,u.avatar,u.sign,c.group_id,u.sid  FROM `contacts` c LEFT JOIN `user` u on u.id =c.uid WHERE  c.master_uid=?"

	contact_records := make([]map[string]string, 0)
	rows, err := db.Query(sql_str, uid)
	if err != nil {

		return contact_records
	}
	for rows.Next() {
		//将行数据保存到record字典
		var id, nick, avatar, sign, group_id, sid string
		record := make(map[string]string)
		rows.Scan(&id, &nick, &avatar, &sign, &group_id, &sid)

		record["id"] = id
		record["username"] = nick
		record["avatar"] = avatar
		record["sign"] = sign
		record["group_id"] = group_id
		record["sid"] = sid
		contact_records = append(contact_records, record)

	}
	return contact_records

}

func getMyGroup(db *sql.DB, uid int) []map[string]string {

	sql_str := "SELECT  id,title as groupname  FROM `contact_group` WHERE uid=? "
	my_group_records := make([]map[string]string, 0)
	rows, err := db.Query(sql_str, uid)
	if err != nil {
		return my_group_records
	}
	for rows.Next() {
		//将行数据保存到record字典
		var gid, groupname string
		record := make(map[string]string)
		err = rows.Scan(&gid, &groupname)
		if err != nil {
			fmt.Println("服务器错误@" + err.Error())
			return my_group_records
		}
		record["id"] = gid
		record["groupname"] = groupname
		fmt.Println(record)
		my_group_records = append(my_group_records, record)
	}
	return my_group_records
}

func getFriends(db *sql.DB, uid int) []FriendType {

	friends := make([]FriendType, 0)

	// 获取所属的联系人列表（未分组）
	contact_records := getMyContacts(db, uid)

	// 获取分组
	my_group_records := getMyGroup(db, uid)
	var friend FriendType
	for _, group := range my_group_records {
		friend = FriendType{}
		friend.Groupname = group[`groupname`]
		friend.Id, _ = strconv.Atoi(group[`id`])
		friend.Online = 1
		tmp_list := make([]map[string]string, 0)

		for _, c := range contact_records {
			group_id, _ := strconv.Atoi(c[`group_id`])
			if group_id == friend.Id {
				tmp_list = append(tmp_list, c)
				//contact_records = append(contact_records[:_k], contact_records[_k+1:]...)
			}
		}
		friend.List = tmp_list
		friends = append(friends, friend)
	}

	return friends
}

func getMyGroups(db *sql.DB, uid int) []map[string]string {

	sql_str := "SELECT id,channel_id,pic as avatar,title  FROM `global_group` WHERE  id in( SELECT `group_id` FROM `user_join_group` WHERE `uid`=? )"
	join_group_records := make([]map[string]string, 0)
	rows, err := db.Query(sql_str, uid)
	if err != nil {
		fmt.Println(504, "服务器错误@" + err.Error())
		return join_group_records
	}
	for rows.Next() {
		//将行数据保存到record字典
		var cid, channel_id, avatar, title string
		record := make(map[string]string)
		err = rows.Scan(&cid, &channel_id, &avatar, &title)
		if err != nil {
			fmt.Println(505, "服务器错误@" + err.Error())
			return join_group_records
		}
		record["id"] = cid
		record["channel_id"] = channel_id
		record["avatar"] = avatar
		record["groupname"] = title
		record["title"] = title
		//fmt.Println(record)
		join_group_records = append(join_group_records, record)
	}
	fmt.Println(join_group_records)
	return join_group_records

}

func getMembers(db  *sql.DB, member_id int) []map[string]string {

	sql_str := "SELECT U.id,U.nick, U.sign, U.avatar,U.sid  FROM `user_join_group` G LEFT JOIN user U on G.uid=U.id WHERE  group_id=?"
	members := make([]map[string]string, 0)
	rows, err := db.Query(sql_str, member_id)
	if err != nil {
		fmt.Println(504, "服务器错误@" + err.Error())
		return members
	}
	for rows.Next() {
		//将行数据保存到record字典
		var id, nick, sign, avatar, sid string
		record := make(map[string]string)
		err = rows.Scan(&id, &nick, &sign, &avatar, &sid)
		if err != nil {
			fmt.Println(505, "服务器错误@" + err.Error())
			return members
		}
		record["id"] = id
		record["sign"] = sign
		record["avatar"] = avatar
		record["username"] = nick
		//fmt.Println(record)
		members = append(members, record)
	}
	fmt.Println(members)
	return members

}
func getRecommendUser(db  *sql.DB, uid int) []map[string]string {

	sql_str := "SELECT id,nick,avatar,sign FROM `user` WHERE id  not in( select uid from contacts where master_uid=? ) AND id!=?"
	members := make([]map[string]string, 0)
	rows, err := db.Query(sql_str, uid, uid,)
	if err != nil {
		fmt.Println(504, "服务器错误@" + err.Error())
		return members
	}
	for rows.Next() {
		//将行数据保存到record字典
		var id, nick, sign, avatar string
		record := make(map[string]string)
		err = rows.Scan(&id, &nick, &avatar, &sign)
		if err != nil {
			fmt.Println(505, "服务器错误@" + err.Error())
			return members
		}
		record["id"] = id
		record["sign"] = sign
		record["avatar"] = avatar
		record["username"] = nick
		//fmt.Println(record)
		members = append(members, record)
	}
	return members
}

func searchGroup(db  *sql.DB, uid int, name string ) []map[string]string {

	datas := make([]map[string]string, 0)
	var rows *sql.Rows
	var err error
	sql_str:=""
	if name==""{
		sql_str = "SELECT id, title, channel_id, pic, remark  FROM `global_group` WHERE " +
			"id not in ( SELECT group_id from  user_join_group where uid=?)"
		rows, err  = db.Query(sql_str, uid,)

	}else{
		sql_str = "SELECT id, title, channel_id, pic, remark  FROM `global_group` WHERE" +
			"( id not in    ( SELECT group_id from  user_join_group where uid=?)  )" +
			" AND (locate (? , title) > 0 )"

		rows, err  = db.Query(sql_str, uid, name,)
	}
	//fmt.Println( sql_str )
	if err != nil {
		fmt.Println(504, "服务器错误@" + err.Error())
		return datas
	}
	for rows.Next() {
		//将行数据保存到record字典
		var id, title, channel_id, pic, remark string
		record := make(map[string]string)
		err = rows.Scan(&id, &title, &channel_id, &pic, &remark )
		if err != nil {
			fmt.Println(505, "服务器错误@" + err.Error())
			return datas
		}
		record["id"] = id
		record["title"] = title
		record["channel_id"] = channel_id
		record["pic"] = pic
		record["remark"] = remark
		datas = append(datas, record)
	}
	return datas
}

func getSysMsgs(db  *sql.DB, uid int) []SysMsgType {

	sql_str := "  SELECT  r.id,R.type, R.from_uid,R.add_group, R.remark,R.readed,R.time,R.status,U.avatar,U.nick ,U.sign  " +
		"            FROM `req_friend` R   " +
		"            LEFT JOIN user U on R.from_uid=u.id  " +
		"            WHERE  R.req_uid=? Order by `up_time` DESC"
	msgs := make([]SysMsgType, 0)
	rows, err := db.Query(sql_str, uid)
	if err != nil {
		fmt.Println(504, "服务器错误@" + err.Error())
		return msgs
	}
	for rows.Next() {
		//将行数据保存到record字典
		var id, from_uid, from_group, _type, readed, add_time, status int64
		var remark, avatar, nick, sign string

		err = rows.Scan(&id, &_type, &from_uid, &from_group, &remark, &readed, &add_time, &status, &avatar, &nick, &sign)
		if err != nil {
			fmt.Println(505, "服务器错误@" + err.Error())
			return msgs
		}
		content := ""
		if ( _type == 1 ) {
			content = nick+" 申请添加你为好友"
		}
		sysMsgType := SysMsgType{}
		sysMsgType.Id = id
		sysMsgType.Type = _type
		sysMsgType.Content = content
		sysMsgType.Uid = from_uid
		sysMsgType.From_group = from_group
		sysMsgType.Remark = remark
		sysMsgType.Href = ""
		sysMsgType.Read = readed
		sysMsgType.Status = status
		if( add_time>0 ){
			sysMsgType.Time = time.Unix(add_time, 0).Format("2006-01-02 03:04:05")
		}else{
			sysMsgType.Time = ""
		}

		user := make(map[string]string)
		user["id"] = strconv.FormatInt(from_uid, 10)
		user["avatar"] = avatar
		user["username"] = nick
		user["sign"] = sign
		sysMsgType.User = user

		//fmt.Println(record)
		msgs = append(msgs, sysMsgType)
	}
	return msgs
}



func  JoinChannel(db *sql.DB, uid int, sid string ) {

	mysql := new(lib.Mysql)
	_, err := mysql.Connect()
	if err != nil {
		fmt.Println("数据库连接失败:" + err.Error())
		return
	}
	sql_str := "SELECT   g.channel_id,g.title   FROM `user_join_group`  J LEFT JOIN global_group G  ON J.group_id=G.id   WHERE  J.uid=?"
	rows, err := mysql.Db.Query( sql_str,uid )
	if err != nil {
		fmt.Println(504, "服务器错误@" + err.Error())
	}
	for rows.Next() {
		//将行数据保存到record字典
		var title, area_id string
		err = rows.Scan( &area_id , &title)
		if err != nil {
			fmt.Println(505, "服务器错误@" + err.Error())
			return
		}
		sdk:=new(Sdk).InitCmd("JoinArea",sid,0,[]byte("") )
		sdk.AreaAddSid( sid ,area_id )
		//fmt.Println( "JoinArea...",  sid ,area_id )
	}

}


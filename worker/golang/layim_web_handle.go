package golang

import (
	_"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/go-sql-driver/mysql"
	"masterlab_socket/lib"
	"masterlab_socket/area"

)


// 初始化 http请求
func InitHandler(){
	http.HandleFunc("/upload_image", UploadImageHandler)
	http.HandleFunc("/upload_file", UploadFileHandler)
	http.HandleFunc("/reg", RegHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/get_list", GetListHandler)
	http.HandleFunc("/get_member", GetMemberHandler)
	http.HandleFunc("/get_recommend_user",GetRecommendUserHandler)
	http.HandleFunc("/req_add_friend", ReqAddFriendHandler)
	http.HandleFunc("/sysmsg", SystemMsgHandler)
	http.HandleFunc("/agree", AgreeHandler)
	http.HandleFunc("/reject", RejectHandler)
	http.HandleFunc("/search_group", searchGroupHandler)
	http.HandleFunc("/add_group", ReqAddGroupHandler)


}

// 初始化群组
func InitGlobalGroup(){
	mysql := new(lib.Mysql)
	_, err := mysql.Connect()
	if err != nil {
		fmt.Println("数据库连接失败:" + err.Error())
		return
	}
	sql_str := "SELECT   `id` ,`title`,`channel_id`     FROM `global_group` WHERE 1"
	rows, err := mysql.Db.Query( sql_str )
	if err != nil {
		fmt.Println(504, "服务器错误@" + err.Error())
		return
	}
	for rows.Next() {
		//将行数据保存到record字典
		var id int64
		var title, channel_id string
		err = rows.Scan(&id, &title, &channel_id )
		if err != nil {
			fmt.Println(505, "服务器错误@" + err.Error())
			return
		}
		area.Create( channel_id,title )
	}
}


func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	format_str := `{ "code":%d ,"msg": "%s","data": {"src":"%s","name":"%s" }} `

	if r.Method == "GET" {
		resp := fmt.Sprintf(format_str, 401, "GET no support!", "", "")
		w.Write([]byte(resp))
		return

	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			//fmt.Println(err)
			resp := fmt.Sprintf(format_str, 400, err.Error(), "", "")
			w.Write([]byte(resp))
			return
		}
		defer file.Close()

		//fmt.Fprintf(w, "%v", handler.Header)
		wd, _ := os.Getwd()
		upload_dir := fmt.Sprintf("%s/web/wwwroot/data/images/", wd)
		f, err := os.OpenFile(upload_dir+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		code := 0
		err_str := ""
		src := "http://" + r.Host + "/data/images/" + handler.Filename
		if err != nil {
			fmt.Println(err)
			code = 500
			err_str = err.Error()
		} else {
			defer f.Close()
			io.Copy(f, file)
		}
		resp := fmt.Sprintf(format_str, code, err_str, src, handler.Filename)
		w.Write([]byte(resp))
	}

}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	format_str := `{ "code":%d ,"msg": "%s","data": {"src":"%s","name":"%s" }} `

	if r.Method == "GET" {
		resp := fmt.Sprintf(format_str, 401, "GET no support!", "", "")
		w.Write([]byte(resp))
		return

	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			//fmt.Println(err)
			resp := fmt.Sprintf(format_str, 400, err.Error(), "", "")
			w.Write([]byte(resp))
			return
		}
		defer file.Close()

		//fmt.Fprintf(w, "%v", handler.Header)
		wd, _ := os.Getwd()
		upload_dir := fmt.Sprintf("%s/web/wwwroot/data/files/", wd)
		f, err := os.OpenFile(upload_dir+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		code := 0
		err_str := ""
		src := "http://" + r.Host + "/data/files/" + handler.Filename
		if err != nil {
			fmt.Println(err)
			code = 500
			err_str = err.Error()
		} else {
			defer f.Close()
			io.Copy(f, file)
		}
		resp := fmt.Sprintf(format_str, code, err_str, src, handler.Filename)
		w.Write([]byte(resp))
	}

}

func RegHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	format_str := `{ "code":%d ,"msg": "%s","data": { }} `

	if r.Method == "GET" {
		resp := fmt.Sprintf(format_str, 401, "GET no support!")
		w.Write([]byte(resp))
		return

	} else {

		r.ParseForm()

		user := r.PostForm.Get(`user`)
		pwd := r.PostForm.Get(`pwd`)
		age := r.PostForm.Get(`age`)
		nick := r.PostForm.Get(`nick`)
		sign := r.PostForm.Get(`sign`)
		avatar := r.PostForm.Get(`avatar`)
		reg_time := time.Now().Unix()
		sid := area.CreateSid()

		db := new(lib.Mysql)
		db.Connect()

		row := db.GetRow(`select user from user where user=? `, user)

		if _, ok := row[`user`]; ok {
			resp := fmt.Sprintf(format_str, 0, "用户名已经存在!")
			w.Write([]byte(resp))
			return
		}

		insert_id, err := db.Insert(`INSERT user (user,pwd,nick,sign,age,sid,avatar,reg_time)
						    values (?,?,?,?,?,?,?,?)`,
			user, pwd, nick, sign, age,sid,avatar, reg_time)
		if err != nil {
			resp := fmt.Sprintf(format_str, 500, "db.Insert err:", err.Error())
			w.Write([]byte(resp))
			return
		}

		// 添加一个默认的分组
		_, err = db.Insert("INSERT INTO `contact_group` (  `uid`, `title` ) VALUES (  ?, ? ) " +
			"VALUES (  ?, ? )", insert_id,"默认")
		if err != nil {
			fmt.Println( "INSERT contact_group err:", err.Error() )
		}

		fmt.Println("insertid:", insert_id)
		resp := fmt.Sprintf(format_str, 1, "注册成功")
		w.Write([]byte(resp))
	}

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	format_str := `{ "code":%d ,"msg": "%s", "data": {}} `
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if r.Method == "GET" {
		resp := fmt.Sprintf(format_str, 401, "GET no support!")
		w.Write([]byte(resp))
		return

	} else {
		r.ParseForm()
		user := r.PostForm.Get(`user`)
		pwd := r.PostForm.Get(`pwd`)
		fmt.Println(user, pwd, mysql.MySQLDriver{})
		db := new(lib.Mysql)
		db.Connect()

		resp := ""
		sql_str := `select id,user,sign,sid ,avatar from user  where user=? and pwd=?`
		var id, sign, sid, avatar string
		record := make(map[string]string)
		scan_err := db.Db.QueryRow(sql_str, user, pwd).Scan(&id, &user, &sign, &sid, &avatar)
		if scan_err != nil {
			resp = fmt.Sprintf(format_str, 500, "用户名密码错误"+scan_err.Error())
			w.Write([]byte(resp))
			return
		}
		record["id"] = id
		record["user"] = user
		record["sign"] = sign
		record["sid"] = sid
		record["avatar"] = avatar
		token := area.CreateSid()
		affect_num,_:=db.Update( `Update user set token=? Where id=?`,token,id)
		if affect_num>0 {
			record["token"] = token
		}

		fmt.Println(record)
		json_encode, _ := json.Marshal(record)

		uid, _ := strconv.Atoi( record["id"] )
		friends := getMyContacts( db.Db, uid )
		friends_encode, _ := json.Marshal( friends )
		groups := getMyGroups( db.Db, uid )
		groups_encode, _ := json.Marshal( groups )

		if id != "" {
			resp = fmt.Sprintf(`{ "code":%d ,"msg": "%s","data":%s,"contacts":%s,"groups":%s} `,
				1, "验证成功", string(json_encode) ,string(friends_encode), string(groups_encode) )
		} else {
			resp = fmt.Sprintf(format_str, 404, "用户名密码错误")
		}
		w.Write([]byte(resp))
	}

}


func GetListHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	root := new(Root)
	_list := new(ListType)
	root.Data = &_list

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		r.ParseForm()
		id_str := r.Form.Get(`id`)
		id, _ := strconv.Atoi(id_str)
		sid := r.Form.Get(`sid`)
		fmt.Println(id, sid, mysql.MySQLDriver{})
		db := new(lib.Mysql)
		_, err := db.Connect()
		if err != nil {
			root.Code = 500
			root.Msg = "数据库连接失败:" + err.Error()
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}

		// 获取当前用户信息
		my_record := GetUserRow(db.Db, sid)
		_, ok := my_record[`id`]
		if !ok {
			root.Code = 400
			root.Msg = "用户验证失败"
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}
		uid, _ := strconv.Atoi(my_record["id"])

		_list.Mine = my_record
		_list.Friend = getFriends(db.Db, uid)
		_list.Group = getMyGroups( db.Db,uid)

		root.Code = 0
		root.Msg = ""
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
	}
}


func GetMemberHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	root := new(Root)
	member := new(MemberType)
	root.Data = &member

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		r.ParseForm()
		id_str := r.Form.Get(`id`)
		id, _ := strconv.Atoi(id_str)
		sid := r.Form.Get(`sid`)
		fmt.Println(id, sid, mysql.MySQLDriver{})
		db := new(lib.Mysql)
		_, err := db.Connect()
		if err != nil {
			root.Code = 500
			root.Msg = "数据库连接失败:" + err.Error()
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}

		// 获取当前用户信息
		my_record := GetUserRow(db.Db, sid)
		_, ok := my_record[`id`]
		if !ok {
			root.Code = 400
			root.Msg = "用户验证失败"
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}

		member.Owner = my_record
		member.List = getMembers(db.Db, id)
		member.Members = len( member.List  )

		root.Code = 0
		root.Msg = ""
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
	}

}

func GetRecommendUserHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	root := new(Root)

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		r.ParseForm()
		sid := r.Form.Get(`sid`)
		fmt.Println(  sid, mysql.MySQLDriver{})
		db := new(lib.Mysql)
		_, err := db.Connect()
		if err != nil {
			root.Code = 500
			root.Msg = "数据库连接失败:" + err.Error()
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}

		// 获取当前用户信息
		my_record := GetUserRow(db.Db, sid)
		_, ok := my_record[`id`]
		if !ok {
			root.Code = 400
			root.Msg = "用户验证失败"
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}
		id, _ := strconv.Atoi(my_record["id"])
		root.Data = getRecommendUser(db.Db, id)

		root.Code = 0
		root.Msg = ""
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
	}

}


func searchGroupHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	root := new(Root)
	if r.Method != "GET" {
		root.Code = 400
		root.Msg = "请使用GET请求"
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
		return

	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	r.ParseForm()
	name := ""
	sid := r.Form.Get(`sid`)
	name = r.Form.Get(`name`)
	fmt.Println(  sid, mysql.MySQLDriver{})
	db := new(lib.Mysql)
	_, err := db.Connect()
	if err != nil {
		root.Code = 500
		root.Msg = "数据库连接失败:" + err.Error()
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
		return
	}

	// 获取当前用户信息
	my_record := GetUserRow(db.Db, sid)
	_, ok := my_record[`id`]
	if !ok {
		root.Code = 400
		root.Msg = "用户验证失败"
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
		return
	}
	id, _ := strconv.Atoi(my_record["id"])
	root.Data = searchGroup(db.Db, id, name )

	root.Code = 0
	root.Msg = ""
	root.Pages = 1   // @todo 未做分页
	json_encode ,_:=json.Marshal( root )
	w.Write( json_encode )


}

func ReqAddGroupHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	root := new(Root)
	if r.Method != "GET" {
		root.Code = 400
		root.Msg = "请使用GET请求"
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
		return

	}
	r.ParseForm()
	sid := r.Form.Get(`sid`)
	group_id := r.Form.Get(`group_id`)

	db := new(lib.Mysql)
	db.Connect()

	// 获取当前用户信息
	my_record := GetUserRow(db.Db, sid)
	uid := my_record["id"]

	sql_str :="SELECT title, pic, channel_id, remark  FROM `global_group` WHERE   `id`=?"
	var  title, pic, area_id, remark string
	err := db.Db.QueryRow( sql_str, group_id ).Scan( &title, &pic, &area_id, &remark )
	if err != nil {
		root.Code = 500
		root.Msg = "群组不存在:"+err.Error()
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
		return
	}


	sql_str ="SELECT id FROM `user_join_group` WHERE  `uid` =? AND `group_id`=?"
	rows, err := db.Db.Query(sql_str, uid,group_id)
	if err != nil {
		root.Code = 500
		root.Msg = "服务器错误:"+err.Error()
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
		return
	}
	if rows.Next() {
		root.Code = 505
		root.Msg = "您已经加入过该群组"
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
		return
	}

	_, err = db.Insert("INSERT INTO `user_join_group` (  `uid`, `group_id` ) " +
		"   VALUES (  ?, ? );", uid, group_id)

	if err != nil {
		root.Code = 500
		root.Msg = "服务器错误:"+err.Error()
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
		return
	}
	record := make(map[string]string)
	record["id"] = group_id
	record["pic"] = pic
	record["title"] = title
	record["channel_id"] = area_id
	record["remark"] = remark

	// 订阅群组消息
	sdk:=new(Sdk).InitCmd("JoinArea",sid,0,[]byte("") )

	sdk.AreaAddSid( sid ,area_id )

	root.Code = 0
	root.Msg = "添加成功"
	root.Data = record
	json_encode ,_:=json.Marshal( root )
	w.Write( json_encode )

}


func SystemMsgHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	root := new(Root)

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		r.ParseForm()

		sid := r.Form.Get(`sid`)
		fmt.Println(  sid, mysql.MySQLDriver{})
		db := new(lib.Mysql)
		_, err := db.Connect()
		if err != nil {
			root.Code = 500
			root.Msg = "数据库连接失败:" + err.Error()
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}

		// 获取当前用户信息
		my_record := GetUserRow(db.Db, sid)
		_, ok := my_record[`id`]
		if !ok {
			root.Code = 400
			root.Msg = "用户验证失败"
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}
		id, _ := strconv.Atoi(my_record["id"])
		root.Data = getSysMsgs(db.Db, id)

		root.Code = 0
		root.Msg = ""
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
	}

}


func ReqAddFriendHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	format_str := `{ "code":%d ,"msg": "%s","data": { }} `

	if r.Method != "GET" {
		resp := fmt.Sprintf(format_str, 401, "POST no support!")
		w.Write([]byte(resp))
		return

	} else {
		r.ParseForm()
		sid := r.PostForm.Get(`sid`)
		req_uid := r.PostForm.Get(`uid`)
		remark := r.PostForm.Get(`remark`)
		add_group := r.PostForm.Get(`add_group`)
		add_time := time.Now().Unix()

		db := new(lib.Mysql)
		db.Connect()

		// 获取当前用户信息
		my_record := GetUserRow(db.Db, sid)

		from_uid := my_record["id"]
		status := "1"
		readed := "0"

		insert_id, err := db.Insert("INSERT INTO `req_friend` (  `from_uid`, `add_group`, `readed`, `req_uid`, `status`, `remark`, `time`,`up_time`) " +
			"   VALUES (  ?, ?, ?, ?, ?, ?, ?,?);", from_uid,add_group,readed,req_uid,status,remark,add_time,add_time)

		if err != nil {
			resp := fmt.Sprintf(format_str, 500, "db.Insert err:", err.Error())
			w.Write([]byte(resp))
			return
		}
		fmt.Println("insertid:", insert_id)
		resp := fmt.Sprintf(format_str, 0, "申请成功")
		w.Write([]byte(resp))
	}
}

func AgreeHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	format_str := `{ "code":%d ,"msg": "%s","data": { }} `
	root := new(Root)
	if r.Method == "GET" {
		resp := fmt.Sprintf(format_str, 401, "GET no support!")
		w.Write([]byte(resp))
		return

	} else {

		r.ParseForm()


		id :=r.PostForm.Get("id")
		add_group := r.PostForm.Get(`group`)
		sid :=r.PostForm.Get("sid")

		db := new(lib.Mysql)
		db.Connect()

		// 获取当前用户信息
		my_record := GetUserRow(db.Db, sid)
		_, ok := my_record[`id`]
		if !ok {
			root.Code = 400
			root.Msg = "用户验证失败"
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}
		var  from_uid,from_group,req_uid int64
		rows, err := db.Db.Query( "SELECT  `from_uid`,`add_group`,`req_uid` FROM `req_friend` WHERE `id`=? ", id)
		if err != nil {
			root.Code = 401
			root.Msg = "服务器错误@" + err.Error()
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}
		for rows.Next() {
			err = rows.Scan(&from_uid, &from_group, &req_uid)
			if err != nil {
				root.Code = 500
				root.Msg = "服务器错误@" + err.Error()
				json_encode ,_:=json.Marshal( root )
				w.Write( json_encode )
				return
			}
			break
		}
		_uid, _ := strconv.ParseInt(my_record["id"],10,0)
		if( _uid!= req_uid){
			root.Code = 401
			root.Msg = "非当前用户消息"
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}

		//修改状态
		_,err = db.Update( "UPDATE `req_friend` SET status=2 ,up_time= ? WHERE id=?",time.Now().Unix(),id )
		if err != nil {
			root.Code = 501
			root.Msg = "服务器错,err:"+err.Error()
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}

		// 请求者处理
		_, err = db.Insert("INSERT INTO `contacts` (  `master_uid`, `group_id`, `uid` ) " +
			                      "VALUES ( ?, ?, ? )", from_uid,from_group,req_uid)
		if err != nil {
			root.Code = 502
			root.Msg = "服务器错,err:"+err.Error()
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}

		// 被请求者处理
		_, err = db.Insert("INSERT INTO `contacts` (  `master_uid`, `group_id`, `uid` ) " +
			"VALUES ( ?, ?, ? )", req_uid,add_group,from_uid)
		if err != nil {
			root.Code = 503
			root.Msg = "服务器错,err:"+err.Error()
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}

		root.Code = 0
		root.Msg = "处理成功"
		root.Data = nil
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
	}

}

func RejectHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	format_str := `{ "code":%d ,"msg": "%s","data": { }} `
	root := new(Root)
	if r.Method == "GET" {
		resp := fmt.Sprintf(format_str, 401, "GET no support!")
		w.Write([]byte(resp))
		return

	} else {

		r.ParseForm()

		id :=r.PostForm.Get("id")
		sid :=r.PostForm.Get("sid")

		db := new(lib.Mysql)
		db.Connect()

		// 获取当前用户信息
		my_record := GetUserRow(db.Db, sid)
		_, ok := my_record[`id`]
		if !ok {
			root.Code = 400
			root.Msg = "用户验证失败"
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}
		var  from_uid,from_group,req_uid int64
		rows, err := db.Db.Query( "SELECT  `from_uid`,`add_group`,`req_uid` FROM `req_friend` WHERE `id`=?", id)
		if err != nil {
			root.Code = 401
			root.Msg = "服务器错误@" + err.Error()
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}
		for rows.Next() {
			err = rows.Scan( &from_uid, &from_group, &req_uid )
			if err != nil {
				root.Code = 500
				root.Msg = "服务器错误@" + err.Error()
				json_encode ,_:=json.Marshal( root )
				w.Write( json_encode )
				return
			}
			break
		}
		if err != nil {
			root.Code = 500
			root.Msg = "服务器错误@" + err.Error()
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}
		_uid, _ := strconv.ParseInt(my_record["id"],10,0)
		if( _uid!= req_uid){
			root.Code = 401
			root.Msg = "非当前用户消息"
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}

		//修改状态
		_,err = db.Update( "UPDATE `req_friend` SET status=2 ,up_time=? WHERE id=?",time.Now().Unix(),id )
		if err != nil {
			root.Code = 501
			root.Msg = "服务器错,err:"+err.Error()
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}

		root.Code = 0
		root.Msg = "处理成功"
		root.Data = nil
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
	}

}




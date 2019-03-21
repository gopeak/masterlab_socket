package lib

import (
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	Db        *sql.DB
	Sql       string
	Config    MysqlConfigType
	Connected bool
}

type MysqlConfigType struct {
	Database     string `toml:"database"`
	User         string `toml:"user"`
	Password     string `toml:"password"`
	Host         string `toml:"host"`
	Port         string `toml:"port"`
	Charset      string `toml:"charset"`
	Timeout      string `toml:"timeout"`
	MaxOpenConns int    `toml:"max_open_conns"`
	MaxIdleConns int    `toml:"max_idle_conns"`
}

func (this *Mysql) Connect() (bool, error) {
	var err error
	var config MysqlConfigType
	if (!this.Connected) {
		_, err = toml.DecodeFile("worker.toml", &config)
		if err != nil {
			fmt.Println("toml.DecodeFile error:", err.Error())
			this.Connected = false
			return false, err
		}
		fmt.Println("config:", config)
		connect_str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=%ss&collation=%s", config.User,
			config.Password, config.Host, config.Port, config.Database, config.Timeout, config.Charset)
		fmt.Println("connect_str:", connect_str)
		this.Db, err = sql.Open("mysql", connect_str)
		if err != nil {
			fmt.Println("sql.Open err:", err.Error())
			this.Connected = false
			return false, err
		}
		this.Db.SetMaxOpenConns(config.MaxOpenConns)
		this.Db.SetMaxIdleConns(config.MaxIdleConns)
		this.Db.Ping()
		this.Connected = true
	}
	return true, nil
}

func (this *Mysql) ShortConnect() (bool, error) {
	var err error
	var config MysqlConfigType
	//if( !this.Connected ){
	_, err = toml.DecodeFile("worker.toml", &config)
	if err != nil {
		fmt.Println("toml.DecodeFile error:", err.Error())
		this.Connected = false
		return false, err
	}
	fmt.Println("config:", config)
	connect_str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=%ss&collation=%s", config.User,
		config.Password, config.Host, config.Port, config.Database, config.Timeout, config.Charset)
	fmt.Println("connect_str:", connect_str)
	this.Db, err = sql.Open("mysql", connect_str)
	if err != nil {
		fmt.Println("sql.Open err:", err.Error())
		this.Connected = false
		return false, err
	}
	this.Db.SetMaxOpenConns(0)
	this.Db.SetMaxIdleConns(0)
	this.Connected = true
	//}
	return true, nil
}

//插入 封装
func (this *Mysql) Insert(sql string, args ...interface{}) (int64, error) {

	stmt, err := this.Db.Prepare(sql)
	if err != nil {
		fmt.Println("Insert err:" + err.Error())
		return 0, err
	}
	res, err := stmt.Exec(args...)
	if err != nil {
		fmt.Println("Insert err:" + err.Error())
		return 0, err
	}
	return res.LastInsertId()

}

//查询多行封装
func (this *Mysql) GetRows(sql string, args ...interface{}) []map[string]string {
	db := this.Db
	this.Sql = sql
	rets := make([]map[string]string, 0)
	rows, err := db.Query(sql, args)
	if err != nil {
		fmt.Println("Insert err:" + err.Error())
		return rets
	}

	//字典类型
	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		rets = append(rets, record)
		fmt.Println(record)
	}
	return rets
}

//查询单行封装
func (this *Mysql) GetRow(sql string, args ...interface{}) map[string]string {
	db := this.Db
	this.Sql = sql
	record := make(map[string]string)
	rows, err := db.Query(sql, args...)
	if err != nil {
		fmt.Println("query err:" + err.Error())
		return record
	}

	//字典类型
	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		fmt.Println(record)
		break
	}
	return record
}

//更新数据
func (this *Mysql) Update(sql string, args ...interface{}) (int64, error) {

	db := this.Db
	this.Sql = sql
	stmt, err := db.Prepare(sql)
	if err != nil {
		fmt.Println("Update err:" + err.Error())
		return 0, err
	}
	res, err := stmt.Exec(args...)
	if err != nil {
		fmt.Println("Update err:" + err.Error())
		return 0, err
	}
	return res.RowsAffected()

}

//删除数据
func (this *Mysql) Remove(sql string, args ...interface{}) (int64, error) {

	db := this.Db
	this.Sql = sql
	stmt, err := db.Prepare(sql)
	if err != nil {
		fmt.Println("Remove err:" + err.Error())
		return 0, err
	}
	res, err := stmt.Exec(args...)
	if err != nil {
		fmt.Println("Remove err:" + err.Error())
		return 0, err
	}
	return res.RowsAffected()

}

func (this *Mysql) checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

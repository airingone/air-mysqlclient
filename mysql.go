package air_mysqlclient

import (
	"database/sql"
	"errors"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	_ "github.com/go-sql-driver/mysql"
)

//mysql client封装

//mysql client
type MysqlClient struct {
	db           *sql.DB //db对象，协程安全的
	addr         string  //地址，如"root:test.2020@tcp(127.0.0.1:3306)/air_test?charset=utf8"
	maxIdleConns uint32  //最大空闲连接
	maxOpenConns uint32  //最大打开连接数
}

//创建db client
//config: 配置
func NewMysqlClient(addr string, maxIdleConns uint32, maxOpenConns uint32) (*MysqlClient, error) {
	client := &MysqlClient{
		addr:         addr,
		maxIdleConns: maxIdleConns,
		maxOpenConns: maxOpenConns,
	}

	err := client.open()
	if err != nil {
		return nil, err
	}

	return client, nil
}

//open
func (cli *MysqlClient) open() error {
	db, err := sql.Open("mysql", cli.addr)
	if err != nil {
		return err
	}
	if cli.maxIdleConns > 0 {
		db.SetMaxIdleConns(int(cli.maxIdleConns))
	}
	if cli.maxOpenConns > 0 {
		db.SetMaxOpenConns(int(cli.maxOpenConns))
	}

	cli.db = db

	return nil
}

//close
func (cli *MysqlClient) Close() {
	if cli.db != nil {
		_ = cli.db.Close()
	}
}

//Insert
//tableName: 数据库表名
//values: 数据key为数据表列名称，values为写入值，如{"c_userid", "123456"}
func (cli *MysqlClient) Insert(tableName string, values map[string]interface{}) error {
	if values == nil || len(values) == 0 {
		return errors.New("value is empty")
	}

	//build insert sql
	var datas []map[string]interface{}
	datas = append(datas, values)
	sql, vals, err := builder.BuildInsert(tableName, datas)
	if err != nil {
		return err
	}

	//mysql insert
	/*res*/
	_, err = cli.db.Exec(sql, vals...)
	if err != nil {
		return err
	}
	//res.LastInsertId()

	return nil
}

//Query
//tableName: 数据库表名
//where: 数据key为数据表列名称，values查询条件值，如{"c_userid", "123456"}
//fileds: 需要读取的数据列，如["c_id", "c_userid"],	全部则["*"]
//offset: limit参数的开始位置
//limit: limit参数的拉取数
//result: 返回对象，不同表或获取的列不同则需要不同的struct对象，数据项必须有ddb标志，如 UserId string `ddb:"c_userid"`
//_orderby实现: where["_orderby"] = "c_id desc"，where["_orderby"] = "c_id asc"
//_groupby实现: where["_groupby"] = "c_userid"
//_having实现:
func (cli *MysqlClient) Query(tableName string, where map[string]interface{}, fileds []string, offset uint32, limit uint32, result interface{}) error {
	if limit == 0 {
		limit = 1
	}
	where["_limit"] = []uint{uint(offset), uint(limit)}

	sql, vals, err := builder.BuildSelect(tableName, where, fileds)
	if err != nil {
		return err
	}
	rows, err := cli.db.Query(sql, vals...)
	if err != nil {
		return err
	}

	err = scanner.ScanClose(rows, result)
	if err != nil {
		return err
	}

	return nil
}

//Count
//tableName: 数据库表名
//where: 数据key为数据表列名称，values查询条件值，如{"c_userid", "123456"}
func (cli *MysqlClient) QueryCount(tableName string, where map[string]interface{}) (uint32, error) {
	fileds := []string{"count(*)"}
	sql, vals, err := builder.BuildSelect(tableName, where, fileds)
	if err != nil {
		return 0, err
	}
	rows, err := cli.db.Query(sql, vals...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return uint32(count), err
}

//Update
//tableName: 数据库表名
//where: 数据key为数据表列名称，values查询条件值，如{"c_userid", "123456"}
//values: 需要修改的值，如{"c_userid", "123457"}
func (cli *MysqlClient) Update(tableName string, where map[string]interface{}, values map[string]interface{}) error {
	if where == nil || len(where) == 0 {
		return errors.New("where is empty")
	}
	sql, vals, err := builder.BuildUpdate(tableName, where, values)
	if err != nil {
		return err
	}

	/*res*/
	_, err = cli.db.Exec(sql, vals...)
	if err != nil {
		return err
	}
	//res.RowsAffected()

	return nil
}

//Delete
//tableName: 数据库表名
//where: 数据key为数据表列名称，values查询条件值，如{"c_userid", "123456"}
func (cli *MysqlClient) Delete(tableName string, where map[string]interface{}) error {
	if where == nil || len(where) == 0 {
		return errors.New("where is empty")
	}
	sql, vals, err := builder.BuildDelete(tableName, where)
	if err != nil {
		return err
	}

	/*res*/
	_, err = cli.db.Exec(sql, vals...)
	if err != nil {
		return err
	}
	//res.RowsAffected()

	return nil
}

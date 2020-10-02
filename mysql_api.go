package air_mysqlclient

//Insert
//configName: 配置文件名
//tableName: 数据库表名
//values: 数据key为数据表列名称，values为写入值，如{"c_userid", "123456"}
func MysqlInsert(configName string, tableName string, values map[string]interface{}) error {
	cli, err := GetMysqlClient(configName)
	if err != nil {
		return err
	}

	return cli.Insert(tableName, values)
}

//Query
//configName: 配置文件名
//tableName: 数据库表名
//where: 数据key为数据表列名称，values查询条件值，如{"c_userid", "123456"}
//fileds: 需要读取的数据列，如["c_id", "c_userid"],	全部则["*"]
//offset: limit参数的开始位置
//limit: limit参数的拉取数
//result: 返回对象，不同表或获取的列不同则需要不同的struct对象，数据项必须有ddb标志，如 UserId string `ddb:"c_userid"`
//_orderby实现: where["_orderby"] = "c_id desc"，where["_orderby"] = "c_id asc"
//_groupby实现: where["_groupby"] = "c_userid"
//_having实现:
func MysqlQuery(configName string, tableName string, where map[string]interface{}, fileds []string, offset uint32, limit uint32, result interface{}) error {
	cli, err := GetMysqlClient(configName)
	if err != nil {
		return err
	}

	return cli.Query(tableName, where, fileds, offset, limit, result)
}

//Count
//configName: 配置文件名
//tableName: 数据库表名
//where: 数据key为数据表列名称，values查询条件值，如{"c_userid", "123456"}
func MysqlQueryCount(configName string, tableName string, where map[string]interface{}) (uint32, error) {
	cli, err := GetMysqlClient(configName)
	if err != nil {
		return 0, err
	}

	return cli.QueryCount(tableName, where)
}

//Update
//configName: 配置文件名
//tableName: 数据库表名
//where: 数据key为数据表列名称，values查询条件值，如{"c_userid", "123456"}
//values: 需要修改的值，如{"c_userid", "123457"}
func MysqlUpdate(configName string, tableName string, where map[string]interface{}, values map[string]interface{}) error {
	cli, err := GetMysqlClient(configName)
	if err != nil {
		return err
	}

	return cli.Update(tableName, where, values)
}

//Delete
//configName: 配置文件名
//tableName: 数据库表名
//where: 数据key为数据表列名称，values查询条件值，如{"c_userid", "123456"}
func MysqlDelete(configName string, tableName string, where map[string]interface{}) error {
	cli, err := GetMysqlClient(configName)
	if err != nil {
		return err
	}

	return cli.Delete(tableName, where)
}

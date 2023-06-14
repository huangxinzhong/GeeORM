package dialect

import "reflect"

// 数据库名 -> 接口抽象
var dialectsMap = map[string]Dialect{}

// Dialect 多数据库兼容的抽象
type Dialect interface {
	// DataTypeOf 将 golang 类型转换为数据库的数据类型
	DataTypeOf(typ reflect.Value) string
	// TableExistSQL 返回某个表是否存在的 SQL 语句
	TableExistSQL(tablename string) (string, []interface{})
}

// RegisterDialect 注册数据库
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// GetDialect 获取已注册的数据库
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}

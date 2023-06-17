package schema

import (
	"GeeORM/dialect"
	"go/ast"
	"reflect"
)

// Field represents a column of database
type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema represents a table of database
type Schema struct {
	Model      interface{}       // 元数据，存储结构体
	Name       string            // 表名
	Fields     []*Field          // 字段类型
	FieldNames []string          // 字段名
	fieldMap   map[string]*Field // 字段名和字段类型的映射
}

// GetField returns field by name
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// RecordValues Values return the values of dest's member variables
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {

		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}

	return fieldValues
}

// ITableName 自定义表名接口
type ITableName interface {
	TableName() string
}

// Parse a struct to a Schema instance
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()

	// 自定义表名接口，如果实现，则使用自动以表名，为实现则使用结构体名作为表名
	var tableName string
	t, ok := dest.(ITableName)
	if !ok {
		tableName = modelType.Name()
	} else {
		tableName = t.TableName()
	}

	schema := &Schema{
		Model:    dest,
		Name:     tableName,
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		// p.Anonymous 是否嵌入字段
		// ast.IsExported(p.Name) 是否大写字母开头
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

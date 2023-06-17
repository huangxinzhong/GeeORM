package session

import (
	"GeeORM/log"
	"reflect"
)

// Hooks constants

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

// CallMethod calls the registered hooks
func (s *Session) CallMethod(method string, value interface{}) {
	// 为什么不使用 interface 实现 hook 机制
	/*
			 *type IBeforeQuery interface {
			 *      BeforeQuery(s *Session) error
			 *}
		     *
			 *type IAfterQuery interface {
			 *      AfterQuery(s *Session) error
			 *}
			 *.....
			 *等等
		     *
			 *然后修改CallMethod
			 *func (s *Session) CallMethod(method string, value interface{}) {
			 *	 ...
			 *     if i, ok := dest.(IBeforQuery); ok == true {
			 *        i. BeforeQuery(s)
			 *     }
			 *     ...
			 *	return
			 *}
	*/
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		if v := fm.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return
}

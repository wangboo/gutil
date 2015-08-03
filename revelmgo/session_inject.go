package revelmgo

import (
	"github.com/revel/revel"
	"labix.org/v2/mgo"
	"reflect"
)

// 创建revel_session 注入过滤器
func MgoSessionInjectFilterFunc(session *mgo.Session) func(c *revel.Controller, fc []revel.Filter) {
	return func(c *revel.Controller, fc []revel.Filter) {
		appCtrl := c.AppController
		typeOfC := reflect.TypeOf(appCtrl).Elem()
		_, ok := typeOfC.FieldByName("MSession")
		if !ok {
			fc[0](c, fc[1:])
			return
		}
		valueOfC := reflect.ValueOf(appCtrl).Elem()
		// 注入 session
		newSession := session.Clone()
		defer newSession.Close()
		valueOfSession := reflect.ValueOf(newSession)
		valoeOfElem := valueOfC.FieldByName("MSession")
		valoeOfElem.Set(valueOfSession)
		fc[0](c, fc[1:])
	}
}

// 重新生成带revel_session 注入过滤器的Filters
func NewFiltersWithMgoInjecter(filters []revel.Filter, session *mgo.Session) []revel.Filter {
	newFilters := []revel.Filter{}
	typeOfActionInvoker := reflect.TypeOf(revel.ActionInvoker)
	for _, o := range filters {
		if reflect.TypeOf(o) == typeOfActionInvoker {
			newFilters = append(newFilters, MgoSessionInjectFilterFunc(session))
		}
		newFilters = append(newFilters, o)
	}
	return newFilters
}

package lib

import (
	"errors"
	"reflect"

	"jryghq.cn/utils"
)

//服务定位器
type ServiceLocator struct {
	Object
	dict *Dict
}

func (self *ServiceLocator) ServiceLocator() *ServiceLocator {
	self.Object.Object(self)
	self.dict = new(Dict).Dict()
	return self
}

func (self *ServiceLocator) CheckService(com interface{}) bool {
	return self.dict.Check(reflect.TypeOf(com).Elem().Name())
}

//加入服务
func (self *ServiceLocator) AddService(obj interface{}) {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		utils.Log().WriteError("AddComponent != reflect.Ptr")
	}
	self.dict.Set(t.Elem().Name(), obj)
}

//移除服务
func (self *ServiceLocator) RemoveService(obj interface{}) {
	t := reflect.TypeOf(obj)
	self.dict.Del(t.Name())
}

//获取服务
func (self *ServiceLocator) Service(obj interface{}) error {
	t := reflect.TypeOf(obj)
	return self.dict.Get(t.Elem().Elem().Name(), obj)
}

//广播定位器内所有实现method的方法
func (self *ServiceLocator) Broadcast(method string, arg interface{}) error {
	list := self.dict.Keys()
	call := false
	for _, v := range list {
		com := self.dict.GetInterface(v)
		if com == nil {
			continue
		}
		value := reflect.ValueOf(com).MethodByName(method)
		if value.Kind() != reflect.Invalid {
			value.Call([]reflect.Value{reflect.ValueOf(arg)})
			call = true
		}
	}
	if !call {
		return errors.New("未找到 '" + method + "'方法")
	}
	return nil
}
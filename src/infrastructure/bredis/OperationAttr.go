package bredis

import (
	"fmt"
	"time"
)

type empty struct {

}

const (
	ATTR_EXPIRE="expr"
	ATTR_NX="nx"
	ATTR_XX="xx"
)

/**
	属性
 */
type OperationAttr struct {
	Name string
	Value interface{}
}

type OperationAttrs []*OperationAttr

//查找
func (this OperationAttrs) Find(name string) *InterfaceResult {
	for _,attr:=range this{
		if attr.Name == name {
			return NewInterfaceResult(attr.Value,nil)
		}
	}
	return NewInterfaceResult(nil,fmt.Errorf("OperationAttrs Found error:%s",name))
}

/**
	设置有效期
 */
func WithExpire(t time.Duration) *OperationAttr {
	return &OperationAttr{Name: ATTR_EXPIRE,Value: t}
}

//加锁

func WithNX() *OperationAttr {
	return &OperationAttr{Name: ATTR_NX,Value: empty{}}
}

//修锁
func WithXX() *OperationAttr {
	return &OperationAttr{Name: ATTR_XX,Value: empty{}}
}
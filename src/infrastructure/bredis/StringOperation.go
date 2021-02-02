package bredis

import (
	"context"
	"time"
)

/**
	专门处理string类型的库
*/
type StringOperation struct {
	ctx context.Context
}



func NewStringOperation() *StringOperation {
	return &StringOperation{ctx: context.Background()}
}
/**
	查找
 */

func (this *StringOperation) Set(key string,value interface{},attrs ...*OperationAttr) *InterfaceResult {
	exp:= OperationAttrs(attrs).Find(ATTR_EXPIRE).Unwrap_Or(time.Second*0).(time.Duration) //查找
	nx:= OperationAttrs(attrs).Find(ATTR_NX).Unwrap_Or(nil)
	if nx!=nil {
		return NewInterfaceResult(Redis().SetNX(this.ctx,key,value,exp).Result())
	}
	xx:= OperationAttrs(attrs).Find(ATTR_XX).Unwrap_Or(nil)
	if xx!=nil{
		return NewInterfaceResult(Redis().SetXX(this.ctx,key,value,exp).Result())
	}
	return NewInterfaceResult(Redis().Set(this.ctx,key,value,exp).Result())
}
/**
	进行解偶 StringResult
	方法：单key 取值
 */
func (this *StringOperation) Get(Key string) *StringResult {
	return NewStringResult(Redis().Get(this.ctx,Key).Result())
}
/**
	方法：多key 取值
	*SliceResult []切片
 */
func (this *StringOperation) MGet(Key... string) *SliceResult {
	return NewSliceResult(Redis().MGet(this.ctx,Key...).Result())
}
/**
	获取剩余时间
 */
func (this *StringOperation) Ttl(key string) *InterfaceResult {
	return NewInterfaceResult(Redis().TTL(this.ctx,key).Result())
}


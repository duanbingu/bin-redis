package DBCache

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/duanbingu/bin-redis/src/infrastructure/bredis"
	"time"
)
/**
	缓存组件
 */
const(
	Serilizer_JSON ="json"
	Serilizer_GOB="gob" //可以实现json所不能支持的struct方法序列化 只能go使用
)

type DBGetterFunc func() interface{} //存放DB sql 数据函数

type SimpleDBCache struct {
	Operation *bredis.StringOperation // 操作类库
	Expire    time.Duration           //过期时间
	DBGetter  DBGetterFunc            //缓存不存在则获取DB方法
	Serilizer string                  //序列化
	Policy    CachePolicy             //防穿透策略
}
// 构造函数
func NewSimpleDBCache(operation *bredis.StringOperation, expire time.Duration, serilizer string,policy CachePolicy) *SimpleDBCache {
	if policy !=nil {
		policy.SetOperation(operation) //设置 也可依赖注入，懒得加了。 没必要很复杂

	}
	return &SimpleDBCache{Operation: operation, Expire: expire, Serilizer: serilizer,Policy: policy}
}

//设置缓存
func (this *SimpleDBCache) SetDBCache(key string,value interface{}) {
	if this.Policy!=nil { //检查策略
		this.Policy.Before(key)
	}
	this.Operation.Set(key,value, bredis.WithExpire(this.Expire)).Unwrap()
}

//获取缓存
func (this *SimpleDBCache) GetStringDBCache(key string) (ret interface{}) {
	if this.Serilizer== Serilizer_JSON {
		//如果是序列化JSON
		f:=func() string{
			obj:=this.DBGetter()
			b,err:=json.Marshal(obj)
			if err!=nil {
				return ""
			}
			return string(b)
		}
		ret = this.Operation.Get(key).Unwrap_Or_Else(f)
	}else if this.Serilizer== Serilizer_GOB {
		f:= func() string{
			obj:=this.DBGetter()
			var buf=&bytes.Buffer{}
			enc:=gob.NewEncoder(buf)
			if err:=enc.Encode(obj);err!=nil {
				return ""
			}
			return buf.String()
		}
		ret = this.Operation.Get(key).Unwrap_Or_Else(f)
	}
	//校验数据库数据是否为空 并且开启了策略
	if ret.(string)=="" && this.Policy!=nil {  //执行ifnil策略
		this.Policy.IfNil(key,"")
	}else{
		this.SetDBCache(key, ret)
	}
	return
}

//转结构体
func(this *SimpleDBCache) GetDBCacheForObject(key string,obj interface{})  interface{} {
	ret:=this.GetStringDBCache(key)
	if ret==nil{
		return nil
	}

	if this.Serilizer== Serilizer_JSON {
		err:=json.Unmarshal([]byte(ret.(string)),obj)
		if err!=nil{
			return nil
		}
	}else if   this.Serilizer== Serilizer_GOB {
		var buf =&bytes.Buffer{}
		buf.WriteString(ret.(string))
		dec:=gob.NewDecoder(buf)
		if dec.Decode(obj)!=nil{
			return nil
		}
	}
	return nil

}

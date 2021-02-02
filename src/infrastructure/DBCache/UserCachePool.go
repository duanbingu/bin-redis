package DBCache

import (
	"github.com/duanbingu/bin-redis/src/infrastructure/bredis"
	"sync"
	"time"
)

/**
	初始化 NewSimpleDBCache 函数
	UserCache 策略
	减少GC Cache
 */

var NewCachePool *sync.Pool

func init()  {
	NewCachePool =&sync.Pool{
		New: func()	interface{} {
			//创建
			return NewSimpleDBCache(
				bredis.NewStringOperation(), //指定操作string类库，后期升级扩展
				time.Second*150,             //指定默认缓存过期时间
				Serilizer_JSON,              //指定序列化方式为json  Serilizer_JSON  gob 更改为Serilizer_GOB
				NewCrossPolicy("^user\\d{1,5}$",time.Second*30),//空包｜｜空key  设定缓存30秒 ID 最多5位数 开启策略
				//nil,// 关闭策略
			)
		},
	}
}
//获取
func NewUserDBCache() *SimpleDBCache {
	return NewCachePool.Get().(*SimpleDBCache)
}
//放回池中
func ReleaseUserNewDBCache(cache *SimpleDBCache)  {
	NewCachePool.Put(cache)
}




package DBCache

import (
	"github.com/duanbingu/bin-redis/src/infrastructure/bredis"
	"regexp"
	"time"
)

type CachePolicy interface {
	Before(key string) //检查keyID
	IfNil(key string,v interface{}) //检测是否为空
	SetOperation(opt *bredis.StringOperation) //操作类库
}

//缓存穿透策略
type CrossPolicy struct {
	KeyRegx string //检查key正则
	Expire time.Duration
	opt *bredis.StringOperation
}

func NewCrossPolicy(keyRegx string, expire time.Duration) *CrossPolicy {
	return &CrossPolicy{KeyRegx: keyRegx, Expire: expire}
}


//检测正则是否匹配
func (this *CrossPolicy) Before(key string) {
	if !regexp.MustCompile(this.KeyRegx).MatchString(key){
		panic("error cache key")
	}
}
/**
	设置包或参数为空 处理
 */
func(this *CrossPolicy) IfNil(key string,v interface{})  {

	this.opt.Set(key,v,bredis.WithExpire(this.Expire)).Unwrap()

}

func(this *CrossPolicy) SetOperation(opt *bredis.StringOperation){
	this.opt=opt
}
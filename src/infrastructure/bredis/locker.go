package bredis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Locker struct{
	key string
	//value interface{}
	expire time.Duration
	unlock bool
	incrScript *redis.Script
}

const incrLua=`
if redis.call('get', KEYS[1]) == ARGV[1] then
  return redis.call('expire', KEYS[1],ARGV[2]) 				
 else
   return '0' 					
end`

func NewLocker(key string, expire time.Duration) *Locker {
	if expire.Seconds()<=0 {
		panic("error expire")
	}
	return &Locker{key: key, expire: expire,incrScript: redis.NewScript(incrLua)}
}

//加锁🔒
func (this *Locker) Lock() *Locker {
	res:=NewStringOperation().Set(this.key,"",WithExpire(this.expire),WithNX())
	if res.Result!=true||res.Err!=nil {
		panic(fmt.Sprint("lock error with key %s",this.key))
	}
	this.expandLockTime()
	return this
}
//解锁🔒
func (this *Locker) Unlock() {
	this.unlock = true
	NewStringOperation().Del(this.key).Unwrap()
}
//协层续签时间
func (this *Locker) expandLockTime() {
	sleepTime:=this.expire.Seconds()*2 / 3 //每间隔3分之2
	go func() {
		for  {
			time.Sleep(time.Second*time.Duration(sleepTime)) //每间隔一秒重新设置过期时间
			if this.unlock{
				break
			}
			this.resetExpire()
		}
	}()
}
//重新设置过期时间
func (this *Locker) resetExpire() {
	cmd:=this.incrScript.Run(context.Background(),RedisClient,[]string{this.key},1,this.expire.Seconds())
	_,err:=cmd.Result()
	if err!=nil {
		panic(fmt.Sprint("lock expire error",err))
	}
	//log.Printf("key=%s,续期结果:%v,%v\n",this.key,err,v)
}

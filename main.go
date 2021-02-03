package main

import (
	"fmt"
	"github.com/duanbingu/bin-redis/src/infrastructure/DBCache"
	"github.com/duanbingu/bin-redis/src/infrastructure/bredis"
	"github.com/duanbingu/bin-redis/src/infrastructure/dao"
	"github.com/duanbingu/bin-redis/src/infrastructure/gorm"
	"github.com/duanbingu/bin-redis/src/models"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func main()  {
	/**
		锁说明：
		bredis.NewLocker("key",time.Second*3).Lock() 加锁 //保证同一时间只执行一次
		defer locker.Unlock() //解锁

		bin-redis string操作 说明
		bredis.RedisClient redis原生操作对象
		bredis.NewStringOperation().Set()
		set 设置 key value
		bredis.WithExpire(time.Second*20)过期时间
		bredis.WithNX() 加锁
		bredis.WithXX() //修锁
		bredis.NewStringOperation().Del(this.key) //删除
		bredis.NewStringOperation().Expire("key",time.Second*10) 续签过期时间
		bredis.NewStringOperation().Ttl(key) //获取剩余时间
		bredis.NewStringOperation().Get(key) 取值
		bredis.NewStringOperation().MGet("key","key") 取多值
		Unwrap 只取值
		Unwrap_Or() 值为空则复值
		bredis.NewStringOperation().MGet("key","key","key").Iter()
		for mget.HasNext(){
			fmt.Println(iter.Next())
		} 迭代器
	*/
	//操作案例
	fmt.Println(bredis.NewStringOperation().Set("name","bin-redis",bredis.WithExpire(time.Second*20))) //进行设置过去时间 未加锁
	fmt.Println(bredis.NewStringOperation().
		Set("names","bin-redis",
			bredis.WithExpire(time.Second*60),
			bredis.WithNX()).
		Unwrap()) //加锁SetNx
	fmt.Println(bredis.NewStringOperation().
		Set("names","bin-redis",
			bredis.WithExpire(time.Second*60),
			bredis.WithXX()).
		Unwrap()) //修锁SetXX
	fmt.Println(bredis.NewStringOperation().Ttl("names")) //获取剩余时间 ttl
	fmt.Println(bredis.NewStringOperation().Get("names")) //get
	fmt.Println(bredis.NewStringOperation().MGet("name","age")) //mget
	iter:=bredis.NewStringOperation().MGet("bin-redis","name","age").Iter() //迭代器
	for iter.HasNext(){
		fmt.Println(iter.Next())
	}
	/**
		缓存组件使用
		主要用于DB缓存操作
		策略加载机制
	 */
	r:=gin.New()
	r.Use(func(context *gin.Context) {
		defer func() {
			if e:=recover();e!=nil{
				context.JSON(400,gin.H{"message":e})
			}
		}()
		context.Next()
	})
	r.Handle("GET","/user/:user_id", func(context *gin.Context) {
		userID:= context.Param("user_id") //参数
		//1.对象池中取出缓存对象cache
		Cache:= DBCache.NewUserDBCache()       		//获取 SimpleCache 对象
		defer DBCache.ReleaseUserNewDBCache(Cache)    //pool 机制需PUT
		Cache.DBGetter= dao.NewUser().FindByID(userID) //设置DBGetter  sql
		//3.取缓存输出 如果没有该数据，上面的DBGetter会被调用去取数据库数据放倒缓存中并返回回来。
		//fmt.Print(Cache.GetStringDBCache("user"+userID))// 取string
		userModel:=&models.UserModel{} //结构体
		Cache.GetDBCacheForObject("user"+userID,userModel) //结构体赋值
		context.JSON(200,userModel)
	})
	r.Run(":8022")
	/**
		redis 加锁案例  注:etcd 也可加锁
		使用redis中 setnx 实现：指定key不存在时,为key设定指定的值 SETNX key value 设置成功 1 设置失败 0
		主要核心实现：抢锁, 使用lau实现原子性 续签时间 并通过协程进行检测
		使用方法：
			bredis.NewLocker("key",time.Second*3).Lock() 加锁 //保证同一时间只执行一次
			defer locker.Unlock() //解锁
		场景：
			问题：3个定时任务 同一时间进行 对 view 加1
			实现： 不论N个任务 只能每5秒 针对view +1
	*/
		Job("job1")
		Job("job2")
		select {} //不退出
}

// 任务 给user表中view 每间隔5秒加1
func Job(name string)  {

	c:=cron.New(cron.WithSeconds())
	//执行加锁
	id,err:=c.AddFunc("0/5 * * * * *", func() {
		defer func() {
			if e:=recover();e!=nil {
				log.Println(name,"执行失败",e)
			}
		}()
		lock:=bredis.NewLocker("job",time.Second*5).Lock()  //加锁
		defer lock.Unlock() //解锁
		time.Sleep(time.Second*2)
		db:=gorm.GormDB.Exec("update by_user set view=1+view where user_id =1737")
		if db.Error !=nil{
			panic(fmt.Sprint("db user view error",db.Error))
		}
		log.Println(name,"任务执行完成")
	})
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Printf("%s任务ID是:%d 启动\n",name,id)
	c.Start()

}


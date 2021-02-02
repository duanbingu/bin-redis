package main

import (
	"github.com/duanbingu/bin-redis/src/infrastructure/DBCache"
	"github.com/duanbingu/bin-redis/src/infrastructure/dao"
	"github.com/duanbingu/bin-redis/src/models"
	"github.com/gin-gonic/gin"
)

func main()  {
	//user:=&models.UserModel{}
	//gorm.GormDB.Table("by_user").Find(&user)
	//fmt.Println(user)
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
		bin-redis string操作
	 */
	//fmt.Println(bin_redis.NewStringOperation().Set("name","bin",bin_redis.WithExpire(time.Second*20))) //进行设置过去时间 未加锁
	//fmt.Println(bin_redis.NewStringOperation().
	//	Set("names","bin",
	//		bin_redis.WithExpire(time.Second*60),
	//		bin_redis.WithNX()).
	//	Unwrap()) //加锁SetNx
	//fmt.Println(bin_redis.NewStringOperation().
	//	Set("names","bin",
	//		bin_redis.WithExpire(time.Second*60),
	//		bin_redis.WithXX()).
	//	Unwrap()) //修锁SetXX
	//fmt.Println(bin_redis.NewStringOperation().Ttl("names")) //获取剩余时间 ttl
	//fmt.Println(bin_redis.NewStringOperation().Get("names")) //get
	//fmt.Println(bin_redis.NewStringOperation().MGet("name","age").Unwrap()) //mget
	//iter:=bin_redis.NewStringOperation().MGet("阿斯顿","name","age").Iter() //迭代器
	//for iter.HasNext(){
	//	fmt.Println(iter.Next())
	//}
}
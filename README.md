package bin_redis

安装命令

go get -u github.com/duanbingu/bin-redis@V0.2

目前实现了 redis字符串操作 锁 lua原子性 续期等功能 ,mysql缓存组件
  
规划：增加集群,Hash List 等支持 

具体看main.go 实例 手册年后出。

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
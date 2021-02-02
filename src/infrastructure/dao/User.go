package dao

import (
	"github.com/duanbingu/bin-redis/src/infrastructure/DBCache"
	"github.com/duanbingu/bin-redis/src/infrastructure/gorm"
	"github.com/duanbingu/bin-redis/src/models"
	"log"
)
/**
	用户DAO
 */
type User struct {

}

func NewUser() *User {
	return &User{}
}
//redis 操作
func (this *User) FindByID(UserID string) DBCache.DBGetterFunc {
	return func() interface{} {
		log.Println("get form db")
		user:=&models.UserModel{}
		if gorm.GormDB.Table("by_user").Where("user_id=?",UserID).Find(user).Error!=nil || user.UserID<=0{
			return nil //查询为空走策略缓存 防穿透
		}
		return user
	}
}
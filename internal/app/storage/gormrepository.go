package storage

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type GormRep struct {}

func (g *GormRep) OpenGorm() *gorm.DB {
	//dsn := "host=localhost user=selectel password=selectel dbname=selectel port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := "postgres://selectel:selectel@127.0.0.1:5432/selectel?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Gorm() *gorm.DB {
	g := GormRep{}
	return g.OpenGorm()
}

func GetById(o interface{}, id int32) *gorm.DB  {
	return Gorm().Debug().Table("event").First(o, "id = ?", id)
}

func Create(o interface{}) *gorm.DB {
	return Gorm().Debug().Table("event").Create(o)
}

func Delete(o interface{}, id int32) error {
	if err := Gorm().Debug().Table("event").Where("id = ?", id).Delete(o).Error; err != nil {
		return err
	}
	return nil
}

func All(o interface{}) *gorm.DB {
	return Gorm().Debug().Table("event").Find(o)
}

func Update(newData interface{}, id int) *gorm.DB {
	return Gorm().Debug().Table("event").Where("id = ?", id).Updates(newData)
}
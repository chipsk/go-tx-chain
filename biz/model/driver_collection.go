package model

import (
	"chipsk/go-tx-chain/common/mysql"
	"fmt"
)

type DriverCollection struct {
	ID     int64  `gorm:"column:id" db:"id" json:"id" form:"id"`
	Uid    int64  `gorm:"column:uid" db:"uid" json:"uid" form:"uid"`
	Phone  string `gorm:"column:phone" db:"phone" json:"phone" form:"phone"`
	Status int    `gorm:"column:status" db:"status" json:"status" form:"status"`
}

func (d *DriverCollection) TableName() string {
	//分库分表场景
	index := mysql.TableIndexByUid(d.Uid)
	return fmt.Sprintf("driver_collection_%d", index)
}

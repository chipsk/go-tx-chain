package mysql

import (
	"chipsk/go-tx-chain/common/mysql/infra"
	"chipsk/go-tx-chain/conf"
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"time"
)

var (
	Client *PoolClient

	NewClient *PoolClient
)

type PoolClient struct {
	name string
	pool *infra.Pool
}

type PooledConnection struct {
	address string
	Conn    *gorm.DB
}

func (p *PooledConnection) Close() {
	sqlDB, _ := p.Conn.DB()
	sqlDB.Close()
}

func Init() error {
	host := conf.Viper.GetString("mysql.host")
	Client = &PoolClient{name: "mysql"}
	pool := infra.NewPool(&infra.Setting{
		Hosts: []string{host},
	}, Client.Dial)

	Client.pool = pool
	err := Client.pool.Run()
	if err != nil {
		log.Fatalf("mysql init err %v ", err)
		os.Exit(1)
	}
	return err
}

// pool 动态调用Dial 获取连接对象
func (d *PoolClient) Dial(address string) (pc infra.NodeShard, err error) {
	dsn := conf.Viper.GetString("mysql.dsn")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// log error
		return nil, err
	}
	pc = &PooledConnection{
		address: address,
		Conn:    db,
	}

	return pc, nil
}

func (p *PooledConnection) Validate() (err error) {
	conn, err := net.DialTimeout("tcp", p.address, 1*time.Second)
	if err == nil {
		conn.Close()
	}
	return err
}

type Option func(db *gorm.DB) *gorm.DB

func (d *PoolClient) GetDB(ctx context.Context, option ...Option) (*gorm.DB, error) {
	conner, err := d.pool.Get()
	if err != nil {
		//log error
		return nil, err
	}
	db := conner.(*PooledConnection)
	conn := db.Conn

	for _, op := range option {
		conn = op(db.Conn)
	}

	return conn, nil
}

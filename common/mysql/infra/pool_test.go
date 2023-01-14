package infra

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"testing"
	"time"
)

// Implement pool.Pooled & redis.Conn interface
type PooledConnection struct {
	address string
	Conn    *gorm.DB
}

func (p *PooledConnection) Validate() (err error) {
	conn, err := net.DialTimeout("tcp", p.address, 1*time.Second)
	if err == nil {
		conn.Close()
	}
	return err
}

func (p *PooledConnection) Close() {
	sqlDB, _ := p.Conn.DB()
	sqlDB.Close()
}

func New(address string) (pc NodeShard, err error) {
	fmt.Printf("address %v \n", address)
	dsn := fmt.Sprintf("root:root@(%s)/sk?charset=utf8&parseTime=True&loc=Local", address)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("conn mysql %v \n", err)
	}

	pc = &PooledConnection{
		Conn: db,
	}
	return pc, nil
}

func TestNewDefaultPool(t *testing.T) {

	p := NewDefaultPool([]string{"127.0.0.1:3306"}, New)

	err := p.Run()

	if err != nil {
		log.Fatalf("pool init err %v \n", err)
	}

	res, err := p.Get()
	if err != nil || res == nil {
		t.Fatal()
	}

	//test retry
	p.services = []string{"127.0.0.1:3306"}
	p.checkAndUpdateMap()
	p.Get()
}

func TestNewPool(t *testing.T) {

	p := NewPool(&Setting{
		Hosts: []string{"127.0.0.1:3306"},
	}, New)

	err := p.Run()

	if err != nil {
		log.Fatalf("pool init err %v \n", err)
	}

	res, err := p.Get()
	if err != nil || res == nil {
		t.Fatal()
	}
}

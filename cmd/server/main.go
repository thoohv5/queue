package main

import (
	"fmt"
	"math/rand"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/thoohv5/queue/model"
	"github.com/thoohv5/queue/service/trick"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {

	dsn := "root:password@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		QueryFields: true,
	})
	if nil != err {
		panic(err)
	}

	// q := queue.New(db)
	// for i := 0; i < 1000000; i++ {
	// 	fmt.Println(q.SendMessage(fmt.Sprintf("未婚夫%d", i)))
	// }

	t := trick.New(db)

	t.Register(func(msg *model.Queue) error {
		fmt.Println(msg)
		// time.Sleep(time.Duration(rand.Int63n(1000) * 1000))
		// if r := rand.Intn(10); r/2 == 0 {
		return nil
		// }
		// return errors.New("fail")
	})

	for i := 0; i < 5; i++ {
		go t.Run()
	}
	select {}
}

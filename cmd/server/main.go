package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/thoohv5/queue/service/queue"
	"github.com/thoohv5/queue/service/trick"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {

	dsn := "root:password@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if nil != err {
		panic(err)
	}

	// q := queue.New(db)
	// for i := 0; i < 10000; i++ {
	// 	fmt.Println(q.SendMessage(context.Background(), fmt.Sprintf("未婚夫%d", i)))
	// }

	t := trick.New(db)

	t.Register(func(msg *queue.Entity) error {
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

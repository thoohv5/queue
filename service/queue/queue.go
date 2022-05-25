package queue

import (
	"errors"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"github.com/thoohv5/queue/model"
	"github.com/thoohv5/queue/util"
)

type (
	IQueue interface {
		SendMessage(msg string) (msgId uint64, err error)

		Pull() (msg *model.Queue, err error)
		Success(queueId uint64) (err error)
		Fail(queueId uint64) (err error)

		Clear() (err error)
		Reset(queueId uint64) (err error)
	}
	queue struct {
		db *gorm.DB
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	NotMsg = "queue not msg"
)

func New(db *gorm.DB) IQueue {
	return &queue{
		db: db,
	}
}

func (q *queue) SendMessage(msg string) (msgId uint64, err error) {
	qm := &model.Queue{
		Data: msg,
	}
	return qm.ID, q.db.Create(qm).Error
}

func (q *queue) Pull() (msg *model.Queue, err error) {

	// 取出100组数据
	matchMsgIdList := make([]uint64, 0, 100)
	err = q.db.Model(&model.Queue{}).Where("lock_id", "").Where("retry_times < ?", 3).Order("id asc, retry_times asc").Limit(100).Pluck("id", &matchMsgIdList).Error
	if nil != err && err != gorm.ErrRecordNotFound {
		return
	}

	if len(matchMsgIdList) == 0 {
		err = errors.New(NotMsg)
		return
	}

	// 随机选取一个
	chooseMsgId := matchMsgIdList[rand.Intn(len(matchMsgIdList))]

	// 数据加锁
	err = q.db.Model(&model.Queue{}).Where("id", chooseMsgId).Update("lock_id", util.GoID()).Error
	if nil != err {
		return
	}

	// 获取数据
	msg = new(model.Queue)
	err = q.db.Model(&model.Queue{}).Where("id", chooseMsgId).Where("lock_id", util.GoID()).First(msg).Error
	if nil != err {
		return
	}

	return
}

func (q *queue) Clear() (err error) {
	return q.db.Model(&model.Queue{}).Where("lock_id > ?", 0).Where("updated_at < ?", time.Now().Add(-3*time.Minute)).Updates(map[string]interface{}{
		"lock_id":     "",
		"retry_times": gorm.Expr("retry_times + ?", 1),
	}).Error
}

func (q *queue) Reset(queueId uint64) (err error) {
	return q.db.Model(&model.Queue{}).Where("id", queueId).Updates(map[string]interface{}{
		"lock_id":     "",
		"retry_times": 0,
	}).Error
}

func (q *queue) Success(queueId uint64) (err error) {
	return q.db.Where("id", queueId).Delete(&model.Queue{}).Error
}

func (q *queue) Fail(queueId uint64) (err error) {
	return q.db.Model(&model.Queue{}).Where("id", queueId).Updates(map[string]interface{}{
		"lock_id":     "",
		"retry_times": gorm.Expr("retry_times + ?", 1),
	}).Error
}

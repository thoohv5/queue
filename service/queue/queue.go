package queue

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/thoohv5/queue/util"
)

type (
	IQueue interface {
		SendMessage(ctx context.Context, msg string) (msgId uint64, err error)

		Pull(ctx context.Context) (msg *Entity, err error)
		Success(ctx context.Context, queueId uint64) (err error)
		Fail(ctx context.Context, queueId uint64) (err error)

		Clear(ctx context.Context) (err error)
		Reset(ctx context.Context, queueId uint64) (err error)
	}
	queue struct {
		db *sql.DB
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	NotMsg = "queue not msg"
)

func New(db *sql.DB) IQueue {
	return &queue{
		db: db,
	}
}

func (q *queue) SendMessage(ctx context.Context, msg string) (msgId uint64, err error) {
	res, err := q.db.ExecContext(ctx, "INSERT INTO queue(`data`) values(?)", msg)
	if nil != err {
		return
	}
	mi, err := res.LastInsertId()
	return uint64(mi), err
}

func (q *queue) Pull(ctx context.Context) (msg *Entity, err error) {

	// 取出100组数据
	matchMsgIdList := make([]uint64, 0, 100)
	rrows, err := q.db.QueryContext(ctx, "SELECT id FROM `queue` WHERE lock_id = '' AND retry_times < ? ORDER BY id asc, retry_times asc LIMIT 100", 3)
	if nil != err {
		return
	}

	for rrows.Next() {
		var matchMsgId uint64
		if err = rrows.Scan(&matchMsgId); nil != err {
			return
		}
		matchMsgIdList = append(matchMsgIdList, matchMsgId)
	}

	if len(matchMsgIdList) == 0 {
		err = errors.New(NotMsg)
		return
	}

	// 随机选取一个
	chooseMsgId := matchMsgIdList[rand.Intn(len(matchMsgIdList))]

	// 数据加锁
	_, err = q.db.ExecContext(ctx, "UPDATE `queue` SET  lock_id = ? WHERE id = ?", util.GoID(), chooseMsgId)
	if nil != err {
		return
	}

	// 获取数据
	msg = new(Entity)
	qrow := q.db.QueryRowContext(ctx, "SELECT id, data FROM queue WHERE id = ? AND lock_id = ?", chooseMsgId, util.GoID())
	if err = qrow.Scan(&msg.ID, &msg.Data); nil != err {
		return
	}

	return
}

func (q *queue) Clear(ctx context.Context) (err error) {
	_, err = q.db.ExecContext(ctx, "UPDATE queue SET lock_id = '', retry_times = retry_times + ? WHERE lock_id > ? AND updated_at < ? ", 1, 0, time.Now().Add(-3*time.Minute))
	return
}

func (q *queue) Reset(ctx context.Context, queueId uint64) (err error) {
	_, err = q.db.ExecContext(ctx, "UPDATE queue SET lock_id = '', retry_times = 0 WHERE id = ?", queueId)
	return
}

func (q *queue) Success(ctx context.Context, queueId uint64) (err error) {
	_, err = q.db.ExecContext(ctx, "DELETE FROM queue WHERE id = ?", queueId)
	return
}

func (q *queue) Fail(ctx context.Context, queueId uint64) (err error) {
	_, err = q.db.ExecContext(ctx, "UPDATE queue SET lock_id = '', retry_times = retry_times + ?, WHERE id = ?", 1, queueId)
	return
}

package trick

import (
	"context"
	"database/sql"
	"time"

	"github.com/thoohv5/queue/log"
	"github.com/thoohv5/queue/service/queue"
	"github.com/thoohv5/queue/util"
)

type (
	ITrick interface {
		Register(execute func(msg *queue.Entity) error)
		Run()
	}
	Execute func(msg *queue.Entity) error
	trick   struct {
		q   queue.IQueue
		e   Execute
		log log.ILog
	}
)

func New(sos ...ServerOption) ITrick {
	opts := new(Options)
	for _, so := range sos {
		so(opts)
	}

	return &trick{
		q:   queue.New(opts.db),
		log: log.GetLogger(opts.log, opts.enableLogger),
	}
}

// Options 可选参数列表
type Options struct {
	db *sql.DB

	enableLogger bool
	log          log.ILog
}

// ServerOption 为可选参数赋值的函数
type ServerOption func(*Options)

// WithDB 是否为写库
func WithDB(db *sql.DB) ServerOption {
	return func(o *Options) {
		o.db = db
	}
}

func WithLogger(log log.ILog) ServerOption {
	return func(o *Options) {
		o.enableLogger = true
		o.log = log
	}
}

func (t *trick) Register(execute func(msg *queue.Entity) error) {
	t.e = execute
}

func (t *trick) Run() {
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			ctx := context.Background()
			t.clear(ctx)
			t.log.Debug("trick start, lock_id:%d", util.GoID())
			if t.e == nil {
				t.log.Debug("trick not execute")
				t.log.Debug("trick end, lock_id:%d", util.GoID())
				continue
			}
			// 拉取消息
			msg, err := t.q.Pull(ctx)
			if nil != err {
				// 不存在消息
				if err.Error() != queue.NotMsg {
					t.log.Error("trick queue pull err, err:+v", err)
				}
				t.log.Debug("trick not msg")
				t.log.Debug("trick end, lock_id:%d", util.GoID())
				continue
			}

			// 保护执行过程
			func(msg *queue.Entity) {
				defer func() {
					if re := recover(); re != nil {
						t.log.Error("exec msg err, re:%+v, msg:%+v", re, msg)
					}
					return
				}()
				// 执行消息
				err = t.e(msg)

				if nil != err {
					// 错误处理
					if err = t.q.Fail(ctx, msg.ID); nil != err {
						t.log.Error("trick queue fail err, err:%+v, msgId:%d", err, msg.ID)
					}
				} else {
					// 成功处理
					if err = t.q.Success(ctx, msg.ID); nil != err {
						t.log.Error("trick queue success err, err:%+v, msgId:%d", err, msg.ID)
					}
				}

			}(msg)

			t.log.Debug("trick end, lock_id:%d", util.GoID())
		}
	}

}

func (t *trick) clear(ctx context.Context) {
	if err := t.q.Clear(ctx); nil != err {
		t.log.Error("trick clear err, err:%+v", err)
	}
}

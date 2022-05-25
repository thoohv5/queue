package trick

import (
	"time"

	"gorm.io/gorm"

	"github.com/thoohv5/queue/log"
	"github.com/thoohv5/queue/model"
	"github.com/thoohv5/queue/service/queue"
	"github.com/thoohv5/queue/util"
)

type (
	ITrick interface {
		Register(execute func(msg *model.Queue) error)
		Run()
	}
	Execute func(msg *model.Queue) error
	trick   struct {
		q queue.IQueue
		e Execute
	}
)

func New(db *gorm.DB) ITrick {
	return &trick{
		q: queue.New(db),
	}
}

func (t *trick) Register(execute func(msg *model.Queue) error) {
	t.e = execute
}

func (t *trick) Run() {
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			t.clear()
			log.GetLogger().Debug("trick start, lock_id:%d", util.GoID())
			if t.e == nil {
				log.GetLogger().Debug("trick not execute")
				log.GetLogger().Debug("trick end, lock_id:%d", util.GoID())
				continue
			}
			// 拉取消息
			msg, err := t.q.Pull()
			if nil != err {
				// 不存在消息
				if err.Error() != queue.NotMsg {
					log.GetLogger().Error("trick queue pull err, err:+v", err)
				}
				log.GetLogger().Debug("trick not msg")
				log.GetLogger().Debug("trick end, lock_id:%d", util.GoID())
				continue
			}

			// 保护执行过程
			func(msg *model.Queue) {
				defer func() {
					if re := recover(); re != nil {
						log.GetLogger().Error("exec msg err, re:%+v, msg:%+v", re, msg)
					}
					return
				}()
				// 执行消息
				err = t.e(msg)

				if nil != err {
					// 错误处理
					if err = t.q.Fail(msg.ID); nil != err {
						log.GetLogger().Error("trick queue fail err, err:%+v, msgId:%d", err, msg.ID)
					}
				} else {
					// 成功处理
					if err = t.q.Success(msg.ID); nil != err {
						log.GetLogger().Error("trick queue success err, err:%+v, msgId:%d", err, msg.ID)
					}
				}

			}(msg)

			log.GetLogger().Debug("trick end, lock_id:%d", util.GoID())
		}
	}

}

func (t *trick) clear() {
	if err := t.q.Clear(); nil != err {
		log.GetLogger().Error("trick clear err, err:%+v", err)
	}
}

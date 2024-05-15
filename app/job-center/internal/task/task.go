package task

import (
	"github.com/go-co-op/gocron"
	"job-center/internal/logic"
	"job-center/internal/svc"
	"time"
)

type Task struct {
	s   *gocron.Scheduler
	ctx *svc.ServiceContext
}

func NewTask(ctx *svc.ServiceContext) *Task {
	return &Task{
		s:   gocron.NewScheduler(time.UTC),
		ctx: ctx,
	}
}

func (t *Task) Run() {
	t.s.Every(1).Minute().Do(func() {
		// 每分钟拉取数据
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.CacheRedis).Do("1m")
	})
	t.s.Every(5).Minute().Do(func() {
		// 每分钟拉取数据
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.CacheRedis).Do("5m")
	})
	t.s.Every(15).Minute().Do(func() {
		// 每分钟拉取数据
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.CacheRedis).Do("15m")
	})
	t.s.Every(30).Minute().Do(func() {
		// 每分钟拉取数据
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.CacheRedis).Do("30m")
	})

	t.s.Every(1).Hour().Do(func() {
		// 每小时拉取数据
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.CacheRedis).Do("1H")
	})
	t.s.Every(3).Hour().Do(func() {
		// 每小时拉取数据
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.CacheRedis).Do("2H")
	})
	t.s.Every(5).Hour().Do(func() {
		// 每小时拉取数据
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.CacheRedis).Do("4H")
	})

	t.s.Every(1).Day().Do(func() {
		// 每小时拉取数据
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.CacheRedis).Do("1D")
	})
	t.s.Every(1).Week().Do(func() {
		// 每小时拉取数据
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.CacheRedis).Do("1W")
	})
	t.s.Every(1).Month().Do(func() {
		// 每小时拉取数据
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.CacheRedis).Do("1M")
	})

	t.s.Every(1).Minute().Do(func() {
		// 每小时拉取数据
		logic.NewRate(t.ctx.Config.Okx, t.ctx.CacheRedis).Do()
	})
	// 十分钟生成一个区块
	t.s.Every(10).Minute().Do(func() {
		logic.NewBitCoin(t.ctx.Config.Bitcoin.Url, t.ctx.CacheRedis, t.ctx.AssetRpc, t.ctx.MongoClient, t.ctx.KafkaClient).Do()
	})
}

// StartBlocking 开启所有任务，同时阻塞当前线程
func (t *Task) StartBlocking() {
	t.s.StartBlocking()
}

func (t *Task) Stop() {
	t.s.Stop()
}

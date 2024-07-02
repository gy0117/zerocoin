package logic

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"job-center/internal/config"
	"job-center/internal/db"
	"job-center/internal/domain"
	"strings"
	"sync"
	"time"

	"zero-common/kafka"
	"zero-common/tools"
)

const success = "0"

const bar_1m = "1m"

type OkxResult struct {
	Code string     `json:"code"`
	Msg  string     `json:"msg"`
	Data [][]string `json:"data"`
}

type Kline struct {
	wg          sync.WaitGroup
	conf        config.OkxConfig
	klineDomain *domain.KlineDomain
	queueDomain *domain.QueueDomain
	cache       cache.Cache
}

func NewKline(conf config.OkxConfig, client *db.MongoClient, kCli *kafka.KafkaClient, redis cache.Cache) *Kline {
	return &Kline{
		wg:          sync.WaitGroup{},
		conf:        conf,
		klineDomain: domain.NewKlineDomain(client),
		queueDomain: domain.NewQueueDomain(kCli),
		cache:       redis,
	}
}

func (k *Kline) Do(duration string) {
	logx.Info("job-center | klineLogic | start pull K-line data, duration: ", duration)

	k.wg.Add(2)
	go k.getKlineData("BTC-USDT", "BTC/USDT", duration)
	go k.getKlineData("ETH-USDT", "ETH/USDT", duration)
	k.wg.Wait()
}

// instId: 产品Id
// bar: 时间粒度
// 示例：/api/v5/market/candles?instId=BTC-USDT&bar=1m
func (k *Kline) getKlineData(instId, symbol, bar string) {
	url := k.conf.Host + "/api/v5/market/candles?instId=" + instId + "&bar=" + bar

	timestamp := tools.ISO(time.Now())
	sign := tools.ComputeHmacSha256(timestamp+"GET"+"/api/v5/market/candles?instId="+instId+"&bar="+bar, k.conf.SecretKey)
	header := make(map[string]string)
	header["OK-ACCESS-KEY"] = k.conf.ApiKey
	header["OK-ACCESS-SIGN"] = sign
	header["OK-ACCESS-TIMESTAMP"] = timestamp
	header["OK-ACCESS-PASSPHRASE"] = k.conf.Passphrase

	// 访问外网需要添加代理
	bytes, err := tools.GetWithHeader(url, header, k.conf.Proxy)
	if err != nil {
		logx.Error(err)
		k.wg.Done()
		return
	}

	result := &OkxResult{}
	if err := json.Unmarshal(bytes, result); err != nil {
		logx.Error(err)
		k.wg.Done()
		return
	}

	if result.Code == success {
		// 将数据存入mongo db
		ctx := context.Background()
		err = k.klineDomain.SaveBatch(ctx, result.Data, symbol, bar)
		if err != nil {
			logx.Error(err)
		}

		if bar == bar_1m {
			// 发送给kafka，将最新的数据 推送给market服务，前端实时变化
			if len(result.Data) > 0 {

				data := result.Data[0]
				go k.queueDomain.Send1mKline(ctx, data, symbol)

				key := strings.ReplaceAll(instId, "-", "::")
				// BTC-USDT，收盘价格 BTC::USDT::RATE
				if err := k.cache.Set(key+"::RATE", data[4]); err != nil {
					logx.Error(err)
				}
			}
		}
	}

	k.wg.Done()
}

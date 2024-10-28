package processor

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/market"
	"market-api/internal/model"
	"market-api/internal/ws"
	"zero-common/tools"
)

var _ MarketHandler = (*WsHandler)(nil)

type WsHandler struct {
	wsServer *ws.WebSocketServer
}

func NewWsHandler(wsServer *ws.WebSocketServer) MarketHandler {
	return &WsHandler{
		wsServer: wsServer,
	}
}

func (h *WsHandler) HandleKline(symbol string, kline *model.Kline, thumbMap map[string]*market.CoinThumb) {
	logx.Info("WsHandler---HandleKline call...")
	logx.Info("symbol: ", symbol, ", closePrice: ", kline.ClosePrice, ", time: ", tools.ToTimeString(kline.Time))

	//logx.Info("WsHandler---HandleKline---end 消费数据")

	thumb := thumbMap[symbol]
	if thumb == nil {
		thumb = kline.InitCoinThumb(symbol)
	}
	coinThumb := kline.ToCoinThumb(symbol, thumb)
	result := &market.CoinThumb{}
	if err := copier.Copy(result, coinThumb); err != nil {
		logx.Error(err)
	}
	marshal, _ := json.Marshal(result)
	h.wsServer.BroadcastToRoom("/", "/topic/market/thumb", string(marshal))

	// 推送给币币交易-k线图
	bytes, _ := json.Marshal(kline)
	h.wsServer.BroadcastToRoom("/", "/topic/market/kline/"+symbol, string(bytes))

}

func (h *WsHandler) HandleTradePlate(symbol string, plate *model.TradePlateResult) {
	marshal, err := json.Marshal(plate)
	if err != nil {
		logx.Error(err)
		return
	}
	h.wsServer.BroadcastToRoom("/", "/topic/market/trade-plate/"+symbol, string(marshal))

	// 买卖盘展示，log
	showBuySellPlate(plate)
}

func showBuySellPlate(plate *model.TradePlateResult) {
	if len(plate.Items) == 0 {
		return
	}
	logx.Info("====== start 买卖盘展示 start ======")
	logx.Infof("Symbol: %s, Direction: %s, MaxAmount: %f, MinAmount: %f, HighestPrice: %f, LowestPrice: %f",
		plate.Symbol, plate.Direction, plate.MaxAmount, plate.MinAmount, plate.HighestPrice, plate.LowestPrice)
	for i := 0; i < len(plate.Items); i++ {
		item := plate.Items[i]
		logx.Infof("item-%d, price: %f, amount: %f", i, item.Price, item.Amount)
	}
	logx.Info("====== end 买卖盘展示 end ======")
}

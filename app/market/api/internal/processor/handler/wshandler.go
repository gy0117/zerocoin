package handler

import (
	"encoding/json"
	"fmt"
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
	fmt.Printf("买卖盘信息，推送到前端 | HandleTradePlate marshal: %s\n", string(marshal))
	h.wsServer.BroadcastToRoom("/", "/topic/market/trade-plate/"+symbol, string(marshal))
	logx.Info("WsHandler | HandleTradePlate 买卖盘通知: ", plate.Direction, " ", symbol, ":", fmt.Sprintf("%d", len(plate.Items)))
}

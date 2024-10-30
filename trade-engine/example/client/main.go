package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-common/trade-engine/types/match"
	"log"

	"time"
)

func main() {
	conn, _ := grpc.Dial("localhost:9083", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	client := match.NewMatchServiceClient(conn)

	for {
		sendBuyOrder(client)
		sendSellOrder(client)
		time.Sleep(time.Second * 2)
	}

	time.Sleep(time.Hour)
}

func sendBuyOrder(client match.MatchServiceClient) {
	req := &match.AddOrderRequest{Order: &match.Order{
		Id:        "1",
		Uid:       123,
		TradePair: "BTC/USDT",
		Price:     "120",
		Quantity:  "2",
		Side:      0,
		Type:      2,
	}}
	resp, err := client.AddOrder(context.Background(), req)
	log.Printf("【send buy】resp: {%v}, err: %v\n", resp, err)
}

func sendSellOrder(client match.MatchServiceClient) {
	req := &match.AddOrderRequest{Order: &match.Order{
		Id:        "2",
		Uid:       123,
		TradePair: "BTC/USDT",
		Price:     "120",
		Quantity:  "1",
		Side:      1,
		Type:      2,
	}}
	resp, err := client.AddOrder(context.Background(), req)
	log.Printf("【send sell】resp: {%v}, err: %v\n", resp, err)
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: order.proto

package order

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Order_GetHistoryOrder_FullMethodName = "/order.Order/GetHistoryOrder"
	Order_GetCurrentOrder_FullMethodName = "/order.Order/GetCurrentOrder"
	Order_AddOrder_FullMethodName        = "/order.Order/AddOrder"
	Order_FindByOrderId_FullMethodName   = "/order.Order/FindByOrderId"
	Order_CancelOrder_FullMethodName     = "/order.Order/CancelOrder"
)

// OrderClient is the client API for Order service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderClient interface {
	GetHistoryOrder(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderResp, error)
	GetCurrentOrder(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderResp, error)
	AddOrder(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*AddOrderResp, error)
	FindByOrderId(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*ExchangeOrder, error)
	CancelOrder(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*CancelOrderResp, error)
}

type orderClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderClient(cc grpc.ClientConnInterface) OrderClient {
	return &orderClient{cc}
}

func (c *orderClient) GetHistoryOrder(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderResp, error) {
	out := new(OrderResp)
	err := c.cc.Invoke(ctx, Order_GetHistoryOrder_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) GetCurrentOrder(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderResp, error) {
	out := new(OrderResp)
	err := c.cc.Invoke(ctx, Order_GetCurrentOrder_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) AddOrder(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*AddOrderResp, error) {
	out := new(AddOrderResp)
	err := c.cc.Invoke(ctx, Order_AddOrder_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) FindByOrderId(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*ExchangeOrder, error) {
	out := new(ExchangeOrder)
	err := c.cc.Invoke(ctx, Order_FindByOrderId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) CancelOrder(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*CancelOrderResp, error) {
	out := new(CancelOrderResp)
	err := c.cc.Invoke(ctx, Order_CancelOrder_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderServer is the server API for Order service.
// All implementations must embed UnimplementedOrderServer
// for forward compatibility
type OrderServer interface {
	GetHistoryOrder(context.Context, *OrderReq) (*OrderResp, error)
	GetCurrentOrder(context.Context, *OrderReq) (*OrderResp, error)
	AddOrder(context.Context, *OrderReq) (*AddOrderResp, error)
	FindByOrderId(context.Context, *OrderReq) (*ExchangeOrder, error)
	CancelOrder(context.Context, *OrderReq) (*CancelOrderResp, error)
	mustEmbedUnimplementedOrderServer()
}

// UnimplementedOrderServer must be embedded to have forward compatible implementations.
type UnimplementedOrderServer struct {
}

func (UnimplementedOrderServer) GetHistoryOrder(context.Context, *OrderReq) (*OrderResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHistoryOrder not implemented")
}
func (UnimplementedOrderServer) GetCurrentOrder(context.Context, *OrderReq) (*OrderResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCurrentOrder not implemented")
}
func (UnimplementedOrderServer) AddOrder(context.Context, *OrderReq) (*AddOrderResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddOrder not implemented")
}
func (UnimplementedOrderServer) FindByOrderId(context.Context, *OrderReq) (*ExchangeOrder, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByOrderId not implemented")
}
func (UnimplementedOrderServer) CancelOrder(context.Context, *OrderReq) (*CancelOrderResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelOrder not implemented")
}
func (UnimplementedOrderServer) mustEmbedUnimplementedOrderServer() {}

// UnsafeOrderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderServer will
// result in compilation errors.
type UnsafeOrderServer interface {
	mustEmbedUnimplementedOrderServer()
}

func RegisterOrderServer(s grpc.ServiceRegistrar, srv OrderServer) {
	s.RegisterService(&Order_ServiceDesc, srv)
}

func _Order_GetHistoryOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).GetHistoryOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Order_GetHistoryOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).GetHistoryOrder(ctx, req.(*OrderReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_GetCurrentOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).GetCurrentOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Order_GetCurrentOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).GetCurrentOrder(ctx, req.(*OrderReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_AddOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).AddOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Order_AddOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).AddOrder(ctx, req.(*OrderReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_FindByOrderId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).FindByOrderId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Order_FindByOrderId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).FindByOrderId(ctx, req.(*OrderReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_CancelOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).CancelOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Order_CancelOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).CancelOrder(ctx, req.(*OrderReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Order_ServiceDesc is the grpc.ServiceDesc for Order service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Order_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "order.Order",
	HandlerType: (*OrderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetHistoryOrder",
			Handler:    _Order_GetHistoryOrder_Handler,
		},
		{
			MethodName: "GetCurrentOrder",
			Handler:    _Order_GetCurrentOrder_Handler,
		},
		{
			MethodName: "AddOrder",
			Handler:    _Order_AddOrder_Handler,
		},
		{
			MethodName: "FindByOrderId",
			Handler:    _Order_FindByOrderId_Handler,
		},
		{
			MethodName: "CancelOrder",
			Handler:    _Order_CancelOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "order.proto",
}

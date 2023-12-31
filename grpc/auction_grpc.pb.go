// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.0
// source: auction.proto

package auction

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
	Auction_Bid_FullMethodName               = "/auction.Auction/Bid"
	Auction_Result_FullMethodName            = "/auction.Auction/Result"
	Auction_CallElection_FullMethodName      = "/auction.Auction/CallElection"
	Auction_AssertCoordinator_FullMethodName = "/auction.Auction/AssertCoordinator"
)

// AuctionClient is the client API for Auction service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuctionClient interface {
	Bid(ctx context.Context, in *SendBidMessage, opts ...grpc.CallOption) (*ResponseBidMessage, error)
	Result(ctx context.Context, in *RequestResultMessage, opts ...grpc.CallOption) (*ResultResponseMessage, error)
	CallElection(ctx context.Context, in *CallElectionMessage, opts ...grpc.CallOption) (*CallElectionResponseMessage, error)
	AssertCoordinator(ctx context.Context, in *AssertCoordinatorMessage, opts ...grpc.CallOption) (*AssertCoordinatorResponseMessage, error)
}

type auctionClient struct {
	cc grpc.ClientConnInterface
}

func NewAuctionClient(cc grpc.ClientConnInterface) AuctionClient {
	return &auctionClient{cc}
}

func (c *auctionClient) Bid(ctx context.Context, in *SendBidMessage, opts ...grpc.CallOption) (*ResponseBidMessage, error) {
	out := new(ResponseBidMessage)
	err := c.cc.Invoke(ctx, Auction_Bid_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionClient) Result(ctx context.Context, in *RequestResultMessage, opts ...grpc.CallOption) (*ResultResponseMessage, error) {
	out := new(ResultResponseMessage)
	err := c.cc.Invoke(ctx, Auction_Result_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionClient) CallElection(ctx context.Context, in *CallElectionMessage, opts ...grpc.CallOption) (*CallElectionResponseMessage, error) {
	out := new(CallElectionResponseMessage)
	err := c.cc.Invoke(ctx, Auction_CallElection_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionClient) AssertCoordinator(ctx context.Context, in *AssertCoordinatorMessage, opts ...grpc.CallOption) (*AssertCoordinatorResponseMessage, error) {
	out := new(AssertCoordinatorResponseMessage)
	err := c.cc.Invoke(ctx, Auction_AssertCoordinator_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuctionServer is the server API for Auction service.
// All implementations must embed UnimplementedAuctionServer
// for forward compatibility
type AuctionServer interface {
	Bid(context.Context, *SendBidMessage) (*ResponseBidMessage, error)
	Result(context.Context, *RequestResultMessage) (*ResultResponseMessage, error)
	CallElection(context.Context, *CallElectionMessage) (*CallElectionResponseMessage, error)
	AssertCoordinator(context.Context, *AssertCoordinatorMessage) (*AssertCoordinatorResponseMessage, error)
	mustEmbedUnimplementedAuctionServer()
}

// UnimplementedAuctionServer must be embedded to have forward compatible implementations.
type UnimplementedAuctionServer struct {
}

func (UnimplementedAuctionServer) Bid(context.Context, *SendBidMessage) (*ResponseBidMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bid not implemented")
}
func (UnimplementedAuctionServer) Result(context.Context, *RequestResultMessage) (*ResultResponseMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Result not implemented")
}
func (UnimplementedAuctionServer) CallElection(context.Context, *CallElectionMessage) (*CallElectionResponseMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CallElection not implemented")
}
func (UnimplementedAuctionServer) AssertCoordinator(context.Context, *AssertCoordinatorMessage) (*AssertCoordinatorResponseMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssertCoordinator not implemented")
}
func (UnimplementedAuctionServer) mustEmbedUnimplementedAuctionServer() {}

// UnsafeAuctionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuctionServer will
// result in compilation errors.
type UnsafeAuctionServer interface {
	mustEmbedUnimplementedAuctionServer()
}

func RegisterAuctionServer(s grpc.ServiceRegistrar, srv AuctionServer) {
	s.RegisterService(&Auction_ServiceDesc, srv)
}

func _Auction_Bid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendBidMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServer).Bid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auction_Bid_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServer).Bid(ctx, req.(*SendBidMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auction_Result_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestResultMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServer).Result(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auction_Result_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServer).Result(ctx, req.(*RequestResultMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auction_CallElection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CallElectionMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServer).CallElection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auction_CallElection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServer).CallElection(ctx, req.(*CallElectionMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auction_AssertCoordinator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssertCoordinatorMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServer).AssertCoordinator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auction_AssertCoordinator_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServer).AssertCoordinator(ctx, req.(*AssertCoordinatorMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// Auction_ServiceDesc is the grpc.ServiceDesc for Auction service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Auction_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auction.Auction",
	HandlerType: (*AuctionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Bid",
			Handler:    _Auction_Bid_Handler,
		},
		{
			MethodName: "Result",
			Handler:    _Auction_Result_Handler,
		},
		{
			MethodName: "CallElection",
			Handler:    _Auction_CallElection_Handler,
		},
		{
			MethodName: "AssertCoordinator",
			Handler:    _Auction_AssertCoordinator_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auction.proto",
}

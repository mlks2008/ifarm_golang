// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.0
// source: sdks/gts_shop/pb/gts_shop.proto

package pb

import (
	vo "components/sdks/gts_shop/pb/vo"
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
	GtsShop_AuthCode_FullMethodName             = "/components.sdks.gts_shop.pb.GtsShop/AuthCode"
	GtsShop_SendVerificationCode_FullMethodName = "/components.sdks.gts_shop.pb.GtsShop/SendVerificationCode"
	GtsShop_VerifyCode_FullMethodName           = "/components.sdks.gts_shop.pb.GtsShop/VerifyCode"
	GtsShop_UserProfile_FullMethodName          = "/components.sdks.gts_shop.pb.GtsShop/UserProfile"
	GtsShop_ListAsset_FullMethodName            = "/components.sdks.gts_shop.pb.GtsShop/ListAsset"
	GtsShop_AddAsset_FullMethodName             = "/components.sdks.gts_shop.pb.GtsShop/AddAsset"
	GtsShop_FreezeAsset_FullMethodName          = "/components.sdks.gts_shop.pb.GtsShop/FreezeAsset"
	GtsShop_SubAsset_FullMethodName             = "/components.sdks.gts_shop.pb.GtsShop/SubAsset"
	GtsShop_ReturnAsset_FullMethodName          = "/components.sdks.gts_shop.pb.GtsShop/ReturnAsset"
)

// GtsShopClient is the client API for GtsShop service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GtsShopClient interface {
	// 鉴权登录
	AuthCode(ctx context.Context, in *vo.AuthorizationCodeRequest, opts ...grpc.CallOption) (*vo.GetAccountResponse, error)
	// 发送验证码
	SendVerificationCode(ctx context.Context, in *vo.SendVerificationCodeRequest, opts ...grpc.CallOption) (*vo.SendVerificationCodeResponse, error)
	// 验证验证码
	VerifyCode(ctx context.Context, in *vo.VerifyCodeRequest, opts ...grpc.CallOption) (*vo.VerifyCodeResponse, error)
	// 获取用户信息
	UserProfile(ctx context.Context, in *vo.UserProfileRequest, opts ...grpc.CallOption) (*vo.GetAccountResponse, error)
	// 账户列表
	ListAsset(ctx context.Context, in *vo.ListAssetRequest, opts ...grpc.CallOption) (*vo.ListAssetResponse, error)
	// 增加资产
	AddAsset(ctx context.Context, in *vo.AddAssetRequest, opts ...grpc.CallOption) (*vo.AddAssetResponse, error)
	// 冻结资产
	FreezeAsset(ctx context.Context, in *vo.FreezeAssetRequest, opts ...grpc.CallOption) (*vo.FreezeAssetResponse, error)
	// 扣除资产
	SubAsset(ctx context.Context, in *vo.SubAssetRequest, opts ...grpc.CallOption) (*vo.SubAssetResponse, error)
	// 解冻资产
	ReturnAsset(ctx context.Context, in *vo.ReturnAssetRequest, opts ...grpc.CallOption) (*vo.ReturnAssetResponse, error)
}

type gtsShopClient struct {
	cc grpc.ClientConnInterface
}

func NewGtsShopClient(cc grpc.ClientConnInterface) GtsShopClient {
	return &gtsShopClient{cc}
}

func (c *gtsShopClient) AuthCode(ctx context.Context, in *vo.AuthorizationCodeRequest, opts ...grpc.CallOption) (*vo.GetAccountResponse, error) {
	out := new(vo.GetAccountResponse)
	err := c.cc.Invoke(ctx, GtsShop_AuthCode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gtsShopClient) SendVerificationCode(ctx context.Context, in *vo.SendVerificationCodeRequest, opts ...grpc.CallOption) (*vo.SendVerificationCodeResponse, error) {
	out := new(vo.SendVerificationCodeResponse)
	err := c.cc.Invoke(ctx, GtsShop_SendVerificationCode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gtsShopClient) VerifyCode(ctx context.Context, in *vo.VerifyCodeRequest, opts ...grpc.CallOption) (*vo.VerifyCodeResponse, error) {
	out := new(vo.VerifyCodeResponse)
	err := c.cc.Invoke(ctx, GtsShop_VerifyCode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gtsShopClient) UserProfile(ctx context.Context, in *vo.UserProfileRequest, opts ...grpc.CallOption) (*vo.GetAccountResponse, error) {
	out := new(vo.GetAccountResponse)
	err := c.cc.Invoke(ctx, GtsShop_UserProfile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gtsShopClient) ListAsset(ctx context.Context, in *vo.ListAssetRequest, opts ...grpc.CallOption) (*vo.ListAssetResponse, error) {
	out := new(vo.ListAssetResponse)
	err := c.cc.Invoke(ctx, GtsShop_ListAsset_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gtsShopClient) AddAsset(ctx context.Context, in *vo.AddAssetRequest, opts ...grpc.CallOption) (*vo.AddAssetResponse, error) {
	out := new(vo.AddAssetResponse)
	err := c.cc.Invoke(ctx, GtsShop_AddAsset_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gtsShopClient) FreezeAsset(ctx context.Context, in *vo.FreezeAssetRequest, opts ...grpc.CallOption) (*vo.FreezeAssetResponse, error) {
	out := new(vo.FreezeAssetResponse)
	err := c.cc.Invoke(ctx, GtsShop_FreezeAsset_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gtsShopClient) SubAsset(ctx context.Context, in *vo.SubAssetRequest, opts ...grpc.CallOption) (*vo.SubAssetResponse, error) {
	out := new(vo.SubAssetResponse)
	err := c.cc.Invoke(ctx, GtsShop_SubAsset_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gtsShopClient) ReturnAsset(ctx context.Context, in *vo.ReturnAssetRequest, opts ...grpc.CallOption) (*vo.ReturnAssetResponse, error) {
	out := new(vo.ReturnAssetResponse)
	err := c.cc.Invoke(ctx, GtsShop_ReturnAsset_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GtsShopServer is the server API for GtsShop service.
// All implementations must embed UnimplementedGtsShopServer
// for forward compatibility
type GtsShopServer interface {
	// 鉴权登录
	AuthCode(context.Context, *vo.AuthorizationCodeRequest) (*vo.GetAccountResponse, error)
	// 发送验证码
	SendVerificationCode(context.Context, *vo.SendVerificationCodeRequest) (*vo.SendVerificationCodeResponse, error)
	// 验证验证码
	VerifyCode(context.Context, *vo.VerifyCodeRequest) (*vo.VerifyCodeResponse, error)
	// 获取用户信息
	UserProfile(context.Context, *vo.UserProfileRequest) (*vo.GetAccountResponse, error)
	// 账户列表
	ListAsset(context.Context, *vo.ListAssetRequest) (*vo.ListAssetResponse, error)
	// 增加资产
	AddAsset(context.Context, *vo.AddAssetRequest) (*vo.AddAssetResponse, error)
	// 冻结资产
	FreezeAsset(context.Context, *vo.FreezeAssetRequest) (*vo.FreezeAssetResponse, error)
	// 扣除资产
	SubAsset(context.Context, *vo.SubAssetRequest) (*vo.SubAssetResponse, error)
	// 解冻资产
	ReturnAsset(context.Context, *vo.ReturnAssetRequest) (*vo.ReturnAssetResponse, error)
	mustEmbedUnimplementedGtsShopServer()
}

// UnimplementedGtsShopServer must be embedded to have forward compatible implementations.
type UnimplementedGtsShopServer struct {
}

func (UnimplementedGtsShopServer) AuthCode(context.Context, *vo.AuthorizationCodeRequest) (*vo.GetAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthCode not implemented")
}
func (UnimplementedGtsShopServer) SendVerificationCode(context.Context, *vo.SendVerificationCodeRequest) (*vo.SendVerificationCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendVerificationCode not implemented")
}
func (UnimplementedGtsShopServer) VerifyCode(context.Context, *vo.VerifyCodeRequest) (*vo.VerifyCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyCode not implemented")
}
func (UnimplementedGtsShopServer) UserProfile(context.Context, *vo.UserProfileRequest) (*vo.GetAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserProfile not implemented")
}
func (UnimplementedGtsShopServer) ListAsset(context.Context, *vo.ListAssetRequest) (*vo.ListAssetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAsset not implemented")
}
func (UnimplementedGtsShopServer) AddAsset(context.Context, *vo.AddAssetRequest) (*vo.AddAssetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAsset not implemented")
}
func (UnimplementedGtsShopServer) FreezeAsset(context.Context, *vo.FreezeAssetRequest) (*vo.FreezeAssetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FreezeAsset not implemented")
}
func (UnimplementedGtsShopServer) SubAsset(context.Context, *vo.SubAssetRequest) (*vo.SubAssetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubAsset not implemented")
}
func (UnimplementedGtsShopServer) ReturnAsset(context.Context, *vo.ReturnAssetRequest) (*vo.ReturnAssetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReturnAsset not implemented")
}
func (UnimplementedGtsShopServer) mustEmbedUnimplementedGtsShopServer() {}

// UnsafeGtsShopServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GtsShopServer will
// result in compilation errors.
type UnsafeGtsShopServer interface {
	mustEmbedUnimplementedGtsShopServer()
}

func RegisterGtsShopServer(s grpc.ServiceRegistrar, srv GtsShopServer) {
	s.RegisterService(&GtsShop_ServiceDesc, srv)
}

func _GtsShop_AuthCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vo.AuthorizationCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GtsShopServer).AuthCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GtsShop_AuthCode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GtsShopServer).AuthCode(ctx, req.(*vo.AuthorizationCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GtsShop_SendVerificationCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vo.SendVerificationCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GtsShopServer).SendVerificationCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GtsShop_SendVerificationCode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GtsShopServer).SendVerificationCode(ctx, req.(*vo.SendVerificationCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GtsShop_VerifyCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vo.VerifyCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GtsShopServer).VerifyCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GtsShop_VerifyCode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GtsShopServer).VerifyCode(ctx, req.(*vo.VerifyCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GtsShop_UserProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vo.UserProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GtsShopServer).UserProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GtsShop_UserProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GtsShopServer).UserProfile(ctx, req.(*vo.UserProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GtsShop_ListAsset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vo.ListAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GtsShopServer).ListAsset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GtsShop_ListAsset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GtsShopServer).ListAsset(ctx, req.(*vo.ListAssetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GtsShop_AddAsset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vo.AddAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GtsShopServer).AddAsset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GtsShop_AddAsset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GtsShopServer).AddAsset(ctx, req.(*vo.AddAssetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GtsShop_FreezeAsset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vo.FreezeAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GtsShopServer).FreezeAsset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GtsShop_FreezeAsset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GtsShopServer).FreezeAsset(ctx, req.(*vo.FreezeAssetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GtsShop_SubAsset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vo.SubAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GtsShopServer).SubAsset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GtsShop_SubAsset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GtsShopServer).SubAsset(ctx, req.(*vo.SubAssetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GtsShop_ReturnAsset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vo.ReturnAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GtsShopServer).ReturnAsset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GtsShop_ReturnAsset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GtsShopServer).ReturnAsset(ctx, req.(*vo.ReturnAssetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GtsShop_ServiceDesc is the grpc.ServiceDesc for GtsShop service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GtsShop_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "components.sdks.gts_shop.pb.GtsShop",
	HandlerType: (*GtsShopServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AuthCode",
			Handler:    _GtsShop_AuthCode_Handler,
		},
		{
			MethodName: "SendVerificationCode",
			Handler:    _GtsShop_SendVerificationCode_Handler,
		},
		{
			MethodName: "VerifyCode",
			Handler:    _GtsShop_VerifyCode_Handler,
		},
		{
			MethodName: "UserProfile",
			Handler:    _GtsShop_UserProfile_Handler,
		},
		{
			MethodName: "ListAsset",
			Handler:    _GtsShop_ListAsset_Handler,
		},
		{
			MethodName: "AddAsset",
			Handler:    _GtsShop_AddAsset_Handler,
		},
		{
			MethodName: "FreezeAsset",
			Handler:    _GtsShop_FreezeAsset_Handler,
		},
		{
			MethodName: "SubAsset",
			Handler:    _GtsShop_SubAsset_Handler,
		},
		{
			MethodName: "ReturnAsset",
			Handler:    _GtsShop_ReturnAsset_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sdks/gts_shop/pb/gts_shop.proto",
}

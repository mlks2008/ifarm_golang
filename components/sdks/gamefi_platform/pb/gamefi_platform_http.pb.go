// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.1
// - protoc             v4.25.0
// source: sdks/gamefi_platform/pb/gamefi_platform.proto

package pb

import (
	vo "components/sdks/gamefi_platform/pb/vo"
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationPlatformDustAddAsset = "/components.sdks.gamefi_platform.pb.PlatformDust/AddAsset"
const OperationPlatformDustAddHero = "/components.sdks.gamefi_platform.pb.PlatformDust/AddHero"
const OperationPlatformDustFreezeAsset = "/components.sdks.gamefi_platform.pb.PlatformDust/FreezeAsset"
const OperationPlatformDustFreezeHero = "/components.sdks.gamefi_platform.pb.PlatformDust/FreezeHero"
const OperationPlatformDustGetUserHero = "/components.sdks.gamefi_platform.pb.PlatformDust/GetUserHero"
const OperationPlatformDustReturnAsset = "/components.sdks.gamefi_platform.pb.PlatformDust/ReturnAsset"
const OperationPlatformDustReturnHero = "/components.sdks.gamefi_platform.pb.PlatformDust/ReturnHero"
const OperationPlatformDustSubAsset = "/components.sdks.gamefi_platform.pb.PlatformDust/SubAsset"
const OperationPlatformDustSubHero = "/components.sdks.gamefi_platform.pb.PlatformDust/SubHero"

type PlatformDustHTTPServer interface {
	// AddAsset 增加资产
	AddAsset(context.Context, *vo.DustAddAssetRequest) (*vo.DustAddAssetResponse, error)
	// AddHero 增加Hero
	AddHero(context.Context, *vo.AddHeroRequest) (*vo.AddHeroResponse, error)
	// FreezeAsset 冻结资产
	FreezeAsset(context.Context, *vo.DustFreezeAssetRequest) (*vo.DustFreezeAssetResponse, error)
	// FreezeHero 冻结Hero
	FreezeHero(context.Context, *vo.FreezeHeroRequest) (*vo.FreezeHeroResponse, error)
	// GetUserHero 查询用户Hero
	GetUserHero(context.Context, *vo.GetUserHeroRequest) (*vo.GetUserHeroResponse, error)
	// ReturnAsset 解冻资产
	ReturnAsset(context.Context, *vo.DustReturnAssetRequest) (*vo.DustReturnAssetResponse, error)
	// ReturnHero 解冻Hero
	ReturnHero(context.Context, *vo.ReturnHeroRequest) (*vo.ReturnHeroResponse, error)
	// SubAsset 扣除资产
	SubAsset(context.Context, *vo.DustSubAssetRequest) (*vo.DustSubAssetResponse, error)
	// SubHero 扣除Hero
	SubHero(context.Context, *vo.SubHeroRequest) (*vo.SubHeroResponse, error)
}

func RegisterPlatformDustHTTPServer(s *http.Server, srv PlatformDustHTTPServer) {
	r := s.Route("/")
	r.POST("/user/galaxy_dust/increase", _PlatformDust_AddAsset0_HTTP_Handler(srv))
	r.POST("/user/galaxy_dust/holding", _PlatformDust_FreezeAsset0_HTTP_Handler(srv))
	r.POST("/user/galaxy_dust/confirm", _PlatformDust_SubAsset0_HTTP_Handler(srv))
	r.POST("/user/galaxy_dust/return", _PlatformDust_ReturnAsset0_HTTP_Handler(srv))
	r.POST("/hero/user/transfer/in", _PlatformDust_AddHero0_HTTP_Handler(srv))
	r.POST("/hero/user/lock", _PlatformDust_FreezeHero0_HTTP_Handler(srv))
	r.POST("/hero/user/transfer/out", _PlatformDust_SubHero0_HTTP_Handler(srv))
	r.POST("/hero/user/unlock", _PlatformDust_ReturnHero0_HTTP_Handler(srv))
	r.GET("/hero/user/{id}", _PlatformDust_GetUserHero0_HTTP_Handler(srv))
}

func _PlatformDust_AddAsset0_HTTP_Handler(srv PlatformDustHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.DustAddAssetRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlatformDustAddAsset)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AddAsset(ctx, req.(*vo.DustAddAssetRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.DustAddAssetResponse)
		return ctx.Result(200, reply)
	}
}

func _PlatformDust_FreezeAsset0_HTTP_Handler(srv PlatformDustHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.DustFreezeAssetRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlatformDustFreezeAsset)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.FreezeAsset(ctx, req.(*vo.DustFreezeAssetRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.DustFreezeAssetResponse)
		return ctx.Result(200, reply)
	}
}

func _PlatformDust_SubAsset0_HTTP_Handler(srv PlatformDustHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.DustSubAssetRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlatformDustSubAsset)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SubAsset(ctx, req.(*vo.DustSubAssetRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.DustSubAssetResponse)
		return ctx.Result(200, reply)
	}
}

func _PlatformDust_ReturnAsset0_HTTP_Handler(srv PlatformDustHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.DustReturnAssetRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlatformDustReturnAsset)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ReturnAsset(ctx, req.(*vo.DustReturnAssetRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.DustReturnAssetResponse)
		return ctx.Result(200, reply)
	}
}

func _PlatformDust_AddHero0_HTTP_Handler(srv PlatformDustHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.AddHeroRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlatformDustAddHero)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AddHero(ctx, req.(*vo.AddHeroRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.AddHeroResponse)
		return ctx.Result(200, reply)
	}
}

func _PlatformDust_FreezeHero0_HTTP_Handler(srv PlatformDustHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.FreezeHeroRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlatformDustFreezeHero)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.FreezeHero(ctx, req.(*vo.FreezeHeroRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.FreezeHeroResponse)
		return ctx.Result(200, reply)
	}
}

func _PlatformDust_SubHero0_HTTP_Handler(srv PlatformDustHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.SubHeroRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlatformDustSubHero)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SubHero(ctx, req.(*vo.SubHeroRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.SubHeroResponse)
		return ctx.Result(200, reply)
	}
}

func _PlatformDust_ReturnHero0_HTTP_Handler(srv PlatformDustHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.ReturnHeroRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlatformDustReturnHero)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ReturnHero(ctx, req.(*vo.ReturnHeroRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.ReturnHeroResponse)
		return ctx.Result(200, reply)
	}
}

func _PlatformDust_GetUserHero0_HTTP_Handler(srv PlatformDustHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.GetUserHeroRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPlatformDustGetUserHero)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetUserHero(ctx, req.(*vo.GetUserHeroRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.GetUserHeroResponse)
		return ctx.Result(200, reply)
	}
}

type PlatformDustHTTPClient interface {
	AddAsset(ctx context.Context, req *vo.DustAddAssetRequest, opts ...http.CallOption) (rsp *vo.DustAddAssetResponse, err error)
	AddHero(ctx context.Context, req *vo.AddHeroRequest, opts ...http.CallOption) (rsp *vo.AddHeroResponse, err error)
	FreezeAsset(ctx context.Context, req *vo.DustFreezeAssetRequest, opts ...http.CallOption) (rsp *vo.DustFreezeAssetResponse, err error)
	FreezeHero(ctx context.Context, req *vo.FreezeHeroRequest, opts ...http.CallOption) (rsp *vo.FreezeHeroResponse, err error)
	GetUserHero(ctx context.Context, req *vo.GetUserHeroRequest, opts ...http.CallOption) (rsp *vo.GetUserHeroResponse, err error)
	ReturnAsset(ctx context.Context, req *vo.DustReturnAssetRequest, opts ...http.CallOption) (rsp *vo.DustReturnAssetResponse, err error)
	ReturnHero(ctx context.Context, req *vo.ReturnHeroRequest, opts ...http.CallOption) (rsp *vo.ReturnHeroResponse, err error)
	SubAsset(ctx context.Context, req *vo.DustSubAssetRequest, opts ...http.CallOption) (rsp *vo.DustSubAssetResponse, err error)
	SubHero(ctx context.Context, req *vo.SubHeroRequest, opts ...http.CallOption) (rsp *vo.SubHeroResponse, err error)
}

type PlatformDustHTTPClientImpl struct {
	cc *http.Client
}

func NewPlatformDustHTTPClient(client *http.Client) PlatformDustHTTPClient {
	return &PlatformDustHTTPClientImpl{client}
}

func (c *PlatformDustHTTPClientImpl) AddAsset(ctx context.Context, in *vo.DustAddAssetRequest, opts ...http.CallOption) (*vo.DustAddAssetResponse, error) {
	var out vo.DustAddAssetResponse
	pattern := "/user/galaxy_dust/increase"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPlatformDustAddAsset))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *PlatformDustHTTPClientImpl) AddHero(ctx context.Context, in *vo.AddHeroRequest, opts ...http.CallOption) (*vo.AddHeroResponse, error) {
	var out vo.AddHeroResponse
	pattern := "/hero/user/transfer/in"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPlatformDustAddHero))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *PlatformDustHTTPClientImpl) FreezeAsset(ctx context.Context, in *vo.DustFreezeAssetRequest, opts ...http.CallOption) (*vo.DustFreezeAssetResponse, error) {
	var out vo.DustFreezeAssetResponse
	pattern := "/user/galaxy_dust/holding"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPlatformDustFreezeAsset))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *PlatformDustHTTPClientImpl) FreezeHero(ctx context.Context, in *vo.FreezeHeroRequest, opts ...http.CallOption) (*vo.FreezeHeroResponse, error) {
	var out vo.FreezeHeroResponse
	pattern := "/hero/user/lock"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPlatformDustFreezeHero))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *PlatformDustHTTPClientImpl) GetUserHero(ctx context.Context, in *vo.GetUserHeroRequest, opts ...http.CallOption) (*vo.GetUserHeroResponse, error) {
	var out vo.GetUserHeroResponse
	pattern := "/hero/user/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationPlatformDustGetUserHero))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *PlatformDustHTTPClientImpl) ReturnAsset(ctx context.Context, in *vo.DustReturnAssetRequest, opts ...http.CallOption) (*vo.DustReturnAssetResponse, error) {
	var out vo.DustReturnAssetResponse
	pattern := "/user/galaxy_dust/return"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPlatformDustReturnAsset))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *PlatformDustHTTPClientImpl) ReturnHero(ctx context.Context, in *vo.ReturnHeroRequest, opts ...http.CallOption) (*vo.ReturnHeroResponse, error) {
	var out vo.ReturnHeroResponse
	pattern := "/hero/user/unlock"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPlatformDustReturnHero))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *PlatformDustHTTPClientImpl) SubAsset(ctx context.Context, in *vo.DustSubAssetRequest, opts ...http.CallOption) (*vo.DustSubAssetResponse, error) {
	var out vo.DustSubAssetResponse
	pattern := "/user/galaxy_dust/confirm"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPlatformDustSubAsset))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *PlatformDustHTTPClientImpl) SubHero(ctx context.Context, in *vo.SubHeroRequest, opts ...http.CallOption) (*vo.SubHeroResponse, error) {
	var out vo.SubHeroResponse
	pattern := "/hero/user/transfer/out"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPlatformDustSubHero))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

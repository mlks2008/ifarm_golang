// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.1
// - protoc             v4.25.0
// source: api/in/v1/equipment.proto

package v1

import (
	context "context"
	vo "gamefi_equipment/api/in/v1/vo"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationEquipmentServiceAddEquipment = "/api.in.v1.EquipmentService/AddEquipment"
const OperationEquipmentServiceAddFightEquipment = "/api.in.v1.EquipmentService/AddFightEquipment"
const OperationEquipmentServiceBatchHeroEquipment = "/api.in.v1.EquipmentService/BatchHeroEquipment"
const OperationEquipmentServiceBreakDownEquipment = "/api.in.v1.EquipmentService/BreakDownEquipment"
const OperationEquipmentServiceClearFightEquipment = "/api.in.v1.EquipmentService/ClearFightEquipment"
const OperationEquipmentServiceListEquipment = "/api.in.v1.EquipmentService/ListEquipment"
const OperationEquipmentServiceListHeroEquipment = "/api.in.v1.EquipmentService/ListHeroEquipment"
const OperationEquipmentServiceUpgradeEquipment = "/api.in.v1.EquipmentService/UpgradeEquipment"

type EquipmentServiceHTTPServer interface {
	// AddEquipment 添加装备(掉落装备)
	AddEquipment(context.Context, *vo.AddEquipmentRequest) (*vo.AddEquipmentResponse, error)
	// AddFightEquipment 穿装备
	AddFightEquipment(context.Context, *vo.AddFightEquipmentRequest) (*vo.AddFightEquipmentResponse, error)
	// BatchHeroEquipment 批量获取英雄总属性
	BatchHeroEquipment(context.Context, *vo.BatchHeroEquipmentRequest) (*vo.BatchHeroEquipmentResponse, error)
	// BreakDownEquipment 分解装备
	BreakDownEquipment(context.Context, *vo.BreakDownEquipmentRequest) (*vo.BreakDownEquipmentResponse, error)
	// ClearFightEquipment 脱装备
	ClearFightEquipment(context.Context, *vo.ClearFightEquipmentRequest) (*vo.ClearFightEquipmentResponse, error)
	// ListEquipment 装备列表
	ListEquipment(context.Context, *vo.ListEquipmentRequest) (*vo.ListEquipmentResponse, error)
	// ListHeroEquipment 英雄总属性
	ListHeroEquipment(context.Context, *vo.ListHeroEquipmentRequest) (*vo.ListHeroEquipmentResponse, error)
	// UpgradeEquipment 强化装备
	UpgradeEquipment(context.Context, *vo.UpgradeEquipmentRequest) (*vo.UpgradeEquipmentResponse, error)
}

func RegisterEquipmentServiceHTTPServer(s *http.Server, srv EquipmentServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/in/user/equipment-add", _EquipmentService_AddEquipment0_HTTP_Handler(srv))
	r.POST("/v1/in/user/equipment-upgrade", _EquipmentService_UpgradeEquipment0_HTTP_Handler(srv))
	r.POST("/v1/in/user/equipment-addfight", _EquipmentService_AddFightEquipment0_HTTP_Handler(srv))
	r.POST("/v1/in/user/equipment-clearfight", _EquipmentService_ClearFightEquipment0_HTTP_Handler(srv))
	r.POST("/v1/in/user/equipment-list", _EquipmentService_ListEquipment0_HTTP_Handler(srv))
	r.POST("/v1/in/user/equipment-heroequipment-list", _EquipmentService_ListHeroEquipment0_HTTP_Handler(srv))
	r.POST("/v1/in/user/equipment-heroequipment-batch", _EquipmentService_BatchHeroEquipment0_HTTP_Handler(srv))
	r.POST("/v1/in/user/equipment-breakdown", _EquipmentService_BreakDownEquipment0_HTTP_Handler(srv))
}

func _EquipmentService_AddEquipment0_HTTP_Handler(srv EquipmentServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.AddEquipmentRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationEquipmentServiceAddEquipment)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AddEquipment(ctx, req.(*vo.AddEquipmentRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.AddEquipmentResponse)
		return ctx.Result(200, reply)
	}
}

func _EquipmentService_UpgradeEquipment0_HTTP_Handler(srv EquipmentServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.UpgradeEquipmentRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationEquipmentServiceUpgradeEquipment)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpgradeEquipment(ctx, req.(*vo.UpgradeEquipmentRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.UpgradeEquipmentResponse)
		return ctx.Result(200, reply)
	}
}

func _EquipmentService_AddFightEquipment0_HTTP_Handler(srv EquipmentServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.AddFightEquipmentRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationEquipmentServiceAddFightEquipment)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AddFightEquipment(ctx, req.(*vo.AddFightEquipmentRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.AddFightEquipmentResponse)
		return ctx.Result(200, reply)
	}
}

func _EquipmentService_ClearFightEquipment0_HTTP_Handler(srv EquipmentServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.ClearFightEquipmentRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationEquipmentServiceClearFightEquipment)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ClearFightEquipment(ctx, req.(*vo.ClearFightEquipmentRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.ClearFightEquipmentResponse)
		return ctx.Result(200, reply)
	}
}

func _EquipmentService_ListEquipment0_HTTP_Handler(srv EquipmentServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.ListEquipmentRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationEquipmentServiceListEquipment)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListEquipment(ctx, req.(*vo.ListEquipmentRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.ListEquipmentResponse)
		return ctx.Result(200, reply)
	}
}

func _EquipmentService_ListHeroEquipment0_HTTP_Handler(srv EquipmentServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.ListHeroEquipmentRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationEquipmentServiceListHeroEquipment)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListHeroEquipment(ctx, req.(*vo.ListHeroEquipmentRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.ListHeroEquipmentResponse)
		return ctx.Result(200, reply)
	}
}

func _EquipmentService_BatchHeroEquipment0_HTTP_Handler(srv EquipmentServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.BatchHeroEquipmentRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationEquipmentServiceBatchHeroEquipment)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.BatchHeroEquipment(ctx, req.(*vo.BatchHeroEquipmentRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.BatchHeroEquipmentResponse)
		return ctx.Result(200, reply)
	}
}

func _EquipmentService_BreakDownEquipment0_HTTP_Handler(srv EquipmentServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in vo.BreakDownEquipmentRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationEquipmentServiceBreakDownEquipment)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.BreakDownEquipment(ctx, req.(*vo.BreakDownEquipmentRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*vo.BreakDownEquipmentResponse)
		return ctx.Result(200, reply)
	}
}

type EquipmentServiceHTTPClient interface {
	AddEquipment(ctx context.Context, req *vo.AddEquipmentRequest, opts ...http.CallOption) (rsp *vo.AddEquipmentResponse, err error)
	AddFightEquipment(ctx context.Context, req *vo.AddFightEquipmentRequest, opts ...http.CallOption) (rsp *vo.AddFightEquipmentResponse, err error)
	BatchHeroEquipment(ctx context.Context, req *vo.BatchHeroEquipmentRequest, opts ...http.CallOption) (rsp *vo.BatchHeroEquipmentResponse, err error)
	BreakDownEquipment(ctx context.Context, req *vo.BreakDownEquipmentRequest, opts ...http.CallOption) (rsp *vo.BreakDownEquipmentResponse, err error)
	ClearFightEquipment(ctx context.Context, req *vo.ClearFightEquipmentRequest, opts ...http.CallOption) (rsp *vo.ClearFightEquipmentResponse, err error)
	ListEquipment(ctx context.Context, req *vo.ListEquipmentRequest, opts ...http.CallOption) (rsp *vo.ListEquipmentResponse, err error)
	ListHeroEquipment(ctx context.Context, req *vo.ListHeroEquipmentRequest, opts ...http.CallOption) (rsp *vo.ListHeroEquipmentResponse, err error)
	UpgradeEquipment(ctx context.Context, req *vo.UpgradeEquipmentRequest, opts ...http.CallOption) (rsp *vo.UpgradeEquipmentResponse, err error)
}

type EquipmentServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewEquipmentServiceHTTPClient(client *http.Client) EquipmentServiceHTTPClient {
	return &EquipmentServiceHTTPClientImpl{client}
}

func (c *EquipmentServiceHTTPClientImpl) AddEquipment(ctx context.Context, in *vo.AddEquipmentRequest, opts ...http.CallOption) (*vo.AddEquipmentResponse, error) {
	var out vo.AddEquipmentResponse
	pattern := "/v1/in/user/equipment-add"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationEquipmentServiceAddEquipment))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *EquipmentServiceHTTPClientImpl) AddFightEquipment(ctx context.Context, in *vo.AddFightEquipmentRequest, opts ...http.CallOption) (*vo.AddFightEquipmentResponse, error) {
	var out vo.AddFightEquipmentResponse
	pattern := "/v1/in/user/equipment-addfight"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationEquipmentServiceAddFightEquipment))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *EquipmentServiceHTTPClientImpl) BatchHeroEquipment(ctx context.Context, in *vo.BatchHeroEquipmentRequest, opts ...http.CallOption) (*vo.BatchHeroEquipmentResponse, error) {
	var out vo.BatchHeroEquipmentResponse
	pattern := "/v1/in/user/equipment-heroequipment-batch"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationEquipmentServiceBatchHeroEquipment))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *EquipmentServiceHTTPClientImpl) BreakDownEquipment(ctx context.Context, in *vo.BreakDownEquipmentRequest, opts ...http.CallOption) (*vo.BreakDownEquipmentResponse, error) {
	var out vo.BreakDownEquipmentResponse
	pattern := "/v1/in/user/equipment-breakdown"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationEquipmentServiceBreakDownEquipment))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *EquipmentServiceHTTPClientImpl) ClearFightEquipment(ctx context.Context, in *vo.ClearFightEquipmentRequest, opts ...http.CallOption) (*vo.ClearFightEquipmentResponse, error) {
	var out vo.ClearFightEquipmentResponse
	pattern := "/v1/in/user/equipment-clearfight"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationEquipmentServiceClearFightEquipment))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *EquipmentServiceHTTPClientImpl) ListEquipment(ctx context.Context, in *vo.ListEquipmentRequest, opts ...http.CallOption) (*vo.ListEquipmentResponse, error) {
	var out vo.ListEquipmentResponse
	pattern := "/v1/in/user/equipment-list"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationEquipmentServiceListEquipment))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *EquipmentServiceHTTPClientImpl) ListHeroEquipment(ctx context.Context, in *vo.ListHeroEquipmentRequest, opts ...http.CallOption) (*vo.ListHeroEquipmentResponse, error) {
	var out vo.ListHeroEquipmentResponse
	pattern := "/v1/in/user/equipment-heroequipment-list"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationEquipmentServiceListHeroEquipment))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *EquipmentServiceHTTPClientImpl) UpgradeEquipment(ctx context.Context, in *vo.UpgradeEquipmentRequest, opts ...http.CallOption) (*vo.UpgradeEquipmentResponse, error) {
	var out vo.UpgradeEquipmentResponse
	pattern := "/v1/in/user/equipment-upgrade"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationEquipmentServiceUpgradeEquipment))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

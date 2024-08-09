package server

import (
	"components/common/global"
	in "gamefi_equipment/api/in/v1"
	"gamefi_equipment/internal/conf"
	"gamefi_equipment/internal/service"
	"github.com/go-kratos/grpc-gateway/v2/protoc-gen-openapiv2/generator"
	"github.com/go-kratos/swagger-api/openapiv2"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, equipmentService *service.EquipmentService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
		http.ResponseEncoder(global.ResponseEncoder),
		http.ErrorEncoder(global.ErrorEncoder),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	in.RegisterEquipmentServiceHTTPServer(srv, equipmentService)

	//h := openapiv2.NewHandler()
	h := openapiv2.NewHandler(openapiv2.WithGeneratorOptions(generator.UseJSONNamesForFields(false), generator.EnumsAsInts(true)))
	//将/q/路由放在最前匹配
	srv.HandlePrefix("/q/", h)

	return srv
}

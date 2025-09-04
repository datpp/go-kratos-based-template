package server

import (
	"encoding/json"

	"github.com/go-kratos/kratos/v2/middleware/metrics"

	v1 "github.com/datpp/go-kratos-based-template/api/healthcheck/v1"
	"github.com/datpp/go-kratos-based-template/internal/conf"
	"github.com/datpp/go-kratos-based-template/internal/service"
	"github.com/datpp/go-kratos-based-template/packages/types"
	"github.com/datpp/go-kratos-based-template/packages/utils"
	"google.golang.org/grpc/encoding"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	http2 "net/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, healthcheck *service.HealthcheckService, m *Metrics, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(),
			metrics.Server(
				metrics.WithSeconds(m.Seconds),
				metrics.WithRequests(m.Requests),
			),
		),
		http.RequestDecoder(RequestDecoder),
		http.ResponseEncoder(func(w http.ResponseWriter, r *http.Request, v interface{}) error {
			response := types.StandardResponse{
				Code: http2.StatusOK,
				Data: v,
			}
			data, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http2.StatusOK)
			w.Write(data)
			return nil
		}),
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
	srv.Handle("/metrics", promhttp.Handler())
	v1.RegisterHealthCheckHTTPServer(srv, healthcheck)
	return srv
}

func RequestDecoder(r *http.Request, v interface{}) error {
	return nil
}

// CodecForRequest get encoding.Codec via http.Request
func CodecForRequest(r *http.Request, name string) (encoding.Codec, bool) {
	for _, accept := range r.Header[name] {
		codec := encoding.GetCodec(utils.ContentSubtype(accept))
		if codec != nil {
			return codec, true
		}
	}
	return encoding.GetCodec("json"), false
}

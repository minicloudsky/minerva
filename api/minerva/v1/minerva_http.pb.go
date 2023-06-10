// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.3.1
// - protoc             v4.24.3
// source: minerva/v1/minerva.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationMineRvaParseSqlType = "/minerva.v1.MineRva/ParseSqlType"

type MineRvaHTTPServer interface {
	ParseSqlType(context.Context, *ParseSqlTypeRequest) (*ParseSqlTypeReply, error)
}

func RegisterMineRvaHTTPServer(s *http.Server, srv MineRvaHTTPServer) {
	r := s.Route("/")
	r.POST("/minerva/auditSqlType", _MineRva_ParseSqlType0_HTTP_Handler(srv))
}

func _MineRva_ParseSqlType0_HTTP_Handler(srv MineRvaHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ParseSqlTypeRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationMineRvaParseSqlType)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ParseSqlType(ctx, req.(*ParseSqlTypeRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ParseSqlTypeReply)
		return ctx.Result(200, reply)
	}
}

type MineRvaHTTPClient interface {
	ParseSqlType(ctx context.Context, req *ParseSqlTypeRequest, opts ...http.CallOption) (rsp *ParseSqlTypeReply, err error)
}

type MineRvaHTTPClientImpl struct {
	cc *http.Client
}

func NewMineRvaHTTPClient(client *http.Client) MineRvaHTTPClient {
	return &MineRvaHTTPClientImpl{client}
}

func (c *MineRvaHTTPClientImpl) ParseSqlType(ctx context.Context, in *ParseSqlTypeRequest, opts ...http.CallOption) (*ParseSqlTypeReply, error) {
	var out ParseSqlTypeReply
	pattern := "/minerva/auditSqlType"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationMineRvaParseSqlType))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

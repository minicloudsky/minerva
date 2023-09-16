package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "minerva/api/minerva/v1"
	"minerva/internal/biz"
)

// MinervaService is a minerva service.
type MinervaService struct {
	v1.UnimplementedMinervaServer

	uc          *biz.MinervaUsecase
	minervaRepo biz.MinervaRepo
}

// NewMinervaService new a minerva service.
func NewMinervaService(uc *biz.MinervaUsecase, minervaRepo biz.MinervaRepo) *MinervaService {
	return &MinervaService{uc: uc, minervaRepo: minervaRepo}
}

func (service *MinervaService) ParseSqlType(ctx context.Context, in *v1.ParseSqlTypeRequest) (out *v1.ParseSqlTypeReply, err error) {

	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "InvalidArgument")
	}
	sqlTypeItems, err := service.minervaRepo.ParseSqlType(ctx, in.Sql)
	types := make([]*v1.ParseSqlTypeReply_SqlCheckResult, 0)
	for _, item := range sqlTypeItems {
		var typeStrings []string
		for _, t := range item.Type {
			typeStrings = append(typeStrings, string(t))
		}
		types = append(types, &v1.ParseSqlTypeReply_SqlCheckResult{
			Sql:     item.Sql,
			SqlType: typeStrings,
			Risk:    string(item.Risk),
		})
	}
	if err != nil {
		return nil, err
	}

	return &v1.ParseSqlTypeReply{
		Code:    int32(codes.OK),
		Message: "OK",
		Data:    types,
	}, nil
}

package service

import (
	"context"
	v1 "minerva/api/minerva/v1"
	"minerva/internal/biz"
)

const SuccessCode = 200
const ErrorCode = 400

// MinervaService is a minerva service.
type MinervaService struct {
	v1.UnimplementedMineRvaServer

	uc *biz.MinervaUsecase
}

// NewMinervaService new a minerva service.
func NewMinervaService(uc *biz.MinervaUsecase) *MinervaService {
	return &MinervaService{uc: uc}
}

func (s *MinervaService) ParseSqlType(ctx context.Context, in *v1.ParseSqlTypeRequest) (out *v1.ParseSqlTypeReply, err error) {
	if err := in.Validate(); err != nil {
		return &v1.ParseSqlTypeReply{
			Code:    ErrorCode,
			Message: err.Error(),
			Data:    &v1.ParseSqlTypeReply_Data{SqlTypes: nil},
		}, err
	}
	sqlTypeItems, err := s.uc.ParseSqlType(ctx, in.Sql)
	types := make([]*v1.ParseSqlTypeReply_Data_SqlTypes, 0)
	for _, item := range sqlTypeItems {
		types = append(types, &v1.ParseSqlTypeReply_Data_SqlTypes{
			Sql:     item.Sql,
			SqlType: item.Type,
		})
	}
	if err != nil {
		return &v1.ParseSqlTypeReply{
			Code:    ErrorCode,
			Message: err.Error(),
			Data:    &v1.ParseSqlTypeReply_Data{SqlTypes: types},
		}, err
	}

	return &v1.ParseSqlTypeReply{
		Code:    SuccessCode,
		Message: "OK",
		Data:    &v1.ParseSqlTypeReply_Data{SqlTypes: types},
	}, nil
}

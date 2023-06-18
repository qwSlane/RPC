package services

import (
	"log"
	"rpc/internal/database"
	"rpc/internal/transport/proto_types"
)

type RecordsService struct {
	Storage storage.Storage
}

func NewRecordsService(storage storage.Storage) *RecordsService {

	return &RecordsService{
		Storage: storage,
	}
}

func (h *RecordsService) Welcome(str string) {
	log.Println("Hello", str)
}

func (h *RecordsService) SetNewRecord(level float64, username string, score float64) *proto_types.RpcResponse {
	err := h.Storage.SetNewRecord(int(level), username, int(score))
	if err != nil {
		return &proto_types.RpcResponse{
			Result: nil,
			Error:  err.Error(),
		}
	}

	return &proto_types.RpcResponse{
		Result: nil,
		Error:  "ok",
	}
}

func (h *RecordsService) GetBestN(count float64, level float64) *proto_types.RpcResponse {
	users, err := h.Storage.GetBestN(int(count), int(level))
	if err != nil {
		return &proto_types.RpcResponse{
			Result: nil,
			Error:  err.Error(),
		}
	}

	return &proto_types.RpcResponse{
		Result: users,
		Error:  "",
	}

}

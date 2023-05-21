package api

import (
	"log"
	"main/storage"
	"main/types"
)

type RPCHandler struct {
	Storage storage.Storage
}

func NewRPCHandler(storage storage.Storage) *RPCHandler {
	return &RPCHandler{
		Storage: storage,
	}
}

func (h *RPCHandler) Welcome(str string) {
	log.Println("Hello", str)
}

func (h *RPCHandler) SetNewRecord(level float64, username string, score float64) *types.RpcResponse {
	err := h.Storage.SetNewRecord(int(level), username, int(score))
	if err != nil {
		return &types.RpcResponse{
			Result: nil,
			Error:  err.Error(),
		}
	}

	return &types.RpcResponse{
		Result: nil,
		Error:  "ok",
	}
}

func (h *RPCHandler) GetBestN(count float64, level float64) *types.RpcResponse {
	users, err := h.Storage.GetBestN(int(count), int(level))
	if err != nil {
		return &types.RpcResponse{
			Result: nil,
			Error:  err.Error(),
		}
	}

	return &types.RpcResponse{
		Result: users,
		Error:  "",
	}

}

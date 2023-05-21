package types

type RpcRequest struct {
	Method string        `json:"method"`
	Args   []interface{} `json:"args"`
	Id     int           `json:"id"`
}

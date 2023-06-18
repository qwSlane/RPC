package proto_types

type RpcResponse struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error"`
	Id     int         `json:"id"`
}

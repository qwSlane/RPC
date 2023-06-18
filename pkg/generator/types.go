package main

type MethodData struct {
	ID         int32
	Name       string
	Params     string
	ResultType string
}

type ServiceData struct {
	Package     string
	ServiceName string
	Methods     []MethodData
}

package transport

import (
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"reflect"
	"rpc/internal/transport/types"
)

type methodHandler func(srv interface{}, args *anypb.Any) (*anypb.Any, error)

type MethodDescription struct {
	MethodId int32
	Handler  interface{}
	Method   methodHandler
}

type ServiceDescription struct {
	ServiceName string
	HandlerType interface{}
	Methods     []MethodDescription
}

type ServiceManager struct {
	Api map[int32]MethodDescription
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		Api: make(map[int32]MethodDescription),
	}

}

func (s *ServiceManager) Invoke(method int32, args *anypb.Any) (*types.Response, error) {

	md := s.Api[method]

	result, err := md.Method(md.Handler, args)
	if err != nil {
		return nil, err
	}

	response := &types.Response{
		Result: result,
		Error:  "",
	}

	return response, nil
}

func (s *ServiceManager) RegisterService(sd *ServiceDescription, cs interface{}) {

	if cs != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(cs)
		if !st.Implements(ht) {
			log.Fatalf("Server.RegisterService found the handler type %v that doesn't satisfy %v", st, ht)
		}
	}

	s.register(sd, cs)
}

func (s *ServiceManager) register(sd *ServiceDescription, cs interface{}) {

	for _, method := range sd.Methods {
		if _, ok := s.Api[method.MethodId]; ok {
			log.Fatalf("Server.Register service %s has ID collision %d", sd.ServiceName, method.MethodId)
		}

		method.Handler = cs
		s.Api[method.MethodId] = method
	}

}

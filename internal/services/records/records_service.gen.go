package records

import (
	"google.golang.org/protobuf/types/known/anypb"
	"rpc/internal/app"
	"rpc/internal/models"
	"rpc/internal/services/records/types"
	"rpc/internal/transport"
)

type RecordsServer interface {
	SetNewRecord(*types.Record) error
	GetBestN(*types.BestLevelCount) (*models.Level, error)
}

func RegisterRecordsService(s app.Server, srv RecordsServer) {
	s.ServiceManager.RegisterService(&Records_ServiceDesc, srv)
}

func _Records_SetNewRecordHandler(src interface{}, args *anypb.Any) (*anypb.Any, error) {

	params := new(types.Record)
	err := args.UnmarshalTo(params)
	if err != nil {
		return nil, err
	}

	err = src.(RecordsServer).SetNewRecord(params)
	if err != nil {
		return nil, err
	}

	return nil, nil

}

func _Records_GetBestN(src interface{}, args *anypb.Any) (*anypb.Any, error) {

	params := new(types.BestLevelCount)
	err := args.UnmarshalTo(params)
	if err != nil {
		return nil, err
	}

	result, err := src.(RecordsServer).GetBestN(params)
	if err != nil {
		return nil, err
	}

	anyResult, err := anypb.New(result)
	if err != nil {
		return nil, err
	}

	return anyResult, nil
}

var Records_ServiceDesc = transport.ServiceDescription{
	ServiceName: "Records",
	HandlerType: (*RecordsServer)(nil),
	Methods: []transport.MethodDescription{
		{
			MethodId: 1,
			Handler:  nil,
			Method:   _Records_SetNewRecordHandler,
		},
		{
			MethodId: 2,
			Handler:  nil,
			Method:   _Records_GetBestN,
		},
	}}

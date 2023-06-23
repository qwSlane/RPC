package records
// Generated code.
// DO NOT EDIT.

import(
	"google.golang.org/protobuf/types/known/anypb"
	"rpc/internal/app"
	"rpc/internal/services/records/types"
	"rpc/internal/transport"
)

func RegisterrecordsService (s app.Server, srv records){
	s.ServiceManager.RegisterService(&records_ServiceDesc, srv)
}
func _SetNewRecord_Handler(src interface{}, args *anypb.Any) (*anypb.Any, error){

	params := new(types.Record)
	err := args.UnmarshalTo(params)
	if err != nil{
		return nil, err
	}
	err = src.(records).SetNewRecord(params)
    if err != nil {
        return nil, err
    }

	return nil, nil
}
func _GetBestN_Handler(src interface{}, args *anypb.Any) (*anypb.Any, error){

	params := new(types.BestLevelCount)
	err := args.UnmarshalTo(params)
	if err != nil{
		return nil, err
	}
	result, err := src.(records).GetBestN(params)
	if err != nil {
		return nil, err
	}

	anyResult, err := anypb.New(result)
	if err != nil {
		return nil, err
	}

	return anyResult, nil
}

var records_ServiceDesc = transport.ServiceDescription{
	ServiceName: "records",
	HandlerType: (*records)(nil),
	Methods: []transport.MethodDescription{
		{
			MethodId: 0,
			Handler:  nil,
			Method:   _SetNewRecord_Handler,
		},
		{
			MethodId: 1,
			Handler:  nil,
			Method:   _GetBestN_Handler,
		},
	},
}
	
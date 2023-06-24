package templates

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"text/template"
)

func Replace(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

func Split(old, symbol string) string {
	return strings.Split(old, symbol)[1]
}

func Title(old string) string {
	t := cases.Title(language.English)
	return t.String(old)
}

var funcMap = template.FuncMap{
	"Replace": Replace,
	"Split":   Split,
	"Title":   Title,
}

var Serv = template.Must(template.New("").Funcs(funcMap).Parse(`package {{ .Package }}
// Generated code.
// DO NOT EDIT.

import(
	"google.golang.org/protobuf/types/known/anypb"
	"rpc/internal/app"
	"rpc/internal/services/{{ .ServiceName }}/types"
	"rpc/internal/transport"
)

func Register{{ .ServiceName }}Service (s app.Server, srv {{ .ServiceName}}){
	s.ServiceManager.RegisterService(&{{ .ServiceName}}_ServiceDesc, srv)
}

{{- range .Methods }}
func _{{ .Name }}_Handler(src interface{}, args *anypb.Any) (*anypb.Any, error){

	params := new({{ Replace .Params "*" ""}})
	err := args.UnmarshalTo(params)
	if err != nil{
		return nil, err
	}


	{{- if .ResultType }}
	result, err := src.({{ $.ServiceName }}).{{ .Name }}(params)
	if err != nil {
		return nil, err
	}

	anyResult, err := anypb.New(result)
	if err != nil {
		return nil, err
	}

	return anyResult, nil
	{{- else }}
	err = src.({{ $.ServiceName }}).{{ .Name }}(params)
    if err != nil {
        return nil, err
    }

	return nil, nil
	{{- end }}
}
{{- end }}

var {{ .ServiceName }}_ServiceDesc = transport.ServiceDescription{
	ServiceName: "{{ .ServiceName }}",
	HandlerType: (*{{ .ServiceName }})(nil),
	Methods: []transport.MethodDescription{
		{{- range .Methods }}
		{
			MethodId: {{ .ID }},
			Handler:  nil,
			Method:   _{{ .Name }}_Handler,
		},
		{{- end }}
	},
}
	`))

var Client = template.Must(template.New("").Funcs(funcMap).Parse(`
// Generated code.
// DO NOT EDIT.

using Client;
using Google.Protobuf.WellKnownTypes;
using Transport;
using Types;

public class {{ Title .ServiceName }}Client
{
{{- range .Methods }}

   {{- if .ResultType }}
   public async Task<{{ Split .ResultType "." }}> {{ .Name }}({{ Replace .Params "*types." ""}} args)
   {
      Request request = new Request
      {
         Args = Any.Pack(args),
         Method = 0,
      };
      
      var result = await RpcClient.Instance.Invoke(request);

      if (String.IsNullOrEmpty(result.Error) == false)
      {
         throw new Exception(result.Error);
      }

      return result.Result.Unpack<{{ Split .ResultType "."}}>();
   }
   {{- else }}
   public async Task {{ .Name }}({{ Replace .Params "*types." ""}} args)
   {
      Request request = new Request
      {
         Args = Any.Pack(args),
         Method = 0,
      };
      
      var result = await RpcClient.Instance.Invoke(request);

      if (String.IsNullOrEmpty(result.Error) == false)
      {
         throw new Exception(result.Error);
      }
   }
   {{- end }}

{{- end }}
}

	`))

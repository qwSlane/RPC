package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"html/template"
	"log"
	"os"
	"strings"
)

type ServiceGenerator struct {
	TypeSpec      *ast.TypeSpec
	InterfaceType *ast.InterfaceType
}

func expr2string(expr ast.Expr) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), expr)
	if err != nil {
		log.Fatalf("error print expression to string: #{err}")
	}
	return buf.String()
}

func (g ServiceGenerator) Generate(md *ast.File, prevId int32) {

	var methods []MethodData

	for i, method := range g.InterfaceType.Methods.List {

		funcType := method.Type.(*ast.FuncType)

		result := expr2string(funcType.Results.List[0].Type)
		if expr2string(funcType.Results.List[0].Type) == "error" {
			result = ""
		}

		methods = append(methods, MethodData{
			ID:         prevId + int32(i),
			Name:       method.Names[0].Name,
			Params:     expr2string(funcType.Params.List[0].Type),
			ResultType: result,
		})
	}

	params := &ServiceData{
		Package:     md.Name.Name,
		ServiceName: g.TypeSpec.Name.Name,
		Methods:     methods,
	}

	buf := new(bytes.Buffer)

	err := tmpl.Execute(buf, params)
	if err != nil {
		panic(err)
	}

	filename := strings.ToLower(g.TypeSpec.Name.Name) + "_service_gen.go"

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = buf.WriteTo(file)
	if err != nil {
		panic(err)
	}

	fmt.Println("File created successfully!")
}

var tmpl = template.Must(template.New("").Parse(`package {{ .Package }}
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

	params := new({{ .Params }})
	err := args.UnmarshalTo(*params)
	if err != nil{
		return nil, err
	}


	{{- if .ResultType }}
	result, err := src.({{ $.ServiceName }}).{{ .Name }}(*params)
	if err != nil {
		return nil, err
	}

	anyResult, err := anypb.New(result)
	if err != nil {
		return nil, err
	}

	return anyResult, nil
	{{- else }}
	err = src.({{ $.ServiceName }}).{{ .Name }}(*params)
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

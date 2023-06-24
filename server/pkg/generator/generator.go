package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os"
	"rpc/pkg/generator/templates"
	"strconv"
	"strings"
)

const FILENAME = "../../../pkg/generator/data.dat"

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

func (g ServiceGenerator) Generate(md *ast.File) {

	var methods []MethodData

	prevId := GetLastIndex()

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

	UpdateIndex(prevId + int32(len(g.InterfaceType.Methods.List)))

	params := &ServiceData{
		Package:     md.Name.Name,
		ServiceName: g.TypeSpec.Name.Name,
		Methods:     methods,
	}

	GenerateServer(params, g.TypeSpec.Name.Name)
	GenerateClient(params, g.TypeSpec.Name.Name)

}

func UpdateIndex(i int32) {

	data := []byte(fmt.Sprintf("%d", i))
	err := os.WriteFile(FILENAME, data, 0666)
	if err != nil {
		fmt.Println("Error writing int in file:", err)
	}
}

func GetLastIndex() int32 {

	var prev int
	file, err := os.Open(FILENAME)
	if err == nil {
		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			prev, err = strconv.Atoi(scanner.Text())
		}

		if err != nil {
			log.Fatal("Error reading int from file")
		}
		file.Close()

		return int32(prev)
	} else {
		file, err = os.Create(FILENAME)
		if err != nil {
			log.Fatal("Error creating file")
		}
		defer file.Close()
		return 0
	}

}

func GenerateServer(params *ServiceData, name string) {
	buf := new(bytes.Buffer)

	err := templates.Serv.Execute(buf, params)
	if err != nil {
		panic(err)
	}

	filename := strings.ToLower(name) + "_service_gen.go"

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = buf.WriteTo(file)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server file created successfully!")
}

func GenerateClient(params *ServiceData, name string) {
	buf := new(bytes.Buffer)

	err := templates.Client.Execute(buf, params)
	if err != nil {
		panic(err)
	}

	t := cases.Title(language.English)
	filename := t.String(name) + "Client.gen.cs"

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = buf.WriteTo(file)
	if err != nil {
		panic(err)
	}

	fmt.Println("Client file created successfully!")
}

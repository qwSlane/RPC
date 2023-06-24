package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/ast/inspector"
)

func main() {

	path := os.Getenv("GOFILE")
	if path == "" {
		log.Fatal("GOFILE must be set")
	}

	astInFile, err := parser.ParseFile(
		token.NewFileSet(),
		path,
		nil,
		parser.ParseComments,
	)
	if err != nil {
		log.Fatalf("parse file: %v", err)
	}

	i := inspector.New([]*ast.File{astInFile})

	iFilter := []ast.Node{

		&ast.GenDecl{},
	}

	var genTasks []ServiceGenerator

	i.Nodes(iFilter, func(node ast.Node, push bool) (proceed bool) {

		genDecl := node.(*ast.GenDecl)
		if genDecl.Doc == nil {
			return false
		}

		typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			return false
		}

		interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			return false
		}

		for _, comment := range genDecl.Doc.List {
			switch comment.Text {
			case "//generator:gen":
				genTasks = append(genTasks, ServiceGenerator{
					TypeSpec:      typeSpec,
					InterfaceType: interfaceType,
				})
			}
		}
		return false
	})

	astOutFile := &ast.File{
		Name: astInFile.Name,
	}

	for _, task := range genTasks {
		task.Generate(astOutFile)
	}
}

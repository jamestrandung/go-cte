package main

import (
	"context"
	"fmt"
	"github.com/jamestrandung/go-cte/sample/dto"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/endpoint"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/jamestrandung/go-cte/sample/config"
	"github.com/jamestrandung/go-cte/sample/server"
	"golang.org/x/tools/go/packages"
)

func main() {
	err := config.Engine.VerifyConfigurations()
	fmt.Println("Engine configuration error:", err)

	server.Serve()

	testEngine()
	testPlainGo()
}

func testEngine() {
	p := endpoint.NewPlan(
		dto.CostRequest{
			PointA: "Clementi",
			PointB: "Changi Airport",
		},
		server.Dependencies,
	)

	if err := p.Execute(context.Background()); err != nil {
		fmt.Println(err)
	}

	//config.Print(p.GetTravelCost())
	config.Print(p.GetTotalCost())
	config.Print(p.GetVATAmount())
}

func testPlainGo() {
	quote, err := server.Handler.Handle(context.Background(), dto.CostRequest{
		PointA: "Clementi",
		PointB: "Changi Airport",
	})

	if err != nil {
		fmt.Println("Plain Go error:", err)
	}

	config.Print(quote.TotalCost)
	config.Print(quote.VATAmount)
}

const loadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedDeps |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo

func testParsePackage() {
	loadConfig := new(packages.Config)
	loadConfig.Mode = loadMode
	loadConfig.Fset = token.NewFileSet()
	pkgs, err := packages.Load(loadConfig, "github.com/jamestrandung/go-cte/sample/...")
	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		for _, syn := range pkg.Syntax {
			for _, dec := range syn.Decls {
				if gen, ok := dec.(*ast.GenDecl); ok && gen.Tok == token.TYPE {
					// print doc comment of the type
					// fmt.Println(gen.Doc.List[0])
					for _, spec := range gen.Specs {
						if ts, ok := spec.(*ast.TypeSpec); ok {
							obj, ok := pkg.TypesInfo.Defs[ts.Name]
							if !ok {
								continue
							}

							typeName, ok := obj.(*types.TypeName)
							if !ok {
								continue
							}

							named, ok := typeName.Type().(*types.Named)
							if !ok {
								continue
							}

							// print the full name of the type
							fmt.Println(named)
							fmt.Println(pkg.TypesInfo.Types[ts.Type].Type)

							s, ok := named.Underlying().(*types.Struct)
							if !ok {
								continue
							}

							// print the struct's fields and tags
							for i := 0; i < s.NumFields(); i++ {
								idx := fmt.Sprint(i)
								fmt.Println("s.Field(", idx, ").Name(): ", s.Field(i).Name())
								fmt.Println("s.Tag(", idx, "): ", s.Tag(i))
							}
						}
					}
				}
			}
		}
	}

	// pkg, err := importer.Default().Import("github.com/jamestrandung/go-cte/core")
	// if err != nil {
	// 	fmt.Printf("error: %s\n", err.Error())
	// 	return
	// }
	// for _, declName := range pkg.Scope().Names() {
	// 	fmt.Println(declName)
	// }
}

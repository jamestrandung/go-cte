package main

import (
	"context"
	"fmt"
	"github.com/jamestrandung/go-cte/sample/service/components/platformfee"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/jamestrandung/go-cte/sample/dto"

	"github.com/jamestrandung/go-cte/sample/config"
	"github.com/jamestrandung/go-cte/sample/server"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/parallel"
	"github.com/jamestrandung/go-cte/sample/service/scaffolding/sequential"
	"golang.org/x/tools/go/packages"
)

type customPostHook struct{}

func (customPostHook) PostExecute(p any) error {
	config.Print("After sequential plan custom hook")

	return nil
}

type prePlatformFeeHook struct{}

func (prePlatformFeeHook) PreExecute(p any) error {
	config.Print("Before platform fee computer custom hook")

	return nil
}

type postPlatformFeeHook struct{}

func (postPlatformFeeHook) PostExecute(p any) error {
	config.Print("After platform fee computer custom hook")

	return nil
}

func main() {
	testEngine()
}

const loadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedDeps |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo

func testParsePackage2() {
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

func testEngine() {
	server.Serve()

	config.Engine.ConnectPostHook(sequential.SequentialPlan{}, customPostHook{})
	config.Engine.ConnectPreHook(platformfee.PlatformFee{}, prePlatformFeeHook{})
	config.Engine.ConnectPostHook(platformfee.PlatformFee{}, postPlatformFeeHook{})

	p := parallel.NewPlan(
		dto.CostRequest{
			PointA: "Clementi",
			PointB: "Changi Airport",
		},
		server.Dependencies,
	)

	if err := p.Execute(context.Background()); err != nil {
		fmt.Println(err)
	}

	config.Print(p.GetTravelCost())
	config.Print(p.GetTotalCost())
	config.Print(p.GetVATAmount())
}

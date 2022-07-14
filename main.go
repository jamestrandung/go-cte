// Copyright (c) 2012-2022 Grabtaxi Holdings PTE LTD (GRAB), All Rights Reserved. NOTICE: All information contained herein
// is, and remains the property of GRAB. The intellectual and technical concepts contained herein are confidential, proprietary
// and controlled by GRAB and may be covered by patents, patents in process, and are protected by trade secret or copyright law.
//
// You are strictly forbidden to copy, download, store (in any medium), transmit, disseminate, adapt or change this material
// in any way unless prior written permission is obtained from GRAB. Access to the source code contained herein is hereby
// forbidden to anyone except current GRAB employees or contractors with binding Confidentiality and Non-disclosure agreements
// explicitly covering such access.
//
// The copyright notice above does not evidence any actual or intended publication or disclosure of this source code,
// which includes information that is confidential and/or proprietary, and is a trade secret, of GRAB.
//
// ANY REPRODUCTION, MODIFICATION, DISTRIBUTION, PUBLIC PERFORMANCE, OR PUBLIC DISPLAY OF OR THROUGH USE OF THIS SOURCE
// CODE WITHOUT THE EXPRESS WRITTEN CONSENT OF GRAB IS STRICTLY PROHIBITED, AND IN VIOLATION OF APPLICABLE LAWS AND
// INTERNATIONAL TREATIES. THE RECEIPT OR POSSESSION OF THIS SOURCE CODE AND/OR RELATED INFORMATION DOES NOT CONVEY
// OR IMPLY ANY RIGHTS TO REPRODUCE, DISCLOSE OR DISTRIBUTE ITS CONTENTS, OR TO MANUFACTURE, USE, OR SELL ANYTHING
// THAT IT MAY DESCRIBE, IN WHOLE OR IN PART.

package main

import (
	"context"
	"fmt"
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

type parent struct {
	child
}

type child struct {
	field float64
}

func (c *child) set(num float64) {
	c.field = num
}

func (c *child) get() float64 {
	return c.field
}

type in interface {
	set(num float64)
}

func doSet(i in) {
	i.set(3)
}

func main() {
	// p := parent{}
	// fmt.Println(p.get())
	// p.set(2)
	// fmt.Println(p.get())
	// doSet(p)
	// fmt.Println(p.get())

	// method, ok := reflect.ValueOf(dummy{}).Type().MethodByName("Do")
	// if ok {
	// 	fmt.Println(method)
	// }
	//
	// method.Func.Call([]reflect.Value{reflect.ValueOf(dummy{})})

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

	config.Engine.ConnectPostHook(&sequential.SequentialPlan{}, customPostHook{})

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

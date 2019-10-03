// MIT License
//
// Copyright (c) 2019 Oncilla
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package serrorscheck

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer checks all calls on the serrors package.
var Analyzer = &analysis.Analyzer{
	Name:             "serrorslint",
	Doc:              "reports invalid serrors calls",
	Run:              run,
	RunDespiteErrors: true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		tgtPkg := findPkgName(file)
		if tgtPkg == "" {
			continue
		}
		ast.Inspect(file, func(n ast.Node) bool {
			ce, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			se, ok := ce.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			pkg, ok := se.X.(*ast.Ident)
			if !ok || pkg.Name != tgtPkg {
				return true
			}
			var varargs []ast.Expr
			switch se.Sel.Name {
			case "New":
				if len(ce.Args) < 2 {
					return true
				}
				varargs = ce.Args[1:]
			case "WithCtx":
				if len(ce.Args) < 2 {
					pass.Reportf(ce.Pos(), "should have context: expr=%q", render(pass.Fset, ce))
					return true
				}
				varargs = ce.Args[1:]
			case "Wrap", "WrapStr":
				if len(ce.Args) < 3 {
					return true
				}
				varargs = ce.Args[2:]
			}
			// We cannot check if varargs with ellipsis.
			if ce.Ellipsis != token.NoPos {
				return true
			}
			if len(varargs)%2 != 0 {
				pass.Reportf(varargs[0].Pos(), "context should be even: len=%d ctx=%s",
					len(varargs), renderCtx(pass.Fset, varargs))
			}
			for i := 0; i < len(varargs); i += 2 {
				lit := varargs[i]
				if !isString(pass, lit) {
					pass.Reportf(lit.Pos(), "key should be string: type=%q name=%q",
						pass.TypesInfo.TypeOf(lit), render(pass.Fset, lit))
				}
			}
			return true
		})
	}
	return nil, nil
}

func findPkgName(file *ast.File) string {
	var tgtPkg string
	for _, imp := range file.Imports {
		if imp.Path.Value == `"github.com/scionproto/scion/go/lib/serrors"` {
			tgtPkg = "serrors"
			if imp.Name != nil {
				tgtPkg = imp.Name.Name
			}
		}
	}
	return tgtPkg
}

func isString(pass *analysis.Pass, lit ast.Expr) bool {
	t, ok := pass.TypesInfo.TypeOf(lit).Underlying().(*types.Basic)
	return ok && t.Info()&types.IsString != 0
}

func renderCtx(fset *token.FileSet, varargs []ast.Expr) string {
	var p []string
	for _, arg := range varargs {
		p = append(p, render(fset, arg))
	}
	return fmt.Sprintf("[%s]", strings.Join(p, ","))
}

func render(fset *token.FileSet, x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}

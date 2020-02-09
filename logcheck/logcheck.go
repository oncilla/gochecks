// MIT License
//
// Copyright (c) 2020 Oncilla
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

package logcheck

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

// Analyzer checks all calls on the log package.
var Analyzer = &analysis.Analyzer{
	Name:             "logcheck",
	Doc:              "reports invalid log calls",
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
			if !isTarget(se, tgtPkg) {
				return true
			}
			var varargs []ast.Expr
			switch se.Sel.Name {
			case "Trace", "Debug", "Info", "Warn", "Error", "Crit":
				if len(ce.Args) < 2 {
					return true
				}
				varargs = ce.Args[1:]
			}
			// We cannot check if varargs with ellipsis.
			if ce.Ellipsis != token.NoPos {
				return true
			}
			if len(varargs)%2 != 0 {
				pass.Reportf(varargs[0].Pos(), "context should be even: len=%d ctx=%s expr=%q",
					len(varargs), renderCtx(pass.Fset, varargs), render(pass.Fset, ce))
			}
			for i := 0; i < len(varargs); i += 2 {
				lit := varargs[i]
				if !isString(pass, lit) {
					pass.Reportf(lit.Pos(), "key should be string: type=%q name=%q expr=%q",
						pass.TypesInfo.TypeOf(lit), render(pass.Fset, lit), render(pass.Fset, ce))
				}
			}
			return true
		})
	}
	return nil, nil
}

func isTarget(se *ast.SelectorExpr, tgtPkg string) bool {
	switch x := se.X.(type) {
	case *ast.Ident:
		if x.Name == tgtPkg && x.Obj == nil {
			return true
		}
		if x.Obj == nil {
			return false
		}
		decl, ok := x.Obj.Decl.(*ast.AssignStmt)
		if !ok {
			return false
		}
		for _, data := range decl.Rhs {
			if loggerConstructor(data, tgtPkg) {
				return true
			}
		}
		return false
	case *ast.CallExpr:
		return loggerConstructor(x, tgtPkg)
	default:
		return false
	}
}

func loggerConstructor(exp ast.Expr, tgtPkg string) bool {
	ce, ok := exp.(*ast.CallExpr)
	if !ok {
		return false
	}
	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	id, ok := se.X.(*ast.Ident)
	if !ok {
		return false
	}
	if id.Name != tgtPkg {
		return false
	}
	switch se.Sel.Name {
	case "FromCtx", "New", "Root":
		return true
	}
	return false
}

func findPkgName(file *ast.File) string {
	var tgtPkg string
	for _, imp := range file.Imports {
		if imp.Path.Value == `"github.com/scionproto/scion/go/lib/log"` {
			tgtPkg = "log"
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

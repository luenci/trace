package ast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
)

type instrumenter struct {
	traceImport string
	tracePkg    string
	traceFunc   string
}

func NewInstrument(traceImport, tracePkg, traceFunc string) *instrumenter {
	return &instrumenter{
		traceImport: traceImport,
		tracePkg:    tracePkg,
		traceFunc:   traceFunc,
	}
}

func hasFuncDecl(f *ast.File) bool {
	if len(f.Decls) == 0 {
		return false
	}

	for _, decl := range f.Decls {
		_, ok := decl.(*ast.FuncDecl)
		if ok {
			return true
		}
	}

	return false
}

func (i *instrumenter) Instrument(filename string) ([]byte, error) {
	fset := token.NewFileSet()
	curAST, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", filename, err)
	}
	if !hasFuncDecl(curAST) { // 如果整个源码都不包含函数声明，则无需注入操作，直接返回
		return nil, nil
	}
	// 在AST上添加包导入语句
	astutil.AddImport(fset, curAST, i.traceImport)
	// 向AST上的所有函数注入Trace函数
	i.addDeferTraceIntoFuncDecls(curAST)
	buf := &bytes.Buffer{}
	err = format.Node(buf, fset, curAST) // 将修改后的AST转换回Go源码
	if err != nil {
		return nil, fmt.Errorf("error formatting new code: %w", err)
	}
	return buf.Bytes(), nil // 返回转换后的Go源码
}

func (i *instrumenter) addDeferTraceIntoFuncDecls(f *ast.File) {
	for _, decl := range f.Decls { // 遍历所有声明语句
		fd, ok := decl.(*ast.FuncDecl) // 类型断言：是否为函数声明
		if ok {
			// 如果是函数声明，则注入跟踪设施
			i.addDeferStmt(fd)
		}
	}
}

func (i *instrumenter) addDeferStmt(fd *ast.FuncDecl) (added bool) {
	stmts := fd.Body.List
	// 判断"defer trace.Trace()()"语句是否已经存在
	for _, stmt := range stmts {
		ds, ok := stmt.(*ast.DeferStmt)
		if !ok {
			// 如果不是defer语句，则继续for循环
			continue
		}
		// 如果是defer语句，则要进一步判断是否是 defer trace.Trace()()
		ce, ok := ds.Call.Fun.(*ast.CallExpr)
		if !ok {
			continue
		}
		se, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		x, ok := se.X.(*ast.Ident)
		if !ok {
			continue
		}
		if (x.Name == i.tracePkg) && (se.Sel.Name == i.traceFunc) {
			// defer trace.Trace()()已存在，返回
			return false
		}
	}
	// 没有找到"defer trace.Trace()()"，注入一个新的跟踪语句
	// 在AST上构造一个defer trace.Trace()()
	ds := &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: i.tracePkg,
					},
					Sel: &ast.Ident{
						Name: i.traceFunc,
					},
				},
			},
		},
	}
	newList := make([]ast.Stmt, len(stmts)+1)
	copy(newList[1:], stmts)
	newList[0] = ds // 注入新构造的defer语句
	fd.Body.List = newList
	return true
}

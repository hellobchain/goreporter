// Copyright 2013 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Gocyclo calculates the cyclomatic complexities of functions and
// methods in Go source code.
//
// Usage:
//      gocyclo [<flag> ...] <Go file or directory> ...
//
// Flags
//      -over N   show functions with complexity > N only and
//                return exit code 1 if the output is non-empty
//      -top N    show the top N most complex functions only
//      -avg      show the average complexity
//
// The output fields for each line are:
// <complexity> <package> <function> <file:row:column>
package cyclo

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hellobchain/wswlog/wlogging"
)

var logger = wlogging.MustGetLoggerWithoutName()

const usageDoc = `Calculate cyclomatic complexities of Go functions.
usage:
        gocyclo [<flag> ...] <Go file or directory> ...

Flags
        -over N   show functions with complexity > N only and
                  return exit code 1 if the set is non-empty
        -top N    show the top N most complex functions only
        -avg      show the average complexity over all functions,
                  not depending on whether -over or -top are set

The output fields for each line are:
<complexity> <package> <function> <file:row:column>
`

func usage() {
	fmt.Fprintf(os.Stderr, usageDoc)
	os.Exit(2)
}

var (
	over = 0
	top  = 10
	avg  = true
)

func Cyclo(packagePath, except string) ([]string, string) {
	args := []string{packagePath}
	if len(args) == 0 {
		usage()
	}

	stats := analyze(args, except)
	sort.Sort(byComplexity(stats))
	// written := writeStats(os.Stdout, stats)
	packageAvg := "0"
	if avg {
		packageAvg = getAverage(stats)
	}
	result := make([]string, 0)

	if over > 0 {
		return result, packageAvg
	}

	for _, stat := range stats {
		result = append(result, stat.String())
	}

	return result, packageAvg
}

func analyze(paths []string, except string) []stat {
	stats := make([]stat, 0)
	for _, path := range paths {
		if isDir(path) && !checkExcept(path, except) {
			stats = analyzeDir(path, stats)
		} else if !checkExcept(path, except) {
			stats = analyzeFile(path, stats)
		}
	}
	return stats
}

func checkExcept(path, except string) bool {
	if except == "" || except == " " {
		return false
	}
	excepts := strings.Split(except, ",")
	for _, val := range excepts {
		if val != "" && val != " " {
			if strings.Contains(path, val) {
				return true
			}
		}
	}
	return false
}

func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

func analyzeFile(fname string, stats []stat) []stat {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, 0)
	if err != nil {
		exitError(err)
	}
	return buildStats(f, fset, stats)
}

func analyzeDir(dirname string, stats []stat) []stat {
	files, _ := filepath.Glob(filepath.Join(dirname, "*.go"))
	for _, file := range files {
		stats = analyzeFile(file, stats)
	}
	return stats
}

func exitError(err error) {
	fmt.Fprintln(os.Stderr, err)
	// os.Exit(1)
}

func writeStats(w io.Writer, sortedStats []stat) int {
	for i, stat := range sortedStats {
		if i == top {
			return i
		}
		if stat.Complexity <= over {
			return i
		}
		fmt.Fprintln(w, stat)
	}
	return len(sortedStats)
}

func showAverage(stats []stat) {
	logger.Infof("Average: %.3g\n", average(stats))
}

func getAverage(stats []stat) string {
	return fmt.Sprintf("%.2f", average(stats))
}

func average(stats []stat) float64 {
	total := 0
	for _, s := range stats {
		total += s.Complexity
	}
	return float64(total) / float64(len(stats))
}

type stat struct {
	PkgName    string
	FuncName   string
	Complexity int
	Pos        token.Position
}

func (s stat) String() string {
	return fmt.Sprintf("%d %s %s %s", s.Complexity, s.PkgName, s.FuncName, s.Pos)
}

type byComplexity []stat

func (s byComplexity) Len() int      { return len(s) }
func (s byComplexity) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byComplexity) Less(i, j int) bool {
	return s[i].Complexity >= s[j].Complexity
}

func buildStats(f *ast.File, fset *token.FileSet, stats []stat) []stat {
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			stats = append(stats, stat{
				PkgName:    f.Name.Name,
				FuncName:   funcName(fn),
				Complexity: complexity(fn),
				Pos:        fset.Position(fn.Pos()),
			})
		}
	}
	return stats
}

// funcName returns the name representation of a function or method:
// "(Type).Name" for methods or simply "Name" for functions.
func funcName(fn *ast.FuncDecl) string {
	if fn.Recv != nil {
		if fn.Recv.NumFields() > 0 {
			typ := fn.Recv.List[0].Type
			return fmt.Sprintf("(%s).%s", recvString(typ), fn.Name)
		}
	}
	return fn.Name.Name
}

// recvString returns a string representation of recv of the
// form "T", "*T", or "BADRECV" (if not a proper receiver type).
func recvString(recv ast.Expr) string {
	switch t := recv.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + recvString(t.X)
	}
	return "BADRECV"
}

// complexity calculates the cyclomatic complexity of a function.
func complexity(fn *ast.FuncDecl) int {
	v := complexityVisitor{}
	ast.Walk(&v, fn)
	return v.Complexity
}

type complexityVisitor struct {
	// Complexity is the cyclomatic complexity
	Complexity int
}

// Visit implements the ast.Visitor interface.
func (v *complexityVisitor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.FuncDecl, *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause:
		v.Complexity++
	case *ast.BinaryExpr:
		if n.Op == token.LAND || n.Op == token.LOR {
			v.Complexity++
		}
	}
	return v
}

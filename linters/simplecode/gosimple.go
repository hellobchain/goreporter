package simplecode // import "github.com/hellobchain/goreporter/linters/simplecode"

import (
	"github.com/hellobchain/goreporter/linters/simplecode/lint/lintutil"
	"github.com/hellobchain/goreporter/linters/simplecode/simple"
)

func Simple(path map[string]string, except string) []string {
	var res []string
	for _, p := range path {
		res = append(res, lintutil.ProcessArgs(except, "gosimple", simple.Funcs, []string{p})...)
	}
	return res
}

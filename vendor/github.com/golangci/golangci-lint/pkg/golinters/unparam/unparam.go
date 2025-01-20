package unparam

import (
<<<<<<< HEAD
=======
	"sync"

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/packages"
	"mvdan.cc/unparam/check"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
<<<<<<< HEAD
=======
	"github.com/golangci/golangci-lint/pkg/result"
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
)

const linterName = "unparam"

func New(settings *config.UnparamSettings) *goanalysis.Linter {
<<<<<<< HEAD
=======
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	analyzer := &analysis.Analyzer{
		Name:     linterName,
		Doc:      goanalysis.TheOnlyanalyzerDoc,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
		Run: func(pass *analysis.Pass) (any, error) {
<<<<<<< HEAD
			err := runUnparam(pass, settings)
=======
			issues, err := runUnparam(pass, settings)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			if err != nil {
				return nil, err
			}

<<<<<<< HEAD
=======
			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		linterName,
		"Reports unused function parameters",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		if settings.Algo != "cha" {
			lintCtx.Log.Warnf("`linters-settings.unparam.algo` isn't supported by the newest `unparam`")
		}
<<<<<<< HEAD
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runUnparam(pass *analysis.Pass, settings *config.UnparamSettings) error {
=======
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runUnparam(pass *analysis.Pass, settings *config.UnparamSettings) ([]goanalysis.Issue, error) {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	ssa := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	ssaPkg := ssa.Pkg

	pkg := &packages.Package{
		Fset:      pass.Fset,
		Syntax:    pass.Files,
		Types:     pass.Pkg,
		TypesInfo: pass.TypesInfo,
	}

	c := &check.Checker{}
	c.CheckExportedFuncs(settings.CheckExported)
	c.Packages([]*packages.Package{pkg})
	c.ProgramSSA(ssaPkg.Prog)

	unparamIssues, err := c.Check()
	if err != nil {
<<<<<<< HEAD
		return err
	}

	for _, i := range unparamIssues {
		pass.Report(analysis.Diagnostic{
			Pos:     i.Pos(),
			Message: i.Message(),
		})
	}

	return nil
=======
		return nil, err
	}

	var issues []goanalysis.Issue
	for _, i := range unparamIssues {
		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Pos:        pass.Fset.Position(i.Pos()),
			Text:       i.Message(),
			FromLinter: linterName,
		}, pass))
	}

	return issues, nil
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

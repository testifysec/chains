package recvcheck

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

<<<<<<< HEAD
// NewAnalyzer returns a new analyzer to check for receiver type consistency.
func NewAnalyzer(s Settings) *analysis.Analyzer {
	a := &analyzer{
		excluded: map[string]struct{}{},
	}

	if !s.DisableBuiltin {
		// Default excludes for Marshal/Encode methods https://github.com/raeperd/recvcheck/issues/7
		a.excluded = map[string]struct{}{
			"*.MarshalText":   {},
			"*.MarshalJSON":   {},
			"*.MarshalYAML":   {},
			"*.MarshalXML":    {},
			"*.MarshalBinary": {},
			"*.GobEncode":     {},
		}
	}

	for _, exclusion := range s.Exclusions {
		a.excluded[exclusion] = struct{}{}
	}

	return &analysis.Analyzer{
		Name:     "recvcheck",
		Doc:      "checks for receiver type consistency",
		Run:      a.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

// Settings is the configuration for the analyzer.
type Settings struct {
	// DisableBuiltin if true, disables the built-in method excludes.
	// Built-in excluded methods:
	//   - "MarshalText"
	//   - "MarshalJSON"
	//   - "MarshalYAML"
	//   - "MarshalXML"
	//   - "MarshalBinary"
	//   - "GobEncode"
	DisableBuiltin bool

	// Exclusions format is `struct_name.method_name` (ex: `Foo.MethodName`).
	// A wildcard `*` can use as a struct name (ex: `*.MethodName`).
	Exclusions []string
}

type analyzer struct {
	excluded map[string]struct{}
}

func (r *analyzer) run(pass *analysis.Pass) (any, error) {
=======
var Analyzer = &analysis.Analyzer{
	Name:     "recvcheck",
	Doc:      "checks for receiver type consistency",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	structs := map[string]*structType{}
	inspector.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Recv == nil || len(funcDecl.Recv.List) != 1 {
			return
		}

<<<<<<< HEAD
		recv, isStar := recvTypeIdent(funcDecl.Recv.List[0].Type)
		if recv == nil {
			return
		}

		if r.isExcluded(recv, funcDecl) {
			return
		}

		st, ok := structs[recv.Name]
		if !ok {
			structs[recv.Name] = &structType{}
=======
		var recv *ast.Ident
		var isStar bool
		switch recvType := funcDecl.Recv.List[0].Type.(type) {
		case *ast.StarExpr:
			isStar = true
			if recv, ok = recvType.X.(*ast.Ident); !ok {
				return
			}
		case *ast.Ident:
			recv = recvType
		default:
			return
		}

		var st *structType
		st, ok = structs[recv.Name]
		if !ok {
			structs[recv.Name] = &structType{recv: recv.Name}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			st = structs[recv.Name]
		}

		if isStar {
<<<<<<< HEAD
			st.starUsed = true
		} else {
			st.typeUsed = true
		}
	})

	for recv, st := range structs {
		if st.starUsed && st.typeUsed {
			pass.Reportf(pass.Pkg.Scope().Lookup(recv).Pos(), "the methods of %q use pointer receiver and non-pointer receiver.", recv)
=======
			st.numStarMethod++
		} else {
			st.numTypeMethod++
		}
	})

	for _, st := range structs {
		if st.numStarMethod > 0 && st.numTypeMethod > 0 {
			pass.Reportf(pass.Pkg.Scope().Lookup(st.recv).Pos(), "the methods of %q use pointer receiver and non-pointer receiver.", st.recv)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		}
	}

	return nil, nil
}

<<<<<<< HEAD
func (r *analyzer) isExcluded(recv *ast.Ident, f *ast.FuncDecl) bool {
	if f.Name == nil || f.Name.Name == "" {
		return true
	}

	_, found := r.excluded[recv.Name+"."+f.Name.Name]
	if found {
		return true
	}

	_, found = r.excluded["*."+f.Name.Name]

	return found
}

type structType struct {
	starUsed bool
	typeUsed bool
}

func recvTypeIdent(r ast.Expr) (*ast.Ident, bool) {
	switch n := r.(type) {
	case *ast.StarExpr:
		if i, ok := n.X.(*ast.Ident); ok {
			return i, true
		}

	case *ast.Ident:
		return n, false
	}

	return nil, false
=======
type structType struct {
	recv          string
	numStarMethod int
	numTypeMethod int
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

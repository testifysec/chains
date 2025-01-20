package goanalysis

import (
	"fmt"
	"runtime/debug"

	"github.com/golangci/golangci-lint/internal/errorutil"
)

type actionAllocator struct {
	allocatedActions []action
	nextFreeIndex    int
}

func newActionAllocator(maxCount int) *actionAllocator {
	return &actionAllocator{
		allocatedActions: make([]action, maxCount),
		nextFreeIndex:    0,
	}
}

func (actAlloc *actionAllocator) alloc() *action {
	if actAlloc.nextFreeIndex == len(actAlloc.allocatedActions) {
		panic(fmt.Sprintf("Made too many allocations of actions: %d allowed", len(actAlloc.allocatedActions)))
	}
	act := &actAlloc.allocatedActions[actAlloc.nextFreeIndex]
	actAlloc.nextFreeIndex++
	return act
}

func (act *action) waitUntilDependingAnalyzersWorked() {
<<<<<<< HEAD
	for _, dep := range act.Deps {
		if dep.Package == act.Package {
=======
	for _, dep := range act.deps {
		if dep.pkg == act.pkg {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			<-dep.analysisDoneCh
		}
	}
}

func (act *action) analyzeSafe() {
	defer func() {
		if p := recover(); p != nil {
<<<<<<< HEAD
			if !act.IsRoot {
				// This line allows to display "hidden" panic with analyzers like buildssa.
				// Some linters are dependent of sub-analyzers but when a sub-analyzer fails the linter is not aware of that,
				// this results to another panic (ex: "interface conversion: interface {} is nil, not *buildssa.SSA").
				act.runner.log.Errorf("%s: panic during analysis: %v, %s", act.Analyzer.Name, p, string(debug.Stack()))
			}

			act.Err = errorutil.NewPanicError(fmt.Sprintf("%s: package %q (isInitialPkg: %t, needAnalyzeSource: %t): %s",
				act.Analyzer.Name, act.Package.Name, act.isInitialPkg, act.needAnalyzeSource, p), debug.Stack())
		}
	}()

	act.runner.sw.TrackStage(act.Analyzer.Name, act.analyze)
=======
			if !act.isroot {
				// This line allows to display "hidden" panic with analyzers like buildssa.
				// Some linters are dependent of sub-analyzers but when a sub-analyzer fails the linter is not aware of that,
				// this results to another panic (ex: "interface conversion: interface {} is nil, not *buildssa.SSA").
				act.r.log.Errorf("%s: panic during analysis: %v, %s", act.a.Name, p, string(debug.Stack()))
			}

			act.err = errorutil.NewPanicError(fmt.Sprintf("%s: package %q (isInitialPkg: %t, needAnalyzeSource: %t): %s",
				act.a.Name, act.pkg.Name, act.isInitialPkg, act.needAnalyzeSource, p), debug.Stack())
		}
	}()

	act.r.sw.TrackStage(act.a.Name, act.analyze)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

func (act *action) markDepsForAnalyzingSource() {
	// Horizontal deps (analyzer.Requires) must be loaded from source and analyzed before analyzing
	// this action.
<<<<<<< HEAD
	for _, dep := range act.Deps {
		if dep.Package == act.Package {
=======
	for _, dep := range act.deps {
		if dep.pkg == act.pkg {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
			// Analyze source only for horizontal dependencies, e.g. from "buildssa".
			dep.needAnalyzeSource = true // can't be set in parallel
		}
	}
}

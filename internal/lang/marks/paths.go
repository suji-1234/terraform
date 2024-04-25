// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package marks

import (
	"github.com/zclconf/go-cty/cty"
)

// PathsWithMark produces a list of paths identified as having a specified
// mark in a given set of [cty.PathValueMarks] that presumably resulted from
// deeply-unmarking a [cty.Value].
//
// This is for situations where a subsystem needs to give special treatment
// to one specific mark value, as opposed to just handling all marks
// generically as cty operations would. The second return value is a
// subset of the given [cty.PathValueMarks] values which contained marks
// other than the one requested, so that a caller that can't preserve other
// marks at all can more easily return an error explaining that.
func PathsWithMark(pvms []cty.PathValueMarks, wantMark any) (withWanted []cty.Path, withOthers []cty.PathValueMarks) {
	if len(pvms) == 0 {
		// No-allocations path for the common case where there are no marks at all.
		return nil, nil
	}

	for _, pvm := range pvms {
		if _, ok := pvm.Marks[wantMark]; ok {
			withWanted = append(withWanted, pvm.Path)
		}
		for mark := range pvm.Marks {
			if mark != wantMark {
				withOthers = append(withOthers, pvm)
			}
		}
	}

	return withWanted, withOthers
}

// MarkPaths transforms the given value by marking each of the given paths
// with the given mark value.
func MarkPaths(val cty.Value, mark any, paths []cty.Path) cty.Value {
	if len(paths) == 0 {
		// No-allocations path for the common case where there are no marked paths at all.
		return val
	}

	// For now we'll use cty's slightly lower-level function to achieve this
	// result. This is a little inefficient due to an additional dynamic
	// allocation for the intermediate data structure, so if that becomes
	// a problem in practice then we may wish to write a more direct
	// implementation here.
	markses := make([]cty.PathValueMarks, len(paths))
	marks := cty.NewValueMarks(mark)
	for i, path := range paths {
		markses[i] = cty.PathValueMarks{
			Path:  path,
			Marks: marks,
		}
	}
	return val.MarkWithPaths(markses)
}

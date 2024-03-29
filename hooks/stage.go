// Code generated by "stringer -type=HookStage -linecomment -output=stage.go"; DO NOT EDIT.

package hooks

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[StagePreReplace-0]
	_ = x[StagePreCommit-1]
	_ = x[StagePreTag-2]
	_ = x[StagePrePush-3]
}

const _HookStage_name = "Pre ReplacePre CommitPre TagPre Push"

var _HookStage_index = [...]uint8{0, 11, 21, 28, 36}

func (i HookStage) String() string {
	if i < 0 || i >= HookStage(len(_HookStage_index)-1) {
		return "HookStage(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _HookStage_name[_HookStage_index[i]:_HookStage_index[i+1]]
}

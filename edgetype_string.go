// Code generated by "stringer -linecomment -trimprefix=EdgeType -type=EdgeType"; DO NOT EDIT.

package myModel

import "strconv"

const _EdgeType_name = "MultiMultiMultiOneMultiOneParentOneMultiOneOneOneOneParent"

var _EdgeType_index = [...]uint8{0, 10, 18, 32, 40, 46, 58}

func (i EdgeType) String() string {
	if i >= EdgeType(len(_EdgeType_index)-1) {
		return "EdgeType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _EdgeType_name[_EdgeType_index[i]:_EdgeType_index[i+1]]
}
func EdgeTypeValue(name string) EdgeType {
	for _, i := range _EdgeType_index {
		if name == _EdgeType_name[_EdgeType_index[i]:_EdgeType_index[i+1]] {
			return EdgeType(i)
		}
	}
	panic("not a constant of EdgeType")
}

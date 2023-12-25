package utils

import (
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnique(t *testing.T) {
	Convey("test-unique", t, func() {
		So(IntsUnique([]int{1, 2, 3, 3, 3, 3}), ShouldResemble, []int{1, 2, 3})

		// uinque will break order
		s := StringsUnique([]string{"a", "b", "c", "c", "d", "d"})
		sort.Sort(sort.StringSlice(s))
		So(s, ShouldResemble, []string{"a", "b", "c", "d"})

		So(Int64Unique([]int64{3, 2, 1, 1, 6}), ShouldResemble, []int64{1, 2, 3, 6})

		So(IsInString("xyz", []string{"a", "xz", "xyz"}), ShouldBeTrue)
		So(IsInString("y", []string{"a", "xz", "xyz"}), ShouldBeFalse)
	})
}

package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLatency(t *testing.T) {
	Convey("test-latency-shuffer", t, func() {
		sf := NewShuffle(3)
		idx := sf.Get()
		So(idx, ShouldBeIn, [...]int64{0, 1, 2})
		sf.TryNext()
		idx2 := sf.Get()
		So(idx2, ShouldBeIn, [...]int64{0, 1, 2})
		So(idx == idx2, ShouldBeFalse)
	})
}

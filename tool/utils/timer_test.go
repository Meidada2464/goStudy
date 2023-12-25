package utils

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTimer(t *testing.T) {
	Convey("test-timer", t, func() {
		Convey("timer-loop", func() {
			var counter int
			cancelFn := TimeLoop(time.Second/10, func(now time.Time) {
				counter++
			})
			time.Sleep(time.Second)
			cancelFn()
			So(counter, ShouldBeGreaterThanOrEqualTo, 10)
		})
		Convey("timer-loop-then", func() {
			var counter int
			cancelFn := TimeLoopThen(time.Second/10, func(now time.Time) {
				counter++
			})
			time.Sleep(time.Second)
			cancelFn()
			So(counter, ShouldBeGreaterThanOrEqualTo, 9) // once lost
			So(counter, ShouldBeLessThan, 11)
		})
	})
	Convey("test-intime", t, func() {
		So(IsInMinute(3600*2, 3600, 3600*10), ShouldBeTrue)   // 02:00 in 01:00-10:00
		So(IsInMinute(3600*12, 3600, 3600*10), ShouldBeFalse) // 12:00 not int 01:00-10:00
		So(IsInMinute(3600*2, 3600*10, 3600), ShouldBeFalse)  //02:00 not in 10:00-01:00(nextday)
	})
}

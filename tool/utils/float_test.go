package utils

import (
	"errors"
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFloat(t *testing.T) {
	Convey("test-float", t, func() {
		Convey("fixfloat", func() {
			So(FixFloat(8.233334), ShouldEqual, 8.23)
			So(FixFloat(8.233834, 4), ShouldEqual, 8.2338)
		})
		Convey("fixzero", func() {
			So(FixZero(8.22), ShouldEqual, 8.22)
			So(FixZero(-7), ShouldEqual, 0)
		})
		Convey("to-float", func() {
			_, err := ToFloat(nil)
			So(err, ShouldNotBeNil)

			f, _ := ToFloat("1.23")
			So(f, ShouldEqual, 1.23)
			_, err = ToFloat("1.23x")
			So(err, ShouldNotBeNil)

			f, _ = ToFloat(float32(1.23))
			So(math.Abs(f-1.23), ShouldBeLessThanOrEqualTo, 10e-5)
			f, _ = ToFloat(float64(1.23))
			So(math.Abs(f-1.23), ShouldBeLessThanOrEqualTo, 10e-5)

			f, _ = ToFloat(int(11))
			So(math.Abs(f-11), ShouldBeLessThanOrEqualTo, 10e-5)
			f, _ = ToFloat(int32(11))
			So(math.Abs(f-11), ShouldBeLessThanOrEqualTo, 10e-5)
			f, _ = ToFloat(int64(11))
			So(math.Abs(f-11), ShouldBeLessThanOrEqualTo, 10e-5)
			f, _ = ToFloat(uint(11))
			So(math.Abs(f-11), ShouldBeLessThanOrEqualTo, 10e-5)
			f, _ = ToFloat(uint32(11))
			So(math.Abs(f-11), ShouldBeLessThanOrEqualTo, 10e-5)
			f, _ = ToFloat(uint64(11))
			So(math.Abs(f-11), ShouldBeLessThanOrEqualTo, 10e-5)

			_, err = ToFloat(errors.New("wrong"))
			So(err, ShouldNotBeNil)
		})
	})
}

package utils

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCommand(t *testing.T) {
	Convey("test-command", t, func() {
		out, err := RunCommand("echo", []string{"a"}, time.Second)
		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "a\n")

		out, err = RunCommand("sleep", []string{"10s"}, time.Second)
		So(len(out) == 0, ShouldBeTrue)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "killed")
	})
}

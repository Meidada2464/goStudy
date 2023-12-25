package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMD5(t *testing.T) {
	Convey("test-md5", t, func() {
		hash := "e807f1fcf82d132f9bb018ca6738a19f"
		So(MD5String("1234567890"), ShouldEqual, hash)
		So(MD5Bytes([]byte("1234567890")), ShouldEqual, hash)
		fileHash, _ := MD5File("md5_test.txt")
		So(fileHash, ShouldEqual, hash)
	})
}

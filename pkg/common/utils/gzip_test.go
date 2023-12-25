package utils

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGzip(t *testing.T) {
	Convey("test-gzip", t, func() {
		values := map[string]interface{}{
			"abc":  "anc",
			"next": 99.88,
			"eext": []interface{}{"xuz", 2.1},
		}

		Convey("gzip", func() {
			gzipData, err := GzipJSONBytes(values)
			So(err, ShouldBeNil)
			So(gzipData, ShouldNotBeEmpty)

			data, _ := json.Marshal(values)
			gzipData2, err := Gzip(data)
			So(err, ShouldBeNil)
			So(bytes.Equal(gzipData, gzipData2), ShouldBeTrue)

			Convey("ungzip", func() {
				value2 := make(map[string]interface{})
				err := UnGzipJSONBytes(gzipData, &value2)
				So(err, ShouldBeNil)
				So(value2, ShouldResemble, values)

				value3 := make(map[string]interface{})
				err = UnGzipJSON(bytes.NewReader(gzipData), &value3)
				So(err, ShouldBeNil)
				So(value3, ShouldResemble, values)

				value2 = make(map[string]interface{})
				data2, err := UnGzip(gzipData2)
				So(err, ShouldBeNil)
				So(data2, ShouldResemble, data)
			})
		})
	})
}

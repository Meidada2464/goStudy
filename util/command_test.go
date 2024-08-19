/**
 * Package util
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2024/1/18
 */

package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRunCommand(t *testing.T) {
	stdout, stderr, err := RunCommand("echo", []string{"hello"}, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "hello\n", string(stdout))
	assert.Equal(t, "", string(stderr))
}

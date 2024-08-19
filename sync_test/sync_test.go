package sync_test

import (
	"sync"
	"testing"
)

type ss struct {
	ab string
	cd string
}

func TestSyncMap(t *testing.T) {
	var a sync.Map

	a.Store("a", ss{
		ab: "ab",
		cd: "cd",
	})

	a.Store("b", ss{
		ab: "aa",
		cd: "cc",
	})

	a.Store("c", ss{
		ab: "aac",
		cd: "ccc",
	})

	a.Store("d", ss{
		ab: "aacd",
		cd: "cccd",
	})

	value, ok := a.Load("a")
	if ok {
		ns := value.(ss)
		ns.cd = "newcd"
		ns.ab = "newab"
		a.Store("a", ns)
	}

	a.Range(func(key, value interface{}) bool {
		t.Log(key, value)
		if key == "b" {
			return true
		}
		return true
	})
}

func FuzzName(f *testing.F) {
	f.Fuzz(func(t *testing.T) {

	})
}

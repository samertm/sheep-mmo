package engine

import (
	"testing"
)

func TestProximity1(t *testing.T) {
	if !(proximate(box{1,2,3,4}, box{2,3,4,5}, 0)) {
		t.Fail()
	}
}

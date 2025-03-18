package utils

import "testing"

func TestUntarzip(t *testing.T) {
	Untargz("utils-test-data/test.tar.gz", "utils-test-data/ttt")
}

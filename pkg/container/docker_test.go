package container

import "testing"

func TestLoadAndPush(t *testing.T) {
	d := newDocker()
	err := d.LoadAndPush("container-test-data/test.tar", "reg.safedog.cn/safedog")
	if err != nil {
		panic(err)
	}
}

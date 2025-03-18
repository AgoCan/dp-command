package helm

import (
	"fmt"
	"testing"
)

func TestRender3(t *testing.T) {
	cr := NewChartRenderer("")
	res, err := cr.RenderChart("./helm-test-data/passport", "passport", []string{}, map[string]any{})
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

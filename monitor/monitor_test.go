package monitor

import (
	"fmt"
	"testing"
	"time"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xtest"
)

func TestMonitor(t *testing.T) {
	var m = New()
	_, err := xtest.DoPerfV_(200, 30, "",
		func(i int) error {
			if i%10 == 0 {
				m.State()
				return nil
			}
			var id = m.Start(fmt.Sprintf("_%v", i%3))
			time.Sleep(time.Duration(i) * time.Millisecond)
			m.Done(id)
			return nil
		})
	if err != nil {
		t.Error(err)
	}
	val, _ := m.State()
	fmt.Println(converter.JSON(val))
}

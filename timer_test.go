package timer

import (
	"time"
	"testing"
	"fmt"
)

func TestC(t *testing.T) {
	go Start()

	fn := func() {
		fmt.Println("ok")
	}
	AddOnce("1", time.Second*2, fn)
	select {}
}

package queue

import (
	"time"
)

func RunQueue() {
	for {
		if size != 0 {
			e := Top()
			e.process()
			pop()
			e.remove()
		}
		time.Sleep(500 * time.Millisecond)
	}
}

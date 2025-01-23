package queue

import (
	"fmt"
	"sync"

	"github.com/rzaf/youtube-clone/file/models"
)

type Element interface {
	process() error
	name() string
	url() string
}

var (
	queue       = make([]Element, 0)
	size        = 0
	queueMutext sync.Mutex
)

func Push(e Element) {
	queueMutext.Lock()
	err := models.CreateProcess(e.name(), e.url())
	fmt.Printf("err: %v\n", err)
	queue = append(queue, e)
	size++
	queueMutext.Unlock()
}

func Top() Element {
	queueMutext.Lock()
	defer queueMutext.Unlock()
	return queue[0]
}

func pop() {
	queueMutext.Lock()
	queue = queue[1:]
	size--
	queueMutext.Unlock()
}

func Size() int {
	queueMutext.Lock()
	defer queueMutext.Unlock()
	return size
}

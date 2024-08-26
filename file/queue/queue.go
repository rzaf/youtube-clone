package queue

import (
	"sync"
)

// type MediaType uint8

// const (
//
//	Video MediaType = iota
//	Music
//	Photo
//
// )
type Element interface {
	process()
	remove()
}

var (
	queue       = make([]Element, 0)
	size        = 0
	queueMutext sync.Mutex
)

func Push(e Element) {
	queueMutext.Lock()
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

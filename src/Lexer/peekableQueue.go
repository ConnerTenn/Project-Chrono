package Lexer

import (
	"fmt"
	"os"
)

func displayError(str string) {
	fmt.Println("Error: ", str)
	os.Exit(-1)
}

// if there is only one reader thread, the lock is not needed
// move to list for peeking ahead multiple tokens?
type PeekableQueue struct {
	queue  chan Token
	peeked Token
}

func NewQueue() PeekableQueue {
	var q PeekableQueue
	q.queue = make(chan Token, 20)
	q.peeked.Value = ""
	return q
}

func (q PeekableQueue) Close() {
	close(q.queue)
}

func (q PeekableQueue) PushBack(t Token) {
	q.queue <- t
}

func (q *PeekableQueue) PeekNext() Token {
	var retVal Token

	if q.peeked.Value == "" {
		var ok bool
		q.peeked, ok = <-q.queue

		if !ok {
			displayError("Unexpected end of file")
		}
	}

	retVal = q.peeked

	return retVal
}

func (q *PeekableQueue) GetNext() Token {
	var retVal Token

	if q.peeked.Value == "" {
		var ok bool
		retVal, ok = <-q.queue

		if !ok {
			displayError("Unexpected end of file")
		}
	} else {
		retVal = q.peeked
		q.peeked.Value = ""
	}

	return retVal
}

func (q *PeekableQueue) IsEmpty() bool {
	ok := true

	if q.peeked.Value == "" {
		q.peeked, ok = <-q.queue
	}

	return !ok
}

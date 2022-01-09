package main

import "sync"

// if there is only one reader thread, the lock is not needed
// move to list for peeking ahead multiple tokens?
type PeekableQueue struct {
  queue chan Token
  peeked Token
  ok bool
  mut sync.Mutex
}

func NewQueue() PeekableQueue {
  var q PeekableQueue
  q.queue = make(chan Token, 20)
  q.peeked.Type = -1
  return q
}

func (q PeekableQueue) Close() {
  close(q.queue)
}

func (q PeekableQueue) PushBack(t Token) {
  q.queue <- t
}

func (q PeekableQueue) PeekNext() Token {
  q.mut.Lock()
  var retVal Token
  if q.peeked.Type == -1 {
    q.peeked, q.ok = <- q.queue
  }

  retVal = q.peeked

  q.mut.Unlock()
  return retVal
}

func (q PeekableQueue) GetNext() (Token, bool) {
  q.mut.Lock()
  var (
    retVal Token
    ok bool
  )

  if q.peeked.Type == -1 {
    retVal, ok = <- q.queue
  } else {
    retVal = q.peeked
    ok = q.ok
    q.peeked.Type = -1
  }

  q.mut.Unlock()
  return retVal, ok
}

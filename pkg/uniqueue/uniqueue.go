package uniqueue

import (
	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	errMultipleCallsToStart = errors.New("unable to start more than once")
	log                     = logging.LoggerForModule()
)

// UniQueue is a queue of unique elements.
// It uses channels to push and pop from the internal queue.
// It is thread safe.
type UniQueue[T comparable] struct {
	stopper      concurrency.Stopper
	mu           sync.Mutex
	queueSize    int
	backChannel  chan T
	frontChannel chan T
	queueChannel chan T
	inQueue      map[T]struct{}
}

// NewUniQueue creates a new UniQueue.
func NewUniQueue[T comparable](size int) *UniQueue[T] {
	return &UniQueue[T]{
		stopper:      concurrency.NewStopper(),
		queueSize:    size,
		backChannel:  nil,
		queueChannel: nil,
		frontChannel: nil,
		inQueue:      nil,
	}
}

// PushC returns the channel from which the client can push elements to the UniQueue.
func (q *UniQueue[T]) PushC() chan<- T {
	if q.inQueue == nil {
		log.Panic("Start must be called before PushC")
	}
	return q.backChannel
}

// PopC returns the channel from which the client can pop elements from the UniQueue.
func (q *UniQueue[T]) PopC() <-chan T {
	if q.inQueue == nil {
		log.Panic("Start must be called before PopC")
	}
	return q.frontChannel
}

// Start initializes and runs the UniQueue.
// It must be called before attempting to push or pop any elements.
func (q *UniQueue[T]) Start() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.inQueue != nil {
		return errMultipleCallsToStart
	}
	q.stopper.LowLevel().ResetStopRequest()
	q.backChannel = make(chan T, 1)
	q.queueChannel = make(chan T, q.queueSize)
	q.frontChannel = make(chan T, 1)
	q.inQueue = make(map[T]struct{})
	go q.run()
	return nil
}

func (q *UniQueue[T]) run() {
	defer q.stopper.Flow().ReportStopped()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go q.pushFromBack(wg)
	go q.popToFront(wg)
	// Wait for the push and pop goroutines to finish
	wg.Wait()
	// Close channels
	close(q.backChannel)
	close(q.queueChannel)
	close(q.frontChannel)
	q.inQueue = nil
}

// Stop the UniQueue.
func (q *UniQueue[T]) Stop() {
	if !q.stopper.Client().Stopped().IsDone() {
		defer func() {
			_ = q.stopper.Client().Stopped().Wait()
		}()
	}
	q.stopper.Client().Stop()
}

func (q *UniQueue[T]) pushFromBack(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-q.stopper.Flow().StopRequested():
			return
		case el, ok := <-q.backChannel:
			if !ok {
				return
			}
			if q.maybeAddToQueue(el) {
				q.queueChannel <- el
			}
		}
	}
}

func (q *UniQueue[T]) maybeAddToQueue(el T) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, alreadyInQueue := q.inQueue[el]; alreadyInQueue {
		return false
	}
	q.inQueue[el] = struct{}{}
	return true
}

func (q *UniQueue[T]) popToFront(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-q.stopper.Flow().StopRequested():
			return
		case el, ok := <-q.queueChannel:
			if !ok {
				return
			}
			select {
			case <-q.stopper.Flow().StopRequested():
				return
			case q.frontChannel <- el:
				q.removeFromQueue(el)
			}
		}
	}
}

func (q *UniQueue[T]) removeFromQueue(el T) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.inQueue[el]; ok {
		delete(q.inQueue, el)
		return true
	}
	return false
}

// IsEmpty indicates whether the UniQueue is empty or not.
func (q *UniQueue[T]) IsEmpty() bool {
	return len(q.queueChannel) == 0 && len(q.frontChannel) == 0
}

// Size returns the number of elements in the UniQueue.
func (q *UniQueue[T]) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.inQueue) + len(q.frontChannel)
}

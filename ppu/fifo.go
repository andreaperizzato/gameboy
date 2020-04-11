package ppu

// FIFO defines a FIFO queue.
type FIFO interface {
	Clear()
	Push(uint8) bool
	Pop() (uint8, bool)
	Size() int
}

// FIFOQueue is an implementation of FIFO.
type FIFOQueue struct {
	data []uint8
	idx  int
	w    int
	r    int
}

// NewFIFOQueue returns a new FIFOQueue with the given length.
func NewFIFOQueue(length int) *FIFOQueue {
	return &FIFOQueue{
		data: make([]uint8, length),
		idx:  -1,
	}
}

// Push puses a value in the queue and returns true if the queue wasn't full.
func (q *FIFOQueue) Push(v uint8) bool {
	if q.idx == len(q.data)-1 {
		return false
	}
	q.idx++
	q.data[q.idx] = v
	return true
}

// Pop pops a value from the queue and returns true when the queue wasn't empty.
func (q *FIFOQueue) Pop() (uint8, bool) {
	if q.idx == -1 {
		return 0, false
	}
	v := q.data[0]
	q.idx--
	for i := 0; i < len(q.data)-1; i++ {
		q.data[i] = q.data[i+1]
	}
	return v, true
}

// Clear clears the queue.
func (q *FIFOQueue) Clear() {
	q.idx = -1
}

// Size returns the number of elements in the queue.
func (q *FIFOQueue) Size() int {
	return q.idx + 1
}

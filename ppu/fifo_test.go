package ppu_test

import (
	"testing"

	"github.com/andreaperizzato/gameboy/ppu"
	"github.com/stretchr/testify/assert"
)

func TestFIFOQueue(t *testing.T) {
	q := ppu.NewFIFOQueue(3)

	assert.Equal(t, 0, q.Size(), "should be initially empty")
	_, ok := q.Pop()
	assert.False(t, ok, "should not pop from an empty queue")

	ok = q.Push(0xAA)
	assert.True(t, ok, "should push")
	assert.Equal(t, 1, q.Size(), "should have a size of 1")

	ok = q.Push(0xBB)
	assert.True(t, ok, "should push")
	assert.Equal(t, 2, q.Size(), "should have a size of 2")

	ok = q.Push(0xCC)
	assert.True(t, ok, "should push")
	assert.Equal(t, 3, q.Size(), "should have a size of 3")

	ok = q.Push(0xDD)
	assert.False(t, ok, "should not push when full")
	assert.Equal(t, 3, q.Size(), "should still have a size of 3")

	v, ok := q.Pop()
	assert.True(t, ok, "should pop")
	assert.Equal(t, uint8(0xAA), v, "should pop the first added")
	assert.Equal(t, 2, q.Size(), "should have a size of 2")

	v, ok = q.Pop()
	assert.True(t, ok, "should pop")
	assert.Equal(t, uint8(0xBB), v, "should pop the second added")
	assert.Equal(t, 1, q.Size(), "should have a size of 1")

	v, ok = q.Pop()
	assert.True(t, ok, "should pop")
	assert.Equal(t, uint8(0xCC), v, "should pop the third added")
	assert.Equal(t, 0, q.Size(), "should have a size of 0")

	_, ok = q.Pop()
	assert.False(t, ok, "should not pop from an empty queue")

	ok = q.Push(0x11)
	assert.True(t, ok, "should push")
	ok = q.Push(0x22)
	assert.True(t, ok, "should push")
	assert.Equal(t, 2, q.Size(), "should have a size of 2")
	q.Clear()
	assert.Equal(t, 0, q.Size(), "should be empty")
	_, ok = q.Pop()
	assert.False(t, ok, "should not pop from an empty queue")
}

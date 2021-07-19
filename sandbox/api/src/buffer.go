package main

import "fmt"

type buffer struct {
	buffer []string
	size   int
}

func NewBuffer(size int) *buffer {

	buf := make([]string, 0, size)
	return &buffer{buf, size}
}

func (buffer *buffer) Push(token string) bool {
	if len(buffer.buffer) == buffer.size {
		return false
	}
	buffer.buffer = append(buffer.buffer, token)
	return true
}

func (buffer *buffer) Skim() (string, bool) {
	length := len(buffer.buffer)
	if length != buffer.size {
		return "", false
	}
	word := buffer.buffer[0]
	buffer.buffer = buffer.buffer[1:]
	return word, true
}

func (buffer *buffer) Shift() (string, bool) {
	if len(buffer.buffer) == 0 {
		return "", false
	}
	word := buffer.buffer[0]
	buffer.buffer = buffer.buffer[1:]
	return word, true
}

func (buffer *buffer) Replace(offset int, newWord string) bool {
	last := len(buffer.buffer) - 1
	if last-offset < 0 {
		return false
	}
	buffer.buffer[last-offset] = newWord
	return true
}

func (buffer *buffer) Read(offset int) (string, bool) {
	last := len(buffer.buffer) - 1
	if last-offset < 0 {
		return "", false
	}
	return buffer.buffer[last-offset], true
}

func (buffer *buffer) Print() {
	fmt.Println("Printing buffer:")
	for _, word := range buffer.buffer {
		fmt.Println(word)
	}
}

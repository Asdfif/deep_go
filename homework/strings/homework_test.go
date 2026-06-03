package main

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type COWBuffer struct {
	data    []byte
	refs    *int
	refData *[]byte
	// need to implement
}

func NewCOWBuffer(data []byte) COWBuffer {
	newRefs := new(int)
	*newRefs = 1
	return COWBuffer{
		data: data,
		refs: newRefs,
	}
}

// func (b *COWBuffer) Clone() COWBuffer                  // создать новую копию буфера
func (b *COWBuffer) Clone() COWBuffer {
	*b.refs++

	return COWBuffer{
		data: b.data,
		refs: b.refs,
	} // need to implement
}

// func (b *COWBuffer) Close()                            // перестать использовать копию буффера
func (b *COWBuffer) Close() {
	if b.refs != nil && *b.refs > 0 {
		*b.refs--
	}
}

// func (b *COWBuffer) Update(index int, value byte) bool // изменить определенный байт в буффере
func (b *COWBuffer) Update(index int, value byte) bool {
	if index > len(b.data)-1 || index < 0 {
		return false
	}

	if *b.refs > 1 {
		*b.refs--
		copyData := make([]byte, len(b.data))
		copy(copyData, b.data)
		*b = NewCOWBuffer(copyData)
	}
	b.data[index] = value

	return true
}

// func (b *COWBuffer) String() string                    // сконвертировать буффер в строку
func (b *COWBuffer) String() string {
	return unsafe.String(unsafe.SliceData(b.data), len(b.data))
}

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()
	fmt.Printf("buffer %+v %v\n", buffer, *buffer.refs)

	copy1 := buffer.Clone()
	fmt.Printf("copy1 %+v %v\n", copy1, *copy1.refs)
	fmt.Printf("buffer %+v %v\n", buffer, *buffer.refs)

	copy2 := buffer.Clone()
	fmt.Printf("copy2 %+v %v\n", copy2, *copy2.refs)
	fmt.Printf("buffer %+v %v\n", buffer, *buffer.refs)

	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data))
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	fmt.Printf("1) %v :: %v\n", (*byte)(unsafe.SliceData(data)), unsafe.StringData(buffer.String()))
	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()))
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()))
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()))

	fmt.Printf("2) buffer %+v %v\n", buffer, *buffer.refs)
	fmt.Printf("2) copy1 %+v %v\n", copy1, *copy1.refs)
	fmt.Printf("2) copy2 %+v %v\n", copy2, *copy2.refs)
	print("\n")
	assert.True(t, buffer.Update(0, 'g'))
	fmt.Printf("3) buffer %+v %v\n", buffer, *buffer.refs)
	fmt.Printf("3) copy1 %+v %v\n", copy1, *copy1.refs)
	fmt.Printf("3) copy2 %+v %v\n", copy2, *copy2.refs)
	print("\n")
	assert.False(t, buffer.Update(-1, 'g'))
	fmt.Printf("4) buffer %+v %v\n", buffer, *buffer.refs)
	fmt.Printf("4) copy1 %+v %v\n", copy1, *copy1.refs)
	fmt.Printf("4) copy2 %+v %v\n", copy2, *copy2.refs)
	print("\n")
	assert.False(t, buffer.Update(4, 'g'))
	fmt.Printf("5) buffer %+v %v\n", buffer, *buffer.refs)

	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	copy1.Close()

	previous := copy2.data
	copy2.Update(0, 'f')
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current))

	copy2.Close()
}

package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node struct {
	key    int
	value  int
	parent *Node
	left   *Node
	right  *Node
}

type OrderedMap struct {
	size int
	root *Node
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		size: 0,
		root: nil,
	}
}

// func (m *OrderedMap) Insert(key, value int)          // добавить элемент в словарь
func (m *OrderedMap) Insert(key, value int) {
	if m.size == 0 {
		m.root = &Node{key: key, value: value}
		m.size++
		return
	}

	currentNode := m.root

	fmt.Printf("INSERT %v %v\n", key, value)

	for {
		if currentNode.key == key {
			return
		} else if key > currentNode.key {
			if currentNode.right == nil {
				currentNode.right = &Node{key: key, value: value}
				m.size++
				return
			}
			currentNode = currentNode.right
		} else if key < currentNode.key {
			if currentNode.left == nil {
				currentNode.left = &Node{key: key, value: value}
				m.size++
				return
			}
			currentNode = currentNode.left
		}
	}
}

func (m *OrderedMap) Find(key int) *Node {
	if m.root == nil {
		return nil
	}

	return FindNode(m.root, key)
}

func FindNode(node *Node, key int) *Node {
	if node == nil {
		return nil
	}

	fmt.Printf("FIND NODE %v key %v\n", node, key)

	if node.key == key {
		return node
	} else if key > node.key {
		return FindNode(node.right, key)
	} else {
		return FindNode(node.left, key)
	}
}

// func (m *OrderedMap) Erase(key int)                  // удалить элемент из словари
func (m *OrderedMap) Erase(key int) {
	if m.root == nil {
		return
	}

	m.root = m.EraseNode(m.root, key)
}

func (m *OrderedMap) EraseNode(node *Node, key int) *Node {
	if node == nil {
		return nil
	}

	if key > node.key {
		node.right = m.EraseNode(node.right, key)
		return node
	}

	if key < node.key {
		node.left = m.EraseNode(node.left, key)
		return node
	}

	if node.left == nil && node.right == nil {
		m.size--
		return nil
	}

	if node.left == nil {
		m.size--
		return node.right
	}

	if node.right == nil {
		m.size--
		return node.left
	}

	minNode := Left(node.right)

	node.key = minNode.key
	node.value = minNode.value

	node.right = m.EraseNode(node.right, minNode.key)
	return node
}

// func (m *OrderedMap) Contains(key int) bool          // проверить существование элемента в словаре
func (m *OrderedMap) Contains(key int) bool {
	return m.Find(key) != nil
}

// func (m *OrderedMap) Size() int                      // получить количество элементов в словаре
func (m *OrderedMap) Size() int {
	return m.size
}

// func (m *OrderedMap) ForEach(action func(int, int))  // применить функцию к каждому элементу словаря от меньшего к большему
func (m *OrderedMap) ForEach(action func(int, int)) {
	walk(m.root, action)
}

func walk(node *Node, action func(int, int)) {
	if node == nil {
		return // дошли до пустого места — возвращаемся
	}

	walk(node.left, action)
	action(node.key, node.value)
	walk(node.right, action)
}

func Left(currentNode *Node) *Node {
	if currentNode.left == nil {
		return currentNode
	}

	return Left(currentNode.left)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	// 		10
	// 	5	  	 15
	// 2		4 12
	// 			14

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, v int) {
		fmt.Printf("NODE %v %v\n", key, v)

		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	fmt.Printf("%v %v\n", keys, expectedKeys)

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}

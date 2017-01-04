package bytequeue

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Size constants
const (
	KB = 1 << 10
	MB = 1 << 20
)

const (
	// Number of bytes used to keep the entry size information in the header.
	headerEntrySize = 4
)

// ByteQueue is a non-thread safe queue.
type ByteQueue struct {
	byteArr      []byte
	head         int
	tail         int
	count        int // number of entries
	capacity     int
	headerBuffer []byte
}

type queueError struct {
	message string
}

// NewByteQueue initializes new ByteQueue.
// capacity: in MB.
func NewByteQueue(capacityMB int) *ByteQueue {
	return &ByteQueue{
		byteArr:      make([]byte, capacityMB),
		capacity:     capacityMB,
		headerBuffer: make([]byte, headerEntrySize),
	}
}

func (bq *ByteQueue) GetByteArr() []byte {
	return bq.byteArr
}

func (bq *ByteQueue) GetHead() int {
	return bq.head
}
func (bq *ByteQueue) GetTail() int {
	return bq.tail
}

func (bq *ByteQueue) getNextHead() {
	// get header
	for i := 0; i < headerEntrySize; i++ {
		bq.headerBuffer[i] = bq.byteArr[bq.head]
		bq.byteArr[bq.head] = 'X'
		bq.head++

		if bq.head == bq.capacity {
			bq.head = 0
		}
	}

	// data
	dataLen := int(binary.BigEndian.Uint32(bq.headerBuffer))

	for i := 0; i < dataLen; i++ {
		bq.byteArr[bq.head] = 'X'
		bq.head++

		if bq.head == bq.capacity {
			bq.head = 0
		}
	}
	//bq.head += dataLen
	//return bq.head
}

func (bq *ByteQueue) Pop() {
	fmt.Printf("Pop: h:%d\tt:%d\ta:%d\n", bq.head, bq.tail, bq.availableSpaceAfterTail())
	fmt.Printf("byteArr (befor pop): %02v\n", bq.GetByteArr())
	bq.getNextHead()
	fmt.Printf("Pop: h:%d\tt:%d\ta:%d\n", bq.head, bq.tail, bq.availableSpaceAfterTail())
	fmt.Printf("byteArr (after pop): %02v\n", bq.GetByteArr())
}

// Push ...
// return the number of bytes copied
func (bq *ByteQueue) Push(data []byte) (int, error) {
	dataLen := len(data)
	entryLen := headerEntrySize + dataLen

	if entryLen > bq.capacity {
		return 0, errors.New("Entry size is bigger than capacity.")
	}

	popCount := 0

	for {
		if entryLen > bq.availableSpaceAfterTail() {
			// pop some entries until the space is enough
			// also check do not exceed the size.
			bq.Pop()

			popCount++
		} else {
			//fmt.Printf("There are %d pops\n", popCount)
			break
		}
	}

	// copy header
	binary.BigEndian.PutUint32(bq.headerBuffer, uint32(dataLen))

	bq.setByteArr(bq.headerBuffer)

	// copy data
	bq.setByteArr(data)

	//
	bq.count++

	return bq.tail, nil
}

func (bq *ByteQueue) setByteArr(data []byte) {
	for _, v := range data {
		bq.byteArr[bq.tail] = v
		bq.tail++
		if bq.tail == bq.capacity {
			bq.tail = 0
		}
	}
}

func (bq *ByteQueue) availableSpaceAfterTail() int {
	if bq.tail >= bq.head {
		//return bq.capacity - (bq.tail - bq.head)
		return bq.capacity - bq.tail + bq.head
	}

	return bq.head - bq.tail
}

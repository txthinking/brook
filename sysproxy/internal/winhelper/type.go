package winhelper

import (
	"bytes"
	"unsafe"
)

/*
	Constructed based on Windows Data Types
	Reference: https://docs.microsoft.com/en-us/windows/desktop/winprog/windows-data-types
*/

// RawPointerReader reads Windows Data Types using raw pointer reads
type RawPointerReader struct {
	underlyingPtr uintptr
	lBound        uintptr
	rBound        uintptr
}

// unsafePtr is a convenient function to convert the underlying pointer to a unsafe.Pointer
func (c RawPointerReader) unsafePtr() unsafe.Pointer {
	return unsafe.Pointer(c.underlyingPtr)
}

// offsetPtr returns a new RawPointerReader with the offset of the size of a pointer
func (c RawPointerReader) offsetPtr() RawPointerReader {
	return c.offset(pSize)
}

// offset8 is a convenient function to offset the reader by 1 byte
func (c RawPointerReader) offset8() RawPointerReader {
	return c.offset(1)
}

// offset16 is a convenient function to offset the reader by 2 byte
func (c RawPointerReader) offset16() RawPointerReader {
	return c.offset(2)
}

// offset32 is a convenient function to offset the reader by 4 byte
func (c RawPointerReader) offset32() RawPointerReader {
	return c.offset(4)
}

// offset64 is a convenient function to offset the reader by 8 byte
func (c RawPointerReader) offset64() RawPointerReader {
	return c.offset(8)
}

// offset returns a new RawPointerReader with x bytes of offset, reading boundaries remain unchanged
func (c RawPointerReader) offset(x int) RawPointerReader {
	return RawPointerReader{c.underlyingPtr + uintptr(x), c.lBound, c.rBound}
}

// Byte reads a BYTE
func (c RawPointerReader) Byte() (ret byte, next RawPointerReader) {
	return c.Char()
}

// Bytes reads a series of bytes of a fixed length
func (c RawPointerReader) Bytes(length int) (ret []byte, next RawPointerReader) {
	ret = make([]byte, length)
	for i := 0; i < length; i++ {
		ret[i], c = c.Byte()
	}
	next = c
	return
}

// PChar reads a PCHAR into a string
func (c RawPointerReader) PChar() (ret string, next RawPointerReader) {
	next = c.offsetPtr()

	r := bytes.NewBuffer([]byte{})
	c, next = c.Deref()
	var b byte
	for {
		b, c = c.Char()
		if b == 0 {
			break
		}
		r.WriteByte(b)
	}
	return r.String(), next
}

// Char reads a CHAR
func (c RawPointerReader) Char() (ret byte, next RawPointerReader) {
	return *(*byte)(c.unsafePtr()), c.offset8()
}

// PWChar reads a PWCHAR into a string
func (c RawPointerReader) PWChar() (ret string, next RawPointerReader) {
	next = c.offsetPtr()

	r := bytes.NewBuffer([]byte{})
	c, _ = c.Deref()
	var b rune
	for {
		b, c = c.WChar()
		if b == 0 {
			break
		}
		r.WriteRune(b)
	}
	return r.String(), next
}

// WChar reads a WCHAR
func (c RawPointerReader) WChar() (ret rune, next RawPointerReader) {
	val, next := c.Word()
	return rune(val), next
}

// Int reads a INT
func (c RawPointerReader) Int() (ret int32, next RawPointerReader) {
	return *(*int32)(c.unsafePtr()), c.offset32()
}

// Ushort reads a USHORT
func (c RawPointerReader) UShort() (ret uint16, next RawPointerReader) {
	return c.Word()
}

// Word reads a WORD
func (c RawPointerReader) Word() (ret uint16, next RawPointerReader) {
	return *(*uint16)(c.unsafePtr()), c.offset16()
}

// DWord reads a DWORD
func (c RawPointerReader) DWord() (ret uint32, next RawPointerReader) {
	return c.ULong()
}

// ULong reads a ULONG
func (c RawPointerReader) ULong() (ret uint32, next RawPointerReader) {
	return *(*uint32)(c.unsafePtr()), c.offset32()
}

// ULong64 reads a ULONG64
func (c RawPointerReader) ULong64() (ret uint64, next RawPointerReader) {
	return c.ULongLong()
}

// ULongLong reads a ULONGLONG
func (c RawPointerReader) ULongLong() (ret uint64, next RawPointerReader) {
	return *(*uint64)(c.unsafePtr()), c.offset64()
}

// Deref interprets the current pointer as a pointer value and dereferences to a new RawPointerReader, reading boundaries remain unchanged
// the returned reader might not be valid, so check for validity using IsValid before accessing them, or an unrecoverable panic would occur
func (c RawPointerReader) Deref() (ret RawPointerReader, next RawPointerReader) {
	ptr := *((*uintptr)(c.unsafePtr()))
	return RawPointerReader{ptr, c.lBound, c.rBound}, c.offsetPtr()
}

// IsValid checks if the current pointer position is in the reading boundaries
func (c RawPointerReader) IsValid() bool {
	return !(c.lBound > c.underlyingPtr || c.rBound < c.underlyingPtr)
}

// IsNilPointer is a convenient function to check if the memory address referenced at the current pointer is in the reading boundaries
func (c RawPointerReader) IsNilPointer() bool {
	ptr, _ := c.Deref()
	return !ptr.IsValid()
}

// Ptr returns the underlying pointer of the RawPointerReader as an unsafe.Pointer
func (c RawPointerReader) Ptr() (ret unsafe.Pointer, next RawPointerReader) {
	return c.unsafePtr(), c.offsetPtr()
}

// New creates a RawPointerReader of using a pointer to an existing type
func New(ptr uintptr, size int) RawPointerReader {
	return RawPointerReader{ptr, ptr, ptr + uintptr(size-1)}
}

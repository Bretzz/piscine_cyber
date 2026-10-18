package exif

import (
	"encoding/binary"
	"fmt"
	"math"
)

// type ExifDataType = uint16

type ProcessFunc = func([]byte, binary.ByteOrder) any

type DataType struct {
	ID       uint16 // uint id found in IFD entries
	ByteSize uint32 // size in bytes
	// RawType  reflect.Type
	Name    string // welp..
	process ProcessFunc
}

func (t DataType) toString() string {
	return t.Name
}

// PROCESSING FUNCTIONS
// ------------------------------------------
func processBYTE(b []byte, endian binary.ByteOrder) any {
	if len(b) == 0 {
		return byte(0)
	}
	return b[0]
}

func processASCII(b []byte, endian binary.ByteOrder) any {
	length := len(b)
	for i, v := range b {
		if v == 0 {
			length = i
			break
		}
	}
	return string(b[:length])
}

func processSHORT(b []byte, endian binary.ByteOrder) any {
	return endian.Uint16(b[:2])
}

func processLONG(b []byte, endian binary.ByteOrder) any {
	return endian.Uint32(b[:4])
}

func processRATIONAL(b []byte, endian binary.ByteOrder) any {
	numerator := float32(endian.Uint32(b[:4]))
	denominator := float32(endian.Uint32(b[4:8]))
	if denominator == 0 {
		return float32(0)
	}
	return numerator / denominator
}

func processSBYTE(b []byte, endian binary.ByteOrder) any {
	if len(b) == 0 {
		return int8(0)
	}
	return int8(b[0])
}

func processUNDEFINED(b []byte, endian binary.ByteOrder) any {
	if len(b) == 0 {
		return byte(0)
	}
	return b[0]
}

func processSSHORT(b []byte, endian binary.ByteOrder) any {
	return int16(endian.Uint16(b[:2]))
}

func processSLONG(b []byte, endian binary.ByteOrder) any {
	return int32(endian.Uint32(b[:4]))
}

func processSRATIONAL(b []byte, endian binary.ByteOrder) any {
	numerator := float32(int32(endian.Uint32(b[:4])))
	denominator := float32(int32(endian.Uint32(b[4:8])))
	if denominator == 0 {
		return float32(0)
	}
	return numerator / denominator
}

// processFLOAT parses a 4-byte IEEE-754 single-precision float
func processFLOAT(b []byte, endian binary.ByteOrder) any {
	bits := endian.Uint32(b[:4])
	return math.Float32frombits(bits)
}

// processDOUBLE parses an 8-byte IEEE-754 double-precision float
func processDOUBLE(b []byte, endian binary.ByteOrder) any {
	bits := endian.Uint64(b[:8])
	return math.Float64frombits(bits)
}

func processIFD(b []byte, endian binary.ByteOrder) any {
	return endian.Uint32(b[:4])
}

// PARSE
// ----------------------

func (t *DataType) parse(id uint16) error {
	t.ID = id
	switch id {
	case 0x001:
		t.ByteSize = 1
		t.Name = "BYTE"
		t.process = processBYTE
	case 0x002:
		t.ByteSize = 1
		t.Name = "ASCII string"
		t.process = processASCII
	case 0x003:
		t.ByteSize = 2
		t.Name = "SHORT"
		t.process = processSHORT
	case 0x004:
		t.ByteSize = 4
		t.Name = "LONG"
		t.process = processLONG
	case 0x005:
		t.ByteSize = 8
		t.Name = "RATIONAL"
		t.process = processRATIONAL
	case 0x006:
		t.ByteSize = 1
		t.Name = "SBYTE"
		t.process = processSBYTE
	case 0x007:
		t.ByteSize = 1
		t.Name = "UNDEFINED"
		t.process = processUNDEFINED
	case 0x008:
		t.ByteSize = 2
		t.Name = "SSHORT"
		t.process = processSSHORT
	case 0x009:
		t.ByteSize = 4
		t.Name = "SLONG"
		t.process = processSLONG
	case 0x00A:
		t.ByteSize = 8
		t.Name = "SRATIONAL"
		t.process = processSRATIONAL
	case 0x00B:
		t.ByteSize = 4
		t.Name = "FLOAT"
		t.process = processFLOAT
	case 0x00C:
		t.ByteSize = 8
		t.Name = "DOUBLE"
		t.process = processDOUBLE
	case 0x00D:
		t.ByteSize = 4
		t.Name = "IFD"
		t.process = processIFD
	default:
		return fmt.Errorf("Invalid DataType %d", id)
	}
	return nil
}

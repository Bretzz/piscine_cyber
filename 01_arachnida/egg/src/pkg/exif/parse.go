package exif

import (
	"encoding/binary"
	"fmt"
	"log"

	"scorpion/pkg/find"
)

// type IFDEntryValue struct {
// 	size uint // size of the
// 	vtype
// 	data []any	// actual data
// 	byte
// }

// Single entry for an IFD
type IFDEntry struct {
	tag      uint16   // Identifies the specific property (e.g., 0x0110 for camera model, 0x829A for exposure time)
	dataType DataType // Specifies format (e.g., BYTE, ASCII string, SHORT, LONG, RATIONAL)
	count    uint32   // The number of items of the specified type
	value    uint32   // Either the Value or the offset to the value (if it's bigger than 32bits)
	pvalue   []any    // parsed value
	bytes    []byte   // Pointer to the raw entry bytes
}

type IFD struct {
	length  uint       // How many entries
	entries []IFDEntry // Array of entries
	next    uint       // Offset of th next IFD
	bytes   []byte     // Pointer to the raw bytes of the IFD struct
}

type EXIF struct {
	offset    uint             // Offset of the "Exif\0\0" tag
	zero      uint             // Where the offset count from
	endian    binary.ByteOrder // Little or Big endian byte converter
	corrupted bool             // Result of the check for the '42' byte
	ifd0      uint             // offset of the first IFD
	bytes     []byte           // Pointer to the raw bytes of the EXIF chunk (starting from ZERO)
}

// APP1 data struct
type APP1 struct {
	offset uint   // offset from the start of the file
	length uint   // length of the APP1 byte chunk
	bytes  []byte // Pointer to the raw bytes of the APP1 chunk
}

// PARSING FUNCTIONS
// -------------------------

// #prod add multiple file standards
func (app1 *APP1) parseJPEG(b []byte) error {
	app1_offset, err := find.ByteSequence(b, []uint8{0xff, 0xe1})
	if err != nil {
		log.Fatal("Error while searching APP1 marker: ", err)
	}

	// Extract length (Always BIG ENDIAN. as per JPEG standard)
	app1_length := binary.BigEndian.Uint16(b[app1_offset+2 : app1_offset+4])

	// Assign to the struct
	app1.offset = app1_offset
	app1.length = uint(app1_length)
	app1.bytes = b[app1.offset : app1.offset+app1.length]

	return nil
}

func (exif *EXIF) parse(app1 APP1) error {
	b := app1.bytes

	// Find the "Exif\0\0" delimiter
	exif_offset, err := find.ByteSequence(b, []uint8{0x45, 0x78, 0x69, 0x66, 0x00, 0x00})
	if err != nil {
		return fmt.Errorf("Error while searching EXIF header: %s\n", err)
	}

	// where all the offset in the TIFF header count from
	zero_offset := exif_offset + 6

	// Get the Byte order
	little, err := find.IsLittle(b[zero_offset : zero_offset+2])
	if err != nil {
		fmt.Print("[Warn] ", err.Error(), "[Warn] Default: LITTLE\n")
		little = true
	}

	// Save the endian
	var endian binary.ByteOrder
	if little == true {
		endian = binary.LittleEndian
	} else {
		endian = binary.BigEndian
	}

	// Check TIFF fixed bytes
	fourthytwo := endian.Uint16(b[zero_offset+2 : zero_offset+4])
	if fourthytwo != 42 {
		fmt.Print("[Warn] Corrupted TIFF fixed value")
	}

	// Get the IFD0 offset
	first_offset := endian.Uint32(b[zero_offset+4 : zero_offset+8])

	// Assign the values
	exif.offset = exif_offset
	exif.zero = zero_offset
	exif.endian = endian
	exif.corrupted = fourthytwo != 42
	exif.ifd0 = uint(first_offset)
	exif.bytes = b[exif.zero:]

	return nil
}

/*
 * (2 bytes) Tag ID
 * (2 bytes) Data Type
 * (4 bytes) Count
 * (4 bytes) Value / Offset
 */
func (e *IFDEntry) parse(b []byte, offset uint, endian binary.ByteOrder) error {
	// Tag ID
	tag := endian.Uint16(b[offset : offset+2])
	// Data Type
	typeId := endian.Uint16(b[offset+2 : offset+4])
	// Count
	count := endian.Uint32(b[offset+4 : offset+8])
	// Value
	value := endian.Uint32(b[offset+8 : offset+12])

	// Assign values
	e.tag = tag
	e.dataType.parse(typeId)
	e.count = count
	e.value = value
	e.bytes = b

	// Parse the value
	err := e.parseValue(b, offset, endian)
	return err
}

// parse the value after the first metadata parsing
func (e *IFDEntry) parseValue(b []byte, offset uint, endian binary.ByteOrder) error {
	bsize := uint32(e.dataType.ByteSize)
	size := bsize * e.count

	// const Uint32_Max uint32 = 2147483648 // (2^32) - 1

	if size <= 4 { // read from there
		of := uint32(offset) + uint32(8)
		for range e.count {
			v := e.dataType.process(b[of:of+bsize], endian)
			e.pvalue = append(e.pvalue, v)
			of += bsize
		}
	} else { // read from pointer
		// e.pvalue = append(e.pvalue, "offset")
		of := uint32(e.value)
		for range e.count {
			// fmt.Printf("COUNT %d, OFFSET %d\n", e.count, of)
			v := e.dataType.process(b[of:of+bsize], endian)
			e.pvalue = append(e.pvalue, v)
			of += bsize
		}
	}

	return nil
}

// #todo pass offset to error msg?
func (ifd *IFD) parse(b []byte, offset uint, endian binary.ByteOrder) error {

	// #todo maybe check output?
	length := uint(endian.Uint16(b[offset : offset+2]))
	var entries []IFDEntry

	// Parse evry entry
	for i := range length {
		of := 2 + i*12
		var entry IFDEntry

		// check if the slice is fucked
		if int(of) > len(b) || int(of+12) > len(b) {
			break
		}

		// slice
		err := entry.parse(b, offset+of, endian)

		if err != nil {
			fmt.Printf("[Warn] Error while parsing IFD entry (offset %d)\n", of)
		} else {
			entries = append(entries, entry)
		}
	}

	// Calculat pointer to next IFD
	next_offset := offset + 2 + length*12
	next := endian.Uint32(b[next_offset : next_offset+4])

	// Assign values
	ifd.length = uint(length)
	ifd.entries = entries
	ifd.next = uint(next)
	ifd.bytes = b

	return nil
}

// -------------------------

func slice128kb(bytes []byte) []byte {
	const toslice = 128 * 1024
	if len(bytes) < toslice {
		return bytes
	}
	return bytes[:toslice]
}

var app1 APP1
var exif EXIF
var ifds []IFD

var parsed bool = false

// Parse the file
func Parse(b []byte) error {

	// clear stuff
	if parsed != false {
		app1 = APP1{}
		exif = EXIF{}
		ifds = []IFD{}
		parsed = false
	}

	// Slice only 128kb of data
	b = slice128kb(b)

	err := app1.parseJPEG(b)
	if err != nil {
		return err
	}

	err = exif.parse(app1)
	if err != nil {
		return err
	}

	next := exif.ifd0
	// loop all the IFDs
	for {
		// fmt.Println("next ", next)
		if next == 0 {
			break
		}
		var ifd IFD

		err = ifd.parse(exif.bytes, next, exif.endian)
		if err != nil {
			return err
		}

		next = ifd.next
		ifds = append(ifds, ifd)
	}

	parsed = true
	return nil
}

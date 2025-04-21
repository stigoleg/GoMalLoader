//go:build windows

package utils

import (
	"encoding/binary"
	"log"
	"testing"
	"unsafe"
)

type IMAGE_DOS_HEADER struct {
	E_magic    uint16
	E_cblp     uint16
	E_cp       uint16
	E_crlc     uint16
	E_cparhdr  uint16
	E_minalloc uint16
	E_maxalloc uint16
	E_ss       uint16
	E_sp       uint16
	E_csum     uint16
	E_ip       uint16
	E_cs       uint16
	E_lfarlc   uint16
	E_ovno     uint16
	E_res      [4]uint16
	E_oemid    uint16
	E_oeminfo  uint16
	E_res2     [10]uint16
	E_lfanew   int32
}

func ParsePEAndLoad(dll []byte) uintptr {
	dos := *(*IMAGE_DOS_HEADER)(unsafe.Pointer(&dll[0]))
	if dos.E_magic != 0x5A4D {
		log.Fatal("Invalid DOS header")
	}

	ntHeadersOffset := dos.E_lfanew
	optHeaderOffset := ntHeadersOffset + 24 // after PE Signature + FileHeader
	entryOffset := optHeaderOffset + 16     // AddressOfEntryPoint is offset 16 into Optional Header

	entryRVA := binary.LittleEndian.Uint32(dll[entryOffset : entryOffset+4])

	mem := NewMemoryOps()
	addr := mem.AllocRWX(len(dll))
	for i := 0; i < len(dll); i++ {
		*(*byte)(unsafe.Pointer(addr + uintptr(i))) = dll[i]
	}

	return addr + uintptr(entryRVA)
}

func TestParsePEAndLoad_InvalidDOSHeader(t *testing.T) {
	invalid := make([]byte, 64)
	// Not a valid MZ header
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid DOS header")
		}
	}()
	ParsePEAndLoad(invalid)
}

func TestParsePEAndLoad_ValidMinimalPE(t *testing.T) {
	// Minimal valid PE: MZ header, e_lfanew at 0x3C, PE signature, minimal optional header
	pe := make([]byte, 128)
	pe[0] = 'M'
	pe[1] = 'Z'
	pe[0x3C] = 0x40 // e_lfanew = 0x40
	pe[0x40] = 'P'
	pe[0x41] = 'E'
	pe[0x42] = 0
	pe[0x43] = 0
	// Set AddressOfEntryPoint at offset 0x54 (0x40 + 24 + 16)
	pe[0x54] = 0x10
	pe[0x55] = 0x00
	pe[0x56] = 0x00
	pe[0x57] = 0x00
	// Should not panic, but returns an address
	addr := ParsePEAndLoad(pe)
	if addr == 0 {
		t.Error("expected non-zero address for valid PE")
	}
}

func TestParsePEAndLoad_TruncatedPE(t *testing.T) {
	truncated := make([]byte, 10) // Too short for any header
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for truncated PE file")
		}
	}()
	ParsePEAndLoad(truncated)
}

func TestParsePEAndLoad_MissingEntryPoint(t *testing.T) {
	pe := make([]byte, 128)
	pe[0] = 'M'
	pe[1] = 'Z'
	pe[0x3C] = 0x40 // e_lfanew = 0x40
	pe[0x40] = 'P'
	pe[0x41] = 'E'
	pe[0x42] = 0
	pe[0x43] = 0
	// Do not set AddressOfEntryPoint (leave as zero)
	addr := ParsePEAndLoad(pe)
	if addr == 0 {
		t.Error("expected non-zero address even with missing entry point (should default to base)")
	}
}

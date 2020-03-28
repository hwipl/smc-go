package clc

import "fmt"

// RMBESize stores the SMC RMBE size
type RMBESize uint8

// String converts rmbeSize to a string
func (s RMBESize) String() string {
	size := 1 << (s + 14)
	return fmt.Sprintf("%d (%d)", s, size)
}

package clc

import "fmt"

// rmbeSize stores the SMC RMBE size
type rmbeSize uint8

// String converts rmbeSize to a string
func (s rmbeSize) String() string {
	size := 1 << (s + 14)
	return fmt.Sprintf("%d (%d)", s, size)
}

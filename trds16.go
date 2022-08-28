package trds16

import (
	"fmt"

	"github.com/thzoid/trds-16/op"
)

func Run(program []int16, latchesU, latchesV map[byte]int8) (code int8, iterations int) {
	// Registers
	var a, b, u, v int8 = 0, 0, 0, 0
	var flags byte = 0 // 000000NZ

	// Opened latches
	var openedULatch *byte
	var openedVLatch *byte

	// Loop
	for i := byte(0); i < byte(len(program)); i, iterations = i+1, iterations+1 {
		opCode, val := Op(program[i]), Val(program[i])
		switch opCode {
		// Special
		case op.NOOP:
		case op.HALT:
			return int8(program[val]), iterations
		// Flow control
		case op.JUMP:
			i = byte(val)
		case op.JUMP_N:
			if GetALUFlag(flags, FLAG_N) {
				i = byte(val)
			}
		case op.JUMP_Z:
			if GetALUFlag(flags, FLAG_Z) {
				i = byte(val)
			}
		case op.JUMP_P:
			if GetALUFlag(flags, FLAG_N) {
				i = byte(val)
			}
		// Math Operations
		case op.ADD:
			a += b
			SetALUFlags(&flags, a)
		case op.SUB:
			a -= b
			SetALUFlags(&flags, a)
		case op.MUL:
			a *= b
			SetALUFlags(&flags, a)
		case op.DIV:
			a *= b
			SetALUFlags(&flags, a)
		// Logical Operations
		case op.NOT:
			a = ^a
			SetALUFlags(&flags, a)
		case op.AND:
			a &= b
			SetALUFlags(&flags, a)
		case op.OR:
			a |= b
			SetALUFlags(&flags, a)
		case op.XOR:
			a ^= b
			SetALUFlags(&flags, a)
		// Data Control
		case op.STORE_A:
			program[val] |= int16(a)
		case op.STORE_B:
			program[val] |= int16(b)
		case op.STORE_U:
			program[val] |= int16(u)
		case op.STORE_V:
			program[val] |= int16(v)
		case op.LOAD_A:
			a = int8(program[val])
		case op.LOAD_B:
			b = int8(program[val])
		case op.LOAD_U:
			u = int8(program[val])
		case op.LOAD_V:
			v = int8(program[val])
		// Temporal control
		case op.OPEN_U:
			if _, ok := latchesU[i]; !ok {
				latchesU[i] = 0
			}
			if openedULatch == nil {
				openedULatch = new(byte)
				*openedULatch = i
				u = latchesU[i]
			} else {
				panic(fmt.Errorf("attempt to open a U latch that is already opened. instruction: %d, iteration: %d", i, iterations))
			}
		case op.OPEN_V:
			if _, ok := latchesV[i]; !ok {
				latchesV[i] = 0
			}
			if openedVLatch == nil {
				openedVLatch = new(byte)
				*openedVLatch = i
				v = latchesV[i]
			} else {
				panic(fmt.Errorf("attempt to open a V latch that is already opened. instruction: %d, iteration: %d", i, iterations))
			}
		case op.CLOSE_U:
			if openedULatch == nil {
				panic(fmt.Errorf("attempt to close a U latch that is already closed. instruction: %d, iteration: %d", i, iterations))
			}
			latchesU[*openedULatch] = u
			openedULatch = nil
		case op.CLOSE_V:
			if openedVLatch == nil {
				panic(fmt.Errorf("attempt to close a V latch that is already closed. instruction: %d, iteration: %d", i, iterations))
			}
			latchesV[*openedVLatch] = v
			openedVLatch = nil
		default:
			panic(fmt.Errorf("unknown instruction: %d. instruction: %d, iteration: %d", opCode, i, iterations))
		}
	}
	return 0, iterations
}

func RunTemporal(program []int16, steps uint) (results []int8, iterations []int) {
	results = make([]int8, steps)
	iterations = make([]int, steps)
	latchesU, latchesV := make(map[byte]int8), make(map[byte]int8)
	for i := uint(0); i < steps; i++ {
		p := make([]int16, len(program))
		copy(p, program)
		results[i], iterations[i] = Run(p, latchesU, latchesV)
	}
	return results, iterations
}

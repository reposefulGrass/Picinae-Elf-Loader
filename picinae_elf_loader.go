package main

import (
	"debug/elf"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: picinae_byte_loader <path of executable> <number of opcodes>")
		os.Exit(1)
	}

	filepath := os.Args[1]
	num_opcodes, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Usage: picinae_byte_loader <path of executable> <number of opcodes>")
		os.Exit(1)
	}

	file, err := elf.Open(filepath)
	if err != nil {
		panic(err)
	}

	bytes := grab_main_bytes(file)
	opcodes := extend_to_opcodes(bytes)
	format_into_program(opcodes, num_opcodes)
}

func grab_main_bytes(file *elf.File) []byte {
	var bytes []byte

	for _, section := range file.Sections {
		if (*section).SectionHeader.Name == ".text" {
			bytes, _ = (*section).Data()
		}
	}

	return bytes
}

func extend_to_opcodes(bytes []byte) []uint32 {
	var opcodes []uint32

	for i := 0; i < len(bytes)/4; i += 4 {
		var opcode uint32

		opcode = opcode | uint32(bytes[i+0])<<24
		opcode = opcode | uint32(bytes[i+1])<<16
		opcode = opcode | uint32(bytes[i+2])<<8
		opcode = opcode | uint32(bytes[i+3])<<0

		opcodes = append(opcodes, opcode)
	}

	return opcodes
}

func format_into_program(opcodes []uint32, num_opcodes int) {
	fmt.Println("Definition program offset : N := ")
	fmt.Println("    match offset with")

	for i, opcode := range opcodes {
		fmt.Println("    |", i*4, "=>", opcode)
		if i == num_opcodes {
			break
		}
	}

	fmt.Println("    | _ => 0")
	fmt.Println("    end.")
}

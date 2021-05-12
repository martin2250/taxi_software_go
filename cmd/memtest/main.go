package main

import (
	"fmt"
	"log"

	"periph.io/x/host/v3/pmem"
)

const TEST_1 = 0x1faaa / 2
const TEST_2 = 0x1f554 / 2

func main() {
	var fpga *[131072 / 2]uint16
	if err := pmem.MapAsPOD(0x30000000, &fpga); err != nil {
		log.Fatal(err)
	}
	errors := 0
	for i := 0; i < 5000; i += 1 {
		fpga[TEST_1] = 0xaaaa
		fpga[TEST_2] = 0x5555
		if fpga[TEST_1] != 0xaaaa {
			errors += 1
		}
		if fpga[TEST_2] != 0x5555 {
			errors += 1
		}
		fpga[TEST_1] = 0x5555
		fpga[TEST_2] = 0xaaaa
		if fpga[TEST_1] != 0x5555 {
			errors += 1
		}
		if fpga[TEST_2] != 0xaaaa {
			errors += 1
		}
	}
	fmt.Printf("%d errors\n", errors)
}

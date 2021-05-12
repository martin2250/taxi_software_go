package main

import (
	"flag"
	"fmt"
	"log"
	"unsafe"

	"github.com/martin2250/taxi_software_go/pkg/fpga"
)

func main() {
	rw := flag.String("rw", "r", "a string")
	reg := flag.Int("reg", 0, "")
	value := flag.Int("val", 0, "")
	flag.Parse()
	fpga, err := fpga.OpenFpga()
	if err != nil {
		log.Fatal(err)
	}
	ptr := unsafe.Pointer(fpga)
	regs := (*[2048 * 1024]uint16)(ptr)
	if *rw == "r" {
		fmt.Printf("regs[0x%04X] = %04X\n", *reg, regs[(*reg)/2])
	} else {
		regs[(*reg)/2] = uint16(*value)
	}
	fpga.WhiteRabbit.CounterPeriod = 25
}

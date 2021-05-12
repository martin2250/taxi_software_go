package main

import (
	"fmt"
	"log"

	"periph.io/x/host/v3/pmem"
)

const OFFSET_PMC = 0xFFFFFC00
const PMC_SCER_PCK1 = (1 << 9)
const PMC_SCER = 0x0000 / 4
const PMC_SCDR = 0x0004 / 4
const PMC_SCSR = 0x0008 / 4
const PMC_PCK1 = 0x0044 / 4
const PMC_SR = 0x0068 / 4

const OFFSET_PIO = 0xFFFFF400
const PIO_PDR = 0x0004 / 4
const PIO_OER = 0x0010 / 4
const PIO_BSR = 0x0074 / 4

func main() {
	var pmc *[128]uint32
	var piob *[128]uint32
	if err := pmem.MapAsPOD(OFFSET_PMC, &pmc); err != nil {
		log.Fatal(err)
	}
	if err := pmem.MapAsPOD(OFFSET_PMC, &piob); err != nil {
		log.Fatal(err)
	}
	// CSS = 0b01
	// PRES = 0
	pmc[PMC_PCK1] = 0x02
	pmc[PMC_SCER] = PMC_SCER_PCK1
	fmt.Printf("PMC_SCSR = %08X\n", pmc[PMC_SCSR])
	fmt.Printf("PMC_SR = %08X\n", pmc[PMC_SR])

	piob[PIO_PDR] = (1 << 31) // disable PIO -> enable perpheral
	piob[PIO_BSR] = (1 << 31) // select peripheral b
	piob[PIO_OER] = (1 << 31) // enable output
}

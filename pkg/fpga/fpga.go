package fpga

import (
	"periph.io/x/host/v3/pmem"
)

//go:generate python3 generate_struct.py

const SMC_MEM_START = 0x30000000
const SMC_MEM_LEN = 0x4000000 // 64 M

func OpenFpga() (*Fpga, error) {
	var fpga *Fpga
	if err := pmem.MapAsPOD(SMC_MEM_START, &fpga); err != nil {
		return &Fpga{}, err
	}
	return fpga, nil
}

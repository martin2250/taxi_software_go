#!/usr/bin/python
from dataclasses import dataclass
from typing import List, Union
import subprocess

@dataclass
class Register:
    offset: int
    name: str
    comment: str = ''
    size: int = 2

@dataclass
class Section:
    name: str
    registers: List[Register]

info: List[Union[Register, Section]] = [
    Section('TestRegisters', [
        Register(0x0004, 'Test5'),
        Register(0x0006, 'TestA'),
    ]),
    Section('BuildInfo', [
        Register(0x0008, 'Year'),
        Register(0x000a, 'Month'),
        Register(0x000c, 'Day'),
        Register(0x000e, 'Version'),
    ]),
    Register(0x0010, 'Modus'),
    Register(0x0012, 'DummyCnt'),
    Section('EventFifo', [
        Register(0x0100, 'DmaBuffer'),
        Register(0x0104, 'FifoWordsAligned'),
        Register(0x0106, 'FifoWordsPerSlice'),
        Register(0x0108, 'IrqStall'),
        Register(0x010a, 'DeviceId'),
    ]),
    Section('GpsTiming', [
        Register(0x0200, 'CounterPeriod'),
        Register(0x0202, 'NewDataLatched'),
        Register(0x0204, 'DiffGpsToLocal'),
        Register(0x0206, 'Week'),
        Register(0x0208, 'QuantizationErrorM'),
        Register(0x020a, 'QuantizationErrorL'),
        Register(0x020c, 'TowMsM'),
        Register(0x020e, 'TowMsL'),
    ]),
    Section('Tmp05', [
        Register(0x0300, 'StartBusy'),
        Register(0x0302, 'TempL'),
        Register(0x0304, 'TempH'),
        Register(0x0306, 'DebugCtrL'),
        Register(0x0308, 'DebugCtrH'),
    ]),
    Section('WhiteRabbit', [
        Register(0x1400, 'NewDataLatched'),
        Register(0x1402, 'IrigDataA'),
        Register(0x1404, 'IrigDataB'),
        Register(0x1406, 'IrigDataC'),
        Register(0x1408, 'IrigDataD'),
        Register(0x140a, 'BinarySeconds'),
        Register(0x140c, 'BinaryDays'),
        Register(0x140e, 'BinaryYears'),
        Register(0x1410, 'RawIrigA'),
        Register(0x1412, 'RawIrigB'),
        Register(0x1414, 'RawIrigC'),
        Register(0x1416, 'RawIrigD'),
        Register(0x1418, 'RawIrigE'),
        Register(0x141a, 'RawIrigF'),
        Register(0x141c, 'BitCounter'),
        Register(0x141e, 'ErrorCounter'),
        Register(0x1420, 'CounterPeriod'),
    ]),
]

class StructWriter:
    def __init__(self, name: str):
        self.current_offset = 0
        self.lines: List[str] = []
        self.name = name
    
    def ensure_offset(self, reg: Register):
        # check offset
        offset_diff = reg.offset-self.current_offset
        if offset_diff < 0:
            raise ValueError(f'register {reg.name}[0x{reg.offset:04X}] declared at position 0x{self.current_offset:04X}')
        if offset_diff > 0:
            self.lines.append(f'_ [0x{offset_diff:X}]uint8')
        self.current_offset = reg.offset

    def add_register(self, reg: Register):
        self.ensure_offset(reg)
        self.current_offset = reg.offset + 2
        line = f'{reg.name} uint16'
        if reg.comment:
            line += f' // {reg.comment}'
        self.lines.append(line)

    def add_section(self, section: Section):
        self.ensure_offset(section.registers[0])
        self.lines.append(f'{section.name} struct {{')
        for reg in section.registers:
            self.add_register(reg)
        self.lines.append('}')

    def print(self, file):
        print(f'type {self.name} struct {{', file=file)
        for line in self.lines:
            print(line, file=file)
        print('}', file=file)


fpga = StructWriter('Fpga')
for item in info:
    if isinstance(item, Section):
        fpga.add_section(item)
    elif isinstance(item, Register):
        fpga.add_register(item)

import pathlib

filename = pathlib.Path(__file__).parent / 'struct_gen.go'

with open(filename, 'w') as file:
    print('package fpga', file=file)
    fpga.print(file)

print('running go fmt')
subprocess.run(['go', 'fmt', filename])
print('done')

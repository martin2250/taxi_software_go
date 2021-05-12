package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/martin2250/taxi_software_go/pkg/fpga"
)

// https://stackoverflow.com/questions/48906483/how-to-set-systems-date-and-time-via-go-syscall
func SetSystemDate(newTime time.Time) error {
	_, lookErr := exec.LookPath("date")
	if lookErr != nil {
		fmt.Printf("Date binary not found, cannot set system date: %s\n", lookErr.Error())
		return lookErr
	} else {
		dateString := newTime.Format("2006-01-2 15:04:05")
		// dateString := newTime.Format("2 Jan 2006 15:04:05")
		fmt.Printf("Setting system date to: %s\n", dateString)
		args := []string{"--set", dateString}
		return exec.Command("date", args...).Run()
	}
}

func main() {
	fpga, err := fpga.OpenFpga()
	if err != nil {
		log.Fatal(err)
	}
	// FIXME: replace the entire wr module inside the FPGA and check version to select correct go code
	// read entire memory region at once -> new data is latched
	wr := fpga.WhiteRabbit
	retries := 5
	for {
		retries--
		if retries == 0 {
			log.Fatal("retry limit reached")
		}
		// separate days and seconds
		days := int32(wr.BinaryDays & 0x1FF)
		seconds := (int32(wr.BinaryDays&0x8000) << 1) | int32(wr.BinarySeconds)
		if wr.BinaryYears < 20 || wr.BinaryYears > 50 {
			log.Printf("invalid year %d - retrying\n", wr.BinaryYears)
		} else if days > 365 {
			log.Printf("invalid day %d - retrying\n", days)
		} else if seconds >= 24*3600 {
			log.Printf("invalid second %d - retrying\n", seconds)
		} else {
			log.Printf("20%02d - %d %d", wr.BinaryDays, wr.BinaryYears, wr.BinarySeconds)
			log.Printf("20%02d-%d %d", wr.BinaryYears, days, seconds)
			date := time.Date(2000+int(wr.BinaryYears), 1, 0, 0, 0, 0, 0, time.UTC)
			date = date.Add(24 * time.Hour * time.Duration(days))
			date = date.Add(time.Second * time.Duration(seconds))
			if err := SetSystemDate(date); err != nil {
				log.Fatal("error setting system time", err)
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
		wr = fpga.WhiteRabbit
	}
}

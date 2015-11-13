// fa-tail.go
//
// tail(1)-like utility that gets the last record and display the clock element
//
// Copyright 2015 © by Ollivier Robert for the EEC

// Implement the fa-tail application, a FA-aware tail(1) clone.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
	"github.com/keltia/flightaware-go/flightaware"
)

const (
	BSIZE = 8192
	VERSION = "1.0"
)

type FArecord struct {
	Type       string
	AirGround  string
	AltChange  string
	Clock      string
	Gs         string
	Heading    string
	Hexid      string
	Id         string
	Ident      string
	Lat        string
	Lon        string
	Reg        string
	Squawk     string
	UpdateType string
}

var (
	fVerbose bool
	fCount   bool
	fileStat os.FileInfo
)

func main() {
	flag.BoolVar(&fCount, "c", false, "Count records.")
	flag.BoolVar(&fVerbose, "v", false, "Be verbose")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Error: You must specify a file!\n")
		os.Exit(1)
	}
	fn := flag.Arg(0)

	// Get the size to seek near the end
	fileStat, err := os.Stat(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error stat(2) on %s: %v\n", fn, err.Error())
	}

	//Seek and read
	fh, err := os.Open(fn)
	if err != nil {
		fmt.Printf("Unable to open %s\n", fn)
		os.Exit(1)
	}

	// Obviously, if we want the number of records, do not seek
	if !fCount {
		// Go forward fast
		_, err = fh.Seek(fileStat.Size() - BSIZE, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to seek into the file %s at %d\n", fn, fileStat.Size() - BSIZE)
			os.Exit(1)
		}
	}

	// Then we go with the usual scanning thingy
	scanner := bufio.NewScanner(fh)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read %s\n", fn)
		os.Exit(1)
	}

	var (
		lastRecord string
		nbRecords  int64
	)

	for scanner.Scan() {
		lastRecord = scanner.Text()

		// In verbose mode we display the last few records read
		if fVerbose {
			fmt.Printf("%d: %s\n", nbRecords, lastRecord)
		}
		nbRecords++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading: %v", err)
	} else {
		// EOF
		if fVerbose {
			fmt.Printf("Last record: %s\n", lastRecord)
		}
	}

	var (
		lastFA FArecord
	)

	recType := flightaware.GetType(lastRecord)

	if fCount {
		fmt.Printf("%s: records %d size %d bytes\n", fn, nbRecords, fileStat.Size())
	} else {
		fmt.Printf("%s: size %d bytes\n", fn, fileStat.Size())
	}

	lastOne, err := flightaware.DecodeRecord(lastRecord)
	switch recType {
	case "position":
		lastPS := lastOne.(flightaware.FAposition)
		fmt.Printf("Last record is a flightplan for %s (%s):\n",
			lastPS.Ident, lastPS.AircraftType)

		if lastPS.Status == "Z" {
			running, _ := strconv.ParseInt(lastPS.Ete, 10, 64)
			fmt.Printf("  At %s (completed, running time: %d s\n",
				lastPS.Dest, running)
		} else {
			time_edt, _ := strconv.ParseInt(lastPS.Edt, 10, 64)
			time_eta, _ := strconv.ParseInt(lastPS.Eta, 10, 64)
			fmt.Printf("  From %s (%v) to %s (%v)\n",
				lastPS.Orig, time.Unix(time_edt, 0),
				lastPS.Dest, time.Unix(time_eta, 0))
		}
	default:
		iClock, _ := strconv.ParseInt(lastFA.Clock, 10, 64)
		fmt.Printf("Last record: %v\n", time.Unix(iClock, 0))
	}
}

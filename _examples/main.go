package main

import (
	"bytes"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/Arhat109/logger/pkg/logger"
)

// Global vars. Waiting this will be allocated by compiler
var outbuf [1024]byte                       // array! allocating may be by compiler
var outWriter = bytes.NewBuffer(outbuf[:0]) // create slice from array and generate new Writer may be by compiler too

type Mem struct {
	Mallocs uint64
	Frees   uint64
	Objects uint64
	Mems    uint64
	Point   time.Time
}

const Repeats = 1000

// for use into for, fllocate max mem
var memData runtime.MemStats
var memStats = [Repeats*2 + 4]Mem{}
var memLen int

func AddPoint() {
	if memLen >= Repeats*2+4 {
		memLen = 0
	}
	memStats[memLen].Point = time.Now()
	runtime.ReadMemStats(&memData) // read to global!
	memStats[memLen].Mallocs = memData.Mallocs
	memStats[memLen].Frees = memData.Frees
	memStats[memLen].Objects = memData.HeapObjects
	memStats[memLen].Mems = memData.Alloc
	memLen++
}

func main() {
	var err error
	//var lgr = &logger.BaseLogger{Out: outWriter}
	var lgr = &logger.BaseLogger{}
	lgr.Init(&logger.LogConfig{
		Out:    "./example.log",
		IsJson: false,
		Flags:  logger.LogDate,
		Level:  logger.LogDebugLevel,
	}, &err)
	if err != nil {
		return
	}

	runtime.GC()
	AddPoint() // [0]
	for i := 0; i < Repeats; i++ {
		outWriter.Reset()
		//AddPoint()
		lgr.Debug("This a simple message")
		//AddPoint() // [2] for first cycle
	}
	AddPoint()

	// clear out flags and switch for see
	log2 := log.New(os.Stdout, "", 0)

	totalTime := memStats[memLen-1].Point.Sub(memStats[0].Point).Nanoseconds()
	totalMallocs := memStats[memLen-1].Mallocs - memStats[0].Mallocs
	totalFrees := memStats[memLen-1].Frees - memStats[0].Frees
	totalObjects := memStats[memLen-1].Objects - memStats[0].Objects
	totalMems := memStats[memLen-1].Mems - memStats[0].Mems
	log2.Printf("Total time=%d ns", totalTime)
	log2.Printf("Total allocates=%d, frees=%d, objects=%d, mem=%d",
		totalMallocs, totalFrees, totalObjects, totalMems,
	)
	log2.Printf("\nAvg time=%d ns/op", totalTime/Repeats)
	log2.Printf("Avg allocates=%d, frees=%d, objects=%d, mem=%d /op",
		totalMallocs/Repeats, totalFrees/Repeats, totalObjects/Repeats, totalMems/Repeats,
	)
}

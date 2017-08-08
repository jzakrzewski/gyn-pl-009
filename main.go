// Copyright © 2017 Jakub Zakrzewski

package main

import (
	"fmt"
	"os"

	"sync"

	"github.com/mjibson/go-dsp/wav"
	"github.com/r9y9/gossp/f0"
)

const samplingFreq = 44100

func main() {
	if len(os.Args) < 2 {
		panic("Gimme a *.wav file!")
	}

	testWav, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	wr, err := wav.New(testWav)
	if err != nil {
		panic(err)
	}

	samples := make([][]float64, wr.Samples/samplingFreq)

	for sec := range samples {
		s, _ := wr.ReadFloats(samplingFreq)
		samples[sec] = make([]float64, len(s))
		for i, v := range s {
			samples[sec][i] = float64(v)
		}
	}

	freqs := make([]float64, len(samples))
	var wg sync.WaitGroup
	wg.Add(len(freqs))
	for i, s := range samples {
		go func(ind int, smpl []float64) { // speed thing up by doing in parallel
			// work with a shorter sample to speed up the computation and improve accuracy by taking a stable slice
			freqs[ind], _ = f0.NewYIN(samplingFreq).ComputeF0(smpl[samplingFreq/4 : samplingFreq/2])
			wg.Done()
		}(i, s)
	}

	codes := make([]byte, len(freqs)/8)
	wg.Wait()
	for i, f := range freqs {
		codes[i/8] |= byte(f/1000.0) << uint(i%8) // bytes are in-order, bits are reversed
	}

	fmt.Println(string(codes))
}

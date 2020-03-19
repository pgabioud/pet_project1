package main

import (
	"math"
	"math/rand"
	"time"
)

func countMultGate(circuitID CircuitID) uint64 {
	var counter uint64 = 0
	for _, op := range TestCircuits[circuitID-1].Circuit {
		_, ok := op.(*Mult)
		if ok {
			counter++
		}
	}
	return counter
}

func genBeavers(multGateCount uint64) [][3]uint64 {
	var s = int64(math.Pow(2, 16)) + 1
	var beavers [][3]uint64

	var i uint64
	for i = 0; i < multGateCount; i++ {
		rand.Seed(time.Now().UTC().UnixNano() / 10000000)
		var singleBeavers = [3]uint64{uint64(rand.Int63n(s)), uint64(rand.Int63n(s))}
		singleBeavers[2] = uint64(mod(int64(singleBeavers[0]*singleBeavers[1]), s))
		beavers = append(beavers, singleBeavers)
	}

	return beavers
}

func genSharedBeavers(beaverTriplet *[][3]uint64, nbPeers int) [][][3]uint64 {
	var s = int64(math.Pow(2, 16)) + 1
	var sharedBeavers = make([][][3]uint64, len(*beaverTriplet))

	rand.Seed(time.Now().UTC().UnixNano() / 10000000)

	for x := range sharedBeavers {
		sharedBeavers[x] = make([][3]uint64, nbPeers)
	}
	for j := 0; j < len(*beaverTriplet); j++ {
		var totA, totB, totC uint64 = 0, 0, 0
		for i := 0; i < nbPeers-1; i++ {
			var singleSharedBeavers [3]uint64
			singleSharedBeavers[0], singleSharedBeavers[1], singleSharedBeavers[2] = uint64(rand.Int63n(s)), uint64(rand.Int63n(s)), uint64(rand.Int63n(s))
			totA += singleSharedBeavers[0]
			totB += singleSharedBeavers[1]
			totC += singleSharedBeavers[2]
			sharedBeavers[i][j] = singleSharedBeavers

		}
		var singleSharedBeavers [3]uint64
		singleSharedBeavers[0] = uint64(mod(int64((*beaverTriplet)[j][0])-int64(totA), s))
		singleSharedBeavers[1] = uint64(mod(int64((*beaverTriplet)[j][1])-int64(totB), s))
		singleSharedBeavers[2] = uint64(mod(int64((*beaverTriplet)[j][2])-int64(totC), s))
		sharedBeavers[nbPeers-1][j] = singleSharedBeavers
	}
	return sharedBeavers
}

func genAllBeaverTriplets(circuitID CircuitID, nbPeers int) [][][3]uint64 {
	nbTriplets := countMultGate(circuitID)
	beavers := genBeavers(nbTriplets)
	sharedBeavers := genSharedBeavers(&beavers, nbPeers)

	return sharedBeavers
}

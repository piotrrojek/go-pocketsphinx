package pocketsphinx

/*
#cgo pkg-config: pocketsphinx sphinxbase
#include <pocketsphinx.h>
*/
import "C"

type segments struct {
	ps *C.ps_decoder_t
	nb *C.ps_nbest_t
}

//Segment represents a word segment contains frame infomations and probabirity
type Segment struct {
	Word       string
	StartFrame int64
	EndFrame   int64
	Prob       int64
	AcProb     int64
	LmProb     int64
	LbackProb  int64
}

// GetSegments returns word segment list for best hypotesis
func GetSegments(ps *C.ps_decoder_t) []Segment {
	s := segments{ps: ps}
	return s.getBesyHypSegments()
}

// GetSegments returns word segment list for nbest_t
func GetSegmentsForNbest(nb *C.ps_nbest_t) []Segment {
	s := segments{nb: nb}
	return s.getNbestHypSegments()
}

func (s segments) getBesyHypSegments() []Segment {
	var score C.int32
	segIt := C.ps_seg_iter(s.ps, &score)
	return s.getSegmentsFromIter(segIt)
}

func (s segments) getNbestHypSegments() []Segment {
	var score C.int32
	segIt := C.ps_nbest_seg(s.nb, &score)
	return s.getSegmentsFromIter(segIt)
}

func (s segments) getSegmentsFromIter(segIt *C.ps_seg_t) []Segment {
	ret := make([]Segment, 0, 10)
	for {
		if segIt == nil {
			break
		}
		seg := s.getCurrentSegment(segIt)
		ret = append(ret, seg)
		segIt = C.ps_seg_next(segIt)
	}
	return ret
}

func (s segments) getCurrentSegment(segIt *C.ps_seg_t) Segment {
	var start, end C.int
	word := C.GoString(C.ps_seg_word(segIt))
	C.ps_seg_frames(segIt, &start, &end)

	var acousticProb, lmProb, lbackProb C.int32
	segProb := C.ps_seg_prob(segIt, &acousticProb, &lmProb, &lbackProb)

	seg := Segment{word, int64(start), int64(end),
		int64(segProb), int64(acousticProb), int64(lmProb), int64(lbackProb),
	}
	return seg

}

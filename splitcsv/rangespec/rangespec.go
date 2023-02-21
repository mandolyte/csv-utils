package rangespec

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// RangeSpec parses a range specification, such as:
// 1,3,5-8,12-
// It will return a slice of RangeSpec, being two ints,
// a start and a stop
// Ranges start at 1, not zero.
type RangeSpec struct {
	pairs []pair
	spec  string
	Max   uint64
}
type pair struct {
	start, stop uint64
}

// New takes a range specification string and
// returns a slice of RangeSpec structs
func New(r string) (*RangeSpec, error) {
	// remove any whitespace as a convenience
	r = strings.Replace(r, " ", "", -1)
	ret := new(RangeSpec)
	ret.pairs = make([]pair, 0)
	ret.spec = r
	tokens := strings.Split(r, ",")
	for n, val := range tokens {
		//fmt.Printf("Working on %v at %v\n", val, n)
		// does val have a dash?
		if strings.Contains(val, "-") {
			// split on dash
			ends := strings.Split(val, "-")
			if len(ends) > 2 {
				return nil, fmt.Errorf("RangeSpec: malformed specification:%v", val)
			}
			if ends[1] != "" {
				//fmt.Print("ends[] greater than 1\n")
				end1, err := strconv.ParseUint(ends[0], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("RangeSpec: not a number:%v\n%v", ends[0], err)
				}
				end2, err := strconv.ParseUint(ends[1], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("RangeSpec: not a number:%v\n%v", ends[1], err)
				}
				var rs pair
				rs.start = end1
				rs.stop = end2
				ret.pairs = append(ret.pairs, rs)
			} else {
				//fmt.Print("ends[] == 1\n")
				if n+1 != len(tokens) {
					return nil, fmt.Errorf("RangeSpec: open range must be last:%v", val)
				}
				end1, err := strconv.ParseUint(ends[0], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("RangeSpec: not a number:%v\n%v", ends[0], err)
				}
				var rs pair
				rs.start = end1
				rs.stop = math.MaxUint64
				ret.pairs = append(ret.pairs, rs)
			}
			continue
		} else {
			end1, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("RangeSpec: not a number:%v\n%v", val, err)
			}
			var rs pair
			rs.start = end1
			rs.stop = end1
			ret.pairs = append(ret.pairs, rs)
		}
	}
	// ensure ascending specification
	for i := 0; i < len(ret.pairs); i++ {
		if i == 0 {
			if ret.pairs[i].start == 0 {
				return nil, fmt.Errorf("RangeSpec: range must be larger zero: %v", ret.pairs[i].start)
			}
		}
		if ret.pairs[i].start > ret.pairs[i].stop {
			return nil, fmt.Errorf("RangeSpec: start (%v) must be equal or less than stop (%v)", ret.pairs[i].start, ret.pairs[i].stop)
		}
		if i > 0 {
			if ret.pairs[i].start <= ret.pairs[i-1].stop {
				return nil, fmt.Errorf("RangeSpec: start (%v) must be greater than previous stop (%v)", ret.pairs[i].start, ret.pairs[i-1].stop)
			}
		}
	}
	// set the maximum row number
	ret.Max = ret.pairs[len(ret.pairs)-1].stop
	return ret, nil
}

// InRange will test whehter a number is in the range specification
func (rs *RangeSpec) InRange(num uint64) bool {
	for _, val := range rs.pairs {
		if num >= val.start && num <= val.stop {
			return true
		}
	}
	return false
}

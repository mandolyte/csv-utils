package rangespec

import (
	"testing"
	"math"
	"fmt"
)

func TestRangeSpec(t *testing.T) {
  input := " 1, 3, 5 - 8 , 12 - "
  expected := []pair{{1,1},{3,3},{5,8},{12,math.MaxUint64}}
	result,err := New(input)
	if err != nil {
		t.Fatalf("[fail] expected no error for: %s, got %#v\n", input, err)
	}

	for n := range result.pairs {
		if result.pairs[n].start != expected[n].start {
			t.Fatalf("[fail] Start[%v] should be: %v, got %v\n", n, result.pairs[n].start, expected[n].start)
		}
		if result.pairs[n].stop != expected[n].stop {
			t.Fatalf("[fail] Stop[%v] should be: %v, got %v\n", n, result.pairs[n].stop, expected[n].stop)
		}
	}

	input = "1,3,5-,12-"
  expected = nil
	result,err = New(input)
	if err == nil {
		t.Fatalf("[fail] expected error for invalid input: %s\n", input)
	}
	if result != nil {
		t.Fatalf("[fail] expected result to be nil for invalid input: %s\n", input)
	}
	fmt.Printf("Invalid error message for %v was:\n	%v\n\n",input,err)

	input = "1,A,5-,12-"
  expected = nil
	result,err = New(input)
	if err == nil {
		t.Fatalf("[fail] expected error for invalid input: %s\n", input)
	}
	if result != nil {
		t.Fatalf("[fail] expected result to be nil for invalid input: %s\n", input)
	}
	fmt.Printf("Invalid error message for %v was:\n	%v\n\n",input,err)

	input = "1,5,3-,12-"
  expected = nil
	result,err = New(input)
	if err == nil {
		t.Fatalf("[fail] expected error for invalid input: %s\n", input)
	}
	if result != nil {
		t.Fatalf("[fail] expected result to be nil for invalid input: %s\n", input)
	}
	fmt.Printf("Invalid error message for %v was:\n	%v\n\n",input,err)

	input = "1,3,5-4,12-"
  expected = nil
	result,err = New(input)
	if err == nil {
		t.Fatalf("[fail] expected error for invalid input: %s\n", input)
	}
	if result != nil {
		t.Fatalf("[fail] expected result to be nil for invalid input: %s\n", input)
	}
	fmt.Printf("Invalid error message for %v was:\n	%v\n\n",input,err)

	input = "1,5,5-6,12-"
  expected = nil
	result,err = New(input)
	if err == nil {
		t.Fatalf("[fail] expected error for invalid input: %s\n", input)
	}
	if result != nil {
		t.Fatalf("[fail] expected result to be nil for invalid input: %s\n", input)
	}
	fmt.Printf("Invalid error message for %v was:\n	%v\n\n",input,err)

	input = " 1,3,5-8,12-"
	result,err = New(input)
	if err != nil {
		t.Fatalf("[fail] expected no error for: %s, got %#v\n", input, err)
	}

	tests := []uint64{1,2,3,5,6,8,12,13,99}
	exptd := []bool{true,false,true,true,true,true,true,true,true}
	for n,x := range tests {
		if result.InRange(x) != exptd[n] {
			t.Fatalf("[fail] InRange error, range %v, for %v got %v\n", input, x, result.InRange(x))
		}
	}

}

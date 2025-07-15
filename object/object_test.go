package object

import (
	"testing"
)

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}

	// name1 := &object.String{Value: "name"}
	// monkey := &object.String{Value: "Monkey"}
	// pairs := map[object.HashKey]object.Object{}
	//
	// pairs[name1.HashKey()] = monkey
	//
	// fmt.Printf("pairs[name1.HashKey()]=%+v\n", pairs[name1.HashKey()])
	//
	// name2 := &object.String{Value: "name"}
	//
	// fmt.Printf("pairs[name2.HashKey()]=%+v\n", pairs[name2.HashKey()])
	// fmt.Printf("(name1 == name2)=%t\n", name1 == name2)
	// fmt.Printf("(name1.HashKey() == name2.HashKey())=%t\n", name1.HashKey() == name2.HashKey())
}

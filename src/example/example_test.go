package example

import (
	"testing"
	"fmt"
)

func TestExample(t *testing.T) {
	l := make([]Animal, 2)
	l[0] = MakeDog("Henry", 4)
	l[1] = MakeHusky("Harry", 5, 2)
	
	fmt.Println(l[0].Sound())
	fmt.Println(l[1].Sound())
}
package main

import (
	"fmt"

	gobrex "github.com/kujtimiihoxha/go-brace-expansion"
)

func main() {
	permutations := gobrex.Expand("ocp-master-01{a,b,c}.{lab1,rdu1,dfw1,lax2,jfk1}.bandwidthclec.local")
	for _, permutation := range permutations {
		fmt.Println(permutation)
	}
}

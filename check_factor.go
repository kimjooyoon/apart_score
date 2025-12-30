package main

import (
	"fmt"
	"apart_score/pkg/metadata"
)

func main() {
	fmt.Println("층수 팩터 타입:", metadata.FloorLevel.FactorType())
	
	internal := metadata.GetMetadataByFactorType(metadata.FactorInternal)
	external := metadata.GetMetadataByFactorType(metadata.FactorExternal)
	
	internalCount := 0
	for _, mt := range internal {
		if mt != 0 {
			internalCount++
		}
	}
	
	externalCount := 0
	for _, mt := range external {
		if mt != 0 {
			externalCount++
		}
	}
	
	fmt.Printf("Internal: %d, External: %d, Total: %d\n", internalCount, externalCount, internalCount+externalCount)
}

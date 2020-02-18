package a52

import "utils"

type A52_builder struct {

}

type CipherGenerationStrategy int

const (
	KeyBasedGenerationStrategy    CipherGenerationStrategy = 0
	CustomInitialStatesGenerationStrategy  CipherGenerationStrategy = 1
)


func(b A52_builder) Build(key [64] int, generation_strategy CipherGenerationStrategy) *A52{
	
	
	a52 :=  A52{
		Session_key: key,
		R1:          R1{},
		R2:          R2{},
		R3:          R3{},
		R4:          R4{},
		UsedBitsOfCurrentCipherOutput: 0,
		OutputCache:                   nil,
		CurrentFrameNumber:            utils.Frame{},
		InitialStates:                 utils.InitialStates{},
	}

	if generation_strategy == KeyBasedGenerationStrategy{
		a52.NextCipherbitStragegy = &KeyBasedNextCipherbitStrategy{&a52}
	}
	if generation_strategy == CustomInitialStatesGenerationStrategy{
		a52.NextCipherbitStragegy = &CustomInitialStatesNextCipherbitStrategy{&a52}
	}

	return &a52
}

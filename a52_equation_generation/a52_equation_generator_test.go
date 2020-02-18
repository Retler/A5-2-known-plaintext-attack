package a52_equation_generation

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"os"
	. "utils"
	. "a52"
)

func TestR1EqMajorityFunction(t *testing.T){
	register := make19IdentityMatrix()
	var values [19]int
	values[12] = 1
	r1 := R1_eg{register,values, make([]int, 0)}
	r1Input1 := r1.Register[15][:]
	r1Input2 := r1.Register[14][:]
	r1Input3 := r1.Register[12][:]
	r1Values := r1.Values[:]
	r1Len := int(19)
	equationIndices,result := eqMajority(r1Input1, r1Input2, r1Input3, r1Values, r1Len)
	//equationIndices,result := R1.equation_maj()
	indexMap, _ := MakeIndexMapAndLookup(19)
	expectedEquationsIndices := []int{12,15,indexMap[12][14],indexMap[12][15],indexMap[14][15]}

	// The bit at position 14 is flipped before majority function is invoked, which is why the result should be 1
	assert.Equal(t,int(1),result)
	assert.True(t, areSetsEqual(equationIndices, expectedEquationsIndices))

	r1.Values[12] = 0
	r1.Register[12][13] = 1
	equationIndices,result = eqMajority(r1Input1, r1Input2, r1Input3, r1Values, r1Len)
	expectedEquationsIndices = []int{12,13,15,indexMap[12][14],indexMap[12][15],indexMap[14][15],indexMap[13][14],indexMap[13][15]}

	assert.Equal(t,int(0),result)
	assert.True(t, areSetsEqual(equationIndices,expectedEquationsIndices))

	r1.Values[14] = 1
	r1.Values[15] = 1
	r1.Register[14][13] = 1

	equationIndices,result = eqMajority(r1Input1, r1Input2, r1Input3, r1Values, r1Len)
	expectedEquationsIndices = []int{12,15,indexMap[12][13],indexMap[12][14],indexMap[12][15],indexMap[14][15],indexMap[13][14]}

	assert.Equal(t,int(0),result)
	assert.True(t, areSetsEqual(equationIndices,expectedEquationsIndices))
}

func TestR2EqMajorityFunction(t *testing.T){
	register := make22IdentityMatrix()
	var values [22]int
	values[9] = 1
	r2 := R2_eg{register,values, make([]int, 0)}
	r2Input1 := r2.Register[9][:]
	r2Input2 := r2.Register[16][:]
	r2Input3 := r2.Register[13][:]
	r2Values := r2.Values[:]
	r2Len := int(22)
	equationIndices,result := eqMajority(r2Input1, r2Input2, r2Input3, r2Values, r2Len)
	//equationIndices,result := R1.equation_maj()
	indexMap, _ := MakeIndexMapAndLookup(22)
	expectedEquationsIndices := []int{9,13,indexMap[9][16],indexMap[13][16],indexMap[9][13]}
	// The bit at position 16 is flipped before majority function is invoked, which is why the result should be 1
	assert.Equal(t,int(1),result)
	//fmt.Println("expectedEquationsIndices", expectedEquationsIndices)
	//fmt.Println("equationIndices", equationIndices)
	assert.True(t, areSetsEqual(equationIndices, expectedEquationsIndices))

	r2.Values[9] = 0
	r2.Register[9][15] = 1
	equationIndices,result = eqMajority(r2Input1, r2Input2, r2Input3, r2Values, r2Len)
	expectedEquationsIndices = []int{9,13,15,indexMap[9][13],indexMap[13][15],indexMap[9][16],indexMap[15][16],indexMap[13][16]}

	assert.Equal(t,int(0),result)
	assert.True(t, areSetsEqual(equationIndices,expectedEquationsIndices))

	r2.Values[9] = 0
	r2.Values[13] = 1
	r2.Values[15] = 1
	r2.Values[17] = 1

	r2.Register[16][17] = 1
	equationIndices,result = eqMajority(r2Input1, r2Input2, r2Input3, r2Values, r2Len)
	expectedEquationsIndices = []int{9,13,15,indexMap[9][13],indexMap[13][15],indexMap[9][16],indexMap[9][17],indexMap[15][16],indexMap[15][17], indexMap[13][16], indexMap[13][17]}
	assert.Equal(t,int(1),result)
	assert.True(t, areSetsEqual(equationIndices,expectedEquationsIndices))
}

func TestR3EqMajorityFunction(t *testing.T){
	register := make23IdentityMatrix()
	var values [23]int
	values[16] = 1
	r3 := R3_eg{register,values, make([]int, 0)}
	r3Input1 := r3.Register[16][:]
	r3Input2 := r3.Register[13][:]
	r3Input3 := r3.Register[18][:]
	r3Values := r3.Values[:]
	r3Len := int(23)
	equationIndices,result := eqMajority(r3Input1, r3Input2, r3Input3, r3Values, r3Len)
	//equationIndices,result := R1.equation_maj()
	indexMap, _ := MakeIndexMapAndLookup(23)
	expectedEquationsIndices := []int{16,18,indexMap[13][16],indexMap[13][18],indexMap[16][18]}
	// The bit at position 16 is flipped before majority function is invoked, which is why the result should be 1
	assert.Equal(t,int(1),result)
	assert.True(t, areSetsEqual(equationIndices, expectedEquationsIndices))


	r3.Values[16] = 0
	r3.Register[16][15] = 1
	equationIndices,result = eqMajority(r3Input1, r3Input2, r3Input3, r3Values, r3Len)
	expectedEquationsIndices = []int{15,16,18,indexMap[13][16],indexMap[13][15],indexMap[13][18],indexMap[15][18],indexMap[16][18]}

	assert.Equal(t,int(0),result)
	assert.True(t, areSetsEqual(equationIndices,expectedEquationsIndices))

	r3.Values[18] = 1
	r3.Values[15] = 1
	r3.Values[17] = 1

	r3.Register[16][17] = 1
	equationIndices,result = eqMajority(r3Input1, r3Input2, r3Input3, r3Values, r3Len)
	expectedEquationsIndices = []int{15,16,17,18,indexMap[13][16],indexMap[13][15],indexMap[13][18],indexMap[13][17],indexMap[15][18],indexMap[16][18], indexMap[17][18]}
	assert.Equal(t,int(1),result)
	assert.True(t, areSetsEqual(equationIndices,expectedEquationsIndices))

	// Lets set x15 in all 3 inputs.
	r3.Register[13][15] = 1
	r3.Register[18][15] = 1

	equationIndices,result = eqMajority(r3Input1, r3Input2, r3Input3, r3Values, r3Len)
	expectedEquationsIndices = []int{15,16,17,18,indexMap[13][16],indexMap[13][18],indexMap[13][17],indexMap[16][18], indexMap[17][18]}
	assert.Equal(t,int(0),result)
	assert.True(t, areSetsEqual(equationIndices,expectedEquationsIndices))
}

func TestR1EquationClocking(t *testing.T){
	register := make19IdentityMatrix()
	var values [19]int
	r1 := R1_eg{register,values, make([]int, 0)}

	// Since we have the identity matrix at start we want the feedback to be an equation consisting of only the x's of the xor'ed registers
	var expectedFeedback1 [19]int
	expectedFeedback1[18] = 1
	expectedFeedback1[17] = 1
	expectedFeedback1[16] = 1
	expectedFeedback1[13] = 1

	var expectedFeedback2 [19]int
	expectedFeedback2[17] = 1
	expectedFeedback2[16] = 1
	expectedFeedback2[15] = 1
	expectedFeedback2[12] = 1

	var expectedLastRegister [19]int
	expectedLastRegister[16] = 1

	r1.clock()

	assert.Equal(t,expectedFeedback1, r1.Register[0])

	r1.clock()

	assert.Equal(t, expectedFeedback2, r1.Register[0])
	assert.Equal(t, expectedFeedback1, r1.Register[1])
	assert.Equal(t, expectedLastRegister, r1.Register[18])

	r1.Register[18][15] = 1
	r1.Register[17][14] = 1
	r1.Register[13][14] = 1

	// Change som Register Values, to emulate further clocking
	var expectedFeedback3 [19]int
	expectedFeedback3[16] = 1
	expectedFeedback3[14] = 1
	expectedFeedback3[11] = 1

	r1.clock()

	assert.Equal(t, r1.Register[0], expectedFeedback3)
}

func TestR2EquationClocking(t *testing.T){
	register := make22IdentityMatrix()
	var values [22]int
	r2 := R2_eg{register,values, make([]int, 0)}

	// Since we have the identity matrix at start we want the feedback to be an equation consisting of only the x's of the xor'ed registers
	var expectedFeedback1 [22]int
	expectedFeedback1[21] = 1
	expectedFeedback1[20] = 1

	var expectedFeedback2 [22]int
	expectedFeedback2[20] = 1
	expectedFeedback2[19] = 1

	var expectedFeedback3 [22]int
	expectedFeedback3[0] = 1
	expectedFeedback3[13] = 1
	expectedFeedback3[19] = 1

	var expectedLastRegister [22]int
	expectedLastRegister[18] = 1
	expectedLastRegister[0] = 1 // We set x_0 this to zero at Register 20 before the last clocking

	r2.clock()

	assert.Equal(t, expectedFeedback1, r2.Register[0])

	r2.clock()

	assert.Equal(t, expectedFeedback2, r2.Register[0])

	r2.Register[21][18] = 1
	r2.Register[21][13] = 1
	r2.Register[20][0] = 1

	r2.clock()

	assert.Equal(t, expectedFeedback3, r2.Register[0])
	assert.Equal(t, expectedFeedback2, r2.Register[1])
	assert.Equal(t, expectedFeedback1, r2.Register[2])
	assert.Equal(t, expectedLastRegister, r2.Register[21])
}

func TestR3EquationClocking(t *testing.T){
	register := make23IdentityMatrix()
	var values [23]int
	r3 := R3_eg{register,values, make([]int, 0)}

	var expectedFeedback1 [23]int
	expectedFeedback1[22] = 1
	expectedFeedback1[21] = 1
	expectedFeedback1[20] = 1
	expectedFeedback1[7] = 1

	r3.clock()

	assert.Equal(t, expectedFeedback1, r3.Register[0])

	var expectedFeedback2 [23]int

	// Registers 22,21,20,7 have Values x21,x20,x19,x6 after the last clocking
	// Now we add some extra x's
	r3.Register[22][5] = 1
	r3.Register[21][5] = 1
	r3.Register[20][5] = 1
	r3.Register[7][13] = 1
	r3.Register[22][3] = 1

	expectedFeedback2[21] = 1
	expectedFeedback2[20] = 1
	expectedFeedback2[19] = 1
	expectedFeedback2[13] = 1
	expectedFeedback2[6] = 1
	expectedFeedback2[5] = 1 // Since we XOR three x5's - there should be 1 left
	expectedFeedback2[3] = 1

	r3.clock()

	assert.Equal(t, expectedFeedback2, r3.Register[0])
	assert.Equal(t, expectedFeedback1, r3.Register[1])
}

func TestR4EquationClocking(t *testing.T){
	// The state of the a52 cipher after the initilization where Register x_i contains only x_i. In this case, we set the Values of all x's manually.
	a52 := A52_equation_generator{
		Session_key: [64]int{}, // 0 session key
		R1:          R1_eg{
			Register: make19IdentityMatrix(),
			Values:   [19]int{},
		},
		R2:          R2_eg{
			Register: make22IdentityMatrix(),
			Values:   [22]int{},
		},
		R3:          R3_eg{
			Register: make23IdentityMatrix(),
			Values:   [23]int{},
		},
		R4:          R4_eg{
			Register: [17]int{},
		},
	}

	/*
	R4[3],R4[7],R4[10] control the clocking of R1,R2,R3 in the following way:
	R1 is clocked if R4[10] == maj(R4[3],R4[7],R4[10])
	R2 is clocked if R4[3] == maj(R4[3],R4[7],R4[10])
	R3 is clocked if R4[7] == maj(R4[3],R4[7],R4[10])
	*/
	// Since all Values at all registers are 0, all registers should be clocked in this first clocking.
	//fmt.Println("a52 before clocking: ", a52)
	a52.R4.clock(&a52)
	//fmt.Println("a52 after clocking: ", a52)

	assert.True(t,a52.R1.Register[18][17]==1) // Since R1 should have been shifted R1[18] should now contain x_17
	assert.True(t,a52.R1.Register[18][18]==0) // Since R1 should have been shifted R1[18] should no longer contain x_18
	assert.True(t,a52.R1.Register[0][18]==1)  // Since R1 should have been shifted R1[0] should now contain x_18

	assert.True(t,a52.R2.Register[21][20]==1) // Since R2 should have been shifted R2[21] should now contain x_20
	assert.True(t,a52.R2.Register[21][21]==0) // Since R2 should have been shifted R2[21] should no longer contain x_21
	assert.True(t,a52.R2.Register[0][21]==1)  // Since R2 should have been shifted R2[0] should now contain x_21

	assert.True(t,a52.R3.Register[22][21]==1) // Since R3 should have been shifted R3[22] should now contain x_21
	assert.True(t,a52.R3.Register[22][22]==0) // Since R3 should have been shifted R3[22] should no longer contain x_22
	assert.True(t,a52.R3.Register[0][22]==1)  // Since R3 should have been shifted R3[0] should now contain x_22

	a52.R4.Register[10] = 1

	// Restart the registers
	a52.R1.Register = make19IdentityMatrix()
	a52.R2.Register = make22IdentityMatrix()
	a52.R3.Register = make23IdentityMatrix()

	a52.R4.clock(&a52)

	assert.True(t,a52.R1.Register[18][18]==1) // Since R1 should have been shifted R1[18] should now contain x_17
	assert.True(t,a52.R1.Register[18][17]==0) // Since R1 should have been shifted R1[18] should no longer contain x_18
	assert.True(t,a52.R1.Register[0][18]==0)  // Since R1 should have been shifted R1[0] should now contain x_18

	assert.True(t,a52.R2.Register[21][20]==1) // Since R2 should have been shifted R2[21] should now contain x_20
	assert.True(t,a52.R2.Register[21][21]==0) // Since R2 should have been shifted R2[21] should no longer contain x_21
	assert.True(t,a52.R2.Register[0][21]==1)  // Since R2 should have been shifted R2[0] should now contain x_21

	assert.True(t,a52.R3.Register[22][21]==1) // Since R3 should have been shifted R3[22] should now contain x_21
	assert.True(t,a52.R3.Register[22][22]==0) // Since R3 should have been shifted R3[22] should no longer contain x_22
	assert.True(t,a52.R3.Register[0][22]==1)  // Since R3 should have been shifted R3[0] should now contain x_22
}

func TestGetOutput(t *testing.T) {
	a52 := new(A52_equation_generator)
	var r1Values [19]int
	var r2Values [22]int
	var r3Values [23]int
	r1register := make19IdentityMatrix()
	r2register := make22IdentityMatrix()
	r3register := make23IdentityMatrix()

	a52.R1 = R1_eg{r1register,r1Values, make([]int, 0)}
	a52.R2 = R2_eg{r2register,r2Values, make([]int, 0)}
	a52.R3 = R3_eg{r3register,r3Values, make([]int, 0)}

	a52.R1.Values[12] = 1
	a52.R1.Values[18] = 1
	a52.R2.Values[21] = 1


	output := a52.getOutput()

	// 1 ^ 1 ^ 1 = 1
	assert.Equal(t, int(1), output.keystream)

	a52.R2.Values[13] = 1
	output = a52.getOutput()

	// 1 ^ 1 ^ 1 ^ 1 = 0
	assert.Equal(t, int(0), output.keystream)

	a52.R3.Values[16] = 1
	output = a52.getOutput()
	// 1 ^ 1 ^ 1 ^ 1 ^ 1 = 1
	assert.Equal(t, int(1), output.keystream)

	a52.R3.Values[22] = 1
	output = a52.getOutput()
	// 1 ^ 1 ^ 1 ^ 1 ^ 1 ^ 1 = 0
	assert.Equal(t, int(0), output.keystream)
}

func TestFinalEquationResult(t *testing.T){
	var key [64]int
	key[63] = 1
	key[0] = 1
	key[11] = 1
	var frameNumber [21]int
	//frameNumber[1] = 1

	r1, r2, r3, r4 := MakeRegisters1234()
	gsm := A52_equation_generator{key, r1, r2, r3, r4}
	keyStreamAndEquations := gsm.GenerateCipherForFrame(frameNumber, false)

	length := 114
	equations := make([][]int, length)
	for i := range equations {
		equations[i] = make([]int, length)
	}

	for i := 0; i < 114; i++ {
		equations[i] = makeEquation(keyStreamAndEquations[i])
	}

	solutions := make([]int, length)
	for i := 0; i < 114; i++ {
		solutions[i] = int(keyStreamAndEquations[i].keystream)
	}


	varValues := append(gsm.R1.Values[:], gsm.R1.linVarValues...)
	varValues = append(varValues, gsm.R2.Values[:]...)
	varValues = append(varValues, gsm.R2.linVarValues[:]...)
	varValues = append(varValues, gsm.R3.Values[:]...)
	varValues = append(varValues, gsm.R3.linVarValues[:]...)

	results := []Tuple{}

	for i := range(equations){
		result := 0
		for j := range(equations[i]){
			if(equations[i][j]==1){
				result ^= varValues[j]
			}
		}
		results = append(results, Tuple{result,solutions[i]})
	}

	for i := range(results){
		assert.True(t, results[i].One == results[i].Two)
	}

}

func TestXORIndices(t *testing.T)  {
	input1 := []int{0, 1,4,66,3,89, 97}
	input2 := []int{0,4,66,3,89, 90}

	expectedResult := []int{1, 90, 97}
	result := XorIndices(input1, input2)
	assert.True(t, areSetsEqual(result, expectedResult))
}

func TestTuplesToIndices(t *testing.T)  {
	tuples := []Tuple{Tuple{3,4},Tuple{4,5},Tuple{4,3}}
	indices := TupleToIndices(tuples,19)
	//fmt.Println("indices", indices)
	expectedIndices := []int{85}
	assert.True(t, TestEq(expectedIndices,indices))

}

func TestA52vsA52_arr(t *testing.T){
	var key [64]int
	key[63] = 1
	key[0] = 1
	key[11] = 1
	key[15] = 1
	key[7] = 1
	key[8] = 1
	key[19] = 1
	key[33] = 1
	key[34] = 1

	r1, r2, r3, r4 := MakeRegisters1234()
	gsm := A52_equation_generator{key, r1, r2, r3, r4}
	keyStreamAndEquations := gsm.GenerateCipherForFrame([21]int{}, false)

	gsm_old := A52{
		Session_key: key,
		R1:          R1{
			Register: [19]int{},
		},
		R2:          R2{
			Register: [22]int{},
		},
		R3:          R3{
			Register: [23]int{},
		},
		R4:          R4{
			Register: [17]int{},
		},
	}

	gsm_old_result := gsm_old.GenerateCipherForFrame([21]int{},false)

	for i := range(gsm_old_result){
		assert.True(t, gsm_old_result[i] == keyStreamAndEquations[i].keystream)
	}
}

func TestMakeLinVarValues(t *testing.T)  {
	values := []int{1,1,0,1}
	expectedValues := []int{1,0,1,0,1,0}
	smallLinVarValues := makeLinVarValues(values)
	assert.True(t, areSetsEqual(smallLinVarValues, expectedValues))
	//fmt.Println("smallLinVarValues: ", smallLinVarValues)
	//fmt.Println("expectedValues: ", expectedValues)

	values = []int{0,0,0,1}
	expectedValues = []int{0,0,0,0,0,0}
	smallLinVarValues = makeLinVarValues(values)
	assert.True(t, areSetsEqual(smallLinVarValues, expectedValues))

	//fmt.Println("smallLinVarValues: ", smallLinVarValues)
	//fmt.Println("expectedValues: ", expectedValues)

	values = []int{0,1,0,1}
	expectedValues = []int{0,0,0,0,1,0}
	smallLinVarValues = makeLinVarValues(values)
	assert.True(t, areSetsEqual(smallLinVarValues, expectedValues))
}


func TestXor(t *testing.T){
	r1BeforeXor :=[]int{0, 7, 13, 15, 1, 9, 10, 23, 24, 29, 36, 99, 116, 123, 100, 112, 127, 134, 106, 118, 163, 179, 108, 120, 165, 186, 111, 168, 38, 75, 76, 82, 40, 71, 102, 103, 109, 46, 77, 146, 154, 166, 51, 151, 159, 19, 22, 27, 28, 39, 70, 4,89, 90, 92, 94, 96, 72, 86, 117, 119, 121, 42, 87, 125, 126, 78, 12, 170, 172, 48, 93, 148, 156, 80, 14, 181,
	50, 95, 150, 158, 84, 98, 174, 183, 188}

	r1_18_indices := []int{1,2,3,7,9,10,17}

	result := XorIndices(r1_18_indices,r1BeforeXor)

	fmt.Println("Result test: ", result)

}


func TestGenerateCipherForFrameEquations(t *testing.T){
	a52 := A52_equation_generator{
		Session_key: [64]int{},
		R1:          R1_eg{
			Register:     [19][19]int{},
			Values:       [19]int{},
			linVarValues: nil,
		},
		R2:          R2_eg{
			Register:     [22][22]int{},
			Values:       [22]int{},
			linVarValues: nil,
		},
		R3:          R3_eg{
			Register:     [23][23]int{},
			Values:       [23]int{},
			linVarValues: nil,
		},
		R4:          R4_eg{
			Register: [17]int{},
		},
	}

	a52.Session_key[0] = 1
	a52.Session_key[0] = 1
	a52.Session_key[11] = 1
	frameNumber := [21]int{}
	output := a52.GenerateCipherForFrame(frameNumber, false)

	r1vectorVals := append(a52.R1.Values[:],a52.R1.linVarValues[:]...)
	r2vectorVals := append(a52.R2.Values[:],a52.R2.linVarValues[:]...)
	r3vectorVals := append(a52.R3.Values[:],a52.R3.linVarValues[:]...)
	allVals := append(r1vectorVals, r2vectorVals...)
	allVals = append(allVals, r3vectorVals...)

	for i := range(output){
		var r1vectorVars [190]int
		var r2vectorVars [253]int
		var r3vectorVars [276]int

		for _,val := range(output[i].r1Equations){
			r1vectorVars[val] = 1
		}

		for _,val := range(output[i].r2Equations){
			r2vectorVars[val] = 1
		}

		for _,val := range(output[i].r3Equations){
			r3vectorVars[val] = 1
		}

		allVars := append(r1vectorVars[:], r2vectorVars[:]...)
		allVars = append(allVars,r3vectorVars[:]...)
		equationResult := 0

		for j := range(allVars){
			if(allVars[j] == 1){
				equationResult ^= allVals[j]
			}
		}

		//fmt.Println("Equation result: ", equationResult, " Output result: ", output[i].keystream)
		//fmt.Println("R1 Values: ", a52.R1.Values)
		//fmt.Println("R2 Values: ", a52.R2.Values)
		//fmt.Println("R3 Values: ", a52.R3.Values)
		assert.Equal(t,output[i].keystream, equationResult)

	}
}

func TestFrame1And2Equations(t *testing.T){
	var key [64]int
	key[0] = 1
	r1, r2, r3, r4 := MakeRegisters1234()
	gsm := A52_equation_generator{key, r1, r2, r3, r4}


	frameNumber1 := [21]int{}
	frameNumber2 := [21]int{}
	frameNumber2[0] = 1
	DiffArrays12 := GetDiffOfFrames(frameNumber1,frameNumber2)

	keyStreamAndEquationsF1 := gsm.GenerateCipherForFrame(frameNumber1, false)

	varValues := append(gsm.R1.Values[:], gsm.R1.linVarValues...)
	varValues = append(varValues, gsm.R2.Values[:]...)
	varValues = append(varValues, gsm.R2.linVarValues[:]...)
	varValues = append(varValues, gsm.R3.Values[:]...)
	varValues = append(varValues, gsm.R3.linVarValues[:]...)

	keyStreamAndEquationsF2 := gsm.GenerateCipherForFrame(frameNumber2, false)

	length := 114
	equations1 := make([][]int, length)
	equations2 := make([][]int, length)

	for i := range equations1 {
		equations1[i] = make([]int, length)
		equations2[i] = make([]int, length)
	}

	for i := 0; i < 114; i++ {
		equations1[i] = makeEquation(keyStreamAndEquationsF1[i])
		equations2[i] = makeEquation(keyStreamAndEquationsF2[i])
	}

	solutions1 := make([]int, length)
	solutions2 := make([]int, length)

	for i := 0; i < 114; i++ {
		solutions1[i] = int(keyStreamAndEquationsF1[i].keystream)
		solutions2[i] = int(keyStreamAndEquationsF2[i].keystream)
	}

	TransformEquationsToMatchTheFirstFrame(DiffArrays12, equations2, solutions2)

	allEquations := append(equations1,equations2...)
	allSolutions := append(solutions1, solutions2...)
	//output := append(keyStreamAndEquationsF1[:114],keyStreamAndEquationsF2[:114]...)

	for i := range(allEquations){
		result := 0

		for j := range(allEquations[i]){
			if(allEquations[i][j] == 1){
				result ^= varValues[j]
			}
		}

		//fmt.Println("Equation result: ", equationResult, " Output result: ", output[i].keystream)
		//fmt.Println("R1 Values: ", a52.R1.Values)
		//fmt.Println("R2 Values: ", a52.R2.Values)
		//fmt.Println("R3 Values: ", a52.R3.Values)
		if(allSolutions[i] != result){
			fmt.Println("Equation ", i)

		}
		assert.Equal(t,allSolutions[i], result)

	}
}

func TestR4Difference(t *testing.T) {

	var key [64]int
	key[63] = 1
	key[60] = 1
	r1, r2, r3, r4 := MakeRegisters1234()
	r4Reg := [17]int{1,0,0,0,1,0,0,1,0,1,0,0,1,1,1,0,0}

	frameNumber1 := [21]int{}
	frameNumber2 := [21]int{}
	frameNumber2[0] = 1

	DiffArrays12 := GetDiffOfFrames(frameNumber1,frameNumber2)


	gsm := A52_equation_generator{key, r1, r2, r3, r4}
	//gsm.R4.Register = r4Reg

	gsm.initialize(frameNumber1)
	gsm.forceBits()

	fmt.Println("r4_gsm_frameNumbeR1: ", gsm.R4.Register)


	gsm = A52_equation_generator{key, r1, r2, r3, r4}
	//gsm.R4.Register = r4Reg


	gsm.initialize(frameNumber2)
	gsm.forceBits()

	fmt.Println("r4_gsm_frameNumbeR2: ", gsm.R4.Register)



	gsm = A52_equation_generator{key, r1, r2, r3, r4}
	gsm.R4.Register = r4Reg

	//gsm.GenerateCipherForFrame(frameNumber1, false)

	gsm.R4.Register = r4Reg


	//_ = gsm.GenerateCipherForFrame(frameNumber1, true)
	gsm.R1.trackValues()
	gsm.R2.trackValues()
	gsm.R3.trackValues()
	gsm.forceBits()

	fmt.Println("r4_gsm_frameNumbeR1: ", gsm.R4.Register)

	gsm.generateLinVarValues()
	gsm.run99()
	_ = gsm.makeKeystream()


	newR4 := [17]int{}
	for i := range(r4Reg){
		newR4[i] = r4Reg[i] ^ DiffArrays12.R4Diff[i]
	}
	gsm.R4.Register = newR4

	gsm.R1.trackValues()
	gsm.R2.trackValues()
	gsm.R3.trackValues()
	gsm.forceBits()

	fmt.Println("r4_gsm_frameNumbeR2: ", gsm.R4.Register)
}



func TestXorRegisters(t *testing.T){
	slice1 := []int{0,0,1,1,0}
	slice2 := []int{1,0,1,0,1}
	expectedResult := []int{1,0,0,1,1}

	result := XorRegisters(slice1,slice2)

	assert.True(t, TestEq(expectedResult,result))
}



func TestTransformFrames(t *testing.T){
	frame1 := Frame{0}.AsBinary()
	frame2 := Frame{1}.AsBinary()
	frame3 := Frame{2}.AsBinary()
	frame4 := Frame{3}.AsBinary()
	frame5 := Frame{4}.AsBinary()

	solutionsForFrame2 := make([]int,114)
	solutionsForFrame3 := make([]int,114)
	solutionsForFrame4 := make([]int,114)
	solutionsForFrame5 := make([]int,114)

	a52 := A52_equation_generator{
		Session_key: [64]int{0,1,0,0,1,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0},
	}

	outputFrame1 := a52.GenerateCipherForFrame(frame1, false)
	keystreamFrame1 := getKeystreamFromOutput(outputFrame1)
	equationsFrame1 := MakeEquationsFromOutput(outputFrame1)
	varValues := a52.getVarValues()

	fmt.Println("VarValues: ", varValues)

	for i, equation := range (equationsFrame1) {
		result := 0
		for j := range (equation) {
			if equation[j] == 1 {
				result ^= varValues[j]
			}
		}
		assert.Equal(t, keystreamFrame1[i], result)
	}

	outputFrame2 := a52.GenerateCipherForFrame(frame2, false)
	keystreamFrame2 := getKeystreamFromOutput(outputFrame2)
	equationsFrame2 := MakeEquationsFromOutput(outputFrame2)
	TransformEquationsToMatchTheFirstFrame(GetDiffOfFrames(frame1, frame2), equationsFrame2[:], solutionsForFrame2[:])

	fmt.Println("xor differences frame 2: ", solutionsForFrame2)

	solutionsForFrame2 = XorRegisters(solutionsForFrame2[:], keystreamFrame2[:])

	for i, equation := range (equationsFrame2) {
		result := 0
		for j := range (equation) {
			if equation[j] == 1 {
				result ^= varValues[j]
			}
		}
		assert.Equal(t, solutionsForFrame2[i], result)
	}

	outputFrame3 := a52.GenerateCipherForFrame(frame3, false)
	keystreamFrame3 := getKeystreamFromOutput(outputFrame3)
	equationsFrame3 := MakeEquationsFromOutput(outputFrame3)

	TransformEquationsToMatchTheFirstFrame(GetDiffOfFrames(frame1, frame3), equationsFrame3[:], solutionsForFrame3[:])

	solutionsForFrame3 = XorRegisters(solutionsForFrame3[:], keystreamFrame3[:])

	for i, equation := range (equationsFrame3) {
		result := 0
		for j := range (equation) {
			if equation[j] == 1 {
				result ^= varValues[j]
			}
		}
		assert.Equal(t, solutionsForFrame3[i], result)
	}

	outputFrame4 := a52.GenerateCipherForFrame(frame4, false)
	keystreamFrame4 := getKeystreamFromOutput(outputFrame4)
	equationsFrame4 := MakeEquationsFromOutput(outputFrame4)

	TransformEquationsToMatchTheFirstFrame(GetDiffOfFrames(frame1, frame4), equationsFrame4[:], solutionsForFrame4[:])

	solutionsForFrame4 = XorRegisters(solutionsForFrame4[:], keystreamFrame4[:])

	for i, equation := range (equationsFrame4) {
		result := 0
		for j := range (equation) {
			if equation[j] == 1 {
				result ^= varValues[j]
			}
		}
		assert.Equal(t, solutionsForFrame4[i], result)
	}

	outputFrame5 := a52.GenerateCipherForFrame(frame5, false)
	keystreamFrame5 := getKeystreamFromOutput(outputFrame5)
	equationsFrame5 := MakeEquationsFromOutput(outputFrame5)

	TransformEquationsToMatchTheFirstFrame(GetDiffOfFrames(frame1, frame5), equationsFrame5[:], solutionsForFrame5[:])

	solutionsForFrame5 = XorRegisters(solutionsForFrame5[:], keystreamFrame5[:])

	for i, equation := range (equationsFrame5) {
		result := 0
		for j := range (equation) {
			if equation[j] == 1 {
				result ^= varValues[j]
			}
		}
		assert.Equal(t, solutionsForFrame5[i], result)
	}

	keystream := append(keystreamFrame1[:], keystreamFrame2[:]...)
	keystream = append(keystream, keystreamFrame3[:]...)
	keystream = append(keystream, keystreamFrame4[:]...)
	keystream = append(keystream, keystreamFrame5[:]...)

	keystreamAsString := SliceElementsToString(keystream)

	writeStringToFile(keystreamAsString, "C:\\GolangKey\\keystream.txt")

	finalSolutions := append(keystreamFrame1[:],solutionsForFrame2...)
	finalSolutions = append(finalSolutions, solutionsForFrame3...)
	finalSolutions = append(finalSolutions, solutionsForFrame4...)
	finalSolutions = append(finalSolutions, solutionsForFrame5...)

	fmt.Println("final solutions: ", finalSolutions)


	writeStringToFile(keystreamAsString, "C:\\GolangKey\\golang_final_solutions.txt")
}

func TestFrame(t *testing.T){
	expected := [5][21]int{}

	expected[0] = [21]int{}
	expected[1] = [21]int{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1}
	expected[2] = [21]int{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0}
	expected[3] = [21]int{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,1}
	expected[4] = [21]int{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0}

	for i := 0; i < 5; i++{
		result := Frame{int64(i)}.AsBinary()

		for j := range(result){
			assert.Equal(t, expected[i][j],result[j])
		}
	}
}

func TestReverseLookUp(t *testing.T) {
	indexMap, reverseLookup := MakeIndexMapAndLookup(19)

	assert.True(t, len(reverseLookup) == 171)

	imRes := indexMap[4][7]
	rlRes := reverseLookup[imRes-19]
	assert.Equal(t,  4, rlRes.One)
	assert.Equal(t,  7, rlRes.Two)

}

func TestFrameBasedGetNextCipherbitStrategy(t *testing.T){
	a52_by_frame := A52{
		Session_key: [64]int{0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
		R1:          R1{},
		R2:          R2{},
		R3:          R3{},
		R4:          R4{},
		OutputCache:                   nil,
		CurrentFrameNumber:            Frame{0},
	}

	a52_by_bit := A52{
		Session_key: [64]int{0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
		R1:          R1{},
		R2:          R2{},
		R3:          R3{},
		R4:          R4{},
		OutputCache:                   nil,
		CurrentFrameNumber:            Frame{0},
	}

	a52_by_bit.NextCipherbitStragegy = &KeyBasedNextCipherbitStrategy{&a52_by_bit}

	output_by_frame1 := a52_by_frame.GenerateCipherForFrame(Frame{0}.AsBinary(),false)
	output_by_frame2 := a52_by_frame.GenerateCipherForFrame(Frame{1}.AsBinary(),false)
	output_by_frame3 := a52_by_frame.GenerateCipherForFrame(Frame{2}.AsBinary(),false)
	output_by_frame4 := a52_by_frame.GenerateCipherForFrame(Frame{3}.AsBinary(),false)
	output_by_frame5 := a52_by_frame.GenerateCipherForFrame(Frame{4}.AsBinary(),false)
	output_by_frame6 := a52_by_frame.GenerateCipherForFrame(Frame{5}.AsBinary(),false)


	keystream_by_frame := append(output_by_frame1[:114], output_by_frame2[:114]...)
	keystream_by_frame = append(keystream_by_frame, output_by_frame3[:114]...)
	keystream_by_frame = append(keystream_by_frame,output_by_frame4[:114]...)
	keystream_by_frame = append(keystream_by_frame, output_by_frame5[:114]...)
	keystream_by_frame = append(keystream_by_frame, output_by_frame6[:114]...)

	keystream_by_bits := a52_by_bit.GetNextCipherBits(113)
	keystream_by_bits = append(keystream_by_bits, a52_by_bit.GetNextCipherBits(115)...)
	keystream_by_bits = append(keystream_by_bits, a52_by_bit.GetNextCipherBits(112)...)
	keystream_by_bits = append(keystream_by_bits, a52_by_bit.GetNextCipherBits(116)...)
	keystream_by_bits = append(keystream_by_bits, a52_by_bit.GetNextCipherBits(8)...)
	keystream_by_bits = append(keystream_by_bits, a52_by_bit.GetNextCipherBits(220)...)

	assert.Equal(t, len(keystream_by_bits), len(keystream_by_frame))

	for i := range(keystream_by_bits){
		assert.Equal(t, keystream_by_bits[i], keystream_by_frame[i])
	}
}

func TestCustomInitialStatesNextCiphebitStrategy(t *testing.T){
	a52_by_frame := A52{
		Session_key: [64]int{0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
		R1:          R1{},
		R2:          R2{},
		R3:          R3{},
		R4:          R4{},
		OutputCache:                   nil,
		CurrentFrameNumber:            Frame{0},
	}

	// When generating the keystream based on initial states, we dont need the key
	a52_by_initial_states := A52{
		OutputCache:        nil,
		CurrentFrameNumber: Frame{0},
		InitialStates:      ReadInitialStates(),
	}

	a52_by_initial_states.NextCipherbitStragegy = &CustomInitialStatesNextCipherbitStrategy{&a52_by_initial_states}

	output_by_frame1 := a52_by_frame.GenerateCipherForFrame(Frame{0}.AsBinary(),false)
	output_by_frame2 := a52_by_frame.GenerateCipherForFrame(Frame{1}.AsBinary(),false)
	output_by_frame3 := a52_by_frame.GenerateCipherForFrame(Frame{2}.AsBinary(),false)
	output_by_frame4 := a52_by_frame.GenerateCipherForFrame(Frame{3}.AsBinary(),false)
	output_by_frame5 := a52_by_frame.GenerateCipherForFrame(Frame{4}.AsBinary(),false)
	output_by_frame6 := a52_by_frame.GenerateCipherForFrame(Frame{5}.AsBinary(),false)


	keystream_by_frame := append(output_by_frame1[:114], output_by_frame2[:114]...)
	keystream_by_frame = append(keystream_by_frame, output_by_frame3[:114]...)
	keystream_by_frame = append(keystream_by_frame,output_by_frame4[:114]...)
	keystream_by_frame = append(keystream_by_frame, output_by_frame5[:114]...)
	keystream_by_frame = append(keystream_by_frame, output_by_frame6[:114]...)

	keystream_by_initial_states := a52_by_initial_states.GetNextCipherBits(113)
	keystream_by_initial_states = append(keystream_by_initial_states, a52_by_initial_states.GetNextCipherBits(115)...)
	keystream_by_initial_states = append(keystream_by_initial_states, a52_by_initial_states.GetNextCipherBits(112)...)
	keystream_by_initial_states = append(keystream_by_initial_states, a52_by_initial_states.GetNextCipherBits(116)...)
	keystream_by_initial_states = append(keystream_by_initial_states, a52_by_initial_states.GetNextCipherBits(8)...)
	keystream_by_initial_states = append(keystream_by_initial_states, a52_by_initial_states.GetNextCipherBits(220)...)
	assert.Equal(t, len(keystream_by_initial_states), len(keystream_by_frame))

	for i := range(keystream_by_initial_states){
		assert.Equal(t, keystream_by_initial_states[i], keystream_by_frame[i])
	}
}

func TestAddLeadingZerosToBitString(t *testing.T){
	assert.Equal(t, "00000101", AddLeadingZerosToBitString("101"))
	assert.Equal(t, "00001010", AddLeadingZerosToBitString("1010"))
	assert.Equal(t, "10101001", AddLeadingZerosToBitString("10101001"))
	assert.Equal(t, "00000000", AddLeadingZerosToBitString("00000"))
}



/******************************************
			helper functions
 ******************************************/
func make19IdentityMatrix() ([19][19]int) {
	var result [19][19]int
	for i := 0; i < 19; i++{
		result[i][i] = 1
	}
	return result
}

func make22IdentityMatrix() ([22][22]int) {
	var result [22][22]int
	for i := 0; i < 22; i++{
		result[i][i] = 1
	}
	return result
}

func make23IdentityMatrix() ([23][23]int) {
	var result [23][23]int
	for i := 0; i < 23; i++{
		result[i][i] = 1
	}
	return result
}

// Checking equality of two sets
func areSetsEqual(slice1 []int, slice2 []int) bool{
	for _,val := range(slice1){
		if(!Contains(slice2,val)){
			return false
		}
	}

	if len(slice1) != len(slice2){
		return false
	}

	return true
}

func writeStringToFile(varValuesAsString string, filePath string) {
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	l, err := f.WriteString(varValuesAsString)
	if err != nil {
		fmt.Println(err)
		f.Close()
		os.Exit(1)
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


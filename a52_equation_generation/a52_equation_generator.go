package a52_equation_generation

import (
	"sort"
	. "utils"
	. "a52"
)

type A52_equation_generator struct {
	Session_key [64]int
	R1          R1_eg
	R2          R2_eg
	R3          R3_eg
	R4          R4_eg
}

type R1_eg struct {
	Register     [19][19]int
	Values       [19]int
	linVarValues []int
}

type R2_eg struct {
	Register     [22][22]int
	Values       [22]int
	linVarValues []int

}

type R3_eg struct {
	Register     [23][23]int
	Values       [23]int
	linVarValues []int

}

type R4_eg struct {
	Register [17]int
}

func (r1 *R1_eg) clock() {
	var feedbackresult [19]int

	for idx := range feedbackresult {
		feedbackresult[idx] = r1.Register[18][idx] ^ r1.Register[17][idx] ^ r1.Register[16][idx] ^ r1.Register[13][idx]
	}

	for i := 18; i > 0; i-- {
		r1.Register[i] = r1.Register[i-1]
	}
	r1.Register[0] = feedbackresult
}

func (r2 *R2_eg) clock() {
	var feedbackresult [22]int
	for idx := range feedbackresult {
		feedbackresult[idx] = r2.Register[21][idx] ^ r2.Register[20][idx]
	}

	for i := 21; i > 0; i-- {
		r2.Register[i] = r2.Register[i-1]
	}
	r2.Register[0] = feedbackresult
}

func (r3 *R3_eg) clock() {
	var feedbackresult [23]int

	for idx := range feedbackresult {
		feedbackresult[idx] = r3.Register[22][idx] ^ r3.Register[21][idx] ^ r3.Register[20][idx] ^ r3.Register[7][idx]
	}

	for i := 22; i > 0; i-- {
		r3.Register[i] = r3.Register[i-1]
	}
	r3.Register[0] = feedbackresult
}

func (r4 *R4_eg) clock(a52 *A52_equation_generator) {

	a := r4.Register[3]
	b := r4.Register[7]
	c := r4.Register[10]

	majresult := Maj(a, b, c)

	if a == majresult {
		a52.R2.clock()
	}
	if b == majresult {
		a52.R3.clock()
	}
	if c == majresult {
		a52.R1.clock()
	}

	var feedbackresult int

	feedbackresult = r4.Register[16] ^ r4.Register[11]

	for i := 16; i > 0; i-- {
		r4.Register[i] = r4.Register[i-1]
	}

	r4.Register[0] = feedbackresult
}

func (a52 *A52_equation_generator) GenerateCipherForFrame(frame_number [21]int, skipInitialise bool) [228]Output {
	if !skipInitialise {
		a52.initialize(frame_number)
	}

	a52.forceBits()

	a52.R1.trackValues()
	a52.R2.trackValues()
	a52.R3.trackValues()

	a52.generateLinVarValues()
	a52.run99()
	keyStreamAndEquations := a52.makeKeystream()
	return keyStreamAndEquations
}

// Initialise the internal state with K_c and frame number and save Values
func (a52 *A52_equation_generator) initialize(frame_number [21]int) {

	vanillaA52 := new(A52)
	vanillaA52.Session_key = a52.Session_key
	vanillaA52.Initialize(frame_number)

	a52.R1.Values = vanillaA52.R1.Register
	a52.R2.Values = vanillaA52.R2.Register
	a52.R3.Values = vanillaA52.R3.Register

	a52.R4.Register = vanillaA52.R4.Register
}

func (a52 *A52_equation_generator) generateLinVarValues() {
	a52.R1.linVarValues = makeLinVarValues(a52.R1.Values[:])
	a52.R2.linVarValues = makeLinVarValues(a52.R2.Values[:])
	a52.R3.linVarValues = makeLinVarValues(a52.R3.Values[:])
}

// make array of Values associated with the linearised Values
func makeLinVarValues(values []int) []int  {
	valuesLength := len(values)
	linVarValues := make([]int, valuesLength*(valuesLength-1)/2)
	idx := 0
	for i:=0; i<valuesLength; i++  {
		for j:= i + 1; j<valuesLength ;j++  {

			linVarValues[idx] = int(values[i])*int(values[j])
			idx++
		}
	}
	return linVarValues[:]
}

// Force bits R1[15], R2[16], R3[18] and R4[10] to 1 and start tracking Values
func (a52 *A52_equation_generator) forceBits() {

	a52.R1.setBit(15, 1)
	a52.R2.setBit(16, 1)
	a52.R3.setBit(18, 1)
	a52.R4.setBit(10, 1)
}

func (a52 *A52_equation_generator) run99() {
	for i := 0; i < 99; i++ {
		a52.R4.clock(a52)
	}
}

func (a52 *A52_equation_generator) makeKeystream() [228]Output {
	var result [228]Output
	for i := 0; i < 228; i++ {
		a52.R4.clock(a52)
		result[i] = a52.getOutput()
	}
	return result
}

// Wrapper for keystream output
type Output struct {
	keystream   int
	r1Equations []int
	r2Equations []int
	r3Equations []int
}

func (a52 *A52_equation_generator) getOutput() Output {
	//r1Equations,keyStream := a52.R1.equation_maj()
	r1Equations,r1result := a52.R1.getMajorityOutputEquationAndBit()
	r2Equations,r2result := a52.R2.getMajorityOutputEquationAndBit()
	r3Equations,r3result := a52.R3.getMajorityOutputEquationAndBit()

	r1_18_indices, r1_18 := getVarIndicesAndVal(a52.R1.Register[18][:], a52.R1.Values[:])
	r2_21_indices, r2_21 := getVarIndicesAndVal(a52.R2.Register[21][:], a52.R2.Values[:])
	r3_22_indices, r3_22 := getVarIndicesAndVal(a52.R3.Register[22][:], a52.R3.Values[:])
	keyStream := r1result ^ r2result ^ r3result ^ r1_18 ^ r2_21 ^ r3_22

	sort.Ints(r1Equations)
	sort.Ints(r1_18_indices)
	r1Equations = XorIndices(r1_18_indices,r1Equations)
	sort.Ints(r1Equations)

	r2Equations = XorIndices(r2_21_indices, r2Equations)
	r3Equations = XorIndices(r3_22_indices, r3Equations)

	return Output{keyStream, r1Equations, r2Equations, r3Equations}
}

func getVarIndicesAndVal(Register []int, RegisterValues []int) ([]int, int) {
	var Register_var_indexes []int
	var value int = 0
	for idx := range Register {
		if Register[idx] == 1 {
			value ^= RegisterValues[idx]
			Register_var_indexes = append(Register_var_indexes, idx)
		}
	}
	return Register_var_indexes, value
}

func(r3 R3_eg) getMajorityOutputEquationAndBit() ([]int, int) {
	// inputs from R3 Register
	r3Input1 := r3.Register[16][:]
	r3Input2 := r3.Register[13][:]
	r3Input3 := r3.Register[18][:]
	r3Values := r3.Values[:]
	r3Len := int(23)
	r3equations, r3majResult := eqMajority(r3Input1, r3Input2, r3Input3, r3Values, r3Len)
	return r3equations, r3majResult
}

func(r2 R2_eg) getMajorityOutputEquationAndBit() ([]int, int) {
	// inputs from R2 Register
	r2Input1 := r2.Register[9][:]
	r2Input2 := r2.Register[16][:]
	r2Input3 := r2.Register[13][:]
	r2Values := r2.Values[:]
	r2Len := int(22)
	// R2 equation and result of majority function
	r2equations, r2majResult := eqMajority(r2Input1, r2Input2, r2Input3, r2Values, r2Len)
	return r2equations, r2majResult
}

func(r1 R1_eg) getMajorityOutputEquationAndBit() ([]int, int) {
	// inputs from R1 Register
	r1Input1 := r1.Register[15][:]
	r1Input2 := r1.Register[14][:]
	r1Input3 := r1.Register[12][:]

	r1Values := r1.Values[:]
	r1Len := int(19)
	// R1 equation and result of majority function
	r1equations, r1majResult := eqMajority(r1Input1, r1Input2, r1Input3, r1Values, r1Len)
	return r1equations, r1majResult
}

// Majority function
// input2 is the one XOR'ed with 1
func eqMajority(input1 []int, input2 []int, input3 []int, values []int, len int) ([]int, int) {

	// extract varibles and value from input1, input2 and input3
	input1_var_indexes, input1_val := getVarIndicesAndVal(input1, values)
	input2_var_indexes, input2_val := getVarIndicesAndVal(input2, values)
	input3_var_indexes, input3_val := getVarIndicesAndVal(input3, values)

	// Results of majority function
	result := input1_val*(input2_val^1) ^ (input2_val^1)*input3_val ^ input1_val*input3_val

	// multiply elements in Registers
	tuples12 := GetTuples(input1_var_indexes, input2_var_indexes)
	tuples23 := GetTuples(input2_var_indexes, input3_var_indexes)
	tuples13 := GetTuples(input1_var_indexes, input3_var_indexes)

	// convert tuples into indices in equation array
	indices1 := TupleToIndices(tuples12, len)

	indices2 := TupleToIndices(tuples23, len)

	indices3 := TupleToIndices(tuples13, len)

	// XOR input 1 and 3 to get the indices of all original non-linearized variables
	indices4 := XorIndices(input1_var_indexes, input3_var_indexes)

	// XOR input1,2,input2,3,input1,3
	indices5 := XorIndices(indices1, indices2)
	indices6 := XorIndices(indices5, indices3)

	// XOR x_0...x_18 and x_19x_20...(x_n-1)(x_n)
	equations := XorIndices(indices4, indices6)

	return equations, result
}

// Set bit in array at idx in Register at idx1 to value
func (r *R1_eg) setBit(idx1 int, value int) {
	r.Values[idx1] = value
}

func (r *R2_eg) setBit(idx1 int, value int) {
	r.Values[idx1] = value
}
func (r *R3_eg) setBit(idx1 int, value int) {
	r.Values[idx1] = value
}

func (r *R4_eg) setBit(idx1 int, value int) {
	r.Register[idx1] = value
}

// Set the equation matrix equal to the identity matrix
func (r *R1_eg) trackValues() {
	r.setZeros()
	for idx := range r.Register {
		r.Register[idx][idx] = 1
	}
}

func (r *R2_eg) trackValues() {
	r.setZeros()
	for idx := range r.Register {
		r.Register[idx][idx] = 1
	}
}

func (r *R3_eg) trackValues() {
	r.setZeros()
	for idx := range r.Register {
		r.Register[idx][idx] = 1
	}
}

// Set all elements to zero - ikke ødvendigt
func (r *R1_eg) setZeros() {
	for idx1 := range r.Register {
		for idx2 := range r.Register {
			r.Register[idx1][idx2] = 0
		}
	}

}

func (r *R2_eg) setZeros() {
	for idx1 := range r.Register {
		for idx2 := range r.Register {
			r.Register[idx1][idx2] = 0
		}
	}

}

func (r *R3_eg) setZeros() {
	for idx1 := range r.Register {
		for idx2 := range r.Register {
			r.Register[idx1][idx2] = 0
		}
	}

}

func (r *R4_eg) setZeros() {
	for idx := range r.Register {
		r.Register[idx] = 0
	}
}

func makeEquation(output Output) []int {
	r1Indices := output.r1Equations
	r2Indices := output.r2Equations
	r3Indices := output.r3Equations

	r1Equation := [190]int{}
	r2Equation := [253]int{}
	r3Equation := [276]int{}

	for _, val := range r1Indices {
		r1Equation[val] = 1
	}

	for _, val := range r2Indices {
		r2Equation[val] = 1
	}

	for _, val := range r3Indices {
		r3Equation[val] = 1
	}

	equation := append(r1Equation[:], r2Equation[:]...)
	equation = append(equation, r3Equation[:]...)
	return equation

}

func getKeystreamFromOutput(output [228]Output) [114]int{
	result := [114]int{}

	for i := 0; i < 114; i++{
		result[i] = output[i].keystream
	}
	return result
}

func MakeEquationsFromOutput(output [228]Output) [114][]int{
	result := [114][]int{}

	for i := 0; i < 114; i++{
		result[i] = makeEquation(output[i])
	}

	return result
}

func(a52 *A52_equation_generator) getVarValues() []int{
	result := []int{}
	result = a52.R1.Values[:]
	result = append(result,a52.R1.linVarValues...)
	result = append(result, a52.R2.Values[:]...)
	result = append(result, a52.R2.linVarValues...)
	result = append(result, a52.R3.Values[:]...)
	result = append(result, a52.R3.linVarValues[:]...)

	return result
}

func TransformEquationsToMatchTheFirstFrame(DiffArrays DiffArrays, equations [][]int, solutions []int) {
	_, lu1 := MakeIndexMapAndLookup(19)
	_, lu2 := MakeIndexMapAndLookup(22)
	_, lu3 := MakeIndexMapAndLookup(23)

	for i := range equations {

		for j:= 0; j<=18; j++ {
			if equations[i][j] == 1 {
				solutions[i] ^= DiffArrays.R1Diff[j]
			}
		}

		for j:= 190; j<=211; j++ {
			if equations[i][j] == 1 {
				solutions[i] ^= DiffArrays.R2Diff[j-190]
			}
		}

		for j:= 443; j<=465; j++ {
			if equations[i][j] == 1 {
				solutions[i] ^= DiffArrays.R3Diff[j-443]
			}
		}

		for j:= 19; j<190; j++ {
			// a_f+1 * b_f+1
			if equations[i][j] == 1 {
				tuple := lu1[j-19]
				// a_f * d_a
				if(DiffArrays.R1Diff[tuple.One] == 1){
					if(equations[i][tuple.Two] == 1){
						equations[i][tuple.Two] = 0
					}else{
						equations[i][tuple.Two] = 1
					}
				}

				if(DiffArrays.R1Diff[tuple.Two] == 1){
					if(equations[i][tuple.One] == 1){
						equations[i][tuple.One] = 0
					}else{
						equations[i][tuple.One] = 1
					}
				}

				// result = result ^ d_a * d_b
				solutions[i] ^= DiffArrays.R1Diff[tuple.One] * DiffArrays.R1Diff[tuple.Two]
			}
		}

		for j:= 212; j<443; j++ {
			// a_f+1 * b_f+1
			if equations[i][j] == 1 {
				tuple := lu2[j-212]

				if(DiffArrays.R2Diff[tuple.One] == 1){
					if(equations[i][tuple.Two+190] == 1){
						equations[i][tuple.Two+190] = 0
					}else{
						equations[i][tuple.Two+190] = 1
					}
				}

				if(DiffArrays.R2Diff[tuple.Two] == 1){
					if(equations[i][tuple.One+190] == 1){
						equations[i][tuple.One+190] = 0
					}else{
						equations[i][tuple.One+190] = 1
					}
				}

				// res = res ^ d_a * d_b
				solutions[i] ^= DiffArrays.R2Diff[tuple.One] * DiffArrays.R2Diff[tuple.Two]
			}
		}

		for j:= 466; j<719; j++ {
			// a_f+1 * b_f+1
			if equations[i][j] == 1 {
				tuple := lu3[j-466]

				if(DiffArrays.R3Diff[tuple.One] == 1){
					if(equations[i][tuple.Two+443] == 1){
						equations[i][tuple.Two+443] = 0
					}else{
						equations[i][tuple.Two+443] = 1
					}
				}

				if(DiffArrays.R3Diff[tuple.Two] == 1){
					if(equations[i][tuple.One+443] == 1){
						equations[i][tuple.One+443] = 0
					}else{
						equations[i][tuple.One+443] = 1
					}
				}// res = res ^ d_a * d_b

				solutions[i] ^= DiffArrays.R3Diff[tuple.One] * DiffArrays.R3Diff[tuple.Two]
			}
		}
	}
}

// make empty Registers
func MakeRegisters1234() (R1_eg, R2_eg, R3_eg, R4_eg) {
	r1 := R1_eg{[19][19]int{}, [19]int{}, make([]int, 190-19)}
	r2 := R2_eg{[22][22]int{}, [22]int{}, make([]int, 253-22)}
	r3 := R3_eg{[23][23]int{}, [23]int{}, make([]int, 276-23)}
	r4 := R4_eg{[17]int{}}
	return r1, r2, r3, r4
}

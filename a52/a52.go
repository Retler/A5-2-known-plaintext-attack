package a52

import (
	"strconv"
	. "utils"
)

type A52 struct {
	Session_key                   [64]int
	R1                            R1
	R2                            R2
	R3                            R3
	R4                            R4
	UsedBitsOfCurrentCipherOutput int
	OutputCache                   []int
	CurrentFrameNumber            Frame
	NextCipherbitStragegy         NextCipherbitStrategy
	InitialStates                 InitialStates
}

type NextCipherbitStrategy interface {
	GetNextCipherBits(ammount int) []int
}

type KeyBasedNextCipherbitStrategy struct {
	A52 *A52
}

type CustomInitialStatesNextCipherbitStrategy struct{
	A52 *A52
}

func(f *KeyBasedNextCipherbitStrategy) GetNextCipherBits(ammount int)[]int{
	skipInitialize := false
	if(f.A52.OutputCache == nil){ // We have not yet generated any output
		f.A52.GenerateCipherForFrame(f.A52.CurrentFrameNumber.AsBinary(), skipInitialize)
	}

	if(f.A52.UsedBitsOfCurrentCipherOutput == 114){ // We have used all bits, lets generate new cipher
		f.A52.GenerateCipherForFrame(f.A52.CurrentFrameNumber.AsBinary(), skipInitialize)
	}

	result := make([]int, 0)
	if ammount > 114-f.A52.UsedBitsOfCurrentCipherOutput { // The number of cipherbits we want is bigger than those left in the current output, so we generate a new output
		restOfBits := f.A52.OutputCache[f.A52.UsedBitsOfCurrentCipherOutput:]
		f.A52.UsedBitsOfCurrentCipherOutput = 114

		result = append(result, restOfBits...)

		missingBits := ammount - len(restOfBits)

		moreBits := f.GetNextCipherBits(missingBits)

		result = append(result, moreBits...)

		return result
	}

	if ammount <= 114-f.A52.UsedBitsOfCurrentCipherOutput {
		result = append(result, f.A52.OutputCache[f.A52.UsedBitsOfCurrentCipherOutput:f.A52.UsedBitsOfCurrentCipherOutput+ammount]...)
		f.A52.UsedBitsOfCurrentCipherOutput = f.A52.UsedBitsOfCurrentCipherOutput + ammount
		return result
	}

	return result
}

func GetDiffOfFrames(frameNumber [21]int, frameNumber2 [21]int) (DiffArrays){
	a52Diff := A52{
		Session_key: [64]int{},
		R1:          R1{},
		R2:          R2{},
		R3:          R3{},
		R4:          R4{},
	}
	for idx, _ := range frameNumber {
		a52Diff.ClockAll()
		a52Diff.R1.Register[0] = a52Diff.R1.Register[0] ^ frameNumber[idx]
		a52Diff.R2.Register[0] = a52Diff.R2.Register[0] ^ frameNumber[idx]
		a52Diff.R3.Register[0] = a52Diff.R3.Register[0] ^ frameNumber[idx]
		a52Diff.R4.Register[0] = a52Diff.R4.Register[0] ^ frameNumber[idx]
	}

	r1_1 := [19]int{}
	r2_1 := [22]int{}
	r3_1 := [23]int{}
	r4_1 := [17]int{}

	copy(r1_1[:],a52Diff.R1.Register[:])
	copy(r2_1[:],a52Diff.R2.Register[:])
	copy(r3_1[:],a52Diff.R3.Register[:])
	copy(r4_1[:],a52Diff.R4.Register[:])

	a52Diff.R1.Register = [19]int{}
	a52Diff.R2.Register = [22]int{}
	a52Diff.R3.Register = [23]int{}
	a52Diff.R4.Register = [17]int{}

	for idx, _ := range frameNumber2 {
		a52Diff.ClockAll()
		a52Diff.R1.Register[0] = a52Diff.R1.Register[0] ^ frameNumber2[idx]
		a52Diff.R2.Register[0] = a52Diff.R2.Register[0] ^ frameNumber2[idx]
		a52Diff.R3.Register[0] = a52Diff.R3.Register[0] ^ frameNumber2[idx]
		a52Diff.R4.Register[0] = a52Diff.R4.Register[0] ^ frameNumber2[idx]
	}

	r1Diff := XorRegisters(r1_1[:],a52Diff.R1.Register[:])
	r2Diff := XorRegisters(r2_1[:],a52Diff.R2.Register[:])
	r3Diff := XorRegisters(r3_1[:],a52Diff.R3.Register[:])
	r4Diff := XorRegisters(r4_1[:],a52Diff.R4.Register[:])

	return  DiffArrays{r1Diff, r2Diff, r3Diff,r4Diff}
}

type DiffArrays struct {
	R1Diff []int
	R2Diff []int
	R3Diff []int
	R4Diff []int
}

func(c *CustomInitialStatesNextCipherbitStrategy) GetNextCipherBits(ammount int)[]int{
	skipInitialize := true
	result := make([]int, 0)

	// We have not yet generated any output - so lets do that
	if(c.A52.OutputCache == nil){
		c.A52.R1.Register = c.A52.InitialStates.R1
		c.A52.R2.Register = c.A52.InitialStates.R2
		c.A52.R3.Register = c.A52.InitialStates.R3
		c.A52.R4.Register = c.A52.InitialStates.R4
		c.A52.GenerateCipherForFrame(c.A52.CurrentFrameNumber.AsBinary(), skipInitialize)
	}

	// We have used all bits of the output cache - lets generate 114 new bits
	if(c.A52.UsedBitsOfCurrentCipherOutput == 114){ // We have used all bits, lets generate new cipher

		// First, we need to calculate new initial states of current frame by xoring the initial states by the differences between the first frame and current frame
		diffs := GetDiffOfFrames(Frame{0}.AsBinary(),c.A52.CurrentFrameNumber.AsBinary())
		copy(c.A52.R1.Register[:],XorRegisters(c.A52.InitialStates.R1[:],diffs.R1Diff))
		copy(c.A52.R2.Register[:],XorRegisters(c.A52.InitialStates.R2[:],diffs.R2Diff))
		copy(c.A52.R3.Register[:],XorRegisters(c.A52.InitialStates.R3[:],diffs.R3Diff))
		copy(c.A52.R4.Register[:],XorRegisters(c.A52.InitialStates.R4[:],diffs.R4Diff))

		c.A52.GenerateCipherForFrame(c.A52.CurrentFrameNumber.AsBinary(), skipInitialize)
	}

	// The number of cipherbits we want is bigger than those left in the current output, so we generate a new output
	if ammount > 114-c.A52.UsedBitsOfCurrentCipherOutput {
		restOfBits := c.A52.OutputCache[c.A52.UsedBitsOfCurrentCipherOutput:]
		c.A52.UsedBitsOfCurrentCipherOutput = 114

		result = append(result, restOfBits...)

		missingBits := ammount - len(restOfBits)

		moreBits := c.GetNextCipherBits(missingBits)

		result = append(result, moreBits...)

		return result
	}

	// The number of bits we want is smaller than what is left in the cache - so we get the ammount that we need from output cache
	if ammount <= 114-c.A52.UsedBitsOfCurrentCipherOutput {
		result = append(result, c.A52.OutputCache[c.A52.UsedBitsOfCurrentCipherOutput:c.A52.UsedBitsOfCurrentCipherOutput+ammount]...)
		c.A52.UsedBitsOfCurrentCipherOutput = c.A52.UsedBitsOfCurrentCipherOutput + ammount
		return result
	}

	return result
}

type R1 struct {
	Register [19]int
}

type R2 struct {
	Register [22]int
}

type R3 struct {
	Register [23]int
}

type R4 struct {
	Register [17]int
}

func(a52 *A52) GetNextCipherBits(ammount int)[]int{
	return a52.NextCipherbitStragegy.GetNextCipherBits(ammount)
}

func (r1 *R1) clock() {
	feedbackresult := r1.Register[18] ^ r1.Register[17] ^ r1.Register[16] ^ r1.Register[13]
	for i := 18; i > 0; i-- {
		r1.Register[i] = r1.Register[i-1]
	}
	r1.Register[0] = feedbackresult
}

func (r2 *R2) clock() {
	feedbackresult := r2.Register[21] ^ r2.Register[20]
	for i := 21; i > 0; i-- {
		r2.Register[i] = r2.Register[i-1]
	}
	r2.Register[0] = feedbackresult
}

func (r3 *R3) clock() {
	feedbackresult := r3.Register[22] ^ r3.Register[21] ^ r3.Register[20] ^ r3.Register[7]
	for i := 22; i > 0; i-- {
		r3.Register[i] = r3.Register[i-1]
	}
	r3.Register[0] = feedbackresult
}

func (r4 *R4) clock(a52 *A52, clockR4Only bool) {

	// In initialization phase, we dont want R4 to clock the other registers.
	if(!clockR4Only){
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
	}

	feedbackresult := r4.Register[16] ^ r4.Register[11]
	for i := 16; i > 0; i-- {
		r4.Register[i] = r4.Register[i-1]
	}
	r4.Register[0] = feedbackresult

}

func (a52 *A52) ClockAll() {
	a52.R1.clock()
	a52.R2.clock()
	a52.R3.clock()

	clockOnlyR4 := true

	a52.R4.clock(a52,clockOnlyR4)

}

func (a52 *A52) GenerateCipherForFrame(frame_number [21]int, skipInitialize bool) [228]int{
	if !skipInitialize{
		a52.Initialize(frame_number)
	}
	a52.forceBits()
	a52.run99()
	result := a52.makeCipher()
	a52.OutputCache = result[:114]
	a52.UsedBitsOfCurrentCipherOutput = 0
	a52.CurrentFrameNumber.Value = a52.CurrentFrameNumber.Value + 1
	return result
}

func (a52 *A52) Initialize(frame_number [21]int) {
	a52.SetAllLsfrsToZero()

	for idx, _ := range a52.Session_key {
		a52.ClockAll()
		a52.R1.Register[0] = a52.R1.Register[0] ^ a52.Session_key[idx]
		a52.R2.Register[0] = a52.R2.Register[0] ^ a52.Session_key[idx]
		a52.R3.Register[0] = a52.R3.Register[0] ^ a52.Session_key[idx]
		a52.R4.Register[0] = a52.R4.Register[0] ^ a52.Session_key[idx]
	}

	for idx, _ := range frame_number {
		a52.ClockAll()
		a52.R1.Register[0] = a52.R1.Register[0] ^ frame_number[idx]
		a52.R2.Register[0] = a52.R2.Register[0] ^ frame_number[idx]
		a52.R3.Register[0] = a52.R3.Register[0] ^ frame_number[idx]
		a52.R4.Register[0] = a52.R4.Register[0] ^ frame_number[idx]
	}
}

func (a52 *A52) SetAllLsfrsToZero() {
	for idx, _ := range a52.R1.Register {
		a52.R1.Register[idx] = 0
	}
	for idx, _ := range a52.R2.Register {
		a52.R2.Register[idx] = 0
	}
	for idx, _ := range a52.R3.Register {
		a52.R3.Register[idx] = 0
	}
	for idx, _ := range a52.R4.Register {
		a52.R4.Register[idx] = 0
	}
}

func (a52 *A52) forceBits() {
	a52.R1.Register[15] = 1
	a52.R2.Register[16] = 1
	a52.R3.Register[18] = 1
	a52.R4.Register[10] = 1
}

func (a52 *A52) run99() {
	clockR4Only := false
	for i := 0; i < 99; i++ {
		a52.R4.clock(a52,clockR4Only)
	}
}

func (a52 *A52) makeCipher() [228]int {
	var result [228]int
	clockR4Only := false

	for i := 0; i < 228; i++ {
		a52.R4.clock(a52,clockR4Only)
		result[i] = a52.getOutput()
	}
	return result
}

func (a52 *A52) getOutput() int {
	bit1 := Maj(a52.R1.Register[15], a52.R1.Register[14]^1, a52.R1.Register[12])
	bit2 := Maj(a52.R2.Register[16]^1, a52.R2.Register[13], a52.R2.Register[9])
	bit3 := Maj(a52.R3.Register[18], a52.R3.Register[16], a52.R3.Register[13]^1)

	return bit1 ^ a52.R1.Register[18] ^ bit2 ^ a52.R2.Register[21] ^ bit3 ^ a52.R3.Register[22]
}

// Given bytes of a plaintext/ciphertext, outputs bytes of the ciphertext/plaintext
func (a52 *A52) EncryptDecryptFile(content []byte) []byte{
	result := make([]byte,len(content))
	result_byte := [8]int{}

	for j,b := range(content){
		current_byte_keystream := a52.GetNextCipherBits(8)

		keystream_as_string := SliceElementsToString(current_byte_keystream)
		byte_as_string_bits := strconv.FormatInt(int64(b), 2)
		byte_as_string_bits = AddLeadingZerosToBitString(byte_as_string_bits)

		for i := range(byte_as_string_bits){
			result_byte[i] = int(byte_as_string_bits[i] ^ keystream_as_string[i]) // Do the bitwise XOR of the current byte ciphertext and the keystream
		}

		plaintext_byte_string := SliceElementsToString(result_byte[:])
		plaintext_byte_number, _ := strconv.ParseUint(plaintext_byte_string, 2, 8) // Convert ciphertext to a byte
		result[j] = byte(plaintext_byte_number)
	}

	return result
}

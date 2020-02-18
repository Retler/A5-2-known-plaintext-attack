package utils

import (
	"strconv"
	"strings"
	"bytes"
	"os"
	"bufio"
	"log"
	"math/big"
	"crypto/rand"
	"fmt"
)

type Frame struct{
	Value int64
}

type InitialStates struct {
	R1 [19]int
	R2 [22]int
	R3 [23]int
	R4 [17]int
}

func(f Frame) AsBinary()[21]int{
	bitsAsString := strconv.FormatInt(f.Value, 2)
	bitsAsStringArray := strings.Split(bitsAsString, "")
	result := [21]int{}
	idx1 := len(bitsAsStringArray)-1
	idx2 := len(result)-1

	for idx1 >= 0{
		bit,_ := strconv.Atoi(bitsAsStringArray[idx1])
		result[idx2] = bit
		idx1--
		idx2--
	}

	return result
}

// maps linearized values to a map
func MakeIndexMapAndLookup(len int) ([][]int, []Tuple) {
	index := int(len)
	reverseLookUp := make([]Tuple, (len*(len-1)/2))
	indexLookup := 0

	// Golang demands we make the slices in this way
	indexMap := make([][]int, len)
	for i := range indexMap {
		indexMap[i] = make([]int, len)
	}

	for i := 0; i < int(len); i++ {
		for j := i + 1; j < int(len); j++ {
			indexMap[i][j] = index
			reverseLookUp[indexLookup] = Tuple{i,j}
			index++
			indexLookup++
		}
	}
	return indexMap, reverseLookUp
}

// XOR two arrays of indices
func XorIndices(indices1 []int, indices2 []int) []int {
	var result []int

	for _, val := range indices1 {
		if Contains(indices2, val) { // duplicate values will not be added to the result
			continue
		}
		result = append(result, val)
	}

	for _, val := range indices2 {
		if Contains(indices1, val) { // duplicate values will not be added to the result
			continue
		}
		result = append(result, val)
	}
	return result
}


// maps intermediary tuples to indices in final equation array
func TupleToIndices(tuples []Tuple, indexMapLen int) []int {
	indexMap, _ := MakeIndexMapAndLookup(indexMapLen)
	var indices []int

	// Sorts tuples, x_14x_15 = x_15x_14
	for idx := range tuples {
		tuple := tuples[idx]
		if tuple.One > tuple.Two {
			a := tuple.One
			b := tuple.Two
			tuples[idx] = Tuple{b, a}
		}
	}

	//fmt.Println("tuples after sort: ", tuples)

	for idx := range tuples {
		tuple := tuples[idx]
		if tuple.One == tuple.Two { // Is not a linearized variable since x_1*x_1 = x_1 when all x's are binary
			indices = append(indices, tuple.One) // Then the indexOfTap is just one of the tuple values
		} else {
			indices = append(indices, indexMap[tuple.One][tuple.Two])
		}
		//fmt.Println("indexMap[",tuple.One,tuple.Two,"] is ", indexMap[tuple.One][tuple.Two])

	}

	return Unique(indices)
}

// multiplication of two sets
// e.g [1,2],[3] => (x1 + x2) * x3 = x1x3 + x2x3 => [(1,3), (2,3)]
func GetTuples(indices1 []int, indices2 []int) []Tuple {
	var tuples []Tuple
	for idx1 := range indices1 {
		for idx2 := range indices2 {
			tuples = append(tuples, Tuple{indices1[idx1], indices2[idx2]})
		}
	}
	return tuples
}

type Tuple struct {
	One int
	Two int
}

// returns true if int array contains element e
func Contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Removes duplicates from slice
func Unique(indices []int) []int {
	var indices_to_remove []int
	var indices_result []int

	for i := range(indices){
		for j := i+1; j < len(indices); j++{
			if(indices[i] == indices[j]){
				indices_to_remove = append(indices_to_remove, indices[i])
			}
		}
	}

	for i := range(indices){
		if(!Contains(indices_to_remove,indices[i])){
			indices_result = append(indices_result,indices[i])
		}
	}

	return indices_result
}

// Used to increment a binary Value by one (the Value is contained in an int slice)
func Increment(A []int) ([]int) {
	carry := true;
	for i := len(A) - 1; i >= 0; i-- {
		if (carry) {
			if (A[i] == 0) {
				A[i] = 1
				carry = false
			} else {
				A[i] = 0;
				carry = true
			}
		}
	}

	return A;
}

// converts the elements of a slice of ints into a single string
// int[]{0,1,0,0} => "0100"
func SliceElementsToString(slice []int) string {
	var buffer bytes.Buffer
	for _, val := range slice {
		buffer.WriteString(strconv.Itoa(val))
	}
	return buffer.String()
}

type FileReaderByLine struct{
	file *os.File
	scanner *bufio.Scanner
}

func(f FileReaderByLine) ReadFrom(fileName string) *FileReaderByLine{
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	f.scanner = bufio.NewScanner(file)

	return &f
}

func(f *FileReaderByLine) NextLine() string{
	if f.scanner.Scan() {
		return f.scanner.Text()
	}else{
		f.file.Close()
		return "EOF"
	}
}

//Precondition: 's' is a binary string representation of a byte (max 8 characters)
func AddLeadingZerosToBitString(s string) string {
	missingZeros := 8-len(s)

	for i := 0; i < missingZeros; i++{
		s = "0" + s
	}

	return s
}

func Maj(i int, i2 int, i3 int) int {
	return i*i2 ^ i*i3 ^ i2*i3
}

// Bitwise XOR of 2 int slices
// Precondition: len(slice1) == len(slice2)
func XorRegisters(slice1 []int, slice2[]int) []int{
	var result []int

	for i := range(slice1){
		result = append(result, slice1[i]^slice2[i])
	}

	return result
}

// Tests equality of 2 int slices
func TestEq(a, b []int) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false;
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}


func ReadInitialStates(intial_states_filename string) InitialStates{
	fileReader := FileReaderByLine{}.ReadFrom("../"+intial_states_filename)
	r1_init_state := fileReader.NextLine()
	r2_init_state := fileReader.NextLine()
	r3_init_state := fileReader.NextLine()
	r4_init_state := fileReader.NextLine()

	r1_initial := [19]int{}
	r2_initial := [22]int{}
	r3_initial := [23]int{}
	r4_initial := [17]int{}
	r4_temp := make([]int, 0)

	for i,v := range(r1_init_state){
		r1_initial[i] = int(v - '0') // Awkward way of converting a rune to an int
	}

	for i,v := range(r2_init_state){
		r2_initial[i] = int(v - '0') // Awkward way of converting a rune to an int
	}

	for i,v := range(r3_init_state){
		r3_initial[i] = int(v - '0') // Awkward way of converting a rune to an int
	}

	r4_as_int,_ := strconv.Atoi(r4_init_state)

	r4_as_string_in_bits := strconv.FormatInt(int64(r4_as_int), 2)

	// Append missing leading zeros to R4
	for i := 0; i < 17-len(r4_as_string_in_bits); i++{
		r4_temp = append(r4_temp, 0)
	}

	for i := 0; i < len(r4_as_string_in_bits); i++{
		r4_temp = append(r4_temp, int(r4_as_string_in_bits[i] - '0'))
	}

	copy(r4_initial[:], r4_temp)

	states := InitialStates{
		R1: r1_initial,
		R2: r2_initial,
		R3: r3_initial,
		R4: r4_initial,
	}

	return states
}

// Precondition: len(a) = len(b)
// Returns: a byte array containing the bits of keystream
func XorBytesBitwise(a,b []byte) []int{
	result := make([]int,len(a)*8)

	for i := range(a){
		byte1_as_string_bits := strconv.FormatInt(int64(a[i]), 2)
		byte1_as_string_bits = AddLeadingZerosToBitString(byte1_as_string_bits)

		byte2_as_string_bits := strconv.FormatInt(int64(b[i]), 2)
		byte2_as_string_bits = AddLeadingZerosToBitString(byte2_as_string_bits)

		for j := range(byte1_as_string_bits){
			result[i*8+j] = int(byte1_as_string_bits[j] ^ byte2_as_string_bits[j])
		}

	}

	return result
}

type KeyGenerationStrategy interface {
	GenerateSessionKey() [64]int
}

type FixedSessionKeyStrategy struct {

}

type RandomSessionKeyStrategy struct {

}

func(f *FixedSessionKeyStrategy) GenerateSessionKey() [64]int{
	return [64]int{0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0}
}

func(r *RandomSessionKeyStrategy) GenerateSessionKey() [64]int{
	result := [64]int{}

	max := new(big.Int).SetInt64(1)

	for i := range(result){
		n,_ := rand.Int(rand.Reader,max)
		result[i] = int(n.Int64())
	}

	return result
}

func ThrowOnError(err error, msg string){
	if err != nil{
		fmt.Println(msg)
		panic(err)
	}
}

package a52

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	. "utils"
)

func TestMajorityFunction(t *testing.T){
	majresult := Maj(0,1,0)
	assert.Equal(t,int(0),majresult,"Error in majority function. Result is: ", majresult, " but should be: 0")

	majresult = Maj(1,1,0)
	assert.Equal(t,int(1),majresult,"Error in majority function. Result is: ", majresult, " but should be: 1")

	majresult = Maj(1,1,1)
	assert.Equal(t,int(1),majresult,"Error in majority function. Result is: ", majresult, " but should be: 1")

	majresult = Maj(0,0,1)
	assert.Equal(t,int(0),majresult,"Error in majority function. Result is: ", majresult, " but should be: 0")

	majresult = Maj(0,0,0)
	assert.Equal(t,int(0),majresult,"Error in majority function. Result is: ", majresult, " but should be: 0")

	majresult = Maj(1,0,1)
	assert.Equal(t,int(1),majresult,"Error in majority function. Result is: ", majresult, " but should be: 1")
}

func TestR1Clocking(t *testing.T) {
	r1 := R1{[19]int{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}}
	r1.clock()
	expected_result := [19]int{0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	assert.Equal(t,expected_result,r1.Register,"Clocking error. Clocked Register should be: ", expected_result, ", but is: ", r1.Register)

	r1.clock()
	r1.clock()
	expected_result = [19]int{0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}

	assert.Equal(t,expected_result,r1.Register,"Clocking error. Clocked Register should be: ", expected_result, ", but is: ", r1.Register)

	r1 = R1{[19]int{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}}
	r1.Register[13] = 1
	expected_result = [19]int{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	expected_result[14] = 1
	expected_result[0] = 1
	r1.clock()

	assert.Equal(t,expected_result,r1.Register,"Clocking error. Clocked Register should be: ", expected_result, ", but is: ", r1.Register)

	r1 = R1{[19]int{1,0,1,0,1,1,1,0,0,1,0,1,0,0,1,1,1,0,1}}
	r1.clock()
	expected_result = [19]int{0,1,0,1,0,1,1,1,0,0,1,0,1,0,0,1,1,1,0}
	assert.Equal(t,expected_result,r1.Register,"Clocking error. Clocked Register should be: ", expected_result, ", but is: ", r1.Register)
}

func TestR2Clocking(t *testing.T){
	r2 := R2{[22]int{1,0,1,0,1,1,1,0,0,1,0,1,0,0,1,1,1,0,1,0,0,1}}
	r2.clock()
	expected_result := [22]int{1,1,0,1,0,1,1,1,0,0,1,0,1,0,0,1,1,1,0,1,0,0}

	assert.Equal(t,expected_result,r2.Register,"Clocking error. Clocked Register should be: ", expected_result, ", but is: ", r2.Register)
}

func TestR3Clocking(t *testing.T){
	r3 := R3{[23]int{1,0,1,0,1,1,1,0,0,1,0,1,0,0,1,1,1,0,1,0,0,1,0}}
	r3.clock()
	expected_result := [23]int{1,1,0,1,0,1,1,1,0,0,1,0,1,0,0,1,1,1,0,1,0,0,1}

	assert.Equal(t,expected_result,r3.Register,"Clocking error. Clocked Register should be: ", expected_result, ", but is: ", r3.Register)
}

func TestR4Clocking(t *testing.T){
	r4 := R4{[17]int{1,0,1,0,1,1,1,0,0,1,0,1,0,1,1,1,0}}
	clockR4Only := false
	r4.clock(&A52{
		Session_key: [64]int{},
		R1:          R1{},
		R2:          R2{},
		R3:          R3{},
		R4:          R4{},
		UsedBitsOfCurrentCipherOutput: 0,
		OutputCache:                   nil,
		CurrentFrameNumber:            Frame{},
	},clockR4Only)
	expected_result := [17]int{1,1,0,1,0,1,1,1,0,0,1,0,1,0,1,1,1}

	assert.Equal(t,expected_result,r4.Register,"Clocking error. Clocked Register should be: ", expected_result, ", but is: ", r4.Register)
}

func TestGenerateCipherForFrame(t *testing.T){
	a52 := makeEmptyA52()
	a52.GenerateCipherForFrame([21]int{},false)
}

func TestIncrement(t *testing.T) {
	TestArray := [17]int{0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0}
	TestArrayPlusOne := Increment(TestArray[:])
	ExpectedArray := [17]int{0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,1}
	assert.Equal(t, true, TestEq(ExpectedArray[:], TestArrayPlusOne[:]))

	TestArrayPlusTwo := Increment(TestArrayPlusOne)

	ExpectedArray = [17]int{0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,1,0}
	assert.Equal(t, true, TestEq(TestArrayPlusTwo[:], ExpectedArray[:]))
}

func TestSliceElementsToString(t *testing.T) {
	TestArray := [5]int{0,0,1,0,0}
	actualString := SliceElementsToString(TestArray[:])
	expectedString := "00100"
	fmt.Println(expectedString)
	assert.Equal(t, actualString, expectedString)

}

func makeEmptyA52() A52{
	return A52{
		Session_key: [64]int{},
		R1:          R1{},
		R2:          R2{},
		R3:          R3{},
		R4:          R4{},
		UsedBitsOfCurrentCipherOutput: 0,
		OutputCache:                   nil,
		CurrentFrameNumber:            Frame{},
	}
}

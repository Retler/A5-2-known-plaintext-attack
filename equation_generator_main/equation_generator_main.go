package main

import (
	"fmt"
	"strconv"
	"strings"
	. "utils"
	. "a52_equation_generation"
	. "a52"
)

// Precondition: Theres exists a MySql database on localhost with user 'root', passowrd 'a52a52a52' and schema 'equationguesses'
// An SQL script is provided to set up the database (though not the user)
func GenerateEquationsForDatabase(){
	fmt.Println("Generating equations for database. This may take a while...")
	frame1 := Frame{0}.AsBinary()
	frame2 := Frame{1}.AsBinary()
	frame3 := Frame{2}.AsBinary()
	frame4 := Frame{3}.AsBinary()
	frame5 := Frame{4}.AsBinary()

	r4Guess := [17]int{} // We start guessing at 0
	fullR4 := 131072

	for j := 0; j<fullR4; j++{
		fmt.Println("Generating equations for R4 = " + strconv.Itoa(j))
		solutionsForFrame1 := make([]int,114)
		solutionsForFrame2 := make([]int,114)
		solutionsForFrame3 := make([]int,114)
		solutionsForFrame4 := make([]int,114)
		solutionsForFrame5 := make([]int,114)

		a52 := A52_equation_generator{}
		copy(a52.R4.Register[:],r4Guess[:])
		//r4Guess := [17]int{1,1,0,0,0,1,1,1,1,0,0,1,1,0,0,1,1} // let our guess start at 102100
		skipInitialize := true // We dont want to run the initialization phase when we guess R4. Instead we want to skip it and set R4 = guess, such that we can produce an equation system

		// Generate equations and solutions for frame 1
		outputFrame1 := a52.GenerateCipherForFrame(frame1, skipInitialize)
		equationsFrame1 := MakeEquationsFromOutput(outputFrame1)

		// Generate equations and solutions for frame 2
		copy(a52.R4.Register[:], XorRegisters(GetDiffOfFrames(frame1,frame2).R4Diff,r4Guess[:]))
		outputFrame2 := a52.GenerateCipherForFrame(frame2,skipInitialize)
		equationsFrame2 := MakeEquationsFromOutput(outputFrame2)
		TransformEquationsToMatchTheFirstFrame(GetDiffOfFrames(frame1,frame2),equationsFrame2[:],solutionsForFrame2)

		// Generate equations and solutions for frame 3
		copy(a52.R4.Register[:], XorRegisters(GetDiffOfFrames(frame1,frame3).R4Diff,r4Guess[:]))
		outputFrame3 := a52.GenerateCipherForFrame(frame3,skipInitialize)
		equationsFrame3 := MakeEquationsFromOutput(outputFrame3)
		TransformEquationsToMatchTheFirstFrame(GetDiffOfFrames(frame1,frame3),equationsFrame3[:],solutionsForFrame3[:])

		// Generate equations and solutions for frame 4
		copy(a52.R4.Register[:], XorRegisters(GetDiffOfFrames(frame1,frame4).R4Diff,r4Guess[:]))
		outputFrame4 := a52.GenerateCipherForFrame(frame4,skipInitialize)
		equationsFrame4 := MakeEquationsFromOutput(outputFrame4)
		TransformEquationsToMatchTheFirstFrame(GetDiffOfFrames(frame1,frame4),equationsFrame4[:],solutionsForFrame4[:])

		// Generate equations and solutions for frame 5
		copy(a52.R4.Register[:], XorRegisters(GetDiffOfFrames(frame1,frame5).R4Diff,r4Guess[:]))
		outputFrame5 := a52.GenerateCipherForFrame(frame5,skipInitialize)
		equationsFrame5 := MakeEquationsFromOutput(outputFrame5)
		TransformEquationsToMatchTheFirstFrame(GetDiffOfFrames(frame1,frame5),equationsFrame5[:],solutionsForFrame5[:])

		// Append solutions to equations (note that they keystream )
		for i := 0; i < 114; i++{
			equationsFrame1[i] = append(equationsFrame1[i], solutionsForFrame1[i])
			equationsFrame2[i] = append(equationsFrame2[i], solutionsForFrame2[i])
			equationsFrame3[i] = append(equationsFrame3[i], solutionsForFrame3[i])
			equationsFrame4[i] = append(equationsFrame4[i], solutionsForFrame4[i])
			equationsFrame5[i] = append(equationsFrame5[i], solutionsForFrame5[i])
		}

		// Append equations to one system
		equationSystem := append(equationsFrame1[:], equationsFrame2[:]...)
		equationSystem = append(equationSystem, equationsFrame3[:]...)
		equationSystem = append(equationSystem, equationsFrame4[:]...)
		equationSystem = append(equationSystem, equationsFrame5[:]...)

		// Build SQL-insert string and insert equations into database
		insertStatement := ""
		for i,equation := range(equationSystem){
			equationAsString := SliceElementsToString(equation)
			insertStatement += "( NULL,'" + equationAsString + "',"+ strconv.Itoa(j) + "," + strconv.Itoa(i) + "),"
		}
		insertStatement = strings.TrimRight(insertStatement, ",") // Trim trailing comma
		InsertEquations(insertStatement)
		copy(r4Guess[:],Increment(r4Guess[:]))
	}
}

func main(){
	fmt.Println("Generating equations for database.")
	fmt.Println("Make sure there is a user with credentials 'root:a52a52a52' that has access to schema 'EquationGuesses'")
	fmt.Println("Run the provided DB setup script to make a new database")
	fmt.Println("")

	GenerateEquationsForDatabase()
}

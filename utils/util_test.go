package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFileReader(t *testing.T){

	fileReader := FileReaderByLine{}.ReadFrom("../Initial states.txt")
	r1_init_state := fileReader.NextLine()
	r2_init_state := fileReader.NextLine()
	r3_init_state := fileReader.NextLine()
	r4_init_state := fileReader.NextLine()
	EOF := fileReader.NextLine()

	assert.Equal(t, "0001000010111001100", r1_init_state)
	assert.Equal(t, "1000010010000011100111", r2_init_state)
	assert.Equal(t, "00101111100011000111011", r3_init_state)
	assert.Equal(t, "102195", r4_init_state)
	assert.Equal(t, "EOF", EOF)
}

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	. "utils"
	. "a52"
)

// The function read a file with name "Initial states.txt" where states of R1,R2,R3 is listed in bits and R4 by decimal (values of R1,R2,R3,R4 are separated by newline)
func CrackCipherBasedOnInitialStates(file_name, initial_states_filename string){

	a52 := A52{
		UsedBitsOfCurrentCipherOutput: 0,
		OutputCache:                   nil,
		CurrentFrameNumber:            Frame{0} ,
		NextCipherbitStragegy:         nil,
		InitialStates: ReadInitialStates(initial_states_filename),
	}

	a52.NextCipherbitStragegy = &CustomInitialStatesNextCipherbitStrategy{A52: &a52}

	ciphertext,_ := ioutil.ReadFile("../"+file_name)

	plaintext := a52.EncryptDecryptFile(ciphertext)

	err := ioutil.WriteFile("Deciphered " + file_name, plaintext, 0644)
	ThrowOnError(err,"Could not write plaintext to file")

	fmt.Println("Succesfully deciphered " + file_name)
}

func ConvertKnownPlaintextToKeystream(){
	//Say we have the first 570 bits or approx 72 bytes of plaintext:
	known_plaintext,_:= ioutil.ReadFile("../Lorem ipsum known plaintext.txt")

	ciphertext,_ := ioutil.ReadFile("../Lorem ipsum cipher.txt")

	// By XOR'ing the known plaintext with the ciphertext we get the bits of the keystream
	keystream_as_bits := XorBytesBitwise(known_plaintext,ciphertext[:72])

	// Save 570 bits of keystream to be used by our solver
	f, _ := os.Create("../known keystream.txt")
	defer f.Close()
	f.WriteString(SliceElementsToString(keystream_as_bits[:570]))
}

func main(){

	// Checking command line arguments and providing documentation
	if  len(os.Args) != 3 && len(os.Args) != 4 {
		fmt.Println("Error, the program must run with 2 or 3 arguments.")
		fmt.Println("Argument 1:")
		fmt.Println("'encrypt': To encrypt the file")
		fmt.Println("'decrypt': To decrypt the file")
		fmt.Println()
		fmt.Println("Argument 2:")
		fmt.Println("<file name.txt>: The file to encrypt/decrypt (with extension).")
		fmt.Println()
		fmt.Println("Argument 3 (optional):")
		fmt.Println("<initial states filename.txt>: If decryption is to be done by using initial states. Default is decrypting using session key.")
		os.Exit(1)
	}

	enc_dec_flag := os.Args[1]
	file_name := os.Args[2]
	session_key_generator := FixedSessionKeyStrategy{}
	session_key := session_key_generator.GenerateSessionKey()

	if enc_dec_flag == "encrypt"{
		a52 := A52_builder{}.Build(session_key, KeyBasedGenerationStrategy)

		file,err := ioutil.ReadFile(file_name)

		ThrowOnError(err, "Could not read file " + file_name + ". Make sure you have read permission.")

		cipher_text := a52.EncryptDecryptFile(file)

		err = ioutil.WriteFile(file_name+" encrypted.txt", cipher_text, 0644)

		ThrowOnError(err, "Could not write " + file_name + " to disk. Make sure your console has write permission.")

		f, err := os.Create("session key.txt")

		ThrowOnError(err, "Could not create file for sesion key")

		_,err = f.WriteString(SliceElementsToString(a52.Session_key[:]))

		ThrowOnError(err, "Could not write the session key to the file")

		fmt.Println("")
	}

	if enc_dec_flag == "decrypt"{
		if len(os.Args) == 4{ // Then decrypt using initial states output by the solver
			init_states_filename := os.Args[3]
			CrackCipherBasedOnInitialStates(file_name, init_states_filename)

		}else{ // Else decrypt using the session key
			key, err := ioutil.ReadFile("session key.txt")
			ThrowOnError(err, "Could not read session key file. Make sure you have 'session key.txt in the same directory")
			file, err := ioutil.ReadFile(file_name)
			ThrowOnError(err, "Could not read session key file. Make sure you have 'session key.txt in the same directory'")

			if len(key) != 64{
				fmt.Println("The provided key is not a valid 64 bit key. Make sure there is no whitespace in the file.")
			}

			keystring := string(key)

			session_key := [64]int{}

			for i := range(keystring){
				session_key[i] = int(keystring[i] - '0')
			}

			a52 := A52_builder{}.Build(session_key,KeyBasedGenerationStrategy)

			plaintext := a52.EncryptDecryptFile(file)

			ioutil.WriteFile(file_name+" decrypted.txt", plaintext, 0644)
		}
	}
}

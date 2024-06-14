package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Defines the alphabet used in Enigma
const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Mapping a rotor
type Rotor struct {
	forwardMapping  string
	backwardMapping string
	position        int
	notch           byte
}

// Create a new rotor
func NewRotor(mapping string, notch byte) *Rotor {
	backwardMapping := make([]byte, 26)
	for i, char := range mapping {
		backwardMapping[char-'A'] = byte(i) + 'A'
	}
	return &Rotor{mapping, string(backwardMapping), 0, notch}
}

// Advance rotor position
func (r *Rotor) Rotate() bool {
	r.position = (r.position + 1) % 26
	return byte(r.position)+'A' == r.notch
}

// Passa o sinal através do rotor para frente
func (r *Rotor) Forward(c byte) byte {
	index := (int(c-'A') + r.position) % 26
	return (r.forwardMapping[index]-'A'-byte(r.position)+26)%26 + 'A'
}

// Passes the signal through the rotor forward
func (r *Rotor) Backward(c byte) byte {
	index := (int(c-'A') + r.position) % 26
	return (r.backwardMapping[index]-'A'-byte(r.position)+26)%26 + 'A'
}

// Reflector mapping
type Reflector struct {
	mapping string
}

// Create a new reflector
func NewReflector(mapping string) *Reflector {
	return &Reflector{mapping}
}

// Reflete o sinal
func (r *Reflector) Reflect(c byte) byte {
	return r.mapping[c-'A']
}

// Reflect the signal
type Plugboard struct {
	wiring map[byte]byte
}

// Create a new plugboard
func NewPlugboard(wiring map[byte]byte) *Plugboard {
	return &Plugboard{wiring}
}

// Maps the character through the plugboard
func (p *Plugboard) Swap(c byte) byte {
	if mapped, ok := p.wiring[c]; ok {
		return mapped
	}
	return c
}

// Enigma Machine
type Enigma struct {
	rotors    []*Rotor
	reflector *Reflector
	plugboard *Plugboard
}

// Create a new Enigma machine
func NewEnigma(rotors []*Rotor, reflector *Reflector, plugboard *Plugboard) *Enigma {
	return &Enigma{rotors, reflector, plugboard}
}

// Encrypt a message
func (e *Enigma) Encrypt(message string) string {
	var encryptedMessage strings.Builder
	for _, char := range strings.ToUpper(message) {
		if char < 'A' || char > 'Z' {
			continue
		}

		// Go through the plugboard
		c := e.plugboard.Swap(byte(char))

		// Pass through the rotors forward
		for i := len(e.rotors) - 1; i >= 0; i-- {
			c = e.rotors[i].Forward(c)
		}

		// Reflect the signal
		c = e.reflector.Reflect(c)

		// Pass through the rotors backwards
		for _, rotor := range e.rotors {
			c = rotor.Backward(c)
		}

		// Go through the plugboard again
		c = e.plugboard.Swap(c)

		// Add the encrypted character to the message
		encryptedMessage.WriteByte(c)

		// Rotate the rotors
		for i := len(e.rotors) - 1; i >= 0; i-- {
			if !e.rotors[i].Rotate() {
				break
			}
		}
	}

	return encryptedMessage.String()
}

func main() {
	// Define rotors and reflector
	rotor1 := NewRotor("EKMFLGDQVZNTOWYHXUSPAIBRCJ", 'Q')
	rotor2 := NewRotor("AJDKSIRUXBLHWTMCQGZNPYFVOE", 'E')
	rotor3 := NewRotor("BDFHJLCPRTXVZNYEIWGAKMUSQO", 'V')
	reflector := NewReflector("YRUHQSLDPXNGOKMIEBFZCWVJAT")

	// Define plugboard
	plugboardWiring := map[byte]byte{
		'A': 'B', 'B': 'A', 'C': 'D', 'D': 'C',
	}
	plugboard := NewPlugboard(plugboardWiring)

	// Create Enigma machine
	enigma := NewEnigma([]*Rotor{rotor1, rotor2, rotor3}, reflector, plugboard)

	fmt.Println(`
	###############################################
	#                                             #
	#                    ENIGMA                   #
	#                                             #
	###############################################
	`)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Type '1' to encrypt or '2' to decrypt (or 'exit' to finish): ")
		if !scanner.Scan() {
			break
		}
		option := scanner.Text()
		if strings.ToLower(option) == "exit" {
			break
		}
		if option != "1" && option != "2" {
			fmt.Println("Option invalid, try (1 or 2)")
			continue
		}

		fmt.Print("Enter a message: ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		// Resets the position of the rotors before each operation
		rotor1.position, rotor2.position, rotor3.position = 0, 0, 0

		if option == "1" {
			encryptedMessage := enigma.Encrypt(input)
			fmt.Printf("Encrypted message: %s\n", encryptedMessage)
		} else {
			decryptedMessage := enigma.Encrypt(input)
			fmt.Printf("Decrypted message: %s\n", decryptedMessage)
		}
	}

	fmt.Println("Program closed.")
}

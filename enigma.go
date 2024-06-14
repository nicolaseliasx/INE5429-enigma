package main

import (
	"fmt"
	"strings"
)

// Define o alfabeto usado na Enigma
const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Mapeamento de um rotor
type Rotor struct {
	forwardMapping  string
	backwardMapping string
	position        int
	notch           byte
}

// Cria um novo rotor
func NewRotor(mapping string, notch byte) *Rotor {
	backwardMapping := make([]byte, 26)
	for i, char := range mapping {
		backwardMapping[char-'A'] = byte(i) + 'A'
	}
	return &Rotor{mapping, string(backwardMapping), 0, notch}
}

// Avança a posição do rotor
func (r *Rotor) Rotate() bool {
	r.position = (r.position + 1) % 26
	return byte(r.position)+'A' == r.notch
}

// Passa o sinal através do rotor para frente
func (r *Rotor) Forward(c byte) byte {
	index := (int(c-'A') + r.position) % 26
	return (r.forwardMapping[index]-'A'-byte(r.position)+26)%26 + 'A'
}

// Passa o sinal através do rotor para trás
func (r *Rotor) Backward(c byte) byte {
	index := (int(c-'A') + r.position) % 26
	return (r.backwardMapping[index]-'A'-byte(r.position)+26)%26 + 'A'
}

// Mapeamento do refletor
type Reflector struct {
	mapping string
}

// Cria um novo refletor
func NewReflector(mapping string) *Reflector {
	return &Reflector{mapping}
}

// Reflete o sinal
func (r *Reflector) Reflect(c byte) byte {
	return r.mapping[c-'A']
}

// Plugboard da Enigma
type Plugboard struct {
	wiring map[byte]byte
}

// Cria um novo plugboard
func NewPlugboard(wiring map[byte]byte) *Plugboard {
	return &Plugboard{wiring}
}

// Mapeia o caractere através do plugboard
func (p *Plugboard) Swap(c byte) byte {
	if mapped, ok := p.wiring[c]; ok {
		return mapped
	}
	return c
}

// Máquina Enigma
type Enigma struct {
	rotors     []*Rotor
	reflector  *Reflector
	plugboard  *Plugboard
}

// Cria uma nova máquina Enigma
func NewEnigma(rotors []*Rotor, reflector *Reflector, plugboard *Plugboard) *Enigma {
	return &Enigma{rotors, reflector, plugboard}
}

// Encripta uma mensagem
func (e *Enigma) Encrypt(message string) string {
	var encryptedMessage strings.Builder
	for _, char := range strings.ToUpper(message) {
		if char < 'A' || char > 'Z' {
			continue
		}

		// Passa pelo plugboard
		c := e.plugboard.Swap(byte(char))

		// Passa pelos rotores para frente
		for i := len(e.rotors) - 1; i >= 0; i-- {
			c = e.rotors[i].Forward(c)
		}

		// Reflete o sinal
		c = e.reflector.Reflect(c)

		// Passa pelos rotores para trás
		for _, rotor := range e.rotors {
			c = rotor.Backward(c)
		}

		// Passa pelo plugboard novamente
		c = e.plugboard.Swap(c)

		// Adiciona o caractere criptografado à mensagem
		encryptedMessage.WriteByte(c)

		// Roda os rotores
		for i := len(e.rotors) - 1; i >= 0; i-- {
			if !e.rotors[i].Rotate() {
				break
			}
		}
	}

	return encryptedMessage.String()
}

func main() {
	// Definir rotores e refletor
	rotor1 := NewRotor("EKMFLGDQVZNTOWYHXUSPAIBRCJ", 'Q')
	rotor2 := NewRotor("AJDKSIRUXBLHWTMCQGZNPYFVOE", 'E')
	rotor3 := NewRotor("BDFHJLCPRTXVZNYEIWGAKMUSQO", 'V')
	reflector := NewReflector("YRUHQSLDPXNGOKMIEBFZCWVJAT")

	// Definir plugboard
	plugboardWiring := map[byte]byte{
		'A': 'B', 'B': 'A', 'C': 'D', 'D': 'C', // Exemplo de configuração de plugboard
	}
	plugboard := NewPlugboard(plugboardWiring)

	// Criar máquina Enigma
	enigma := NewEnigma([]*Rotor{rotor1, rotor2, rotor3}, reflector, plugboard)

	// Mensagem de exemplo
	message := "HELLOENIGMA"
	encryptedMessage := enigma.Encrypt(message)

	fmt.Printf("Mensagem original: %s\n", message)
	fmt.Printf("Mensagem criptografada: %s\n", encryptedMessage)
}

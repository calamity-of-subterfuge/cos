package client_test

import (
	"fmt"
	"math/rand"

	"github.com/calamity-of-subterfuge/cos/pkg/client"
)

func ExampleTypoify() {
	rand.Seed(74)
	og := "this is a test! look at me go: I'm doing great"
	fmt.Printf(" %q\n", og)
	fmt.Printf("%q", client.Typoify(og, 5))
	// Output:
	//  "this is a test! look at me go: I'm doing great"
	// ["this is a t3st! look at me go: I\"mdoibggreat"]
}

func ExampleTypoify_diffseed() {
	rand.Seed(5)
	og := "this is a test! look at me go: I'm doing great"
	fmt.Printf(" %q\n", og)
	fmt.Printf("%q", client.Typoify(og, 0.5))
	// Output:
	//  "this is a test! look at me go: I'm doing great"
	// ["this is a tesf! look at me go: I'm doing great"]
}

func ExampleTypoify_diffseed2() {
	rand.Seed(67)
	og := "this is a test! look at me go: I'm doing great"
	fmt.Printf(" %q\n", og)
	fmt.Printf("%q", client.Typoify(og, 0.5))
	// Output:
	//  "this is a test! look at me go: I'm doing great"
	// ["this is at est! look at me go: I'm doing great"]
}

func ExampleTypoify_diffseed3() {
	rand.Seed(15)
	og := "this is a test! look at me go: I'm doing great"
	fmt.Printf(" %q\n", og)
	fmt.Printf("%q", client.Typoify(og, 1))
	// Output:
	//  "this is a test! look at me go: I'm doing great"
	// ["this is a test! look at me go? I'm doing great"]
}

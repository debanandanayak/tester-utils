package internal

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var randInstance = rand.New(rand.NewSource(time.Now().UnixNano()))

var randomWords = []string{
	"apple",
	"orange",
	"banana",
	"pear",
	"grape",
	"pineapple",
	"mango",
	"strawberry",
	"raspberry",
	"blueberry",
}

// Init must be called at the start of every program.
//
// If CODECRAFTERS_RANDOM_SEED is set, it will be used to generate predictable random numbers.
func Init() {
	if seed := os.Getenv("CODECRAFTERS_RANDOM_SEED"); seed != "" {
		seedInt, err := strconv.Atoi(seed)
		if err != nil {
			panic(err)
		}

		randInstance.Seed(int64(seedInt))
	} else {
		randInstance.Seed(time.Now().UnixNano())
	}
}

// RandomWord returns a random word from the list of words.
func RandomWord() string {
	return randomWords[randInstance.Intn(len(randomWords))]
}

// RandomWords returns a random list of n words.
func RandomWords(n int) []string {
	return shuffle(randomWords, n)
}

// RandomString returns a random string of 6 words.
func RandomString() string {
	return strings.Join(RandomWords(6), " ")
}

// RandomStrings returns a random list of n strings.
func RandomStrings(n int) []string {
	l := make([]string, n)

	for i := range l {
		l[i] = RandomString()
	}

	return l
}

func shuffle(vals []string, n int) []string {
	if n > len(vals) {
		panic("don't have so many words")
	}

	ret := make([]string, n)

	for i, randIndex := range randInstance.Perm(len(vals))[:n] {
		ret[i] = vals[randIndex]
	}

	return ret
}

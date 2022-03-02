package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PatrikOlin/gordle/api"
)

func rules() {
	fmt.Println(`
	Rules:
		- You have 5 guesses
		- You can only guess using 5 letters words
		- You can only guess using lowercase letters
		- Green means correct at the right position
		- Yellow means correct at the wrong position
	Give a filename in argument to have more words

	Good luck!
	`)
}

func getWord() string {
	var (
		words = []string{
			"trove",
			"nasty",
			"tasty",
			"hasty",
			"place",
			"right",
			"black",
			"white",
			"happy",
			"thing",
			"child",
			"night",
			"world",
		}
		word string
		i    int
	)
	if len(os.Args) > 1 {
		words, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		buff := strings.Split(string(words), "\n")
		rand.Seed(time.Now().UnixNano())
		i = rand.Intn(len(buff))
		word = buff[i]
		if len(word) != 5 {
			fmt.Println("Error: the file must contain only five letter words")
			os.Exit(1)
		}
	} else {
		rand.Seed(time.Now().UnixNano())
		i = rand.Intn(len(words))
		word = words[i]
	}

	return word
}

func game() {
	var (
		word   string = getWord()
		guess  string
		count  int    = 0
		end    bool   = false
		yellow string = "\033[33m"
		green  string = "\033[32m"
		reset  string = "\033[0m"
	)
	for !end {
		fmt.Print("> ")
		fmt.Scanf("%s", &guess)
		if len(guess) != 5 {
			fmt.Println("Please enter 5 letters")
			continue
		}
		for i := 0; i < len(word); i++ {
			if word[i] == guess[i] {
				fmt.Print(green, string(guess[i]), green)
			} else if strings.ContainsAny(word, string(guess[i])) {
				fmt.Print(yellow, string(guess[i]), yellow)
			} else {
				fmt.Print(reset, string(guess[i]))
			}
		}
		fmt.Println()
		if word == guess {
			fmt.Println(green + "You won!" + reset)
			end = true
		} else if count == 5 {
			fmt.Println("\033[31m" + "You lost!" + reset)
			fmt.Println("The word was " + word)
			end = true
		}

		count++
	}
}

func main() {
	// rules()
	// game()
	// fmt.Println("Thanks for playing")

	r := api.GetRouter()

	http.ListenAndServe(":4040", r)
}

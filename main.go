package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var n_letter = [][]int{
	{0, 0, 0, 0, 1, 1, 1},
	{0, 0, 0, 0, 1, 0, 0},
	{0, 0, 0, 0, 1, 1, 1},
}

var s_letter = [][]int{
	{0, 0, 1, 1, 1, 0, 1},
	{0, 0, 1, 0, 1, 1, 1},
}

var e_letter = [][]int{
	{0, 0, 1, 1, 1, 1, 1},
	{0, 0, 1, 0, 1, 0, 1},
}

var d_letter = [][]int{
	{0, 0, 0, 1, 1, 1, 1},
	{0, 0, 0, 1, 0, 0, 1},
	{0, 1, 1, 1, 1, 1, 1},
}

var u_letter = [][]int{
	{0, 0, 0, 0, 1, 1, 1},
	{0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 1, 1, 1},
}

func main() {
	fmt.Println("Github nudes project is starting...")

	// load .env file from given path
	// we keep it empty it will load .env from current directory
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//Checking that an environment variable is present or not.
	githubToken, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		fmt.Println("GITHUB_TOKEN is not present")
	}

	fmt.Println(githubToken)
	fmt.Println(u_letter)
}

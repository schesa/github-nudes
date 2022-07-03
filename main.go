package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var ctx = context.Background()
var client *github.Client

// this should be on a sunday to be aligned with the letters format
var start_date = time.Date(2022, time.Month(7), 3, 0, 0, 0, 0, time.UTC)

var (
	sourceOwner   = flag.String("source-owner", "schesa", "Name of the owner (user or org) of the repo to create the commit in.")
	sourceRepo    = flag.String("source-repo", "github-nudes", "Name of repo to create the commit in.")
	commitMessage = flag.String("commit-message", "Auto-Generated Commit", "Content of the commit message.")
	baseBranch    = flag.String("base-branch", "main", "Name of branch to create the `commit-branch` from.")
	changedFile   = flag.String("files", "diff.log", `File to commit`)
	authorName    = flag.String("author-name", "Sebastian Chesa", "Name of the author of the commit.")
	authorEmail   = flag.String("author-email", "chesasebastian1997@gmail.com", "Email of the author of the commit.")
)

var n_letter = [][]string{
	{"0", "0", "0", "0", "1", "1", "1"},
	{"0", "0", "0", "0", "1", "0", "0"},
	{"0", "0", "0", "0", "1", "1", "1"},
}

var s_letter = [][]string{
	{"0", "0", "1", "0", "1", "1", "1"},
	{"0", "0", "1", "1", "1", "0", "1"},
}

var e_letter = [][]string{
	{"0", "0", "1", "0", "1", "0", "1"},
	{"0", "0", "1", "1", "1", "1", "1"},
}

var d_letter = [][]string{
	{"0", "1", "1", "1", "1", "1", "1"},
	{"0", "0", "0", "1", "0", "0", "1"},
	{"0", "0", "0", "1", "1", "1", "1"},
}

var u_letter = [][]string{
	{"0", "0", "0", "0", "1", "1", "1"},
	{"0", "0", "0", "0", "0", "0", "1"},
	{"0", "0", "0", "0", "1", "1", "1"},
}

var empty_space = [][]string{
	{"0", "0", "0", "0", "0", "0", "0"},
}

func main() {
	log.Println("Github nudes project is starting...")

	// load .env file from given path
	// we keep it empty it will load .env from current directory
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	action_matrix := getData()
	// generate_commit_csv()
	action := getTodaysAction(action_matrix)
	if action == "1" {
		for i := 1; i < 7; i++ {
			logToFile()
			commit()
		}
	} else {
		log.Println("nothing to do")
	}
}

func getTodaysAction(action_matrix [][]string) string {
	current_date := time.Now()
	days := current_date.Sub(start_date).Hours() / 24
	log.Println(math.Floor(days))
	week := math.Floor(math.Floor(days) / 7)
	day := math.Mod(math.Floor(days), 7)
	log.Println("week")
	log.Println(week)
	log.Println(day)
	log.Println("Action")
	log.Println(action_matrix[int(week)][int(day)])
	return action_matrix[int(week)][int(day)]
}

// This is another solution: to use csv file to keep track of the actions
// func generate_commit_csv() {
// 	file, err := os.Create("actions.csv")
// 	if err != nil {
// 		log.Fatalln("failed to create actions.csv file", err)
// 	}
// 	defer file.Close()
// 	csvwriter := csv.NewWriter(file)
// 	defer csvwriter.Flush()

// 	// Using WriteAll
// 	var data [][]string = getData()

// 	csvwriter.WriteAll(data)
// }

func getData() [][]string {
	x := [][]string{}
	x = append(s_letter, empty_space...) // Can't concatenate more than 2 slice at once
	x = append(x, e_letter...)
	x = append(x, empty_space...)
	x = append(x, n_letter...)
	x = append(x, empty_space...)
	x = append(x, d_letter...)
	x = append(x, empty_space...)
	x = append(x, n_letter...)
	x = append(x, empty_space...)
	x = append(x, u_letter...)
	x = append(x, empty_space...)
	x = append(x, d_letter...)
	x = append(x, empty_space...)
	x = append(x, e_letter...)
	x = append(x, empty_space...)
	x = append(x, s_letter...)

	fmt.Printf("\n######### After Concatenation #########\n")
	for i := 0; i < len(x); i++ {
		for j := 0; j < 7; j++ {
			fmt.Printf("%s ", x[i][j])
		}
		fmt.Printf("\n")
	}
	return x
}

func commit() {
	//Checking that an environment variable is present or not.
	token, ok := os.LookupEnv("GITHUB_AUTH_TOKEN")
	if !ok {
		log.Println("GITHUB_AUTH_TOKEN is not present")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
	ref, refErr := getRef()
	if refErr != nil {
		log.Fatalf("Error loading branch ref")
		return
	}
	// log.Println(ref)
	// tree, _, err := client.Git.GetTree(ctx, *sourceOwner, *sourceRepo, *ref.Object.SHA, true)
	// log.Println(tree)
	// log.Println(err)

	entries := []github.TreeEntry{}
	// Load each file into the tree.
	file, content, err := getFileContent("diff.log")
	log.Println(err)
	log.Println("file")
	log.Println(file)
	log.Println(content)

	entries = append(entries, github.TreeEntry{Path: github.String(file), Type: github.String("blob"), Content: github.String(string(content)), Mode: github.String("100644")})
	log.Println(entries)

	tree, _, err := client.Git.CreateTree(ctx, *sourceOwner, *sourceRepo, *ref.Object.SHA, entries)
	log.Println("tree")
	log.Println(tree)
	log.Println(err)

	err = pushCommit(ref, tree)
	log.Println(err)
}

// pushCommit creates the commit in the given reference using the given tree.
func pushCommit(ref *github.Reference, tree *github.Tree) (err error) {
	// Get the parent commit to attach the commit to.
	parent, _, err := client.Repositories.GetCommit(ctx, *sourceOwner, *sourceRepo, *ref.Object.SHA)
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()

	author := &github.CommitAuthor{Date: &date, Name: authorName, Email: authorEmail}
	commit := &github.Commit{Author: author, Message: commitMessage, Tree: tree, Parents: []github.Commit{*parent.Commit}}
	newCommit, _, err := client.Git.CreateCommit(ctx, *sourceOwner, *sourceRepo, commit)
	if err != nil {
		return err
	}

	// Attach the commit to the master branch.
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, *sourceOwner, *sourceRepo, ref, false)
	return err
}

// getFileContent loads the local content of a file and return the target name
// of the file in the target repository and its contents.
func getFileContent(fileArg string) (targetName string, b []byte, err error) {
	var localFile string
	files := strings.Split(fileArg, ":")
	switch {
	case len(files) < 1:
		return "", nil, errors.New("empty `-files` parameter")
	case len(files) == 1:
		localFile = files[0]
		targetName = files[0]
	default:
		localFile = files[0]
		targetName = files[1]
	}

	b, err = ioutil.ReadFile(localFile)
	return targetName, b, err
}

// getRef returns the commit branch reference object if it exists or creates it
// from the base branch before returning it.
func getRef() (ref *github.Reference, err error) {
	ref, _, err = client.Git.GetRef(ctx, *sourceOwner, *sourceRepo, "refs/heads/"+*baseBranch)
	return ref, err
}

func logToFile() {
	f, err := os.OpenFile(*changedFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("%s\n", time.Now().String())); err != nil {
		log.Println(err)
	}
}

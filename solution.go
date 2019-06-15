package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func findSolution(args []string) ([]string, error) {
	result := []string{}
	for _, name := range args {
		if strings.HasSuffix(strings.ToLower(name), ".sln") {
			result = append(result, name)
		}
	}
	if len(result) > 0 {
		return result, nil
	}
	fd, err := os.Open(".")
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	files, err := fd.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, file1 := range files {
		if strings.HasSuffix(strings.ToLower(file1.Name()), ".sln") {
			result = append(result, file1.Name())
		}
	}
	return result, nil
}

func FindSolution(args []string) (string, error) {
	sln, err := findSolution(args)
	if err != nil {
		return "", err
	}
	if len(sln) < 1 {
		return "", errors.New("no solution files")
	}
	if len(sln) >= 2 {
		return "", fmt.Errorf("%s: too may solution files", strings.Join(sln, ", "))
	}
	return sln[0], nil
}

type Solution struct {
	Path          string
	Version       string
	Configuration []string
}

func NewSolution(fname string) (*Solution, error) {
	fd, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	sln := &Solution{Path: fname}

	var block func([]string)
	block = func(f []string) {
		if f[0] == "#" && f[1] == "Visual" && f[2] == "Studio" && len(f) >= 4 {
			sln.Version = f[3]
		} else if f[0] == "GlobalSection(SolutionConfigurationPlatforms)" {
			save := block
			block = func(f []string) {
				if f[0] == "EndGlobalSection" {
					block = save
				} else {
					sln.Configuration = append(sln.Configuration, f[0])
				}
			}
		}
	}

	sc := bufio.NewScanner(fd)
	for sc.Scan() {
		block(strings.Fields(sc.Text()))
	}
	return sln, nil
}
package main

import (
	"strings"
)

func badwordCheck(s string) string {
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	splitString := strings.Split(s, " ")

	for i, word := range splitString {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				splitString[i] = "****"
			}
		}
	}
	return strings.Join(splitString, " ")
}

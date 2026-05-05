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
	newString := []string{}

	for _, word := range splitString {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				word = "****"
			}
		}
		newString = append(newString, word)
	}
	return strings.Join(newString, " ")
}

package domain

import "strings"

type Text struct {
	Message string
}

type MsgStatistics struct {
	LettersCount int
}

func CountingLetters(t Text) MsgStatistics {
	var count int
	for _, r := range strings.ToLower(t.Message) {
		if ('a' < r && r < 'z') || ('а' < r && 'я' > r) {
			count++
		}
	}
	return MsgStatistics{LettersCount: count}
}

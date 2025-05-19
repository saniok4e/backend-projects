package utils

import (
	"fmt"
	"strings"
)

// LevenshteinDistance вычисляет расстояние Левенштейна между двумя строками
func LevenshteinDistance(s1, s2 string) int {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	m := make([][]int, len(s1)+1)
	for i := range m {
		m[i] = make([]int, len(s2)+1)
		m[i][0] = i
	}
	for j := range m[0] {
		m[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			if s1[i-1] == s2[j-1] {
				m[i][j] = m[i-1][j-1]
			} else {
				m[i][j] = min(m[i-1][j-1]+1, min(m[i-1][j]+1, m[i][j-1]+1))
			}
		}
	}
	return m[len(s1)][len(s2)]
}

// Min возвращает минимальное из двух чисел
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// FormatScoreWord возвращает правильное склонение слова "балл"
func FormatScoreWord(score int) string {
	lastDigit := score % 10
	lastTwoDigits := score % 100

	// Исключение для чисел 11-19
	if lastTwoDigits >= 11 && lastTwoDigits <= 19 {
		return "баллов"
	}

	// Склонение по последней цифре
	switch lastDigit {
	case 1:
		return "балл"
	case 2, 3, 4:
		return "балла"
	default:
		return "баллов"
	}
}

// FormatScore форматирует сообщение о счете
func FormatScore(score int) string {
	return fmt.Sprintf("У тебя %d %s", score, FormatScoreWord(score))
}

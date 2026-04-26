package main

import (
	"unicode"
)

func splitWhitespaceN(s string, n int) []string {
	var (
		fields []string = []string{""}
		field  *string  = &fields[0]
	)

	for _, chr := range s {
		if unicode.IsSpace(chr) {
			if *field != "" {
				if n > 0 && len(fields) >= n {
					return fields
				}

				fields = append(fields, "")
				field = &fields[len(fields)-1]
			}
		} else {
			*field += string(chr)
		}
	}

	if *field == "" {
		return fields[:len(fields)]
	} else {
		return fields
	}
}

func getWord(s string, n int) string {
	var (
		word_i     int
		word_start int = -1
		word_end   int = -1
	)

	for o, chr := range s {
		if unicode.IsSpace(chr) {
			if word_i >= n {
				break
			} else {
				word_start = -1
				word_end = -1
				word_i++
			}

		} else {
			if word_start == -1 {
				word_start = o
			}

			word_end = o
		}
	}

	return s[word_start : word_end+1]
}

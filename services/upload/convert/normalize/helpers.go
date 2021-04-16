package normalize

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

// Atoi32 converts a string directly to a int32, baseline golang parses string always into int64 and have to be converted
// to int32. You can however transform a string easily to int, which is somehow the same, but the parquet package expects int32
// specifically
func Atoi32(input string) (int32, error) {
	bigint, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		return 0, err
	}

	num := int32(bigint)
	return num, nil
}

// Atoi64 is just a small wrapper around ParseInt
func Atoi64(input string) (int64, error) {
	num, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, err
	}

	return num, nil
}

// TODO: check that it is surrounded by quotes and fail otherwise
//		to make it fail early.
//		Because the columns must be surrounded by quotes otherwise it is a wrong column
// 		not sure that is the case everywhere I use this function
func trimQuotes(input string) string {
	output := strings.TrimSuffix(input, "\"")
	output = strings.TrimPrefix(output, "\"")
	return output
}

// TODO: test and fix the year problem
// this will break during new year, because go assumes UTC,
// but the combatlog has the time of the player afaik
func parseTimestamp(input *string) (*string, error) {
	input = aws.String(fmt.Sprintf("%v/%s", time.Now().Year(), *input))
	stupidTimeFormat := "2006/1/2 15:04:05.000"
	t, err := time.Parse(stupidTimeFormat, *input)
	if err != nil {
		return aws.String(""), fmt.Errorf("failed to parse time: %v", err)
	}

	timeAsInt64 := t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	timeAsString := fmt.Sprintf("%s", strconv.FormatInt(timeAsInt64, 10))
	return aws.String(timeAsString), nil
}

// copying code from stackoverflow like a pro
// https://stackoverflow.com/questions/59297737/go-split-string-by-comma-but-ignore-comma-within-double-quotes
// at least I added tests^^ and switched to string pointers to reduce memory
func splitAtCommas(s *string) []string {
	var res []string
	var beg int
	var inString bool

	for i := 0; i < len(*s); i++ {
		if (*s)[i] == ',' && !inString {
			res = append(res, (*s)[beg:i])
			beg = i + 1
		} else if (*s)[i] == '"' {
			if !inString {
				inString = true
			} else if i > 0 && (*s)[i-1] != '\\' {
				inString = false
			}
		}
	}
	return append(res, (*s)[beg:])
}

// wrapper around strings.Split to test functionality
// doesn't really make sense to test the standard library, but
// i want to refactor this later
func splitString(s, sep string) []string {
	return strings.Split(s, sep)
}

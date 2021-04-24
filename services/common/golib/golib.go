package golib

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// MinSecToMilliseconds converts time in the "minute:seconds" format to milliseconds
func MinSecToMilliseconds(input string) (int64, error) {
	input = fmt.Sprintf("1970 %s", input)
	t, err := time.Parse("2006 04:05", input)
	if err != nil {
		return 0, fmt.Errorf("failed to parse time input: %v", err)
	}
	milliseconds := t.UnixNano() / 1e6
	return milliseconds, nil
}

// TimedAsPercent determines if a key was intime, deplete, two chest or three chest and
// returns the quotient of the expected intime duration and the actual duration in milliseconds
// this is used to order the keys, so the fastest within a key level is first
// TODO: add table test and check 0 duration
func TimedAsPercent(dungeonID int, durationInMilliseconds float64) (durAsPercent float64, intime int, err error) {
	var intimeDuration, twoChestDuration, threeChestDuration float64

	if durationInMilliseconds == 0 {
		// if a key is abandoned the duration will always be 0
		return float64(1000), 0, nil
	}

	switch dungeonID {
	case 2291: // De Other Side
		ms, err := MinSecToMilliseconds("43:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = MinSecToMilliseconds("34:25")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = MinSecToMilliseconds("25:49")
		if err != nil {
			return 0, 0, err
		}
		threeChestDuration = float64(ms)

	case 2289: // Plaguefall
		ms, err := MinSecToMilliseconds("38:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = MinSecToMilliseconds("30:24")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = MinSecToMilliseconds("22:48")
		if err != nil {
			return 0, 0, err
		}
		threeChestDuration = float64(ms)
	case 2284: // Sanguine Depths
		ms, err := MinSecToMilliseconds("41:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = MinSecToMilliseconds("32:48")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = MinSecToMilliseconds("24:36")
		if err != nil {
			return 0, 0, err
		}
	// TODO: parse time and convert to milli seconds do it in TDD
	/*
		https://www.wowhead.com/mythic-keystones-and-dungeons-guide
		Dungeon	Timer	+2	+3
		De Other Side	43:00	34:25	25:49
		Plaguefall	38:00	30:24	22:48
		Halls of Atonement	31:00	24:48	18:36
		Mists of Tirna Scithe	30:00	24:00	18:00
		Spires of Ascension	39:00	31:12	23:24
		Sanguine Depths	41:00	32:48	24:36
		Necrotic Wake	36:00	28:48	21:36
		Theater of Pain	37:00	29:36	22:12
	*/
	case 2287: // Halls of Atonement
		ms, err := MinSecToMilliseconds("31:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = MinSecToMilliseconds("24:48")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = MinSecToMilliseconds("18:36")
		if err != nil {
			return 0, 0, err
		}
		threeChestDuration = float64(ms)

	case 2290: // Mists of Tirna Scithe
		ms, err := MinSecToMilliseconds("30:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = MinSecToMilliseconds("24:00")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = MinSecToMilliseconds("18:00")
		if err != nil {
			return 0, 0, err
		}
		threeChestDuration = float64(ms)

	case 2285: // Spires of Ascension
		ms, err := MinSecToMilliseconds("39:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = MinSecToMilliseconds("31:12")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = MinSecToMilliseconds("23:24")
		if err != nil {
			return 0, 0, err
		}
		threeChestDuration = float64(ms)

	case 2286: // Necrotic Wake
		ms, err := MinSecToMilliseconds("36:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = MinSecToMilliseconds("28:48")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = MinSecToMilliseconds("21:36")
		if err != nil {
			return 0, 0, err
		}
		threeChestDuration = float64(ms)

	case 2293: // Theater of Pain
		ms, err := MinSecToMilliseconds("37:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = MinSecToMilliseconds("29:36")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = MinSecToMilliseconds("22:12")
		if err != nil {
			return 0, 0, err
		}
		threeChestDuration = float64(ms)
	}
	intime = timed(durationInMilliseconds, intimeDuration, twoChestDuration, threeChestDuration)

	return durationAsPercent(intimeDuration, durationInMilliseconds), intime, err
}

// timed determines if a key was intime, deplete, two chest or three chest
func timed(durationInMilliseconds, intimeDuration, twoChestDuration, threeChestDuration float64) int {
	if durationInMilliseconds <= threeChestDuration {
		return 3 // three chest
	} else if durationInMilliseconds > threeChestDuration && durationInMilliseconds <= twoChestDuration {
		return 2 // two chest
	} else if durationInMilliseconds > twoChestDuration && durationInMilliseconds <= intimeDuration {
		return 1 // timed
	} else {
		return 0 // deplete
	}
}

// durationAsPercent returns the quotient of the expected intime duration and the actual duration in milliseconds
// this is used to order the keys, so the fastest within a key level is first
func durationAsPercent(dungeonIntimeDuration, durationInMilliseconds float64) float64 {
	return (dungeonIntimeDuration / durationInMilliseconds) * 100
}

// InitLogging sets up the logging for every lambda and should be called before the handler
func InitLogging() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

// CanonicalLog writes a structured message to stdout if the log level is atleast INFO
func CanonicalLog(msg map[string]interface{}) {
	logrus.WithFields(msg).Info()
}

// PrettyStruct converts a struct, slice of a struct or map into a readable string
func PrettyStruct(input interface{}) (string, error) {
	prettyJSON, err := json.MarshalIndent(input, "", "    ")
	if err != nil {
		return "", err
	}
	return string(prettyJSON), nil
}

// Atoi64 is just a small wrapper around ParseInt
func Atoi64(input string) (int64, error) {
	num, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, err
	}

	return num, nil
}

var levelTwoAffixes = map[int]string{
	9:  "Tyrannical",
	10: "Fortified",
}

var levelFourAffixes = map[int]string{
	123: "Spiteful",
	7:   "Bolstering",
	11:  "Bursting",
	8:   "Sanguine",
	6:   "Raging",
	122: "Inspiring",
}

var levelSevenAffixes = map[int]string{
	13:  "Explosive",
	4:   "Necrotic",
	3:   "Volcanic",
	124: "Storming",
	14:  "Quaking",
	12:  "Grievous",
}

var levelTenAffixes = map[int]string{
	121: "Prideful",
}

// AffixIDsToString takes an array of affix ids and converts them to a readable list of array names
// separated by commas
// TODO:
// 	- rethink if I actually need this, displaying a list of affixes as string takes a lot of space,
// 	should just return an array of ids and string and display icons and display the name on hover
func AffixIDsToString(levelTwoID, levelFourID, levelSevenID, levelTenID int) string {
	affixes := levelTwoAffixes[levelTwoID]

	if levelFourID == 0 {
		return affixes
	}

	affixes += ", " + levelFourAffixes[levelFourID]

	if levelSevenID == 0 {
		return affixes
	}

	affixes += ", " + levelSevenAffixes[levelSevenID]

	if levelTenID == 0 {
		return affixes
	}

	affixes += ", " + levelTenAffixes[levelTenID]

	return affixes
}

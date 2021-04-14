package golib

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

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
// 	- rethink if I actually need this, displaying a list of affixes as string takes a lot of space, should just return an array of ids and display icons
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

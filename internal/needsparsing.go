package internal

import (
	"encoding/json"
	"log"
	"os"
	//"github.com/yassinebenaid/godump"
)

func ParseNeedsJson(needsPath string, logger *log.Logger) NeedsJsonInfo {
	needsJsonFile, err := os.ReadFile(needsPath)
	if err != nil {
		logger.Printf("could not open needsJsonFile. Path: %s Error: %s", needsPath, err.Error())
	}

	var needsJson NeedsJsonInfo
	if err := json.Unmarshal(needsJsonFile, &needsJson); err != nil {
		logger.Printf("could not parse stuff: %s", err.Error())
	}
	// DEBUGGING PRINTS
	//t := needsJson.Versions["0.1"].Needs["feat_req__example__some_title"]
	//var d godump.Dumper
	//logger.Println(d.Sprintln(t))
	//logger.Printf("This is one needsJson parsed: %v\n", needsJson.Versions["0.1"].Needs["feat_req__example__some_title"])
	return needsJson
}

// HACK: Version is assumed 0.1 for now. This might not be ideal
func GetNeedsList(needsJSON NeedsJsonInfo) NeedsInfo {
	return needsJSON.Versions["0.1"].NeedsInfo
}

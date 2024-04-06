package bazica

//
//import (
//	"bazica/internal/coresvc"
//	"encoding/json"
//	"fmt"
//	"log"
//	"os"
//)
//
//type SolarTermSvc struct {
//	svc *coresvc.CoreSvc
//}
//
//func GetSolarTermByYear(year string) *json.Decoder {
//
//	fileToRead := "internal/datasvc/solar-terms-data/" + year + ".json"
//	file, err := os.Open(fileToRead)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Sprintf("Reading file: %s.json", year)
//	defer file.Close()
//
//	// Open the JSON file
//	decoder := json.NewDecoder(file)
//
//	return decoder
//	//err := json.Unmarshal(decoder, &datasvc.SolarTerms{})
//}

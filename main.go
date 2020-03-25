package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"time"

	"github.com/DusanKasan/parsemail"
	"github.com/ghodss/yaml"
	"github.com/kardianos/osext"
)

// config parser base struct
type config struct {
	categories struct {
		category []string
	}
}

// person object for dict
type person struct {
	email      string
	categories []string
	double     bool
}

// helper functions for array
func find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}
func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

// analyse E-Mail file (.eml)
func analyseFile(filePath string, persons map[string]person, finderConfig map[string]string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	email, err := parsemail.Parse(file) // returns Email struct and error
	if err != nil {
		log.Println("could not parse email appropriately:", filePath)
	}

	var regexChecker map[string]*regexp.Regexp
	regexChecker = make(map[string]*regexp.Regexp)

	for fcK, fcV := range finderConfig {
		regexChecker[fcK] = regexp.MustCompile(`(?i)(` + fcV + `)`)
	}

	var mailAddress string
	if len(email.From) < 1 {
		mailAddress = "unknown"
		fmt.Printf("HEADER: %v", email.Header)
	} else {
		mailAddress = email.From[0].Address
	}
	for k, v := range regexChecker {
		match := v.MatchString(email.TextBody)
		if match == true {
			if val, ok := persons[mailAddress]; ok {
				if index := find(val.categories, "unmatched"); index < len(val.categories) {
					val.categories = remove(val.categories, index)
				}
				val.categories = append(val.categories, k)
				if len(val.categories) > 1 {
					val.double = true
				}
			} else {
				persons[mailAddress] = person{
					email:      mailAddress,
					categories: []string{k},
					double:     false,
				}
			}
		} else {
			if _, ok := persons[mailAddress]; !ok {
				persons[mailAddress] = person{
					email:      mailAddress,
					categories: []string{"unmatched"},
					double:     false,
				}
			}
		}
	}
}

func main() {
	startOutput := `
----------------------
- E-Mail Categorizer -
----------------------
`
	fmt.Print(startOutput)

	// get folder of execution or if in development, ignore it
	dir, err := osext.ExecutableFolder()
	if dir[0:4] == "/var" {
		dir = ""
	} else {
		dir += "/"
	}
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	// parse config
	yamlDoc, err := ioutil.ReadFile(dir + "config.yml")
	if err != nil {
		log.Fatal(err)
	}
	jsonDoc, err := yaml.YAMLToJSON(yamlDoc)
	if err != nil {
		fmt.Printf("Error converting YAML to JSON: %s\n", err.Error())
		return
	}
	var someStruct map[string]interface{}
	err = json.Unmarshal(jsonDoc, &someStruct)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %s\n", err.Error())
		return
	}
	var finderConfig map[string]string
	finderConfig = make(map[string]string)
	for elem := range someStruct {
		if reflect.TypeOf(elem).Name() == "string" {
			elemsSlice := someStruct[elem].(map[string]interface{})
			for elem2 := range elemsSlice {
				if reflect.TypeOf(elem2).Name() == "string" {
					elems2Slice := elemsSlice[elem2].([]interface{})
					for i, cat := range elems2Slice {
						if i < len(elems2Slice)-1 {
							finderConfig[elem2] += fmt.Sprintf(`%s\w{0,2}|`, cat)
						} else {
							finderConfig[elem2] += fmt.Sprintf(`%s\w{0,2}`, cat)
						}
					}
				}
			}
		}
	}

	// e-mail folder
	emaildir := dir + "emails"
	emailpath := filepath.FromSlash(dir + "emails")
	files, err := ioutil.ReadDir(emailpath)
	if err != nil {
		log.Fatal(err)
	}

	// analyse
	persons := make(map[string]person)
	for _, file := range files {
		if file.Name()[len(file.Name())-4:len(file.Name())] == ".eml" {
			analyseFile(emaildir+"/"+file.Name(), persons, finderConfig)
		}
	}

	// write CSV
	rows := [][]string{
		{"E-Mail", "Category", "Double-Category (Y|N)}"},
	}

	tn := time.Now()
	fn := fmt.Sprintf("%d-%02d-%02d_%02d_%02d_%02d-results.csv", tn.Year(), tn.Month(), tn.Day(), tn.Hour(), tn.Minute(), tn.Second())
	csvpath := filepath.FromSlash(dir + fn)
	csvfile, err := os.Create(csvpath)
	if err != nil {
		log.Fatal(err)
	}

	csvwriter := csv.NewWriter(csvfile)
	for _, v := range persons {
		for _, e := range v.categories {
			if v.double == true {
				rows = append(rows, []string{v.email, e, "Y"})
			} else {
				rows = append(rows, []string{v.email, e, "N"})
			}
		}
	}
	for _, row := range rows {
		_ = csvwriter.Write(row)
	}

	csvwriter.Flush()
	csvfile.Close()

	// end with some information
	fmt.Printf("> Analysed %d E-Mails\n", len(files))
	fmt.Printf("> Wrote output to %s\n\n", csvpath)
}

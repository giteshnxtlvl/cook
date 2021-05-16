package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

var m = make(map[interface{}]map[string][]string)
var params = make(map[string]string)
var pattern = []string{}
var version = "1.5"

// var verbose = false
var min int
var total = 0

const (
	blue   = "\u001b[38;5;14m"
	green  = "\u001b[38;5;46m"
	purple = "\u001b[38;5;207m"
	red    = "\u001b[38;5;196m"
	bold   = "\u001b[1m"
	white  = "\u001b[38;5;255m"
	reset  = "\u001b[0m"
)

var banner = `

                             
  ░            ░ ░      ░ ░  ░  ░            Created by a person who
  ░ ░        ░ ░ ░ ▒  ░ ░ ░ ▒  ░             got frustated creating 
░░▒ ▒░    ░ ▒ ▒░   ░ ▒ ▒░ ░ ░ ▒ ░            permutation and combination 
░ ░▒ ▒  ░░ ▒░▒░▒░ ░ ▒░▒░▒░ ▒ ▒▒ ▓▒           of words manually.
 ▄████▄   ▒█████   ▒█████   ██ ▄█▀           
▒██▀ ▀█  ▒██▒  ██▒▒██▒  ██▒ ██▄█▒            How the fk you guys were
▒▓█    ▄ ▒██░  ██▒▒██░  ██▒▓███▄░            working without this till yet?
▒▓▓▄ ▄██▒▒██   ██░▒██   ██░▓██ █▄             
 ▒▓███▀ ░░ ████▓▒░░ ████▓▒░▒██▒ █▄ ` + version + `       -Gitesh Sharma @giteshnxtlvl

`

// Goona remove this in future as flag pkg done
func valueInSlice(list []string, val string) bool {
	for _, l := range list {
		if l == val {
			return true
		}
	}
	return false
}

func findRegex(file, expresssion string) []string {
	founded := []string{}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return []string{file + ":" + expresssion}
	}

	r, err := regexp.Compile(expresssion)
	if err != nil {
		panic(err)
	}

	e := make(map[string]bool)
	for _, found := range r.FindAllString(string(content), -1) {
		e[found] = true
	}

	for k := range e {
		founded = append(founded, k)
	}
	return founded
}

func fileValues(file string) []string {
	tmp := []string{}
	readFile, err := os.Open(file)

	if err != nil {
		fmt.Println("Err: Opening File ", file)
	}

	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	for fileScanner.Scan() {
		tmp = append(tmp, fileScanner.Text())
	}

	return tmp
}

func getConfigFile() []byte {

	res, err := http.Get("https://raw.githubusercontent.com/giteshnxtlvl/cook/main/cook.yaml")
	if err != nil {
		log.Fatal(err)
	}

	data, _ := ioutil.ReadAll(res.Body)

	res.Body.Close()

	return data
}

var content []byte
var home, _ = os.UserHomeDir()
var configFile = path.Join(home, ".config", "cook", "cook.yaml")

func applyColumnCases(columnValues []string, columnNum int) {
	temp := []string{}

	// Using cases for columnValues
	if len(columnCases[columnNum]) > 0 {

		//All cases
		if columnCases[columnNum]["A"] {
			for _, t := range final {
				for _, v := range columnValues {
					temp = append(temp, t+strings.ToUpper(v))
					temp = append(temp, t+strings.ToLower(v))
					temp = append(temp, t+strings.Title(v))
				}
			}
		} else {

			if columnCases[columnNum]["U"] {
				for _, t := range final {
					for _, v := range columnValues {
						temp = append(temp, t+strings.ToUpper(v))
					}
				}
			}

			if columnCases[columnNum]["L"] {
				for _, t := range final {
					for _, v := range columnValues {
						temp = append(temp, t+strings.ToLower(v))
					}
				}
			}

			if columnCases[columnNum]["T"] {
				for _, t := range final {
					for _, v := range columnValues {
						temp = append(temp, t+strings.Title(v))
					}
				}
			}
		}

	} else {
		for _, t := range final {
			for _, v := range columnValues {
				temp = append(temp, t+v)
			}
		}
	}

	final = temp
}

//Initializing with empty string, so loops will run for 1st column
var final = []string{""}

func main() {
	fmt.Fprintln(os.Stderr, banner)
	cookConfig()
	parseInput(os.Args[1:])

	for columnNum, param := range pattern {

		columnValues := []string{}

		for _, p := range strings.Split(param, ",") {

			val, success := parseRanges(p)
			if success {
				columnValues = append(columnValues, val...)
				continue
			}
			if val, exists := params[p]; exists {
				columnValues = append(columnValues, parseValue(val)...)
				continue
			}
			if val, exists := m["charSet"][p]; exists {
				chars := strings.Split(val[0], "")
				columnValues = append(columnValues, chars...)
				continue
			}
			if val, exists := m["files"][p]; exists {
				columnValues = append(columnValues, fileValues(val[0])...)
				continue
			}
			if val, exists := m["lists"][p]; exists {
				columnValues = append(columnValues, val...)
				continue
			}
			if val, exists := m["extensions"][p]; exists {
				for _, ext := range val {
					columnValues = append(columnValues, "."+ext)
				}
				continue
			}

			columnValues = append(columnValues, p)
		}

		applyColumnCases(columnValues, columnNum)

		if columnNum >= min {
			for _, v := range final {
				total++
				fmt.Println(v)
			}
		}
	}
	fmt.Fprintln(os.Stderr, "Total Words Generated", total)
}

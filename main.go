// Generate a set of cyclus input files from an xml template and candidate
// values.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var params = [][]interface{}{}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "    %s [xml-template] [json-vals]", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	fname := flag.Arg(0)
	paramfile := flag.Arg(1)

	data, err := ioutil.ReadFile(paramfile)
	if err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(data, &params); err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles(fname)
	if err != nil {
		log.Fatal(err)
	}

	dims := []int{}
	for _, param := range params {
		dims = append(dims, len(param)-1)
	}
	perms := Permute(dims...)

	fbase := filepath.Base(fname)
	fbase = fbase[:len(fbase)-len(filepath.Ext(fbase))]
	fbase = fbase[:len(fbase)-len(filepath.Ext(fbase))]

	legend := map[string]map[string]interface{}{}

	for i, perm := range perms {
		vals := map[string]interface{}{}
		for j, index := range perm {
			vals[params[j][0].(string)] = params[j][index+1]
		}
		name := fmt.Sprintf("%v-%v.xml", fbase, i+1)
		f, err := os.Create(name)
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(f, vals)
		f.Close()

		legend[name] = vals
	}

	data, err = json.MarshalIndent(legend, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(fmt.Sprintf("%v-legend.json", fbase))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		log.Fatal(err)
	}
}

func Permute(dimensions ...int) [][]int {
	return permute(dimensions, []int{})
}

func permute(dimensions []int, prefix []int) [][]int {
	set := make([][]int, 0)

	if len(dimensions) == 1 {
		for i := 0; i < dimensions[0]; i++ {
			val := append(append([]int{}, prefix...), i)
			set = append(set, val)
		}
		return set
	}

	max := dimensions[0]
	for i := 0; i < max; i++ {
		set = append(set, permute(dimensions[1:], append(prefix, i))...)
	}
	return set
}

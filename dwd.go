package main

import "bufio"
import "gopkg.in/alecthomas/kingpin.v2"
import "fmt"
import "encoding/csv"
import "io"
import "os"

var (
	filename = kingpin.Arg("filename", "Filename to load").Required().String()
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func parse_file(filename string) ([]string, []string, []string) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	var column_names []string
	var units []string
	var descriptions []string

	r := csv.NewReader(bufio.NewReader(f))
	r.Comma = ';'
RECORDS:
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		headers := []*[]string{
			&column_names,
			&units,
			&descriptions,
		}

		for _, target := range headers {
			if *target == nil {
				*target = make([]string, len(record))
				copy(*target, record)
				continue RECORDS
			}
		}
	}

	return column_names, units, descriptions
}

func main() {
	kingpin.Parse()
	column_names, units, descriptions := parse_file(*filename)
	fmt.Println(column_names)
	fmt.Println(units)
	fmt.Println(descriptions)
}

package main

import "bufio"
import "fmt"
import "encoding/csv"
import "io"
import "os"

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
	column_names, units, descriptions := parse_file("example_headers")
	fmt.Println(column_names)
	fmt.Println(units)
	fmt.Println(descriptions)
}

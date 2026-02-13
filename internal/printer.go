/*
Copyright © 2026 Amanda Hager Lopes de Andrade Katz amandahla@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package internal

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Printable interface {
	Header() []string
	Row() []interface{}
}

func Print[P Printable](output []P, csv bool) {
	if len(output) == 0 {
		fmt.Println("No data to display")
		return
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	header := dataToRow(output[0].Header())
	t.AppendHeader(header)
	for _, o := range output {
		row := dataToRow(o.Row())
		t.AppendRow(row)
	}
	if csv {
		t.RenderCSV()
	} else {
		t.Render()
	}
}

func dataToRow[T any](d []T) table.Row {
	row := make(table.Row, len(d))
	for i, s := range d {
		row[i] = s
	}
	return row
}

package common

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func Output(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorder(true)
	table.AppendBulk(data)
	table.Render()
}
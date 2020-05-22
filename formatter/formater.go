package formatter

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func Output(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorder(true)
	table.AppendBulk(data)
	table.Render()
}
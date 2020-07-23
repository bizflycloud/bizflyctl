package formatter

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/olekukonko/tablewriter"
)

// Output is func support string data
func Output(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(true)
	table.AppendBulk(data)
	table.Render()
}

// SimpleOutput is func support any type of data
func SimpleOutput(header table.Row, rows []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(header)
	t.AppendRows(rows)
	t.Render()
}

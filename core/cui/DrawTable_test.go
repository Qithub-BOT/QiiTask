package cui_test

import (
	"github.com/Qithub-BOT/QiiTask/core/cui"
	"github.com/jedib0t/go-pretty/v6/table"
)

func ExampleUI_DrawTable() {
	tableTmp := table.NewWriter()

	tableTmp.AppendHeader(table.Row{"#", "title"})
	tableTmp.AppendRow([]interface{}{1, "uno"})
	tableTmp.AppendRow([]interface{}{2, "dos"})
	tableTmp.AppendSeparator()
	tableTmp.AppendRow([]interface{}{3, "tres"})

	// Available table styles:
	//     cui.AsDefaultTable (= cui.AsSimpleTable)
	//     cui.AsSimpleTable
	//     cui.AsColoredTable
	//     cui.AsMarkdownTable
	//     cui.AsCSVTable
	//     cui.AsHTMLTable
	ui := cui.New()
	ui.DrawTable(tableTmp, cui.AsDefaultTable)

	// Output:
	// +---+-------+
	// | # | TITLE |
	// +---+-------+
	// | 1 | uno   |
	// | 2 | dos   |
	// +---+-------+
	// | 3 | tres  |
	// +---+-------+
}

func ExampleUI_DrawTable_colored() {
	tableTmp := table.NewWriter()

	tableTmp.AppendHeader(table.Row{"#", "title"})
	tableTmp.AppendRow([]interface{}{1, "uno"})
	tableTmp.AppendRow([]interface{}{2, "dos"})
	tableTmp.AppendSeparator()
	tableTmp.AppendRow([]interface{}{3, "tres"})

	ui := cui.New()
	ui.DrawTable(tableTmp, cui.AsColoredTable)

	// Output:
	// [44;37m # [0m[44;37m TITLE [0m
	// [40;97m 1 [0m[40;97m uno   [0m
	// [100;37m 2 [0m[100;37m dos   [0m
	// [100;37m---[0m[100;37m-------[0m
	// [40;97m 3 [0m[40;97m tres  [0m
}

func ExampleUI_DrawTable_markdown() {
	tableTmp := table.NewWriter()

	tableTmp.AppendHeader(table.Row{"#", "title"})
	tableTmp.AppendRow([]interface{}{1, "uno"})
	tableTmp.AppendRow([]interface{}{2, "dos"})
	tableTmp.AppendSeparator()
	tableTmp.AppendRow([]interface{}{3, "tres"})

	ui := cui.New()
	ui.DrawTable(tableTmp, cui.AsMarkdownTable)

	// Output:
	// | # | title |
	// | ---:| --- |
	// | 1 | uno |
	// | 2 | dos |
	// | 3 | tres |
}

func ExampleUI_DrawTable_csv() {
	tableTmp := table.NewWriter()

	tableTmp.AppendHeader(table.Row{"#", "title", "description"})
	tableTmp.AppendRow([]interface{}{1, "uno", "uno, one, いち are the same"})
	tableTmp.AppendRow([]interface{}{2, "dos", "dos, two, に are the same"})
	tableTmp.AppendSeparator()
	tableTmp.AppendRow([]interface{}{3, "tres", "tres and tree are not same"})

	ui := cui.New()
	ui.DrawTable(tableTmp, cui.AsCSVTable)

	// Output:
	// #,title,description
	// 1,uno,"uno\, one\, いち are the same"
	// 2,dos,"dos\, two\, に are the same"
	// 3,tres,tres and tree are not same
}

func ExampleUI_DrawTable_html() {
	tableTmp := table.NewWriter()

	tableTmp.AppendHeader(table.Row{"#", "title"})
	tableTmp.AppendRow([]interface{}{1, "uno"})
	tableTmp.AppendRow([]interface{}{2, "dos"})
	tableTmp.AppendSeparator()
	tableTmp.AppendRow([]interface{}{3, "tres"})

	ui := cui.New()
	ui.DrawTable(tableTmp, cui.AsHTMLTable)

	// Output:
	// <table class="go-pretty-table">
	//   <thead>
	//   <tr>
	//     <th align="right">#</th>
	//     <th>title</th>
	//   </tr>
	//   </thead>
	//   <tbody>
	//   <tr>
	//     <td align="right">1</td>
	//     <td>uno</td>
	//   </tr>
	//   <tr>
	//     <td align="right">2</td>
	//     <td>dos</td>
	//   </tr>
	//   <tr>
	//     <td align="right">3</td>
	//     <td>tres</td>
	//   </tr>
	//   </tbody>
	// </table>
}

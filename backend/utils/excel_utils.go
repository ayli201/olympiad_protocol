package utils

import "github.com/xuri/excelize/v2"

func CreateStyle(
	f *excelize.File,
	borders PBorders,
	halign PAlign,
	valign PAlign,
	bold bool,
	italic bool,
	fontSize float64,
	isPercent bool,
) (int, error) {

	style := &excelize.Style{}

	// borders
	switch borders {
	case BorderAll:
		style.Border = []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		}
	case BorderBottom:
		style.Border = []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 1},
		}
	}

	// alignment
	style.Alignment = &excelize.Alignment{
		WrapText: true,
	}

	switch halign {
	case AlignHorCenter:
		style.Alignment.Horizontal = "center"
	case AlignHorLeft:
		style.Alignment.Horizontal = "left"
	case AlignHorRight:
		style.Alignment.Horizontal = "right"
	}

	if valign == AlignVertCenter {
		style.Alignment.Vertical = "center"
	}

	// font
	style.Font = &excelize.Font{
		Family: "Times New Roman",
		Size:   fontSize,
		Bold:   bold,
		Italic: italic,
	}

	if isPercent {
		style.NumFmt = 10 // 0.00%
	}

	return f.NewStyle(style)
}

func ApplyStyle(
	f *excelize.File,
	sheet string,
	styleID int,
	col1, row1, col2, row2 int,
) error {

	from, _ := excelize.CoordinatesToCellName(col1, row1)
	to, _ := excelize.CoordinatesToCellName(col2, row2)

	return f.SetCellStyle(sheet, from, to, styleID)
}

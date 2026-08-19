package services

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type StyleOptions struct {
	FontName     string            // По умолчанию: "Times New Roman"
	FontSize     float64           // По умолчанию: 11
	Bold         bool              // Жирный
	Italic       bool              // Курсив
	FontColor    string            // Цвет текста в HEX (напр., "FFFFFF" или "#FF0000")
	BgColor      string            // Цвет фона в HEX (напр., "#1F4E78")
	HAlign       string            // "left", "center", "right"
	VAlign       string            // "top", "center", "bottom"
	WrapText     *bool             // Перенос текста
	CustomNumFmt *string           // Числовой формат (напр., "#,##0.00 ₴")
	Border       []excelize.Border // Границы ячейки (если нужны)
	Underline    bool              // Подчеркивание
}

type ExcelStreamBuilder struct {
	file       *excelize.File
	sheet      string
	writer     *excelize.StreamWriter
	currentRow int64
}

func NewExcelStreamBuilder(f *excelize.File, sheet string) (*ExcelStreamBuilder, error) {
	sw, err := f.NewStreamWriter(sheet)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания StreamWriter: %w", err)
	}

	return &ExcelStreamBuilder{
		file:       f,
		sheet:      sheet,
		writer:     sw,
		currentRow: 1,
	}, nil
}

// func (sb *ExcelStreamBuilder) AddStyle(style *excelize.Style) (int, error) {
// 	return sb.file.NewStyle(style)
// }

// WriteRowWithOffset записывает строку, начиная с указанной колонки (1 = A, 2 = B и т.д.)
func (sb *ExcelStreamBuilder) WriteRowWithOffset(startCol int, values []interface{}, height float64, styleIDs ...int) error {
	rowCells := make([]interface{}, len(values))

	for i, val := range values {
		cell := excelize.Cell{Value: val}

		if i < len(styleIDs) && styleIDs[i] > 0 {
			cell.StyleID = styleIDs[i]
		}

		rowCells[i] = cell
	}

	cellRef, err := excelize.CoordinatesToCellName(startCol, int(sb.currentRow))
	if err != nil {
		return err
	}

	// Готовим опции строки
	var rowOpts []excelize.RowOpts
	if height > 0 {
		rowOpts = append(rowOpts, excelize.RowOpts{Height: height})
	} else {
		rowOpts = append(rowOpts, excelize.RowOpts{Height: 30})
	}

	if err := sb.writer.SetRow(cellRef, rowCells, rowOpts...); err != nil {
		return fmt.Errorf("ошибка записи строки %d: %w", sb.currentRow, err)
	}

	sb.currentRow++
	return nil
}

func (sb *ExcelStreamBuilder) MergeCell(hCell, vCell string) error {
	return sb.writer.MergeCell(hCell, vCell)
}

func (sb *ExcelStreamBuilder) AddStyle(opts StyleOptions) (int, error) {
	// 1. Устанавливаем дефолты для шрифта
	fontName := "Times New Roman"
	if opts.FontName != "" {
		fontName = opts.FontName
	}

	fontSize := 11.0
	if opts.FontSize > 0 {
		fontSize = opts.FontSize
	}

	underline := ""
	if opts.Underline {
		underline = "single"
	}

	font := &excelize.Font{
		Family:    fontName,
		Size:      fontSize,
		Bold:      opts.Bold,
		Italic:    opts.Italic,
		Underline: underline,
	}

	if opts.FontColor != "" {
		font.Color = opts.FontColor
	}

	// 2. Формируем итоговую структуру excelize.Style
	halign := "center"
	if opts.HAlign != "" {
		halign = opts.HAlign
	}
	valign := "center"
	if opts.VAlign != "" {
		valign = opts.VAlign
	}

	wrapText := true
	if opts.WrapText != nil {
		wrapText = *opts.WrapText
	}

	style := &excelize.Style{
		Font: font,
		Alignment: &excelize.Alignment{
			Horizontal: halign,
			Vertical:   valign,
			WrapText:   wrapText,
		},
		Border: opts.Border,
	}

	// 3. Задаем заливку фона (если указана)
	if opts.BgColor != "" {
		style.Fill = excelize.Fill{
			Type:    "pattern",
			Color:   []string{opts.BgColor},
			Pattern: 1,
		}
	}

	// 4. Задаем числовой формат (если указан)
	if opts.CustomNumFmt != nil {
		style.CustomNumFmt = opts.CustomNumFmt
	}

	// Регистрируем стиль в документе и возвращаем его ID
	return sb.file.NewStyle(style)
}

func (sb *ExcelStreamBuilder) SetDefaultRowHeight(height float64) error {
	return sb.file.SetSheetProps(sb.sheet, &excelize.SheetPropsOptions{
		DefaultRowHeight: &height,
		CustomHeight:     boolPtr(true), // Сообщаем Excel, что используется пользовательская высота
	})
}

func boolPtr(b bool) *bool { return &b }

// SetColWidth устанавливает ширину диапазона колонок (1 = A, 2 = B и т.д.)
func (sb *ExcelStreamBuilder) SetColWidth(startCol, endCol int, width float64) error {
	return sb.writer.SetColWidth(startCol, endCol, width)
}

func (sb *ExcelStreamBuilder) BorderFull() []excelize.Border {
	return []excelize.Border{
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "left", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}
}

func (sb *ExcelStreamBuilder) Flush() error {
	return sb.writer.Flush()
}

func strPtr(s string) *string {
	return &s
}

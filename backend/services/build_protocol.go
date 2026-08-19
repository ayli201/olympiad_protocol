package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func (s *ExportService) buildReport(ctx context.Context, path string) error {
	setting, err := s.settingsService.GetByName(ctx, "tasks_count")
	if err != nil {
		return err
	}
	tasksCount, err := strconv.Atoi(setting.Value)
	if err != nil {
		log.Fatalf("Ошибка преобразования: %v", err)
	}

	// Создаем файл Excel
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Протокол"

	f.SetSheetName("Sheet1", sheetName)

	s.buildPrintSettings(f, sheetName)

	// Создаем поток для записи
	builder, err := NewExcelStreamBuilder(f, sheetName)
	if err != nil {
		log.Fatal(err)
	}
	s.builder = builder

	// s.builder.SetDefaultRowHeight(30)

	colWidths := make([]float64, 9+tasksCount)
	colWidths[0] = 6.2   // Колонка №
	colWidths[1] = 29.6  // Колонка ФИО
	colWidths[2] = 12.2  // Колонка Шифр
	colWidths[3] = 26.1  // Колонка ОО
	colWidths[4] = 6.9   // Колонка Класс
	pointColWidth := 6.0 // Колонки баллов
	if tasksCount*6 < 30 {
		pointColWidth = 30 / float64(tasksCount)
	}
	for i := range tasksCount {
		colWidths[i+5] = pointColWidth
	}
	nextCell := 5 + tasksCount
	colWidths[nextCell] = 16   // Колонка Итог
	colWidths[nextCell+1] = 16 // Колонка % выполнения
	colWidths[nextCell+2] = 16 // Колонка Рейтинг
	colWidths[nextCell+3] = 16 // Колонка % выполнения

	s.setColumnWidths(colWidths)
	s.buildTitleRow(ctx, tasksCount)
	s.buildMaxScoreRow(ctx)
	s.buildHeaderRows(tasksCount)
	s.buildData(ctx, tasksCount, colWidths)
	s.buildFooterRows()

	// Завершаем поток
	if err := builder.Flush(); err != nil {
		log.Fatalf("Ошибка Flush: %v", err)
	}

	// Сохраняем файл
	if err := f.SaveAs(path); err != nil {
		log.Fatalf("Ошибка сохранения: %v", err)
	}

	fmt.Println("Файл успешно сформирован: protocol.xlsx")
	return nil
}

func (s *ExportService) setColumnWidths(colWidths []float64) {
	// Настраиваем ширину конкретных колонок:
	for i, width := range colWidths {
		s.builder.SetColWidth(i+1, i+1, width)
	}
}

func (s *ExportService) buildTitleRow(ctx context.Context, tasksCount int) error {
	// Стиль основного заголовка (B1:S1)
	titleStyle, _ := s.builder.AddStyle(StyleOptions{
		FontSize: 14,
		Bold:     true,
	})
	setting, err := s.settingsService.GetByName(ctx, "year_start")
	if err != nil {
		return err
	}
	yearStart := setting.Value
	setting, err = s.settingsService.GetByName(ctx, "year_end")
	if err != nil {
		return err
	}
	yearEnd := setting.Value
	setting, err = s.settingsService.GetByName(ctx, "discipline")
	if err != nil {
		return err
	}
	discipline := setting.Value

	titleValues := make([]interface{}, 9+tasksCount)
	parts := make([]string, 6)
	parts[0] = "Протокол муниципального этапа Всероссийской Олимпиады школьников в"
	parts[1] = yearStart
	parts[2] = "-"
	parts[3] = yearEnd
	parts[4] = "по"
	parts[5] = discipline
	titleValues[0] = strings.Join(parts, " ")

	titleStyles := make([]int, 9+tasksCount)
	for i := range titleStyles {
		titleStyles[i] = titleStyle
	}

	// Записываем строку 1 с отступом в 2-ю колонку (B)
	if err := s.builder.WriteRowWithOffset(1, titleValues, 0, titleStyles...); err != nil {
		log.Fatal(err)
	}
	titleLastCell, _ := excelize.CoordinatesToCellName(9+tasksCount, 1)
	s.builder.MergeCell("A1", titleLastCell)

	return nil
}

func (s *ExportService) buildMaxScoreRow(ctx context.Context) error {
	setting, err := s.settingsService.GetByName(ctx, "max_points")
	if err != nil {
		return err
	}
	maxPoints := setting.Value

	// Стиль строки максимальных баллов
	maxScoreStyles := make([]int, 3)
	maxScoreStyles[0], _ = s.builder.AddStyle(StyleOptions{
		Italic: true,
		Bold:   true,
		HAlign: "left",
		VAlign: "bottom",
	})
	maxScoreStyles[2], _ = s.builder.AddStyle(StyleOptions{
		Bold:     true,
		FontSize: 14,
		Border:   s.builder.BorderFull(),
		VAlign:   "bottom",
	})

	// Максимальное количество баллов
	maxPointsRow := []interface{}{
		"Максимальное количество баллов", "",
		maxPoints,
	}
	_ = s.builder.MergeCell("B2", "C2")
	if err := s.builder.WriteRowWithOffset(2, maxPointsRow, 17.25, maxScoreStyles...); err != nil {
		log.Fatal(err)
		return err
	}

	if err := s.builder.WriteRowWithOffset(1, []interface{}{}, 5); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *ExportService) buildHeaderRows(tasksCount int) error {
	// Стиль первого уровня шапки (родительские объединения)
	parentHeaderStyle, _ := s.builder.AddStyle(StyleOptions{
		Italic: true,
		Bold:   true,
		Border: s.builder.BorderFull(),
	})

	// Стиль второго уровня шапки (подколонки)
	childHeaderStyle, _ := s.builder.AddStyle(StyleOptions{
		Italic: true,
		Bold:   true,
		Border: s.builder.BorderFull(),
	})

	// первая строка шапки
	headerL1 := make([]interface{}, 9+tasksCount)
	headerL1[0] = "№"
	headerL1[1] = "Фамилия, имя, отчество"
	headerL1[2] = "Шифр"
	headerL1[3] = "ОО"
	headerL1[4] = "Класс"
	headerL1[5] = "Количество баллов за задание"
	for i := range tasksCount {
		headerL1[6+i] = ""
	}
	headerL1[5+tasksCount] = "Итог"
	headerL1[6+tasksCount] = "% выполнения"
	headerL1[7+tasksCount] = "Рейтинг"
	headerL1[8+tasksCount] = "Статус"

	headerL1Styles := make([]int, 9+tasksCount)
	for i := range headerL1Styles {
		headerL1Styles[i] = parentHeaderStyle
	}

	if err := s.builder.WriteRowWithOffset(1, headerL1, 21, headerL1Styles...); err != nil {
		log.Fatal(err)
		return err
	}

	// вторая строка шапки
	headerL2 := make([]interface{}, 9+tasksCount)
	for i := 0; i < 5; i++ {
		headerL2[i] = ""
	}
	for i := 5; i < 5+tasksCount; i++ {
		headerL2[i] = i - 4
	}
	for i := 5 + tasksCount; i < 9+tasksCount; i++ {
		headerL2[i] = ""
	}

	headerL2Styles := make([]int, 9+tasksCount)
	for i := range headerL2Styles {
		headerL2Styles[i] = childHeaderStyle
	}

	if err := s.builder.WriteRowWithOffset(1, headerL2, 25, headerL2Styles...); err != nil {
		log.Fatal(err)
		return err
	}

	// Объединение ячеек
	s.builder.MergeCell("A4", "A5") // ID
	s.builder.MergeCell("B4", "B5") // ФИО
	s.builder.MergeCell("C4", "C5") // Шифр
	s.builder.MergeCell("D4", "D5") // ОО
	s.builder.MergeCell("E4", "E5") // Класс
	firstCell, _ := excelize.CoordinatesToCellName(6+tasksCount, 4)
	lastCell, _ := excelize.CoordinatesToCellName(6+tasksCount, 5)
	s.builder.MergeCell(firstCell, lastCell) // Итог
	firstCell, _ = excelize.CoordinatesToCellName(7+tasksCount, 4)
	lastCell, _ = excelize.CoordinatesToCellName(7+tasksCount, 5)
	s.builder.MergeCell(firstCell, lastCell) // % выполнения
	firstCell, _ = excelize.CoordinatesToCellName(8+tasksCount, 4)
	lastCell, _ = excelize.CoordinatesToCellName(8+tasksCount, 5)
	s.builder.MergeCell(firstCell, lastCell) // Рейтинг
	firstCell, _ = excelize.CoordinatesToCellName(9+tasksCount, 4)
	lastCell, _ = excelize.CoordinatesToCellName(9+tasksCount, 5)
	s.builder.MergeCell(firstCell, lastCell) // Статус

	// Б. Горизонтальное объединение динамических столбцов во 2-й строке:

	firstCell, _ = excelize.CoordinatesToCellName(6, 4)
	lastCell, _ = excelize.CoordinatesToCellName(6+tasksCount-1, 4)
	s.builder.MergeCell(firstCell, lastCell)
	return nil
}

func (s *ExportService) buildData(ctx context.Context, tasksCount int, colWidths []float64) error {
	participants, err := s.participantService.GetAll(ctx)
	if err != nil {
		return err
	}

	centerStyle, _ := s.builder.AddStyle(StyleOptions{
		FontSize: 12,
		Border:   s.builder.BorderFull(),
	})

	// Настройка стилей для колонок данных
	rowStyles := make([]int, 9+tasksCount)
	rowStyles[0] = centerStyle
	rowStyles[1] = centerStyle
	rowStyles[2] = centerStyle
	rowStyles[3] = centerStyle
	rowStyles[4] = centerStyle
	for i := 5; i < 5+tasksCount; i++ {
		rowStyles[i] = centerStyle
	}
	rowStyles[5+tasksCount] = centerStyle
	rowStyles[6+tasksCount] = centerStyle
	rowStyles[7+tasksCount] = centerStyle
	rowStyles[8+tasksCount] = centerStyle

	for idx, participant := range participants {
		row := make([]interface{}, 9+tasksCount)
		row[0] = idx + 1
		row[1] = participant.FullName
		row[2] = participant.Cipher
		res, _ := s.schoolService.Get(ctx, int64(participant.SchoolId))
		if res != nil {
			row[3] = res.Label
		} else {
			row[3] = ""
		}
		row[4] = participant.ClassName

		for i, score := range participant.Tasks {
			if i > tasksCount-1 {
				break
			}
			row[5+i] = int(score.Value)
		}

		row[5+tasksCount] = *participant.Total
		row[6+tasksCount] = *participant.Percent
		row[7+tasksCount] = *participant.Rating
		row[8+tasksCount] = *participant.Status

		rowHeight := calculateMaxRowHeight(row, ConvertToColumnConfigs(colWidths, 12))

		if err := s.builder.WriteRowWithOffset(1, row, rowHeight, rowStyles...); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func (s *ExportService) buildFooterRows() {
	if err := s.builder.WriteRowWithOffset(1, []interface{}{}, 20.3); err != nil {
		log.Fatal(err)
	}

	signatureStyle, _ := s.builder.AddStyle(StyleOptions{
		HAlign:   "right",
		VAlign:   "bottom",
		FontSize: 14,
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	dateStyle, _ := s.builder.AddStyle(StyleOptions{
		HAlign:   "center",
		VAlign:   "bottom",
		FontSize: 14,
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	styles := make([]int, 3)
	styles[0] = signatureStyle
	styles[1] = 0
	styles[2] = dateStyle

	now := time.Now()

	row := []interface{}{
		"Дата проведения:",
		"",
		now.Format("02.01.2006"),
	}
	if err := s.builder.WriteRowWithOffset(2, row, 20.3, styles...); err != nil {
		log.Fatal(err)
	}

	// firstCell, _ := excelize.CoordinatesToCellName(4, currentRowId)
	// lastCell, _ := excelize.CoordinatesToCellName(9, currentRowId)
	// s.builder.MergeCell(firstCell, lastCell)

	row = []interface{}{
		"Председатель жюри:",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
	}

	fioStyle, _ := s.builder.AddStyle(StyleOptions{
		HAlign:   "left",
		VAlign:   "bottom",
		FontSize: 14,
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	styles = make([]int, 8)
	styles[0] = signatureStyle
	styles[1] = 0
	styles[2] = fioStyle
	styles[3] = fioStyle
	styles[4] = fioStyle
	styles[5] = fioStyle
	styles[6] = fioStyle
	styles[7] = fioStyle
	if err := s.builder.WriteRowWithOffset(2, row, 20.3, styles...); err != nil {
		log.Fatal(err)
	}

	firstCell, _ := excelize.CoordinatesToCellName(4, int(s.builder.currentRow-1))
	lastCell, _ := excelize.CoordinatesToCellName(9, int(s.builder.currentRow-1))
	s.builder.MergeCell(firstCell, lastCell)

	row = []interface{}{
		"Члены жюри:",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
	}
	if err := s.builder.WriteRowWithOffset(2, row, 20.3, styles...); err != nil {
		log.Fatal(err)
	}

	firstCell, _ = excelize.CoordinatesToCellName(4, int(s.builder.currentRow-1))
	lastCell, _ = excelize.CoordinatesToCellName(9, int(s.builder.currentRow-1))
	s.builder.MergeCell(firstCell, lastCell)
}

func (s *ExportService) buildPrintSettings(f *excelize.File, sheet string) {
	enable := true
	if err := f.SetSheetProps(sheet, &excelize.SheetPropsOptions{
		FitToPage: &enable,
	}); err != nil {
		fmt.Println(err)
	}

	var (
		size        = 9
		orientation = "landscape"
		fitToWidth  = 1
		fitToHeight = 0
	)

	// 1. Ориентация: Альбомная (Landscape), Масштаб: Fit to page (1x1), Размер: A4 (PaperSize 9)
	if err := f.SetPageLayout(sheet, &excelize.PageLayoutOptions{
		Size:        &size,
		Orientation: &orientation,
		FitToWidth:  &fitToWidth,
		FitToHeight: &fitToHeight,
	}); err != nil {
		fmt.Println(err)
	}

	top := 0.35
	bottom := 0.35
	left := 0.25
	right := 0.25
	horizontally := true
	vertically := false
	header := 0.0
	footer := 0.0

	if err := f.SetPageMargins(sheet, &excelize.PageLayoutMarginsOptions{
		Left:         &left,
		Right:        &right,
		Top:          &top,
		Bottom:       &bottom,
		Horizontally: &horizontally,
		Vertically:   &vertically,
		Header:       &header,
		Footer:       &footer,
	}); err != nil {
		fmt.Println(err)
	}

	if err := f.SetDefinedName(&excelize.DefinedName{
		Name:     "_xlnm.Print_Titles",
		RefersTo: fmt.Sprintf("'%s'!$4:$5", sheet),
		Scope:    sheet,
	}); err != nil {
		fmt.Println(err)
	}
}

func calculateCellHeight(text string, colWidth int, fontSize float64) float64 {
	if text == "" {
		return fontSize + 6.0
	}

	totalLines := 0

	// Разбиваем текст по явным переносам строки '\n'
	paragraphs := strings.Split(text, "\n")
	for _, p := range paragraphs {
		if len(p) == 0 {
			totalLines++
			continue
		}
		// Считаем, сколько строк займет параграф с учетом ширины колонки
		linesInParagraph := int(math.Ceil(float64(len([]rune(p))) / float64(colWidth)))
		if linesInParagraph == 0 {
			linesInParagraph = 1
		}
		totalLines += linesInParagraph
	}

	// Высота одной строки текста = размер шрифта + межстрочный интервал (4pt)
	singleLineHeight := fontSize + 4.0

	// Общая высота = (кол-во строк * высота строки) + отступы сверху/снизу (6pt)
	return (float64(totalLines) * singleLineHeight) + 6.0
}

type ColumnConfig struct {
	Width    int     // Ширина колонки в символах (как в Excel)
	FontSize float64 // Размер шрифта в pt (например, 10.0, 12.0)
}

func calculateMaxRowHeight(rowCells []interface{}, colConfigs []ColumnConfig) float64 {
	// Базовая минимальная высота (стандарт для 10-11pt шрифта)
	maxHeight := 30.0

	for i, text := range rowCells {
		// Если для колонки нет настроек, используем дефолтные
		config := ColumnConfig{Width: 15, FontSize: 12.0}
		if i < len(colConfigs) && colConfigs[i].Width > 0 {
			config = colConfigs[i]
		}

		// Считаем высоту текущей ячейки
		cellHeight := calculateCellHeight(fmt.Sprint(text), config.Width, config.FontSize)

		// Находим максимальную высоту среди всех ячеек строки
		if cellHeight > maxHeight {
			maxHeight = cellHeight + 12
		}
	}

	// Округляем до ближайшего целого/полуцелого (Правило Кратности для печати)
	return math.Ceil(maxHeight*2) / 2
}

func ConvertToColumnConfigs(colWidths []float64, defaultFontSize float64) []ColumnConfig {
	configs := make([]ColumnConfig, len(colWidths))

	for i, w := range colWidths {
		// Округляем в меньшую сторону (Floor), чтобы случайно не переоценить ширину ячейки.
		// Если ширина будет чуть меньше реальной, алгоритм с запасом рассчитает высоту строки и текст не обрежется.
		widthInt := int(math.Floor(w))
		if widthInt < 1 {
			widthInt = 1 // Защита от нулевой ширины
		}

		configs[i] = ColumnConfig{
			Width:    widthInt,
			FontSize: defaultFontSize,
		}
	}

	return configs
}

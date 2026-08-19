package services

import (
	"context"
	"fmt"
	"log"
	"protocol/backend/models"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/xuri/excelize/v2"
)

type ImportService struct {
	db            *sqlx.DB
	schoolService *SchoolService
}

func NewImportService(db *sqlx.DB) *ImportService {
	schoolService := NewSchoolService(db)
	return &ImportService{db: db, schoolService: schoolService}
}

func (s *ImportService) openFile(ctx context.Context, path string) ([]models.Participant, error) {
	file, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 1. Получаем имя самого первого листа (индекс 0)
	firstSheet := file.GetSheetName(0)
	if firstSheet == "" {
		log.Fatal("В файле нет листов")
		return nil, fmt.Errorf("В файле нет листов")
	}

	rows, err := file.Rows(firstSheet)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var (
		currentRowNum = 0
		headerRowNum  = -1
		headerL1      []string // Буфер для 4-й строки Excel
		headerL2      []string // Буфер для 5-й строки Excel
		normalizedL1  []string // Нормализованная L1 (протянутые смерженные заголовки)
		participants  []models.Participant
	)

	for rows.Next() {
		currentRowNum++

		// Пропускаем мета-заголовок (строка 1) и двухстрочную шапку таблицы (строки 2 и 3)
		// if currentRowNum < 4 {
		// 	continue
		// }

		rowCells, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		if len(rowCells) == 0 {
			continue
		}

		if headerRowNum == -1 && rowCells[0] != "№" && currentRowNum < 5 {
			continue
		}

		if rowCells[0] == "№" {
			headerRowNum = currentRowNum
		}

		if headerRowNum != -1 && currentRowNum == headerRowNum {
			headerL1 = make([]string, len(rowCells))
			copy(headerL1, rowCells)
			continue
		}

		if headerRowNum != -1 && currentRowNum == headerRowNum+1 {
			headerL2 = make([]string, len(rowCells))
			copy(headerL2, rowCells)

			// В момент чтения 5-й строки мы Готовы нормализовать L1!
			normalizedL1 = s.buildNormalizedL1(headerL1, headerL2)
			continue
		}

		// -------------------------------------------------------------
		// ОБРАБОТКА ДАННЫХ
		// -------------------------------------------------------------
		if len(rowCells) == 0 {
			continue // Пропускаем пустые строки
		}

		if rowCells[1] == "Дата проведения:" {
			break
		}

		participant := s.parseStudentRow(ctx, rowCells, normalizedL1, headerL2)
		if participant.FullName != "" {
			participants = append(participants, participant)
		}
	}

	if len(participants) == 0 {
		return nil, fmt.Errorf("Не найдено ни одного участника")
	}

	return participants, nil
}

func normalizeString(s string) string {
	// 1. Приводим к нижнему регистру (русский + английский + любые Unicode символы)
	s = strings.ToLower(s)

	// 2. Заменяем "ё" на "е" (частая проблема при ручном вводе)
	// s = strings.ReplaceAll(s, "ё", "е")

	// 3. Заменяем неразрывные пробелы (\u00a0) на обычные
	s = strings.ReplaceAll(s, "\u00a0", " ")

	// 4. Удаляем пробелы по краям и схлопываем множественные пробелы внутри
	// FieldsSplit разбивает по ЛЮБЫМ пробельным символам (пробел, табуляция, \n)
	words := strings.Fields(s)

	return strings.Join(words, " ")
}

func (s *ImportService) buildNormalizedL1(headerL1, headerL2 []string) []string {
	maxCols := len(headerL1)
	if len(headerL2) > maxCols {
		maxCols = len(headerL2)
	}

	normalized := make([]string, maxCols)
	currentGroupHeader := ""

	for colIdx := 0; colIdx < maxCols; colIdx++ {
		val := ""
		if colIdx < len(headerL1) {
			val = strings.TrimSpace(headerL1[colIdx])
		}

		// Если встречаем новое название группы — запоминаем его.
		// Если ячейка пустая (из-за MergeCell) — используем текущую группу!
		if val != "" {
			currentGroupHeader = val
		}
		normalized[colIdx] = normalizeString(currentGroupHeader)
	}

	return normalized
}

func (s *ImportService) parseStudentRow(ctx context.Context, rowCells []string, normalizedL1, headerL2 []string) models.Participant {
	participant := models.Participant{}
	var tasks []models.Task

	maxCols := len(normalizedL1)

	for colIdx := 0; colIdx < maxCols; colIdx++ {
		var cellValue string
		if colIdx < len(rowCells) {
			cellValue = strings.TrimSpace(rowCells[colIdx])
		}

		groupTitle := normalizedL1[colIdx]

		subTitle := ""
		if colIdx < len(headerL2) {
			subTitle = strings.TrimSpace(headerL2[colIdx])
		}

		participant.CreatedAt = time.Now()

		// Определяем поле по названию из нормализованной шапки L1 или подзаголовка L2
		switch {
		case groupTitle == "фамилия, имя, отчество":
			participant.FullName = cellValue
			break

		case groupTitle == "шифр":
			participant.Cipher = cellValue
			break

		case groupTitle == "оо":
			res, _ := s.schoolService.FindByTitle(ctx, cellValue)
			if res != nil {
				participant.SchoolId = res.Value
			} else {
				participant.SchoolId = 0
			}
			break

		case groupTitle == "класс":
			participant.ClassName = cellValue
			break

		// Группа столбцов для связанной сущности (Задания/Оценки)
		case groupTitle == "количество баллов за задание":
			if subTitle != "" {
				subValue, _ := strconv.Atoi(subTitle)
				score, _ := parseFloat(cellValue)
				tasks = append(tasks, models.Task{
					Number: subValue,
					Value:  float32(score),
				})
			}
			break
		}
	}

	participant.Tasks = tasks
	return participant
}

// parseFloat обработает случай, если в Excel число записано с запятой (4,5)
func parseFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	s = strings.ReplaceAll(s, ",", ".")
	return strconv.ParseFloat(s, 64)
}

func (s *ImportService) ImportData(ctx context.Context) ([]models.Participant, error) {
	path, err := application.Get().Dialog.OpenFile().
		SetTitle("Открыть файл").
		AddFilter("Excel Files", "*.xlsx").
		AddFilter("All Files", "*.*").
		PromptForSingleSelection()

	if err != nil {
		return nil, err
	}
	if path == "" {
		return nil, fmt.Errorf("Загрузка отменена")
	}

	participants, err := s.ImportDataFromFile(ctx, path)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (s *ImportService) ImportDataFromFile(ctx context.Context, path string) ([]models.Participant, error) {
	participants, err := s.openFile(ctx, path)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

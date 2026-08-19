import type { HeaderRow } from "./components/DataTable/types";

export let headerRows: (tasksCount: number) => HeaderRow[] = (
  tasksCount: number,
) => [
  {
    cells: [
      { value: "№", rowspan: 2, columnName: "index" },
      { value: "Ф.И.О.", rowspan: 2, columnName: "fullName" },
      { value: "Шифр", rowspan: 2, columnName: "cipher" },
      { value: "ОО", rowspan: 2, columnName: "schoolId" },
      { value: "Класс", rowspan: 2, columnName: "className" },
      {
        value: "Количество баллов за задание",
        colspan: tasksCount,
        isGroup: true,
        columnName: "tasks",
      },
      { value: "Итог", rowspan: 2, columnName: "total" },
      { value: "% выполнения", rowspan: 2, columnName: "percent" },
      { value: "Рейтинг", rowspan: 2, columnName: "rating" },
      { value: "Статус", rowspan: 2, columnName: "status" },
    ],
  },
  {
    cells: [...new Array(tasksCount)].map((_, idx) => ({
      value: idx + 1,
      parentName: "tasks",
      columnName: `task_${idx}`,
    })),
  },
];

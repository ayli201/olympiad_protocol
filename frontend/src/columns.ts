import { type School } from "./types/school";
import type { Column, SelectOption } from "./components/DataTable/types";

export function getColumns(
  options: SelectOption[],
  maxTasks: number,
): Column[] {
  return [
    {
      name: "index",
      isSortable: false,
      isEditable: false,
      width: 40,
      align: "center",
    },
    {
      name: "fullName",
      isSortable: true,
      isEditable: true,
      width: 120,
      align: "left",
      required: true,
    },
    {
      name: "cipher",
      isSortable: true,
      isEditable: true,
      width: 70,
      align: "left",
      required: true,
    },
    {
      name: "schoolId",
      isSortable: true,
      isEditable: true,
      type: "select",
      options,
      width: 100,
      align: "left",
      required: true,
    },
    {
      name: "className",
      isSortable: true,
      isEditable: true,
      width: 70,
      align: "left",
      required: true,
    },
    {
      name: "tasks",
      isSortable: true,
      isEditable: true,
      isGroup: true,
      width: 50,
      type: "number",
      multiply: maxTasks,
      align: "center",
      defaultValue: 0,
    },
    {
      name: "total",
      isSortable: true,
      width: 50,
      align: "left",
    },
    {
      name: "percent",
      isSortable: true,
      width: 70,
      align: "left",
    },
    {
      name: "rating",
      isSortable: true,
      width: 60,
      align: "left",
    },
    {
      name: "status",
      isSortable: false,
      width: 80,
      align: "left",
    },
  ];
}

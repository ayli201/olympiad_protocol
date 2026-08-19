import type { Column, ColumnType, Row } from "./types";

export function mapItemsToRowsForDataTable({
  items,
  columns,
  primaryKey,
  draft,
  isValidating = false,
}: {
  items: Record<string, any>[];
  columns: Record<string, Column>;
  primaryKey: string;
  draft?: Record<string, any>;
  isValidating: boolean;
}): Row[] | { rows: Row[]; isValid: boolean } {
  // console.log(items);
  const rows: Row[] = [];

  let isValid = true;

  for (const item of items) {
    const res = mapItemToRowForDataTable(
      item,
      columns,
      primaryKey,
      isValidating,
      !item?.[primaryKey],
    );
    if (isValidating && "isValid" in res) {
      rows.push(res.row as Row);
      if (!res.isValid) isValid = false;
    } else {
      rows.push(res as Row);
    }
  }

  const newRows = processDrafItems(
    items,
    rows,
    draft,
    columns,
    primaryKey,
    isValidating,
  );
  // console.log(rows);
  if (isValidating) {
    return { rows: newRows, isValid };
  }
  return newRows;
}

function processDrafItems(
  items: Record<string, any>[],
  rows: Row[],
  draft: Record<string, any> | undefined,
  columns: Record<string, Column>,
  primaryKey: string,
  isValidating: boolean,
): Row[] {
  // let res = processDraftCreated(rows, draft, columns, primaryKey, isValidating);
  let res = processDraftUpdated(
    items,
    rows,
    draft,
    columns,
    primaryKey,
    isValidating,
  );
  // res = processDraftDeleted(res, draft, columns, primaryKey, isValidating);
  return res;
}

function processDraftUpdated(
  items: Record<string, any>[],
  rows: Row[],
  draft: Record<string, any> | undefined,
  columns: Record<string, Column>,
  primaryKey: string,
  isValidating: boolean,
): Row[] {
  if (
    !draft?.updated ||
    (draft?.updated && Object.keys(draft.updated).length === 0)
  )
    return rows;

  for (const id in draft.updated) {
    const itemIndex = rows.findIndex((row) => +row.id === +id);
    if (itemIndex === -1) continue;
    // console.log(itemIndex);
    const updatedItem = draft.updated[id];
    const newRow = mapItemToRowForDataTable(
      items[itemIndex],
      columns,
      primaryKey,
      isValidating,
      false,
      updatedItem,
    );

    if (isValidating && "isValid" in newRow) {
      rows[itemIndex] = newRow.row as Row;
      // if (!res.isValid) isValid = false;
    } else {
      rows[itemIndex] = newRow as Row;
    }
  }
  return rows;
}

export function mapItemToRowForDataTable(
  item: Record<string, any>,
  columns: Record<string, Column>,
  primaryKey: string,
  isValidating: boolean = false,
  isDraft: boolean = false,
  updated: Record<string, any> = {},
): Row | { row: Row; isValid: boolean } {
  // console.log(items);
  const row: Row = {
    id: item?.[primaryKey] ?? 0,
    cells: [],
    isEditing: false,
  };
  let isValid = true;

  for (const key in item) {
    let value = item[key];
    if (Array.isArray(value)) {
      for (let idx = 0; idx < value.length; idx++) {
        const itemValue = value[idx];
        const isItemEmpty = isEmpty(itemValue, columns[key]?.type);
        // if (isItemEmpty) continue;
        const isItemValid =
          isValidating && columns[key]?.required ? !isEmpty : true;
        if (!isItemValid) isValid = false;

        const fieldDraft = updated?.[key + "_" + idx];

        const cell = {
          id: itemValue.id ?? 0,
          parentName: key,
          columnName: key + "_" + idx,
          value: itemValue.value,
          type: columns[key]?.type ?? "text",
          align: columns[key]?.align ?? "left",
          isEditable: columns[key]?.isEditable ?? false,
          isError: !isItemValid,
          isDraft: isDraft || fieldDraft || false,
        };
        row.cells.push(cell);
      }
    } else {
      const isItemValid =
        isValidating && columns[key]?.required
          ? !isEmpty(item?.[key], columns[key]?.type)
          : true;
      if (!isItemValid) isValid = false;

      const fieldDraft = updated?.[key];

      const cell = {
        id: item?.[primaryKey] ?? 0,
        columnName: key,
        value,
        type: columns[key]?.type ?? "text",
        align: columns[key]?.align ?? "left",
        isEditable: columns[key]?.isEditable ?? false,
        isError: !isItemValid,
        isDraft: isDraft || fieldDraft || false,
      };
      row.cells.push(cell);
    }
  }

  if (isValidating) {
    return { row, isValid };
  }
  return row;
}

function isEmpty(value: any, type: ColumnType = "text"): boolean {
  if (value === undefined || value === null) return true;
  if (type === "select") return value === 0;
  if (type === "number") return value === "";
  return value === "";
}

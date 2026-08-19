export type ColumnType = "text" | "number" | "date" | "select";

export interface SelectOption {
  label: string;
  value: string | number;
}

export interface Column {
  name: string;
  isSortable?: boolean;
  isResizable?: boolean;
  isEditable?: boolean;
  isVisible?: boolean;
  multiply?: number;
  align?: "left" | "center" | "right";
  type?: ColumnType;
  options?: SelectOption[];
  width?: number;
  isGroup?: boolean;
  calculate?: (row: Row, rawRows: Record<string, any>[]) => void;
  required?: boolean;
  defaultValue?: any;
  onFormat?: (
    value: any,
    cell: Cell,
    row: Row,
    rawRows: Record<string, any>[],
  ) => any;
}

export interface CellBase {
  value: any;
  rowspan?: number;
  colspan?: number;
  isGroup?: boolean;
  isHidden?: boolean;
  columnName?: string;
  parentName?: string;
  isError?: boolean;
  isDraft?: boolean;
}

export interface Cell extends CellBase {
  id: string;
  isEditable?: boolean;
  type: ColumnType;
}

export interface HeaderCell extends CellBase {
  width?: number;
  minWidth?: number;
  maxWidth?: number;
}

export interface Row {
  id: string;
  cells: Cell[];
  isEditing: boolean;
}

export interface HeaderRow {
  cells: HeaderCell[];
}

export enum SortOrder {
  asc = 1,
  desc = -1,
}

export interface SortState {
  columnName: string;
  direction: SortOrder;
}

export interface SaveResult {
  success: boolean;
  message?: string;
  data?: Record<string, any>[];
}

export interface Status {
  message?: string;
  type?: "success" | "error" | "warning" | "info";
}

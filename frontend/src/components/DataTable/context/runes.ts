import { DataTable } from "./dataTableContext.svelte";
import { getContext, setContext } from "svelte";
import type { Column, HeaderRow, Row, SaveResult } from "../types";

const STATE_KEY = Symbol("table_state");

export function createContext({
  headerRows = () => [],
  // items = () => [],
  fetchItems,
  fetchDraft,
  onDraftDelete,
  onRowRecalculate,
  columns = () => undefined,
  isSortable = () => true,
  isPagination = () => true,
  onSort,
  onInput,
  onRowClick,
  onSave,
  onDraftSave,
  // onRowAdd,
  // onRowDelete,
  onPaginationClick,
  primaryKey = () => "id",
  sortKey = () => "id",
  isRemovable = () => false,
  isAddable = () => false,
  isEditable = () => false,
}: {
  headerRows?: () => HeaderRow[];
  // items?: () => Record<string, any>[];
  fetchItems: () => () => Promise<Record<string, any>[]>;
  fetchDraft?: () => () => Promise<Record<string, any>>;
  onDraftDelete?: () => () => Promise<void>;
  columns: () => Column[] | undefined;
  isSortable?: () => boolean;
  isPagination?: () => boolean;
  onSort?: () => ({
    columnName,
    context,
    parentKey,
  }: {
    columnName: string;
    context: any;
    parentKey?: string;
  }) => void | Promise<void>;
  onSave?: () => ({
    created,
    updated,
    deletedIds,
  }: {
    created: Record<string, any>[];
    updated: Record<string, any>[];
    deletedIds: number[];
  }) => SaveResult | Promise<SaveResult>;
  onDraftSave?: () => ({
    created,
    updated,
    deletedIds,
  }: {
    created: Record<string, any>[];
    updated: Record<string, any>;
    deletedIds: number[];
  }) => SaveResult | Promise<SaveResult>;
  // draft?: () => Record<string, any>;
  onInput?: () => ({
    row,
    columnName,
    value,
    index,
    rowsUpdate,
    rawRow,
    rawRows,
    groupIndex,
  }: {
    row: Row;
    columnName: string;
    value: any;
    index: number;
    rowsUpdate: (rows: Record<string, any>[]) => void;
    rawRow: Record<string, any>;
    rawRows: Record<string, any>[];
    groupIndex?: number;
  }) => void | Promise<void>;
  onRowClick?: () => ({
    row,
    rawRow,
  }: {
    row: Row;
    rawRow: Record<string, any>;
  }) => void | Promise<void>;
  // onRowAdd?: () => ({
  //   rawRow,
  //   rawRows,
  //   rowsUpdate,
  //   clearNewRow,
  // }: {
  //   rawRow: Record<string, any>;
  //   rawRows: Record<string, any>[];
  //   rowsUpdate: (rows: Record<string, any>[]) => void;
  //   clearNewRow: () => void;
  // }) => void | Promise<void>;
  // onRowDelete?: () => ({
  //   index,
  //   id,
  // }: {
  //   index: number;
  //   id: string;
  // }) => void | Promise<void>;
  onPaginationClick?: () => ({
    page,
  }: {
    page: number;
  }) => void | Promise<void>;
  onRowRecalculate?: () => ({
    rawRow,
  }: {
    rawRow: Record<string, any>;
  }) => Record<string, any>;
  primaryKey?: () => string;
  sortKey?: () => string;
  isRemovable?: () => boolean;
  isAddable?: () => boolean;
  isEditable?: () => boolean;
}) {
  const tableInstance = new DataTable({
    headerRows: headerRows(),
    // items: items(),
    fetchItems: fetchItems(),
    columns: columns(),
    isSortable: isSortable(),
    isPagination: isPagination(),
    onSort: onSort ? onSort() : undefined,
    onInput: onInput ? onInput() : undefined,
    onRowClick: onRowClick ? onRowClick() : undefined,
    onSave: onSave ? onSave() : undefined,
    onDraftSave: onDraftSave ? onDraftSave() : undefined,
    onRowRecalculate: onRowRecalculate ? onRowRecalculate() : undefined,
    onDraftDelete: onDraftDelete ? onDraftDelete() : undefined,
    fetchDraft: fetchDraft ? fetchDraft() : undefined,
    // onRowAdd: onRowAdd ? onRowAdd() : undefined,
    // onRowDelete: onRowDelete ? onRowDelete() : undefined,
    onPaginationClick: onPaginationClick ? onPaginationClick() : undefined,
    primaryKey: primaryKey(),
    sortKey: sortKey(),
    isRemovable: isRemovable(),
    isAddable: isAddable(),
    isEditable: isEditable(),
  });
  setContext(STATE_KEY, tableInstance);
  // tableInstance.handleFetch();
  return tableInstance;
}

export function getDataTableState(): DataTable<any> {
  const state = getContext<DataTable<any>>(STATE_KEY);
  if (!state) {
    throw new Error(
      "useDataTableState must be used within a parent component initializing it",
    );
  }
  return state;
}

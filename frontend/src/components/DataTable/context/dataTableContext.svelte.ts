import { mapItemsToRowsForDataTable, mapItemToRowForDataTable } from "../utils";
import type {
  Column,
  HeaderRow,
  Row,
  SortState,
  SelectOption,
  SaveResult,
} from "../types";

import { SortOrder, type Status } from "../types";

export class DataTable<T extends Record<string, any>> {
  primaryKey: string = "id";
  #headerRows: HeaderRow[] = $state([]);
  #rawRows: Record<string, any>[] = $state([]);
  #rawColumns?: Column[] = $state();
  #status: Status | null = $state(null);
  #statusTimer: any | null = null;
  #draftSaveTimer: any | null = null;
  // #cachedPromises = $state<Record<string, Promise<any>>>({});
  #cachedColumnOptions: Record<string, SelectOption[]> = $state({});
  // #columns: Record<string, any> = $state({});
  #sortState: SortState = $state({
    direction: SortOrder.asc,
    columnName: "id",
  });
  isInitialLoading?: boolean = $state();
  isLoading: boolean = $state(false);
  isRemovable: boolean = true;
  isAddable: boolean = true;
  isEditable: boolean = true;
  #isSortable: boolean = false;
  #isPagination: boolean = false;
  #isValidating: boolean = $state(false);
  #isNewRowValidating: boolean = $state(false);
  #totalPages: number = 0;
  #rowsPerPage: number = 2;
  #currentPage: number = 1;
  #originalData: Record<string, any>[] = $state([]);
  searchQuery = $state("");
  // #paginatedData: Row[] = $state([]);
  newRowRaw?: Record<string, any> = $state();
  #draft?: Record<string, any>;

  deletedIds: number[] = $state([]);
  // defaultRawValue: Record<string, any> = $state({});
  // newRowRaw: Record<string, any> = $state({});

  fetchDraft?: () => Promise<Record<string, any>>;
  onRowRecalculate?: ({
    rawRow,
  }: {
    rawRow: Record<string, any>;
  }) => Record<string, any>;
  onDraftDelete?: () => Promise<void>;

  get status(): Status | null {
    return this.#status;
  }

  get draft(): Record<string, any> | undefined {
    return this.#draft;
  }

  get updatedItems() {
    return this.paginatedDataRaw.filter((row) => {
      if (!row[this.primaryKey]) return false; // Новые элементы отфильтруем отдельно
      const original = this.#originalData.find(
        (o) => o[this.primaryKey] === row[this.primaryKey],
      );
      // Глубокое сравнение (простой вариант через JSON, либо по полям)
      return JSON.stringify(original) !== JSON.stringify(row);
    });
  }

  get createdItems() {
    return this.paginatedDataRaw.filter((row) => !row[this.primaryKey]);
  }

  fetchItems: () => Promise<Record<string, any>[]>;
  onSort?: ({
    columnName,
    context,
    parentKey,
  }: {
    columnName: string;
    context: any;
    parentKey?: string;
  }) => void | Promise<void>;

  constructor({
    headerRows = [],
    // items = [],
    fetchItems,
    columns,
    isSortable = true,
    isPagination = false,
    onSort,
    onInput,
    onRowClick,
    onPaginationClick,
    onRowRecalculate,
    // onRowAdd,
    // onRowDelete,
    onSave,
    onDraftSave,
    fetchDraft,
    onDraftDelete,
    primaryKey = "id",
    sortKey = "id",
    isRemovable = true,
    isAddable = true,
    isEditable = true,
  }: {
    headerRows?: HeaderRow[];
    // items?: Record<string, any>[];
    fetchItems: () => Promise<Record<string, any>[]>;
    columns?: Column[];
    isSortable?: boolean;
    isPagination?: boolean;
    onRowRecalculate?: ({
      rawRow,
    }: {
      rawRow: Record<string, any>;
    }) => Record<string, any>;
    onSort?: ({
      columnName,
      context,
      parentKey,
    }: {
      columnName: string;
      context: any;
      parentKey?: string;
    }) => void;
    onInput?: ({
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
    onRowClick?: ({
      row,
      rawRow,
    }: {
      row: Row;
      rawRow: Record<string, any>;
    }) => void | Promise<void>;
    onSave?: ({
      created,
      updated,
      deletedIds,
    }: {
      created: Record<string, any>[];
      updated: Record<string, any>[];
      deletedIds: number[];
    }) => SaveResult | Promise<SaveResult>;
    onDraftSave?: ({
      created,
      updated,
      deletedIds,
    }: {
      created: Record<string, any>[];
      updated: Record<string, any>;
      deletedIds: number[];
    }) => SaveResult | Promise<SaveResult>;
    fetchDraft?: () => Promise<Record<string, any>>;
    onDraftDelete?: () => Promise<void>;
    // onRowAdd?: ({
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
    onPaginationClick?: ({ page }: { page: number }) => void | Promise<void>;
    // onRowDelete?: ({
    //   index,
    //   id,
    // }: {
    //   index: number;
    //   id: string;
    // }) => void | Promise<void>;
    primaryKey: string;
    sortKey?: string;
    isRemovable?: boolean;
    isAddable?: boolean;
    isEditable?: boolean;
  }) {
    this.headerRows = headerRows;
    // this.#rawRows = items;
    this.fetchItems = fetchItems;

    this.#rawColumns = columns;
    // this.columns = columns;

    this.#isSortable = isSortable;
    this.#isPagination = isPagination;
    this.onSort = onSort;
    this.onInput = onInput;
    this.onRowClick = onRowClick;
    this.onSave = onSave;
    this.onDraftSave = onDraftSave;
    this.onDraftDelete = onDraftDelete;
    this.onRowRecalculate = onRowRecalculate;
    this.fetchDraft = fetchDraft;
    this.isRemovable = isRemovable;
    this.isAddable = isAddable;
    this.isEditable = isEditable;
    // this.onRowAdd = onRowAdd;
    this.onPaginationClick = onPaginationClick;
    // this.onDelete = onRowDelete;
    this.primaryKey = primaryKey;
    this.#sortState = {
      direction: SortOrder.asc,
      columnName: sortKey,
    };
    this.handleFetch();
    // this.handleFetchDraft();

    // this.#totalPages = Math.ceil(items.length / this.#rowsPerPage);
    this.clearNewRow();
  }

  // handleFetchData = async () => {
  //   const res = await this.handleFetch();
  //   this.handleFetchDraft();
  // };

  clearNewRow = () => {
    // this.newRowRaw = Object.keys(this.newRowRaw).reduce((res, key) => {
    //   res[key] = undefined;
    //   return res;
    // }, this.newRow);
    if (!this.#rawColumns) return;
    this.newRowRaw = this.#rawColumns.reduce(
      (res: Record<string, any>, item) => {
        res[item.name] = item.isGroup
          ? Array.from({ length: item.multiply ?? 0 }, (_, index) => ({
              id: 0,
              number: index + 1,
              value: item.defaultValue ?? "",
            }))
          : item.type === "select"
            ? 0
            : "";
        return res;
      },
      { uuid: crypto.randomUUID() },
    );
  };

  get newRow(): Row | null {
    if (!this.newRowRaw) return null;
    const res = mapItemToRowForDataTable(
      this.newRowRaw,
      this.columns,
      this.primaryKey,
      this.#isNewRowValidating,
    );
    if ("isValid" in res) return res.row;
    return res;
  }

  get filteredRows(): Record<string, any>[] {
    if (!this.searchQuery) return this.#rawRows;
    return this.#rawRows.filter((row) =>
      Object.values(row).some((value) =>
        value
          ?.toString()
          .toLowerCase()
          .includes(this.searchQuery.toLowerCase()),
      ),
    );
  }

  get paginatedDataRaw(): Record<string, any>[] {
    if (!this.#isPagination) return this.filteredRows || [];
    const startIndex = (this.#currentPage - 1) * this.#rowsPerPage;
    return this.filteredRows.slice(startIndex, startIndex + this.#rowsPerPage);
  }

  get paginatedData(): Row[] {
    if (!this.#isPagination) return this.rows;
    const startIndex = (this.#currentPage - 1) * this.#rowsPerPage;
    return this.rows.slice(startIndex, startIndex + this.#rowsPerPage);
  }

  get rowsPerPage(): number {
    return this.#rowsPerPage;
  }

  get totalPages(): number {
    return this.#totalPages;
  }

  get currentPage(): number {
    return this.#currentPage;
  }

  get isSortable(): boolean {
    return this.#isSortable;
  }

  get isPagination(): boolean {
    return this.#isPagination;
  }

  get sortState(): SortState {
    return this.#sortState;
  }

  get headerRows(): HeaderRow[] {
    return this.#headerRows;
  }

  get rows(): Row[] {
    if (!this.#rawColumns) return [];

    const result = mapItemsToRowsForDataTable({
      items: this.filteredRows,
      columns: this.columns,
      primaryKey: this.primaryKey,
      isValidating: this.#isValidating,
      draft: this.#draft,
    });

    if ("isValid" in result) {
      return result.rows as Row[];
    }

    return result as Row[];
  }

  get rawRows(): Record<string, any>[] {
    return this.#rawRows;
  }

  get columns(): Record<string, Column> {
    if (!this.#rawColumns) return {};

    return this.#rawColumns.reduce(
      (acc, column) => {
        acc[column.name] = column;
        return acc;
      },
      {} as Record<string, Column>,
    );
  }
  // get columns() {
  //   // return this.#rawColumns;
  //   console.log(this.#rawColumns);
  //   return this.#rawColumns.map((col) => {
  //     if (!col.getOptionsAsync) return col;

  //     return {
  //       ...col,
  //       // Возвращаем функцию, которая выдает закэшированный промис для строки
  //       getOptionsFor: (row: Record<string, any>) => {
  //         const cacheKey = `${col.name}`;

  //         // Считываем реактивное значение, от которого зависит список (например, страна)
  //         // Svelte запомнит эту зависимость!
  //         // const triggerValue = row[col.dependsOn];

  //         // Если промиса для этой строки еще нет ИЛИ изменился триггер (например, выбрали другую страну)
  //         if (!this.#cachedPromises[cacheKey] && !!col.getOptionsAsync) {
  //           // Вызываем getOptionsAsync ОДИН раз и сохраняем промис
  //           this.#cachedPromises[cacheKey] = col.getOptionsAsync(col.name);
  //         }

  //         return this.#cachedPromises[cacheKey];
  //       },
  //     };
  //   });
  // }

  get isSaveActive(): boolean {
    return (
      this.createdItems.length > 0 ||
      this.updatedItems.length > 0 ||
      this.deletedIds.length > 0
    );
  }

  get rawColumns(): Column[] | undefined {
    return this.#rawColumns;
  }

  // get newRow(): Record<string, any> {
  //   return this.#newRow;
  // }

  set status(status: Status) {
    this.#status = status;
    if (this.#statusTimer) clearTimeout(this.#statusTimer);
    this.#statusTimer = setTimeout(() => {
      this.#status = null;
    }, 3000);
  }

  set draft(draft: Record<string, any> | undefined) {
    this.#draft = draft;
    // this.handleDraftUpdated(draft);
  }

  set headerRows(rows: HeaderRow[]) {
    this.#headerRows = rows;
  }

  set columns(columns: Column[]) {
    this.#rawColumns = columns;
  }

  set totalPages(totalPages: number) {
    this.#totalPages = totalPages;
  }

  set currentPage(currentPage: number) {
    this.#currentPage = currentPage;
  }

  set rowsPerPage(rowsPerPage: number) {
    this.#rowsPerPage = rowsPerPage;
    this.#totalPages = Math.ceil(this.rows.length / rowsPerPage);
  }

  handleDraftUpdated = (
    draft: Record<string, any> | undefined,
    prev?: Record<string, any>[],
  ) => {
    // console.log(draft);
    let newRows = prev ?? this.#rawRows;
    if (draft?.created?.length) {
      // console.log(draft.created, [...this.#rawRows, ...draft.created]);
      newRows = [...newRows, ...draft.created];
    }
    if (draft?.updated) {
      newRows = newRows.map((row) => {
        let updated = draft.updated[row[this.primaryKey]];
        if (updated && Object.keys(updated).some((key) => key.includes("_"))) {
          updated = Object.keys(updated).reduce(
            (acc: Record<string, any>, key) => {
              const groupKey = key.split("_");
              const parentKey = groupKey[0];
              const childIndex = groupKey?.[1];
              if (!childIndex) {
                acc[key] = updated[key];
              } else {
                if (!acc[parentKey]) acc[parentKey] = row[parentKey] ?? [];
                acc[parentKey][childIndex] = {
                  ...row[parentKey][childIndex],
                  value: updated[key],
                };
              }
              return acc;
            },
            {},
          );
        }
        let updatedRow = updated ? { ...row, ...updated } : row;
        if (this.onRowRecalculate && updated)
          updatedRow = this.onRowRecalculate({
            rawRow: updatedRow,
          });
        return updatedRow;
      });
    }

    if (draft?.deleted?.length) {
      newRows = newRows.filter(
        (row) => !draft.deleted.includes(row[this.primaryKey]),
      );
    }
    this.#rawRows = newRows;
  };

  draftHandler = async () => {
    if (!this.onDraftSave) return;

    try {
      await this.onDraftSave({
        created: this.createdItems,
        updated: this.#draft?.updated ?? {},
        deletedIds: this.deletedIds,
      });
    } catch (error) {
      if (error instanceof Error)
        this.status = { type: "error", message: error.message };
    }
  };

  handleSaveDraft = async (updatedItem?: Record<string, any>) => {
    if (this.#draftSaveTimer) clearTimeout(this.#draftSaveTimer);

    let updated = null;
    let created = null;
    if (updatedItem) {
      if (!updatedItem[this.primaryKey]) {
        created = this.#draft?.created ?? [];
        const editedIndex = created.findIndex(
          (el: Record<string, any>) => el.uuid == updatedItem.uuid,
        );
        // console.log(created);
        if (editedIndex !== -1)
          created[editedIndex] = {
            ...created[editedIndex],
            [updatedItem.columnName]: updatedItem.value,
          };
        // created = created.filter(
        //   (el: Record<string, any>) => el.uuid !== updatedItem.uuid,
        // );
      } else {
        updated = {
          ...(this.#draft?.updated ?? {}),
          [updatedItem[this.primaryKey]]: {
            ...(this.#draft?.updated?.[updatedItem[this.primaryKey]] ?? {}),
            [updatedItem.columnName]: updatedItem.value,
          },
        };
      }
    }
    // console.log({
    //   ...(this.#draft ?? {}),
    //   ...(updatedItem
    //     ? {
    //         ...updated,
    //         ...(created ? { created } : {}),
    //       }
    //     : {}),
    // });
    this.#draft = {
      ...(this.#draft ?? {}),
      ...(updatedItem
        ? {
            ...{ updated },
            ...(created ? { created } : {}),
          }
        : {}),
    };
    this.#draftSaveTimer = setTimeout(() => this.draftHandler(), 400);
  };

  // set newRow(newRow: Record<string, any>) {
  //   this.#newRow = newRow;
  // }

  handleSearch = (event: Event) => {
    event.preventDefault();

    const target = event.target as HTMLInputElement;
    const text = target.value;

    if (typeof this.onSearch === "function") {
      this.onSearch({ text });
    }

    this.searchQuery = text;
  };

  handleClear = (withoutDb: boolean = false) => {
    // this.#rawRows = [];
    this.deletedIds = [
      ...this.deletedIds,
      ...this.#rawRows.reduce(
        (acc: number[], row) =>
          row[this.primaryKey] > 0 ? [...acc, row[this.primaryKey]] : acc,
        [],
      ),
    ];
    if (this.#draft)
      this.#draft.deleted = [
        ...this.deletedIds,
        ...this.#rawRows.map((row) => row[this.primaryKey]),
      ];
    if (this.#draft) this.#draft.created = [];
    if (this.#draft) this.#draft.updated = [];
    if (!withoutDb) {
      this.draftSave();
    }
    this.#rawRows = [];
    //this.#rawRows = [];
  };

  draftSave = () => {
    if (this.#draftSaveTimer) clearTimeout(this.#draftSaveTimer);
    this.#draftSaveTimer = setTimeout(() => this.draftHandler(), 400);
  };

  handleCancel = async (isUpdate?: boolean) => {
    this.#draft = undefined;

    if (!isUpdate && this.onDraftDelete) await this.onDraftDelete();
    const originalData = $state.snapshot(this.#originalData);
    this.#rawRows = structuredClone(originalData);
    // this.handleDraftUpdated(this.#draft, this.#rawRows);
    this.deletedIds = [];
  };

  handleFetchDraft = async (items: Record<string, any>[]) => {
    if (!this.fetchDraft) return;
    try {
      const res = await this.fetchDraft();
      this.#draft = res;
      this.handleDraftUpdated(res, items);
    } catch (error) {
      if (error instanceof Error)
        this.status = { message: error.message, type: "error" };
    }
  };

  handleFetch = async () => {
    if (this.isInitialLoading === undefined) this.isInitialLoading = true;
    else this.isLoading = true;

    return this.fetchItems()
      .then((data) => {
        this.handleCancel(true);
        if (this.fetchDraft) this.handleFetchDraft(data);
        else this.#rawRows = data;
        this.#originalData = structuredClone(data);
        // this.#draft = undefined;
        // console.log(data);
        this.#totalPages = Math.ceil(data.length / this.#rowsPerPage);
      })
      .catch((error) => {
        this.status = { message: error.message, type: "error" };
      })
      .finally(() => {
        if (this.isInitialLoading) this.isInitialLoading = false;
        else this.isLoading = false;
      });
  };

  handleDelete = async (index: number) => {
    if (this.#rawRows[index][this.primaryKey]) {
      this.deletedIds = [
        ...this.deletedIds,
        this.#rawRows[index][this.primaryKey],
      ];
      if (this.#draft) this.#draft.deleted = [...this.deletedIds];
    } else {
      // console.log(this.#rawRows[index].uuid);
      const idx = this.#draft?.created.findIndex(
        (row: Record<string, any>) => this.#rawRows[index].uuid === row.uuid,
      );
      if (idx !== -1) {
        this.#draft?.created.splice(idx, 1);
      }
    }
    if (this.#draftSaveTimer) clearTimeout(this.#draftSaveTimer);
    this.#draftSaveTimer = setTimeout(() => this.draftHandler(), 400);

    this.#rawRows.splice(index, 1);
  };

  handleAdd = async () => {
    if (!this.newRowRaw) return;

    this.#isNewRowValidating = true;

    try {
      const res = mapItemToRowForDataTable(
        this.newRowRaw,
        this.columns,
        this.primaryKey,
        true,
      );
      if ("isValid" in res && !res.isValid) {
        this.status = {
          type: "warning",
          message: "Проверьте что все поля заполнены корректно",
        };
        return;
      }
      const cleanRowRaw = $state.snapshot(this.newRowRaw);
      this.#rawRows.push(cleanRowRaw);
      this.clearNewRow();
      this.#isNewRowValidating = false;
      this.handleSaveDraft();
    } catch (err) {
      this.status = {
        type: "error",
        message: "Ошибка при добавлении строки",
      };
    }

    // if (typeof this.onRowAdd === "function") {
    //   try {
    //     this.clearNewRow();
    //   } catch (err) {
    //     // Сюда попадет ошибка из async-функции или throw из обычной функции
    //     console.error("Ошибка при добавлении строки:", err);
    //     return; // Прерываем выполнение, если нужно не очищать строку при ошибке
    //   }
    // }
  };

  handleSave = async () => {
    if (!this.onSave) return;

    this.#isValidating = true;

    // Валидация данных перед сохранением
    const res = mapItemsToRowsForDataTable({
      items: this.#rawRows,
      columns: this.columns,
      primaryKey: this.primaryKey,
      isValidating: this.#isValidating,
    });
    // console.log(res);
    if ("isValid" in res && !res.isValid) {
      this.status = {
        type: "warning",
        message: "Проверьте что все поля заполнены корректно",
      };
      return;
    }

    try {
      const res = await this.onSave({
        created: this.createdItems,
        updated: this.updatedItems,
        deletedIds: this.deletedIds,
      });
      if (res.success && res?.data) {
        // this.handleCancel();
        this.#rawRows = res.data;
        this.#originalData = structuredClone(res.data);
        this.handleCancel();
        this.status = { type: "success", message: "Сохранено" };
        this.deletedIds = [];

        return;
      }
      this.#isValidating = false;
    } catch (error) {
      if (error instanceof Error)
        this.status = { type: "error", message: error.message };
    }
  };

  sortByKey = (key: string, parentKey?: string) => {
    if (!this.#isSortable) return;

    if (this.#sortState.columnName === key) {
      this.#sortState = {
        direction:
          this.#sortState.direction === SortOrder.asc
            ? SortOrder.desc
            : SortOrder.asc,
        columnName: key,
      };
    } else {
      this.#sortState = {
        direction: SortOrder.asc,
        columnName: key,
      };
    }

    this.#rawRows = this.#rawRows.sort((a, b) => {
      let aItem = a[parentKey || key];
      let bItem = b[parentKey || key];
      console.log(aItem, bItem);
      if (parentKey) {
        const keyArr = key.split("_");
        const keyIndex = keyArr[keyArr.length - 1];
        aItem = a[parentKey]?.[keyIndex]?.value;
        bItem = b[parentKey]?.[keyIndex]?.value;
      }
      if (aItem < bItem) return -1 * this.#sortState.direction;
      if (aItem > bItem) return 1 * this.#sortState.direction;
      return 0;
    });
  };

  rowsUpdate = (newRows: Record<string, any>[]) => {
    this.#rawRows = newRows;
  };

  mapUUID = (newRows: Record<string, any>[]) => {
    return newRows.map((row) => {
      return { ...row, uuid: crypto.randomUUID() };
    });
  };

  onSave?: ({
    created,
    updated,
    deletedIds,
  }: {
    created: Record<string, any>[];
    updated: Record<string, any>[];
    deletedIds: number[];
  }) => SaveResult | Promise<SaveResult>;

  onDraftSave?: ({
    created,
    updated,
    deletedIds,
  }: {
    created: Record<string, any>[];
    updated: Record<string, any>;
    deletedIds: number[];
  }) => SaveResult | Promise<SaveResult>;

  // onDelete?: ({
  //   index,
  //   id,
  // }: {
  //   index: number;
  //   id: string;
  // }) => void | Promise<void>;

  onInput?: ({
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

  onRowClick?: ({
    row,
    rawRow,
  }: {
    row: Row;
    rawRow: Record<string, any>;
  }) => void | Promise<void>;

  // onRowAdd?: ({
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

  onPaginationClick?: ({ page }: { page: number }) => void | Promise<void>;

  onSearch?: ({ text }: { text: string }) => void | Promise<void>;
}

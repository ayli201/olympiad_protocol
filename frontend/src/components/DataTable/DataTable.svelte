<script lang="ts">
    import DataTableFooter from "./DataTableFooter.svelte";
    import DataTableHeader from "./DataTableHeader.svelte";
    // import DataTableCaption from "./DataTableCaption.svelte";
    // import DataTableRow from "./DataTableRow.svelte";
    import type { Column, Row, HeaderRow, SaveResult } from "./types";
    import { createContext } from "./context/runes";
    import DataTableBody from "./DataTableBody.svelte";
    import DataTablePagination from "./DataTablePagination.svelte";
    import Icon from "../Icon.svelte";
    import type { Snippet } from "svelte";
    import { untrack } from "svelte";
    // import { select_multiple_value } from "svelte/internal";

    let {
        headerRows = [],
        onRowClick,
        onSave,
        onDraftSave,
        fetchDraft,
        onDraftDelete,
        onRowRecalculate,
        // onRowAdd = undefined,
        // onRowDelete = undefined,
        onInput,
        onPaginationClick,
        onSort,
        isPagination = true,
        isSortable = true,
        primaryKey = "id",
        sortKey = "id",
        isRemovable = true,
        isAddable = true,
        isEditable = true,
        columns,
        items = [],
        fetchItems,
        header,
        buttonSaveText = "сохранить",
    }: {
        headerRows?: HeaderRow[];
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
        onRowRecalculate?: ({
            rawRow,
        }: {
            rawRow: Record<string, any>;
        }) => Record<string, any>;
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
        onPaginationClick?: ({
            page,
        }: {
            page: number;
        }) => void | Promise<void>;
        onSort?: ({
            columnName,
            context,
            parentKey,
        }: {
            columnName: string;
            context: any;
            parentKey?: string;
        }) => void | Promise<void>;
        isPagination?: boolean;
        isSortable?: boolean;
        primaryKey?: string;
        sortKey?: string;
        isRemovable?: boolean;
        isAddable?: boolean;
        isEditable?: boolean;
        columns: Column[];
        items?: Record<string, any>[];
        fetchItems: () => Promise<Record<string, any>[]>;
        header: Snippet;
        buttonSaveText?: string;
    } = $props();

    let context = createContext({
        headerRows: () => headerRows,
        // items: () => items,
        fetchItems: () => fetchItems,
        columns: () => columns,
        isSortable: () => isSortable,
        isPagination: () => isPagination,
        fetchDraft: fetchDraft && (() => fetchDraft),
        onDraftDelete: onDraftDelete && (() => onDraftDelete),
        onSort: onSort && (() => onSort),
        onInput: onInput && (() => onInput),
        onRowClick: onRowClick && (() => onRowClick),
        // onRowAdd: () => onRowAdd,
        // onRowDelete: () => onRowDelete,
        onSave: onSave && (() => onSave),
        onDraftSave: onDraftSave && (() => onDraftSave),
        onPaginationClick: onPaginationClick && (() => onPaginationClick),
        onRowRecalculate: onRowRecalculate && (() => onRowRecalculate),
        primaryKey: () => primaryKey,
        sortKey: () => sortKey,
        isRemovable: () => isRemovable,
        isAddable: () => isAddable,
        isEditable: () => isEditable,
    });

    $effect(() => {
        // $inspect(context.rawColumns);
        context.columns = columns;
        context.clearNewRow();
        // $inspect(context.rawRows, context.primaryKey)
        // console.log(context.columns);
    });

    $effect(() => {
        if (items?.length) {
            context.mapUUID(items);
            untrack(() => {
                context.handleClear(true);
                context.draft = {
                    ...context.draft,
                    created: items,
                };
                context.draftSave();
                context.rowsUpdate(items);
            });
            // let count = 0;
            // count++;
            // console.log(count),
        }
    });

    // $effect(() => {
    //     let count = 0;
    //     count++;
    //     if (count < 10) {
    //         // Выводим данные только первые 10 раз
    //         console.log(
    //             "Вызов геттера. Текущее состояние:",
    //             $state.snapshot(context.deletedIds),
    //         );
    //     }
    //     if (count === 10) {
    //         console.error("СТОП: Подозрение на бесконечный цикл!");
    //     }
    // });

    // $effect(() => {
    //     // context.draft = draft;
    //     $inspect(context.newRowRaw);
    // });
</script>

<div class="w-full mx-auto">
    <div>
        {@render header()}
    </div>
    <div class="min-h-8">
        {#if context.status}
            <span
                class="flex items-center py-1"
                class:text-green-600={context.status.type === "success"}
                class:text-yellow-600={context.status.type === "warning"}
                class:text-red-600={context.status.type === "error"}
                ><Icon
                    >{context.status.type === "success"
                        ? "done"
                        : "error"}</Icon
                >
                {context.status.message}</span
            >
        {/if}
    </div>
    <div class="flex flex-wrap items-center justify-start gap-2 mb-2">
        {#if context.isEditable || context.isAddable}
            <button
                title="Сохранить"
                disabled={!context.isSaveActive}
                class="btn-success"
                onclick={() => context.handleSave()}
                class:btn-active={context.isSaveActive}
                ><Icon>save</Icon> {buttonSaveText}</button
            >
        {/if}
        <!-- <button
            title="Обновить"
            class="btn-indigo btn-active"
            onclick={context.handleFetch}><Icon>refresh</Icon> Обновить</button
        > -->
        {#if context.isRemovable}
            <button
                title="Очистить"
                class="btn-neutral btn-active"
                onclick={() => context.handleClear()}
                ><Icon>clear</Icon> Очистить</button
            >
        {/if}
        {#if context.isEditable || context.isAddable}
            <button
                title="Отменить"
                class="btn-neutral btn-active"
                onclick={() => context.handleCancel()}
                ><Icon>undo</Icon> Отменить</button
            >
        {/if}

        <input
            class="w-full border border-t-blue-400 rounded-sm p-2 ml-auto"
            type="text"
            placeholder="Поиск..."
            oninput={(e) => context.handleSearch(e)}
        />
    </div>
    {#if !context.rawColumns}
        <p class="m-auto w-20">Загрузка...</p>
    {:else}
        <div class="data-table">
            <table>
                <colgroup>
                    {#each columns as column (column.name)}
                        {#if !column.width}
                            <col width="auto" />
                        {:else if column.multiply && column.multiply > 0}
                            {#each Array.from( { length: column.multiply } ) as _, idx (idx)}
                                <col
                                    style="width: {column.width}px; max-width: {+column.width +
                                        40}px;"
                                />
                            {/each}
                        {:else}
                            <col
                                style="width: {column.width}px; max-width: {+column.width +
                                    40}px;"
                            />
                        {/if}
                    {/each}
                    {#if context.isAddable || context.isRemovable}
                        <col width="40" />
                    {/if}
                </colgroup>
                <DataTableHeader />

                <DataTableBody />

                <DataTableFooter />
            </table>

            {#if context.isPagination && context.totalPages > 1}
                <DataTablePagination />
            {/if}
        </div>
    {/if}
</div>

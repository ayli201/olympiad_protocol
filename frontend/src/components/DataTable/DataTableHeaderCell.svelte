<script lang="ts">
    import type { Column, HeaderCell, HeaderRow, SortState } from "./types";
    import { SortOrder } from "./types";
    import Icon from "../Icon.svelte";
    import { getDataTableState } from "./context/runes";

    const { cell, row }: { cell: HeaderCell; row: HeaderRow } = $props();

    const context = getDataTableState();
    const sortState = $derived(context.sortState);
    const column: Column | null = $derived(
        context.columns[cell.parentName || cell.columnName || ""] ?? null,
    );
    const isSortable = $derived(
        (context.isSortable && column?.isSortable) ?? false,
    );

    function handleSort(key: string, parentKey?: string) {
        context.sortByKey(key, parentKey);
        if (context.onSort) {
            context.onSort({
                columnName: key,
                context: context.sortState,
                parentKey,
            });
        }
    }

    let isActive = $derived(
        sortState && sortState.columnName === cell.columnName,
    );
    // $effect(() => {
    //     $inspect(column);
    //     console.log(cell);
    // });
</script>

<th
    rowspan={cell.rowspan}
    colspan={cell.colspan}
    data-name={cell.parentName || cell.columnName}
    class:cursor-pointer={cell.colspan == 1 || !cell.colspan}
    class:cursor-default={cell.colspan && cell.colspan > 1}
    class={"group border-b border-r border-slate-200 p-1 transition-colors"}
    class:bg-blue-500!={cell.colspan && cell.colspan > 1}
    onclick={() => handleSort(cell.columnName ?? "", cell.parentName)}
>
    <div
        class="flex items-center gap-1 {cell.isGroup
            ? 'justify-center text-center'
            : 'justify-between'}"
    >
        <span
            class="text-[13px]! w-[90%] font-medium whitespace-normal hyphens-auto wrap-break-words"
        >
            {cell.value}
        </span>
        {#if !cell.isGroup && isSortable}
            <Icon
                class="text-[14px]! text-base transition-all duration-200
                {sortState.columnName === cell.columnName
                    ? 'text-white! font-bold opacity-100'
                    : 'text-slate-100 opacity-0 group-hover:opacity-60'}"
            >
                {!isActive || sortState.direction === SortOrder.asc
                    ? "north"
                    : "south"}
            </Icon>
        {/if}
    </div>
</th>

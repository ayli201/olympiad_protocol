<script lang="ts">
    import type { Cell, Column, Row } from "./types";
    import { getDataTableState } from "./context/runes";

    // export let cell: Cell;
    // export let column: Column;
    // export let row: Row;
    // export let rawRow: Record<string, any>;

    let {
        cell,
        column,
        row,
        rawRow,
        index,
        value,
        isError,
        isDraft,
        groupIndex = undefined,
        isNewRow = false,
    } = $props();

    // const isEditing = $derived(row.isEditing && cell.isEditable);

    const context = getDataTableState();

    // export let onChange:
    //     | ((
    //           row: Row,
    //           columnName: string,
    //           value: any,
    //           rawItem: Record<string, any>,
    //       ) => void)
    //     | undefined = undefined;

    function handleInput(
        columnName: string,
        event: Event,
        columnType?: string,
    ) {
        // event.preventDefault();
        // event.stopPropagation();

        const target = event.currentTarget as HTMLInputElement;
        let value: any = target.value;
        if (["number", "select"].includes(columnType ?? "")) {
            value = +value;
        }

        const newRow = isNewRow ? context.newRowRaw : { ...rawRow };
        if (groupIndex !== undefined) {
            if (!newRow[columnName]) newRow[columnName] = [];
            if (!newRow[columnName][groupIndex])
                newRow[columnName][groupIndex] = {};

            newRow[columnName][groupIndex].value = value;
        } else newRow[columnName] = value;

        const draftColumnName =
            columnName + (groupIndex !== undefined ? `_${groupIndex}` : "");
        // console.log(rawRow);
        context.handleSaveDraft({
            [context.primaryKey]: rawRow?.[context.primaryKey],
            uuid: rawRow?.uuid,
            columnName: draftColumnName,
            value,
        });

        if (!isNewRow)
            context.rowsUpdate([
                ...context.rawRows.slice(0, index),
                newRow,
                ...context.rawRows.slice(index + 1),
            ]);

        if (context.onInput) {
            context.onInput({
                row,
                columnName,
                value,
                index,
                rowsUpdate: context.rowsUpdate,
                rawRow,
                rawRows: context.rawRows,
                groupIndex,
            });
        }
    }

    // function handleKeyDown(event: KeyboardEvent) {
    //   // let value = (event.target as HTMLInputElement).value;
    //   if (isNaN(+value)) {
    //     event.preventDefault();
    //     event.stopPropagation();
    //     (event.target as HTMLInputElement).value = value.replace(/[^0-9.]/g, '')
    //   }
    // }
    // $effect(() => {
    //     $inspect(isError);
    // });
</script>

<td
    data-name={column.name}
    class="text-[16px]! text-{column?.align || 'left'} {!column.isEditable
        ? 'px-2! py-1!'
        : 'p-0!'}"
    class:bg-yellow-100={isDraft}
    class:bg-red-100!={isError}
>
    {#if context.isEditable && column?.isEditable}
        {#if column.type === "select"}
            <div class="relative w-full h-full flex items-center">
                <select
                    class="text-[16px]! font-mono"
                    value={`${value}`}
                    onchange={(e) => handleInput(column.name, e, column.type)}
                >
                    <option class="text-[16px]! font-mono" value="0"
                        >—————————</option
                    >
                    {#each column.options || [] as option (option.value)}
                        <option
                            class="text-[16px]! font-mono"
                            value={`${option.value}`}
                        >
                            {option.label}
                        </option>
                    {/each}
                </select>
                <i
                    class="material-icons absolute right-3 pointer-events-none text-slate-400 text-sm"
                    >arrow_drop_down</i
                >
            </div>
        {:else}
            <input
                type={column.type ?? "text"}
                {value}
                oninput={(e) => handleInput(column.name, e, column.type)}
            />
        {/if}
    {:else}
        {value !== undefined && value !== null ? value : ""}
    {/if}
</td>

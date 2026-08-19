<script lang="ts">
    import Icon from "../Icon.svelte";
    import DataTableCell from "./DataTableCell.svelte";
    // import type { Row } from "./types";
    import { getDataTableState } from "./context/runes";
    import { flip } from "svelte/animate";

    const context = getDataTableState();

    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === "Enter") {
            event.preventDefault();
            context.handleAdd();
        }
    }

    // $effect(() => {
    //     $inspect(context.newRow);
    // });
</script>

<tbody>
    {#if context.isAddable && (!!context.newRowRaw || context.isInitialLoading)}
        <tr onkeydown={handleKeyDown}>
            {#each context.rawColumns as column (column.name)}
                {@const cell = context.newRow?.cells.find(
                    (el) => el.columnName == column.name,
                )}
                {#if column.name === "index"}
                    <DataTableCell
                        isNewRow={true}
                        {column}
                        index={0}
                        rawRow={context.newRowRaw}
                        cell={{ value: "+" }}
                        value={``}
                        isError={cell?.isError}
                        isDraft={cell?.isDraft}
                        row={context.newRow}
                    />
                {:else if column.isGroup}
                    {#each [...new Array(column.multiply)] as item, idx (column.name + "_" + idx)}
                        {@const rawGroupValue =
                            context.newRowRaw?.[column.name][idx]}
                        <DataTableCell
                            isNewRow={true}
                            index={0}
                            {column}
                            rawRow={context.newRowRaw}
                            cell={{
                                id: 0,
                                parentName: column.name,
                                columnName: column.name + "_" + idx,
                                value: rawGroupValue.value ?? "",
                                type: column?.type ?? "text",
                                align: column?.align ?? "left",
                            }}
                            isError={cell?.isError}
                            isDraft={cell?.isDraft}
                            value={rawGroupValue.value ?? ""}
                            row={context.newRow}
                            groupIndex={idx}
                        />
                    {/each}
                {:else}
                    <DataTableCell
                        isNewRow={true}
                        {column}
                        index={0}
                        rawRow={context.newRowRaw}
                        cell={{
                            value: context.newRowRaw?.[column.name] ?? "",
                        }}
                        isError={cell?.isError}
                        isDraft={cell?.isDraft}
                        value={`${context.newRowRaw?.[column.name] ?? ""}`}
                        row={context.newRow}
                    />
                {/if}
            {/each}
            <td class="text-center">
                <button
                    title="Добавить"
                    class="btn m-auto text-green-600 cursor-pointer select-none"
                    onclick={() => context.handleAdd()}><Icon>add</Icon></button
                >
            </td>
        </tr>
    {/if}
    {#each context.paginatedData as row, index (index)}
        {@const rawRow = context.paginatedDataRaw[index]}
        <tr
            animate:flip={{ duration: 400 }}
            onclick={() => {
                if (context.onRowClick) {
                    context.onRowClick({
                        row,
                        rawRow: rawRow,
                    });
                }
            }}
        >
            {#each context.rawColumns as column (column.name)}
                {@const rawValue = rawRow[column.name]}
                {@const cell = row.cells.find(
                    (el) => el.columnName == column.name,
                )}
                {#if column.name === "index"}
                    <DataTableCell
                        {row}
                        {index}
                        {rawRow}
                        cell={{ value: index + 1 }}
                        value={`${index + 1}`}
                        isError={false}
                        isDraft={cell?.isDraft}
                        {column}
                    />
                {:else if column.isGroup}
                    {#each [...new Array(column.multiply)] as item, idx (column.name + "_" + idx)}
                        {@const rawGroupValue = rawValue?.[idx]}
                        {@const groupCell = row.cells.find(
                            (el) => el.columnName == column.name + "_" + idx,
                        )}
                        <DataTableCell
                            {row}
                            {index}
                            {rawRow}
                            groupIndex={idx}
                            cell={{
                                id: rawGroupValue?.id ?? 0,
                                parentName: column.name,
                                columnName: column.name + "_" + idx,
                                value: rawGroupValue?.value,
                                type: column?.type ?? "text",
                                align: column?.align ?? "left",
                            }}
                            isError={groupCell?.isError ?? false}
                            isDraft={groupCell?.isDraft ?? false}
                            value={rawGroupValue?.value}
                            {column}
                        />
                    {/each}
                {:else}
                    <DataTableCell
                        {row}
                        {index}
                        {rawRow}
                        cell={{
                            columnName: column.name,
                            type: column?.type ?? "text",
                            align: column?.align ?? "left",
                            value: rawValue,
                        }}
                        isError={cell?.isError ?? false}
                        isDraft={cell?.isDraft ?? false}
                        value={rawValue}
                        {column}
                    />
                {/if}
            {/each}
            {#if context.isRemovable}
                <td class="text-center">
                    <button
                        title="Удалить"
                        class="btn m-auto text-red-500 cursor-pointer select-none"
                        onclick={() => context.handleDelete(index)}
                        ><Icon>delete</Icon></button
                    >
                </td>
            {/if}
        </tr>
    {/each}
</tbody>

<script lang="ts">
    import DataTableHeaderCell from "./DataTableHeaderCell.svelte";
    import { getDataTableState } from "./context/runes";

    const context = getDataTableState();
</script>

<thead>
    {#each context.headerRows as row, index (row)}
        <tr>
            {#each row.cells as cell}
                <DataTableHeaderCell {cell} {row} />
            {/each}
            {#if (context.isAddable || context.isRemovable) && index === 0}
                <th rowspan={row.cells.length}></th>
            {/if}
        </tr>
    {/each}
</thead>

<script lang="ts">
    import DataTable from "../../components/DataTable/DataTable.svelte";
    import { SchoolService } from "../../../bindings/protocol/backend/services";
    import type { SaveResult } from "../../components/DataTable/types";
    import type { School } from "../../../bindings/protocol/backend/models";

    const columns = $state([
        { name: "label", isEditable: true, isSortable: true, required: true },
    ]);

    async function handleSave({
        created,
        updated,
        deletedIds,
    }: {
        created: Record<string, any>[];
        updated: Record<string, any>[];
        deletedIds: number[];
    }) {
        const create = created as School[];
        const update = updated as School[];
        const deleted = deletedIds;

        const res = await SchoolService.BulkSave({
            create,
            update,
            delete: deleted,
        });
        return res as SaveResult;
    }
</script>

<div class="w-2/3 m-auto border border-stone-300 rounded-xl p-2">
    <DataTable
        headerRows={[
            {
                cells: [{ value: "Наименование", columnName: "label" }],
            },
        ]}
        fetchItems={SchoolService.GetAll}
        {columns}
        isSortable
        isPagination={false}
        sortKey="value"
        primaryKey="value"
        onSave={handleSave}
        isAddable={false}
        isRemovable={false}
        isEditable={false}
    >
        {#snippet header()}
            <h3 class="text-center text-lg! w-full p-2">
                <strong>Учебные заведения</strong>
            </h3>
        {/snippet}
    </DataTable>
</div>

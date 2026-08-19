<script lang="ts">
    import DataTable from "../../components/DataTable/DataTable.svelte";
    import { SettingsService } from "../../../bindings/protocol/backend/services";
    import type { SaveResult } from "../../components/DataTable/types";
    import type { Setting } from "../../../bindings/protocol/backend/models";

    const columns = $state([
        { name: "title", isSortable: true },
        { name: "value", isEditable: true, isSortable: true, required: true },
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
        const update = updated as Setting[];

        const res = await SettingsService.BulkSave({
            create: [],
            update,
            delete: [],
        });
        return res as SaveResult;
    }
</script>

<div class="w-2/3 m-auto border border-stone-300 rounded-xl p-2 mb-5">
    <DataTable
        headerRows={[
            {
                cells: [
                    {
                        value: "Описание",
                        columnName: "title",
                    },
                    {
                        value: "Значение",
                        columnName: "value",
                    },
                ],
            },
        ]}
        fetchItems={SettingsService.GetAllVisible}
        {columns}
        isSortable
        isPagination={false}
        sortKey="id"
        isAddable={false}
        isRemovable={false}
        isEditable={false}
        onSave={handleSave}
    >
        {#snippet header()}
            <h3 class="text-center text-lg! w-full p-2">
                <strong>Настройки</strong>
            </h3>
        {/snippet}
    </DataTable>
</div>

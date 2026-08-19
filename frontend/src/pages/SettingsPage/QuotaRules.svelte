<script lang="ts">
    import DataTable from "../../components/DataTable/DataTable.svelte";
    import { QuotaRuleService } from "../../../bindings/protocol/backend/services";
    import type { SaveResult } from "../../components/DataTable/types";
    import type { QuotaRule } from "../../../bindings/protocol/backend/models";

    const columns: {
        name: string;
        isEditable: boolean;
        type?: "number" | "text" | "date" | "select";
        required?: boolean;
    }[] = $state([
        {
            name: "minParticipants",
            isEditable: true,
            isSortable: true,
            type: "number",
        },
        {
            name: "maxParticipants",
            isEditable: true,
            isSortable: true,
            type: "number",
        },
        {
            name: "winnersQuota",
            isEditable: true,
            isSortable: true,
            type: "number",
            required: true,
        },
        {
            name: "winnersAndPrizersQuota",
            isEditable: true,
            isSortable: true,
            type: "number",
            required: true,
        },
        {
            name: "minWinnersPointsPercent",
            isEditable: false,
            isSortable: true,
            type: "number",
        },
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
        const create = created as QuotaRule[];
        const update = updated as QuotaRule[];
        const deleted = deletedIds;

        const res = await QuotaRuleService.BulkSave({
            create,
            update,
            delete: deleted,
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
                        value: "Мин. участников",
                        columnName: "minParticipants",
                    },
                    {
                        value: "Макс. участников",
                        columnName: "maxParticipants",
                    },
                    {
                        value: "Победителей",
                        columnName: "winnersQuota",
                    },
                    {
                        value: "Победителей и призеров",
                        columnName: "winnersAndPrizersQuota",
                    },
                    {
                        value: "Минимальный процент баллов у победителя",
                        columnName: "minWinnersPointsPercent",
                    },
                ],
            },
        ]}
        fetchItems={QuotaRuleService.GetAll}
        {columns}
        isSortable
        isPagination={false}
        sortKey="id"
        onSave={handleSave}
        isAddable={false}
        isRemovable={false}
        isEditable={false}
    >
        {#snippet header()}
            <h3 class="text-center text-lg! w-full p-2">
                <strong>Квоты победителей и призёров</strong>
            </h3>
        {/snippet}
    </DataTable>
</div>

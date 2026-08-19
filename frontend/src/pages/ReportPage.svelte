<script lang="ts">
    import { Events } from "@wailsio/runtime";
    import DataTable from "../components/DataTable/DataTable.svelte";
    import { onMount, setContext } from "svelte";
    import { getColumns } from "./../columns";
    import Icon from "../components/Icon.svelte";
    import type {
        Column,
        Row,
        SaveResult,
        SelectOption,
    } from "../components/DataTable/types";
    import {
        ParticipantService,
        SchoolService,
        ImportService,
        ExportService,
        SettingsService,
        QuotaRuleService,
        DraftService,
    } from "../../bindings/protocol/backend/services";
    import { headerRows } from "../sample";

    import {
        Participant,
        School,
        Setting,
        QuotaRule,
    } from "../../bindings/protocol/backend/models";

    const now = $state(new Date());
    let yearStart = $state(now.getFullYear());
    let yearEnd = $state(now.getFullYear() + 1);
    let yearStartSetting: Setting | undefined = $state();
    let yearEndSetting: Setting | undefined = $state();
    let disciplineSetting: Setting | undefined = $state();
    let discipline = $state("");
    let isLoading = $state(true);
    let errorMessage = $state("");
    let participants: Participant[] = $state([]);
    let schools: School[] = $state([]);
    let settings: Setting[] = $state([]);
    let quotaRules: QuotaRule[] = $state([]);
    let maxPoints = $state(0);
    let maxPointsSetting: Setting | undefined = $state();
    let tasksCount = $state(0);
    let tasksCountSetting: Setting | undefined = $state();
    let items: Participant[] = $state([]);
    let setTimer = $state(0);
    let importTimer = $state(0);
    // let draft: Record<string, any> | undefined = $state();

    let columns: Column[] = $state([]);

    let status:
        { type: "success" | "warning" | "error"; message: string } | undefined =
        $state();

    function handleSettingChange({
        value,
        setting,
    }: {
        value: string;
        setting: Setting | undefined;
    }) {
        if (!setting) return;
        if (setTimer) clearTimeout(setTimer);
        setTimer = setTimeout(() => {
            SettingsService.Update({ ...setting, value });
        }, 200);
    }

    function handleInput({
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
        value: string;
        index: number;
        rowsUpdate: (rows: Record<string, any>[]) => void;
        rawRow: Record<string, any>;
        rawRows: Record<string, any>[];
        groupIndex?: number;
    }) {
        if (groupIndex !== undefined) {
            rawRow[columnName][groupIndex] = {
                ...(rawRow[columnName][groupIndex] || {}),
                id: rawRow[columnName][groupIndex]?.id || 0,
                value: value,
            };
        } else {
            rawRow[columnName] = value;
        }

        if (groupIndex !== undefined) {
            handleRowRecalculate({ rawRow });
        }

        if (index >= 1)
            rowsUpdate([
                ...rawRows.slice(0, index),
                rawRow,
                ...rawRows.slice(index + 1),
            ]);
    }

    function handleRowRecalculate({
        rawRow,
    }: {
        rawRow: Record<string, any>;
    }): Record<string, any> {
        rawRow.total = rawRow.tasks.reduce(
            (acc: number, task: { value: number }) => acc + +task.value,
            0,
        );
        rawRow.percent = ((rawRow.total * 100) / maxPoints).toFixed(2);
        rawRow.rating = "";
        rawRow.status = "";
        return rawRow;
    }

    function handleRowClick({
        row,
        rawRow,
    }: {
        row: Row;
        rawRow: Record<string, any>;
    }) {
        // console.log(row, rawRow);
    }

    // function handleSort({
    //     columnName,
    //     context,
    //     parentKey,
    // }: {
    //     columnName: string;
    //     context: any;
    //     parentKey?: string;
    // }) {
    //     // console.log(columnName, context, parentKey);
    // }

    // async function handleRowAdd({
    //     rawRow,
    //     rawRows,
    //     rowsUpdate,
    //     clearNewRow,
    // }: {
    //     rawRow: Record<string, any>;
    //     rawRows: Record<string, any>[];
    //     rowsUpdate: (rows: Record<string, any>[]) => void;
    //     clearNewRow: () => void;
    // }) {
    //     rawRow.schoolId = +rawRow.schoolId;
    //     try {
    //         console.log(rawRow);
    //         const res = await ParticipantService.Add(rawRow as Participant);
    //         rowsUpdate(res);
    //         clearNewRow();
    //     } catch (err) {
    //         alert("Не удалось добавить участника");
    //         console.error(err);
    //     }
    // }

    // async function handleRowDelete({
    //     index,
    //     id,
    // }: {
    //     index: number;
    //     id: string;
    // }) {
    //     return ParticipantService.Delete(+id);
    // }

    // async function fetchParticipants() {
    //     return ParticipantService.GetAll();
    //     // .then((res) => {
    //     //     // console.log(res);
    //     //     participants = res;
    //     // })
    //     // .catch((err) => {
    //     //     alert("Не удалось загрузить данные из базы");
    //     //     console.error(err);
    //     // })
    //     // .finally(() => {
    //     //     isLoading = false;
    //     // });
    // }

    async function fetchSettings() {
        return SettingsService.GetAll()
            .then((res) => {
                settings = res;
                return res;
            })
            .catch((err) => {
                alert("Не удалось загрузить данные из базы");
                console.error(err);
            });
    }

    async function fetchQuotaRules() {
        return QuotaRuleService.GetAll()
            .then((res) => {
                quotaRules = res;
            })
            .catch((err) => {
                alert("Не удалось загрузить данные из базы");
                console.error(err);
            });
    }

    async function fetchSchools() {
        return SchoolService.GetAll()
            .then((res) => {
                schools = res;
                return res;
            })
            .catch((err) => {
                alert("Не удалось загрузить данные из базы");
                console.error(err);
            });
    }

    onMount(async () => {
        // sample().then((res) => {
        //     console.log(res);
        //     participants = res;
        //     isLoading = false;
        // });
        const settings = await fetchSettings();
        maxPointsSetting = settings?.find((s) => s.name === "max_points");
        maxPoints = +(maxPointsSetting?.value || 0);
        await fetchQuotaRules();
        const schools = await fetchSchools();
        tasksCountSetting = settings?.find((s) => s.name === "tasks_count");
        tasksCount = +(tasksCountSetting?.value || 0);
        columns = getColumns(schools as SelectOption[], tasksCount);

        yearStartSetting = settings?.find((s) => s.name === "year_start");
        if (yearStartSetting) yearStart = +yearStartSetting.value;
        yearEndSetting = settings?.find((s) => s.name === "year_end");
        if (yearEndSetting) yearEnd = +yearEndSetting.value;
        disciplineSetting = settings?.find((s) => s.name === "discipline");
        if (disciplineSetting) discipline = disciplineSetting.value;

        // await fetchDraft();

        // fetchParticipants();
    });

    async function fetchDraft() {
        return DraftService.Get("participants").then((res) => {
            return {
                ...res,
                ...(res.created && {
                    created: res.created.map((el: Record<string, any>) => ({
                        ...el,
                        schoolId: +el.schoolId,
                        tasks: el.tasks.map((t: Record<string, any>) => ({
                            ...t,
                            value: +t.value,
                        })),
                    })),
                }),
            };
        });
    }

    async function handleDraftDelete() {
        return DraftService.Delete("participants");
    }

    async function handleDraftSave({
        created,
        updated,
        deletedIds,
    }: {
        created: Record<string, any>[];
        updated: Record<string, any>;
        deletedIds: number[];
    }) {
        const res = await DraftService.Save("participants", 1, {
            created,
            updated,
            deleted: deletedIds,
        });
        return res as SaveResult;
    }

    async function handleSave({
        created,
        updated,
        deletedIds,
    }: {
        created: Record<string, any>[];
        updated: Record<string, any>[];
        deletedIds: number[];
    }) {
        const create = created as Participant[];
        const update = updated as Participant[];
        const deleted = deletedIds;

        const res = await ParticipantService.BulkSave({
            create,
            update,
            delete: deleted,
        });
        return res as SaveResult;
    }

    function updateItems(array: Participant[]) {
        items = array;
        status = {
            type: "success",
            message: "Успешно загружено",
        };
        setTimeout(() => {
            status = undefined;
        }, 3000);
    }

    async function handleImport() {
        status = undefined;
        try {
            const result = await ImportService.ImportData();
            updateItems(result);
        } catch (error) {
            if (error instanceof Error) {
                if (error.message == "Загрузка отменена") {
                    status = { type: "warning", message: error.message };
                    setTimeout(() => {
                        status = undefined;
                    }, 3000);
                } else {
                    status = { type: "error", message: error.message };
                }
            } else {
                status = { type: "error", message: "Неизвестная ошибка" };
            }
        }
    }

    async function handleExport() {
        status = undefined;
        try {
            const now = new Date();

            const year = now.getFullYear();
            const month = String(now.getMonth() + 1).padStart(2, "0");
            const day = String(now.getDate()).padStart(2, "0");
            const hours = String(now.getHours()).padStart(2, "0");
            const minutes = String(now.getMinutes()).padStart(2, "0");
            const seconds = String(now.getSeconds()).padStart(2, "0");
            const fileName = `Протокол_олимпиады__${day}-${month}-${year}_${hours}-${minutes}-${seconds}.txt`;

            await ExportService.ExportData(fileName);
            status = {
                type: "success",
                message: "Успешно скачано",
            };
            setTimeout(() => {
                status = undefined;
            }, 3000);
        } catch (error) {
            if (error instanceof Error) {
                if (error.message == "Сохранение отменено") {
                    status = { type: "warning", message: error.message };
                    setTimeout(() => {
                        status = undefined;
                    }, 3000);
                } else {
                    status = { type: "error", message: error.message };
                }
            } else {
                status = { type: "error", message: "Неизвестная ошибка" };
            }
        }
    }

    $effect(() => {
        // Подписываемся на событие Wails v3
        const unsubscribe = Events.On(
            "protocol:data-imported-event",
            (event) => {
                if (!event.data) {
                    status = {
                        type: "warning",
                        message: "Не найдено участников",
                    };
                    if (importTimer) clearTimeout(importTimer);
                    importTimer = setTimeout(() => {
                        status = undefined;
                    }, 3000);
                    return;
                }
                updateItems(event.data);
            },
        );

        // Функция очистки (вызовется автоматически при размонтировании компонента)
        return () => {
            unsubscribe();
        };
    });

    // $effect(() => {
    //     console.log($state.snapshot(participants));
    //     console.log($state.snapshot(isLoading));
    // });
</script>

<div class="flex justify-start items-center gap-3 p-3">
    <button
        class="btn-indigo btn-active"
        title="Сохранить в файл Microsoft Excel"
        onclick={handleExport}
        ><Icon>file_download</Icon> Скачать в Excel</button
    >
    <!-- <button
        class="btn-indigo"
        title="Загрузить из файла Microsoft Excel"
        onclick={handleImport}
        ><Icon>file_upload</Icon> Загрузить из Excel</button
    > -->

    <div class="min-h-8">
        {#if status}
            <span
                class="flex items-center py-1"
                class:text-green-600={status.type === "success"}
                class:text-yellow-600={status.type === "warning"}
                class:text-red-600={status.type === "error"}
                ><Icon>{status.type === "success" ? "done" : "error"}</Icon>
                {status.message}</span
            >
        {/if}
    </div>
</div>

<button
    id="upload"
    class="drop-zone"
    data-file-drop-target
    onclick={handleImport}
>
    <strong>Перетащите</strong> Excel-файл (.xlsx) сюда <br />или
    <strong>кликните</strong> для загрузки
</button>

{#if !tasksCount}
    <p class="m-auto w-20">Загрузка...</p>
{:else}
    <DataTable
        headerRows={headerRows(tasksCount)}
        {columns}
        isSortable
        isPagination={false}
        onInput={handleInput}
        onRowClick={handleRowClick}
        fetchItems={ParticipantService.GetAll}
        sortKey="rating"
        {items}
        onSave={handleSave}
        onDraftSave={handleDraftSave}
        onDraftDelete={handleDraftDelete}
        onRowRecalculate={handleRowRecalculate}
        {fetchDraft}
    >
        {#snippet header()}
            <h3 class="text-center text-lg! w-full p-2">
                <strong
                    >Протокол муниципального этапа Всероссийской Олимпиады
                    школьников в
                </strong>
                <input
                    class="bg-white itext text-lg w-14 font-bold border p-1"
                    id="year_start"
                    type="text"
                    oninput={(e: Event) =>
                        handleSettingChange({
                            value: (e.target as HTMLInputElement).value,
                            setting: yearStartSetting,
                        })}
                    bind:value={yearStart}
                />
                -
                <input
                    class="bg-white itext text-lg w-14 font-bold border p-1"
                    id="year_end"
                    bind:value={yearEnd}
                    oninput={(e: Event) =>
                        handleSettingChange({
                            value: (e.target as HTMLInputElement).value,
                            setting: yearEndSetting,
                        })}
                    type="text"
                />
                <strong>учебном году по</strong>
                <input
                    class="bg-white itext text-lg font-bold border p-1"
                    bind:value={discipline}
                    oninput={(e: Event) =>
                        handleSettingChange({
                            value: (e.target as HTMLInputElement).value,
                            setting: disciplineSetting,
                        })}
                    type="text"
                />
            </h3>
            <div>
                <span>Количество заданий: </span>
                <input
                    class="bg-white itext text-lg w-14 font-bold border p-1"
                    id="tasks_count"
                    type="text"
                    oninput={(e: Event) =>
                        handleSettingChange({
                            value: (e.target as HTMLInputElement).value,
                            setting: tasksCountSetting,
                        })}
                    bind:value={tasksCount}
                />
                <span>Максимальный балл: </span>
                <input
                    class="bg-white itext text-lg w-14 font-bold border p-1"
                    id="max_points"
                    type="text"
                    oninput={(e: Event) =>
                        handleSettingChange({
                            value: (e.target as HTMLInputElement).value,
                            setting: maxPointsSetting,
                        })}
                    bind:value={maxPoints}
                />
            </div>
        {/snippet}
    </DataTable>
{/if}

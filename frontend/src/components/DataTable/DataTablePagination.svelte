<script lang="ts">
    import { getDataTableState } from "./context/runes";
    const context = getDataTableState();

    let startIndex = 0;
    let endIndex = context.rowsPerPage;
    const rows = context.rows;

    function goToPage(page: number) {
        if (page >= 1 && page <= context.totalPages) {
            context.currentPage = page;
        }
    }
</script>

<div
    class="bg-slate-50 px-6 py-4 border-t border-slate-200 flex items-center justify-between select-none"
>
    <div class="text-sm text-slate-500">
        <span class="font-medium text-slate-800">{context.currentPage}</span>
        из
        <span class="font-medium text-slate-800">{context.totalPages}</span> страниц
    </div>

    <div class="flex items-center gap-2">
        <button
            onclick={() => goToPage(context.currentPage - 1)}
            disabled={context.currentPage === 1}
            class="material-icons inline-flex items-center justify-center w-9 h-9 text-slate-500 bg-white border border-slate-200 rounded-lg hover:bg-slate-50 disabled:opacity-50 disabled:pointer-events-none transition-colors"
        >
            keyboard_arrow_left</button
        >
        {#each Array(context.totalPages) as _, i}
            <button
                onclick={() => goToPage(i + 1)}
                class="w-9 h-9 text-sm font-medium rounded-lg transition-colors{context.currentPage ===
                i + 1
                    ? 'bg-indigo-600 text-slate-600 shadow-sm shadow-indigo-200'
                    : 'text-slate-600 bg-white border border-slate-200 hover:bg-slate-50'}"
                >{i + 1}</button
            >
        {/each}
        <button
            onclick={() => goToPage(context.currentPage + 1)}
            disabled={context.currentPage === context.totalPages}
            class="material-icons inline-flex items-center justify-center w-9 h-9 text-slate-500 bg-white border border-slate-200 rounded-lg hover:bg-slate-50 disabled:opacity-50 disabled:pointer-events-none transition-colors"
            >keyboard_arrow_right</button
        >
    </div>
</div>

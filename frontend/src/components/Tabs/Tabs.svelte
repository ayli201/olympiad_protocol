<script lang="ts">
    import type { Item } from "./types";

    export let items: Item[] = []; // Список вкладок: [{ id: 'home', label: 'Главная' }, ...]
    export let activeTabName = ""; // ID текущей активной вкладки
</script>

<div class="flex border-b border-gray-200 gap-2 mb-4">
    {#each items as tab}
        <button
            type="button"
            class="px-4 py-2 text-sm font-medium border-b-2 transition-all duration-200 cursor-pointer -mb-[2px]"
            class:border-blue-600={activeTabName === tab.id}
            class:text-blue-600={activeTabName === tab.id}
            class:border-transparent={activeTabName !== tab.id}
            class:text-gray-500={activeTabName !== tab.id}
            class:hover:text-gray-700={activeTabName !== tab.id}
            class:hover:border-gray-300={activeTabName !== tab.id}
            onclick={() => (activeTabName = tab.id)}
        >
            {tab.label}
        </button>
    {/each}
</div>

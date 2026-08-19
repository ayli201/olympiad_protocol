<script lang="ts">
    import { errorStore } from "./errorStore.svelte";
    import ReportPage from "./pages/ReportPage.svelte";
    import SettingsPage from "./pages/SettingsPage/Page.svelte";
    import "./assets/css/app.css";

    let activeTab = $state("home");
</script>

<div class="h-full">
    {errorStore.message}

    <div
        class="flex flex-col h-screen w-screen bg-white text-slate-800 select-none overflow-hidden"
    >
        <!-- Верхняя панель вкладок -->
        <nav
            class="flex items-end bg-slate-100 px-4 pt-2 border-b border-slate-200 h-12 gap-1"
        >
            <button
                class="px-4 py-2 text-sm font-medium rounded-t-lg transition-colors duration-150 border-b-2
          {activeTab === 'home'
                    ? 'bg-white text-blue-600 border-blue-600 shadow-sm'
                    : 'text-slate-500 hover:bg-slate-200/60 hover:text-slate-800 border-transparent'}"
                onclick={() => (activeTab = "home")}
            >
                <span class="mr-1.5">🏠</span> Протокол
            </button>

            <button
                class="px-4 py-2 text-sm font-medium rounded-t-lg transition-colors duration-150 border-b-2
          {activeTab === 'settings'
                    ? 'bg-white text-blue-600 border-blue-600 shadow-sm'
                    : 'text-slate-500 hover:bg-slate-200/60 hover:text-slate-800 border-transparent'}"
                onclick={() => (activeTab = "settings")}
            >
                <span class="mr-1.5">⚙️</span> Настройки
            </button>
        </nav>

        <!-- Область контента -->
        <main class="flex-1 p-6 overflow-y-auto bg-white">
            {#if activeTab === "home"}
                <ReportPage />
            {:else if activeTab === "settings"}
                <SettingsPage />
            {/if}
        </main>
    </div>
</div>

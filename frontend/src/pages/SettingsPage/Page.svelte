<script lang="ts">
    import SettingsBlock from "./Settings.svelte";
    import QuotaRulesBlock from "./QuotaRules.svelte";
    import SchoolsBlock from "./Schools.svelte";

    // import {
    //     SettingsService,
    //     SchoolService,
    //     QuotaRuleService,
    // } from "../../../bindings/protocol/backend/services";
    // import DataTable from "../../components/DataTable/DataTable.svelte";
    // import { onMount, setContext } from "svelte";
    // import {
    //     School,
    //     QuotaRule,
    //     Settings,
    // } from "../../../bindings/protocol/backend/models";
    // import { createQuery, createLazyQuery } from "../../lib/runes/index.svelte";

    // let saveStatus = $state("");

    // let isSchoolLoading = $state(true);
    // let isQuotaRuleLoading = $state(true);
    // let isSettingsLoading = $state(true);

    // let schools = $state<School[]>([]);
    // let quotaRules = $state<QuotaRule[]>([]);
    // let settings = $state<Settings[]>([]);

    // const {
    //     data: schools,
    //     isLoading: isSchoolsLoading,
    //     error: schoolsError,
    // } = createQuery({
    //     key: () => "schools",
    //     queryFn: () => SettingsService.GetAll(),
    // });

    // const {
    //     data: quotaRules,
    //     isLoading: isQuotaRuleLoading,
    //     error: quotaRulesError,
    // } = createQuery({
    //     key: () => "quotaRules",
    //     queryFn: () => QuotaRuleService.GetAll(),
    // });

    // const {
    //     data: settings,
    //     isLoading: isSettingsLoading,
    //     error: settingsError,
    // } = createQuery({
    //     key: () => "settings",
    //     queryFn: () => SettingsService.GetAll(),
    // });

    // function saveSettings() {
    //     saveStatus = "Сохранено";
    //     setTimeout(() => (saveStatus = ""), 2000);
    //     // Здесь логика отправки данных в Go
    // }

    // onMount(() => {
    //     fetchDefaults();
    // });
</script>

<div class="overflow-x-auto">
    <SettingsBlock />
    <QuotaRulesBlock />
    <SchoolsBlock />
</div>

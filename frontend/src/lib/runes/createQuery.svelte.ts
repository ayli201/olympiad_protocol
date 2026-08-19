import type { QueryOptions } from "../types/queryOptions";

export function createQuery<T>(options: QueryOptions<T>) {
  let data = $state<T | null>(null);
  let isLoading = $state(false);
  let error = $state<Error | null>(null);

  // Счетчик для ручного перезапуска данных
  let refetchIndex = $state(0);

  // Функция-условие: если не передана, то запрос всегда разрешен
  const isEnabled = options.enabled ?? (() => true);

  // Собираем все триггеры в один производный стейт:
  // Изменение key, изменение enabled или вызов refetch() запустят эффект
  const trigger = $derived.by(() => ({
    key: options.key(),
    enabled: isEnabled(),
    refetch: refetchIndex,
  }));

  $effect(() => {
    // Подписываем Svelte на изменения триггера
    const current = trigger;

    // Если условие false, сбрасываем загрузку и ничего не делаем
    if (!current.enabled) {
      isLoading = false;
      return;
    }

    let isCurrent = true;
    isLoading = true;
    error = null;

    options
      .queryFn()
      .then((result) => {
        if (isCurrent) {
          data = result;
          isLoading = false;
        }
      })
      .catch((err) => {
        if (isCurrent) {
          error = err instanceof Error ? err : new Error(String(err));
          isLoading = false;
        }
      });

    return () => {
      isCurrent = false;
    };
  });

  // Функция ручного перезапуска просто инкрементирует стейт-счетчик
  function refetch() {
    refetchIndex += 1;
  }

  return {
    get data() {
      return data;
    },
    get isLoading() {
      return isLoading;
    },
    get error() {
      return error;
    },
    refetch,
  };
}

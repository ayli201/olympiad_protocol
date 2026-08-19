import type { LazyQueryOptions } from "../types/queryOptions";

export function createLazyQuery<T, Args extends any[]>(
  options: LazyQueryOptions<T, Args>,
) {
  let data = $state<T | null>(null);
  let isLoading = $state(false);
  let error = $state<Error | null>(null);

  // Функция, которую мы будем вызывать при клике
  async function execute(...args: Args) {
    isLoading = true;
    error = null;

    try {
      const result = await options.queryFn(...args);
      data = result;
    } catch (err) {
      error = err instanceof Error ? err : new Error(String(err));
    } finally {
      isLoading = false;
    }
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
    execute, // метод для запуска запроса
  };
}

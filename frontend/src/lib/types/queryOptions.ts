export interface QueryOptions<T> {
  queryFn: () => Promise<T>;
  key: () => any;
  enabled?: () => boolean; // Функция-условие (по умолчанию всегда true)
}

export interface LazyQueryOptions<T, Args extends any[]> {
  // Функция теперь принимает аргументы (Args), которые вы передадите при вызове
  queryFn: (...args: Args) => Promise<T>;
}

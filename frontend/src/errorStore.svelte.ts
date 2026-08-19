class ErrorController {
  // Svelte 5 автоматически поймет, что это строка, благодаря инициализации
  message: string = $state("");
  private timeoutId: ReturnType<typeof setTimeout> | null = null;

  show(text: string, duration: number = 5000) {
    this.message = text;

    if (this.timeoutId) {
      clearTimeout(this.timeoutId);
    }

    this.timeoutId = setTimeout(() => {
      this.clear();
    }, duration);
  }

  clear() {
    this.message = "";
  }
}

export const errorStore = new ErrorController();

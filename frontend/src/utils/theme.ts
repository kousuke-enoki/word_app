export const setTheme = (dark: boolean) => {
  const root = document.documentElement;
  root.classList.toggle('dark', dark);
}; 
(() => {
    const root = document.documentElement;
    const toggle = document.querySelector("[data-theme-toggle]");
    if (!toggle) {
        return;
    }

    const prefersDark = window.matchMedia("(prefers-color-scheme: dark)");
    const currentTheme = () => root.dataset.theme || (prefersDark.matches ? "dark" : "light");
    const sync = () => toggle.setAttribute("aria-pressed", String(currentTheme() === "dark"));

    toggle.addEventListener("click", () => {
        const next = currentTheme() === "dark" ? "light" : "dark";
        root.dataset.theme = next;
        try {
            localStorage.setItem("theme", next);
        } catch {
            // Storage can be unavailable (private mode); the choice just won't persist.
        }
        sync();
    });
    prefersDark.addEventListener("change", sync);
    sync();
})();

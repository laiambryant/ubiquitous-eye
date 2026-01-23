const getPreferredTheme = () => {
    const stored = localStorage.getItem("theme");
    if (stored === "light" || stored === "dark") {
        return stored;
    }
    return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
};

const setTheme = (theme) => {
    document.body.dataset.theme = theme;
    localStorage.setItem("theme", theme);
};

const initThemeToggle = () => {
    const toggle = document.querySelector("[data-theme-toggle]");
    if (!toggle) {
        return;
    }
    setTheme(getPreferredTheme());
    toggle.addEventListener("click", () => {
        const nextTheme = document.body.dataset.theme === "dark" ? "light" : "dark";
        setTheme(nextTheme);
    });
};

const uniqueSorted = (values) => {
    return Array.from(new Set(values.filter(Boolean))).sort((a, b) => a.localeCompare(b));
};

const buildFilterOptions = (cards, select, extractor) => {
    const options = uniqueSorted(cards.flatMap((card) => extractor(card)));
    options.forEach((option) => {
        const element = document.createElement("option");
        element.value = option;
        element.textContent = option;
        select.appendChild(element);
    });
};

const applyFilters = (cards, state) => {
    let visibleCount = 0;
    cards.forEach((card) => {
        const name = card.querySelector("h3").textContent.toLowerCase();
        const description = card.querySelector(".project-description").textContent.toLowerCase();
        const language = card.dataset.language || "";
        const topics = (card.dataset.topics || "").split(",").filter(Boolean);
        const matchesSearch = !state.search || name.includes(state.search) || description.includes(state.search);
        const matchesLanguage = state.language === "all" || state.language === language;
        const matchesTopic = state.topic === "all" || topics.includes(state.topic);
        const isVisible = matchesSearch && matchesLanguage && matchesTopic;
        card.hidden = !isVisible;
        if (isVisible) {
            visibleCount += 1;
        }
    });
    return visibleCount;
};

const sortCards = (cards, sortBy) => {
    const sorted = [...cards];
    switch (sortBy) {
        case "name-asc":
            sorted.sort((a, b) => a.querySelector("h3").textContent.localeCompare(b.querySelector("h3").textContent));
            break;
        default:
            sorted.sort((a, b) => new Date(b.dataset.updated) - new Date(a.dataset.updated));
            break;
    }
    return sorted;
};

const initProjects = () => {
    const grid = document.querySelector("[data-project-grid]");
    if (!grid) {
        return;
    }
    const cards = Array.from(grid.querySelectorAll("[data-project-card]"));
    const search = document.querySelector("[data-project-search]");
    const languageFilter = document.querySelector("[data-filter-language]");
    const topicFilter = document.querySelector("[data-filter-topic]");
    const sort = document.querySelector("[data-project-sort]");
    const emptyState = document.querySelector("[data-empty-state]");
    const count = document.querySelector("[data-project-count]");
    const state = {
        search: "",
        language: "all",
        topic: "all",
        sort: "updated-desc",
    };

    buildFilterOptions(cards, languageFilter, (card) => [card.dataset.language]);
    buildFilterOptions(cards, topicFilter, (card) => (card.dataset.topics || "").split(","));

    const refresh = () => {
        const sorted = sortCards(cards, state.sort);
        sorted.forEach((card) => grid.appendChild(card));
        const visibleCount = applyFilters(sorted, state);
        if (count) {
            count.textContent = visibleCount;
        }
        if (emptyState) {
            emptyState.hidden = visibleCount !== 0;
        }
    };

    search.addEventListener("input", (event) => {
        state.search = event.target.value.trim().toLowerCase();
        refresh();
    });

    languageFilter.addEventListener("change", (event) => {
        state.language = event.target.value;
        refresh();
    });

    topicFilter.addEventListener("change", (event) => {
        state.topic = event.target.value;
        refresh();
    });

    sort.addEventListener("change", (event) => {
        state.sort = event.target.value;
        refresh();
    });

    refresh();
};

const initReveal = () => {
    const reveals = document.querySelectorAll(".reveal");
    if (!reveals.length) {
        return;
    }
    const observer = new IntersectionObserver(
        (entries, obs) => {
            entries.forEach((entry) => {
                if (entry.isIntersecting) {
                    entry.target.classList.add("is-visible");
                    obs.unobserve(entry.target);
                }
            });
        },
        { threshold: 0.15 }
    );
    reveals.forEach((element, index) => {
        element.style.transitionDelay = `${index * 60}ms`;
        observer.observe(element);
    });
};

document.addEventListener("DOMContentLoaded", () => {
    initThemeToggle();
    initProjects();
    initReveal();
});

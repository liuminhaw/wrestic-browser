window.addEventListener("DOMContentLoaded", () => {
    const selection = document.querySelector("#type");
    const repoConfigs = document.querySelectorAll(".repo-config");
    const requiredOptions = document.querySelectorAll(".required-option");

    selection.addEventListener("change", (e) => {
        repoConfigs.forEach((config) => {
            config.classList.add("hidden");
        });
        requiredOptions.forEach((option) => {
            option.removeAttribute("required");
        });

        let selectedConfig = document.querySelector(
            `#${e.target.value}-config`,
        );
        let selectedOptions = document.querySelectorAll(
            `#${e.target.value}-config .required-option`,
        );
        selectedConfig.classList.remove("hidden");
        selectedOptions.forEach((option) => {
            option.setAttribute("required", "");
        });
    });
});

window.addEventListener("DOMContentLoaded", () => {
  const selection = document.querySelector("#type");
  const repoConfigs = document.querySelectorAll(".repo-config");

  selection.addEventListener("change", (e) => {
    repoConfigs.forEach((config) => {
      config.classList.add("hidden");
    });
    let selectedConfig = document.querySelector(`#${e.target.value}-config`);
    selectedConfig.classList.remove("hidden");
  });
});

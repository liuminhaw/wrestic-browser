window.addEventListener("DOMContentLoaded", () => {
  const selection = document.querySelector("#type");
  const s3Config = document.querySelector("#s3-config");
  const sftpConfig = document.querySelector("#sftp-config");

  selection.addEventListener("change", (e) => {
    switch (e.target.value) {
      case "s3":
        s3Config.classList.remove("hidden");
        sftpConfig.classList.add("hidden");
        break;
      case "sftp":
        s3Config.classList.add("hidden");
        sftpConfig.classList.remove("hidden");
        break;
      default:
        s3Config.classList.add("hidden");
        sftpConfig.classList.add("hidden");
    }
  });
});

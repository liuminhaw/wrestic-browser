window.addEventListener("DOMContentLoaded", () => {
    document.addEventListener("click", (e) => {
        // let isClickOnMenuToggle = false;

        const toggleDots = document.querySelectorAll(".toggle-dots");
        toggleDots.forEach((dot) => {
            const menu = dot.nextElementSibling;
            if (dot.contains(e.target) || menu.contains(e.target)) {
                return;
            }
            if (!menu.classList.contains("hidden")) {
                toggleEditForm(menu);
            }
        });

    });

    const toggleDots = document.querySelectorAll(".toggle-dots");

    toggleDots.forEach((toggleDot) => {
        toggleDot.addEventListener("click", () => {
            // Get the next sibling of the toggle dot
            const nextSibling = toggleDot.nextElementSibling;
            // Toggle the class of the next sibling
            toggleEditForm(nextSibling);
        });
    });
});

function toggleEditForm(obj) {
    if (obj.classList.contains("hidden")) {
        obj.classList.remove("hidden");
        window.setTimeout(function () {
            obj.classList.add("opacity-100", "scale-100", "duration-75");
            obj.classList.remove("opacity-0", "scale-95", "duration-100");
        }, 0);
    } else {
        obj.classList.add("opacity-0", "scale-95", "duration-100");
        obj.classList.remove("opacity-100", "scale-100", "duration-75");
        window.setTimeout(function () {
            obj.classList.add("hidden");
        }, 75);
    }
}

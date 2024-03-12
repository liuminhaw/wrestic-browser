const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
    content: ["../templates/**/*.{gohtml,html}"],
    theme: {
        extend: {
            fontFamily: {
                sans: ['Inter var', ...defaultTheme.fontFamily.sans],
                "satisfy": ["Satisfy"],
            },
        },
    },
    variants: {},
    plugins: [],
};

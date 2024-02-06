// tailwind.config.js
const { nextui } = require('@nextui-org/react');

module.exports = {
    content: ['./node_modules/@nextui-org/theme/dist/**/*.{js,ts,jsx,tsx}'],
    theme: {
        extend: {}
    },
    darkMode: 'class',
    plugins: [
        nextui({
            prefix: 'nextui', // prefix for themes variables
            addCommonColors: false // override common colors (e.g. "blue", "green", "pink").
        })
    ]
};

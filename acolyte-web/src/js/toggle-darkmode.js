const lightTheme = {
    backgroundColor: 'white',
    textColorBold: 'black',
    textColorSubtle: 'darkslategrey',
    linkColor: 'blueviolet',
    accentColor: 'lightgrey',
    shadowColor: 'rgba(0,0,0,0.07)',
    shadowColorIntense: 'rgba(0,0,0,.28)',
};

const darkTheme = {
    backgroundColor: '#021615',
    textColorBold: 'white',
    textColorSubtle: '#e0e0e0',
    linkColor: '#ebf745',
    accentColor: '#092d2c',
    shadowColor: 'rgba(235, 247, 69, 0.07)',
    shadowColorIntense: 'rgba(235, 247, 69, 0.28)',
};

let currentTheme = lightTheme;

function setToTheme(theme) {
    currentTheme = theme;

    document.documentElement.style.setProperty('--background-color', theme.backgroundColor);
    document.documentElement.style.setProperty('--text-color-bold', theme.textColorBold);
    document.documentElement.style.setProperty('--text-color-subtle', theme.textColorSubtle);
    document.documentElement.style.setProperty('--link-color', theme.linkColor);
    document.documentElement.style.setProperty('--accent-color', theme.accentColor);
    document.documentElement.style.setProperty('--shadow-color', theme.shadowColor);
    document.documentElement.style.setProperty('--shadow-color-intense', theme.shadowColorIntense);
}

function setThemeToStorage() {
    let storageTheme = localStorage.getItem("theme");
    if (storageTheme === "dark") {
        setToTheme(darkTheme);
    } else if (storageTheme === "light") {
        setToTheme(lightTheme);
    }
    // if no previously set theme, do nothing
}

setThemeToStorage();

export function toggleDarkMode() {
    if (currentTheme === darkTheme) {
        setToTheme(lightTheme);
        localStorage.setItem("theme", "light");
    } else {
        setToTheme(darkTheme);
        localStorage.setItem("theme", "dark");
    }
}
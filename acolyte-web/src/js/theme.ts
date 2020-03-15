interface Theme {
    backgroundColor: string,
    cardColor: string,
    textColorBold: string,
    textColor: string,
    textColorSubtle: string,
    linkColor: string,
    accentColor: string,
    shadowColor: string,
    shadowColorIntense: string,
    brandBackground: string,
}

const lightTheme = {
    backgroundColor: '#F8F8F8',
    cardColor: 'white',
    textColorBold: 'black',
    textColor: '#0f0f0f',
    textColorSubtle: 'darkslategrey',
    linkColor: 'blueviolet',
    accentColor: 'lightgrey',
    shadowColor: 'rgba(111,111,111,.2)',
    shadowColorIntense: 'rgba(111,111,111,.4)',
    brandBackground: 'linear-gradient(-45deg, #bcf6b3, #b6d2ff)',
};

const darkTheme = {
    backgroundColor: '#021615',
    cardColor: 'black',
    textColorBold: 'white',
    textColor: '#f0f0f0',
    textColorSubtle: '#e0e0e0',
    linkColor: '#35deea',
    accentColor: '#092d2c',
    shadowColor: 'rgba(0,0,0, 0.6)',
    shadowColorIntense: 'rgba(0,0,0, 0.9)',
    brandBackground: 'linear-gradient(-45deg, #0f3538,#240f38)',
};

let currentTheme = lightTheme;

export function setTheme(theme: Theme) {
    let documentElementStyle = document.documentElement.style;

    documentElementStyle.setProperty('--background-color', theme.backgroundColor);

    documentElementStyle.setProperty('--card-color', theme.cardColor);
    documentElementStyle.setProperty('--text-color-bold', theme.textColorBold);
    documentElementStyle.setProperty('--text-color', theme.textColor);
    documentElementStyle.setProperty('--text-color-subtle', theme.textColorSubtle);
    documentElementStyle.setProperty('--link-color', theme.linkColor);
    documentElementStyle.setProperty('--accent-color', theme.accentColor);
    documentElementStyle.setProperty('--shadow-color', theme.shadowColor);
    documentElementStyle.setProperty('--shadow-color-intense', theme.shadowColorIntense);
    documentElementStyle.setProperty('--brand-background', theme.brandBackground);

    currentTheme = theme;
    localStorage.setItem("theme", JSON.stringify(currentTheme));
}

function setThemeFromStorage() {
    try {
        let storageTheme: Theme = JSON.parse(localStorage.getItem("theme"));
        if (storageTheme != null) {
            setTheme(storageTheme);
        }
    } catch (e) {
        let hasDarkColorScheme = window.matchMedia('(prefers-color-scheme: dark)').matches;
        let storageTheme: Theme = hasDarkColorScheme ? darkTheme : lightTheme;
        setTheme(storageTheme);
    }
}

setThemeFromStorage();

export function toggleDarkMode() {
    if (JSON.stringify(currentTheme) === JSON.stringify(darkTheme)) {
        setTheme(lightTheme);
    } else {
        setTheme(darkTheme);
    }
}

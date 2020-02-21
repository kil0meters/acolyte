interface Theme {
    backgroundColor: string,
    cardColor: string,
    textColorBold: string,
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
    textColorSubtle: 'darkslategrey',
    linkColor: 'blueviolet',
    accentColor: 'lightgrey',
    shadowColor: 'rgba(111,111,111,.2)',
    shadowColorIntense: 'rgba(111,111,111,.4)',
    brandBackground: 'linear-gradient(-45deg, #eaf9ef, #f4eaf9)',
};

const darkTheme = {
    backgroundColor: '#021615',
    cardColor: 'black',
    textColorBold: 'white',
    textColorSubtle: '#e0e0e0',
    linkColor: '#35deea',
    accentColor: '#092d2c',
    shadowColor: 'rgba(0,0,0, 0.6)',
    shadowColorIntense: 'rgba(0,0,0, 0.9)',
    brandBackground: 'linear-gradient(-45deg, #0f3538,#240f38)',
};

let currentTheme = lightTheme;

export function setTheme(theme: Theme) {
    currentTheme = theme;

    document.documentElement.style.setProperty('--background-color', theme.backgroundColor);
    document.documentElement.style.setProperty('--card-color', theme.cardColor);
    document.documentElement.style.setProperty('--text-color-bold', theme.textColorBold);
    document.documentElement.style.setProperty('--text-color-subtle', theme.textColorSubtle);
    document.documentElement.style.setProperty('--link-color', theme.linkColor);
    document.documentElement.style.setProperty('--accent-color', theme.accentColor);
    document.documentElement.style.setProperty('--shadow-color', theme.shadowColor);
    document.documentElement.style.setProperty('--shadow-color-intense', theme.shadowColorIntense);
    document.documentElement.style.setProperty('--brand-background', theme.brandBackground);


    let chatFrame = <HTMLIFrameElement>document.getElementById("chat");

    if (chatFrame) {
        let chatDocument = chatFrame.contentDocument;

        chatDocument.documentElement.style.setProperty('--background-color', theme.backgroundColor);
        chatDocument.documentElement.style.setProperty('--card-color', theme.cardColor);
        chatDocument.documentElement.style.setProperty('--text-color-bold', theme.textColorBold);
        chatDocument.documentElement.style.setProperty('--text-color-subtle', theme.textColorSubtle);
        chatDocument.documentElement.style.setProperty('--link-color', theme.linkColor);
        chatDocument.documentElement.style.setProperty('--accent-color', theme.accentColor);
        chatDocument.documentElement.style.setProperty('--shadow-color', theme.shadowColor);
        chatDocument.documentElement.style.setProperty('--shadow-color-intense', theme.shadowColorIntense);
        chatDocument.documentElement.style.setProperty('--brand-background', theme.brandBackground);
    }

    localStorage.setItem("theme", JSON.stringify(theme));
}

function setThemeFromStorage() {
    try {
        let storageTheme: Theme = JSON.parse(localStorage.getItem("theme"));
        if (storageTheme != null) {
            setTheme(storageTheme);
        }
    } catch (e) {
        localStorage.setItem("theme", JSON.stringify(darkTheme));
    }
}

setThemeFromStorage();

export function toggleDarkMode() {
    if (currentTheme === darkTheme) {
        setTheme(lightTheme);
    } else {
        setTheme(darkTheme);
    }
}
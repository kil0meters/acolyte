import FeelsBadMan from "../emotes/FeelsBadMan.png"
import trueEmote from "../emotes/TRUE.png"
import F from "../emotes/F.png"

const emotes = {
    "TRUE": trueEmote,
    "FeelsBadMan": FeelsBadMan,
    "F": F,
}

function regexEscape(str) {
    return str.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&')
}

export function getEmotes() {
    return Object.keys(emotes)
}

export function replaceTextWithEmotes(text) {
    Object.keys(emotes).forEach((emoteName) => {
        let regex = new RegExp(`(${regexEscape(emoteName)})(\\s|$)`, 'g')

        text = text.replace(regex, `<img class="emote" alt="$1" src="/static/${emotes[emoteName]}">$2`)
    })
    return text
}

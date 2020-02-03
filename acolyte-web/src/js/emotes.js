import AwwMan from "../emotes/AwwMan.png"
import FeelsBadMan from "../emotes/FeelsBadMan.png"
import trueEmote from "../emotes/TRUE.png"
import F from "../emotes/F.png"
import Fukushima from "../emotes/Fukushima.png"
import Tetrahydrocannabinol from "../emotes/Tetrahydrocannabinol.png"
import Australia from "../emotes/Australia.png"

const emotes = {
    "AwwMan": AwwMan,
    "Australia": Australia,
    "TRUE": trueEmote,
    "Tetrahydrocannabinol": Tetrahydrocannabinol,
    "FeelsBadMan": FeelsBadMan,
    "Fukushima": Fukushima,
    "F": F,
}

function regexEscape(str) {
    return str.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&')
}

export function getEmotes() {
    return Object.keys(emotes)
}

export function replaceTextWithEmotes(text, onclick) {
    Object.keys(emotes).forEach((emoteName) => {
        let regex = new RegExp(`(${regexEscape(emoteName)})(\\s|$)`, 'g')

        text = text.replace(regex, `<img class="emote" alt="$1" onclick="${onclick}" src="/static/${emotes[emoteName]}">$2`)
    })
    return text
}
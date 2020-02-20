import AwwMan from "../emotes/AwwMan.png";
import FeelsBadMan from "../emotes/FeelsBadMan.png";
import trueEmote from "../emotes/TRUE.png";
import F from "../emotes/F.png";
import Fukushima from "../emotes/Fukushima.png";
import Tetrahydrocannabinol from "../emotes/Tetrahydrocannabinol.png";
import Australia from "../emotes/Australia.png";
import WOW from "../emotes/WOW.png";

import {AutocompletionSuggestion} from "./autocompletion";

const emotes: any = {
    "AwwMan": AwwMan,
    "Australia": Australia,
    "TRUE": trueEmote,
    "Tetrahydrocannabinol": Tetrahydrocannabinol,
    "FeelsBadMan": FeelsBadMan,
    "Fukushima": Fukushima,
    "F": F,
    "WOW": WOW,
};

function regexEscape(str: string): string {
    return str.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&');
}

export function getEmotes(): AutocompletionSuggestion[] {
    return Object.keys(emotes).map((emote) => {
        return {
            name: emote,
            description: replaceTextWithEmotes(emote),
        };
    })
}

export function replaceTextWithEmotes(text: string, onclick?: string): string {
    for (let emoteName of Object.keys(emotes)) {

        let regex = new RegExp(`(${regexEscape(emoteName)})(\\s|$|[,.<>?/!'"])`, 'g');

        text = text.replace(regex, `<img class="emote" alt="$1" title="$1" onclick="${onclick}" src="/static/${emotes[emoteName]}">$2`);
    }
    return text;
}

// this is done in a rather inefficient manner 
export function renderEmotesInElement(element: HTMLElement) {
    element.innerHTML = replaceTextWithEmotes(element.innerHTML);
}
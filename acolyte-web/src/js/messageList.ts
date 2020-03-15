import {getEmotes, renderEmotesInElement, replaceTextWithEmotes} from "./emotes";
import {AutocompletionSuggestion} from "./autocompletion";

export interface Message {
    text: any;
    username: string;
    id?: string;
}

const linkRegex: RegExp = /https?:\/\/(www\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_+.~#?&//=]*)/g;

export function renderLinksInElement(element: HTMLElement) {
    element.innerHTML = element.innerHTML.replace(linkRegex, (url: string) => {
        getLinkPreview(url)
            .then((linkElement: HTMLElement) => {
                if (linkElement != null) {
                    element.appendChild(linkElement);
                }
            });

        let displayURL = (url.length > 50) ? url.substring(1, 50) + "..." : url;

        return `<a class="chat-link" target="_blank" href="${url}">${displayURL}</a>`
    });
}

async function getLinkPreview(linkURL: string): Promise<HTMLElement> {
    let response = await fetch('/api/v1/link-data?link=' + encodeURI(linkURL), {
        credentials: 'include',
    });
    let data = await response.json();

    if (data["title"].length === 0) {
        return
    }

    let linkElement = document.createElement('link-preview');
    linkElement.setAttribute('link', data["link"]);
    linkElement.setAttribute('published-date', data["published_date"]);
    linkElement.setAttribute('title', data["title"]);
    linkElement.setAttribute('content', data["content"]);

    return linkElement;
}

export class MessageList {
    public list: Message[] = [];
    currentCombo: Message[] = [];

    messageListElement: HTMLDivElement;
    emotes: AutocompletionSuggestion[] = getEmotes();
    maxHeight: number;

    username: string;
    moderatorPerms: boolean;

    constructor(messageListID: string, maxHeight: number, username: string, moderatorPerms: boolean) {
        this.messageListElement = <HTMLDivElement>document.getElementById(messageListID);
        this.emotes = getEmotes();
        this.maxHeight = maxHeight;

        this.username = username;
        this.moderatorPerms = moderatorPerms;

        this.messageListElement.addEventListener("scroll", function () {
            // checks if it's scrolled to the ototm
            if ((this.scrollTop + this.clientHeight) >= this.scrollHeight) {
                console.log("false");
                document.getElementById('messages-below').classList.add('hidden');
            } else {
                console.log("true");
                document.getElementById('messages-below').classList.remove('hidden');
            }
        })

        document.getElementById("messages-below").addEventListener("", () => {
            this.scrollToBottom(0);
        })
    }

    buildMessage(message: Message): HTMLElement {
        let messageElement = document.createElement("div");
        messageElement.id = message.id;
        messageElement.classList.add("chat-message");

        let usernameElement = document.createElement("a");
        usernameElement.href = '#' + message.username;
        usernameElement.textContent = message.username;
        usernameElement.classList.add("username");

        let textElement = document.createElement("span");
        textElement.innerHTML = message.text;

        renderLinksInElement(textElement);
        renderEmotesInElement(textElement);

        if (message.text.includes(this.username) && message.username !== this.username) {
            messageElement.classList.add("mentioned")
        } else if (message.username === this.username) {
            messageElement.classList.add("self");
            // if sent by self, scroll to bottom regardless of current scroll position
            this.scrollToBottom(document.body.offsetHeight);
        }

        textElement.classList.add("message-text");

        messageElement.appendChild(usernameElement);

        if (this.moderatorPerms && message.id !== "00000000-0000-0000-0000-000000000000") {
            let removeMessageElement = document.createElement("button");
            removeMessageElement.textContent = "Remove";
            removeMessageElement.classList.add("remove-message-button");
            removeMessageElement.onclick = function () {
                // @ts-ignore
                // the variable `liveChat` is defined in the HTML document
                liveChat.sendMessage(`/remove ${message.id}`)
            };

            messageElement.appendChild(removeMessageElement)
        }

        messageElement.appendChild(textElement);

        return messageElement;
    }


    scrollToBottom(offsetHeight: number) {
        if ((this.messageListElement.clientHeight + this.messageListElement.scrollTop + 64) >= (this.messageListElement.scrollHeight - offsetHeight)) {
            this.messageListElement.scrollTo(0, this.messageListElement.scrollHeight);
        }
    }


    push(message: Message) {
        if (message.username === "delete-message") {
            this.removeByID(message.text);
        } else {
            this.list.push(message);

            let messageElement = this.buildMessage(message);
            this.messageListElement.appendChild(messageElement);

            this.checkForCombos();
            this.replaceComboListWithElement();

            this.scrollToBottom(messageElement.offsetHeight);
            while ((this.messageListElement.offsetHeight + 24) > this.maxHeight) { // offset for padding
                this.removeByIndex(0);
            }
        }
    }

    checkForCombos() {
        this.currentCombo = [];

        let recentMessage = this.list[this.list.length - 1];
        let prevMessage = recentMessage;

        if (!this.emotes.map((name) => name.name).includes(recentMessage.text)) {
            return;
        }

        for (let i = 1; recentMessage.text === prevMessage.text && i <= this.list.length; i++) {
            this.currentCombo.push(this.list[this.list.length - i]);
            prevMessage = this.list[this.list.length - i];
        }

        if (this.currentCombo[this.currentCombo.length - 1].text !== recentMessage.text) {
            this.currentCombo.pop(); // combo broken
        }
    }

    replaceComboListWithElement() {
        if (this.currentCombo.length > 1) {
            for (let message of this.currentCombo) {
                let elementToRemove = document.getElementById(message.id);

                if (typeof (elementToRemove) !== undefined && elementToRemove !== null) {
                    this.messageListElement.removeChild(elementToRemove);
                }
            }

            let currentComboElement = <HTMLDivElement>document.createElement('div');

            let mostRecentMessage = <HTMLDivElement>this.messageListElement.lastElementChild;

            if (mostRecentMessage !== null) {
                if (Array.from(mostRecentMessage.classList).includes('combo') && mostRecentMessage.firstElementChild.getAttribute('alt') === this.currentCombo[0].text) {
                    currentComboElement = mostRecentMessage;
                }
            }

            currentComboElement.classList.add('chat-message', 'combo');
            currentComboElement.innerHTML = `${replaceTextWithEmotes(this.currentCombo[0].text)} ${this.currentCombo.length}x ~combo~`;

            this.messageListElement.appendChild(currentComboElement);
        }
    }

    removeByIndex(index: number) {
        this.list.splice(index, 1);
        this.messageListElement.children[0].remove()
    }

    removeByID(id: string) {
        let i;
        for (i = 0; i < this.list.length; i++) {
            if (this.list[i].id === id) {
                break;
            }
        }
        if (i !== this.list.length) this.list.splice(i, 1);

        this.messageListElement.removeChild(document.getElementById(id))
    }
}

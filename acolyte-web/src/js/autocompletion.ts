import {getEmotes} from "./emotes";

export interface AutocompletionSuggestion {
    name: string,
    description: string,
}

export class Autocompletion {
    entry: HTMLInputElement;
    popup: HTMLDivElement;

    chatCommands: AutocompletionSuggestion[] = [];
    users: AutocompletionSuggestion[] = [];
    suggestions: AutocompletionSuggestion[] = [];
    emotes: AutocompletionSuggestion[] = [];

    tabIndex: number = 1;

    previousMessages: string[] = [];
    messageIndex: number = 0;
    currentValue: string = "";

    constructor() {
        this.entry = <HTMLInputElement>document.getElementById('entry-body');
        this.popup = <HTMLDivElement>document.getElementById('autocompletion-popup');

        this.emotes = getEmotes();

        this.registerEventListeners();
    }

    setCommands(commands: AutocompletionSuggestion[]) {
        this.chatCommands = commands;
    }

    setUsers(users: string[]) {
        this.users = users.map((username) => {
            return {
                name: username,
                description: ""
            }
        })
    }

    setPopupToSuggestions() {
        if (this.suggestions.length === 0) {
            this.popup.classList.add('hidden');
        } else {
            while (this.popup.firstChild) {
                this.popup.removeChild(this.popup.firstChild);
            }

            this.popup.classList.remove('hidden');

            for (let suggestion of this.suggestions) {
                let suggestionElement = document.createElement('p');

                if (this.emotes.includes(suggestion)) {
                    suggestionElement.innerHTML = suggestion.description + suggestion.name;
                } else {
                    suggestionElement.textContent = suggestion.name + ' ' + suggestion.description;
                }

                this.popup.appendChild(suggestionElement);
            }
        }
    }

    setHighlightedSuggestion(index: number) {
        let children = this.popup.children;
        for (let i = 0; i < children.length; i++) {
            children[i].classList.remove('highlighted');
        }

        this.popup.children[index].classList.add('highlighted');
    }

    sentMessage(message: string) {
        if (this.previousMessages[this.previousMessages.length - 1] !== message) {
            this.previousMessages.push(message);
        }
        this.messageIndex = 0;
    }

    registerEventListeners() {
        this.entry.addEventListener("keydown", (event) => {
            if (event.key === "Tab") {
                event.preventDefault();

                if (!event.shiftKey) {
                    this.tabIndex = Math.min(this.suggestions.length, this.tabIndex + 1);
                } else {
                    this.tabIndex = Math.max(1, this.tabIndex - 1);
                }

                if (this.suggestions.length !== 0) {
                    let words = this.entry.value.split(' ');

                    words[words.length - 1] = this.suggestions[this.suggestions.length - this.tabIndex].name;

                    this.entry.value = words.join(' ');
                    this.setHighlightedSuggestion(this.suggestions.length - this.tabIndex);
                }
            } else if (event.key === "ArrowDown") {
                event.preventDefault();
                this.messageIndex = Math.max(this.messageIndex - 1, 0);
                if (this.messageIndex === 0) {
                    this.entry.value = this.currentValue;
                } else {
                    this.entry.value = this.previousMessages[this.previousMessages.length - this.messageIndex];
                }
            } else if (event.key === "ArrowUp") {
                event.preventDefault();
                this.messageIndex = Math.min(this.messageIndex + 1, this.previousMessages.length);
                this.entry.value = this.previousMessages[this.previousMessages.length - this.messageIndex];
            }
        });

        this.entry.addEventListener("keyup", (event) => {
            if (event.key !== "ArrowUp" &&
                event.key !== "ArrowDown" &&
                event.key !== "Tab" &&
                // event.key !== "Enter" &&
                event.key !== "Shift" &&
                event.key !== "Control" &&
                event.key !== "Meta" &&
                event.key !== "Alt") {

                let text = this.entry.value;

                if (this.messageIndex === 0) {
                    this.currentValue = text
                }

                this.tabIndex = 0;
                this.suggestions = [];
                if (text === "") {
                    this.messageIndex = 0;
                    this.setPopupToSuggestions();
                } else {
                    let completionText = text.toLowerCase();

                    for (let suggestion of this.users) {
                        let words = completionText.split(' ');
                        let word = words[words.length - 1];
                        if (word !== '') {
                            if (suggestion.name.toLowerCase().startsWith(words[words.length - 1])) {
                                this.suggestions.push(suggestion);
                            }
                        }
                    }

                    for (let suggestion of this.chatCommands) {
                        if (suggestion.name.toLowerCase().startsWith(completionText)
                            || (completionText.trim().toLowerCase().startsWith(suggestion.name.toLowerCase()))) {
                            this.suggestions.push(suggestion);
                        }
                    }

                    for (let suggestion of this.emotes) {
                        let words = completionText.split(' ');
                        let word = words[words.length - 1];
                        if (word !== '') {
                            if (suggestion.name.toLowerCase().startsWith(words[words.length - 1])) {
                                this.suggestions.push(suggestion);
                            }
                        }
                    }

                    this.setPopupToSuggestions();
                }
            }
        })
    }
}
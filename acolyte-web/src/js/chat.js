// import katex from 'katex'
import {getEmotes, replaceTextWithEmotes} from './emotes.js';
import {MessageList} from "./messageList.js";
import {UserList} from "./userList.js";

export class MBChat {
    constructor(maxHeight, noEntry, username, moderatorPerms) {
        if (location.protocol === 'https:') {
            this.wsProtocol = 'wss:';
        } else {
            this.wsProtocol = 'ws:';
        }

        this.noEntry = noEntry === true;
        this.moderatorPerms = moderatorPerms === true;

        this.maxHeight = maxHeight;
        this.conn = new WebSocket(`${this.wsProtocol}//${window.location.host}/api/v1/chat`);
        this.username = username;
        this.isUnauthorized = false;

        this.timeoutInterval = null;

        this.userList = new UserList();
        this.messageList = new MessageList('message-list', this.maxHeight, this.username, this.moderatorPerms);

        if (this.noEntry === false) {
            this.entryBody = document.getElementById('entry-body');
            let emotePopup = document.getElementById('emote-popup');

            getEmotes().forEach((emote) => {
                emotePopup.innerHTML += replaceTextWithEmotes(emote.name, `document.getElementById('entry-body').value += '${emote.name} '`);
            });

            this.initializeEntryBody();
            this.autocompletionHelper = new Autocompletion();
        }
    }

    initializeConnection() {
        this.conn.addEventListener("message", (m) => {
            let message = JSON.parse(m.data);
            if (message.text === "UNAUTHORIZED") {
                this.isUnauthorized = true;
            } else if (message.username === "user-add") {
                console.log("user joined: " + message.text);
                this.userList.add(message.text);
                this.autocompletionHelper.setUsers(this.userList.list);
            } else if (message.username === "user-remove") {
                console.log("user left: " + message.text);
                this.userList.remove(message.text);
                this.autocompletionHelper.setUsers(this.userList.list)
            } else if (message.username === "command-list" && this.noEntry === false) {
                this.autocompletionHelper.setCommands(message.text);
            } else if (message.username === "user-list") {
                message.text.forEach((username) => {
                    this.userList.add(username);
                });
                this.autocompletionHelper.setUsers(this.userList.userList);
            } else {
                this.messageList.push(message);
            }
        });

        this.conn.addEventListener("close", (_) => {
            this.messageList.push({"username": "Client", "text": "Disconnected. Trying to reconnect in 5 seconds..."});
            this.timeoutInterval = setInterval(() => {
                console.log("Trying to reconnect");
                this.messageList.push({
                    "username": "Client",
                    "text": "Disconnected. Trying to reconnect in 5 seconds..."
                });

                this.conn = new WebSocket(`${this.wsProtocol}//${window.location.host}/api/v1/chat`);

                setTimeout(() => {
                    if (this.conn.readyState === this.conn.OPEN) {
                        this.messageList.push({"username": "Client", "text": "Reconnected."});
                        this.initializeConnection(this.conn);
                        console.log("Reconnected");
                        clearInterval(this.timeoutInterval);
                    }
                }, 100);
            }, 5000);
        })
    }

    sendMessage(message) {
        this.conn.send(JSON.stringify({
            "username": this.username,
            "text": message,
        }));
    }

    initializeEntryBody() {
        document.getElementById('entry-body').addEventListener("keydown", (event) => {
            if (event.key === "Enter" && !event.shiftKey) {
                event.preventDefault();

                if (this.isUnauthorized) {
                    toggleLoginPrompt();
                } else {
                    this.autocompletionHelper.sentMessage(document.getElementById('entry-body').value);

                    this.sendMessage(document.getElementById('entry-body').value);

                    document.getElementById('entry-body').value = "";
                }
            }
        })
    }
}

class Autocompletion {
    constructor(entry) {
        this.entry = entry;
        this.entry = document.getElementById('entry-body');
        this.popup = document.getElementById('autocompletion-popup');

        this.chatCommands = [];
        this.emotes = getEmotes();
        this.users = [];

        this.suggestions = [];
        this.tabIndex = 1;

        this.previousMessages = [];
        this.messageIndex = 0;

        this.currentValue = "";

        this.registerEventListeners();
    }

    setCommands(commands) {
        this.chatCommands = commands;
    }

    setUsers(users) {
        console.log(users);

        this.users = users.map((username) => {
            return {
                "name": username,
                "description": "",
            }
        });
        console.log(this.users);
    }

    setPopupToSuggestions() {
        if (this.suggestions === []) {
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

    setHighlightedSuggestion(index) {
        for (let suggestionElement of this.popup.children) {
            suggestionElement.classList.remove('highlighted');
        }

        this.popup.children[index].classList.add('highlighted');
    }

    sentMessage(message) {
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

let loginPromptShown = false;
let settingsShown = false;
let userListShown = false;
let emotePopupShown = false;

export function toggleLoginPrompt() {
    if (loginPromptShown) {
        document.getElementById("login-prompt").classList.add('hidden');
    } else {
        document.getElementById("login-prompt").classList.remove('hidden');
    }
    loginPromptShown = !loginPromptShown;
}

export function toggleUserList() {
    if (userListShown) {
        document.getElementById("user-list-overlay").classList.add('hidden');
    } else {
        document.getElementById("user-list-overlay").classList.remove('hidden');
    }
    userListShown = !userListShown;
}

export function toggleSettings() {
    if (settingsShown) {
        document.getElementById("settings-overlay").classList.add('hidden');
    } else {
        document.getElementById("settings-overlay").classList.remove('hidden');
    }
    settingsShown = !settingsShown;
}

export function toggleEmotePopup() {
    if (emotePopupShown) {
        document.getElementById("emote-popup").classList.add("hidden");
    } else {
        document.getElementById("emote-popup").classList.remove("hidden");
    }
    emotePopupShown = !emotePopupShown;
}


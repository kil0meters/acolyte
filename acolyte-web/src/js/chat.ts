import {getEmotes, replaceTextWithEmotes} from './emotes';
import {Message, MessageList} from "./messageList";
import {UserList} from "./userList";
import {Autocompletion} from "./autocompletion";

export class MBChat {
    websocketLocation: string;
    conn: WebSocket;

    username: string;
    moderatorPerms: boolean;
    authorized: boolean = true;
    timeoutInterval: number;

    maxHeight: number;

    userList: UserList = new UserList();
    messageList: MessageList;

    entryBody: HTMLElement;
    autocompletionHelper: Autocompletion;

    constructor(websocketLocation: string, maxHeight: number, noEntry: boolean, username: string, moderatorPerms: boolean) {
        this.websocketLocation = websocketLocation;

        this.moderatorPerms = moderatorPerms === true;
        this.messageList = new MessageList('message-list', maxHeight, username, this.moderatorPerms);

        this.maxHeight = maxHeight;
        this.conn = new WebSocket(websocketLocation);

        this.username = username;

        if (!noEntry) {
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
            let message: Message = JSON.parse(m.data);
            if (message.text === "UNAUTHORIZED") {
                this.authorized = false;
            } else if (message.username === "user-add") {
                console.log("user joined: " + message.text);
                this.userList.add(message.text);
                this.autocompletionHelper.setUsers(this.userList.list);
            } else if (message.username === "user-remove" && this.autocompletionHelper) {
                console.log("user left: " + message.text);
                this.userList.remove(message.text);
                this.autocompletionHelper.setUsers(this.userList.list)
            } else if (message.username === "command-list" && this.autocompletionHelper) {
                this.autocompletionHelper.setCommands(message.text);
            } else if (message.username === "user-list") {
                for (let username of message.text) {
                    this.userList.add(username);
                }
                this.autocompletionHelper.setUsers(this.userList.list);
            } else {
                this.messageList.push(message);
            }
        });

        this.conn.addEventListener("close", (_) => {
            this.messageList.push({
                username: "Client",
                text: "Disconnected. Trying to reconnect in 5 seconds..."
            });
            this.timeoutInterval = window.setInterval(() => {
                console.log("Trying to reconnect");
                this.messageList.push({
                    username: "Client",
                    text: "Disconnected. Trying to reconnect in 5 seconds..."
                });

                this.conn = new WebSocket(this.websocketLocation);

                window.setTimeout(() => {
                    if (this.conn.readyState === this.conn.OPEN) {
                        this.messageList.push({"username": "Client", "text": "Reconnected."});
                        this.initializeConnection();
                        console.log("Reconnected");
                        clearInterval(this.timeoutInterval);
                    }
                }, 100);
            }, 5000);
        })
    }

    sendMessage(message: string) {
        this.conn.send(message);
    }

    initializeEntryBody() {
        document.getElementById('entry-body').addEventListener("keydown", (event) => {
            if (event.key === "Enter" && !event.shiftKey) {
                event.preventDefault();

                if (this.authorized) {
                    let entry = <HTMLInputElement>document.getElementById('entry-body');
                    this.sendMessage(entry.value);

                    this.autocompletionHelper.sentMessage(entry.value);
                    entry.value = "";

                } else {
                    toggleLoginPrompt();
                }
            }
        })
    }
}

let loginPromptShown: boolean = false;
let settingsShown: boolean = false;
let userListShown: boolean = false;
let emotePopupShown: boolean = false;

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


function renderLinksInElement(element) {
    element_default()(element, {
        "target": {
            url: "_blank",
        },
        "className": 'chat-link',
        "format": (value, type) => {
            if (type === 'url') {
                getLinkPreview(value)
                    .then(linkElement => {
                        element.appendChild(linkElement);
                        scrollToBottom(element.offsetHeight);
                    });


                return value.length > 50 ? value.slice(0, 50) + '...' : value
            }
        }
    });
}

function scrollToBottom(offsetHeight) {
    if ((window.innerHeight + window.scrollY + 64) >= (document.body.offsetHeight - offsetHeight)) {
        window.scrollTo(0, document.body.scrollHeight);
    }
}

async function getLinkPreview(value) {
    let response = await fetch('/api/v1/link-data?link=' + encodeURI(value), {
        credentials: 'include',
    });
    let data = await response.json();
    console.log(data);

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

class messageList_MessageList {
    constructor(messageListID, maxHeight, username, moderatorPerms) {
        this._list = [];
        this.messageListElement = document.getElementById(messageListID);
        this.emotes = getEmotes();
        this.maxHeight = maxHeight;

        this.username = username;
        this.moderatorPerms = moderatorPerms;

        this.currentCombo = [];
    }

    get list() {
        return this._list
    }

    buildMessage(message) {
        let messageElement = document.createElement("div");
        messageElement.id = message.id;
        messageElement.classList.add("chat-message");
        if (message.text.includes(this.username) && message.username !== this.username) {
            messageElement.classList.add("mentioned")
        } else if (message.username === this.username) {
            messageElement.classList.add("self")
        }


        let usernameElement = document.createElement("a");
        usernameElement.href = '#' + message.username;
        usernameElement.textContent = message.username;
        usernameElement.classList.add("message-username");

        let textElement = document.createElement("span");
        textElement.innerHTML = message.text;

        auto_render(textElement, {
            delimiters: [
                {left: "$$", right: "$$", display: true},
                {left: "$", right: "$", display: false}
            ]
        });

        renderLinksInElement(textElement);
        renderEmotesInElement(textElement);

        textElement.classList.add("message-text");

        messageElement.appendChild(usernameElement);

        if (this.moderatorPerms && message.id !== "00000000-0000-0000-0000-000000000000") {
            let removeMessageElement = document.createElement("button");
            removeMessageElement.textContent = "Remove";
            removeMessageElement.classList.add("remove-message-button");
            removeMessageElement.onclick = function () {
                liveChat.sendMessage(`/remove ${message.id}`)
            };

            messageElement.appendChild(removeMessageElement)
        }

        messageElement.appendChild(textElement);

        return messageElement;
    }

    push(message) {
        if (message.username === "delete-message") {
            this.removeByID(message.text);
        } else {
            this._list.push(message);

            let messageElement = this.buildMessage(message);
            this.messageListElement.appendChild(messageElement);

            this.checkForCombos();
            this.replaceComboListWithElement();

            scrollToBottom(messageElement.offsetHeight);
            while ((this.messageListElement.offsetHeight + 24) > this.maxHeight) { // offset for padding
                this.removeByIndex(0);
            }
        }
    }

    checkForCombos() {
        this.currentCombo = [];

        let recentMessage = this._list[this._list.length - 1];
        let prevMessage = recentMessage;

        if (!this.emotes.map((name) => name.name).includes(recentMessage.text)) {
            return;
        }

        for (let i = 1; recentMessage.text === prevMessage.text && i <= this._list.length; i++) {
            this.currentCombo.push(this._list[this._list.length - i]);
            prevMessage = this._list[this._list.length - i];
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

            let currentComboElement = document.createElement('div');

            let mostRecentMessage = this.messageListElement.lastChild;

            if (mostRecentMessage !== null) {
                if (Array.from(mostRecentMessage.classList).includes('combo-message') && mostRecentMessage.firstChild.getAttribute('alt') === this.currentCombo[0].text) {
                    currentComboElement = mostRecentMessage;
                }
            }

            currentComboElement.classList.add('chat-message', 'combo-message');
            currentComboElement.innerHTML = `${replaceTextWithEmotes(this.currentCombo[0].text)} ${this.currentCombo.length}x`;

            this.messageListElement.appendChild(currentComboElement);
        }
    }

    removeByIndex(index) {
        this._list.splice(index, 1);
        this.messageListElement.children[0].remove()
    }

    removeByID(id) {
        let i = this._list.indexOf(id);
        if (i !== -1) this._list.splice(i, 1);

        this.messageListElement.removeChild(document.getElementById(id))
    }
}
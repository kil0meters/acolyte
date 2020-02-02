// import katex from 'katex'
import renderMathInElement from 'katex/dist/contrib/auto-render'
import Split from 'split.js'


class MessageList {
  constructor(messageListElement, maxHeight) {
    this._list = []
    this.messageListElement = messageListElement
    this.maxHeight = maxHeight
  }

  static buildMessage(message) {
    let messageElement = document.createElement("div")
    messageElement.id = message.id
    messageElement.classList.add("chat-message")

    let usernameElement = document.createElement("a")
    usernameElement.href = '#' + message.username
    usernameElement.textContent = message.username
    usernameElement.classList.add("message-username")

    let textElement = document.createElement("span")
    textElement.textContent = message.text
    textElement.classList.add("message-text")

    messageElement.appendChild(usernameElement)
    messageElement.appendChild(textElement)
    // messageElement.appendChild(messageElement)

    return messageElement
  }

  push(message) {
    this._list.push(message)

    let messageElement = MessageList.buildMessage(message)
    this.messageListElement.appendChild(messageElement)

    renderMathInElement(messageElement,{delimiters: [
      {left: "$$", right: "$$", display: true},
      {left: "$", right: "$", display: false}
    ]})

    if ((window.innerHeight + window.scrollY + 64) >= (document.body.offsetHeight - messageElement.offsetHeight)) {
      window.scrollTo(0,document.body.scrollHeight)

      while ((this.messageListElement.offsetHeight + 24) > this.maxHeight) { // offset for padding
        this.removeByIndex(0)
      }
    }
  }

  removeByIndex(index) {
    this._list.splice(index, 1)
    this.messageListElement.children[0].remove()
  }

  removeByID(id) {
    for (i = 0; i < this._list.length; i++) {
      if (this._list[i].id === id) {
        this._list.splice(i, 1)
        this.messageListElement.removeChild(document.getElementById(id))
        break
      }
    }
  }

  get list() { 
    return this._list
  }
}

export class MBChat {
  constructor(maxHeight, noEntry) { 
    this.maxHeight = maxHeight
    this.conn = new WebSocket(`ws://${window.location.host}/api/v1/chat`)
    this.username = "username"
    this.isUnauthorized = false 
    this.entryBody = document.getElementById('entry-body')

    this.timeoutInterval = null

    this.messageList = new MessageList(document.getElementById('message-list'), this.maxHeight)

    if (noEntry != true) { 
      this.initializeEntryBody()
    }

    this.autocompletionHelper = new Autocompletion()
  }

  initializeConnection() {
    this.conn.addEventListener("message", (m) => {
      if (m.data == "UNAUTHORIZED") {
        this.isUnauthorized = true
      }
      else {
        let data = JSON.parse(m.data)

        this.messageList.push(data)
      }
    })

    this.conn.addEventListener("close", (m) => {
      this.messageList.push({ "username": "Client", "text": "Disconnected. Trying to reconnect in 5 seconds..."})
      this.timeoutInterval = setInterval(() => {
        console.log("Trying to reconnect")
        this.messageList.push({ "username": "Client", "text": "Disconnected. Trying to reconnect in 5 seconds..."})

        this.conn = new WebSocket('ws://localhost:8080/api/v1/chat')

        setTimeout(() => {
          if (this.conn.readyState == this.conn.OPEN) {
            this.messageList.push({ "username": "Client", "text": "Reconnected."})
            this.initializeConnection(this.conn)
            console.log("Reconnected")
            clearInterval(this.timeoutInterval)
          }
        }, 100)
      }, 5000)
    })
  }

  initializeEntryBody() {
    document.getElementById('entry-body').addEventListener("keydown", (event) => {
      if (event.key == "Enter" && !event.shiftKey) {
        event.preventDefault()

        if (this.isUnauthorized) {
          toggleLoginPrompt()
        } else {
          this.autocompletionHelper.sentMessage(document.getElementById('entry-body').value)

          this.conn.send(JSON.stringify({
            "username": this.username,
            "text": document.getElementById('entry-body').value,
          }))
          
          document.getElementById('entry-body').value = ""
        }
      }
    })
  }
}

class Autocompletion {
  constructor(entry) {
    this.entry = document.getElementById('entry-body')
    this.popup = document.getElementById('autocompletion-popup')

    this.chatCommands = ["/ban", "/mute", "/addcommand", "/toggle-dark-mode"]

    this.suggestions = []
    this.tabIndex = 1

    this.previousMessages = []
    this.messageIndex = 0

    this.currentValue = ""

    this.registerEventListeners()
  }

  setPopupToSuggestions() {
    if (this.suggestions == []) {
      this.popup.classList.add('hidden') 
    } else {
      while (this.popup.firstChild) {
        this.popup.removeChild(this.popup.firstChild)
      }

      this.popup.classList.remove('hidden')

      for (let suggestion of this.suggestions) {
        let suggestionElement = document.createElement('p')
        suggestionElement.textContent = suggestion
        
        this.popup.appendChild(suggestionElement)
      }
    }
  }

  setHighlightedSuggestion(index) {
    for (let suggestionElement of this.popup.children) {
      suggestionElement.classList.remove('highlighted')
    }

    this.popup.children[index].classList.add('highlighted')
  }

  sentMessage(message) {
    if (this.previousMessages[this.previousMessages.length-1] != message) {
      this.previousMessages.push(message)
    }
    this.messageIndex = 0
  }

  registerEventListeners() {
    this.entry.addEventListener("keydown", (event) => {
      if (event.key == "Tab") {
        event.preventDefault()

        if (!event.shiftKey) {
          this.tabIndex = Math.min(this.suggestions.length, this.tabIndex+1)
        } else {
          this.tabIndex = Math.max(1, this.tabIndex-1)
        }

        if (this.suggestions.length != 0) {
          this.entry.value = this.suggestions[this.suggestions.length-this.tabIndex]
          this.setHighlightedSuggestion(this.suggestions.length-this.tabIndex)
        }
      }

      else if (event.key == "ArrowDown") {
        event.preventDefault()
        this.messageIndex = Math.max(this.messageIndex-1, 0)
        if (this.messageIndex == 0 ) {
          this.entry.value = this.currentValue
        } else {
          this.entry.value = this.previousMessages[this.previousMessages.length - this.messageIndex]
        }
      }
      else if (event.key == "ArrowUp") {
        event.preventDefault()
        this.messageIndex = Math.min(this.messageIndex+1, this.previousMessages.length)
        this.entry.value = this.previousMessages[this.previousMessages.length - this.messageIndex]
      }
    })

    this.entry.addEventListener("keyup", (event) => {
      if (event.key != "ArrowUp" &&
          event.key != "ArrowDown" &&
          event.key != "Tab" &&
          // event.key != "Enter" &&
          event.key != "Shift" &&
          event.key != "Control" &&
          event.key != "Meta" &&
          event.key != "Alt") {

        let text = this.entry.value
      
        if (this.messageIndex == 0) {
          this.currentValue = text
        }

        this.tabIndex = 0
        this.suggestions = []
        if (text == "") {
          this.messageIndex = 0
          this.setPopupToSuggestions()
        } else {
          for (let suggestion of this.chatCommands) {
            if (suggestion.startsWith(text)) {
              this.suggestions.push(suggestion)
            }
          }

          this.setPopupToSuggestions()
        }
      }
    })
  }
}

var loginPromptShown = false
var settingsShown = false

export function toggleLoginPrompt() {
  if (loginPromptShown) {
    document.getElementById("login-prompt").classList.add('hidden')
  } else {
    document.getElementById("login-prompt").classList.remove('hidden')
  }
  loginPromptShown = !loginPromptShown
}

export function toggleSettings() {
  if (settingsShown) {
    document.getElementById("settings-overlay").classList.add('hidden')
  } else {
    document.getElementById("settings-overlay").classList.remove('hidden')
  }
  settingsShown = !settingsShown
}

export function setupSplitpanes(left, right) {
  Split([left, right], {
    sizes: [70, 30],
    gutterSize: 10,
    minSize: 200,
  })
}
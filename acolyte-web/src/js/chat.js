// import katex from 'katex'
import renderMathInElement from 'katex/dist/contrib/auto-render'
import Split from 'split.js'
import { replaceTextWithEmotes, getEmotes } from './emotes.js'

class MessageList {
  constructor(messageListElement, maxHeight) {
    this._list = []
    this.messageListElement = messageListElement
    this.emotes = getEmotes()
    this.maxHeight = maxHeight

    this.currentCombo = []
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
    textElement.innerHTML = replaceTextWithEmotes(message.text)
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

    this.checkForCombos()
    this.replaceComboListWithElement()

    if ((window.innerHeight + window.scrollY + 64) >= (document.body.offsetHeight - messageElement.offsetHeight)) {
      window.scrollTo(0,document.body.scrollHeight)

      while ((this.messageListElement.offsetHeight + 24) > this.maxHeight) { // offset for padding
        this.removeByIndex(0)
      }
    }
  }

  checkForCombos() {
    this.currentCombo = []

    let recentMessage = this._list[this._list.length-1]
    var prevMessage = recentMessage

    if (!this.emotes.includes(recentMessage.text)) {
      return
    }

    for (let i = 1; recentMessage.text === prevMessage.text && i <= this._list.length; i++) {
      this.currentCombo.push(this._list[this._list.length-i])
      prevMessage = this._list[this._list.length-i]
    }

    if (this.currentCombo[this.currentCombo.length-1].text != recentMessage.text) {
      this.currentCombo.pop() // combo broken
    }
  }

  replaceComboListWithElement() {
    if (this.currentCombo.length > 1) {
      for (let message of this.currentCombo) {
        let elementToRemove = document.getElementById(message.id)

        if (elementToRemove != undefined) {
          this.messageListElement.removeChild(document.getElementById(message.id))
        }
      }

      let currentComboElement = document.createElement('div')


      let mostRecentMessage = this.messageListElement.lastChild

      if (mostRecentMessage != null) {
        if (Array.from(mostRecentMessage.classList).includes('combo-message') && mostRecentMessage.firstChild.getAttribute('alt') == this.currentCombo[0].text) {
          currentComboElement = mostRecentMessage;
        }
      }
      replaceTextWithEmotes

      currentComboElement.classList.add('chat-message', 'combo-message')
      currentComboElement.innerHTML = `${replaceTextWithEmotes(this.currentCombo[0].text)} ${this.currentCombo.length}x`

      this.messageListElement.appendChild(currentComboElement)
      // if (currentComboElement.classList)
      console.log(currentComboElement);
    }
  }

  removeByIndex(index) {
    this._list.splice(index, 1)
    this.messageListElement.children[0].remove()
  }

  removeByID(id) {
    let i = this._list.indexOf(id)
    if (i != -1) this._list.splice(i, 1)

    this.messageListElement.removeChild(document.getElementById(id))
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

    let emotePopup = document.getElementById('emote-popup')
    getEmotes().forEach((emote) => {
      emotePopup.innerHTML += replaceTextWithEmotes(emote, `document.getElementById('entry-body').value += '${emote} '`)
    })

    this.timeoutInterval = null

    this.messageList = new MessageList(document.getElementById('message-list'), this.maxHeight)

    if (noEntry != true) { 
      this.initializeEntryBody()
      this.autocompletionHelper = new Autocompletion()
    }
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
    this.emotes = getEmotes()

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

        if (this.emotes.includes(suggestion)) {
          suggestionElement.innerHTML = replaceTextWithEmotes(suggestion) + ' ' + suggestion
        } else {
          suggestionElement.textContent = suggestion
        }
        
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
          let words = this.entry.value.split(' ')
          words[words.length-1] = this.suggestions[this.suggestions.length-this.tabIndex]

          this.entry.value = words.join(' ')
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
          let completionText = text.toLowerCase()

          for (let suggestion of this.chatCommands) {
            if (suggestion.startsWith(completionText)) {
              this.suggestions.push(suggestion)
            }
          }

          for (let suggestion of this.emotes) {
            let words = completionText.split(' ')
            let word = words[words.length - 1]
            if (word != '') {
              if (suggestion.toLowerCase().startsWith(words[words.length - 1])) {
                this.suggestions.push(suggestion)
              }
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
var emotePopupShown = false

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

export function toggleEmotePopup() {
  if (emotePopupShown) {
    document.getElementById("emote-popup").classList.add("hidden")
  } else {
    document.getElementById("emote-popup").classList.remove("hidden")
  }
  emotePopupShown = !emotePopupShown
}

export function setupSplitpanes(left, right) {
  Split([left, right], {
    sizes: [70, 30],
    gutterSize: 10,
    minSize: 200,
  })
}


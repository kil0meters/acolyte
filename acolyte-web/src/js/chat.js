// import katex from 'katex'
import renderMathInElement from 'katex/dist/contrib/auto-render'
import { replaceTextWithEmotes, getEmotes, renderEmotesInElement } from './emotes.js'
import linkifyElement from 'linkifyjs/element'

class MessageList {
  constructor(messageListElement, maxHeight, moderatorPerms) {
    this._list = []
    this.messageListElement = messageListElement
    this.emotes = getEmotes()
    this.maxHeight = maxHeight

    this.moderatorPerms = moderatorPerms

    this.currentCombo = []
  }

  static buildMessage(message, moderatorPerms, conn) {
    let messageElement = document.createElement("div")
    messageElement.id = message.id
    messageElement.classList.add("chat-message")

    let usernameElement = document.createElement("a")
    usernameElement.href = '#' + message.username
    usernameElement.textContent = message.username
    usernameElement.classList.add("message-username")

    let textElement = document.createElement("span")
    textElement.innerHTML = message.text

    renderMathInElement(textElement,{delimiters: [
      {left: "$$", right: "$$", display: true},
      {left: "$", right: "$", display: false}
    ]})
    linkifyElement(textElement)
    renderEmotesInElement(textElement)

    textElement.classList.add("message-text")

    messageElement.appendChild(usernameElement)
    messageElement.appendChild(textElement)

    console.log(moderatorPerms);
    if (moderatorPerms && message.id != "00000000-0000-0000-0000-000000000000") {
      let removeMessageElement = document.createElement("button") 
      removeMessageElement.textContent = "Remove"
      removeMessageElement.classList.add("remove-message-button")
      removeMessageElement.onclick = function() {
        liveChat.sendMessage(`/remove ${message.id}`)
      }

      messageElement.appendChild(removeMessageElement)
    }

    return messageElement
  }

  push(message) {
    if (message.username == "delete") {
      this.removeByID(message.text)
      return
    }

    this._list.push(message)

    let messageElement = MessageList.buildMessage(message, this.moderatorPerms)
    this.messageListElement.appendChild(messageElement)

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
  constructor(maxHeight, noEntry, moderatorPerms) { 
    if (location.protocol == 'https:') {
      this.wsProtocol = 'wss:'  
    } else {
      this.wsProtocol = 'ws:'
    }

    this.noEntry = noEntry == true
    this.moderatorPerms = moderatorPerms == true

    this.maxHeight = maxHeight
    this.conn = new WebSocket(`${this.wsProtocol}//${window.location.host}/api/v1/chat`)
    this.username = "username"
    this.isUnauthorized = false 

    this.timeoutInterval = null


    this.messageList = new MessageList(document.getElementById('message-list'), this.maxHeight, this.moderatorPerms)

    if (this.noEntry == false) {
      this.entryBody = document.getElementById('entry-body')
      let emotePopup = document.getElementById('emote-popup')
  
      getEmotes().forEach((emote) => {
        emotePopup.innerHTML += replaceTextWithEmotes(emote.name, `document.getElementById('entry-body').value += '${emote.name} '`)
      })

      this.initializeEntryBody()
      this.autocompletionHelper = new Autocompletion()
    }
  }

  initializeConnection() {
    this.conn.addEventListener("message", (m) => {
      let data = JSON.parse(m.data)
      if (data.constructor == Array) { // this is probably a bad way to select the command list
        if (this.noEntry == false) {
          this.autocompletionHelper.setCommands(data)
        }
      } else if (data.text == "UNAUTHORIZED") {
        this.isUnauthorized = true  
      } else {
        this.messageList.push(data)
      }
    })

    this.conn.addEventListener("close", (m) => {
      this.messageList.push({ "username": "Client", "text": "Disconnected. Trying to reconnect in 5 seconds..."})
      this.timeoutInterval = setInterval(() => {
        console.log("Trying to reconnect")
        this.messageList.push({ "username": "Client", "text": "Disconnected. Trying to reconnect in 5 seconds..."})

        this.conn = new WebSocket(`${this.wsProtocol}//${window.location.host}/api/v1/chat`)

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

  sendMessage(message) {
    this.conn.send(JSON.stringify({
      "username": this.username,
      "text": message,
    }))
  }

  initializeEntryBody() {
    document.getElementById('entry-body').addEventListener("keydown", (event) => {
      if (event.key == "Enter" && !event.shiftKey) {
        event.preventDefault()

        if (this.isUnauthorized) {
          toggleLoginPrompt()
        } else {
          this.autocompletionHelper.sentMessage(document.getElementById('entry-body').value)

          this.sendMessage(document.getElementById('entry-body').value)

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

    this.chatCommands = []
    this.emotes = getEmotes()

    this.suggestions = []
    this.tabIndex = 1

    this.previousMessages = []
    this.messageIndex = 0

    this.currentValue = ""

    this.registerEventListeners()
  }

  setCommands(commands) {
    this.chatCommands = commands
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
          suggestionElement.innerHTML = suggestion.description + suggestion.name
        } else if (this.chatCommands.includes(suggestion)) {
          suggestionElement.textContent = suggestion.name + ' ' + suggestion.description
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

          let highlighted = this.suggestions[this.suggestions.length-this.tabIndex]
          words[words.length-1] = this.suggestions[this.suggestions.length-this.tabIndex].name

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
            if (suggestion.name.toLowerCase().startsWith(completionText)
            || (completionText.trim().toLowerCase().startsWith(suggestion.name.toLowerCase()))) {
              this.suggestions.push(suggestion)
            }
          }

          for (let suggestion of this.emotes) {
            let words = completionText.split(' ')
            let word = words[words.length - 1]
            if (word != '') {
              if (suggestion.name.toLowerCase().startsWith(words[words.length - 1])) {
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


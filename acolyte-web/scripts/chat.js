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

class MBChat {
  constructor(maxHeight, noEntry) { 
    this.maxHeight = maxHeight
    this.conn = new WebSocket('ws://localhost:8080/api/v1/chat')
    this.username = "username"
    this.isUnauthroized = true
    this.entryBody = document.getElementById('entry-body')

    this.chatCommands = ["/ban", "/mute"]
    this.timeoutInterval = null

    this.messageList = new MessageList(document.getElementById('message-list'), this.maxHeight)

    if (noEntry != true) { 
      this.initializeEntryBody()
    }
  }

  initializeConnection() {
    this.conn.addEventListener("message", (m) => {
      console.log("sending message")

      if (m.data == "UNAUTHORIZED") {
        this.isUnauthroized = true
      }
      else {
        this.isUnauthroized = false

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
    document.getElementById('entry-body').addEventListener("keyup", (event) => {
      if (event.key == "Enter" && !event.shiftKey) {
        if (this.isUnauthroized) {
          toggleLoginPrompt()
        } else {
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

  }

  registerEventListeners() {
    document.getElementById('entry-body').addEventListener("keyup", function(event) {
      text = document.getElementById('entry-body').value
  
      suggestions = []
  
      for (autocompleteOption of chatCommands) {
        if (autocompleteOption.startsWith(text)) {
          suggestions.push(autocompleteOption)
        }
      }
  
      if (event.key == "Tab") {
        
      }
  
      console.log(suggestions)
    })
  }
}

loginPromptShown = false
settingsShown = false

function toggleLoginPrompt() {
  if (loginPromptShown) {
    document.getElementById("login-prompt").classList.add('hidden')
  } else {
    document.getElementById("login-prompt").classList.remove('hidden')
  }
  loginPromptShown = !loginPromptShown
}

function toggleSettings() {
  if (settingsShown) {
    document.getElementById("settings-overlay").classList.add('hidden')
  } else {
    document.getElementById("settings-overlay").classList.remove('hidden')
  }
  settingsShown = !settingsShown
}


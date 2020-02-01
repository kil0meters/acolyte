var conn = new WebSocket('ws://localhost:8080/api/v1/chat')
messageListElement = document.getElementById('message-list')
username = "username"
settingsShown = false
loginPromptShown = false
isUnauthroized = true
const entryBody = document.getElementById('entry-body')

var chatCommands = ["/ban", "/mute"]

messageList = {
  _list: [],
  addMessageListener: function(message) {},
  removeMessageListener: function(id) {},
  push(message) {
    this._list.push(message)
    this.addMessageListener(message)
  },
  remove(id) {
    for (i = 0; i < this._list.length; i++) {
      if (this._list[i].id === id) {
        this._list.splice(i, 1)
        this.removeMessageListener(id)
        break
      }
    }
  },
  get list() { 
    return this._list
  },
}

function buildMessage(message) {
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

messageList.addMessageListener = function (message) {
  // messageListElement.children = []
  // for (let message in value) {
  
  messageElement = buildMessage(message)

  messageListElement.appendChild(messageElement)

  console.log('window.scrollY: ' + (window.innerHeight + window.scrollY));
  console.log('document.body.scrollHeight: ' + document.body.scrollHeight);
  console.log('document.body.offsetHeight: ' + document.body.offsetHeight);

  renderMathInElement(messageElement,{delimiters: [
    {left: "$$", right: "$$", display: true},
    {left: "$", right: "$", display: false}
  ]})

  if ((window.innerHeight + window.scrollY) >= (document.body.offsetHeight - messageElement.offsetHeight)) {
    window.scrollTo(0,document.body.scrollHeight)
  }

}

messageList.removeMessageListener = function (id) {
  messageListElement.removeChild(document.getElementById(id))
}

// {"username": "<USERNAME>", "text": "<TEXT>"}

function initializeConnection(_conn) {
  _conn.addEventListener("message", (m) => {
    // console.log(m)

    if (m.data == "UNAUTHORIZED") {
      isUnauthroized = true
    }
    else {
      isUnauthroized = false
      data = JSON.parse(m.data)
      console.log(data)
      // messageData = JSON.parse(data.body)

      // console.log(messageData);
      messageList.push(data)
    }
  })

  _conn.addEventListener("close", (m) => {
    messageList.push({ "username": "Client", "text": "Disconnected. Trying to reconnect in 5 seconds..."})
    interval = setInterval(function() {
      console.log("Trying to reconnect")
      messageList.push({ "username": "Client", "text": "Disconnected. Trying to reconnect in 5 seconds..."})
      conn = new WebSocket('ws://localhost:8080/api/v1/chat')

      setTimeout(function () {
        if (conn.readyState == conn.OPEN) {
          messageList.push({ "username": "Client", "text": "Reconnected."})
          initializeConnection(conn)
          console.log("Reconnected")
          clearInterval(interval)
        }
      }, 100)
    }, 5000)
  })
}

initializeConnection(conn)

document.getElementById('entry-body').addEventListener("keyup", function(event) {
  if (event.key == "Enter" && !event.shiftKey) {
    if (isUnauthroized) {
      toggleLoginPrompt()
    } else {
      conn.send(JSON.stringify({
        "username": username,
        "text": entryBody.value,
        }))
      entryBody.value = ""
    }
  }
})

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

  console.log(suggestions);
})

document.getElementById('settings-username-input').addEventListener("keyup", function(event) {
  username = document.getElementById('settings-username-input').value
})

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
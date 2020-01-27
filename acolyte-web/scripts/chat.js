var conn = new WebSocket('ws://localhost:3000/api/v1/chat')
messageListElement = document.getElementById('message-list')
const entryBody = document.getElementById('entry-body')

messageList = {
  _list: [],
  messageListener: function(message) {},
  push(message) {
    this._list.push(message)
    this.messageListener(message)
  },
  get list() { 
    return this._list
  },
}

function buildMessage(message) {
  let messageElement = document.createElement("div")
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

messageList.messageListener = function (message) {
  // messageListElement.children = []
  // for (let message in value) {
  messageListElement.appendChild(buildMessage(message))
  window.scrollTo(0,document.body.scrollHeight)
  // }
}

// {"username": "<USERNAME>", "text": "<TEXT>"}

function initializeConnection(_conn) {
  _conn.addEventListener("message", (m) => {
    console.log(m)
  
    data = JSON.parse(m.data)
    let text = document.createTextNode(data.body)
    messageList.push({"username": "username", "text": data.body})
  })

  _conn.addEventListener("close", (m) => {
    messageList.push({"username": "System", "text": "Disconnected. Trying to reconnect in 5 seconds..."})
    interval = setInterval(function() {
      console.log("Trying to reconnect")
      messageList.push({"username": "System", "text": "Disconnected. Trying to reconnect in 5 seconds..."})
      conn = new WebSocket('ws://localhost:3000/api/v1/chat')

      setTimeout(function () {
        if (conn.readyState == conn.OPEN) {
          messageList.push({"username": "System", "text": "Reconnected."})
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
    conn.send(entryBody.value)
    entryBody.value = ""
  }
})
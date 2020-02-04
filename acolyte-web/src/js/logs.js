export class LogSearch {
  constructor() {
    this.searchbar = document.getElementById('logs-search-bar')
    this.results = document.getElementById('logs-results')

    this.urlParams = new URLSearchParams(window.location.search)

    this.requestURL = location.protocol + '//' + location.host + '/api/v1/search-logs?search='
  }

  updateSearch() {
    if (this.searchbar.value != "") {   
      fetch(this.requestURL + encodeURIComponent(this.searchbar.value))
        .then(response => response.json())
        .then((messages) => {
          this.results.innerHTML = ''
          for (let message of messages) {
            let logElement = document.createElement('p')
            let logLink = document.createElement('a')

            let date = new Date(message.time)
            let dateStr = date.getFullYear() + "-" +
                          ("00" + (date.getMonth() + 1)).slice(-2) + "-" +
                          ("00" + date.getDate()).slice(-2) + " " +
                          ("00" + date.getHours()).slice(-2) + ":" +
                          ("00" + date.getMinutes()).slice(-2) + ":" +
                          ("00" + date.getSeconds()).slice(-2)

            logLink.textContent = `${dateStr} [${message.username}] ${message.message}`
            logLink.href = `/logs/messages/${message.message_id}`

            logElement.appendChild(logLink)


            this.results.appendChild(logElement)
          }
          console.log(JSON.stringify(messages))
        })
    } else {

    }
  }

  initializeCallbacks() {
    this.searchbar.value = this.urlParams.get('search')
    this.updateSearch()

    this.searchbar.addEventListener('keydown', () => {
      this.updateSearch()
    })
  }
}
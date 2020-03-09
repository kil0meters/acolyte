import {renderEmotesInElement} from './emotes'

export class LogSearch {

    searchbar: HTMLInputElement;
    results: HTMLDivElement;

    urlParams: URLSearchParams;
    requestURL: string;

    constructor() {
        this.searchbar = <HTMLInputElement>document.getElementById('logs-search-bar');
        this.results = <HTMLDivElement>document.getElementById('logs-results');

        this.urlParams = new URLSearchParams(window.location.search);

        this.requestURL = location.protocol + '//' + location.host + '/api/v1/search-logs?search='
    }

    updateSearch() {
        if (this.searchbar.value !== "") {
            fetch(this.requestURL + encodeURIComponent(this.searchbar.value))
                .then(response => response.json())
                .then((messages) => {
                    this.results.innerHTML = '';
                    for (let message of messages) {
                        let logElement = document.createElement('p');
                        let logLink = document.createElement('a');

                        let date = new Date(message.time);
                        let dateStr = date.getFullYear() + "-" +
                            ("00" + (date.getMonth() + 1)).slice(-2) + "-" +
                            ("00" + date.getDate()).slice(-2) + " " +
                            ("00" + date.getHours()).slice(-2) + ":" +
                            ("00" + date.getMinutes()).slice(-2) + ":" +
                            ("00" + date.getSeconds()).slice(-2);

                        logLink.innerHTML = `${dateStr} [${message.username}] ${message.message}`;
                        logLink.href = `/logs/messages/${dateStr.split(' ')[0]}#${message.message_id}`;

                        logElement.appendChild(logLink);

                        this.results.appendChild(logElement);
                    }
                    renderEmotesInElement(document.getElementById("logs-results"))
                })
        } else {

        }
    }

    initializeCallbacks() {
        this.searchbar.value = this.urlParams.get('search');
        this.updateSearch();

        this.searchbar.addEventListener('keyup', () => {
            this.updateSearch();
        })
    }
}

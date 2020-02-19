class UserList {
    constructor() {
        this.userList = [];
    }

    get list() {
        return this.userList;
    }

    static buildListElement(username) {
        let listElement = document.createElement("span");
        listElement.classList.add("list-card");

        console.log("username: " + username);
        listElement.textContent = username;

        return listElement;
    }

    add(username) {
        let index = this.userList.indexOf(username);
        if (index === -1) this.userList.push(username);
        this.updateList();
    }

    remove(username) {
        let index = this.userList.indexOf(username);
        if (index !== -1) this.userList.splice(index, 1);
        this.updateList();
    }

    updateList() {
        let overlay = document.getElementById("user-list-list");
        overlay.innerHTML = "";

        this.userList.forEach((username) => {
            overlay.appendChild(UserList.buildListElement(username))
        })
    }
}
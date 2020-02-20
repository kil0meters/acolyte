export class UserList {

    public list: string[] = [];

    static buildListElement(username: string) {
        let listElement = document.createElement("span");
        listElement.classList.add("list-card");

        console.log("username: " + username);
        listElement.textContent = username;

        return listElement;
    }

    add(username: string) {
        let index = this.list.indexOf(username);
        if (index === -1) this.list.push(username);
        this.updateList();
    }

    remove(username: string) {
        let index = this.list.indexOf(username);
        if (index !== -1) this.list.splice(index, 1);
        this.updateList();
    }

    updateList() {
        let overlay = document.getElementById("user-list-list");
        overlay.innerHTML = "";

        for (let username of this.list) {
            overlay.appendChild(UserList.buildListElement(username))
        }
    }
}
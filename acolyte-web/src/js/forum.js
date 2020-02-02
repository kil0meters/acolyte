export class TransientHeader {
    constructor(id) {
        this.oldWindowScrollY = 0
        this.id = id


        window.addEventListener("scroll", () => {
            console.log(window.scrollY - this.oldWindowScrollY);
            if (window.scrollY > this.oldWindowScrollY && window.scrollY > 64) {
                document.getElementById(this.id).classList.add('header-hidden')
            } else {
                document.getElementById(this.id).classList.remove('header-hidden')
            }
            this.oldWindowScrollY = window.scrollY
        })
    }
}

class ForumPost extends HTMLElement {
    constructor() {
        super();
    }

    connectedCallback() {
        let title = this.getAttribute('title')
        let postID = this.getAttribute('post-id')

        this.innerHTML = `
            <div class="post-thumbnail"></div>
            <a class="post-title" href="forum/posts/${postID}">${title}</a>
            <ul class="post-options">
                <li class="post-option expander">+</li>
                <li class="post-option">Comments</li>
                <li class="post-option">Report</li>
            </ul>
        `;
    }

}

customElements.define('forum-post', ForumPost)

oldWindowScrollY = 0

window.addEventListener("scroll", function () {
    if (window.scrollY > oldWindowScrollY && window.scrollY > 64) {
        document.getElementById('forum-header').classList.add('header-hidden')
    } else {
        document.getElementById('forum-header').classList.remove('header-hidden')
    }
    oldWindowScrollY = window.scrollY
})

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

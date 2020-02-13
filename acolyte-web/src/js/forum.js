export class TransientHeader {
    constructor(id) {
        this.oldWindowScrollY = 0;
        this.id = id;


        window.addEventListener("scroll", () => {
            if (window.scrollY > this.oldWindowScrollY && window.scrollY > 64) {
                document.getElementById(this.id).classList.add('header-hidden');
            } else {
                document.getElementById(this.id).classList.remove('header-hidden');
            }
            this.oldWindowScrollY = window.scrollY;
        })
    }
}

class ForumPost extends HTMLElement {
    constructor() {
        super();
    }

    connectedCallback() {
        let title = this.getAttribute('title');
        let postID = this.getAttribute('post-id');

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

class CommentEditor extends HTMLElement {
    constructor() {
        super();
    }

    connectedCallback() {
        let username = this.getAttribute('username');
        let parentID = this.getAttribute('parent-id');

        this.innerHTML = `
        <div class="comment-editor">
            <span>Comment as <a href="/user/${username}">${username}</a></span>
            <form action="/forum/posts/${parentID}" method="POST">
                <textarea name="body" class="authorize-form-input" id="comment-entry" cols="30" rows="10"
                          placeholder="make a nice comment please"></textarea>
                <button class="authorize-form-button">COMMENT</button>
            </form>
        </div>`
    }

}

export function toggleReplyEditorVisibility(id, username) {
    let editor = document.querySelector(`#${id} .comment-container comment-editor`);
    let commentContainer = document.querySelector(`#${id} .comment-container`);

    if (editor === undefined || editor === null) {

        editor = document.createElement('comment-editor');
        editor.setAttribute('username', username);
        editor.setAttribute('parent-id', id);

        commentContainer.appendChild(editor)
    } else {
        editor.remove()
    }
}

customElements.define('forum-post', ForumPost);
customElements.define('comment-editor', CommentEditor);

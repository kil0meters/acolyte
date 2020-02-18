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

function formatDate(date) {
    let monthNames = [
        "January", "February", "March",
        "April", "May", "June", "July",
        "August", "September", "October",
        "November", "December"
    ];

    let day = date.getDate();
    let monthIndex = date.getMonth();
    let year = date.getFullYear();

    return day + ' ' + monthNames[monthIndex] + ', ' + year;
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
            <div class="post-content">
                <a class="post-title" href="forum/posts/${postID}">${title}</a>
                <ul class="post-options">
                    <li class="post-option expander">+</li>
                    <li class="post-option">Comments</li>
                    <li class="post-option">Report</li>
                </ul>
            </div>
        `;
    }
}

class LinkPreview extends HTMLElement {
    constructor() {
        super();
    }

    connectedCallback() {
        let title = this.getAttribute('title');
        let publishedDate = this.getAttribute('published-date');
        let link = this.getAttribute('link');
        let content = this.getAttribute('content');

        publishedDate = formatDate(new Date(publishedDate));

        if (title === "") {
            this.innerHTML = `
                <div class="link-preview-container card">
                    <a class="link-preview-link" href="${link}">${link}</a>
                </div>
            `
        } else {
            this.innerHTML = `
            <div class="link-preview-container card">
                <div>
                    <a class="article-title" href="${link}">${title}</a>
                    <span class="article-date"> ${publishedDate}</span>
                </div>
                <div class="article-description">
                    ${content}
                </div>
            </div>
            `
        }
    }
}

class CommentEditor extends HTMLElement {
    constructor() {
        super();
    }

    connectedCallback() {
        let username = this.getAttribute('username');
        let parentID = this.getAttribute('parent-id');
        let postID = this.getAttribute('post-id');

        this.innerHTML = `
        <div class="comment-editor">
            <span>Comment as <a href="/user/${username}">${username}</a></span>
            <form action="/forum/posts/${parentID}" method="POST">
                <input type="hidden" name="post-id" value="${postID}">
                <textarea name="body" class="authorize-form-input" id="comment-entry" cols="30" rows="10"
                          placeholder="make a nice comment please"></textarea>
                <button class="authorize-form-button">COMMENT</button>
            </form>
        </div>`
    }
}

export function toggleReplyEditorVisibility(id, postID, username) {
    let editor = document.querySelector(`#${id} .comment-container .comment-actions comment-editor`);
    let commentActions = document.querySelector(`#${id} .comment-container .comment-actions`);

    if (editor === undefined || editor === null) {

        editor = document.createElement('comment-editor');
        editor.setAttribute('username', username);
        editor.setAttribute('parent-id', id);
        editor.setAttribute('post-id', postID);

        commentActions.appendChild(editor)
    } else {
        editor.remove()
    }
}

customElements.define('forum-post', ForumPost);
customElements.define('comment-editor', CommentEditor);
customElements.define('link-preview', LinkPreview);

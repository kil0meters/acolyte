function formatDate(date: Date) {
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
                    <a class="article-title" target="_blank" href="${link}">${title}</a>
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
        let threadID = this.getAttribute('post-id');

        this.innerHTML = `
        <div class="comment-editor">
            <span>Comment as <a href="/user/${username}">${username}</a></span>
            <form action="/forum/create-comment" method="POST">
                <input type="hidden" name="thread_id" value="${threadID}">
                <input type="hidden" name="parent_id" value="${parentID}">
                <textarea name="body" class="authorize-form-input" id="comment-entry" cols="30" rows="10"
                          placeholder="make a nice comment please"></textarea>
                <button class="authorize-form-button">COMMENT</button>
            </form>
        </div>`
    }
}

export function toggleReplyEditorVisibility(parentID: string, postID: string, username: string) {
    let editor = document.querySelector(`#${parentID} .comment-container .comment-actions comment-editor`);
    let commentActions = document.querySelector(`#${parentID} .comment-container .comment-actions`);

    if (editor === undefined || editor === null) {

        editor = document.createElement('comment-editor');
        editor.setAttribute('username', username);
        editor.setAttribute('parent-id', parentID);
        editor.setAttribute('post-id', postID);

        commentActions.appendChild(editor)
    } else {
        editor.remove()
    }
}

customElements.define('comment-editor', CommentEditor);
customElements.define('link-preview', LinkPreview);

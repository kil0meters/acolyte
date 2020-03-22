let username = "";

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
        let threadId = this.getAttribute('thread-id');
        let parentId = this.getAttribute('parent-id');

        this.innerHTML = `
        <div class="comment-editor">
            <span>Comment as <a href="/user/${username}">${username}</a></span>
            <form action="/forum/create-comment" method="POST">
                <input type="hidden" name="thread_id" value="${threadId}">
                <input type="hidden" name="parent_id" value="${parentId}">
                <textarea name="body" class="authorize-form-input" id="comment-entry" cols="30" rows="10"
                          placeholder="make a nice comment please"></textarea>
                <button class="authorize-form-button">COMMENT</button>
            </form>
        </div>`
    }
}

export function toggleReplyEditorVisibility(parentId: string, commentID: string) {
    let editor = document.querySelector(`#${commentID} .content .actions comment-editor`);

    if (editor === undefined || editor === null) {
        let commentActions = document.querySelector(`#${commentID} .content .actions`);

        editor = document.createElement('comment-editor');
        editor.setAttribute('username', username);
        editor.setAttribute('thread-id', parentId.split('-')[0]);
        editor.setAttribute('parent-id', parentId);

        commentActions.appendChild(editor)
    } else {
        editor.remove()
    }
}

customElements.define('comment-editor', CommentEditor);
customElements.define('link-preview', LinkPreview);

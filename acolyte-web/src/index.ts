import {MBChat, toggleEmotePopup, toggleLoginPrompt, toggleSettings, toggleUserList} from './js/chat';
import {toggleReplyEditorVisibility} from './js/forum';
import {toggleDarkMode} from './js/theme';
import {renderEmotesInElement} from './js/emotes';
import {renderMarkdownInElement} from './js/markdown';
import {renderLinksInElement} from "./js/messageList";
import {LogSearch} from './js/logs';

import Split from 'split.js';

import renderMathInElement from 'katex/dist/contrib/auto-render';

import './css/chat.css';
import './css/fonts.css';
import './css/homepage.css';
import './css/livestream.css';
import './css/markdown.css';
import './css/forum.css';
import './css/post-editor.css';
import './css/logs.css';

import '../node_modules/katex/dist/katex.css';

export {
    MBChat,
    toggleLoginPrompt,
    toggleUserList,
    toggleSettings,
    toggleEmotePopup,

    renderMathInElement,
    toggleReplyEditorVisibility,
    renderLinksInElement,
    renderEmotesInElement,
    renderMarkdownInElement,

    LogSearch,

    toggleDarkMode,
    Split,
}
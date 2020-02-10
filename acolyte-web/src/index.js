import {MBChat, toggleEmotePopup, toggleLoginPrompt, toggleSettings, toggleUserList} from './js/chat.js';
import {TransientHeader} from './js/forum.js';
import {toggleDarkMode} from './js/toggle-darkmode.js';
import {renderEmotesInElement} from './js/emotes.js';
import {renderMarkdownInElement} from './js/markdown.js';
import {LogSearch} from './js/logs.js';
import Split from 'split.js';

import renderMathInElement from 'katex/dist/contrib/auto-render';

import './css/chat.css';
import './css/fonts.css';
import './css/homepage.css';
import './css/livestream.css';
import './css/markdown.css';
import './css/forum/forum.css';
import './css/forum/post-editor.css';
import './css/logs.css';

import '../node_modules/katex/dist/katex.css';

export {
    MBChat,
    toggleLoginPrompt,
    toggleUserList,
    toggleSettings,
    toggleEmotePopup,

    TransientHeader,
    renderMathInElement,
    renderEmotesInElement,
    renderMarkdownInElement,

    LogSearch,

    toggleDarkMode,
    Split,
}
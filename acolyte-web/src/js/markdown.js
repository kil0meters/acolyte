import highlightjs from 'highlight.js';
import marked, {Renderer} from 'marked';

const escapeMap = {
  "&": "&amp;",
  "<": "&lt;",
  ">": "&gt;",
  '"': "&quot;",
  "'": "&#39;"
};

function escapeForHTML(input) {
  return input.replace(/([&<>'"])/g, char => escapeMap[char]);
}

const renderer = new Renderer();
renderer.code = (code, language) => {
  const validLang = !!(language && highlightjs.getLanguage(language));

  const highlighted = validLang
    ? highlightjs.highlight(language, code).value
    : escapeForHTML(code);

  return `<pre><code class="hljs ${language}">${highlighted}</code></pre>`;
};

marked.setOptions({renderer});

export function renderMarkdownInElement(element) {
  console.log(element.innerHTML);
  element.innerHTML = marked(element.innerHTML);
}

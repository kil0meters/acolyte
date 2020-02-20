import marked, {Renderer} from 'marked';

const renderer = new Renderer();
marked.setOptions({renderer});

export function renderMarkdownInElement(element: HTMLElement) {
  element.innerHTML = marked(element.innerHTML.trim().replace(/&gt;+/g, '>'));
}

import marked from 'marked';

const article_regex = /---\n((?:.|\n)*)---\n((?:.|\n)*)/;

// function render_markdown() {
// }

export function renderArticle(content: string): string {
    let matches = article_regex.exec(content);

    let metadata_str = matches[1];
    let content_md = matches[2];

    let metadata: any = {}
    for (let line of metadata_str.split('\n')) {
        let item: string[] = line.split(': ');

        metadata[item[0]] = item[1];
    }

    return `
<article class="wrapper">
  <header class="post-meta card">
    <span class="post-title">${metadata['title']}</span>
    <span class="post-byline">
      ${metadata['date']}
    </span>
  </header>
  ${marked(content_md)}
</article>`;
}

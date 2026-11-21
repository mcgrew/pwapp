
// redefine fetch to handle x-origin
_fetch = fetch
window.fetch = function(url, args) {
    if (/^https?:\/\//i.test(url))
        url = '/_pw/fetch/' + url
    return _fetch(url, args);
}

export class File {
    constructor(filename, mode) {
    }
}

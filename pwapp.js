
export class LocalFile {
    constructor(filename, mode) {
        this.filename = filename;
        this.mode = mode;
    }

    async read() {
        return fetch('/_pwfile/' + encodeURIComponent(this.filename));
    }

    async write(data) {
        return fetch('/_pwfile/' + encodeURIComponent(this.filename), {
            method: "PUT",
            body: data,
        });
    }
    async delete() {
        return fetch('/_pwfile/' + encodeURIComponent(this.filename), {
            method: "DELETE",
        });
    }
    async move(newFilename) {
        return fetch('/_pwfile/' + encodeURIComponent(this.filename), {
            method: "POST",
            headers: {
              "X-Operation-Type": "move",
            },
            body: data,
        });
    }
}

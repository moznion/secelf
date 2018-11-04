export class APIClient {
  static search(query: string, callback: (result: Array<any>, err: Error) => void): void {
    const encodedQuery = encodeURI(query);
    fetch('/api/search?q=' + encodedQuery)
      .then(response => {
        response.json().then(r => callback(r, null));
      })
      .catch(e => {
        callback(null, e);
      });
  }
}

import request from 'superagent';

import {Bookmark} from 'types/model';

import pluginId from './plugin_id';

export default class Client {
    constructor() {
        this.url = `/plugins/${pluginId}/`;

        // good idea, but not setup this way yet
        // this.url = `/plugins/${pluginId}/api/v1`;
    }

    fetchBookmark = async (postID: string) => {
        return this.doGet(`${this.url}/get?postID=${postID}`);
    }

    saveBookmark = async (postID: string, bookmark: Bookmark) => {
        return this.doPost(`${this.url}/add?postID=${postID}`, bookmark);
    }

    doGet = async (url, headers = {}) => {
        headers['X-Requested-With'] = 'XMLHttpRequest';
        headers['X-Timezone-Offset'] = new Date().getTimezoneOffset();

        const response = await request.
            get(url).
            set(headers).
            accept('application/json');

        return response.body;
    }

    doPost = async (url, body, headers = {}) => {
        headers['X-Requested-With'] = 'XMLHttpRequest';
        headers['X-Timezone-Offset'] = new Date().getTimezoneOffset();

        const response = await request.
            post(url).
            send(body).
            set(headers).
            type('application/json').
            accept('application/json');

        return response.body;
    }
}

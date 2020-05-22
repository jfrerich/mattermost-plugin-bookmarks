import {Client4} from 'mattermost-redux/client';
import {ClientError} from 'mattermost-redux/client/client4';

import {Bookmark} from 'types/model';

import pluginId from './plugin_id';

export default class Client {
    constructor() {
        this.url = `/plugins/${pluginId}/api/v1`;
    }

    fetchBookmark = async (postID: string) => {
        return this.doGet(`${this.url}/get?postID=${postID}`);
    }

    saveBookmark = async (bookmark: Bookmark, channelId: string) => {
        return this.doPost(`${this.url}/add`, {bookmark, channelId});
    }

    postEphemeralBookmarks = async (channelId: string) => {
        return this.doPost(`${this.url}/view`, {channelId});
    }

    addLabelByName = async (labelName: string) => {
        return this.doPost(`${this.url}/labels/add?labelName=${labelName}`);
    }

    fetchLabels = async () => {
        return this.doGet(`${this.url}/labels/get`);
    }

    doGet = async (url: string, headers = {}) => {
        headers['X-Timezone-Offset'] = new Date().getTimezoneOffset();

        const options = {
            method: 'get',
            headers,
        };

        const response = await fetch(url, Client4.getOptions(options));

        if (response.ok) {
            return response.json();
        }

        const text = await response.text();

        throw new ClientError(Client4.url, {
            message: text || '',
            status_code: response.status,
            url,
        });
    }

    doPost = async (url: string, body, headers = {}) => {
        headers['X-Timezone-Offset'] = new Date().getTimezoneOffset();

        const options = {
            method: 'post',
            body: JSON.stringify(body),
            headers,
        };

        const response = await fetch(url, Client4.getOptions(options));

        if (response.ok) {
            return response.json();
        }

        const text = await response.text();

        throw new ClientError(Client4.url, {
            message: text || '',
            status_code: response.status,
            url,
        });
    }
}

// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// import {PostTypes} from 'mattermost-redux/action_types';
// import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';

import {Dispatch} from 'redux';

import ActionTypes from 'action_types';
import {Bookmark} from 'types/model';

import Client from 'client';

export function fetchBookmark(postID: string) {
    return async (dispatch: Dispatch) => {
        let data;
        try {
            data = await (new Client()).fetchBookmark(postID);
        } catch (error) {
            return {error};
        }

        dispatch({
            type: ActionTypes.RECEIVED_BOOKMARK,
            data,
        });

        return {data};
    };
}

export function fetchLabels() {
    return async (dispatch: Dispatch) => {
        let data;
        try {
            data = await (new Client()).fetchLabels();
        } catch (error) {
            return {error};
        }

        dispatch({
            type: ActionTypes.RECEIVED_LABELS,
            data,
        });

        return {data};
    };
}

export function saveBookmark(bookmark: Bookmark, channelId: string) {
    let data;
    try {
        data = (new Client()).saveBookmark(bookmark, channelId);
    } catch (error) {
        return {error};
    }

    return {
        type: ActionTypes.CLOSE_ADD_BOOKMARK_MODAL,
        data,
    };
}

export function postEphemeralBookmarks(channelId: string) {
    let data;
    try {
        data = (new Client()).postEphemeralBookmarks(channelId);
    } catch (error) {
        return {error};
    }

    return {
        type: ActionTypes.CLOSE_ADD_BOOKMARK_MODAL,
        data,
    };
}

export function addLabelByName(labelName: string) {
    return async (dispatch: Dispatch) => {
        let data;
        try {
            data = await (new Client()).addLabelByName(labelName);
        } catch (error) {
            return {error};
        }

        dispatch({
            type: ActionTypes.ADDED_LABEL_BY_NAME,
            data,
        });

        return {data};
    };
}

export const openAddBookmarkModal = (postID: string) => {
    return {
        type: ActionTypes.OPEN_ADD_BOOKMARK_MODAL,
        data: {
            postID,
        },
    };
};

export const closeAddBookmarkModal = () => {
    return {
        type: ActionTypes.CLOSE_ADD_BOOKMARK_MODAL,
    };
};

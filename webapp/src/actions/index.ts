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

export function saveBookmark(postID: string, bookmark: Bookmark) {
    let data;
    try {
        data = (new Client()).saveBookmark(postID, bookmark);
    } catch (error) {
        return {error};
    }

    return {data};
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

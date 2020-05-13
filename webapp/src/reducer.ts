import {combineReducers} from 'redux';

import {GenericAction} from 'mattermost-redux/types/actions';

import ActionTypes from './action_types';

const addBookmarksModalVisible = (state = false, action: GenericAction) => {
    switch (action.type) {
    case ActionTypes.OPEN_ADD_BOOKMARK_MODAL:
        return true;
    case ActionTypes.CLOSE_ADD_BOOKMARK_MODAL:
        return false;
    default:
        return state;
    }
};

const addBookmarkModalForPostId = (state = '', action: GenericAction) => {
    switch (action.type) {
    case ActionTypes.OPEN_ADD_BOOKMARK_MODAL:
        return action.data.postID;
    default:
        return state;
    }
};

export default combineReducers({
    addBookmarksModalVisible,
    addBookmarkModalForPostId,
});


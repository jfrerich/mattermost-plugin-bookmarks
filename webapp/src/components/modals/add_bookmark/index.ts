// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';

import {bindActionCreators, Dispatch} from 'redux';

import {getPost} from 'mattermost-redux/selectors/entities/posts';
import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';
import {GlobalState} from 'mattermost-redux/types/store';

import {addBookmarksModalState, getAddBookmarksModalPostId} from 'selectors';
import {
    closeAddBookmarkModal,
    fetchBookmark,
    fetchLabels,
    saveBookmark,
} from 'actions';

import AddBookmarkModal from './add_bookmark';

const mapStateToProps = (state: GlobalState) => {
    const postId = getAddBookmarksModalPostId(state);
    const post = getPost(state, postId);
    const channelId = getCurrentChannelId(state);

    return {
        visible: addBookmarksModalState(state),
        channelId,
        post,
    };
};

const mapDispatchToProps = (dispatch: Dispatch) => bindActionCreators({
    getBookmark: fetchBookmark,
    getAllLabels: fetchLabels,
    close: closeAddBookmarkModal,
    save: saveBookmark,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(AddBookmarkModal);

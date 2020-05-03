// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';

import {bindActionCreators, Dispatch} from 'redux';

import {getPost} from 'mattermost-redux/selectors/entities/posts';
import {GlobalState} from 'mattermost-redux/types/store';
import {Post} from 'mattermost-redux/types/posts';

import {addBookmarksModalState, getAddBookmarksModalPostId} from 'selectors';
import {fetchBookmark, closeAddBookmarkModal} from 'actions';

import AddBookmarkModal from './add_bookmark';

// type Props = {
//     post: Post;
// }

// const mapStateToProps = (state: GlobalState, ownProps: Props) => {
const mapStateToProps = (state: GlobalState) => {
    const postId = getAddBookmarksModalPostId(state);
    const post = getPost(state, postId);

    return {
        visible: addBookmarksModalState(state),
        postId,
        post,
    };
};

const mapDispatchToProps = (dispatch: Dispatch) => bindActionCreators({
    bookmark: fetchBookmark,
    close: closeAddBookmarkModal,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(AddBookmarkModal);

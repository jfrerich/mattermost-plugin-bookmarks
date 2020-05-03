// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators, Dispatch} from 'redux';

import {getPost} from 'mattermost-redux/selectors/entities/posts';

import {openAddBookmarkModal} from 'actions';

import AddBookmarkPostMenuAction from './add_bookmark';

const mapStateToProps = (state, ownProps) => {
    const post = getPost(state, ownProps.postId);
    return {
        post,
    };
};

const mapDispatchToProps = (dispatch: Dispatch) => bindActionCreators({
    open: openAddBookmarkModal,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(AddBookmarkPostMenuAction);

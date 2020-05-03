// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {PureComponent} from 'react';
import {Modal} from 'react-bootstrap';

import {Post} from 'mattermost-redux/types/posts';

import {Bookmark} from 'src/types/model';

import FormButton from 'components/form_button';

export type Props = {
    bookmark: () => void;
    post: Post;
    postId: string;
    visible: boolean;
}

export type State = {
    showModal: boolean;
    submitting: false;
    bookmark: Bookmark;
    fetchError: any;
};

export default class AddBookmarkModal extends PureComponent<Props, State> {
    constructor(props: Props) {
        super(props);

        this.state = {
            showModal: true,
            submitting: false,
        };
    }

    componentDidUpdate(prevProps) {
        if (this.props.post && (!prevProps.post || this.props.post.id !== prevProps.post.id)) {
            const postId = this.props.post.id;
            this.props.bookmark(postId).then((fetched) => {
                // this.setState({bookmark: fetched.error.message, submitting: false});
                this.setState({
                    bookmark: fetched.data,
                    submitting: false}
                );
            });
        }
    }

    handleClose = (e: React.MouseEvent) => {
        if (e && e.preventDefault) {
            e.preventDefault();
        }
        this.setState({showModal: false});
    };

    render() {
        const {showModal, submitting} = this.state;
        const {post} = this.props;

        let postMessageComponent;
        if (post && post.message) {
            const message = post.message;
            postMessageComponent = (
                <div>
                    <h2>
                        {'Post Message'}
                    </h2>
                    {message}
                </div>
            );
        }

        let labelComponent;
        let titleComponent;
        if (this.state && this.state.bookmark) {
            const {bookmark} = this.state;
            if (bookmark.labelIds) {
                const labelMessage = bookmark.labelIds.join();
                labelComponent = (
                    <div>
                        <h2>
                            {'Labels'}
                        </h2>
                        {labelMessage}
                    </div>
                );
            }
            if (bookmark.title) {
                const title = bookmark.title;
                titleComponent = (
                    <div>
                        <h2>
                            {'Title'}
                        </h2>
                        {title}
                    </div>
                );
            }
        }

        return (
            <Modal
                dialogClassName='modal--scroll'
                show={this.props.visible && showModal}
                bsSize='large'
                backdrop='static'
            >
                <Modal.Header closeButton={true}>
                    <Modal.Title>
                        {'Add Your Bookmark'}
                    </Modal.Title>
                </Modal.Header>
                <form
                    role='form'
                    onSubmit={() => null}
                >
                    <Modal.Body
                        ref='modalBody'
                    >
                        {titleComponent}
                        {labelComponent}
                        {postMessageComponent}
                    </Modal.Body>
                    <Modal.Footer >
                        <React.Fragment>
                            <FormButton
                                type='button'
                                btnClass='btn-link'
                                defaultMessage='Cancel'
                                onClick={this.handleClose}
                            />
                            <FormButton
                                id='submit-button'
                                type='submit'
                                btnClass='btn btn-primary'
                                saving={submitting}
                            >
                                {'Create'}
                            </FormButton>
                        </React.Fragment>
                    </Modal.Footer>
                </form>
            </Modal>
        );
    }
}

// const getStyle = (theme) => ({
//     modalBody: {
//         padding: '2em 2em 3em',
//         color: theme.centerChannelColor,
//         backgroundColor: theme.centerChannelBg,
//     },
//     modalFooter: {
//         padding: '2rem 15px',
//     },
//     descriptionArea: {
//         height: 'auto',
//         width: '100%',
//         color: '#000',
//     },
// });

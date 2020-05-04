// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {PureComponent} from 'react';
import {Modal} from 'react-bootstrap';

import {Post} from 'mattermost-redux/types/posts';

import FormButton from 'components/form_button';

export type Props = {
    bookmark: () => void;
    close: () => void;
    save: () => void;
    post: Post;
    postId: string;
    visible: boolean;
}

export type State = {
    showModal: boolean;
    submitting: false;
    bookmark: Bookmark;
    fetchError: any;
    title: string;
    label_ids: string;
};

export default class AddBookmarkModal extends PureComponent<Props, State> {
    constructor(props: Props) {
        super(props);

        this.state = {
            title: '',
            showModal: true,
            submitting: false,
        };
    }

    componentDidUpdate(prevProps) {
        if (this.props.post && (!prevProps.post || this.props.post.id !== prevProps.post.id)) {
            const postId = this.props.post.id;
            this.props.bookmark(postId).then((fetched) => {
                this.setState({
                    bookmark: fetched.data,
                    title: fetched.data.title,
                    label_ids: fetched.data.label_ids,
                    submitting: false}
                );
            });
        }
    }

    handleClose = (e?: Event) => {
        if (e && e.preventDefault) {
            e.preventDefault();
        }
        this.props.close();
    };

    handleSubmit = (e?: Event) => {
        if (e && e.preventDefault) {
            e.preventDefault();
        }

        const timestamp = Date.now();
        const bookmark = {
            postid: this.props.postId,
            title: this.state.title,
            label_ids: this.state.label_ids,
            create_at: timestamp,
            update_at: timestamp,
        };

        this.props.save(bookmark).then((saved) => {
            if (saved.error) {
                this.setState({error: saved.error.message, submitting: false});
            }
        });
        this.props.close();
    };

    handleTitleChange = (e) => {
        this.setState({
            title: e.target.value,
        });
    }

    handleLabelsChange = (e) => {
        this.setState({
            label_ids: e.target.value,
        });
    }

    render() {
        const {showModal, submitting} = this.state;
        const {post} = this.props;
        const style = getStyle();

        let postMessageComponent;
        if (post && post.message) {
            const message = post.message;
            postMessageComponent = (
                <div className='form-group'>
                    <label className='control-label'>{'Post Message'}</label>
                    <textarea
                        style={style.textarea}
                        className='form-control'
                        value={message}
                        resize={'none'}
                        disabled={true}
                    />
                </div>
            );
        }

        const titleComponent = (
            <div className='form-group'>
                <label className='control-label'>{'Title'}</label>
                <input
                    onInput={this.handleTitleChange}
                    className='form-control'
                    value={this.state.title ? this.state.title : ''}
                />
            </div>
        );

        const labelComponent = (
            <div className='form-group'>
                <label className='control-label'>{'Labels'}</label>
                <input
                    onInput={this.handleLabelsChange}
                    className='form-control'
                    value={this.state.label_ids ? this.state.label_ids : []}
                />
            </div>
        );

        // }

        return (
            <Modal
                dialogClassName='modal--scroll'
                show={this.props.visible && showModal}
                bsSize='large'
                backdrop='static'
            >
                <Modal.Header closeButton={true}>
                    <Modal.Title>
                        {'Create or Edit Bookmark'}
                    </Modal.Title>
                </Modal.Header>
                <form
                    role='form'
                    onSubmit={() => null}
                >
                    <Modal.Body ref='modalBody' >
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
                                onClick={this.handleSubmit}
                            >
                                {'Submit'}
                            </FormButton>
                        </React.Fragment>
                    </Modal.Footer>
                </form>
            </Modal>
        );
    }
}
const getStyle = () => ({
    textarea: {
        resize: 'none',
    },
});

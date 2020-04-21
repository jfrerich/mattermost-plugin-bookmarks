// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {PureComponent} from 'react';
import {Modal} from 'react-bootstrap';

import {Post} from 'mattermost-redux/types/posts';

import FormButton from 'components/form_button';

export type Props = {
    visible: boolean;

    post: Post;
    postId: string;
}

export type State = {
    showModal: boolean;
    submitting: false;
};

export default class AddBookmarkModal extends PureComponent<Props, State> {
    constructor(props: Props) {
        super(props);

        this.state = {
            showModal: true,
            submitting: false,
        };
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

        let message = '';
        if (post && post.message) {
            message = post.message;
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
                        {message}
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

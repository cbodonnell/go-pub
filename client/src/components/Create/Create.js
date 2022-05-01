import Modal from 'react-bootstrap/Modal';
import './Create.scss';
import React from 'react';
import { withRouterAndRef } from '../../utils/router';
import { Loading } from '../Loading/Loading';
import { MdAdd, MdEdit, MdUpload } from 'react-icons/md';
import { Tab, Tabs } from 'react-bootstrap';
import CreateNote from '../Note/CreateNote';
import { environment } from '../../environment';
import AuthContext from '../../contexts/AuthContext';
import Publish from '../Publish/Publish';


class Create extends React.Component {

    static contextType = AuthContext;

    constructor(props) {
        super(props);
        this.state = {
            showCreate: false,
            processing: false,
            createError: '',
        }

        this.handleShowCreate = this.handleShowCreate.bind(this);
        this.handleCloseCreate = this.handleCloseCreate.bind(this);

        this.handlePost = this.handlePost.bind(this);
        this.handlePublish = this.handlePublish.bind(this);
    }

    handleCloseCreate = (e) => this.setState({ showCreate: false });

    handleShowCreate(e) { 
        e.preventDefault();
        this.setState({ showCreate: true });
    }

    handlePost(e) {
        this.setState({ showCreate: false })
    };

    handlePublish(e) {
        this.setState({ showCreate: false })
    };

    render() {
        return (
            <>
            {this.context &&
                <div className="flex items-center">
                    <button type="button" title="Create" className="Create-button button button-primary button-icon"
                    onClick={this.handleShowCreate} >
                        <MdAdd />
                    </button>
                    <div className="Create-dialog">
                        <Modal show={this.state.showCreate} onHide={this.handleCloseCreate} centered dialogClassName="Create-dialog">
                            {this.state.processing && 
                                <Loading />
                            }
                            <Modal.Header closeButton>
                                <Modal.Title>Create</Modal.Title>
                            </Modal.Header>
                            <Modal.Body>
                            <div className="w-full h-full relative">
                                <div className="flex flex-col">
                                    <Tabs variant="pills" defaultActiveKey="post" className="Create-tab flex flex-row"
                                    mountOnEnter={true}>
                                        <Tab eventKey="post"
                                        title={
                                        <>
                                            <span className="tab-button-lg">Post</span>
                                            <MdEdit title="Post" className="tab-button-sm" />
                                        </>
                                        }>
                                            <CreateNote to={[`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/followers`]}
                                            onCreate={this.handlePost}/>
                                        </Tab>
                                        <Tab eventKey="publish" 
                                        title={
                                        <>
                                            <span className="tab-button-lg">Publish</span>
                                            <MdUpload title="Publish" className="tab-button-sm" />
                                        </>
                                        }>
                                            <Publish to={[`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/followers`]}
                                            onCreate={this.handlePublish}/>
                                        </Tab>
                                    </Tabs>
                                </div>
                            </div>
                            </Modal.Body>
                            {/* <Modal.Footer>
                            <div className="flex justify-end items-center w-full">
                                <div>
                                    <input type="button" className="button button-secondary" value="Close" onClick={this.handleCloseCreate} />
                                    <input type="submit" className="button button-primary" value="Submit" />
                                </div>
                            </div>
                            </Modal.Footer> */}
                        </Modal>
                    </div>
                </div>
            }
            </>
        );
    }

}

export default withRouterAndRef(Create);
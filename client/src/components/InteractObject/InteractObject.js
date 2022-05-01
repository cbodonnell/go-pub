import './InteractObject.scss';
import React from 'react';
import axios from 'axios';
import logError from '../../utils/errors';
import { apHeaders } from '../../utils/ap';
import CreateNote from '../Note/CreateNote';
import AuthContext from '../../contexts/AuthContext';
import { environment } from '../../environment';
import { Loading } from '../Loading/Loading';
import { MdDelete, MdRecordVoiceOver, MdReply, MdThumbUp } from 'react-icons/md';


export default class InteractObject extends React.Component {
    
    static contextType = AuthContext;

    constructor(props) {
        super(props);
        this.state = { 
            object: null,
            isInteracting: false,
            checkedLiked: false,
            isLiked: false,
            isReplying: false,
            likeIRI: ''
        }

        this.handleDelete = this.handleDelete.bind(this);
        this.toggleInteract = this.toggleInteract.bind(this);
        this.handleCloseInteract = this.handleCloseInteract.bind(this);
        this.checkIfLiked = this.checkIfLiked.bind(this);
        this.handleLike = this.handleLike.bind(this);
        this.handleUnlike = this.handleUnlike.bind(this);
        this.toggleReply = this.toggleReply.bind(this);
        this.handleCreate = this.handleCreate.bind(this);
    }

    componentDidMount() {
        if (this.props.object) {
            this.setState({object: this.props.object});
        }
    }

    async checkIfLiked() {
        this.setState({checkedLiked: false});
        let isLiked = false;
        const result = await axios.get(`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/check?activity=Like&object=${this.state.object.id}`, {
            withCredentials: true
        }).then(res => {
            console.log(res);
            return res.data;
        }).catch(error => {
            logError(error);
        });
        
        if (result.exists) {
            isLiked = true;
            this.setState({likeIRI: result.iri})  
        }
        this.setState({isLiked, checkedLiked: true});
    }

    get isYours() {
        return this.state.object.attributedTo === `${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}`
    }

    handleDelete(e) {
        e.preventDefault();
        if (window.confirm("Are you sure you want to delete this object?")) {
            axios.post(`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/outbox`,
                {
                    "@context": "https://www.w3.org/ns/activitystreams",
                    type: "Delete",
                    actor: `${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}`,
                    object: `${this.state.object.id}`,
                    to: this.state.object.to,
                },
                {
                    headers: {'Content-Type': apHeaders.contentType},
                    withCredentials: true
                }
            ).then(res => {
                console.log(res);
            }).catch(error => {
                logError(error);
            });
        }
    }

    toggleInteract(e) {
        this.setState({ isInteracting: !this.state.isInteracting }, () => {
            if (this.state.isInteracting) {
                this.checkIfLiked();
            }
        });
    };

    handleCloseInteract(e) {
        this.setState({ isInteracting: false });
    };

    handleLike(e) {
        e.preventDefault();
        this.setState({isLiked: true});
        axios.post(`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/outbox`,
            {
                "@context": "https://www.w3.org/ns/activitystreams",
                type: "Like",
                actor: `${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}`,
                object: `${this.state.object.id}`,
                to: `${this.state.object.attributedTo}`,
            },
            {
                headers: {'Content-Type': apHeaders.contentType},
                withCredentials: true
            }
        ).then(res => {
            console.log(res);
        }).catch(error => {
            logError(error);
        });
        // }).finally(() => this.setState({ processing: false }));
    }
    
    handleUnlike(e) {
        e.preventDefault();
        this.setState({isLiked: false});
        axios.post(`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/outbox`,
            {
                "@context": "https://www.w3.org/ns/activitystreams",
                type: "Undo",
                actor: `${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}`,
                object: this.state.likeIRI,
                to: `${this.state.object.attributedTo}`,
            },
            {
                headers: {'Content-Type': apHeaders.contentType},
                withCredentials: true
            }
        ).then(res => {
            console.log(res);
        }).catch(error => {
            logError(error);
        });
    }

    toggleReply(e) {
        this.setState({ isReplying: !this.state.isReplying });
    };

    handleCreate(e) {
        this.setState({ isReplying: false });
    };

    render() {
        return (
        <>
            {this.state.object && 
            <>
                <div className="InteractObject w-full flex justify-between">
                    <div>
                        {this.isYours && this.state.object.type !== 'Tombstone' &&
                            <button type="button" title="Delete"
                            className="ml-0 button button-secondary button-icon" onClick={this.handleDelete} >
                                <MdDelete />
                            </button>
                        }
                    </div>
                    <button type="button" title={!this.state.isInteracting ? 'Interact' : 'Hide Interact'}
                    className={`button ${this.state.isInteracting ? 'button-secondary' : 'button-primary'} button-icon`} onClick={this.toggleInteract} >
                        <MdRecordVoiceOver />
                    </button>
                </div>
                {this.state.isInteracting &&
                <div className="InteractObject-pane">
                    { this.state.checkedLiked &&
                    <div>
                        <div className="w-full flex justify-between">
                            <div>
                                {this.state.checkedLiked && 
                                <>
                                    <div className="w-full flex justify-end">
                                        <button type="button" title={!this.state.isLiked ? 'Like' : 'Unlike'}
                                        className={`ml-0 button ${this.state.isLiked ? 'button-secondary' : 'button-primary'} button-icon`} onClick={!this.state.isLiked ? this.handleLike : this.handleUnlike}>
                                            <MdThumbUp />
                                        </button>
                                    </div>
                                    {/* {!this.state.isLiked &&
                                    }
                                    {this.state.isLiked &&
                                        <div className="w-full flex justify-end">
                                            <input type="button" className="Like-button button button-secondary" value="Unlike" onClick={this.handleUnlike} />
                                        </div>
                                    } */}
                                </>
                                }
                            </div>
                            <div>
                                <div className="w-full flex justify-end">
                                    <button type="button" title={!this.state.isReplying ? 'Reply' : 'Hide Reply'}
                                    className={`button ${this.state.isReplying ? 'button-secondary' : 'button-primary'} button-icon`}  onClick={this.toggleReply}>
                                        <MdReply />
                                    </button>
                                </div>
                            </div>
                        </div>
                        {this.state.isReplying &&
                            <div className="InteractObject-reply w-full flex flex-col">
                                <CreateNote
                                to={this.state.object.attributedTo ? [this.state.object.attributedTo] : []}
                                inReplyTo={this.state.object.id}
                                onCreate={this.handleCreate} />
                            </div>
                        }
                    </div>
                    }
                    { !this.state.checkedLiked &&
                    <div className="InteractObject-loading">
                        <Loading />
                    </div>
                    }
                </div>
                }
            </>
            }
        </>
        );
    }

}
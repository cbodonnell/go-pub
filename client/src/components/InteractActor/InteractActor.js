import './InteractActor.scss';
import React from 'react';
import axios from 'axios';
import logError from '../../utils/errors';
import { apHeaders } from '../../utils/ap';
import AuthContext from '../../contexts/AuthContext';
import { environment } from '../../environment';
import { Loading } from '../Loading/Loading';
import { AiOutlineUserAdd, AiOutlineUserDelete } from 'react-icons/ai';
import { MdRecordVoiceOver } from 'react-icons/md';
import { MdEdit } from 'react-icons/md';
import CreateNote from '../Note/CreateNote';


export default class InteractActor extends React.Component {
    
    static contextType = AuthContext;

    constructor(props) {
        super(props);
        this.state = {
            actor: null,
            loading: false,
            isInteracting: false,
            checkedFollowing: false,
            isFollowing: false,
            followIRI: '',
            isDMing: false
        };

        this.toggleInteract = this.toggleInteract.bind(this);
        this.toggleDM = this.toggleDM.bind(this);
        this.handleFollow = this.handleFollow.bind(this);
        this.handleUnfollow = this.handleUnfollow.bind(this);
        this.handleCreate = this.handleCreate.bind(this);
    }

    componentDidMount() {
        if (this.props.actor) {
            this.setState({actor: this.props.actor});
        }
    }

    toggleInteract() {
        this.setState({isInteracting: !this.state.isInteracting}, () => {
            if (this.state.isInteracting) {
                this.checkIfFollowing()
            }
        });
    }

    toggleDM() {
        this.setState({isDMing: !this.state.isDMing});
    }
    
    async checkIfFollowing() {
        let isFollowing = false;
        if (this.isNotSelf) {
            const result = await axios.get(`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/check?activity=Follow&object=${this.state.actor.id}`, {
                withCredentials: true
            }).then(res => {
                console.log(res);
                return res.data;
            }).catch(error => {
                logError(error);
            });
            
            if (result.exists) {
                isFollowing = true;
                this.setState({followIRI: result.iri})  
            }
        }
        this.setState({isFollowing, checkedFollowing: true});
    }

    get isNotSelf() {
        return this.state.actor.id !== `${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}`
    }

    handleFollow(e) {
        e.preventDefault();
        this.setState({isFollowing: true});
        axios.post(`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/outbox`,
            {
                "@context": "https://www.w3.org/ns/activitystreams",
                type: "Follow",
                actor: `${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}`,
                object: `${this.state.actor.id}`,
                to: `${this.state.actor.id}`,
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
    
    handleUnfollow(e) {
        e.preventDefault();
        this.setState({isFollowing: false});
        axios.post(`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/outbox`,
            {
                "@context": "https://www.w3.org/ns/activitystreams",
                type: "Undo",
                actor: `${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}`,
                object: this.state.followIRI,
                to: `${this.state.actor.id}`,
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

    handleCreate() {
        this.setState({ isDMing: false });
    }

    render() {
        return (
        <>
            {this.state.actor && this.isNotSelf &&
            <div className="InteractActor">
                <div className="w-full flex justify-end">
                    <button type="button" title={!this.state.isInteracting ? 'Interact' : 'Hide Interact'}
                    className={`button ${this.state.isInteracting ? 'button-secondary' : 'button-primary'} button-icon`} onClick={this.toggleInteract} >
                        <MdRecordVoiceOver />
                    </button>
                </div>
                {this.state.isInteracting &&
                <div className="InteractActor-pane">
                    {this.state.checkedFollowing &&
                    <>
                        <div className="InteractActor-buttons">
                            <button type="button" title={!this.state.isFollowing ? 'Follow' : 'Unfollow'}
                            className={`ml-0 button ${this.state.isFollowing ? 'button-secondary' : 'button-primary'} button-icon`}
                            onClick={!this.state.isFollowing ? this.handleFollow : this.handleUnfollow}>
                                {!this.state.isFollowing ? <AiOutlineUserAdd /> : <AiOutlineUserDelete />}
                            </button>
                            <button type="button" title={!this.state.isDMing ? 'Message' : 'Hide Message'}
                            className={`button ${this.state.isDMing ? 'button-secondary' : 'button-primary'} button-icon`}
                            onClick={this.toggleDM}>
                                <MdEdit />
                            </button>
                        </div>
                        {this.state.isDMing &&
                            <div className="InteractActor-reply w-full flex flex-col">
                                <CreateNote to={[this.state.actor.id]} onCreate={this.handleCreate} />
                            </div>
                        }
                    </>
                    }
                    { !this.state.checkedFollowing &&
                    <div className="InteractActor-loading">
                        <Loading />
                    </div>
                    }
                </div>}
            </div>
            }
        </>
        );
    }

}
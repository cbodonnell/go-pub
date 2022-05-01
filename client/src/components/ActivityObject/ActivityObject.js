import './ActivityObject.scss';
import React from 'react';
import logError from '../../utils/errors';
import AuthContext from '../../contexts/AuthContext';
import Note from '../Note/Note';
import { apHeaders, renderIcon } from '../../utils/ap';
import InteractObject from '../InteractObject/InteractObject';
import { Loading } from '../Loading/Loading';
import IRI from '../IRI/IRI';
import { wrapIfRemote } from '../../utils/urls';
import { Divider } from '../Divider/Divider';
import { BiConversation } from 'react-icons/bi';
import { proxyCacheClient } from '../../utils/http';
import Audio from '../Audio/Audio';
import Replies from '../Replies/Replies';


export default class ActivityObject extends React.Component {
    
    static contextType = AuthContext;

    constructor(props) {
        super(props);
        this.state = {
            object: null,
            loading: false,
            attributedTo: [],
            unauthorized: false,
            notfound: false,
            showInReplyTo: false,
            replies: [],
            replyLevel: 1,
            maxReplies: 3,
        };
        
        this._isMounted = false;

        this.handleInteract = this.handleInteract.bind(this);
        this.handleCloseInteract = this.handleCloseInteract.bind(this);
        
        this.toggleShowInReplyTo = this.toggleShowInReplyTo.bind(this);
    }

    componentDidMount() {
        this._isMounted = true;
        if (this.props.object) {
            this.setState({ object: this.props.object })
            this.setObject(this.props.object);
        } else if (this.props.iri) {
            this.fetchObject(this.props.iri);
        }

        if (this.props.replyLevel) {
            this.setState({replyLevel: this.props.replyLevel});
        }
        
        if (!this.props.isReferenced) {
            this.setState({maxReplies: 10});
        } else if (this.props.maxReplies) {
            this.setState({maxReplies: this.props.maxReplies});
        }
    }

    componentWillUnmount() {
        this._isMounted = false;
    }

    setObject(object) {
        if (this._isMounted) {
            this.setState({ object }, () => this.getAttributedTo());
        }
    }

    fetchObject(iri) {
        this.setState({loading: true});
        proxyCacheClient.get(iri, {
            headers: {'Accept': apHeaders.accept},
        } ).then(res => {
            console.log(res);
            const object = res.data;
            this.setObject(object);
        }).catch(error => {
            this.setState({loading: false})
            if (error.response && (error.response.status === 401 || error.response.status === 403)) {
                this.setState({unauthorized: true});
            } else {
                this.setState({notfound: true});
            }
            logError(error);
        });
    }

    getAttributedTo() {
        this.setState({loading: true});
        if (typeof this.state.object.attributedTo !== 'string') {
            if (Array.isArray(this.state.object.attributedTo)) {
                this.state.object.attributedTo.forEach(author => {
                    this.fetchAttributedTo(author.id)
                    .finally(() => this.setState({loading: false}));
                });
            } else {
                this.setState({loading: false});
            }
        } else {
            this.fetchAttributedTo(this.state.object.attributedTo)
            .catch(error => {
                if (error.response && (error.response.status === 401 || error.response.status === 403)) {
                    this.setState({unauthorized: true});
                } else {
                    this.setState({notfound: true});
                }
            }).finally(() => this.setState({loading: false}));
        }
    }

    fetchAttributedTo(iri) {
        return proxyCacheClient.get(iri, {
            headers: {'Accept': apHeaders.accept},
        }).then(res => {
            const actor = res.data;
            const attributedTo = this.state.attributedTo;
            attributedTo.push(actor);
            this.setState({ attributedTo });
        }).catch(error => {
            logError(error);
        });
    }

    renderAttributedTo() { 
        return <span>
            {this.state.attributedTo.map((attributedTo, i) => {
                return <span key={`${this.state.object.id}#attributedTo-${i}`}>
                    {i > 0 ?
                    <span>
                        , <IRI url={wrapIfRemote(attributedTo.id, 'actor')}
                        text={attributedTo.preferredUsername || attributedTo.name} />
                    </span> :
                    <IRI url={wrapIfRemote(attributedTo.id, 'actor')}
                    text={attributedTo.preferredUsername || attributedTo.name} />}
                </span>
            })}
        </span>
    }

    renderObject() {
        switch(this.state.object.type){
            case 'Note': return <Note {...{...this.state.object, ...this.props}} />;
            case 'Audio': return <Audio {...{...this.state.object, ...this.props}} />;
            case 'Tombstone': return <p>This object has been deleted.</p>;
            default: return <IRI url={this.props.isReferenced ? wrapIfRemote(this.state.object.id, 'object'): this.state.object.id}
            text={this.state.object.name || this.state.object.id} />;
        }
    }

    renderAttachment() {
        if (Array.isArray(this.state.object.attachment)) {
            return <div className="ActivityObject-attachment mt-1">
                {this.state.object.attachment.map((item, i) => {
                    return <div className="ActivityObject-attachment-item flex justify-center"
                    key={`${this.state.object.id}#attachment-${i}`}>
                        {this.renderAttachmentItem(item)}
                    </div>
                })}
                {/* TODO: Add next & back buttons for pagination through attachments */}
            </div>;
        }
        return;
    }

    renderAttachmentItem(item) {
        switch (item.mediaType) {
            case "image/jpeg":
            case "image/png":
                return <a href={item.url} target="_blank" rel="noreferrer">
                    <img className="ActivityObject-attachment-image w-full" src={item.url} alt="..." />
                </a>;
            default:
                return
        }
    }

    toggleShowInReplyTo(e) {
        e.preventDefault();
        this.setState({showInReplyTo: !this.state.showInReplyTo});
    }

    get canGoDeeper() {
        return this.state.replyLevel <= this.state.maxReplies;
    }

    handleInteract(e) {
        this.setState({ isInteracting: true });
    };

    handleCloseInteract(e) {
        this.setState({ isInteracting: false });
    };

    render() {
        return (
            <>
                { !this.state.loading &&
                <>
                    {this.state.object && 
                    <div className="ActivityObject">
                        { this.state.showInReplyTo && this.canGoDeeper &&
                            <div>
                                {/* <Divider /> */}
                                <ActivityObject key={`${this.state.object.id}#reply-${this.state.replyLevel}`}
                                iri={this.state.object.inReplyTo}
                                isReferenced={true}
                                replyLevel={this.state.replyLevel + 1}
                                maxReplies={this.state.maxReplies} />
                                <Divider />
                            </div>
                        }
                        <div className="flex justify-between">
                            {this.state.attributedTo.map((actor, i) => {
                                return <span key={`${this.state.object.id}#attributedToIcon-${i}`}>
                                    {renderIcon(actor.icon)}
                                </span>
                            })}
                            {/* {renderIcon(this.state.attributedTo.length > 0 ? this.state.attributedTo[0].icon : null)} */}
                            <div className="ActivityObject-actor flex-auto overflow-hidden">
                                <div className="w-full overflow-hidden">
                                    <p className="ActivityObject-actor-name overflow-hidden">
                                        {this.state.object.attributedTo && this.renderAttributedTo()}
                                    </p>
                                    <p className="m-0">
                                        <IRI url={this.props.isReferenced ? wrapIfRemote(this.state.object.id, 'object'): this.state.object.id}
                                        text={this.state.object.type} />
                                    </p>
                                </div>
                            </div>
                            <div className="ml-1">
                                {!this.props.hideInReplyTo && this.state.object.inReplyTo &&
                                <>
                                    { this.canGoDeeper &&
                                        <button type="button" title={!this.state.showInReplyTo ? 'Show In Reply To' : 'Hide In Reply To'}
                                        className={`m-0 button ${this.state.showInReplyTo ? 'button-secondary' : 'button-tertiary'} button-icon`}
                                        onClick={this.toggleShowInReplyTo} >
                                            <BiConversation />
                                        </button>
                                    }
                                    { !this.canGoDeeper &&
                                        <p className="ActivityObject-view-more">
                                            <IRI url={wrapIfRemote(this.state.object.id, 'object')}
                                            text={"more..."} />
                                        </p>
                                    }
                                </>
                                }
                            </div>
                        </div>
                        {this.renderObject()}
                        {this.state.object.attachment && this.renderAttachment()}
                        {this.context && <InteractObject object={this.state.object} />}
                        {!this.props.isReferenced && this.state.object.replies && <Replies object={this.state.object} />}
                    </div>
                    }
                    { this.state.unauthorized &&
                    <>
                        { this.props.remote && !this.context ? <p>Sign in to browse the fediverse!</p> : <p>Unauthorized.</p> }
                    </>
                    }
                    { this.state.notfound && <p>Not found.</p> }
                </>
                }
                { this.state.loading && <Loading /> }
            </>
        );
    }

}
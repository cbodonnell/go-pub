import './Actor.scss';
import React from 'react';
// import axios from 'axios';
import logError from '../../utils/errors';
import IRI from '../IRI/IRI';
import AuthContext from '../../contexts/AuthContext';
import { apHeaders, renderIcon } from '../../utils/ap';
import { wrapIfRemote } from '../../utils/urls';
import { Tabs, Tab } from 'react-bootstrap';
import OrderedCollection from '../OrderedCollection/OrderedCollection';
import { Divider } from '../Divider/Divider';
import InteractActor from '../InteractActor/InteractActor';
import { Loading } from '../Loading/Loading';
import SanitizedHTML from 'react-sanitized-html';
import { FcCollapse, FcExpand } from 'react-icons/fc';
import { proxyCacheClient } from '../../utils/http';
import { MdInbox, MdOutbox, MdThumbUp } from 'react-icons/md';
import { IoMdPeople, IoMdPersonAdd } from 'react-icons/io';

export default class Actor extends React.Component {

    static contextType = AuthContext;

    constructor(props) {
        super(props);
        this.state = {
            actor: null,
            loading: false,
            isCollapsed: true,
            unauthorized: false,
            notfound: false
        };
        
        this._isMounted = false;

        this.toggleCollapse = this.toggleCollapse.bind(this);
    }

    componentDidMount() {
        this._isMounted = true;
        if (this.props.actor) {
            this.setActor(this.props.actor);
        } else if (this.props.iri) {
            this.fetchActor(this.props.iri);
        } else {
            this.setState({notfound: true});
        }
    }

    componentWillUnmount() {
        this._isMounted = false;
    }

    setActor(actor) {
        if (this._isMounted) {
            this.setState({actor});
        }
    }

    fetchActor() {
        this.setState({loading: true});
        // axios.get(this.props.iri, {
        proxyCacheClient.get(this.props.iri, {
            headers: {'Accept': apHeaders.accept},
        }).then(res => {
            console.log(res);
            const actor = res.data;
            this.setActor(actor);
        }).catch(error => {
            if (error.response && (error.response.status === 401 || error.response.status === 403)) {
                this.setState({unauthorized: true});
            } else {
                this.setState({notfound: true});
            }
            logError(error);
        }).finally(() => this.setState({loading: false}));
    }

    toggleCollapse() {
        this.setState({isCollapsed: !this.state.isCollapsed});
    }

    render() {
        return (
            <>
                { !this.state.loading &&
                <>
                    {this.state.actor &&
                        <>
                            <div className="Actor">
                                <div className="flex w-full items-start">
                                    {renderIcon(this.state.actor.icon)}
                                    <div className="Actor-name-subtitle w-full">
                                        {this.props.small ?
                                        <p className="Actor-name-small">{this.state.actor.name || this.state.actor.id}</p>:
                                        <h5>{this.state.actor.name || this.state.actor.id}</h5>}
                                        <div>
                                            <IRI url={this.props.isReferenced ? wrapIfRemote(this.state.actor.id, 'actor') : this.state.actor.id}
                                            text={this.state.actor.preferredUsername || this.state.actor.id} />
                                        </div>
                                    </div>
                                    { !this.props.small &&
                                        <button type="button" title={!this.state.isCollapsed ? 'Collapse' : 'Expand'}
                                        className={`m-0 button ${this.state.isCollapsed ? 'button-tertiary' : 'button-secondary'} button-icon`} onClick={this.toggleCollapse} >
                                            {!this.state.isCollapsed ? <FcCollapse /> : <FcExpand />}
                                        </button>
                                    }
                                </div>
                            </div>
                            {!this.state.isCollapsed &&
                            <div className="Actor-collapsible">
                                {this.state.actor.summary && 
                                <div className="Actor-summary overflow-hidden">
                                    <SanitizedHTML
                                    allowedAttributes={{ 'a': ['href', 'target', 'ref'] }}
                                    allowedTags={['a', 'p', 'br']}
                                    transformTags={{
                                        'a': (tagName, attribs) => ({ tagName, attribs: {...attribs, target: '_blank', ref: 'noreferrer'}})
                                    }}
                                    html={this.state.actor.summary} />
                                </div>
                                }
                                {this.context && <InteractActor actor={this.state.actor} />}
                            </div>
                            }
                            { !this.props.small &&
                            <>
                                <Divider />
                                <div className="Actor-collections">
                                    <Tabs variant="pills" defaultActiveKey="outbox" className="Feed-tab"
                                    mountOnEnter={true} unmountOnExit={true}>
                                        {this.state.actor.outbox && 
                                            <Tab eventKey="outbox" 
                                            title={
                                            <>
                                                <span className="tab-button-lg">Outbox</span>
                                                <MdOutbox title="Outbox" className="tab-button-sm" />
                                            </>
                                            }>
                                                <OrderedCollection
                                                title={'Outbox'}
                                                type={'Activity'}
                                                link={true}
                                                iri={this.state.actor.outbox} />
                                            </Tab>
                                        }
                                        {this.state.actor.inbox && 
                                            <Tab eventKey="inbox"
                                            title={
                                            <>
                                                <span className="tab-button-lg">Inbox</span>
                                                <MdInbox title="Inbox" className="tab-button-sm" />
                                            </>
                                            }>
                                                <OrderedCollection
                                                title={'Inbox'}
                                                type={'Activity'}
                                                link={true}
                                                iri={this.state.actor.inbox} />
                                            </Tab>
                                        }
                                        {this.state.actor.followers && 
                                            <Tab eventKey="followers"
                                            title={
                                            <>
                                                <span className="tab-button-lg">Followers</span>
                                                <IoMdPeople title="Followers" className="tab-button-sm" />
                                            </>
                                            }>
                                                <OrderedCollection
                                                title={'Followers'}
                                                type={'Actor'}
                                                link={true}
                                                iri={this.state.actor.followers} />
                                            </Tab>
                                        }
                                        {this.state.actor.following && 
                                            <Tab eventKey="following"
                                            title={
                                            <>
                                                <span className="tab-button-lg">Following</span>
                                                <IoMdPersonAdd title="Following" className="tab-button-sm" />
                                            </>
                                            }>
                                                <OrderedCollection
                                                title={'Following'}
                                                type={'Actor'}
                                                link={true}
                                                iri={this.state.actor.following} />
                                            </Tab>
                                        }
                                        {this.state.actor.liked && 
                                            <Tab eventKey="liked"
                                            title={
                                            <>
                                                <span className="tab-button-lg">Liked</span>
                                                <MdThumbUp title="Liked" className="tab-button-sm" />
                                            </>
                                            }>
                                                <OrderedCollection
                                                title={'Liked'}
                                                type={'Object'}
                                                link={true}
                                                iri={this.state.actor.liked} />
                                            </Tab>
                                        }
                                    </Tabs>
                                </div>
                            </>
                            }
                        </>
                    }
                    {this.state.unauthorized && 
                    <>
                        { this.props.remote && !this.context ? <p>Sign in to browse the fediverse!</p> : <p>Unauthorized.</p> }
                    </>
                    }
                    {this.state.notfound && <p>Not found.</p> }
                </>
                }
                { this.state.loading && <div className={this.props.small ? "Actor-loading": ""}>
                    <Loading />
                </div> }
            </>
        );
    }

}
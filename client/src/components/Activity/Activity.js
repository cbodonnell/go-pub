import './Activity.scss';
import React from 'react';
import IRI from '../IRI/IRI';
// import axios from 'axios';
import logError from '../../utils/errors';
import ActivityObject from '../ActivityObject/ActivityObject';
import { apHeaders, isValidActivity, renderIcon } from '../../utils/ap';
import { wrapIfRemote } from '../../utils/urls';
import { Loading } from '../Loading/Loading';
import AuthContext from '../../contexts/AuthContext';
import { proxyCacheClient } from '../../utils/http';
import Actor from '../Actor/Actor';


export default class Activity extends React.Component {

    static contextType = AuthContext;

    constructor(props) {
        super(props);
        this.state = {
            activity: null,
            loading: false,
            actor: null,
            unauthorized: false,
            notfound: false
        };
        
        this._isMounted = false;
    }

    componentDidMount() {
        this._isMounted = true;
        if (this.props.activity) {
            this.setActivity(this.props.activity);
        } else if (this.props.iri) {
            this.fetchActivity(this.props.iri);
        }
    }

    componentWillUnmount() {
        this._isMounted = false;
    }

    setActivity(activity) {
        if (this._isMounted) {
            this.setState({ activity }, () => this.fetchActor());
        }
    }

    fetchActivity(iri) {
        this.setState({loading: true});
        // axios.get(iri, {
        proxyCacheClient.get(iri, {
            headers: {'Accept': apHeaders.accept},
        } ).then(res => {
            console.log(res);
            const activity = res.data;
            if (!isValidActivity(activity)) {
                throw new Error({response: { status: 404}});
            } else {
                this.setActivity(activity);
            }
        }).catch(error => {
            this.setState({loading: false});
            if (error.response && (error.response.status === 401 || error.response.status === 403)) {
                this.setState({unauthorized: true});
            } else {
                this.setState({notfound: true});
            }
            logError(error);
        });
    }

    fetchActor() {
        this.setState({loading: true});
        const actorIRI = typeof this.state.activity.actor === 'string' ?
        this.state.activity.actor :
        this.state.activity.actor.id;

        // axios.get(timestamp(actorIRI), {
        proxyCacheClient.get(actorIRI, {
            headers: {'Accept': apHeaders.accept},
        }).then(res => {
            console.log(res);
            const actor = res.data;
            this.setState({ actor });
            // this.fetchActor();
        }).catch(error => {
            if (error.response && (error.response.status === 401 || error.response.status === 403)) {
                this.setState({unauthorized: true});
            } else {
                this.setState({notfound: true});
            }
            logError(error);
        }).finally(() => this.setState({loading: false}));
    }

    getObjectIRI(object, type) {
        const objectIRI = typeof object === 'string' ? object : object.id;
        return <IRI url={wrapIfRemote(objectIRI, type)} text={objectIRI} />;
    }

    renderActivityObject() {
        if (this.props.isChild)  {
            // TODO: May need to revisit this based on other child Activity types
            const type = this.state.activity.type === 'Follow' ? 'actor' : 'object';
            return <p>Object: {this.getObjectIRI(this.state.activity.object, type)}</p>;
        }
        switch(this.state.activity.type){
            case 'Create':
            case 'Announce':
            case "Like":
            case 'Delete':
                return <>
                    <p>Object:</p>
                    <div className="Activity-object">
                        {typeof this.state.activity.object === 'string' ?
                        <ActivityObject iri={this.state.activity.object} isReferenced={true} />:
                        <ActivityObject object={this.state.activity.object} isReferenced={true} />}
                    </div>
                </>;
            // case "Delete":
            //     return <p>Object: {this.getObjectIRI(this.state.activity.object, 'object')}</p>;
            case "Follow":
                return <>
                    <p>Object:</p>
                    <div className="Activity-object">
                        {typeof this.state.activity.object === 'string' ? 
                        <Actor iri={this.state.activity.object} small={true} isReferenced={true} /> :
                        <Actor actor={this.state.activity.object} small={true} isReferenced={true} />}
                    </div>
                </>;
            case "Undo":
            case "Accept":
                return <>
                    <p>Object:</p>
                    <div className="Activity-object">
                        {typeof this.state.activity.object === 'string' ?
                        <Activity iri={this.state.activity.object} isReferenced={true} isChild={true} />:
                        <Activity object={this.state.activity.object} isReferenced={true} isChild={true} />}
                    </div>
                </>;
            default: return <p>Unsupported Activity type: {this.state.activity.type}</p>
        }
    }

    render() {
        return (
            <>
                { !this.state.loading &&
                <>
                    {this.state.activity &&
                        <div className="Activity">
                            <div className="Activity-actor">
                                {this.state.actor && renderIcon(this.state.actor.icon)}
                                <div className="w-full overflow-hidden">
                                    <p className="Activity-actor-name overflow-hidden">
                                        <IRI url={wrapIfRemote(this.state.activity.actor, 'actor')}
                                        text={this.state.actor ? 
                                            this.state.actor.preferredUsername || this.state.actor.name :
                                            this.state.activity.actor} />
                                    </p>
                                    <p className="m-0">
                                        <IRI url={this.props.isReferenced ? wrapIfRemote(this.state.activity.id, 'activity') : this.state.activity.id}
                                        text={this.state.activity.type} />
                                    </p>
                                </div>
                            </div>
                            {this.state.activity.object ? this.renderActivityObject() : <p>No Object.</p>}
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
                { this.state.loading && <Loading />}
            </>
        );
    }

}
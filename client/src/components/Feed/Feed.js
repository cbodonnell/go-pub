import './Feed.scss';
import React from 'react';
import Tabs from 'react-bootstrap/Tabs'
import Tab from 'react-bootstrap/Tab'
import OrderedCollection from '../OrderedCollection/OrderedCollection';
import AuthContext from '../../contexts/AuthContext';
// import CreateNote from '../Note/CreateNote';
import { environment } from '../../environment';
import { MdInbox , MdOutbox, MdThumbUp} from 'react-icons/md';
import { IoMdPeople, IoMdPersonAdd } from 'react-icons/io';


export default class Feed extends React.Component {

    static contextType = AuthContext;

    constructor(props) {
        super(props);
        this.state = {
            isCreating: false,
            outbox: null,
            inbox: null,
            followers: null,
        }
        
        this.outboxRef = React.createRef();
        
        // this.toggleCreate = this.toggleCreate.bind(this);
        // this.handleCreate = this.handleCreate.bind(this);
        // this.handleCancel = this.handleCancel.bind(this);
    }

    // toggleCreate() {
    //     this.setState({ isCreating: !this.state.isCreating });
    // }

    // handleCreate() {
    //     this.setState({ isCreating: false });
    //     this.outboxRef.current.refreshOrderedCollection();
    // }

    // handleCancel() { }

    clearForm() {
        document.getElementById("note-form").reset();
    }

    render() {
        return (
            <>
                {this.context && 
                    <div className="Feed">
                        {/* <div className="Feed-welcome">
                            <div className="Feed-welcome-titleblock flex justify-between items-center">
                                <h5 className="m-0">Welcome, {this.context.username}!</h5>
                                <button type="button" title={!this.state.isCreating ? 'Post' : 'Hide Post'}
                                className={`m-0 button ${this.state.isCreating ? 'button-secondary' : 'button-primary'} button-icon`} onClick={this.toggleCreate} >
                                    <MdEdit />
                                </button>
                            </div>
                            {this.state.isCreating &&
                                <CreateNote to={[`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/followers`]}
                                onCreate={this.handleCreate}
                                onCancel={this.handleCancel} />
                            }
                        </div> */}
                        <div className="Feed-collections">
                            {/* <Divider /> */}
                            <Tabs variant="pills" defaultActiveKey="inbox" className="Feed-tab"
                            mountOnEnter={true} unmountOnExit={true}>
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
                                    iri={`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/inbox`} />
                                </Tab>
                                <Tab eventKey="outbox" 
                                title={
                                <>
                                    <span className="tab-button-lg">Outbox</span>
                                    <MdOutbox title="Outbox" className="tab-button-sm" />
                                </>
                                }>
                                    <OrderedCollection
                                    ref={this.outboxRef}
                                    title={'Outbox'}
                                    type={'Activity'}
                                    link={true}
                                    iri={`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/outbox`} />
                                </Tab>
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
                                    iri={`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/followers`} />
                                </Tab>
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
                                    iri={`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/following`} />
                                </Tab>
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
                                    iri={`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/liked`} />
                                </Tab>
                            </Tabs>
                        </div>
                    </div>
                }
                { !this.context && 
                    <p>You are not signed in.</p>
                }
            </>
        );
    }

}
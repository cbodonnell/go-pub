import './Replies.scss';
import React from 'react';
import { fetchCollectionWithItems } from '../../utils/ap';
import { Divider } from '../Divider/Divider';
import ActivityObject from '../ActivityObject/ActivityObject';
import { BiConversation } from 'react-icons/bi';
import { Loading } from '../Loading/Loading';


export default class Replies extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            object: null,
            loading: false,
            showReplies: false,
            replies: []
        };
        
        this.toggleShowReplies = this.toggleShowReplies.bind(this);
    }

    componentDidMount() {
        if (this.props.object) {
            this.setState({object: this.props.object}, () => {
                this.getReplies();
            });
        }
    }

    toggleShowReplies(e) {
        this.setState({ showReplies: !this.state.showReplies });
    };

    async getReplies() {
        this.setState({loading: true});
        const iri = typeof this.state.object.replies !== 'string' ?
            this.state.object.replies.id :
            this.state.object.replies;
        const collection = await fetchCollectionWithItems(iri);
        this.setState({replies: collection.orderedItems, loading: false});
    }

    render() {
        return (
            <>
                {this.state.replies.length > 0 &&
                    <div className="Replies">
                        <div className="flex justify-end">
                            <button type="button" title={!this.state.showReplies ? 'Show Replies' : 'Hide Replies'}
                            className={`button ${this.state.showReplies ? 'button-secondary' : 'button-tertiary'} button-icon`} onClick={this.toggleShowReplies} >
                                <BiConversation />
                            </button>
                        </div>
                        {this.state.showReplies &&
                            <div className="Replies-collection">
                                {this.state.replies.map((item, _) => <>
                                    <Divider />
                                    {typeof item === 'string' ?
                                    <ActivityObject iri={item} isReferenced={true} hideInReplyTo={true} /> :
                                    <ActivityObject object={item} isReferenced={true} hideInReplyTo={true} />}
                                </>
                                )}
                            </div>
                        }
                    </div>
                }
                {this.state.loading && 
                    <div className="Replies-loading">
                        <Loading />
                    </div>
                }
            </>
        );
    }

}
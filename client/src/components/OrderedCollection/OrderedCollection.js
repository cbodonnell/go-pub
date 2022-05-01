import './OrderedCollection.scss';
import React from 'react';
import logError from '../../utils/errors';
import Activity from '../Activity/Activity';
import { Divider } from '../Divider/Divider';
import { apHeaders } from '../../utils/ap';
import IRI from '../IRI/IRI';
import ActivityObject from '../ActivityObject/ActivityObject';
import { Loading } from '../Loading/Loading';
import Actor from '../Actor/Actor';
import { wrapIfRemote } from '../../utils/urls';
import { GrRefresh } from "react-icons/gr";
import AuthContext from '../../contexts/AuthContext';
import { MdFirstPage, MdNavigateNext, MdNavigateBefore, MdLastPage } from 'react-icons/md';
import InfiniteScroll from 'react-infinite-scroll-component';
import { proxyClient } from '../../utils/http';


export default class OrderedCollection extends React.Component {

    static contextType = AuthContext;

    constructor(props) {
        super(props);
        this.state = {
            orderedCollection: null,
            loading: false,
            unauthorized: false,
            notfound: false
        };

        this.handleRefresh = this.handleRefresh.bind(this);
        this.handleNext = this.handleNext.bind(this);
        this.handleFirst = this.handleFirst.bind(this);
        this.handleLast = this.handleLast.bind(this);
        this.handlePrevious = this.handlePrevious.bind(this);
    }

    componentDidMount() {
        if (this.props.orderedCollection) {
            console.log(this.props.orderedCollection);
            this.setState(
                { orderedCollection: this.props.orderedCollection },
                () => this.getOrderedItems()
            );
        } else if (this.props.iri) {
            this.fetchOrderedCollection(this.props.iri);
        }
    }

    componentDidUpdate(oldProps) {
        if (this.props.orderedCollection !== oldProps.orderedCollection) {
            console.log(this.props.orderedCollection);
            this.setState(
                { orderedCollection: this.props.orderedCollection },
                () => this.getOrderedItems()
            );
        }
    }

    fetchOrderedCollection(iri) {
        this.setState({loading: true});
        proxyClient.get(iri, {
            headers: {'Accept': apHeaders.accept},
            withCredentials: true
        }).then(res => {
            console.log(iri, res);
            const orderedCollection = res.data;
            if (Object.keys(orderedCollection).length === 0) {
                throw new Error({response: { status: 404}});
            } else {
                this.setState({ orderedCollection }, () => this.getOrderedItems());
            }
            // this.getOrderedItems();
        }).catch(error => {
            this.setState({loading: false});
            if (!error.response || error.response.status === 404) {
                this.setState({notfound: true});
            } else if (error.response.status === 401 || error.response.status === 403) {
                this.setState({unauthorized: true});
            }
            logError(error);
        });
    }

    getOrderedItems() {
        const page = "first";
        if (this.state.orderedCollection[page]) {
            this.setPage("first");
        } else {
            this.setState({loading: false});
        }
    }

    setPage(page) {
        this.setState({loading: true});
        const orderedCollection = this.state.orderedCollection;
        orderedCollection.orderedItems = [];
        this.setState({ orderedCollection }, () => {
            this.fetchPage(page).finally(() => this.setState({loading: false}));
        });
    }

    fetchPage(page, infiniteScroll=false) {
        const iri = typeof this.state.orderedCollection[page] !== 'string' ?
            this.state.orderedCollection[page].id :
            this.state.orderedCollection[page];
        return proxyClient.get(`${iri}`, {
            headers: {'Accept': apHeaders.accept},
            withCredentials: true
        } ).then(res => {
            console.log(res);
            const page = res.data;
            const orderedCollection = this.state.orderedCollection;
            if (page.orderedItems || page.items) {
                orderedCollection.orderedItems = orderedCollection.orderedItems.concat(page.orderedItems || page.items);
            }
            if (page.first) {
                orderedCollection.first = page.first;
            }
            if (page.last) {
                orderedCollection.last = page.last;
            }
            orderedCollection.next = page.next;
            if (!infiniteScroll) {
                orderedCollection.prev = page.prev;
            }
            console.log('orderedCollection', orderedCollection);
            this.setState({ orderedCollection });
        }).catch(error => {
            logError(error);
        });
    }

    handleRefresh(e) {
        e.preventDefault();
        this.refreshOrderedCollection();
    }

    refreshOrderedCollection() {
        this.fetchOrderedCollection(this.state.orderedCollection.id);
    }

    get isPage() {
        return this.state.orderedCollection.first || 
        this.state.orderedCollection.prev ||
        this.state.orderedCollection.next ||
        this.state.orderedCollection.last;
    }

    handleFirst(e) {
        e.preventDefault();
        this.setPage("first");
    }

    handleNext(e) {
        e.preventDefault();
        this.setPage("next");
    }

    handlePrevious(e) {
        e.preventDefault();
        this.setPage("prev");
    }

    handleLast(e) {
        e.preventDefault();
        this.setPage("last");
    }

    renderItem(item) {
        switch (this.props.type) {
            case "Activity":
                return typeof item === 'string' ? 
                    // <Activity iri={wrapIfRemote(item, 'activity')} /> :
                    <Activity iri={item} isReferenced={true} /> :
                    <Activity activity={item} isReferenced={true} />;
            case "Object":
                return typeof item === 'string' ?
                    // <ActivityObject iri={wrapIfRemote(item, 'object')} /> :
                    <ActivityObject iri={item} isReferenced={true} /> :
                    <ActivityObject object={item} isReferenced={true} />
            case "Actor":
                return typeof item === 'string' ? 
                    <Actor iri={item} small={true} isReferenced={true} /> :
                    <Actor actor={item} small={true} isReferenced={true} />
            default:
                return typeof item === 'string' ? <IRI url={item} /> : <p>Unsupported Item</p>
        }
    }

    render() {
        return (
            <>
                {this.state.orderedCollection &&
                    <div className="OrderedCollection">
                        <div className="OrderedCollection-header">
                            <div className="OrderedCollection-titleblock">
                                <h5>
                                    { !this.props.link && 
                                    <>
                                        {this.props.title || this.state.orderedCollection.id}
                                    </>
                                    }
                                    { this.props.link && 
                                    <>
                                        <IRI url={wrapIfRemote(this.state.orderedCollection.id, 'collection', this.props.type)}
                                        text={this.props.title || this.state.orderedCollection.id} />
                                    </>
                                    }
                                </h5>
                                {this.state.orderedCollection.totalItems !== null && 
                                    <p className="m-0">Total items: {this.state.orderedCollection.totalItems}</p>
                                }
                            </div>
                            <div>
                                <button type="button" title="Refresh"
                                className="OrderedCollection-refresh button button-tertiary button-icon" onClick={this.handleRefresh} >
                                    <GrRefresh />
                                </button>
                            </div>

                        </div>
                        { this.state.loading && 
                        <>
                            <Divider />
                            <div className="OrderedCollection-items-loading">
                                <Loading />
                            </div>
                        </>
                        }
                        
                        { !this.state.loading && 
                        <>
                            { this.isPage &&
                            <>
                                <Divider />
                                <div className="w-full flex justify-between">
                                    <div className="page-button-container flex justify-start">
                                        { this.state.orderedCollection.first &&
                                            <button type="button" title="First Page"
                                            className="m-0 button button-secondary button-icon"
                                            onClick={this.handleFirst}>
                                                <MdFirstPage />
                                            </button>
                                        }
                                    </div>
                                    <div className="page-button-container flex justify-center">
                                        { this.state.orderedCollection.prev &&
                                            <button type="button" title="Previous Page"
                                            className="m-0 button button-tertiary button-icon"
                                            onClick={this.handlePrevious}>
                                                <MdNavigateBefore />
                                            </button>
                                        }
                                    </div>
                                    <div className="page-button-container flex justify-center">
                                        { this.state.orderedCollection.next &&
                                            <button type="button" title="Next Page"
                                            className="m-0 button button-tertiary button-icon"
                                            onClick={this.handleNext}>
                                                <MdNavigateNext />
                                            </button>
                                        }
                                    </div>
                                    <div className="page-button-container flex justify-end">
                                        { this.state.orderedCollection.last &&
                                            <button type="button" title="Last Page"
                                            className="m-0 button button-secondary button-icon"
                                            onClick={this.handleLast}>
                                                <MdLastPage />
                                            </button>
                                        }
                                    </div>
                                </div>
                            </>
                            }
                            { this.state.orderedCollection.orderedItems?.length > 0 && 
                            <>
                                <InfiniteScroll
                                className="overflow-visible"
                                dataLength={this.state.orderedCollection.orderedItems.length}
                                next={() => this.fetchPage('next', true)}
                                hasMore={this.state.orderedCollection.next}
                                scrollableTarget="content"
                                loader={<Loading />}>
                                {this.state.orderedCollection.orderedItems.map((item, _) => 
                                    <div key={`${this.state.orderedCollection.id}#${typeof item === 'string' ? item : item.id}`}>
                                        <Divider />
                                        {this.renderItem(item)}
                                    </div>
                                )}
                                </InfiniteScroll>
                            </>
                            }
                            { this.state.orderedCollection.orderedItems?.length === 0 && 
                            <>
                                <Divider />
                                <p>This collection is empty.</p>
                            </>
                            }
                        </>
                        }
                    </div>
                }
                { this.state.loading && !this.state.orderedCollection &&
                    <div className="OrderedCollection">
                        <Loading />
                    </div>
                }
                { this.state.unauthorized &&
                <>
                    { this.props.remote && !this.context ? <p>Sign in to browse the fediverse!</p> : <p>Unauthorized.</p> }
                </>
                }
                { this.state.notfound && <p>Not found.</p> }
            </>
        );
    }

}
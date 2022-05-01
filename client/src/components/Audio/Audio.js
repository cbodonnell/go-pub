import React from 'react';
import './Audio.scss';
import { proxyCacheBlobClient } from '../../utils/http';
import { FiHeadphones } from 'react-icons/fi';
import { Loading } from '../Loading/Loading';


export default class Audio extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            sources: [],
            showPlayer: false,
            loading: false
        };

        this.supportedMediaTypes = [
            'audio/ogg',
            'audio/wav',
            'audio/mp3',
            'audio/mpeg',
            'audio/mp4',
            'audio/acc',
            'audio/flac',
        ];

        this.toggleShowPlayer = this.toggleShowPlayer.bind(this);
    }

    componentDidMount() {
        if (Array.isArray(this.props.url)) {
            const sources = this.props.url.filter(url => this.supportedMediaTypes.includes(url.mediaType));
            this.setState({sources});
        }
    }

    toggleShowPlayer() {
        this.setState({showPlayer: !this.state.showPlayer}, () => {
            if (this.state.showPlayer) {
                this.fetchSources();
            }
        })
    }

    fetchSources() {
        this.setState({loading: true});
        this.state.sources.forEach(source => {
            proxyCacheBlobClient.get(source.href, {responseType: 'blob'}).then(res => {
                const objectUrl = res.headers.cached === true ? res.data : URL.createObjectURL(res.data);
                this.setState({loading: false}, () => this.addSource(objectUrl, source.mediaType));
            }).finally(() => this.setState({loading: false}));
        });
    }

    addSource(objectUrl, mediaType) {
        const audio = document.getElementById(`player-${this.props.id}`);
        const audioSource = document.createElement('source');
        audioSource.type = mediaType;
        audioSource.src = objectUrl;
        audio.appendChild(audioSource);
    }

    render() {
        return (
            <>
                <div className="Audio">
                    <p>
                        {this.props.name || this.props.id}
                        {/* <IRI url={this.props.id} text={this.props.name || this.props.id} /> */}
                    </p>
                    <button type="button" title={!this.state.showPlayer ? 'Play' : 'Hide Play'}
                    className={`m-0 button ${this.state.showPlayer ? 'button-secondary' : 'button-tertiary'} button-icon`} onClick={this.toggleShowPlayer} >
                        <FiHeadphones />
                    </button>
                    {this.state.showPlayer &&
                    <>
                        {this.state.loading && <Loading />}
                        {!this.state.loading &&
                        <>
                            { this.state.sources.length > 0 &&
                            <audio className="w-full mt-1" id={`player-${this.props.id}`} controls>
                                Your browser does not support the audio tag.
                            </audio>
                            }
                            {this.state.sources.length === 0 && <p>No sources found.</p>}
                        </>
                        }
                    </>
                    }
                </div>
            </>
        );
    }

}
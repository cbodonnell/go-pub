import React from 'react';
import './IRI.scss';
import { Link } from 'react-router-dom';
import { validURL } from '../../utils/urls';


export default class IRI extends React.Component {

    constructor(props) {
        super(props);
        this.state = { content: null }
        
        // this.handleSubmit = this.handleSubmit.bind(this);
    }

    componentDidMount() {
        const baseURL = window.location.protocol + '//' + window.location.hostname;
        let content
        if (this.props.url) {
            if (this.props.url.indexOf(baseURL) === 0) {
                // IRI is internal, create a router link
                content = <Link className="IRI link" to={this.props.url.replace(baseURL, '')}>{this.props.text || this.props.url}</Link>;
            } else {
                // IRI is external
                content = validURL(this.props.url) ? (
                    <a className="IRI link" target="_blank" rel="noreferrer" href={this.props.url}>{this.props.text || this.props.url}</a>
                ) : (
                    <span className="IRI">{this.props.text || this.props.url}</span>
                );
            }
        }
        this.setState({ content });
    }

    render() {
        return (
            <>
                { this.state.content }
            </>
        );
    }

}

import React from 'react';
import ReactDOMServer from 'react-dom/server';
import './Note.scss';
import SanitizedHTML from 'react-sanitized-html';
import { routeRemote } from '../../utils/urls';


export default class Note extends React.Component {

    constructor(props) {
        super(props);
        this.state = { content: null };
    }

    componentDidMount() {
        const content = this.proccessContent(this.props.content);
        this.setState({content});
    }

    sanitizeHTML(raw) {
        return <SanitizedHTML className="Note-content"
        allowedAttributes={{ 'a': ['href', 'target', 'ref'] }}
        allowedTags={['a', 'p', 'br']}
        transformTags={{
            'a': (tagName, attribs) => ({ tagName, attribs: {...attribs, target: '_blank', ref: 'noreferrer'}})
        }}
        html={raw} />
    }

    proccessContent(content) {
        const sanitized = this.sanitizeHTML(content);
        const sanitizedString = ReactDOMServer.renderToString(sanitized);
        const htmlDom = new DOMParser().parseFromString(sanitizedString, 'text/html');
        const nodes = htmlDom.documentElement.querySelectorAll('a');
        for (let i = 0; nodes[i]; i++) {
            const node = nodes[i];
            if (node.innerText.includes("@")) {
                const route = routeRemote(node.href, 'actor');
                node.href = route;
                // node.target = "_self";
                // // TODO: This is getting removed after setting...
                // node.onclick = (e) => {
                //     e.preventDefault();
                //     this.props.history.push(route)
                // };
            }
        }
        return htmlDom.documentElement.innerHTML;
    }

    render() {
        return (
            <>
                <div className="Note">
                    {this.state.content && 
                        <>
                            <p>Content:</p>
                            <div dangerouslySetInnerHTML={{ __html: this.state.content }}></div>
                            {/* <SanitizedHTML className="Note-content"
                            allowedAttributes={{ 'a': ['href', 'target', 'ref'] }}
                            allowedTags={['a', 'p', 'br']}
                            transformTags={{
                                'a': (tagName, attribs) => ({ tagName, attribs: {...attribs, target: '_blank', ref: 'noreferrer'}})
                            }}
                            html={this.state.content} /> */}
                        </>
                    }
                </div>
            </>
        );
    }

}
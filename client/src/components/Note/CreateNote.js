import React from 'react';
import './CreateNote.scss';
import axios from 'axios';
import logError from '../../utils/errors';
import AuthContext from '../../contexts/AuthContext';
import { apHeaders } from '../../utils/ap';
import { uniqueID } from '../../utils/forms';
import { environment } from '../../environment';
import { MdSend, MdClear } from 'react-icons/md';


export default class CreateNote extends React.Component {

    static contextType = AuthContext;

    constructor(props) {
        super(props);

        this.handleClear = this.handleClear.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);

        this.id = uniqueID();
    }

    handleClear(e) {
        e.preventDefault();
        this.clearForm();
    };

    handleSubmit(e) {
        e.preventDefault();
        const formData = new FormData(e.target);
        axios.post(`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/outbox`,
            {
                "@context": "https://www.w3.org/ns/activitystreams",
                type: "Note",
                content: formData.get('reply'),
                inReplyTo: this.props.inReplyTo,
                // published: "2015-02-10T15:04:55Z",
                to: this.props.to,
            },
            {
                headers: {'Content-Type': apHeaders.contentType},
                withCredentials: true
            }
        ).then(res => {
            console.log(res);
            this.clearForm();
            this.props.onCreate();
        }).catch(error => {
            logError(error);
        });
        // }).finally(() => this.setState({ processing: false }));
    };

    clearForm() {
        document.getElementById(`note-form-${this.id}`).reset();
    }

    render() {
        return (
            <>
                <div>
                    <form id={`note-form-${this.id}`} onSubmit={this.handleSubmit}>
                        {/* TODO: Add some sort of @mention@host mechanism */}
                        <textarea className="Note-textarea" name="reply"
                        placeholder="What would you like to say?" required></textarea>
                        <div className="Note-buttons w-full">
                            <button type="button" title="Clear"
                            className="button button-secondary button-icon"
                            onClick={this.handleClear}>
                                <MdClear />
                            </button>
                            <button type="submit" title="Send"
                            className="button button-primary button-icon">
                                <MdSend />
                            </button>
                        </div>
                    </form>
                </div>
            </>
        );
    }

}
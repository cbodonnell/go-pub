import React from 'react';
import './Publish.scss';
import AuthContext from '../../contexts/AuthContext';
import { uniqueID } from '../../utils/forms';
import { MdClear, MdUpload } from 'react-icons/md';
import axios from 'axios';
import { environment } from '../../environment';
import logError from '../../utils/errors';


export default class Publish extends React.Component {

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
        const object = {
            "@context": "https://www.w3.org/ns/activitystreams",
            type: "Audio",
            name: formData.get('name'),
            // published: "2015-02-10T15:04:55Z",
            to: this.props.to,
        };
        formData.append('file', document.getElementById('file').files[0]);
        formData.append('object', JSON.stringify(object));
        axios.post(`${environment.REACT_APP_ACTIVITY_URL}/users/${this.context.username}/endpoints/uploadMedia`,
            formData,
            {
                headers: {'Content-Type': "multipart/form-data"},
                withCredentials: true
            }
        ).then(res => {
            console.log(res);
            this.clearForm();
            this.props.onCreate();
        }).catch(error => {
            logError(error);
        });
    };

    clearForm() {
        document.getElementById(`publish-form-${this.id}`).reset();
    }

    render() {
        return (
            <>
                <div>
                    <form id={`publish-form-${this.id}`} onSubmit={this.handleSubmit}>
                        {/* <Form.Group controlId="file" className="mb-3">
                            <Form.Label>File</Form.Label>
                            <Form.Control type="file" />
                        </Form.Group> */}
                        <p className="Publish-label">Name:</p>
                        <input className="Publish-input" type="text" name="name" id="name" required />
                        <p className="Publish-label">File:</p>
                        <input type="file" name="file" id="file" required />
                        <div className="Publish-buttons w-full">
                            <button type="button" title="Clear"
                            className="button button-secondary button-icon"
                            onClick={this.handleClear}>
                                <MdClear />
                            </button>
                            <button type="submit" title="Publish"
                            className="button button-primary button-icon">
                                <MdUpload />
                            </button>
                        </div>
                    </form>
                </div>
            </>
        );
    }

}
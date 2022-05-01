import React from 'react';
import Create from '../Create/Create';
import './Footer.scss';

export default class Footer extends React.Component {

    constructor(props) {
        super(props);
        this.state = {};
    }

    componentDidMount() {

    }
 
    render() {
      return (
        <>
            <div className="Footer w-full flex-initial">
                <Create />
            </div>
        </>
      )
    }
  }
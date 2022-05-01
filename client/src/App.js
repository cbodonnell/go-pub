import './App.scss';
import Header from './components/Header/Header';
import Feed from './components/Feed/Feed';
import Actor from './components/Actor/Actor';
import React from 'react';
import { useState } from 'react';
import { BrowserRouter, Switch, Route, Redirect } from 'react-router-dom';
import Activity from './components/Activity/Activity';
import ActivityObject from './components/ActivityObject/ActivityObject';
import OrderedCollection from './components/OrderedCollection/OrderedCollection';
import AuthContext from './contexts/AuthContext';
import { getSearchParam } from './utils/urls';
import { environment } from './environment';
import { AuthLoading } from './components/Loading/AuthLoading';
import Footer from './components/Footer/Footer';
// import Footer from './components/Footer/Footer';


function App() {
  const [auth, setAuth] = useState(null);
  const [checkingAuth, setCheckingAuth] = useState(true);

  const handleAuth = (auth) => {
    setAuth(auth);
    setCheckingAuth(false);
    if (auth) {
      console.log('Authenticated', auth);
    } else {
      console.warn('Failed to authenticate', auth);
    }
  }

  return (
    <div className="App h-full">
      <header className="App-header w-full h-full flex flex-col">
        <BrowserRouter>
            <AuthContext.Provider value={auth}>
              <Header onAuth={handleAuth}/>
              <div className="App-content flex-auto" id="content">
                { !checkingAuth && 
                  <div className="App-page">
                      <Switch>
                          <Route exact path="/users/:name/outbox" render={(props) => <OrderedCollection {...props} iri={environment.REACT_APP_ACTIVITY_URL+props.match.url} title={`${props.match.params.name}'s Outbox`} type={'Activity'} />}/>
                          <Route exact path="/users/:name/inbox" render={(props) => <OrderedCollection {...props} iri={environment.REACT_APP_ACTIVITY_URL+props.match.url} title={`${props.match.params.name}'s Inbox`} type={'Activity'} />}/>
                          <Route exact path="/users/:name/followers" render={(props) => <OrderedCollection {...props} iri={environment.REACT_APP_ACTIVITY_URL+props.match.url} title={`${props.match.params.name}'s Followers`} />}/>
                          <Route exact path="/users/:name/following" render={(props) => <OrderedCollection {...props} iri={environment.REACT_APP_ACTIVITY_URL+props.match.url} title={`${props.match.params.name}'s Following`} />}/>
                          <Route exact path="/users/:name/liked" render={(props) => <OrderedCollection {...props} iri={environment.REACT_APP_ACTIVITY_URL+props.match.url} title={`${props.match.params.name}'s Liked`} type={'Object'} />}/>
                          <Route exact path="/users/:name" render={(props) => <Actor {...props} iri={`${environment.REACT_APP_ACTIVITY_URL}/users/${props.match.params.name}`} key={props.location.key} />}/>
                          <Route exact path="/activities/:id" render={(props) => <Activity {...props} iri={`${environment.REACT_APP_ACTIVITY_URL}/activities/${props.match.params.id}`} key={props.location.key} />}/>
                          <Route exact path="/objects/:id" render={(props) => <ActivityObject {...props} iri={`${environment.REACT_APP_ACTIVITY_URL}/objects/${props.match.params.id}`} key={props.location.key} />}/>
                          {/* TODO: Component for remote things? */}
                          <Route exact path="/actor" render={(props) => {
                            const remote = getSearchParam(props.location.search, "remote");
                            return <Actor {...props} iri={remote} key={props.location.key} remote={true} />;
                          }}/>
                          <Route exact path="/collection" render={(props) => {
                            const type = getSearchParam(props.location.search, "type");
                            const remote = getSearchParam(props.location.search, "remote");
                            return <OrderedCollection {...props} iri={remote} type={type} key={props.location.key} remote={true} />;
                          }}/>
                          <Route exact path="/activity" render={(props) => {
                            const remote = getSearchParam(props.location.search, "remote");
                            return <Activity {...props} iri={remote} key={props.location.key} remote={true} />;
                          }}/>
                          <Route exact path="/object" render={(props) => {
                            const remote = getSearchParam(props.location.search, "remote");
                            return <ActivityObject {...props} iri={remote} key={props.location.key} remote={true} />;
                          }}/>
                          <Route exact path="/" render={() => <Feed />}/>
                          <Route exact path="/unauthorized" render={() => <p>Unauthorized.</p>}/>
                          <Route exact path="/not-found" render={() => <p>Not found.</p>}/>
                          <Redirect to="/not-found" />
                      </Switch>
                  </div>
                }
                { checkingAuth && <AuthLoading />}
              </div>
              <Footer />
            </AuthContext.Provider>
        </BrowserRouter>
      </header>
    </div>
  );
}

export default App;

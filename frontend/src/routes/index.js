import React from 'react';
import {
	BrowserRouter,
	Route,
	Switch,
	Redirect,
} from 'react-router-dom';

import Home from './Home';
import Register from './Register';
import Login from './Login';
import CreateTeam from './CreateTeam';

const isAuthenticated = () => {
	return true;
}

const PrivateRoute = ({component: Component, ...rest}) => (
	<Route {...rest} render={props => 
		( isAuthenticated() ? (
			<Component {...props} />
		) : (
			<Redirect 
				to={{
					pathname: "/login",
				}}
			/>
		))}
	/>
);

export default () => (
	<BrowserRouter>
		<Switch>
			<PrivateRoute path="/" exact component={Home} />
			<Route path="/register" exact component={Register} />
			<Route path="/login" exact component={Login} />
			<Route path="/team" exact component={ViewTeam} />
			<PrivateRoute path="/create" exact component={CreateTeam} />
			<PrivateRoute path="/logout" exact component={Redirect} loc="https://initium.vc" />
		</Switch>
	</BrowserRouter>
);

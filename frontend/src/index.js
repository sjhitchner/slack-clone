import React from 'react';
import ReactDOM from 'react-dom';
import ApolloClient from "apollo-boost";
import { ApolloProvider } from 'react-apollo';
import 'semantic-ui-css/semantic.min.css';

import Routes from './routes';
import registerServiceWorker from './registerServiceWorker';

const client = new ApolloClient({
	uri: 'http://localhost:8080/graphql',
});

const App = (
	<ApolloProvider client={client}>
		<Routes />
	</ApolloProvider>
)

ReactDOM.render(App, document.getElementById('root'));
registerServiceWorker();

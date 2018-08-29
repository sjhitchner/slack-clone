import React from 'react';
import ReactDOM from 'react-dom';
import { ApolloClient } from 'apollo-client';
import { createHttpLink } from 'apollo-link-http';
import { setContext } from 'apollo-link-context';
import { InMemoryCache } from 'apollo-cache-inmemory';
//import ApolloClient from "apollo-boost";
import { ApolloProvider } from 'react-apollo';
import 'semantic-ui-css/semantic.min.css';

import Routes from './routes';
import registerServiceWorker from './registerServiceWorker';

const httpLink = createHttpLink({
	uri: 'http://localhost:8080/graphql',
});

const authLink = setContext((_, { headers }) => {
	const token = "XXXX"; // localStorage.getItem('token');
	return {
    headers: {
      authorization: token ? `Bearer ${token}` : "",
    }
  }
});

const client = new ApolloClient({
	link: authLink.concat(httpLink),
	cache: new InMemoryCache(),
});

const App = (
	<ApolloProvider client={client}>
		<Routes />
	</ApolloProvider>
)

ReactDOM.render(App, document.getElementById('root'));
registerServiceWorker();

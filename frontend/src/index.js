import React from 'react';
import ReactDOM from 'react-dom';
import { ApolloClient } from 'apollo-client';
import { createHttpLink } from 'apollo-link-http';
import { setContext } from 'apollo-link-context';
import { onError } from 'apollo-link-error'
import { InMemoryCache } from 'apollo-cache-inmemory';
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

const logoutLink = onError(({ response, networkError }, callback ) => {
	if (networkError && networkError.statusCode === 401) {
		console.log("Network Error");	
	}
});

const client = new ApolloClient({
	link: logoutLink.concat(authLink.concat(httpLink)),
	cache: new InMemoryCache(),
});

const App = (
	<ApolloProvider client={client}>
		<Routes />
	</ApolloProvider>
)

ReactDOM.render(App, document.getElementById('root'));
registerServiceWorker();

import React from 'react';
import { Message, Button, Input, Container, Header } from 'semantic-ui-react';
import gql from 'graphql-tag';
import { graphql } from 'react-apollo';

class Register extends React.Component {
  state = {
    username: '',
    email: '',
    password: '',
	usernameError: '',
	emailError: '',
	passwordError: '',
  };

  onSubmit = async () => {
	this.setState({
		usernameError: '',
		emailError: '',
		passwordError: '',
	});

    const {username, email, password} = this.state;
	const response = await this.props.mutate({
      variables: {username, email, password},
    });

    console.log(response);

	const {ok, errors } = response.data.createUser;
	if (!ok) {
		const err = {};
		errors.forEach(({field, message}) => {
			err[`${field}Error`] = message;
		});

		this.setState(err);
	} else {
		this.props.history.push('/');
	}
  };

  onChange = (e) => {
    const { name, value } = e.target;
    this.setState({ [name]: value });
  };

  render() {
    const { username, email, password, usernameError, emailError, passwordError } = this.state;

	const errorList = [];

	if (usernameError) {
		errorList.push(usernameError);
	}

	if (emailError) {
		errorList.push(emailError);
	}

	if (passwordError) {
		errorList.push(passwordError);
	}

    return (
      <Container text>
        <Header as="h2">Register</Header>
		{ usernameError || emailError || passwordError ? (
			<Message error header="Registration Errors" list={errorList} />
		) : null}
        <Input
          name="username"
		  error={!!usernameError}
          onChange={this.onChange}
          value={username}
          placeholder="Username"
          fluid
        />
        <Input
		  name="email"
		  error={!!emailError}
		  onChange={this.onChange}
		  value={email}
		  placeholder="Email"
		  fluid
		/>
        <Input
          name="password"
		  error={!!passwordError}
          onChange={this.onChange}
          value={password}
          type="password"
          placeholder="Password"
          fluid
        />
        <Button onClick={this.onSubmit}>Submit</Button>
      </Container>
    );
  }
}

const registerMutation = gql`
  mutation($username: String!, $email: String!, $password: String!) {
    createUser(input: {
		username: $username,
		email: $email,
		password: $password
	}) {
	  ok
	  errors {
	    type
		field
		message
	  }
	}
  }
`;

export default graphql(registerMutation)(Register);

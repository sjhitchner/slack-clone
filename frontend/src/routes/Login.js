import React from 'react';
import { Message, Button, Input, Container, Header } from 'semantic-ui-react';
import gql from 'graphql-tag';
import { graphql } from 'react-apollo';

class Login extends React.Component {
  state = {
    email: '',
    password: '',
	emailError: '',
	passwordError: '',
  };

  onSubmit = async () => {
	this.setState({
		emailError: '',
		passwordError: '',
	});

    const {email, password} = this.state;
	const response = await this.props.mutate({
      variables: {email, password},
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
    const {email, password, emailError, passwordError } = this.state;

	const errorList = [];

	if (emailError) {
		errorList.push(emailError);
	}

	if (passwordError) {
		errorList.push(passwordError);
	}

    return (
      <Container text>
        <Header as="h2">Register</Header>
		{ emailError || passwordError ? (
			<Message error header="Login Errors" list={errorList} />
		) : null}
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

export default Login;

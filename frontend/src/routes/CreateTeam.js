import React from 'react';
import { Message, Button, Input, Container, Header } from 'semantic-ui-react';
import gql from 'graphql-tag';
import { graphql } from 'react-apollo';

class CreateTeam extends React.Component {
  state = {
    name: '',
	nameError: '',
  };

  onSubmit = async () => {
	this.setState({
		nameError: '',
	});

    const {name} = this.state;
	const response = await this.props.mutate({
      variables: {name},
    });

    console.log(response);

	const {ok, errors} = response.data.createTeam;
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
	  const { name, nameError } = this.state;

	  const errorList = [];
	  if (nameError) {
		  errorList.push(nameError);
	  }

	  return (
	    <Container text>
        <Header as="h2">Create Team</Header>
		{ nameError ? ( <Message error header="Creating Team Errors" list={errorList} />) : null}
        <Input
          name="name"
		  error={!!nameError}
          onChange={this.onChange}
          value={name}
          placeholder="Team Name"
          fluid
        />
        <Button onClick={this.onSubmit}>Submit</Button>
      </Container>
    );
  }
}

const createTeamMutation = gql`
mutation ($name: String!) {
  createTeam(input: {name: $name}) {
    ok
    team {
      id
    }
    errors {
      field
      message
    }
  }
}
`;

export default graphql(createTeamMutation)(CreateTeam);

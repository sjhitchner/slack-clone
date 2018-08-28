import React from 'react';
import gql from 'graphql-tag';
import { graphql } from 'react-apollo';

const userListQuery = gql`
{
  userList {
    id
    username
    email
  }
}
`;

const Home = ({data: {userList = []}}) => userList.map(u => <h1 key={u.id}>{u.username} {u.email}</h1>);

export default graphql(userListQuery)(Home);

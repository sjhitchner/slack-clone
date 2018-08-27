package main

const Query = `
{
  userTeamList(userId: 1) {
    name
    owner {
      username
    }
    members {
      username
    }
    channels {
      name
      owner {
        username
      }
      members {
        username
      }
      messages {
        text
      }
    }
  }
}
`

import React from 'react';
import AppBar from 'material-ui/AppBar';
import Toolbar from 'material-ui/Toolbar';
import Typography from 'material-ui/Typography';

export default (props) => (
  <AppBar position="static" color="default">
      <Toolbar>
      <Typography variant="title" color="inherit">
          GraphQL Subscription Server example
      </Typography>
      </Toolbar>
  </AppBar>
)

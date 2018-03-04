import React from 'react';
import lightBaseTheme from 'material-ui/styles/baseThemes/lightBaseTheme';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import AppBar from 'material-ui/AppBar';

import Name, { NameWatcher } from './Name';
import { NewChannel, ChannelList, ChannelWatcher } from './Channel';
import Notification, { NotificationWatcher } from './Notification';

const nameWatcher = new NameWatcher()
const chanWatcher = new ChannelWatcher();
const notifWatcher = new NotificationWatcher();

nameWatcher.onRegisterCallback = () => {

};

const App = (props) => {
  return (
    <MuiThemeProvider muiTheme={getMuiTheme(lightBaseTheme)}>
      <AppBar title="GraphQL Subscription Server example" />
      <Name watcher={nameWatcher} />
      <Notification watcher={notifWatcher} />
      <NewChannel watcher={chanWatcher} />
      <ChannelList watcher={chanWatcher} />
    </MuiThemeProvider>
  );
}

export default App

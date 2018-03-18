import React from 'react';
import { MuiThemeProvider, createMuiTheme } from 'material-ui/styles';
import AppBar from 'material-ui/AppBar';
import Name, { NameWatcher } from './Name';
import { NewChannel, ChannelList, ChannelWatcher } from './Channel';
import Notification, { NotificationWatcher } from './Notification';

const nameWatcher = new NameWatcher()
const chanWatcher = new ChannelWatcher();
const notifWatcher = new NotificationWatcher();

nameWatcher.onRegisterCallback = () => {

};

const theme = createMuiTheme();

const App = (props) => {
  return (
    <MuiThemeProvider theme={theme}>
      <AppBar title="GraphQL Subscription Server example" />
      <Name watcher={nameWatcher} />
      <Notification watcher={notifWatcher} />
      <NewChannel watcher={chanWatcher} />
      <ChannelList watcher={chanWatcher} />
    </MuiThemeProvider>
  );
}

export default App

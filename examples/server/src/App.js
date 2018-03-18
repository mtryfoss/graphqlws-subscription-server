import React from 'react';
import { MuiThemeProvider, createMuiTheme } from 'material-ui/styles';
import Header from './Header';
import Name, { NameWatcher } from './Name';
import { NewChannel, ChannelList, ChannelWatcher } from './Channel';
import Notification, { NotificationWatcher } from './Notification';

const nameWatcher = new NameWatcher()
const chanWatcher = new ChannelWatcher();
const notifWatcher = new NotificationWatcher();

nameWatcher.onRegisterCallback = () => {

};

const theme = createMuiTheme();

export default (props) => (
  <MuiThemeProvider theme={theme}>
    <Header />
    <Name watcher={nameWatcher} />
    <Notification watcher={notifWatcher} />
    <NewChannel watcher={chanWatcher} />
    <ChannelList watcher={chanWatcher} />
  </MuiThemeProvider>
)

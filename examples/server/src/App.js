import React from 'react';
import { MuiThemeProvider, createMuiTheme } from 'material-ui/styles';
import Header from './Header';
import Name, { NameWatcher } from './Name';
import { NewChannel, ChannelList, ChannelWatcher } from './Channel';
import Notification, { NotificationWatcher } from './Notification';
import styled, { css } from 'styled-components';
import { inject, observer } from 'mobx-react';

const nameWatcher = new NameWatcher()
const chanWatcher = new ChannelWatcher();
const notifWatcher = new NotificationWatcher();

nameWatcher.onRegisterCallback = () => {

};

const theme = createMuiTheme();

const AppContainer = observer(styled.div`
  ${ props => !props.isWalkedIn && css`
  display: none;
  `}
`)

export default (props) => (
  <MuiThemeProvider theme={theme}>
    <Header />
    <Name watcher={nameWatcher} />
    <AppContainer isWalkedIn={nameWatcher.isWalkedIn}>
      <Notification watcher={notifWatcher} />
      <NewChannel watcher={chanWatcher} />
      <ChannelList watcher={chanWatcher} />
    </AppContainer>
  </MuiThemeProvider>
)

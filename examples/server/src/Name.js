import React from 'react'
import { observer } from 'mobx-react'
import { observable } from 'mobx'
import RaisedButton from 'material-ui/RaisedButton';
import TextField from 'material-ui/TextField';
import {Card, CardHeader, CardActions} from 'material-ui/Card';
import Paper from 'material-ui/Paper';
import styled, { css } from 'styled-components';

const style = {
  name: {
    margin: 16,
    width: 400,
    height: 72,
  },
  nameInner: {
    padding: 8,
    textAlign: 'center',
    display: 'inline-block',
  },
  button: { margin: 16 },
};

const EntryForm = styled.div`
  ${props => props.isWalkedIn && css`
  display: none;
  `}
`

const NameDisplay = styled.div`
  margin-top: 16px;
  margin-left: 16px;
  ${props => !props.isWalkedIn && css`
  display: none;
  `}
`

class NameWatcher {
  @observable name = "";
  @observable isWalkedIn = false;
  onRegisterCallback = () => {};

  onNameChanged(e) {
    e.preventDefault();
    this.name = e.target.value;
  }

  walkIn(e) {
    e.preventDefault();
    if (this.name !== "") {
      this.isWalkedIn = true;
      this.onRegisterCallback();
    }
  }
}

const Name = observer((props) => {
  const w = props.watcher
  return (
    <Card style={style.name}>
      <CardActions style={style.nameInner}>
        <EntryForm isWalkedIn={w.isWalkedIn}>
          <TextField hintText="type your name ascii only" value={w.name} onChange={(e) => w.onNameChanged(e)} />
          <RaisedButton label="入室" primary={true} style={style.button} onClick={(e) => w.walkIn(e)} />
        </EntryForm>
        <NameDisplay isWalkedIn={w.isWalkedIn}>{w.name} さんが入室しています</NameDisplay>
      </CardActions>
    </Card>
  );
});

export default Name
export { NameWatcher }
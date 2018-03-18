import React from 'react'
import { observer } from 'mobx-react'
import { observable } from 'mobx'
import Button from 'material-ui/Button';
import TextField from 'material-ui/TextField';
import Card, {CardContent, CardActions} from 'material-ui/Card';
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

export default observer((props) => {
  const { watcher } = props
  return (
    <Card>
      <CardContent>
        <CardActions>
          <EntryForm isWalkedIn={watcher.isWalkedIn}>
            <TextField hintText="type your name ascii only" value={watcher.name} onChange={(e) => watcher.onNameChanged(e)} />
            <Button variant="raised" color="primary" style={style.button} onClick={(e) => watcher.walkIn(e)}>入室</Button>
          </EntryForm>
        </CardActions>
        <NameDisplay isWalkedIn={watcher.isWalkedIn}>{watcher.name} さんが入室しています</NameDisplay>
      </CardContent>
    </Card>
  );
});

export { NameWatcher }
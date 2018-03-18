import React from 'react';
import { observer } from 'mobx-react';
import { observable } from 'mobx';
import Button from 'material-ui/Button';
import TextField from 'material-ui/TextField';
import Card, {CardHeader, CardContent, CardActions} from 'material-ui/Card';
import List, { ListItem, ListItemText } from 'material-ui/List';

const notificationStyle = {
  margin: 16,
  height: 320,
  width: 400,
};

const buttonStyle = {
  margin: 12,
};

class NotificationWatcher {
  @observable message = "";
  @observable messages = [];

  onMessageChanged(e) {
    e.preventDefault();
    this.message = e.target.value;
  }

  post() {
    // TODO
    this.receive(this.message);
    this.message = "";
  }

  receive(msg) {
    const messages = this.messages.slice();
    this.messages = messages.concat(msg);
  }
}

export default observer((props) => {
  const { watcher } = props;
  return (
    <Card>
      <CardHeader
        title="お知らせ"
        subheader="全員に対しての通知はここに入ります"
      />
      <CardContent>
        <List style={{height:160}}>
          {watcher.messages.map((msg) => <ListItem><ListItemText primaryText={msg} /></ListItem>)}
        </List>

        <CardActions>
          <TextField hintText="" value={watcher.message} onChange={(e) => watcher.onMessageChanged(e)} />
          <Button variant="raised" color="primary" style={buttonStyle} onClick={() => watcher.post()}>送信</Button>
        </CardActions>
      </CardContent>
    </Card>
  );
});

export { NotificationWatcher }

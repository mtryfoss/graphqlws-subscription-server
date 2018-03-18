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

const Notification = observer((props) => {
  const w = props.watcher;
  return (
    <Card>
      <CardHeader
        title="お知らせ"
        subheader="全員に対しての通知はここに入ります"
      />
      <CardContent>
        <List style={{height:160}}>
          {w.messages.map((msg) => <ListItem><ListItemText primaryText={msg} /></ListItem>)}
        </List>

        <CardActions>
          <TextField hintText="" value={w.message} onChange={(e) => w.onMessageChanged(e)} />
          <Button variant="raised" color="primary" style={buttonStyle} onClick={() => w.post()}>送信</Button>
        </CardActions>
      </CardContent>
    </Card>
  );
});

export default Notification
export { NotificationWatcher }

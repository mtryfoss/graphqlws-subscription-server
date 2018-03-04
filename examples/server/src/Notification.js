import React from 'react';
import { observer } from 'mobx-react';
import { observable } from 'mobx';
import RaisedButton from 'material-ui/RaisedButton';
import TextField from 'material-ui/TextField';
import {Card, CardHeader, CardActions} from 'material-ui/Card';
import {List, ListItem} from 'material-ui/List';

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
    <Card style={notificationStyle}>
      <CardHeader
        title="お知らせ"
        subtitle="全員に対しての通知はここに入ります"
      />
      <CardActions>
        <List style={{height:160}}>
          {w.messages.map((msg) => <ListItem primaryText={msg} />)}
        </List>

        <div className="Notification-post">
          <TextField hintText="" value={w.message} onChange={(e) => w.onMessageChanged(e)} />
          <RaisedButton label="送信" primary={true} style={buttonStyle} onClick={() => w.post()} />
        </div>
      </CardActions>
    </Card>
  );
});

export default Notification
export { NotificationWatcher }

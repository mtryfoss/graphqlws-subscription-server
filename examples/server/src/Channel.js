import React from 'react';
import { observer } from 'mobx-react';
import { observable } from 'mobx';
import Button from 'material-ui/Button';
import TextField from 'material-ui/TextField';
import Paper from 'material-ui/Paper';

const buttonStyle = {
  margin: 12,
};

const paperStyle = {
  height: 70,
  width: 400,
  margin: 20,
  textAlign: 'center',
};

class ChannelWatcher {
  @observable newName = "";
  @observable channels = [];

  newChannelCallback = function(name) {
    alert("new channel: " + name);
    this.channels.push(<Channel watcher={new ChannelInfo(name)} />)
  };

  onNewNameChanged(e) {
    e.preventDefault();
    this.newName = e.target.value;
  }

  addChannel(e) {
    e.preventDefault();
    if (this.newName !== "") {
      this.newChannelCallback(this.newName);
      this.newName = "";
    }
  }
}

class ChannelInfo {
  name = "";
  @observable message = "";
  @observable messages = [];

  constructor(name) {
    this.name = name;
  }

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

const NewChannel = observer((props) => {
  const w = props.watcher;
  return (
    <Paper style={paperStyle} zDepth={2}>
      <TextField hintText="type name you want to talk to" value={w.newName} onChange={(e) => w.onNewNameChanged(e)} />
      <Button variant="raised" label="+" color="primary" style={buttonStyle} onClick={(e) => w.addChannel(e)} />
    </Paper>
  );
});

const Channel = observer((props) => {
  const w = props.watcher;
  return (
    <div className="Channel">
      <div className="Channel-target">{w.name}</div>
      <div className="Channel-messages-container">
        <ul className="Channel-messages-list">
        {w.messages.map((msg) => <li>{msg}</li> )}
        </ul>
      </div>
      <div className="Channel-post">
        <div className="Channel-post-input">
          <input type="text" name="message" value={w.message} onChange={(e) => w.onMessageChanged(e)} />
        </div>
        <div className="Channel-post-register">
          <button onClick={() => w.post()}>送信</button>
        </div>
      </div>
    </div>
  );
});

const ChannelList = observer((props) => {
  const w = props.watcher;
  return (
    <div className="ChannelList">
      <div className="ChannelList-name">登録チャンネル</div>
      <div className="ChannelList-list">
        {w.channels}
      </div>
    </div>
  );
});

export { NewChannel, Channel, ChannelList, ChannelWatcher }

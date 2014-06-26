/** @jsx React.DOM */
var WCCard = React.createClass({
  addCard: function() {
    this.props.addCard(this.props.player.name);
  },
  removeCard: function() {
    this.props.removeCard(this.props.player.nameAlias);
  },
  render: function() {
    var dob = new Date(this.props.player.birthDate),
    style = {
      'font-size': '14px'
    }
    return (
      <div className="card">
        <h1 style={this.props.player.name.length > 30 ? style : {}}>{this.props.player.name}<i style={this.props.user ? {display:'visible'} : {display:'none'}} className="fa fa-plus add-card" onClick={this.addCard}></i><i style={this.props.user ? {display:'visible'} : {display:'none'}}className="fa fa-minus add-card" onClick={this.removeCard}></i></h1>
        <a className="wc-image" href="#">
          <img src={this.props.player.image} />
        </a>
        <ul className="card-data-long">
          <li>
            <label>National</label>
            <span>{this.props.player.national}&nbsp;</span>
          </li>
          <li>
            <label>Club</label>
            <span>{this.props.player.club}&nbsp;</span>
          </li>
          <li>
            <label>Position</label>
            <span>{this.props.player.position}&nbsp;</span>
          </li>
        </ul>
        <ul className="card-data">
          <li>
            <label>Height</label>
            <span>{this.props.player.height}cm&nbsp;</span>
          </li>
          <li>
            <label>Weight</label>
            <span>{this.props.player.weight}kg&nbsp;</span>
          </li>
          <li>
            <label>Dob</label>
            <span>{dob.getFullYear() + '/' + (dob.getMonth() + 1) + '/' + dob.getDate()}&nbsp;</span>
          </li>
          <li>
            <label>Age</label>
            <span>{this.props.player.age}&nbsp;</span>
          </li>
          <li>
            <label>Foot</label>
            <span>{this.props.player.foot}&nbsp;</span>
          </li>
          <li>
            <label>Goals</label>
            <span>{this.props.player.goals}&nbsp;</span>
          </li>
          <li>
            <label>Assists</label>
            <span>{this.props.player.assists}&nbsp;</span>
          </li>
          <li>
            <label>Penalty Goals</label>
            <span>{this.props.player.penaltyGoals}&nbsp;</span>
          </li>
          <li>
            <label>OwnGoals</label>
            <span>{this.props.player.ownGoals}&nbsp;</span>
          </li>
        </ul>
      </div>
    );
  }
});

var WCCards = React.createClass({
  render: function() {
    var addCard = this.props.addCard,
      removeCard = this.props.removeCard,
      user = this.props.user;

    var cardNodes = this.props.players.map(function(player, arr) {
      return <WCCard player={player} addCard={addCard} removeCard={removeCard} user={user} />;
    }); 

    return <div className="card-wrapper">{cardNodes}</div>;
  }
});

var SearchWCCard = React.createClass({
  mixins: [React.addons.LinkedStateMixin],
  getInitialState: function() {
    return {
      playerName: '',
      nameFilter: '',
      players: WC_PLAYERS
    };
  },

  searchCard: function(e) {
    // check if it's being entered or clicked on button
    if ((e.type === 'keypress' && e.keyCode === 13) || e.type === 'click')
      this.props.searchCard(this.refs.playerName.getDOMNode().value)
  },

  handleFilterChange: function() {
    this.setState({
      nameFilter: this.refs.playerName.getDOMNode().value
    });
  },

  // select players from the list
  selectPlayer: function(e) {
    this.refs.playerName.getDOMNode().value = e.target.dataset.name;
  },

  // gives timeout before blur on input
  blurList: function() {
    var self = this;
    setTimeout(function() {
      self.props.toggle();
    }, 300)
  },

  render: function() {
    var self = this;
    var displayedNames = this.state.players.filter(function(item) {
      var match = item.toLowerCase().indexOf(this.state.nameFilter.toLowerCase());
      return (match !== -1);
    }.bind(this));

    var playersLi = displayedNames.map(function(name) {
      return <li data-name={name} onClick={self.selectPlayer}>{name}</li>;
    });
    var playersUl = <ul>{playersLi}</ul>;

    return (
      <div id="search-card">
        <input id="player-name" type="text" onChange={this.handleFilterChange} onBlur={this.blurList} onFocus={this.props.toggle} ref="playerName" onKeyPress={this.searchCard} />
        <button className="button" onClick={this.searchCard}>Search</button>
        <div id="players-list">
          {playersUl}
        </div>
      </div>
    );
  }
});


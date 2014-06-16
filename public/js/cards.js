/** @jsx React.DOM */
var Card = React.createClass({
  addCard: function() {
    this.props.addCard(this.props.player.nameAlias);
  },
  removeCard: function() {
    this.props.removeCard(this.props.player.nameAlias);
  },
  render: function() {
    return (
      <div className="card">
        <h1>{this.props.player.name}<i className="fa fa-plus add-card" onClick={this.addCard}></i><i className="fa fa-minus add-card" onClick={this.removeCard}></i></h1>
        <a href="#">
          <img src={"http://www.premierleague.com/" + this.props.player.image} />
        </a>
        <ul className="card-data">
          <li>
            <label>Club</label>
            <span>{this.props.player.club}</span>
          </li>
          <li>
            <label>Position</label>
            <span>{this.props.player.position}</span>
          </li>
          <li>
            <label>Appearances</label>
            <span>{this.props.player.appearances}</span>
          </li>
          <li>
            <label>Goals</label>
            <span>{this.props.player.goals}</span>
          </li>
          <li>
            <label>Shots</label>
            <span>{this.props.player.shots}</span>
          </li>
          <li>
            <label>Penalties</label>
            <span>{this.props.player.penalties}</span>
          </li>
          <li>
            <label>Assists</label>
            <span>{this.props.player.assists}</span>
          </li>
          <li>
            <label>Crosses</label>
            <span>{this.props.player.crosses}</span>
          </li>
          <li>
            <label>Offsides</label>
            <span>{this.props.player.offsides}</span>
          </li>
          <li>
            <label>SavesMade</label>
            <span>{this.props.player.savesMade}</span>
          </li>
          <li>
            <label>OwnGoals</label>
            <span>{this.props.player.ownGoals}</span>
          </li>
          <li>
            <label>CleanSheets</label>
            <span>{this.props.player.ownGoals}</span>
          </li>
          <li>
            <label>Blocks</label>
            <span>{this.props.player.blocks}</span>
          </li>
          <li>
            <label>Clearances</label>
            <span>{this.props.player.clearances}</span>
          </li>
          <li>
            <label>Fouls</label>
            <span>{this.props.player.fouls}</span>
          </li>
          <li>
            <label>Cards</label>
            <span>{this.props.player.cards}</span>
          </li>
          <li>
            <label>Dob</label>
            <span>{this.props.player.dob}</span>
          </li>
          <li>
            <label>Height</label>
            <span>{this.props.player.height}</span>
          </li>
          <li>
            <label>Age</label>
            <span>{this.props.player.age}</span>
          </li>
          <li>
            <label>Weight</label>
            <span>{this.props.player.weight}</span>
          </li>
          <li>
            <label>National</label>
            <span>{this.props.player.national}</span>
          </li>
        </ul>
      </div>
    );
  }
});

var Cards = React.createClass({
  render: function() {
    var addCard = this.props.addCard;
    var removeCard = this.props.removeCard;
    var cardNodes = this.props.players.map(function(player, arr) {
      return <Card player={player} addCard={addCard} removeCard={removeCard}/>;
    }); 

    return <div className="card-wrapper">{cardNodes}</div>;
  }
});

var SearchCard = React.createClass({
  mixins: [React.addons.LinkedStateMixin],
  getInitialState: function() {
    return {
      playerName: '',
      nameFilter: '',
      players: PLAYERS
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
      <div>
        <input id="player-name" type="text" onChange={this.handleFilterChange} onBlur={this.blurList} onFocus={this.props.toggle} ref="playerName" onKeyPress={this.searchCard} />
        <button className="button" onClick={this.searchCard}>Search</button>
        <div id="players-list">
          {playersUl}
        </div>
      </div>
    );
  }
});

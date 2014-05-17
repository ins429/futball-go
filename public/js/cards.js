/** @jsx React.DOM */
var Card = React.createClass({
  render: function() {
    return (
      <div className="card-wrapper">
        <div className="card">
          <h1>{ this.props.player.name }<i className="fa fa-plus add-card" onClick={ this.props.addCard(this.props.player.nameAlias) }></i></h1>
          <a href="#">
            <img src={"http://www.premierleague.com/" + this.props.player.image } />
          </a>
          <ul className="card-data">
            <li>
              <label>Club</label>
              <span>{ this.props.player.club }</span>
            </li>
            <li>
              <label>Position</label>
              <span>{ this.props.player.position }</span>
            </li>
            <li>
              <label>Appearances</label>
              <span>{ this.props.player.appearances }</span>
            </li>
            <li>
              <label>Goals</label>
              <span>{ this.props.player.goals }</span>
            </li>
            <li>
              <label>Shots</label>
              <span>{ this.props.player.shots }</span>
            </li>
            <li>
              <label>Penalties</label>
              <span>{ this.props.player.penalties }</span>
            </li>
            <li>
              <label>Assists</label>
              <span>{ this.props.player.assists }</span>
            </li>
            <li>
              <label>Crosses</label>
              <span>{ this.props.player.crosses }</span>
            </li>
            <li>
              <label>Offsides</label>
              <span>{ this.props.player.offsides }</span>
            </li>
            <li>
              <label>SavesMade</label>
              <span>{ this.props.player.savesMade }</span>
            </li>
            <li>
              <label>OwnGoals</label>
              <span>{ this.props.player.ownGoals }</span>
            </li>
            <li>
              <label>CleanSheets</label>
              <span>{ this.props.player.ownGoals }</span>
            </li>
            <li>
              <label>Blocks</label>
              <span>{ this.props.player.blocks }</span>
            </li>
            <li>
              <label>Clearances</label>
              <span>{ this.props.player.clearances }</span>
            </li>
            <li>
              <label>Fouls</label>
              <span>{ this.props.player.fouls }</span>
            </li>
            <li>
              <label>Cards</label>
              <span>{ this.props.player.cards }</span>
            </li>
            <li>
              <label>Dob</label>
              <span>{ this.props.player.dob }</span>
            </li>
            <li>
              <label>Height</label>
              <span>{ this.props.player.height }</span>
            </li>
            <li>
              <label>Age</label>
              <span>{ this.props.player.age }</span>
            </li>
            <li>
              <label>Weight</label>
              <span>{ this.props.player.weight }</span>
            </li>
            <li>
              <label>National</label>
              <span>{ this.props.player.national }</span>
            </li>
          </ul>
        </div>
      </div>
    );
  }
});


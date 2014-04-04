app.controller('Cards', function cardsCtrl($scope, $route, $routeParams, api) {
  $scope.initCards = function() {
    $scope.cards = [];
  };

  $scope.getCards = function() {
    var params = {
      name: $scope.playerName
    }
    api.queryCards(params).then(function(data) {
      console.log(data)
      if (data && data.stats) {
        $scope.cards = data.stats;
        console.log($scope.cards)
      } else {
        // fix me, handle error
      }
    });
  };

  // pop form to create an card
  $scope.addCard = function() {
    $scope.cardData = {
      action: 'add'
    };
  };

  // pop form to edit an card
  $scope.editCard = function(card) {
    $scope.cardData = {
      action: 'edit',
      modalName: card.name,
      id: card.id,
      name: card.name,
      tags: card.tags,
      url: card.url
    }
  };

});
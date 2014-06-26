app.controller('Cards', function cardsCtrl($rootScope, $scope, $route, $routeParams, api) {
  $scope.initCards = function() {
    $rootScope.cards = [];
  };
});

app.directive('searchWcCard', function(api, $rootScope) {
  return {
    restrict: 'E',
    scope: {
    },
    link: function (scope, elem, attrs) {
      scope.searchCard = function(name) {
        // do nothing if name is not there
        if (!name || name === '')
          return;

        var params = {
          names: [name]
        }
        api.getWCCard(params).then(function(data) {
          if (data && data.stats) {
            $rootScope.cards = mergeCards($rootScope.cards, data.stats);
            scope.playerName = '';
          } else {
            // fix me, handle error
          }
        });
      };

      scope.toggle = function() {
        angular.element(document.querySelector('#players-list')).toggleClass('active');
      };

      scope.$watch('init', function(){
        React.renderComponent(window.SearchWCCard({searchCard: scope.searchCard, toggle: scope.toggle}), elem[0]);
      })
    }
  }
});

app.directive('cards', function(api, $rootScope, $timeout) {
  return {
    restrict: 'E',
    scope: {
      cards: '=',
    },
    link: function (scope, elem, attrs) {
      scope.$watch('cards', function(oldVal, newVal) {
        updateView();
      });

      function updateView() {
        if (scope.cards.length > 0) {
          React.renderComponent(window.WCCards({players: scope.cards, removeCard: scope.removeCard, addCard: scope.addCard, user: scope.user}), elem[0]);
        }
      }
    }
  }
});

function removeCard(cards, name) {
  for (var i = 0; i < cards.length; i++) {
    if (cards[i] && cards[i].nameAlias && cards[i].nameAlias === name) {
      delete cards[i];
    }
  }
  return cards;
}

function mergeCards(cards, newCards) {
  for (var i = 0; i < cards.length; i++) {
    for (var j = 0; j < newCards.length; j ++) {
      if (cards[i].nameAlias && cards[i].nameAlias === newCards[j].nameAlias) {
        newCards.splice(j, 1);
      } else if (cards[i].name && cards[i].name === newCards[j].name) {
        newCards.splice(j, 1);
      }
    }
  }
  return newCards.concat(cards);
}


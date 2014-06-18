app.controller('Users', function UsersCtrl($rootScope, $scope, $route, api) {
  $scope.initUser = function() {
    $scope.userForm = {};
    $scope.showMe();
  }

  $scope.login = function() {
    api.login({
      email: $scope.userForm.email,
      password: $scope.userForm.password
    }).then(function(resp) {
      if (resp.status && resp.status === 200) {
        $rootScope.user = resp.user;
      } else if (resp.message) {
        $rootScope.flash = resp.message
      }
    });
  };

  $scope.logout = function() {
    api.logout().then(function(resp) {
      if (resp.status && resp.status === 200) {
        $rootScope.user = null;
      }
    });
  };

  $scope.showMe = function() {
    api.showMe().then(function(resp) {
      if (resp.status && resp.status === 200) {
        $rootScope.user = resp.user;
        var players = resp.user.players,
          names = [];

        // if players not found, do nothing
        if (!players || players.length === 0) {
          return;
        }
        // collect names
        angular.forEach(players, function(value){
          this.push(value.name);
        }, names);

        api.getWCCard({
          names: names
        }).then(function(data) {
          if (data && data.stats) {
            $rootScope.cards = data.stats;
            $scope.playerName = '';
          }
        })
      }
    });
  };

  $scope.fbLogin = function() {
    FB.getLoginStatus(function(response) {
      if (response.status === 'connected') {
        api.fbLogin({
          token: response.authResponse.accessToken
        }).then(function(resp) {
          $rootScope.user = resp.user;
        });
      } else {
        FB.login(function(response) {
          if (response.authResponse) {
            api.fbLogin({
              token: response.authResponse.accessToken
            }).then(function(resp) {
              $rootScope.user = resp.user;
            });
          } else {
            // fix me handle error for fb auth
          }
        });
      }
    });
  }

  $scope.fbSignup = function() {
    FB.getLoginStatus(function(response) {
      if (response.status === 'connected') {
        api.fbSignup({
          token: response.authResponse.accessToken
        }).then(function(resp) {
          if (resp.status === 200) {
            $rootScope.user = resp;
          } else if (resp.message) {
            $rootScope.flash = resp.message
          }
        });
      } else {
        FB.login(function(response) {
          if (response.authResponse) {
            api.fbSignup({
              token: response.authResponse.accessToken
            }).then(function(resp) {
              if (resp.status === 200) {
                $rootScope.user = resp;
              } else if (resp.message) {
                $rootScope.flash = resp.message
              }
            });
          } else {
            // fix me handle error for fb auth
          }
        });
      }
    });
  }

  // pop form to edit an image
  $scope.editUser = function() {
    $scope.togglerUserForm();
    $scope.userData = {
      action: 'Edit',
      first_name: $rootScope.user.first_name,
      last_name: $rootScope.user.last_name,
      email: $rootScope.user.email
    };
  };

  $scope.signOut = function() {
    $scope.showUserMenu = false;

    api.logout().then(function(data) {
      if (data && data.status && data.status === 200 ) {
        $rootScope.user = null;
      }
    });
  };

});

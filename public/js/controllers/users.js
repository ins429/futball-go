app.controller('Users', function MainCtrl($rootScope, $scope, $route, api) {
  $scope.initUser = function(user) {
    $scope.fbLoggedUser = null;
    $rootScope.user = user;
    $scope.showUserMenu = false;
    $scope.userData = {};
    $scope.showUserForm = false;
    $scope.formUrl = '/views/users/form.html';
  }

  $scope.signIn = function() {
    $scope.showUserMenu = false;

    FB.login(function(response) {
      if (response.authResponse) {
        api.fbLogin({
          token: response.authResponse.accessToken
        }).then(function(data) {
          $rootScope.user = data;
        });
      } else {
        // fix me handle error for fb auth
      }
    });

  };

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

  $scope.updateUser = function() {

  }

  $scope.signOut = function() {
    $scope.showUserMenu = false;

    api.logout().then(function(data) {
      if (data && data.status && data.status === 200 ) {
        $rootScope.user = null;
      }
    });
  };

  $scope.toggleUserMenu = function() {
    $scope.showUserMenu = !$scope.showUserMenu;
  };

  $scope.togglerUserForm = function() {
    $scope.showUserForm = !$scope.showUserForm;
  }

});

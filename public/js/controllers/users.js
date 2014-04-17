app.controller('Users', function MainCtrl($rootScope, $scope, $route, api) {
  $scope.initUsers = function() {
    $scope.userForm = {};
  }

  $scope.login = function() {
    api.login({
      email: $scope.userForm.email,
      password: $scope.userForm.password
    }).then(function(data) {
      console.log(data)
      $rootScope.user = data;
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

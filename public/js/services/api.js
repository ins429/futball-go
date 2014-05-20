app.service('api', function($http, $upload) {
  return {
    getCard: function(params) {
      var data = $http({method: 'GET', params: params, url: '/players'}).then(function(resp){
        return resp.data;
      }, function(resp) {
        return resp.data;
      });

      return data;
    },

    fbSignup: function(params) {
      var user = $http({method: 'POST', url: '/fb_signup', data: params}).then(function(resp){
        return resp.data;
      }, function(resp) {
        return resp.data;
      });

      return user;
    },

    fbLogin: function(params) {
      var user = $http({method: 'POST', url: '/fb_login', data: params}).then(function(resp){
        return resp.data;
      }, function(resp) {
        return resp.data;
      });

      return user;
    },

    addCard: function(params) {
      var result = $http({
        method: 'PUT',
        url: '/add_card',
        data: params
      }).then(function(resp){
        return resp.data;
      }, function(resp) {
        return resp.data;
      });

      return result;
    },

    login: function(params) {
      var user = $http({
        method: 'POST',
        url: '/login',
        data: params
      }).then(function(resp){
        return resp.data;
      }, function(resp) {
        return resp.data;
      });

      return user;
    },

    logout: function() {
      var data = $http({method: 'DELETE', url: '/logout'}).then(function(resp){
        return resp.data;
      }, function(resp) {
        return resp.data;
      });

      return data;
    },

    signup: function(params) {
      var user = $http({method: 'POST', url: '/signup', data: params}).then(function(resp){
        return resp.data;
      });

      return user;
    },

    showMe: function() {
      var user = $http({method: 'GET', url: '/showme'}).then(function(resp){
        if (resp && resp.status === 200) {
          return resp.data;
        } else {
          return null;
        }
      });

      return user;
    }
  };

});

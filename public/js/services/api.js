app.service('api', function($http, $upload) {
  return {
    queryImages: function(query) {
      var data = $http({method: 'GET', params: query, url: '/images/query.json'}).then(function(resp){
        return resp.data;
      });

      return data;
    },

    updateImage: function(params) {
      var data = $http({method: 'PUT', params: params, url: '/images/update.json'}).then(function(resp){
        return resp.data;
      });

      return data;
    },

    createImage: function(params, file) {
      var data = $http({method: 'POST', params: params, url: '/images/create.json'}).then(function(resp){
        return resp.data;
      });

      return data;
    },

    createImageWithFile: function(params, file) {
      var data = $upload.upload({
        method: 'POST',
        data: params,
        file: file,
        url: '/images/create.json'
      }).then(function(resp){
        return resp.data;
      });

      return data;
    },

    updateImageWithFile: function(params, file) {
      var data = $upload.upload({
        method: 'PUT',
        data: params,
        file: file,
        url: '/images/update.json'
      }).then(function(resp){
        return resp.data;
      });

      return data;
    },

    deleteImage: function(params) {
      var data = $http({method: 'DELETE', params: params, url: '/images/delete.json'}).then(function(resp){
        return resp.data;
      });

      return data;
    },

    fbLogin: function(params) {
      var user = $http({method: 'POST', url: '/fb_login', data: params}).then(function(resp){
        return resp.data;
      });

      return user;
    },

    login: function(params) {
      var user = $http({method: 'POST', url: '/login', data: params}).then(function(resp){
        return resp.data;
      });

      return user;
    },

    logout: function() {
      var user = $http({method: 'DELETE', url: '/signout'}).then(function(resp){
        return resp.data;
      });

      return user;
    },

    signup: function(params) {
      var user = $http({method: 'POST', url: '/signup', data: params}).then(function(resp){
        return resp.data;
      });

      return user;
    },

    showMe: function() {
      var user = $http({method: 'GET', url: '/showme'}).then(function(resp){
        if (resp && resp.code === 200) {
          return resp.users[0];
        } else {
          return null;
        }
      });

      return user;
    },

    updateUser: function(params) {
      var user = $http({method: 'POST', url: '/users/update.json', data: params}).then(function(resp){
        if (resp && resp.code === 200) {
          return resp.data.users[0];
        } else {
          return null;
        }
      });

      return user;
    }

  };
});

app.service('api', function($http) {
  return {
    getWCCard: function(params) {
      var data = $http({method: 'GET', params: params, url: '/wc_players'}).then(function(resp){
        return resp.data;
      }, function(resp) {
        return resp.data;
      });

      return data;
    }
  };
});

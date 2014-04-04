'use strict';

var app = angular.module('futball-cards', ['ngRoute', 'angularFileUpload']);

app.run(['$rootScope', '$window', function($rootScope, $window) {
  $window.fbAsyncInit = function() {
    FB.init({
      appId      : '739283939424031', // FB App ID dev
      channelUrl : '/channel', // Channel File
      status     : true, // check login status
      cookie     : true, // enable cookies to allow the server to access the session
      xfbml      : true  // parse XFBML
    });
  };

  (function(d){
    var js, id = 'facebook-jssdk', ref = d.getElementsByTagName('script')[0];
    if (d.getElementById(id)) {return;}
    js = d.createElement('script'); js.id = id; js.async = true;
    js.src = "//connect.facebook.net/en_US/all.js";
    ref.parentNode.insertBefore(js, ref);
  }(document));
}]);

app.config(function ($routeProvider, $locationProvider) {
  $routeProvider.when('/',
    {
      templateUrl: '/views/images/index.html',
      controller: 'Images'
    }
  );

  $locationProvider.html5Mode(true);
});

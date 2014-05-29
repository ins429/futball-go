'use strict';

var app = angular.module('futball-cards', ['ngRoute']);

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

  $rootScope.toggleMenu = function($event) {
    var body = angular.element(document.getElementsByTagName('body')[0]),
      target = angular.element($event.target);
    if (target.hasClass('menu-trigger') || (body.hasClass('menu-open') && target.hasClass('pusher')))
      body.toggleClass('menu-open');
  }
}]);

app.config(function ($routeProvider, $locationProvider) {
  $routeProvider.when('/', {
    templateUrl: '/views/cards.html',
    controller: 'Cards'
  });
  
  $locationProvider.html5Mode(true);
});

// directives
app.directive('ngEnter', function() {
  return function (scope, element, attrs) {
    element.bind("keydown keypress", function (event) {
      if(event.which === 13) {
        scope.$apply(function (){
          scope.$eval(attrs.ngEnter);
        });
          event.preventDefault();
      }
    });
  };
});

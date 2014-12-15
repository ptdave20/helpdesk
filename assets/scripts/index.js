var indexApp = angular.module('helpIndex',['ngRoute'])

indexApp.controller('bCtrl', function ($scope,$http) {
	$http.get('/o/user/logged_in',{withCredentials:true}).success(function(data) {
		if(data["result"]) {
			window.location = "/main.html";
		} else {
			// In case we want to directly move them
		}
	});
});

indexApp.controller('indexRoute', function($scope) {
	$scope.header = "Welcome to the Helpdesk";
	$scope.text = "In order to login, please click the login button below. You will be authenticated using your google credentials. After, you will be redirected to the helpdesk once you are authenicated.";
})

indexApp.controller('invalidDomain', function($scope) {
	$scope.header = "Error";
	$scope.text = "Unfortunately, the domain associated with your Google account is not enabled to be used here. Please try again with a different account.";
});
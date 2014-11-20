var indexApp = angular.module('helpIndex',[])

indexApp.controller('bCtrl', function ($scope,$http) {
	$http.get('/o/user/logged_in',{withCredentials:true}).success(function(data) {
		var j = angular.fromJson(data);
		console.log(j["result"]);
		if(j["result"]) {
			window.location = "/main.html";
		} else {
			// In case we want to directly move them
		}
	});
});

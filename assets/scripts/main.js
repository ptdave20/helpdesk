var indexApp = angular.module('helpIndex',['ui.bootstrap.modal'])

indexApp.controller('bCtrl', function ($scope,$http,$modal) {
	$http.get('/o/user/me',{withCredentials:true}).success(function(data) {
		var j = angular.fromJson(data);
		$scope.me = j;
	});


	$scope.newTicket = function() {
		var modalInstance = $modal.open({
			templateUrl: 'ticketModal.html'
		});

		modalInstance.result.then(function($data) {

		});
	}
});
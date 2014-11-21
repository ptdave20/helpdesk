angular.module('helpIndex',['ui.bootstrap.tpls', 'ui.bootstrap.modal','ui.bootstrap.alert'])



angular.module('helpIndex').controller('bCtrl', function ($scope,$http,$modal) {
	$http.get('/o/user/me',{withCredentials:true}).success(function(data) {
		var j = angular.fromJson(data);
		$scope.me = j;
	});


	$scope.newTicket = function() {
		var modalInstance = $modal.open({
			templateUrl: 'ticketModal.html',
			controller: 'newTicketCtrl',
			backdrop: 'static'
		});

		modalInstance.result.then(function($data) {

		});
	}
});

angular.module('helpIndex').controller('newTicketCtrl', function($scope,$modalInstance) {
	$scope.submit = function() {
		if($scope.description==null || $scope.description.trim().length == 0) {
			console.log("description is empty")
		}
		$modalInstance.close();
	}

	$scope.cancel = function() {
		$modalInstance.dismiss();
	}
});
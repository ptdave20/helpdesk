angular.module('helpIndex',['ui.bootstrap.tpls', 'ui.bootstrap.modal','ui.bootstrap.alert'])



angular.module('helpIndex').controller('bCtrl', function ($scope,$http,$modal) {
	$scope.submitted = [];
	$scope.assigned = [];
	$scope.department = [];
	$scope.departments = [];

	$http.get('/o/user/me',{withCredentials:true}).success(function(data) {
		var j = angular.fromJson(data);
		$scope.me = j;
	});

	$http.get('/o/ticket/list/mine',{withCredentials:true}).success(function(data) {
		var j = angular.fromJson(data);
		$scope.submitted = j;
	});

	$http.get('/o/departments/list',{withCredentials:true}).success(function(data) {
		var j = angular.fromJson(data);
		$scope.departments = j;
	});

	$scope.newTicket = function() {
		var modalInstance = $modal.open({
			templateUrl: 'ticketModal.html',
			controller: 'newTicketCtrl',
			backdrop: 'static'
		});

		modalInstance.result.then(function(data) {
			if(data) {
				// Reget our list of tickets
				$http.get('/o/ticket/list/mine',{withCredentials:true}).success(function(data) {
					var j = angular.fromJson(data);
					$scope.submitted = j;
				});
			}
		});
	}
});

angular.module('helpIndex').filter('depCat', function() {
	return function(data,id) {
		data = data || [];
		var out = [];

		for(var i=0; i<data.length; i++) {
			if(data[i].Id==id) {
				out = data[i].Category;
			}
		}

		return out;
	}
});


angular.module('helpIndex').filter('depFilter', function() {
	return function(id,data) {
		for(var i=0; i<data.length; i++) {
			if(data[i].Id == id) {
				return data[i].Name;
			}
		}
		return "Unknown";
	}
});

angular.module('helpIndex').controller('newTicketCtrl', function($scope,$http,$modalInstance) {
	$http.get('/o/departments/list',
			{
				withCredentials:true
			})
	.success(function(data) {
		$scope.deps = angular.fromJson(data);
	});

	$scope.submit = function() {
		$http.post(
			'/o/ticket/insert',
			$scope.tkt,
			{
				withCredentials:true,
				headers: { 
					'Content-Type': 'application/json',
				}
			}
		).success(function(data) {
			j = angular.fromJson(data);

			if(j["Id"]!=null || j["Id"]!=undefined) {
				$modalInstance.close(true);	
			}
		});
	}

	$scope.cancel = function() {
		$modalInstance.dismiss();
	}
});
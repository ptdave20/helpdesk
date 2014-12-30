angular.module('helpIndex').controller('userModal',['$scope','$http','$modalInstance', function($scope,$http,$modalInstance) {
	$scope.buildings = [];
	$scope.user = {};
	$http.get('/o/user/me',{withCredentials:true}).success(function(data) {
		$scope.user = data;
	});

	$http.get('/o/domain/buildings', {withCredentials:true}).success(function(data) {
		if(data!=null || data.length == 0)
			$scope.buildings = data;
	});

	$scope.saveProfile = function() {
		if($scope.user.Firstname=="" || $scope.user.Firstname==null) {
			return;
		}
		if($scope.user.Lastname=="" || $scope.user.Lastname==null) {
			return;
		}
		if($scope.user.Building=="" || $scope.user.Building==null) {
			return;
		}

		$http.post('/o/user/me', $scope.user, {withCredentials:true}).success(function(data) {
			if(data == "success") {
				$modalInstance.close(true);
			}
		});
	}
	$scope.cancel = function() {
		if(!$scope.user.NewUser) {
			$modalInstance.close(true);
		}
	}
}]);
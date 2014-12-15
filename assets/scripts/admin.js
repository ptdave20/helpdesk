var admin = angular.module('adminIndex',['ngRoute','ui.bootstrap.tpls', 'ui.bootstrap.modal','ui.bootstrap']);

admin.controller('bCtrl', function ($scope,$http) {
	$scope.isCollapsed = true;
})

admin.controller('depCtrl', function($scope,$http) {
	$scope.newDep = {};
	$scope.departments = [];

	$scope.getDepartments = function() {
		$http.get('/o/department/list',{withCredentials:true}).success(function(data) {
			$scope.departments = data;
		});
	}

	$scope.addDepartment = function() {
		$http.post('/o/department',$scope.newDep,{withCredentials:true}).success(function(data) {
			$scope.getDepartments();
			$scope.newDep = {};
		});
		//$scope.departments.push($scope.newDep);
		
	}

	$scope.getDepartments();
})
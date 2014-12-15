var admin = angular.module('adminIndex',['ngRoute','ui.bootstrap.tpls', 'ui.bootstrap.modal','ui.bootstrap']);

admin.controller('bCtrl', function ($scope,$http) {
	$scope.isCollapsed = true;
})

admin.controller('depCtrl', function($scope,$http) {
	$scope.newDep = {};
	$scope.departments = [];
	$scope.addDepartment = function() {
		$scope.departments.push($scope.newDep);
		$scope.newDep = {};
	}
})
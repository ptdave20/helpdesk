var admin = angular.module('adminIndex',['ngRoute','ui.bootstrap.tpls', 'ui.bootstrap.modal','ui.bootstrap']);

admin.controller('bCtrl', function ($scope,$http) {
	$scope.isCollapsed = true;
})

admin.controller('depCtrl', function($scope,$http) {
	$scope.newDep = {};
	$scope.departments = [];
	$scope.selDep = null;

	$scope.selectDep = function(id) {
		for(var i=0; i<$scope.departments.length; i++) {
			if($scope.departments[i].Id == id) 
				$scope.selDep = $scope.departments[i];
		}
	}

	$scope.getDepartments = function() {
		$http.get('/o/department/list',{withCredentials:true}).success(function(data) {
			$scope.departments = data;
			$scope.selectDep = $scope.departments[0];
		});
	}

	$scope.addDepartment = function() {
		$http.post('/o/department',$scope.newDep,{withCredentials:true}).success(function(data) {
			$scope.getDepartments();
			$scope.newDep = {};
		});
	}

	$scope.addCategory = function(depId, cat) {
		var cat_data = {
			Name : cat
		}
		$http.post('/o/department/'+depId,cat_data,{withCredentials:true}).success(function(data) {
			$scope.getDepartments();
		});
	}

	$scope.delDepartment = function(depId) {
		$http.delete('/o/department/'+depId, {withCredentials:true}).success(function(data) {
			if(data == "success") {
				$scope.getDepartments();
			}
		})
	}

	$scope.delCategory = function(depId,catId) {
		$http.delete('/o/department/'+depId+"/"+catId, {withCredentials:true}).success(function(data) {
			if(data == "success") {
				$scope.getDepartments();
			}
		})
	}

	$scope.getDepartments();
})

admin.controller('depConfigCtrl', ['$scope','$routeParams','$http', function($scope,$routeParams,$http){
	var id = $routeParams.id;
	$http.get('/o/department/'+id,{withCredentials:true}).success(function(data) {
		$scope.dep = data;
	});
}]);

admin.controller('bldgCtrl', ['$scope','$http', function($scope,$http) {
	$scope.newBldg = {}
	$scope.addBuilding = function() {
		$http.post('/o/domain/building',$scope.newBldg, {withCredentials:true}).success(function(data) {
			$scope.newBldg = {};
			$scope.getBuildings();
		});
	}
	$scope.getBuildings = function() {
		$http.get('/o/domain/buildings', {withCredentials:true}).success(function(data) {
			$scope.buildings = data;
		});
	}
	$scope.getBuildings();
}]);
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
		};
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
});

admin.controller('depConfigCtrl', ['$scope','$routeParams','$http','$location','$q', function($scope,$routeParams,$http,$location,$q){
	var id = $routeParams.id;

	$scope.newCats = [];

	$scope.nameChange = false;
	$scope.buildingChange = false;
	$scope.buildingSpecChange = false;
	
	$scope.getBuildings = function() {
		$http.get('/o/domain/buildings', {withCredentials:true}).success(function(data) {
			$scope.buildings = data;
			$scope.getDepartment();
		});
	}

	$scope.getDepartment = function() {
		$http.get('/o/department/'+id,{withCredentials:true}).success(function(data) {
			$scope.dep = data;
		});
	}


	$scope.noChanges = function() {
		if($scope.nameChange || $scope.buildingChange || $scope.buildingSpecChange) {
			return false;
		}
		if($scope.dep!=undefined && $scope.dep!=null && $scope.dep.Category!=undefined && $scope.dep.Category!=null) {
			for(var i=0; i<$scope.dep.Category.length; i++) {
				if($scope.dep.Category[i].Add || $scope.dep.Category[i].Remove) {
					return false;
				}
			}	
		}
		
		return true;
	}

	$scope.save = function() {
		if($scope.noChanges()) {
			return;
		}
		var upd = {
			Name: $scope.dep.Name,
			IsBuildingSpecific: $scope.dep.IsBuildingSpecific,
			Building: $scope.dep.Building
		}
		$http.post('/o/department/'+id,upd,{withCredentials:true}).success(function(data) {
			if(data == "success") {
				$scope.nameChange = false;
				$scope.buildingChange = false;
				$scope.buildingSpecChange = false;
				if($scope.noChanges()) {
					$location.path("/departments");
				} else {
					// We have categories to add or remove
					var p = [];
					for(var i=0; i<$scope.dep.Category.length; i++) {
						var cat = $scope.dep.Category[i];
						if(cat.Add) {
							// We need to add
							p.push($http.post('/o/department/'+id+'/category',cat,{withCredentials:true}));
						} else {
							if(cat.Remove) {
								p.push($http.delete('/o/department/'+id+'/'+cat.Id,cat,{withCredentials:true}));
							}
						}
					}
					console.log(p.length);
					$q.all(p).then(function(data) {
						var s = true;
						for(var i=0; i<data.length; i++)
							if(data[i].data!="success")
								s = false;

						if(s)
							$location.path("/departments");

					});

				}
			}
		});
	}

	$scope.addCat = function() {
		var newCat = {
			Name: $scope.newCatName,
			Add: true
		}
		$scope.dep.Category = $scope.dep.Category || [];
		$scope.dep.Category.push(newCat);
		$scope.newCatName = "";
		$scope.catAddChange = true;
	}

	$scope.delCat = function(cat) {
		var index = $scope.dep.Category.indexOf(cat);
		if(cat.Add) {
			if(index > -1) {
				$scope.dep.Category.splice(index,1);
			} else {
				console.log("Can't find!");
			}
		} else {
			$scope.dep.Category[index].Remove = true;
		}

	}

	$scope.getBuildings();
}]);

admin.controller('bldgCtrl', ['$scope','$http', function($scope,$http) {
	$scope.newBldg = {}
	$scope.addBuilding = function() {
		$http.post('/o/domain/building',$scope.newBldg, {withCredentials:true}).success(function(data) {
			$scope.newBldg = {}
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
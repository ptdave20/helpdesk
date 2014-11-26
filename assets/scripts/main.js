angular.module('helpIndex',['ui.bootstrap.tpls', 'ui.bootstrap.modal','ui.bootstrap'])



angular.module('helpIndex').controller('bCtrl', function ($scope,$http,$modal,$interval) {
	$scope.submitted = [];
	$scope.assigned = [];
	$scope.department = [];
	$scope.departments = [];
	$scope.ticketSubmitterStatus = "open";
	$scope.ticketDepartmentStatus = "open";
	$scope.submittedHasTickets = false;
	$scope.assignedHasTickets = false;
	$scope.departmentHasTickets = false;


	$scope.selDepartment = null;

	$http.get('/o/user/me',{withCredentials:true}).success(function(data) {
		var j = angular.fromJson(data);
		$scope.me = j;

		$scope.me.Department = $scope.me.Department || [];

		if($scope.me.Department!=null || $scope.me.Department.length == 0) {
			$scope.selDepartment = null;
		} else {
			$scope.selDepartment = $scope.me.Department[0];
		}
		$scope.getSubmittedTickets();
		$scope.getDepTickets();

	});

	$scope.openTicket = function(id) {
		console.log(id);
	}

	$scope.viewSubmitterOpen = function() {
		$scope.ticketSubmitterStatus = "open";
		$scope.getSubmittedTickets();
	}

	$scope.viewSubmitterClosed = function() {
		$scope.ticketSubmitterStatus = "closed";
		$scope.getSubmittedTickets();
	}

	$scope.viewDepartmentOpen = function() {
		$scope.ticketDepartmentStatus = "open";
		$scope.getDepTickets();
	}

	$scope.viewDepartmentClosed = function() {
		$scope.ticketDepartmentStatus = "closed";
		$scope.getDepTickets();
	}

	$scope.getSubmittedTickets = function() {
		$http.get('/o/ticket/list/mine/'+$scope.ticketSubmitterStatus,{withCredentials:true}).success(function(data) {
			var j = angular.fromJson(data);
			j = j || [];
			$scope.submitted = j;
			if($scope.submitted.length > 0) {
				$scope.submittedHasTickets = true;
			} else {
				$scope.submittedHasTickets = false;
			}
		});
	}

	$scope.closeTicket = function(id) {
		$http.post('/o/ticket/close/'+id,{withCredentials:true}).success(function(data){
			var j = angular.fromJson(data);
			if(j["result"]) {
				// Success!
			}
		});
	}

	$scope.getDepTickets = function() {
		// If we don't have a department, then return
		if($scope.selDepartment == undefined || $scope.selDepartment == null) {
			return;
		}
		$http.get('/o/ticket/list/department/'+$scope.selDepartment+"/"+$scope.ticketDepartmentStatus,{withCredentials:true}).success(function(data) {
			var j = angular.fromJson(data);
			$scope.department = j;
			if($scope.department.length > 0) {
				$scope.departmentHasTickets = true;
			} else {
				$scope.departmentHasTickets = false;
			}

		});
	}
	

	$http.get('/o/departments/list',{withCredentials:true}).success(function(data) {
		var j = angular.fromJson(data);
		$scope.departments = j;
	});

	$scope.bg = $interval(function() {
		// Reget your submitted tickets
		$scope.getSubmittedTickets()

		// Reget your department tickets
		$scope.getDepTickets();
	}, 30000);

	

	$scope.newTicket = function() {
		var modalInstance = $modal.open({
			templateUrl: 'ticketSubmitModal.html',
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

angular.module('helpIndex').filter('depName', function() {
	return function(input,scope) {
		if(scope.departments == undefined)
			return "Unknown";
		for(var i=0; i<scope.departments.length; i++) {
			if(scope.departments[i].Id == input)
				return scope.departments[i].Name;
		}
		return "Unknown"
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

angular.module('helpIndex').filter('search', function() {
	return function(input,data) {
		input = input || "";
		data = data || [];

		if(input == "") {
			return data;
		}

		var out = [];
		for(var i = 0; i<data.length; i++) {
			if(data[i].Subject.indexOf(input) > -1) {
				out.append(data[i]);
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

angular.module('helpIndex').controller('ticketModal', function($scope,$http,$modalInstance) {
	$http.get('/o/departments/list',
			{
				withCredentials:true
			})
	.success(function(data) {
		$scope.deps = angular.fromJson(data);
	});
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
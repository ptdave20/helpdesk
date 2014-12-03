var helpdesk = angular.module('helpIndex',['ngRoute','ui.bootstrap.tpls', 'ui.bootstrap.modal','ui.bootstrap'])

helpdesk.config(function($routeProvider) {
	$routeProvider
	.when('/', {
		template: '<h1>Welcome to the helpdesk</h1>'
	})
	.when('/tickets/mine', {
		templateUrl: 'ticketlist.html',
		controller: 'myTicketListCtrl'
	})
	.when('/tickets/department', {
		templateUrl: 'ticketList.html',
		controller: 'depTicketListCtrl'
	})
	.when('/tickets/assigned', {
		templateUrl: 'ticketList.html',
		controller: 'assignedTicketListCtrl'
	})
	.when('/ticket/:id', {
		templateUrl: 'ticket.html'
	})
});



angular.module('helpIndex').controller('bCtrl', function ($scope,$http,$modal,$interval,$location) {
	$scope.submitted = [];
	$scope.assigned = [];
	$scope.departments = [];

	$scope.options = {
		ticket: {
			hasTickets: false,
			selectable : false,
			showDepartment: true,
			showCategory: false,
			showAssignedTo: false,
			showSubmitter: false,
			status: "open",
			hasTickets: false
		},
	};
	

	$scope.isActive=function(route) {
		if(route === '/') {
			return $location.path() === '/';
		}
		return $location.path().indexOf(route) != -1;
	}

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
	});

	$scope.openTicket = function(ticketData) {
		var modalInstance = $modal.open({
			templateUrl: 'ticketViewModal.html',
			controller: 'ticketModal',
			backdrop: 'static',
			resolve: {
				ticket: function() {
					return ticketData;
				},
				departments: function() {
					return $scope.departments;
				},
				options: function() {
					return {
						newTicket: false,
						canEdit: ticketData.Status != "closed",
						canClose: ticketData.Status != "closed"
					}
				}
			}
		});

		modalInstance.result.then(function(data) {
			if(data) {
				// Reget our list of tickets
				$scope.getSubmittedTickets();
				$scope.getDepTickets();
			}
		});
	}

	$scope.newTicket = function() {
		var modalInstance = $modal.open({
			templateUrl: 'ticketViewModal.html',
			controller: 'ticketModal',
			backdrop: 'static',
			resolve: {
				ticket: function() {
					return {};
				},
				departments: function() {
					return $scope.departments;
				},
				options: function() {
					return {
						newTicket: true,
						canEdit: true,
						canClose: false
					}
				}
			}
		});

		modalInstance.result.then(function(data) {
			if(data) {
				// Reget our list of tickets
				$scope.getSubmittedTickets();
				$scope.getDepTickets();
			}
		});
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

	$scope.getDepTickets = function() {
		// If we don't have a department, then return
		if($scope.selDepartment == undefined || $scope.selDepartment == null) {
			$scope.departments = [];
			return;
		}
		$http.get('/o/ticket/list/department/'+$scope.selDepartment+"/"+$scope.ticketDepartmentStatus,{withCredentials:true}).success(function(data) {
			$scope.departments[$scope.selDepartment] = data;
			if(data.length > 0) {
				$scope.departmentHasTickets = true;
			} else {
				$scope.departmentHasTickets = false;
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

	
	

	$http.get('/o/departments/list',{withCredentials:true}).success(function(data) {
		var j = angular.fromJson(data);
		$scope.departments = j;
	});

	


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

angular.module('helpIndex').controller('ticketModal', function($scope,$http,$modalInstance,ticket,departments,options) {
	$scope.ticket = ticket;
	$scope.departments = departments;
	$scope.categories = [];
	$scope.options = options;

	$scope.DepCatChange = function() {
		for(var d=0; d<$scope.departments.length; d++) {
			if($scope.departments[d].Id == $scope.ticket.Department) {
				// Set the uneditable value
				$scope.departmentName = $scope.departments[d].Name;
				
				$scope.departments[d].Category = $scope.departments[d].Category || [];

				// Set the editable categories
				$scope.categories = $scope.departments[d].Category;

				// Set the uneditable category value
				for(var c=0; c<$scope.departments[d].Category.length; c++) {
					if($scope.departments[d].Category[c].Id == $scope.ticket.Category) {
						$scope.categoryName = $scope.departments[d].Category[c].Name;
					}
				}
			}
		}
	}

	$scope.update = function() {

	}

	$scope.submit = function() {
		if(!$scope.options.newTicket)
			return;
		$http.post(
			'/o/ticket/insert',
			$scope.ticket,
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
	$scope.DepCatChange();
	
});


helpdesk.controller('myTicketListCtrl', ['$scope','$http','$interval', function($scope,$http,$interval) {
	$scope.options = {
		ticket: {
			selectable : false,
			showDepartment: true,
			showCategory: false,
			showAssignedTo: false,
			showSubmitter: false,
			status: "open",
			hasTickets: false
		},
	};

	$scope.getTickets = function() {
		$http.get('/o/ticket/list/mine/'+$scope.options.ticket.status,{withCredentials:true}).success(function(data) {
			$scope.tickets = data;
		});
	}
	$scope.viewOpenTickets = function() {
		$scope.options.ticket.status = "open";
		$scope.getTickets();
	}

	$scope.viewClosedTickets = function() {
		$scope.options.ticket.status = "closed";
		$scope.getTickets();
	}

	$scope.options = $scope.$parent.options;
	$scope.departments = $scope.departments;
	$scope.bg = $interval(function() {
		$scope.getTickets();
	}, 30000);


	$scope.getTickets();
	

	console.log($scope);
}]);

helpdesk.controller('depTicketListCtrl', ['$scope','$http', function($scope,$http) {
	$scope.options = {
		ticket: {
			selectable : false,
			showDepartment: true,
			showCategory: false,
			showAssignedTo: false,
			showSubmitter: false,
			status: "open",
			hasTickets: false
		},
	};

	$scope.departments = $scope.departments;
	$scope.tickets = $scope.$parent.mine;
}]);

helpdesk.controller('assignedTicketListCtrl', ['$scope','$http', function($scope,$http) {
	$scope.options = {
		ticket: {
			selectable : false,
			showDepartment: true,
			showCategory: false,
			showAssignedTo: false,
			showSubmitter: false,
			status: "open",
			hasTickets: false
		},
	};
	$scope.departments = $scope.departments;
	$scope.tickets = $scope.$parent.mine;
}]);
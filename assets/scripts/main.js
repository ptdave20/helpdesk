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

helpdesk.factory('Tickets', function() {
	return {
		mine: {
			open:[],
			closed:[],
			status: "open",
			lastOpenCount: -1,
			currentOpenCount: -1,
			getTickets: function() {}
		},
		departments: {
			areas: [],
			status: "open",
			lastOpenCount: -1,
			currentOpenCount: -1,
			activeDepartment: "",
			available:[],
			getTickets: function() {}
		},
		assigned: {
			open:[],
			closed:[],
			status: "open",
			lastOpenCount: -1,
			currentOpenCount: -1,
			getTickets: function() {}
		}
	}
});
helpdesk.factory('Departments', function() {
	return [];
});
helpdesk.factory('Me', function() {
	return {};
});


angular.module('helpIndex').controller('bCtrl', function ($scope,$http,$modal,$interval,$location,Tickets,Departments,Me) {
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
	
	$scope.Tickets = Tickets;

	Tickets.mine.getTickets = function() {
		$http.get('/o/ticket/list/mine/'+Tickets.mine.status,{withCredentials:true}).success(function(data) {
			switch(Tickets.mine.status) {
				case "open":
					Tickets.mine.open = data;
					Tickets.mine.currentOpenCount = data.length;
					break;
				case "closed":
					Tickets.mine.closed = data;
					break;
			}
		});
	}

	Tickets.departments.getTickets = function() {
		// If we don't have a department, then return
		angular.forEach(Tickets.departments.available, function(depValue,depKey) {
			angular.forEach(["open","closed"], function(statValue, statKey) {
				$http.get('/o/ticket/list/department/'+depValue+"/"+statValue,{withCredentials:true}).success(function(data) {
					if(Tickets.departments.areas[depValue] == undefined || Tickets.departments[depValue] == null) {
						Tickets.departments.areas[depValue] = {
							open:[],
							closed:[]
						}
					}
					if(data!=null) {

						Tickets.departments.areas[depValue][statValue] = data;
						//console.log(Tickets.departments.areas[depValue][statValue]);
					}
				});
			});
			
		});

		
	}

	Tickets.mine.getTickets();
	Tickets.departments.getTickets();

	$scope.bg = $interval(function() {
		Tickets.mine.getTickets();
		Tickets.departments.getTickets();
	}, 10000);

	$scope.isActive=function(route) {
		if(route === '/') {
			return $location.path() === '/';
		}
		return $location.path().indexOf(route) != -1;
	}

	$http.get('/o/user/me',{withCredentials:true}).success(function(data) {
		angular.forEach(data.Department, function(value,key) {
			Tickets.departments.available.push(value);
			Tickets.departments.areas[value] = {
				open:[],
				closed:[]
			};

			if(Tickets.departments.activeDepartment == "" || Tickets.departments.activeDepartment == undefined)
				Tickets.departments.activeDepartment = value;
		});
		Me = data;
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


helpdesk.controller('myTicketListCtrl', ['$scope','$http','Tickets', function($scope,$http,Tickets) {
	$scope.options = {
		ticket: {
			selectDepartment: false,
			selectable : false,
			showDepartment: true,
			showCategory: false,
			showAssignedTo: false,
			showSubmitter: false,
			status: "open",
			hasTickets: false
		},
	};

	$scope.Mine = Tickets.mine;
	$scope.status = $scope.Mine.status;
	$scope.tickets = $scope.Mine[$scope.status];
	
	$scope.availDepartments = Tickets.departments.available || [];

	$scope.viewOpenTickets = function() {
		$scope.status = "open";
		$scope.tickets = $scope.Mine[$scope.status];
		//$scope.Mine.getTickets();
	}

	$scope.viewClosedTickets = function() {
		$scope.status = "closed";
		$scope.tickets = $scope.Mine[$scope.status];
		//$scope.Mine.getTickets();
	}
}]);

helpdesk.controller('depTicketListCtrl', ['$scope','$http','Tickets', function($scope,$http,Tickets) {
	$scope.options = {
		ticket: {
			selectDepartment: true,
			selectable : false,
			showDepartment: false,
			showCategory: false,
			showAssignedTo: false,
			showSubmitter: false,
			status: "open",
			hasTickets: false
		},
	};

	$scope.Departments = Tickets.departments;

	$scope.status = $scope.Departments.status;
	$scope.activeDepartment = $scope.Departments.activeDepartment;
	$scope.availDepartments = $scope.Departments.available || [];
	$scope.tickets = $scope.Departments.areas[$scope.activeDepartment]["open"];
	Tickets.departments.getTickets();
	$scope.setDepartment = function(v) {
		$scope.activeDepartment = v;
	}
	$scope.viewOpenTickets = function() {
		$scope.status = "open";
		$scope.tickets = $scope.Departments.areas[$scope.activeDepartment][$scope.status];
		Tickets.departments.getTickets();
		console.log(Tickets.departments);
		//$scope.Tickets.getTickets();
	}

	$scope.viewClosedTickets = function() {
		$scope.status = "closed";
		$scope.tickets = $scope.Departments.areas[$scope.activeDepartment][$scope.status];
		Tickets.departments.getTickets();
		console.log(Tickets.departments);
		//$scope.Tickets.getTickets();
	}
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
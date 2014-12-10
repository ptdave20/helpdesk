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

helpdesk.factory('Tickets', function($http) {
	var obj = {
		mine: {
			open:[],
			closed:[],
			status: "open",
			lastOpenCount: -1,
			currentOpenCount: -1,
			getTickets: function() {}
		},
		departments: {
			open: [],
			closed: [],
			lastOpenCount: -1,
			currentOpenCount: -1,
			activeDepartment: null,
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
	};

	obj.mine.getTickets = function() {
		$http.get('/o/ticket/list/mine/'+obj.mine.status,{withCredentials:true}).success(function(data) {
			switch(obj.mine.status) {
				case "open":
					obj.mine.open = data;
					obj.mine.currentOpenCount = data.length;
					break;
				case "closed":
					obj.mine.closed = data;
					break;
			}
		});
	}

	obj.departments.getTickets = function() {
		// If we don't have a department, then return
		if(obj.departments.activeDepartment==null)
			return;
		angular.forEach(["open"], function(stat, key) {
			$http.get('/o/ticket/list/department/'+obj.departments.activeDepartment+"/"+stat,{withCredentials:true}).success(function(data) {
				while(obj.departments[stat].length > 0) {
					obj.departments[stat].pop();
				}
					
				if(data!=null) {
					for(var i=0; i<data.length; i++) {
						obj.departments[stat].push(data[i]);
					}
					obj.departments[stat] = data;
				}
			});
		});
		
	}

	return  obj;
});
helpdesk.factory('DepartmentsList', function($http) {
	var ret = [];

	$http.get('/o/departments/list',{withCredentials:true}).success(function(data) {
		while(ret.length > 0)
			ret.pop();
		for(var i=0; i<data.length; i++)
			ret.push(data[i]);
		//$scope.departments = j;
	});

	return ret;
});

angular.module('helpIndex').controller('bCtrl', function ($scope,$http,$modal,$interval,$location,Tickets,DepartmentsList) {
	$scope.Tickets = Tickets;

	$scope.Tickets.mine.getTickets();
	$scope.Tickets.departments.getTickets();

	$scope.updateTask = $interval(function() {
		$scope.Tickets.mine.getTickets();
		$scope.Tickets.departments.getTickets();
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

			if(Tickets.departments.activeDepartment == "" || Tickets.departments.activeDepartment == undefined)
				Tickets.departments.activeDepartment = value;
		});
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
				$scope.Tickets.departments.getTickets();
				$scope.Tickets.mine.getTickets();
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
				$scope.Tickets.departments.getTickets();
				$scope.Tickets.mine.getTickets();
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

angular.module('helpIndex').filter('catName', function() {
	return function(input,scope) {
		if(scope.departments == undefined || input=="")
			return "Unknown";
		for(var i=0; i<scope.departments.length; i++) {
			for(var c=0; c<scope.departments[i].Category.length; c++) {
				if(scope.departments[i].Category[c].Id == input)
					return scope.departments[i].Category[c].Name;
			}
		}
		return "Unknown"
	}
});

angular.module('helpIndex').filter('paginate', function() {
	return function(data, start, finish) {
		var out = [];

		if(start < data.length)
			start = 0;

		if(finish > data.length)
			finish = data.length;


		for(var i=start; i<finish; i++) {
			out.push(data[i]);
		}

		return out;
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
			hasTickets: false,
			order: 'date',
			orderReverse: false
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
	$scope.setOrder = function(value) {
		if($scope.order == value) {
			$scope.order = "-"+$scope.order;
		} else {
			$scope.order = value;
		}
	}
}]);

helpdesk.controller('depTicketListCtrl', ['$scope','$http','Tickets','DepartmentsList', function($scope,$http,Tickets,DepartmentsList) {
	$scope.options = {
		ticket: {
			selectDepartment: true,
			selectable : false,
			showDepartment: false,
			showCategory: true,
			showAssignedTo: false,
			showSubmitter: false,
			status: "open",
			hasTickets: false,
		},
	};

	$scope.Service = Tickets;
	$scope.departments = DepartmentsList;

	$scope.status = $scope.Service.departments.status;
	$scope.activeDepartment = $scope.Service.departments.activeDepartment;
	$scope.availDepartments = $scope.Service.departments.available || [];
	$scope.Service.departments.getTickets();

	$scope.tickets = $scope.Service.departments.open;
	$scope.setDepartment = function(v) {
		$scope.activeDepartment = v;
	}
	$scope.viewOpenTickets = function() {
		$scope.tickets = $scope.Service.departments.open;
	}

	$scope.viewClosedTickets = function() {
		$scope.tickets = $scope.Service.departments.closed;
	}

	$scope.setOrder = function(value) {
		if($scope.order == value) {
			if($scope.reverse == undefined)
				$scope.reverse = false;
			$scope.reverse = !$scope.reverse;
		} else {
			$scope.order = value;
			$scope.reverse = false;
		}
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
			hasTickets: false,
			order: 'date',
			orderReverse: false
		},
	};
	$scope.departments = $scope.departments;
	$scope.tickets = $scope.$parent.mine;
}]);
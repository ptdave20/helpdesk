var helpdesk = angular.module('helpIndex',['ngRoute','ui.bootstrap.tpls', 'ui.bootstrap.modal','ui.bootstrap']);


helpdesk.service('DepartmentsService', function($http) {
	function DeptListService() {
		this.getDepartments = function() {
			return $http.get('/o/departments/list',{withCredentials:true});
		}
	}
	var obj = new DeptListService();
	return obj;
});

helpdesk.service('DeptTickets', function($http) {
	var obj = {};
	obj._department = "";
	obj._status = "open";
	obj.setDepartment = function(id) {
		obj._department = id;
	}
	obj.getTickets = function() {
		var promise;
		if(obj._department == "" || obj._department == null || obj._department == undefined)
			return;
		$http.get("/o/tickets/departments/"+obj._department+"/"+obj._status,{withCredentials:true})
		return promise;
	}
	obj.setStatus = function(status) {
		obj._status = status;
	}
	return obj;
});
helpdesk.service('MyTickets', function($http) {
	var obj = {};
	obj._status = "open";
	obj.getTickets = function() {
		var promise;
		promise = $http.get("/o/tickets/submitted/"+obj._status,{withCredentials:true});
		return promise;
	}
	obj.setStatus = function(status) {
		obj._status = status;
	}
	return obj;
});

helpdesk.service('TicketService', ['$http',function($http) {
	function TicketService() {
		this.Get = function(id) {
			console.log(id);
			return $http.get('/o/ticket/'+id,{withCredentials:true});
		}
		this.Update = function(ticket) {
			console.log(ticket);
		}
		this.Close = function(ticket) {
			console.log(ticket);
		}
		this.Assign = function(ticket,user_id) {
			console.log(ticket);
		}
		this.AddNote = function(ticket,private,data) {
			console.log(ticket);
		}
		this.Create = function(ticket) {
			console.log(ticket);
		}
	}

	return new TicketService();
}]);

angular.module('helpIndex').controller('bCtrl', function ($scope,$http,$modal,$interval,$location) {
	$scope.isCollapsed = true;

	$scope.isActive=function(route) {
		if(route === '/') {
			return $location.path() === '/';
		}
		return $location.path().indexOf(route) != -1;
	}

	$scope.openTicket = function(ticketId) {
		var modalInstance = $modal.open({
			templateUrl: 'ticketViewModal.html',
			controller: 'ticketModal',
			backdrop: 'static',
			resolve: {
				departments: function() {
					return $scope.departments;
				},
				options: function() {
					return {
						newTicket: false,
						ticketId: ticketId
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

angular.module('helpIndex').controller('ticketModal', ['$scope','$http','$modalInstance','TicketService', 'DepartmentsService', 'options', function($scope,$http,$modalInstance,TicketService,DepartmentsService,options) {
	$scope.ticket = {};
	$scope.departments = []
	$scope.categories = [];
	$scope.options = options;

	DepartmentsService.getDepartments().then(function(data) { $scope.departments = data.data; });
	TicketService.Get($scope.options.ticketId).then(function(data) { $scope.ticket = data.data; });
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
	
}]);


helpdesk.controller('myTicketListCtrl', ['$scope','$http','$modal','MyTickets','TicketService', function($scope,$http,$modal,MyTickets,TicketService) {
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

	$scope.Service = MyTickets;
	$scope.Service.getTickets().then(function(data) {
		$scope.tickets = data.data;
	});
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

	$scope.openTicket = function(ticketId) {
		var modalInstance = $modal.open({
			templateUrl: 'ticketViewModal.html',
			controller: 'ticketModal',
			backdrop: 'static',
			resolve: {
				departments: function() {
					return $scope.departments;
				},
				options: function() {
					return {
						newTicket: false,
						ticketId: ticketId
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
}]);

helpdesk.controller('depTicketListCtrl', ['$scope','$http','DeptTickets','DepartmentsList', function($scope,$http,DeptTickets,DepartmentsList) {
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

helpdesk.controller('ticketViewCtrl', ['$scope','$http', function($scope,$http) {
	
}]);
angular.module('helpIndex').controller('ticketViewCtrl', ['$scope','$http','$routeParams','TicketService','DepartmentsService', function($scope,$http,$routeParams,TicketService,DepartmentsService) {
	DepartmentsService.getDepartments().then(function(data) { $scope.departments = data.data; });
	TicketService.Get($routeParams.id).then(function(data) { $scope.ticket = data.data; });
}]);

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
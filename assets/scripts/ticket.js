helpdesk.controller('ticketViewCtrl', ['$scope','$http','$routeParams','TicketService', function($scope,$http,$routeParams,TicketService) {
	
	TicketService.Get($routeParams.id).then(function(data) {
		$scope.ticket = data.data;
	});
	console.log($routeParams);
}]);
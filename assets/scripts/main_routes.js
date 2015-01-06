helpdesk
.config(function($routeProvider) {
	$routeProvider
	.when('/', {
		templateUrl: '/templates/home.html',
		controller: 'homeCtrl'
	})
	.when('/tickets/mine', {
		templateUrl: '/templates/ticketlist.html',
		controller: 'myTicketListCtrl'
	})
	.when('/tickets/department', {
		templateUrl: '/templates/ticketList.html',
		controller: 'depTicketListCtrl'
	})
	.when('/tickets/assigned', {
		templateUrl: '/templates/ticketList.html',
		controller: 'assignedTicketListCtrl'
	})
	.when('/ticket/:id', {
		templateUrl: '/templates/ticket.html',
		controller: 'ticketViewCtrl'
	})
});